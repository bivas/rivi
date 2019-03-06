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
	Head headSection `json:"head"`
}

type headSection struct {
	Ref  string `json:"ref"`
	Sha  string `json:"sha"`
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	Repo struct {
		Name   string `json:"name"`
		GitURL string `json:"git_url"`
	} `json:"repo"`
}

type repositorySection struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Private bool `json:"private"`
}

type PullRequestPayload struct {
	Action       string             `json:"action"`
	Number       int                `json:"number"`
	PullRequest  pullRequestSection `json:"pull_request"`
	Repository   repositorySection  `json:"repository"`
	Installation struct {
		ID int `json:"id"`
	} `json:"installation"`
}
