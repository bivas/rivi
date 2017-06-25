package github

type eventData struct {
	client       *ghClient
	number       int
	state        string
	origin       string
	owner        string
	repo         string
	title        string
	changedFiles int
	fileNames    []string
	fileExt      []string
	additions    int
	deletions    int
	labels       []string
	assignees    []string
	comments     []string
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

func (d *eventData) AddAssignees(assignees... string) {
	d.assignees = d.client.AddAssignees(d.number, assignees...)
}

func (d *eventData) RemoveAssignees(assignees... string) {
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

func (d *eventData) GetComments() []string {
	return d.comments
}

func (d *eventData) AddComment(comment string) {
	d.comments = append(d.comments, comment)
	d.client.AddComment(d.number, comment)
}

func (d *eventData) GetNumber() int {
	return d.number
}

func (d *eventData) GetTitle() string {
	return d.title
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
