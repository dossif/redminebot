package bot

import (
	"fmt"
	"strings"
)

const (
	Telegram BotName = iota
	Slack
)

type BotName int

func (b BotName) String() string {
	switch b {
	case Telegram:
		return "telegram"
	case Slack:
		return "slack"
	default:
		return "unknown"
	}
}

func (b BotName) Set(name string) (BotName, error) {
	switch strings.ToLower(name) {
	case "telegram":
		return Telegram, nil
	case "slack":
		return Slack, nil
	default:
		return -1, fmt.Errorf("unknown bot name %v", name)
	}
}
