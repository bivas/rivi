package github

import (
	"fmt"

	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
)

type eventData struct {
	client       *ghClient
	number       int
	state        string
	locked       bool
	origin       string
	owner        string
	repo         string
	ref          string
	title        string
	description  string
	changedFiles int
	fileNames    []string
	fileExt      []string
	additions    int
	deletions    int
	labels       []string
	assignees    []string
	comments     []bot.Comment
	payload      []byte
	reviewers    map[string]string
}

func (d *eventData) GetShortName() string {
	return fmt.Sprintf("%s/%s#%d", d.owner, d.repo, d.number)
}

func (d *eventData) GetLongName() string {
	return fmt.Sprintf("%s/%s#%d [%s]", d.owner, d.repo, d.number, d.title)
}

func (d *eventData) GetReviewers() map[string]string {
	return d.reviewers
}

func (d *eventData) GetApprovals() []string {
	result := util.StringSet{}
	for reviewer, state := range d.reviewers {
		if state == "approved" {
			result.Add(reviewer)
		}
	}
	return result.Values()
}

func (d *eventData) Lock() {
	d.client.Lock(d.number)
	d.locked = true
}

func (d *eventData) Unlock() {
	d.client.Unlock(d.number)
	d.locked = false
}

func (d *eventData) LockState() bool {
	return d.locked
}

func (d *eventData) GetRawPayload() []byte {
	return d.payload
}

func (d *eventData) Merge(mergeMethod string) {
	d.client.Merge(d.number, mergeMethod)
}

func (d *eventData) GetState() string {
	return d.state
}

func (d *eventData) AddLabel(label string) {
	d.labels = d.client.AddLabel(d.number, label)
}

func (d *eventData) RemoveLabel(label string) {
	d.labels = d.client.RemoveLabel(d.number, label)
}

func (d *eventData) GetLabels() []string {
	return d.labels
}

func (d *eventData) HasLabel(label string) bool {
	for _, name := range d.labels {
		if name == label {
			return true
		}
	}
	return false
}

func (d *eventData) AddAssignees(assignees ...string) {
	d.assignees = d.client.AddAssignees(d.number, assignees...)
}

func (d *eventData) RemoveAssignees(assignees ...string) {
	d.assignees = d.client.RemoveAssignees(d.number, assignees...)
}

func (d *eventData) GetAssignees() []string {
	return d.assignees
}

func (d *eventData) HasAssignee(assignee string) bool {
	for _, name := range d.assignees {
		if name == assignee {
			return true
		}
	}
	return false
}

func (d *eventData) GetComments() []bot.Comment {
	return d.comments
}

func (d *eventData) AddComment(comment string) {
	d.comments = append(d.comments, d.client.AddComment(d.number, comment))
}

func (d *eventData) GetNumber() int {
	return d.number
}

func (d *eventData) GetTitle() string {
	return d.title
}

func (d *eventData) GetDescription() string {
	return d.description
}

func (d *eventData) GetOrigin() string {
	return d.origin
}

func (d *eventData) GetOwner() string {
	return d.owner
}

func (d *eventData) GetRepo() string {
	return d.repo
}

func (d *eventData) GetRef() string {
	return d.ref
}

func (d *eventData) GetFileNames() []string {
	return d.fileNames
}

func (d *eventData) GetChangedFiles() int {
	return d.changedFiles
}

func (d *eventData) GetFileExtensions() []string {
	return d.fileExt
}

func (d *eventData) GetChanges() (int, int) {
	return d.additions, d.deletions
}
