package mock

import (
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
)

type MockEventData struct {
	Number           int
	Title            string
	State            string
	Owner            string
	Repo             string
	Ref              string
	Origin           string
	Assignees        []string
	AddedAssignees   []string
	RemovedAssignees []string
	Labels           []string
	AddedLabels      []string
	RemovedLabels    []string
	FileExtensions   []string
	Comments         []bot.Comment
	AddedComments    []string
	ChangedFiles     int
	ChangesAdd       int
	ChangesRemove    int
	RawPayload       []byte
}

func (m *MockEventData) GetRawPayload() []byte {
	return m.RawPayload
}

func (m *MockEventData) GetNumber() int {
	return m.Number
}

func (m *MockEventData) GetTitle() string {
	return m.Title
}

func (m *MockEventData) GetState() string {
	return m.State
}

func (m *MockEventData) GetOrigin() string {
	return m.Origin
}

func (m *MockEventData) GetOwner() string {
	return m.Owner
}

func (m *MockEventData) GetRepo() string {
	return m.Repo
}

func (m *MockEventData) GetRef() string {
	return m.Ref
}

func (m *MockEventData) GetLabels() []string {
	return m.Labels
}

func (m *MockEventData) HasLabel(label string) bool {
	for _, search := range m.Labels {
		if search == label {
			return true
		}
	}
	return false
}

func (m *MockEventData) AddLabel(label string) {
	m.AddedLabels = append(m.AddedLabels, label)
	m.Labels = append(m.Labels, label)
}

func (m *MockEventData) RemoveLabel(label string) {
	m.RemovedLabels = append(m.RemovedLabels, label)
	set := util.StringSet{}
	set.AddAll(m.Labels).Remove(label)
	m.Labels = set.Values()
}

func (m *MockEventData) GetAssignees() []string {
	return m.Assignees
}

func (m *MockEventData) HasAssignee(assignee string) bool {
	for _, search := range m.Assignees {
		if search == assignee {
			return true
		}
	}
	return false
}

func (m *MockEventData) AddAssignees(assignees ...string) {
	m.AddedAssignees = append(m.AddedAssignees, assignees...)
	m.Assignees = append(m.Assignees, assignees...)
}

func (m *MockEventData) RemoveAssignees(assignees ...string) {
	m.RemovedLabels = append(m.RemovedLabels, assignees...)
}

func (m *MockEventData) GetComments() []bot.Comment {
	return m.Comments
}

func (m *MockEventData) AddComment(comment string) {
	m.AddedComments = append(m.AddedComments, comment)
	m.Comments = append(m.Comments, bot.Comment{Commenter: "mock", Comment: comment})
}

func (*MockEventData) GetFileNames() []string {
	panic("implement me")
}

func (m *MockEventData) GetChangedFiles() int {
	return m.ChangedFiles
}

func (m *MockEventData) GetFileExtensions() []string {
	return m.FileExtensions
}

func (m *MockEventData) GetChanges() (int, int) {
	return m.ChangesAdd, m.ChangesRemove
}
