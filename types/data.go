package types

import (
	"net/http"
	"strings"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/util/log"
)

type Data interface {
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

type DataBuilder interface {
	BuildFromHook(config client.ClientConfig, r *http.Request) (Data, bool, error)
	BuildFromPayload(config client.ClientConfig, payload []byte) (Data, bool, error)
}

var builders map[string]DataBuilder = make(map[string]DataBuilder)

func RegisterNewDataBuilder(provider string, builder DataBuilder) {
	search := strings.ToLower(provider)
	_, exists := builders[search]
	if exists {
		log.Error("build for %s exists", provider)
	} else {
		log.Debug("registering builder %s", provider)
		builders[search] = builder
	}
}

func getBuilderFromRequest(r *http.Request) DataBuilder {
	var builder DataBuilder
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

func BuildFromHook(config client.ClientConfig, r *http.Request) (Data, bool) {
	builder := getBuilderFromRequest(r)
	if builder == nil {
		log.Error("No Builder to work with!")
		return nil, false
	}
	result, process, err := builder.BuildFromHook(config, r)
	if err != nil {
		log.ErrorWith(log.MetaFields{log.E(err)}, "Unable to build from request")
		return nil, false
	}
	return result, process
}

func BuildComplete(config client.ClientConfig, r *http.Request, data Data) (Data, bool) {
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
