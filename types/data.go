package types

type Data interface {
	GetShortName() string
	GetLongName() string
	GetNumber() int
	GetTitle() string
	GetDescription() string
	GetState() string
	GetOrigin() Origin
	GetOwner() string
	GetRepo() string
	GetRef() string
	GetLabels() []string
	HasLabel(label string) bool
	AddLabel(label string)
	RemoveLabel(label string)
	GetAssignees() []string
	HasAssignee(assignee string) bool
	AddAssignees(assignees ...string)
	RemoveAssignees(assignees ...string)
	GetComments() []Comment
	AddComment(comment string)
	GetFileNames() []string
	GetChangedFiles() int
	GetFileExtensions() []string
	GetChanges() (int, int)
	GetProvider() string
	GetRawPayload() []byte
}
