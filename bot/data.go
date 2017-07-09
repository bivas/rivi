package bot

import (
	"github.com/bivas/rivi/util"
	"net/http"
	"strings"
)

type Comment struct {
	Commenter string
	Comment   string
}

type EventData interface {
	GetNumber() int
	GetTitle() string
	GetState() string
	GetOrigin() string
	GetOwner() string
	GetRepo() string
	GetRef() string
	GetLabels() []string
	HasLabel(label string) bool
	AddLabel(label string)
	RemoveLabel(label string)
	GetAssignees() []string
	HasAssignee(assignee string) bool
	AddAssignees(assignees ...string)
	RemoveAssignees(assignees ...string)
	GetComments() []Comment
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
	search := strings.ToLower(provider)
	_, exists := builders[search]
	if exists {
		util.Logger.Error("connector for %s exists!", provider)
	} else {
		util.Logger.Debug("registering connector %s", provider)
		builders[search] = builder
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
