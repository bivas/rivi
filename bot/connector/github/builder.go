package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
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

func (builder *eventDataBuilder) readPayload(r *http.Request) (*payload, error) {
	body := r.Body
	defer body.Close()
	raw, err := ioutil.ReadAll(io.LimitReader(body, r.ContentLength))
	if err != nil {
		return nil, err
	}
	if !builder.validate(raw, r) {
		return nil, fmt.Errorf("Payload could not be validated")
	}
	var pr payload
	if e := json.Unmarshal(raw, &pr); e != nil {
		return nil, e
	}
	return &pr, nil
}

func (builder *eventDataBuilder) readFromJson(payload *payload) {
	builder.data.number = payload.Number
	builder.data.title = payload.PullRequest.Title
	builder.data.changedFiles = payload.PullRequest.ChangedFiles
	builder.data.additions = payload.PullRequest.Additions
	builder.data.deletions = payload.PullRequest.Deletions
	builder.data.ref = payload.PullRequest.Base.Ref
	assignees := make([]string, 0)
	for _, assignee := range payload.PullRequest.Assignees {
		assignees = append(assignees, assignee.Login)
	}
	builder.data.assignees = assignees
	builder.data.origin = payload.PullRequest.User.Login
	builder.data.state = payload.PullRequest.State
}

func (builder *eventDataBuilder) readFromClient() {
	id := builder.data.number
	builder.data.state = builder.client.GetState(id)
	builder.data.labels = builder.client.GetLabels(id)
	builder.data.comments = builder.client.GetComments(id)
	fileNames := builder.client.GetFileNames(id)
	builder.data.fileNames = fileNames
	stringSet := util.StringSet{Transformer: filepath.Ext}
	builder.data.fileExt = stringSet.AddAll(fileNames).Values()
}

func (builder *eventDataBuilder) checkProcessState() bool {
	util.Logger.Debug("Current issue [(%d) %s] state is %s",
		builder.data.GetNumber(),
		builder.data.GetTitle(),
		builder.data.state)
	return builder.data.state != "closed"
}

func (builder *eventDataBuilder) BuildFromRequest(config bot.ClientConfig, r *http.Request) (bot.EventData, bool, error) {
	builder.secret = []byte(config.GetSecret())
	pl, err := builder.readPayload(r)
	if err != nil {
		return nil, false, err
	}
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
	repo := pl.Repository.Name
	owner := pl.Repository.Owner.Login
	builder.client = newClient(config, owner, repo)
	builder.data = &eventData{client: builder.client, owner: owner, repo: repo}
	builder.readFromJson(pl)
	builder.readFromClient()
	return builder.data, builder.checkProcessState(), nil
}

func (*eventDataBuilder) Build(config bot.ClientConfig, json string) (bot.EventData, error) {
	panic("implement me")
}

func init() {
	bot.RegisterNewBuilder("github", &eventDataBuilder{})
}
