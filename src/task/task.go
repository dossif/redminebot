package task

import "net/url"

type Task struct {
	Message string
	ChatId  int64
}

type Status struct {
	TaskId  int
	TaskUrl *url.URL
	Error   error
	ChatId  int64
}
