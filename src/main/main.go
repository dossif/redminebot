package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zs5460/art"
	"redminebot/src/bot"
	"redminebot/src/config"
	"redminebot/src/logger"
	"redminebot/src/redmine"
	"redminebot/src/signal"
	"redminebot/src/task"
	"strings"
	"sync"
)

const (
	appName = "RedmineBot"
)

var (
	AppVersion = "0.0.0"
)

func main() {
	fmt.Println(art.String(appName))
	ctx := signal.ContextWithSignal(context.Background())
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to create config: %v", err))
	}
	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}
	lg := log.WithFields(logrus.Fields{"app": strings.ToLower(appName), "ver": AppVersion})
	lg.Infof("start %v ver %v", appName, AppVersion)
	defer lg.Infof("stop %v", appName)
	if strings.ToLower(cfg.LogLevel) == strings.ToLower(logrus.DebugLevel.String()) {
		lg.Debugf("debug config: %v", config.GetTextConfig(cfg))
	}
	newTaskChan := make(chan task.Task)
	taskStatusChan := make(chan task.Status)
	// bot
	var botName bot.BotName
	bn, err := botName.Set(cfg.Bot.Name)
	if err != nil {
		lg.Fatalf("failed to get bot name: %v", err)
	}
	var wg sync.WaitGroup
	b, err := bot.NewClient(lg, cfg.Bot, bn)
	if err != nil {
		lg.Fatalf("failed to create bot client: %v", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		b.Start(ctx, newTaskChan, taskStatusChan)
	}()
	// redmine
	r, err := redmine.NewClient(lg.WithField("client", "redmine"), cfg.Redmine)
	if err != nil {
		lg.Fatalf("failed to create redmine client: %v", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		r.Start(ctx, newTaskChan, taskStatusChan)
	}()
	wg.Wait()
}
