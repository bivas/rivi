package engine

import (
	"fmt"

	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util"
)

type mockData struct {
	Number           int
	Title            string
	Description      string
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
	FileNames        []string
	FileExtensions   []string
	Comments         []types.Comment
	AddedComments    []string
	ChangedFiles     int
	ChangesAdd       int
	ChangesRemove    int
	RawPayload       []byte
}

func (m *mockData) GetProvider() string {
	return "mock"
}

func (m *mockData) GetShortName() string {
	return fmt.Sprintf("%s/%s#%d", m.Owner, m.Repo, m.Number)
}

func (m *mockData) GetLongName() string {
	return fmt.Sprintf("%s/%s#%d [%s]", m.Owner, m.Repo, m.Number, m.Title)
}

func (m *mockData) GetRawPayload() []byte {
	return m.RawPayload
}

func (m *mockData) GetNumber() int {
	return m.Number
}

func (m *mockData) GetTitle() string {
	return m.Title
}

func (m *mockData) GetDescription() string {
	return m.Description
}

func (m *mockData) GetState() string {
	return m.State
}

func (m *mockData) GetOrigin() string {
	return m.Origin
}

func (m *mockData) GetOwner() string {
	return m.Owner
}

func (m *mockData) GetRepo() string {
	return m.Repo
}

func (m *mockData) GetRef() string {
	return m.Ref
}

func (m *mockData) GetLabels() []string {
	return m.Labels
}

func (m *mockData) HasLabel(label string) bool {
	for _, search := range m.Labels {
		if search == label {
			return true
		}
	}
	return false
}

func (m *mockData) AddLabel(label string) {
	m.AddedLabels = append(m.AddedLabels, label)
	m.Labels = append(m.Labels, label)
}

func (m *mockData) RemoveLabel(label string) {
	m.RemovedLabels = append(m.RemovedLabels, label)
	set := util.StringSet{}
	set.AddAll(m.Labels).Remove(label)
	m.Labels = set.Values()
}

func (m *mockData) GetAssignees() []string {
	return m.Assignees
}

func (m *mockData) HasAssignee(assignee string) bool {
	for _, search := range m.Assignees {
		if search == assignee {
			return true
		}
	}
	return false
}

func (m *mockData) AddAssignees(assignees ...string) {
	m.AddedAssignees = append(m.AddedAssignees, assignees...)
	m.Assignees = append(m.Assignees, assignees...)
}

func (m *mockData) RemoveAssignees(assignees ...string) {
	m.RemovedLabels = append(m.RemovedLabels, assignees...)
}

func (m *mockData) GetComments() []types.Comment {
	return m.Comments
}

func (m *mockData) AddComment(comment string) {
	m.AddedComments = append(m.AddedComments, comment)
	m.Comments = append(m.Comments, types.Comment{Commenter: "mock", Comment: comment})
}

func (m *mockData) GetFileNames() []string {
	return m.FileNames
}

func (m *mockData) GetChangedFiles() int {
	return m.ChangedFiles
}

func (m *mockData) GetFileExtensions() []string {
	return m.FileExtensions
}

func (m *mockData) GetChanges() (int, int) {
	return m.ChangesAdd, m.ChangesRemove
}
