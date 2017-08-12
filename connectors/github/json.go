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
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
	Head struct {
		Ref  string `json:"ref"`
		Sha  string `json:"sha"`
		User struct {
			Login string `json:"login"`
		} `json:"user"`
		Repo struct {
			GitURL string `json:"git_url"`
		} `json:"repo"`
	} `json:"head"`
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
