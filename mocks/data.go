package mocks

import (
	"fmt"

	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util"
)

type MockData struct {
	Number            int
	Title             string
	Description       string
	State             string
	Owner             string
	Repo              string
	Ref               string
	Origin            types.Origin
	Assignees         []string
	AddedAssignees    []string
	RemovedAssignees  []string
	Labels            []string
	AddedLabels       []string
	RemovedLabels     []string
	FileNames         []string
	FileExtensions    []string
	Comments          []types.Comment
	AddedComments     []string
	ChangedFiles      int
	ChangesAdd        int
	ChangesRemove     int
	Provider          string
	RawPayload        []byte
	Collaborators     []string
	AvailableLabels   []string
	RulesFileContent  string
	EventType         string
	StatusDescription string
	StatusState       types.State
}

func (m *MockData) SetStatus(desc string, state types.State) {
	m.StatusDescription = desc
	m.StatusState = state
}

func (m *MockData) GetRawType() string {
	return m.EventType
}

func (m *MockData) GetCollaborators() []string {
	return m.Collaborators
}

func (m *MockData) IsCollaborator(name string) bool {
	for _, collaborator := range m.Collaborators {
		if name == collaborator {
			return true
		}
	}
	return false
}

func (m *MockData) GetAvailableLabels() []string {
	return m.AvailableLabels
}

func (m *MockData) GetRulesFile() string {
	return m.RulesFileContent
}

func (m *MockData) GetRepository() types.Repository {
	return m
}

func (m *MockData) GetProvider() string {
	return m.Provider
}

func (m *MockData) GetShortName() string {
	return fmt.Sprintf("%s/%s#%d", m.Owner, m.Repo, m.Number)
}

func (m *MockData) GetLongName() string {
	return fmt.Sprintf("%s/%s#%d [%s]", m.Owner, m.Repo, m.Number, m.Title)
}

func (m *MockData) GetRawPayload() []byte {
	return m.RawPayload
}

func (m *MockData) GetNumber() int {
	return m.Number
}

func (m *MockData) GetTitle() string {
	return m.Title
}

func (m *MockData) GetDescription() string {
	return m.Description
}

func (m *MockData) GetState() string {
	return m.State
}

func (m *MockData) GetOrigin() types.Origin {
	return m.Origin
}

func (m *MockData) GetOwner() string {
	return m.Owner
}

func (m *MockData) GetRepo() string {
	return m.Repo
}

func (m *MockData) GetRef() string {
	return m.Ref
}

func (m *MockData) GetLabels() []string {
	return m.Labels
}

func (m *MockData) HasLabel(label string) bool {
	for _, search := range m.Labels {
		if search == label {
			return true
		}
	}
	return false
}

func (m *MockData) AddLabel(label string) {
	m.AddedLabels = append(m.AddedLabels, label)
	m.Labels = append(m.Labels, label)
}

func (m *MockData) RemoveLabel(label string) {
	m.RemovedLabels = append(m.RemovedLabels, label)
	set := util.StringSet{}
	set.AddAll(m.Labels).Remove(label)
	m.Labels = set.Values()
}

func (m *MockData) GetAssignees() []string {
	return m.Assignees
}

func (m *MockData) HasAssignee(assignee string) bool {
	for _, search := range m.Assignees {
		if search == assignee {
			return true
		}
	}
	return false
}

func (m *MockData) AddAssignees(assignees ...string) {
	m.AddedAssignees = append(m.AddedAssignees, assignees...)
	m.Assignees = append(m.Assignees, assignees...)
}

func (m *MockData) RemoveAssignees(assignees ...string) {
	m.RemovedLabels = append(m.RemovedLabels, assignees...)
}

func (m *MockData) GetComments() []types.Comment {
	return m.Comments
}

func (m *MockData) AddComment(comment string) {
	m.AddedComments = append(m.AddedComments, comment)
	m.Comments = append(m.Comments, types.Comment{Commenter: "mock", Comment: comment})
}

func (m *MockData) GetFileNames() []string {
	return m.FileNames
}

func (m *MockData) GetChangedFiles() int {
	return m.ChangedFiles
}

func (m *MockData) GetFileExtensions() []string {
	return m.FileExtensions
}

func (m *MockData) GetChanges() (int, int) {
	return m.ChangesAdd, m.ChangesRemove
}
