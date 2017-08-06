package github

type pullRequestSection struct {
	Number       int    `json:"number"`
	State        string `json:"state"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	Commits      int    `json:"commits"`
	Additions    int    `json:"additions"`
	Deletions    int    `json:"deletions"`
	ChangedFiles int    `json:"changed_files"`
	Assignees    []struct {
		Login string `json:"login"`
	} `json:"assignees"`
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

type repositorySection struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type payload struct {
	Action      string             `json:"action"`
	Number      int                `json:"number"`
	PullRequest pullRequestSection `json:"pull_request"`
	Repository  repositorySection  `json:"repository"`
}
