package redmine

import (
	"context"
	"fmt"
	rd "github.com/nixys/nxs-go-redmine/v4"
	"github.com/sirupsen/logrus"
	"net/url"
	"redminebot/src/config"
	"redminebot/src/task"
	"strings"
	"unicode/utf8"
)

const taskUrlTpl = "%v/issues/%v"

type Task struct {
	Title       string
	Description string
}

type Redmine struct {
	Api            rd.Context
	Log            *logrus.Entry
	ProjectId      int
	TrackerId      int
	TaskStatusId   int
	SubjectLimiter int
	Host           string
}

func NewClient(log *logrus.Entry, cfg config.RedmineConfig) (Redmine, error) {
	var r rd.Context
	r.SetEndpoint(cfg.Host)
	r.SetAPIKey(cfg.ApiKey.Get())
	projectId, err := getProjectId(r, cfg.Project)
	if err != nil {
		return Redmine{}, fmt.Errorf("failed to get project id: %v", err)
	}
	trackerId, taskStatusId, err := getTrackerId(r, cfg.Tracker)
	if err != nil {
		return Redmine{}, fmt.Errorf("failed to get tracker id: %v", err)
	}
	return Redmine{
		Api:            r,
		Log:            log,
		ProjectId:      projectId,
		TrackerId:      trackerId,
		TaskStatusId:   taskStatusId,
		SubjectLimiter: cfg.SubjectLimiter,
		Host:           cfg.Host,
	}, nil
}

func (r Redmine) Start(ctx context.Context, newTaskChan chan task.Task, taskStatusChan chan task.Status) {
	log := r.Log
	r.Log.Infof("start redmine client")
	for {
		select {
		case <-ctx.Done():
			log.Infof("stop redmine client")
			return
		case t := <-newTaskChan:
			var ts task.Status
			ts.ChatId = t.ChatId
			newTask, err := r.createNewTask(t.Message)
			if err != nil {
				log.Errorf("failed to create new redmine task: %v", err)
				ts.Error = fmt.Errorf("failed to create new redmine task: %v", err)
			} else {
				log.Debugf("redmine new task id %v", newTask.ID)
				ts.TaskId = newTask.ID
				ts.TaskUrl, _ = url.Parse(fmt.Sprintf(taskUrlTpl, r.Host, newTask.ID))
			}
			taskStatusChan <- ts
		}
	}
}

func (r Redmine) createNewTask(tsk string) (rd.IssueObject, error) {
	subject, description := formatTask(tsk, r.SubjectLimiter)
	issue, _, err := r.Api.IssueCreate(rd.IssueCreateObject{
		ProjectID:   r.ProjectId,
		TrackerID:   r.TrackerId,
		Subject:     subject,
		Description: description,
		StatusID:    r.TaskStatusId,
	})
	if err != nil {
		return rd.IssueObject{}, fmt.Errorf("failed to create new task: %v", err)
	}
	return issue, nil
}

func getProjectId(r rd.Context, name string) (int, error) {
	projects, _, err := r.ProjectAllGet(rd.ProjectAllGetRequest{})
	if err != nil {
		return -1, fmt.Errorf("failed to get projects: %v", err)
	}
	for _, p := range projects.Projects {
		if strings.ToLower(p.Name) == strings.ToLower(name) {
			return p.ID, nil
		}
	}
	return -1, fmt.Errorf("unknown project name %v", name)
}

func getTrackerId(r rd.Context, name string) (trackerId int, taskStatusId int, err error) {
	trackers, _, err := r.TrackerAllGet()
	if err != nil {
		return -1, -1, fmt.Errorf("failed to get trackers: %v", err)
	}
	for _, t := range trackers {
		if strings.ToLower(t.Name) == strings.ToLower(name) {
			return t.ID, t.DefaultStatus.ID, nil
		}
	}
	return -1, -1, fmt.Errorf("unknown tracker name %v", name)
}

func formatTask(msg string, limiter int) (subject string, description string) {
	l := strings.Split(msg, "\n")
	s := l[0]
	if utf8.RuneCountInString(s) > limiter {
		subject = fmt.Sprintf("%v...", s[0:limiter])
	} else {
		subject = s
	}
	description = strings.Join(l, "\n")
	return subject, description
}
