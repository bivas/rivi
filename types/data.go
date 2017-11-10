package types

type InfoData interface {
	GetShortName() string
	GetLongName() string
	GetProvider() string
	GetRepository() Repository
}

type RawData interface {
	GetRawPayload() []byte
	GetRawType() string
}

type ReadOnlyData interface {
	InfoData
	RawData
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
	GetAssignees() []string
	HasAssignee(assignee string) bool
	GetComments() []Comment
	GetFileNames() []string
	GetChangedFiles() int
	GetFileExtensions() []string
	GetChanges() (int, int)
}

type MutableData interface {
	AddLabel(label string)
	RemoveLabel(label string)
	AddAssignees(assignees ...string)
	RemoveAssignees(assignees ...string)
	AddComment(comment string)
}

type HookData interface {
	ReadOnlyData
}

type Data interface {
	ReadOnlyData
	MutableData
}
