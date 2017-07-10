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
	"path/filepath"

	"strings"

	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
)

var (
	supportedEventTypes = []string{
		"issue_comment",
		"pull_request",
		"pull_request_review",
		"pull_request_review_comment"}
)

type eventDataBuilder struct {
	secret []byte
	client *ghClient
	data   *eventData
}

func (builder *eventDataBuilder) validate(payload []byte, request *http.Request) bool {
	if len(builder.secret) == 0 {
		return true
	}
	h := hmac.New(sha1.New, builder.secret)
	h.Write(payload)
	result := fmt.Sprintf("sha1=%s", hex.EncodeToString(h.Sum(nil)))
	return request.Header.Get("X-Hub-Signature") == result
}

func (builder *eventDataBuilder) readPayload(r *http.Request) (*payload, []byte, error) {
	body := r.Body
	defer body.Close()
	raw, err := ioutil.ReadAll(io.LimitReader(body, r.ContentLength))
	if err != nil {
		return nil, raw, err
	}
	if !builder.validate(raw, r) {
		return nil, raw, fmt.Errorf("Payload could not be validated")
	}
	var pr payload
	if e := json.Unmarshal(raw, &pr); e != nil {
		return nil, raw, e
	}
	return &pr, raw, nil
}

func (builder *eventDataBuilder) readFromJson(payload *payload) {
	builder.data.number = payload.Number
	builder.data.title = payload.PullRequest.Title
	builder.data.changedFiles = payload.PullRequest.ChangedFiles
	builder.data.additions = payload.PullRequest.Additions
	builder.data.deletions = payload.PullRequest.Deletions
	builder.data.ref = payload.PullRequest.Base.Ref
	builder.data.origin = strings.ToLower(payload.PullRequest.User.Login)
	builder.data.state = payload.PullRequest.State
}

func (builder *eventDataBuilder) readFromClient() {
	id := builder.data.number
	builder.data.assignees = builder.client.GetAssignees(id)
	builder.data.state = builder.client.GetState(id)
	builder.data.labels = builder.client.GetLabels(id)
	builder.data.comments = builder.client.GetComments(id)
	fileNames := builder.client.GetFileNames(id)
	builder.data.fileNames = fileNames
	stringSet := util.StringSet{Transformer: filepath.Ext}
	builder.data.fileExt = stringSet.AddAll(fileNames).Values()
}

func (builder *eventDataBuilder) checkProcessState() bool {
	util.Logger.Debug("Current issue [(%d) %s] state is '%s'",
		builder.data.GetNumber(),
		builder.data.GetTitle(),
		builder.data.state)
	return builder.data.state != "closed"
}

func (builder *eventDataBuilder) PartialBuildFromRequest(config bot.ClientConfig, r *http.Request) (bot.EventData, bool, error) {
	githubEvent := r.Header.Get("X-Github-Event")
	if githubEvent == "ping" {
		util.Logger.Message("Got GitHub 'ping' event")
		return nil, false, nil
	}
	supportedEvent := false
	for _, event := range supportedEventTypes {
		if event == githubEvent {
			supportedEvent = true
		}
	}
	if !supportedEvent {
		util.Logger.Debug("Got GitHub '%s' event", githubEvent)
		return nil, false, nil
	}
	builder.secret = []byte(config.GetSecret())
	pl, raw, err := builder.readPayload(r)
	if err != nil {
		return nil, false, err
	}
	if pl.Number == 0 {
		util.Logger.Warning("Payload appear to have issue id 0")
		util.Logger.Debug("Faulty payload %+v", pl)
		return nil, false, fmt.Errorf("Payload appear to have issue id 0")
	}
	repo := pl.Repository.Name
	owner := pl.Repository.Owner.Login
	builder.data = &eventData{owner: owner, repo: repo, payload: raw}
	builder.readFromJson(pl)
	return builder.data, builder.checkProcessState(), nil
}

func (builder *eventDataBuilder) BuildFromRequest(config bot.ClientConfig, r *http.Request) (bot.EventData, bool, error) {
	_, ok, err := builder.PartialBuildFromRequest(config, r)
	if !ok || err != nil {
		return nil, ok, err
	}
	repo := builder.data.repo
	owner := builder.data.owner
	builder.client = newClient(config, owner, repo)
	builder.data.client = builder.client
	builder.readFromClient()
	return builder.data, builder.checkProcessState(), nil
}

func (builder *eventDataBuilder) BuildFromPayload(config bot.ClientConfig, raw []byte) (bot.EventData, bool, error) {
	var pl payload
	if e := json.Unmarshal(raw, &pl); e != nil {
		return nil, false, e
	}
	if pl.Number == 0 {
		util.Logger.Warning("Payload appear to have issue id 0")
		util.Logger.Debug("Faulty payload %+v", pl)
		return nil, false, fmt.Errorf("Payload appear to have issue id 0")
	}
	repo := pl.Repository.Name
	owner := pl.Repository.Owner.Login
	builder.client = newClient(config, owner, repo)
	builder.data = &eventData{owner: owner, repo: repo, payload: raw, client: builder.client}
	builder.readFromJson(&pl)
	builder.readFromClient()
	return builder.data, builder.checkProcessState(), nil
}

func init() {
	bot.RegisterNewBuilder("github", &eventDataBuilder{})
}
