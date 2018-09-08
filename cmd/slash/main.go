package main

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/ags/guide-slack/pkg/guide"
	"github.com/ags/guide-slack/pkg/slack"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-kit/kit/log"
	rollbar "github.com/rollbar/rollbar-go"
)

func main() {
	var (
		apiKey        = os.Getenv("APP_GUIDE_API_KEY")
		signingSecret = os.Getenv("APP_SLACK_SIGNING_SECRET")
		rollbarToken  = os.Getenv("APP_ROLLBAR_TOKEN")
	)

	rollbar.SetToken(rollbarToken)
	rollbar.SetEnvironment("production")

	rand.Seed(time.Now().Unix())

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	client := guide.NewClient(apiKey)

	var svc slack.Service
	{
		svc = slack.NewService(client)
		svc = slack.NewLoggingService(logger, svc)
	}

	handler := slack.NewHandler([]byte(signingSecret), svc)

	lambda.Start(rollbarHandler(handler.Handle))
}

func rollbarHandler(h slack.HandlerFunc) slack.HandlerFunc {
	return func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		res, err := h(ctx, r)
		if err != nil {
			rollbar.Error(err)
			rollbar.Wait()
		}
		return res, err
	}
}
