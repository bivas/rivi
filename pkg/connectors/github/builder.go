package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/bivas/rivi/pkg/config/client"
	"github.com/bivas/rivi/pkg/types"
	"github.com/bivas/rivi/pkg/types/builder"
	"github.com/bivas/rivi/pkg/util/log"
)

type builderContext struct {
	secret []byte
	client *Client
	data   *data
}

type dataBuilder struct {
	handlers       map[string]EventHandler
	defaultHandler EventHandler
	logger         log.Logger
}

func validate(secret, payload []byte, request *http.Request) bool {
	if len(secret) == 0 {
		return true
	}
	h := hmac.New(sha1.New, secret)
	h.Write(payload)
	result := fmt.Sprintf("sha1=%s", hex.EncodeToString(h.Sum(nil)))
	return request.Header.Get("X-Hub-Signature") == result
}

func ReadPayload(secret []byte, r *http.Request, v interface{}) ([]byte, error) {
	body := r.Body
	defer body.Close()
	raw, err := ioutil.ReadAll(io.LimitReader(body, r.ContentLength))
	if err != nil {
		return raw, err
	}
	if !validate(secret, raw, r) {
		return raw, errors.New("PullRequestPayload could not be validated")
	}
	if e := json.Unmarshal(raw, v); e != nil {
		return raw, e
	}
	return raw, nil
}

func (builder *dataBuilder) findEventHandler(githubEvent string) EventHandler {
	handler, ok := builder.handlers[githubEvent]
	if !ok {
		builder.logger.DebugWith(log.MetaFields{
			log.F("eventType", githubEvent),
		}, "Using default event handler")
		handler = builder.defaultHandler
	}
	return handler
}

func (builder *dataBuilder) BuildFromHook(config client.ClientConfig, r *http.Request) (types.HookData, bool, error) {
	githubEvent := r.Header.Get("X-Github-Event")
	return builder.findEventHandler(githubEvent).FromRequest(config, r)
}

func (builder *dataBuilder) BuildFromPayload(config client.ClientConfig, ofType string, raw []byte) (types.Data, bool, error) {
	return builder.findEventHandler(ofType).FromPayload(config, raw)
}

var DataBuilder dataBuilder

func init() {
	logger := log.Get("GitHub.DataBuilder")
	prHandler := &pullRequestEventHandler{
		logger: logger.Get("PullRequestHandler"),
	}
	DataBuilder = dataBuilder{
		logger: logger,
		handlers: map[string]EventHandler{
			"pull_request":                prHandler,
			"pull_request_review":         prHandler,
			"pull_request_review_comment": prHandler,
		},
		defaultHandler: DefaultEventHandler,
	}
	builder.RegisterNewDataBuilder("github", &DataBuilder)
}
