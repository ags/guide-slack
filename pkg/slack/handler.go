package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

type HandlerFunc func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Handler struct {
	signingSecret []byte
	svc           Service
}

func NewHandler(signingSecret []byte, svc Service) *Handler {
	return &Handler{signingSecret, svc}
}

func (h *Handler) Handle(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !h.verifySignature(r) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}, nil
	}

	cmd, err := decodeCommand(r)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	res, err := h.svc.Execute(ctx, cmd)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	body, err := json.Marshal(res)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}

func (h *Handler) verifySignature(r events.APIGatewayProxyRequest) bool {
	signature, _ := r.Headers["X-Slack-Signature"]
	timestamp, _ := r.Headers["X-Slack-Request-Timestamp"]
	return VerifySignature(h.signingSecret, timestamp, r.Body, signature)
}

func decodeCommand(r events.APIGatewayProxyRequest) (Command, error) {
	q, err := url.ParseQuery(r.Body)
	if err != nil {
		return Command{}, err
	}

	return Command{
		Text:        q.Get("text"),
		TeamID:      q.Get("team_id"),
		TeamDomain:  q.Get("team_domain"),
		UserName:    q.Get("user_name"),
		ChannelName: q.Get("channel_name"),
	}, nil
}
