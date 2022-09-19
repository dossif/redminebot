package slack

import (
	"context"
	"github.com/sirupsen/logrus"
	"redminebot/src/config"
	"sync"
)

type Bot struct {
}

func NewBot(log *logrus.Entry, cfg config.TelegramConfig) (Bot, error) {
	return Bot{}, nil
}

func (b *Bot) Start(ctx context.Context, wg *sync.WaitGroup) {
	// TODO: implement slack bot
}
