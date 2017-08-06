package types

import (
	"net/http"
	"strings"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/util/log"
)

type EventData interface {
	GetShortName() string
	GetLongName() string
	GetNumber() int
	GetTitle() string
	GetDescription() string
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
	PartialBuildFromRequest(config client.ClientConfig, r *http.Request) (EventData, bool, error)
	BuildFromRequest(config client.ClientConfig, r *http.Request) (EventData, bool, error)
	BuildFromPayload(config client.ClientConfig, payload []byte) (EventData, bool, error)
}

var builders map[string]EventDataBuilder = make(map[string]EventDataBuilder)

func RegisterNewDataBuilder(provider string, builder EventDataBuilder) {
	search := strings.ToLower(provider)
	_, exists := builders[search]
	if exists {
		log.Error("build for %s exists", provider)
	} else {
		log.Debug("registering builder %s", provider)
		builders[search] = builder
	}
}

func getBuilderFromRequest(r *http.Request) EventDataBuilder {
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
	return builder
}

func BuildFromHook(config client.ClientConfig, r *http.Request) (EventData, bool) {
	builder := getBuilderFromRequest(r)
	if builder == nil {
		log.Error("No Builder to work with!")
		return nil, false
	}
	result, process, err := builder.PartialBuildFromRequest(config, r)
	if err != nil {
		log.ErrorWith(log.MetaFields{log.E(err)}, "Unable to build from request")
		return nil, false
	}
	return result, process
}

func BuildComplete(config client.ClientConfig, r *http.Request, data EventData) (EventData, bool) {
	builder := getBuilderFromRequest(r)
	if builder == nil {
		log.Error("No Builder to work with!")
		return nil, false
	}
	result, process, err := builder.BuildFromPayload(config, data.GetRawPayload())
	if err != nil {
		log.ErrorWith(log.MetaFields{log.E(err)}, "Unable to build from payload.")
		return nil, false
	}
	return result, process
}
