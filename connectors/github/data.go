package github

import (
	"fmt"

	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util"
)

type data struct {
	client        *ghClient
	number        int
	state         string
	locked        bool
	origin        types.Origin
	owner         string
	repo          string
	ref           string
	title         string
	description   string
	changedFiles  int
	fileNames     []string
	fileExt       []string
	additions     int
	deletions     int
	labels        []string
	assignees     []string
	comments      []types.Comment
	payload       []byte
	reviewers     map[string]string
	collaborators []string
	repoLabels    []string
}

func (d *data) GetAvailableLabels() []string {
	if len(d.repoLabels) == 0 {
		d.repoLabels = d.client.GetAvailableLabels()
	}
	return d.repoLabels
}

func (d *data) GetCollaborators() []string {
	if len(d.collaborators) == 0 {
		d.collaborators = d.client.GetCollaborators()
	}
	return d.collaborators
}

func (d *data) IsCollaborator(name string) bool {
	for _, collaborator := range d.GetCollaborators() {
		if name == collaborator {
			return true
		}
	}
	return false
}

func (d *data) GetRulesFile() string {
	hasUpdatedRules := false
	for _, file := range d.fileNames {
		if file == types.RulesConfigFileName {
			hasUpdatedRules = true
			break
		}
	}
	if hasUpdatedRules && d.IsCollaborator(d.origin.User) {
		return d.client.GetFileContentFromRef(
			types.RulesConfigFileName,
			d.origin.User,
			d.origin.Repo,
			d.origin.Head)
	}
	return d.client.GetFileContent(types.RulesConfigFileName)
}

func (d *data) GetRepository() types.Repository {
	return d
}

func (d *data) GetProvider() string {
	return "github"
}

func (d *data) GetShortName() string {
	return fmt.Sprintf("%s/%s#%d", d.owner, d.repo, d.number)
}

func (d *data) GetLongName() string {
	return fmt.Sprintf("%s/%s#%d [%s]", d.owner, d.repo, d.number, d.title)
}

func (d *data) GetReviewers() map[string]string {
	return d.reviewers
}

func (d *data) GetApprovals() []string {
	result := util.StringSet{}
	for reviewer, state := range d.reviewers {
		if state == "approved" {
			result.Add(reviewer)
		}
	}
	return result.Values()
}

func (d *data) Lock() {
	d.client.Lock(d.number)
	d.locked = true
}

func (d *data) Unlock() {
	d.client.Unlock(d.number)
	d.locked = false
}

func (d *data) LockState() bool {
	return d.locked
}

func (d *data) GetRawPayload() []byte {
	return d.payload
}

func (d *data) Merge(mergeMethod string) {
	d.client.Merge(d.number, mergeMethod)
}

func (d *data) GetState() string {
	return d.state
}

func (d *data) AddLabel(label string) {
	d.labels = d.client.AddLabel(d.number, label)
}

func (d *data) RemoveLabel(label string) {
	d.labels = d.client.RemoveLabel(d.number, label)
}

func (d *data) GetLabels() []string {
	return d.labels
}

func (d *data) HasLabel(label string) bool {
	for _, name := range d.labels {
		if name == label {
			return true
		}
	}
	return false
}

func (d *data) AddAssignees(assignees ...string) {
	d.assignees = d.client.AddAssignees(d.number, assignees...)
}

func (d *data) RemoveAssignees(assignees ...string) {
	d.assignees = d.client.RemoveAssignees(d.number, assignees...)
}

func (d *data) GetAssignees() []string {
	return d.assignees
}

func (d *data) HasAssignee(assignee string) bool {
	for _, name := range d.assignees {
		if name == assignee {
			return true
		}
	}
	return false
}

func (d *data) GetComments() []types.Comment {
	return d.comments
}

func (d *data) AddComment(comment string) {
	d.comments = append(d.comments, d.client.AddComment(d.number, comment))
}

func (d *data) GetNumber() int {
	return d.number
}

func (d *data) GetTitle() string {
	return d.title
}

func (d *data) GetDescription() string {
	return d.description
}

func (d *data) GetOrigin() types.Origin {
	return d.origin
}

func (d *data) GetOwner() string {
	return d.owner
}

func (d *data) GetRepo() string {
	return d.repo
}

func (d *data) GetRef() string {
	return d.ref
}

func (d *data) GetFileNames() []string {
	return d.fileNames
}

func (d *data) GetChangedFiles() int {
	return d.changedFiles
}

func (d *data) GetFileExtensions() []string {
	return d.fileExt
}

func (d *data) GetChanges() (int, int) {
	return d.additions, d.deletions
}
