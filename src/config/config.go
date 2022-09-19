package config

import (
	"encoding/json"
	"fmt"
	"github.com/cristalhq/aconfig"
	"strings"
)

const (
	configPrefix  = "app"
	listDelimiter = ","
)

type MainConfig struct {
	LogLevel string `default:"error" usage:"log level: debug|info|warn|error|fatal"`
	Bot      BotConfig
	Redmine  RedmineConfig
}

type BotConfig struct {
	Name     string `default:"telegram" usage:"witch bot to use"`
	Telegram TelegramConfig
	Slack    SlackBot
}

type RedmineConfig struct {
	ApiHost        string       `default:"http://redmine:3000" usage:"redmine api host"`
	WebHost        string       `default:"https://redmine.example.com" usage:"redmine web url"`
	ApiKey         SecretString `default:"<NOT_SET>" usage:"redmine user rest-api token"`
	Project        string       `default:"Misc" usage:"redmine project name"`
	Tracker        string       `default:"Task" usage:"redmine tracker name"`
	SubjectLimiter int          `default:"20" usage:"redmine task subject length limiter"`
}

type TelegramConfig struct {
	Token       SecretString `default:"<NOT_SET>" usage:"telegram bot token"`
	AllowedUser string       `default:"admin" usage:"comma-separated list of allowed telegram users"`
}

type SlackBot struct {
	//TODO: implement slack bot
}

func NewConfig() (cfg MainConfig, err error) {
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipFiles:        true,
		EnvPrefix:        strings.ToUpper(configPrefix),
		FlagPrefix:       strings.ToLower(configPrefix),
		AllFieldRequired: true,
	})
	err = loader.Load()
	if err != nil {
		return cfg, fmt.Errorf("failed to load config from source: %v", err)
	}
	return cfg, err
}

func SplitValue(val string) ([]string, error) {
	sl := strings.Split(val, listDelimiter)
	if len(sl) == 0 {
		return sl, fmt.Errorf("empty list")
	}
	return sl, nil
}

func GetTextConfig(config interface{}) string {
	js, _ := json.MarshalIndent(config, "", "  ")
	return string(js)
}
