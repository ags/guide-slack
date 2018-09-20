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

const (
	ResponseTypeEphemeral = "ephemeral"
)

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
	cmd := resolveCommand(c.Text)

	switch cmd {
	case commandUnknown:
		return responseNoCommand, nil
	case commandHelp:
		return s.help(ctx, c)
	case commandList:
		return s.list(ctx, c)
	case commandRandom:
		return s.random(ctx, c)
	default:
		return Response{}, fmt.Errorf("unhandled command '%v'", cmd)
	}
}

func (s *service) help(ctx context.Context, c Command) (Response, error) {
	return Response{
		ResponseType: ResponseTypeEphemeral,
		Text:         "sorry buddy, no help yet",
		Attachments: []Attachment{
			Attachment{
				Fallback: "sorry buddy",
				ImageURL: "https://i.kym-cdn.com/photos/images/newsfeed/000/676/045/ffc.gif",
			},
		},
	}, nil
}

func (s *service) list(ctx context.Context, c Command) (Response, error) {
	team, err := s.findTeam(c.TeamID)
	if err != nil {
		return Response{}, err
	}

	hunt := textToHunt(c.Text)

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

	destinations := sample(res.Destinations, 10)

	lines := []string{
		"Our curated list:",
	}

	for _, d := range destinations {
		parts := []string{d.Name}

		addPart := func(str string) {
			s := strings.Trim(str, " ")
			if s != "" {
				parts = append(parts, s)
			}
		}

		addPart(d.Street)
		addPart(d.Suburb)
		addPart(d.Website)

		lines = append(lines, "• "+strings.Join(parts, ", "))
	}

	return Response{
		ResponseType: ResponseTypeEphemeral,
		Text:         strings.Join(lines, "\n"),
	}, nil
}

func (s *service) random(ctx context.Context, c Command) (Response, error) {
	team, err := s.findTeam(c.TeamID)
	if err != nil {
		return Response{}, err
	}

	// TODO consider making this part of the comamnd itself
	hunt := textToHunt(strings.TrimLeft(c.Text, "random "))

	res, err := s.client.FindHunt(ctx, team.key, team.region, hunt)
	if err != nil {
		return Response{}, err
	}

	if len(res.Destinations) == 0 {
		return notFoundResponse(c.Text), nil
	}

	d := sample(res.Destinations, 1)[0]

	a := Attachment{
		Fallback:  d.Name,
		Title:     d.Name,
		TitleLink: d.Website,
		Text:      d.Description,
		Fields: []Field{
			Field{
				Title: "Address",
				Value: fmt.Sprintf("%s, %s", d.Street, d.Suburb), // TODO handle missing part
				Short: false,
			},
		},
	}
	if len(d.BannerImages) > 0 {
		a.ImageURL = "https://guide.app" + d.BannerImages[0]
	}

	return Response{
		ResponseType: ResponseTypeEphemeral,
		Text:         "What about...",
		Attachments:  []Attachment{a},
	}, nil
}

// TODO actually map this
func (s *service) findTeam(teamID string) (team, error) {
	return bluechilli, nil
}

const (
	commandHelp    = "help"
	commandList    = "list"
	commandRandom  = "random"
	commandUnknown = "unknown"
)

func resolveCommand(text string) string {
	switch {
	case text == "":
		return commandUnknown
	case text == "help":
		return commandHelp
	case strings.HasPrefix(text, "random "):
		return commandRandom
	default:
		return commandList
	}
}

func textToHunt(text string) string {
	hunt := text
	hunt = strings.ToLower(hunt)
	hunt = strings.Replace(hunt, " ", "-", -1)
	return hunt
}

func notFoundResponse(text string) Response {
	return Response{
		ResponseType: ResponseTypeEphemeral,
		Text:         fmt.Sprintf("Sorry, could't find anything for '%s' :disappointed:", text),
	}
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
