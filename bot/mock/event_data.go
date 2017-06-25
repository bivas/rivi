package mock

type MockEventData struct {
	Number           int
	Title            string
	State            string
	Owner            string
	Repo             string
	Origin           string
	Assignees        []string
	AddedAssignees   []string
	RemovedAssignees []string
	Labels           []string
	AddedLabels      []string
	RemovedLabels    []string
	FileExtensions   []string
	Comments         []string
	AddedComments    []string
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

func (m *MockEventData) GetComments() []string {
	return m.Comments
}

func (m *MockEventData) AddComment(comment string) {
	m.AddedComments = append(m.AddedComments, comment)
	m.Comments = append(m.Comments, comment)
}

func (*MockEventData) GetFileNames() []string {
	panic("implement me")
}

func (*MockEventData) GetChangedFiles() int {
	panic("implement me")
}

func (m *MockEventData) GetFileExtensions() []string {
	return m.FileExtensions
}

func (*MockEventData) GetChanges() (int, int) {
	panic("implement me")
}
