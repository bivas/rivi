package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util"
	"github.com/bivas/rivi/util/log"
)

type pullRequestEventHandler struct {
	logger log.Logger
}

func (h *pullRequestEventHandler) readFromJson(context *builderContext, payload *payload) {
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

func (h *pullRequestEventHandler) readFromClient(context *builderContext) {
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

func (h *pullRequestEventHandler) checkProcessState(context *builderContext) bool {
	h.logger.DebugWith(
		log.MetaFields{
			log.F("issue", context.data.GetShortName())},
		"Current state is '%s'", context.data.state)
	return context.data.state != "closed"
}

func (h *pullRequestEventHandler) FromRequest(config client.ClientConfig, r *http.Request) (types.HookData, bool, error) {
	context := &builderContext{secret: []byte(config.GetSecret())}
	pl, raw, err := ReadPayload(context.secret, r)
	if err != nil {
		return nil, false, err
	}
	if pl.Number == 0 && pl.PullRequest.Number == 0 {
		h.logger.Warning("Payload appear to have issue id 0")
		h.logger.Debug("Faulty payload %+v", pl)
		return nil, false, fmt.Errorf("Payload appear to have issue id 0")
	}
	repo := pl.Repository.Name
	owner := pl.Repository.Owner.Login
	installation := pl.Installation.ID
	if installation > 0 {
		context.client = NewAppClient(config, owner, repo, installation)
	} else {
		context.client = NewClient(config, owner, repo)
	}
	if context.client == nil {
		return nil, false, errors.New("Unable to initialize github client")
	}
	context.data = &data{
		owner: owner, repo: repo,
		payload: raw, eventType: r.Header.Get("X-Github-Event"),
		client: context.client}
	h.readFromJson(context, pl)
	return context.data, h.checkProcessState(context), nil
}

func (h *pullRequestEventHandler) FromPayload(config client.ClientConfig, raw []byte) (types.Data, bool, error) {
	var pl payload
	if e := json.Unmarshal(raw, &pl); e != nil {
		return nil, false, e
	}
	repo := pl.Repository.Name
	owner := pl.Repository.Owner.Login
	installation := pl.Installation.ID
	context := &builderContext{}
	if installation > 0 {
		context.client = NewAppClient(config, owner, repo, installation)
	} else {
		context.client = NewClient(config, owner, repo)
	}
	if context.client == nil {
		return nil, false, errors.New("Unable to initialize github client")
	}
	context.data = &data{owner: owner, repo: repo, payload: raw, client: context.client}
	h.readFromJson(context, &pl)
	if context.data.GetNumber() == 0 {
		h.logger.Warning("Payload appear to have issue id 0")
		h.logger.Debug("Faulty payload %+v", pl)
		return nil, false, fmt.Errorf("Payload appear to have issue id 0")
	}
	h.readFromClient(context)
	return context.data, h.checkProcessState(context), nil
}
