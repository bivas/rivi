package bot

import (
	"net/http"
	"strings"

	"github.com/bivas/rivi/util"
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
	GetRawPayload() []byte
}

type EventDataBuilder interface {
	PartialBuildFromRequest(config ClientConfig, r *http.Request) (EventData, bool, error)
	BuildFromRequest(config ClientConfig, r *http.Request) (EventData, bool, error)
	BuildFromPayload(config ClientConfig, payload []byte) (EventData, bool, error)
}

var builders map[string]EventDataBuilder = make(map[string]EventDataBuilder)

func RegisterNewBuilder(provider string, builder EventDataBuilder) {
	search := strings.ToLower(provider)
	_, exists := builders[search]
	if exists {
		util.Logger.Error("build for %s exists!", provider)
	} else {
		util.Logger.Debug("registering builder %s", provider)
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
	result, process, err := builder.PartialBuildFromRequest(config, r)
	if err != nil {
		util.Logger.Error("Unable to build from request. %s", err)
		return nil, false
	}
	return result, process
}

func completeBuildFromRequest(config ClientConfig, r *http.Request) (EventData, bool) {
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

func completeBuild(config ClientConfig, r *http.Request, data EventData) (EventData, bool) {
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
	result, process, err := builder.BuildFromPayload(config, data.GetRawPayload())
	if err != nil {
		util.Logger.Error("Unable to build from payload. %s", err)
		return nil, false
	}
	return result, process
}
