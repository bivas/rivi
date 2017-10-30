package builder

import (
	"net/http"
	"strings"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
)

var l log.Logger = log.Get("data.builder")

type DataBuilder interface {
	BuildFromHook(config client.ClientConfig, r *http.Request) (types.HookData, bool, error)
	BuildFromPayload(config client.ClientConfig, ofType string, payload []byte) (types.Data, bool, error)
}

var builders map[string]DataBuilder = make(map[string]DataBuilder)

func RegisterNewDataBuilder(provider string, builder DataBuilder) {
	search := strings.ToLower(provider)
	_, exists := builders[search]
	if exists {
		l.Error("build for %s exists", provider)
	} else {
		l.Debug("registering builder %s", provider)
		builders[search] = builder
	}
}

func getBuilderFromRequest(r *http.Request) DataBuilder {
	userAgent := strings.ToLower(r.UserAgent())
	for provider := range builders {
		if strings.Contains(userAgent, provider) {
			return builders[provider]
		}
	}
	return nil
}

func BuildFromHook(config client.ClientConfig, r *http.Request) (types.HookData, bool) {
	builder := getBuilderFromRequest(r)
	if builder == nil {
		l.Error("No Builder to work with!")
		return nil, false
	}
	result, process, err := builder.BuildFromHook(config, r)
	if err != nil {
		l.ErrorWith(log.MetaFields{log.E(err)}, "Unable to build from request")
		return nil, false
	}
	return result, process
}

func BuildComplete(config client.ClientConfig, data types.ReadOnlyData) (types.Data, bool) {
	builder, exists := builders[data.GetProvider()]
	if !exists {
		l.Error("No existing builder to work with!")
		return nil, false
	}
	result, process, err := builder.BuildFromPayload(config, data.GetRawType(), data.GetRawPayload())
	if err != nil {
		l.ErrorWith(log.MetaFields{log.E(err)}, "Unable to build from payload.")
		return nil, false
	}
	return result, process
}
