package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/types/builder"
	"github.com/bivas/rivi/util/log"
)

type builderContext struct {
	secret []byte
	client *ghClient
	data   *data
}

type dataBuilder struct {
	handlers       map[string]eventHandler
	defaultHandler eventHandler
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

func ReadPayload(secret []byte, r *http.Request) (*payload, []byte, error) {
	body := r.Body
	defer body.Close()
	raw, err := ioutil.ReadAll(io.LimitReader(body, r.ContentLength))
	if err != nil {
		return nil, raw, err
	}
	if !validate(secret, raw, r) {
		return nil, raw, fmt.Errorf("Payload could not be validated")
	}
	var pr payload
	if e := json.Unmarshal(raw, &pr); e != nil {
		return nil, raw, e
	}
	return &pr, raw, nil
}

func (builder *dataBuilder) findEventHandler(githubEvent string) eventHandler {
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
		handlers: map[string]eventHandler{
			"pull_request":                prHandler,
			"pull_request_review":         prHandler,
			"pull_request_review_comment": prHandler,
		},
		defaultHandler: defaultHandler,
	}
	builder.RegisterNewDataBuilder("github", &DataBuilder)
}
