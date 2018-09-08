package slack

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	svc    Service
}

func NewLoggingService(logger log.Logger, svc Service) Service {
	return &loggingService{logger, svc}
}

func (s *loggingService) Execute(ctx context.Context, c Command) (r Response, err error) {
	defer func(begin time.Time) {
		d := time.Since(begin)
		_ = s.logger.Log(
			"command_text", c.Text,
			"command_team_id", c.TeamID,
			"command_team_domain", c.TeamDomain,
			"command_user_name", c.UserName,
			"command_channel_name", c.ChannelName,
			"command_success", err == nil,
			"command_err", err,
			"command_took", d.String(),
			"command_sec", d.Seconds(),
		)
	}(time.Now())
	return s.svc.Execute(ctx, c)
}
