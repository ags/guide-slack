package slack

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/ags/guide-slack/pkg/guide"
)

type Service interface {
	Execute(ctx context.Context, c Command) (Response, error)
}

type Command struct {
	Text        string
	TeamID      string
	TeamDomain  string
	UserName    string
	ChannelName string
}

type Response struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

type Attachment struct {
	Fallback  string  `json:"fallback"`
	Title     string  `json:"title"`
	TitleLink string  `json:"title_link"`
	Text      string  `json:"text"`
	ImageURL  string  `json:"image_url"`
	Fields    []Field `json:"fields"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type team struct {
	region string
	key    string
}

type service struct {
	client *guide.Client
}

var (
	responseNoCommand = Response{
		Text: "Need some guidance? Try `/guide help`!",
	}

	bluechilli = team{
		key:    "7e8ff9e9-a1d9-4503-bf13-05a04603e4fe",
		region: "sydney",
	}
)

func NewService(client *guide.Client) Service {
	return &service{client}
}

func (s *service) Execute(ctx context.Context, c Command) (Response, error) {
	switch c.Text {
	case "":
		return responseNoCommand, nil
	case "help":
		return s.help(ctx, c)
	default:
		return s.sample(ctx, c)
	}
}

func (s *service) help(ctx context.Context, c Command) (Response, error) {
	return Response{
		ResponseType: "ephemeral",
		Text:         "sorry buddy, no help yet",
		Attachments: []Attachment{
			Attachment{
				Fallback: "sorry buddy",
				ImageURL: "https://i.kym-cdn.com/photos/images/newsfeed/000/676/045/ffc.gif",
			},
		},
	}, nil
}

func (s *service) sample(ctx context.Context, c Command) (Response, error) {
	team, err := s.findTeam(c.TeamID)
	if err != nil {
		return Response{}, err
	}

	hunt := c.Text
	hunt = strings.ToLower(hunt)
	hunt = strings.Replace(hunt, " ", "-", -1)

	res, err := s.client.FindHunt(ctx, team.key, team.region, hunt)
	if err != nil {
		return Response{}, err
	}

	if len(res.Destinations) == 0 {
		return Response{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Sorry, could't find anything for '%s' :disappointed:", c.Text),
		}, nil
	}

	destinations := sample(res.Destinations, 2)
	attachments := make([]Attachment, len(destinations))

	for i := 0; i < len(destinations); i++ {
		d := destinations[i]

		a := Attachment{
			Fallback:  d.Name,
			Title:     d.Name,
			TitleLink: d.Website,
			Text:      d.Description,
			Fields: []Field{
				Field{
					Title: "Address",
					Value: fmt.Sprintf("%s, %s", d.Street, d.Suburb),
					Short: false,
				},
			},
		}
		if len(d.BannerImages) > 0 {
			a.ImageURL = "https://guide.app" + d.BannerImages[0]
		}

		attachments[i] = a
	}

	return Response{
		ResponseType: "in_channel",
		Text:         "Check these out:",
		Attachments:  attachments,
	}, nil
}

// TODO actually map this
func (s *service) findTeam(teamID string) (team, error) {
	return bluechilli, nil
}

func sample(ds []guide.Destination, n int) []guide.Destination {
	n = min(len(ds), n)
	rand.Shuffle(len(ds), func(i, j int) {
		ds[i], ds[j] = ds[j], ds[i]
	})
	return ds[:n]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
