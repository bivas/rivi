package bot

import (
	"github.com/bivas/rivi/util"
	"net/http"
	"strings"
)

type EventData interface {
	GetNumber() int
	GetTitle() string
	GetState() string
	GetOrigin() string
	GetOwner() string
	GetRepo() string
	GetLabels() []string
	HasLabel(label string) bool
	AddLabel(label string)
	RemoveLabel(label string)
	GetAssignees() []string
	HasAssignee(assignee string) bool
	AddAssignees(assignees... string)
	RemoveAssignees(assignees... string)
	GetComments() []string
	AddComment(comment string)
	GetFileNames() []string
	GetChangedFiles() int
	GetFileExtensions() []string
	GetChanges() (int, int)
}

type EventDataBuilder interface {
	BuildFromRequest(config ClientConfig, r *http.Request) (EventData, bool, error)
	Build(config ClientConfig, json string) (EventData, error)
}

var builders map[string]EventDataBuilder = make(map[string]EventDataBuilder)

func RegisterNewBuilder(provider string, builder EventDataBuilder) {
	_, exists := builders[provider]
	if exists {
		util.Logger.Error("build for %s exists!", provider)
	} else {
		util.Logger.Debug("registering builder %s", provider)
		builders[provider] = builder
	}
}

func buildFromRequest(config ClientConfig, r *http.Request) (EventData, bool) {
	var builder EventDataBuilder
	for name := range r.Header {
		for provider := range builders {
			if strings.Contains(strings.ToLower(name), provider) {
				builder = builders[provider]
				break
			}
		}
		if builder != nil {
			break
		}
	}
	if builder == nil {
		util.Logger.Error("No Builder to work with!")
		return nil, false
	}
	result, process, err := builder.BuildFromRequest(config, r)
	if err != nil {
		util.Logger.Error("Unable to build from request. %s", err)
		return nil, false
	}
	return result, process
}
