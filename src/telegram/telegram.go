package telegram

import (
	"context"
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"redminebot/src/config"
	"redminebot/src/task"
)

const botTimeout = 60

type Bot struct {
	Bot          *tgbot.BotAPI
	Log          *logrus.Entry
	AllowedUsers []string
}

func NewClient(log *logrus.Entry, cfg config.TelegramConfig) (*Bot, error) {
	bot, err := tgbot.NewBotAPI(cfg.Token.Get())
	if err != nil {
		return &Bot{}, fmt.Errorf("failed to create telegram-bot client: %v", err)
	} else {
		log.Debugf("connected to bot %v", bot.Self.UserName)
	}
	au, _ := config.SplitValue(cfg.AllowedUser)
	return &Bot{
		Bot:          bot,
		Log:          log,
		AllowedUsers: au,
	}, nil
}

func (b *Bot) Start(ctx context.Context, newTaskChan chan task.Task, taskStatusChan chan task.Status) {
	log := b.Log
	log.Info("start telegram-bot client")
	err := tgbot.SetLogger(log)
	if err != nil {
		log.Errorf("failed to set telegrma-bot client logger: %v", err)
	}
	u := tgbot.NewUpdate(0)
	u.Timeout = botTimeout
	updates := b.Bot.GetUpdatesChan(u)
	for {
		select {
		case <-ctx.Done():
			log.Infof("stop telegram-bot client")
			return
		case upd := <-updates:
			if upd.Message != nil {
				if checkUser(upd.Message.From.UserName, b.AllowedUsers) == true {
					log.Debugf("message from %v: %v", upd.Message.From.UserName, upd.Message.Text)
					newTaskChan <- task.Task{
						Message: upd.Message.Text,
						ChatId:  upd.Message.Chat.ID,
					}
				} else {
					log.Warnf("telegram user %v not allowed", upd.Message.From.UserName)
				}
			}
		case st := <-taskStatusChan:
			var statusMsg string
			if st.Error != nil {
				statusMsg = "error: failed to create redmine task"
			} else {
				statusMsg = fmt.Sprintf("New task _%v_ url: [%v](%v)", st.TaskId, st.TaskUrl.String(), st.TaskUrl.String())
			}
			msg := tgbot.NewMessage(st.ChatId, statusMsg)
			msg.ParseMode = "Markdown"
			_, err := b.Bot.Send(msg)
			if err != nil {
				log.Errorf("failed to send message to telegram-bot: %v", err)
			}
		}
	}
}

func checkUser(user string, allowedUser []string) bool {
	for _, u := range allowedUser {
		if u == user {
			return true
		}
	}
	return false
}
