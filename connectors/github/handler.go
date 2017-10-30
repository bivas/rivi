package github

import (
	"net/http"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
)

type eventHandler interface {
	FromRequest(client.ClientConfig, *http.Request) (types.HookData, bool, error)
	FromPayload(client.ClientConfig, []byte) (types.Data, bool, error)
}

type defaultEventHandler struct {
	logger log.Logger
}

func (h *defaultEventHandler) FromRequest(config client.ClientConfig, r *http.Request) (types.HookData, bool, error) {
	githubEvent := r.Header.Get("X-Github-Event")
	h.logger.Info("Got GitHub '%s' event", githubEvent)
	return nil, false, nil
}

func (h *defaultEventHandler) FromPayload(client.ClientConfig, []byte) (types.Data, bool, error) {
	h.logger.Warning("Calling 'FromPayload' of default handler")
	return nil, false, nil
}

var defaultHandler = &defaultEventHandler{log.Get("GitHub.DataBuilder.DefaultHandler")}
