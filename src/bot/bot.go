package bot

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"redminebot/src/config"
	"redminebot/src/task"
	"redminebot/src/telegram"
)

type Bot interface {
	Start(ctx context.Context, newTaskChan chan task.Task, taskStatusChan chan task.Status)
}

func NewClient(log *logrus.Entry, cfg config.BotConfig, name BotName) (Bot, error) {
	switch name {
	case Telegram:
		return telegram.NewClient(log.WithField("client", Telegram.String()), cfg.Telegram)
	case Slack:
		return nil, fmt.Errorf("slack bot is TODO")
	}
	return nil, fmt.Errorf("unknown bot name %v", name)
}
