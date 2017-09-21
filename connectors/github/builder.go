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

	"errors"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/types/builder"
	"github.com/bivas/rivi/util"
	"github.com/bivas/rivi/util/log"
)

var (
	supportedEventTypes = []string{
		"issue_comment",
		"pull_request",
		"pull_request_review",
		"pull_request_review_comment"}
)

type builderContext struct {
	secret []byte
	client *ghClient
	data   *data
}

type dataBuilder struct {
	logger log.Logger
}

func (builder *dataBuilder) validate(context *builderContext, payload []byte, request *http.Request) bool {
	if len(context.secret) == 0 {
		return true
	}
	h := hmac.New(sha1.New, context.secret)
	h.Write(payload)
	result := fmt.Sprintf("sha1=%s", hex.EncodeToString(h.Sum(nil)))
	return request.Header.Get("X-Hub-Signature") == result
}

func (builder *dataBuilder) readPayload(context *builderContext, r *http.Request) (*payload, []byte, error) {
	body := r.Body
	defer body.Close()
	raw, err := ioutil.ReadAll(io.LimitReader(body, r.ContentLength))
	if err != nil {
		return nil, raw, err
	}
	if !builder.validate(context, raw, r) {
		return nil, raw, fmt.Errorf("Payload could not be validated")
	}
	var pr payload
	if e := json.Unmarshal(raw, &pr); e != nil {
		return nil, raw, e
	}
	return &pr, raw, nil
}

func (builder *dataBuilder) readFromJson(context *builderContext, payload *payload) {
	pr := payload.PullRequest
	if pr.Number > 0 {
		context.data.number = pr.Number
	} else {
		context.data.number = payload.Number
	}
	context.data.title = pr.Title
	context.data.description = pr.Body
	context.data.changedFiles = pr.ChangedFiles
	context.data.additions = pr.Additions
	context.data.deletions = pr.Deletions
	context.data.ref = pr.Base.Ref
	head := pr.Head
	context.data.origin = types.Origin{
		User:   strings.ToLower(head.User.Login),
		Repo:   head.Repo.Name,
		Ref:    head.Ref,
		Head:   head.Sha[0:6],
		GitURL: head.Repo.GitURL,
	}
	context.data.state = pr.State
}

func (builder *dataBuilder) readFromClient(context *builderContext) {
	id := context.data.number
	context.data.assignees = context.client.GetAssignees(id)
	context.data.state = context.client.GetState(id)
	context.data.labels = context.client.GetLabels(id)
	context.data.comments = context.client.GetComments(id)
	fileNames := context.client.GetFileNames(id)
	context.data.fileNames = fileNames
	stringSet := util.StringSet{Transformer: filepath.Ext}
	context.data.fileExt = stringSet.AddAll(fileNames).Values()
	context.data.reviewers = context.client.GetReviewers(id)
	context.data.locked = context.client.Locked(id)
}

func (builder *dataBuilder) checkProcessState(context *builderContext) bool {
	builder.logger.DebugWith(log.MetaFields{log.F("issue", context.data.GetShortName())},
		"Current state is '%s'", context.data.state)
	return context.data.state != "closed"
}

func (builder *dataBuilder) BuildFromHook(config client.ClientConfig, r *http.Request) (types.HookData, bool, error) {
	githubEvent := r.Header.Get("X-Github-Event")
	if githubEvent == "ping" {
		builder.logger.Info("Got GitHub 'ping' event")
		return nil, false, nil
	}
	supportedEvent := false
	for _, event := range supportedEventTypes {
		if event == githubEvent {
			supportedEvent = true
		}
	}
	if !supportedEvent {
		builder.logger.Debug("Got GitHub '%s' event", githubEvent)
		return nil, false, nil
	}
	context := &builderContext{secret: []byte(config.GetSecret())}
	pl, raw, err := builder.readPayload(context, r)
	if err != nil {
		return nil, false, err
	}
	if pl.Number == 0 && pl.PullRequest.Number == 0 {
		builder.logger.Warning("Payload appear to have issue id 0")
		builder.logger.Debug("Faulty payload %+v", pl)
		return nil, false, fmt.Errorf("Payload appear to have issue id 0")
	}
	repo := pl.Repository.Name
	owner := pl.Repository.Owner.Login
	installation := pl.Installation.ID
	if installation > 0 {
		context.client = newAppClient(config, owner, repo, installation)
	} else {
		context.client = newClient(config, owner, repo)
	}
	if context.client == nil {
		return nil, false, errors.New("Unable to initialize github client")
	}
	context.data = &data{owner: owner, repo: repo, payload: raw, client: context.client}
	builder.readFromJson(context, pl)
	return context.data, builder.checkProcessState(context), nil
}

func (builder *dataBuilder) BuildFromPayload(config client.ClientConfig, raw []byte) (types.Data, bool, error) {
	var pl payload
	if e := json.Unmarshal(raw, &pl); e != nil {
		return nil, false, e
	}
	repo := pl.Repository.Name
	owner := pl.Repository.Owner.Login
	installation := pl.Installation.ID
	context := &builderContext{}
	if installation > 0 {
		context.client = newAppClient(config, owner, repo, installation)
	} else {
		context.client = newClient(config, owner, repo)
	}
	if context.client == nil {
		return nil, false, errors.New("Unable to initialize github client")
	}
	context.data = &data{owner: owner, repo: repo, payload: raw, client: context.client}
	builder.readFromJson(context, &pl)
	if context.data.GetNumber() == 0 {
		builder.logger.Warning("Payload appear to have issue id 0")
		builder.logger.Debug("Faulty payload %+v", pl)
		return nil, false, fmt.Errorf("Payload appear to have issue id 0")
	}
	builder.readFromClient(context)
	return context.data, builder.checkProcessState(context), nil
}

func init() {
	builder.RegisterNewDataBuilder("github", &dataBuilder{logger: log.Get("GitHub.DataBuilder")})
}
