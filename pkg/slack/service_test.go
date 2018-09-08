package slack_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/ags/guide-slack/pkg/guide"
	"github.com/ags/guide-slack/pkg/slack"
)

func TestService_Sample(t *testing.T) {
	// TODO this is for debugging, not a proper test.
	c := guide.NewClient(os.Getenv("APP_GUIDE_API_KEY"))
	s := slack.NewService(c)

	res, err := s.Execute(context.Background(), slack.Command{
		Text: "coffee",
	})

	if err != nil {
		t.Fatal(err)
	}
	json.NewEncoder(os.Stdout).Encode(res)
}
