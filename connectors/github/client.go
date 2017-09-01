package github

import (
	"context"
	"strings"

	"github.com/bivas/rivi/config/client"
	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
)

type ghClient struct {
	client *github.Client
	secret []byte
	owner  string
	repo   string

	logger log.Logger
}

func (c *ghClient) d() {
}

func (c *ghClient) handleFileContentResponse(file *github.RepositoryContent, err error, fields *log.MetaFields) string {
	if err != nil {
		c.logger.ErrorWith(*fields, "Unable to get file")
	}
	content, err := file.GetContent()
	if err != nil {
		c.logger.ErrorWith(*fields, "Unable to get file content")
		content = ""
	}
	return content
}

func (c *ghClient) GetFileContentFromRef(path, owner, repo, ref string) string {
	opts := &github.RepositoryContentGetOptions{Ref: ref}
	file, _, _, err := c.client.Repositories.GetContents(context.Background(), owner, repo, path, opts)
	return c.handleFileContentResponse(file, err,
		&log.MetaFields{
			log.E(err),
			log.F("repo", repo),
			log.F("owner", owner),
			log.F("path", path),
			log.F("ref", ref)})
}

func (c *ghClient) GetFileContent(path string) string {
	file, _, _, err := c.client.Repositories.GetContents(context.Background(), c.owner, c.repo, path, nil)
	return c.handleFileContentResponse(file, err,
		&log.MetaFields{
			log.E(err),
			log.F("repo", c.repo),
			log.F("owner", c.owner),
			log.F("path", path)})
}

func (c *ghClient) GetCollaborators() []string {
	users, _, err := c.client.Repositories.ListCollaborators(context.Background(), c.owner, c.repo, nil)
	result := make([]string, 0)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("repo", c.repo)}, "Unable to get collaborators")
	} else {
		for _, user := range users {
			result = append(result, strings.ToLower(*user.Login))
		}
	}
	return result
}

func (c *ghClient) GetState(issue int) string {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Getting issue state")
	response, _, err := c.client.Issues.Get(context.Background(), c.owner, c.repo, issue)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to get issue state.")
		return ""
	}
	return *response.State
}

func handleLabelsResult(labels []*github.Label, err error, logError func(error)) []string {
	result := make([]string, 0)
	if err != nil {
		logError(err)
	} else {
		for _, p := range labels {
			result = append(result, *p.Name)
		}
	}
	return result
}

func (c *ghClient) Lock(issue int) {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Locking issue")
	_, err := c.client.Issues.Lock(context.Background(), c.owner, c.repo, issue)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to set issue lock state")
	}
}

func (c *ghClient) Unlock(issue int) {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Unlocking issue")
	_, err := c.client.Issues.Unlock(context.Background(), c.owner, c.repo, issue)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to set issue unlock state")
	}
}

func (c *ghClient) Locked(issue int) bool {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Checking lock state")
	response, _, err := c.client.Issues.Get(context.Background(), c.owner, c.repo, issue)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to get issue lock state")
	}
	return *response.Locked
}

func (c *ghClient) GetAvailableLabels() []string {
	c.logger.Debug("Getting available labels")
	labels, _, e := c.client.Issues.ListLabels(context.Background(), c.owner, c.repo, nil)
	return handleLabelsResult(labels, e, func(err error) {
		c.logger.ErrorWith(log.MetaFields{log.E(err)}, "Unable to get available labels")
	})
}

func (c *ghClient) GetLabels(issue int) []string {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Getting labels")
	labels, _, err := c.client.Issues.ListLabelsByIssue(context.Background(), c.owner, c.repo, issue, nil)
	return handleLabelsResult(labels, err, func(err error) {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to get labels for issue.")
	})
}

func (c *ghClient) AddLabel(issue int, label string) []string {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Adding label '%s'", label)
	labels, _, err := c.client.Issues.AddLabelsToIssue(context.Background(), c.owner, c.repo, issue, []string{label})
	return handleLabelsResult(labels, err, func(err error) {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to add label '%s' for issue", label)
	})
}

func (c *ghClient) RemoveLabel(issue int, label string) []string {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Removing label '%s'", label)
	c.client.Issues.RemoveLabelForIssue(context.Background(), c.owner, c.repo, issue, label)
	return c.GetLabels(issue)
}

func (c *ghClient) GetAssignees(issue int) []string {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Getting assignees")
	result := make([]string, 0)
	issueObject, _, err := c.client.Issues.Get(context.Background(), c.owner, c.repo, issue)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to get assignees for issue")
	} else {
		for _, p := range issueObject.Assignees {
			result = append(result, strings.ToLower(*p.Login))
		}
	}
	return result
}

func (c *ghClient) AddAssignees(issue int, assignees ...string) []string {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Adding assignees %s", assignees)
	response, _, err := c.client.Issues.AddAssignees(context.Background(), c.owner, c.repo, issue, assignees)
	result := make([]string, 0)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to add assignees %s for issue", assignees)
	} else {
		for _, p := range response.Assignees {
			result = append(result, *p.Login)
		}
	}
	return result

}

func (c *ghClient) RemoveAssignees(issue int, assignees ...string) []string {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Removing assignees %s", assignees)
	response, _, err := c.client.Issues.RemoveAssignees(context.Background(), c.owner, c.repo, issue, assignees)
	result := make([]string, 0)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to remove assignees %s for issue", assignees)
	} else {
		for _, p := range response.Assignees {
			result = append(result, *p.Login)
		}
	}
	return result
}

func (c *ghClient) GetFileNames(issue int) []string {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Getting file names")
	files, _, err := c.client.PullRequests.ListFiles(context.Background(), c.owner, c.repo, issue, nil)
	result := make([]string, 0)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to get file names for issue")
	} else {
		for _, p := range files {
			result = append(result, *p.Filename)
		}
	}
	return result
}

func (c *ghClient) GetComments(issue int) []types.Comment {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Getting comments")
	comments, _, err := c.client.Issues.ListComments(context.Background(), c.owner, c.repo, issue, nil)
	result := make([]types.Comment, 0)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to get comments for issue")
	} else {
		for _, p := range comments {
			comment := types.Comment{
				Commenter: *p.User.Login,
				Comment:   *p.Body,
			}
			result = append(result, comment)
		}
	}
	return result
}

func (c *ghClient) AddComment(issue int, comment string) types.Comment {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue), log.F("comment", comment)}, "Adding comment")
	commentObject := &github.IssueComment{Body: github.String(comment)}
	posted, _, err := c.client.Issues.CreateComment(context.Background(), c.owner, c.repo, issue, commentObject)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to add comment for issue")
		return types.Comment{}
	} else {
		return types.Comment{
			Commenter: *posted.User.Login,
			Comment:   *posted.Body,
		}
	}
}

func (c *ghClient) GetReviewers(issue int) map[string]string {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue)}, "Getting reviewers")
	result := make(map[string]string)
	reviews, _, err := c.client.PullRequests.ListReviews(context.Background(), c.owner, c.repo, issue, nil)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Unable to get reviewers for issue")
		return result
	}
	for _, review := range reviews {
		user := strings.ToLower(*review.User.Login)
		state := strings.ToLower(*review.State)
		result[user] = state
	}
	c.logger.DebugWith(
		log.MetaFields{log.F("issue.id", issue), log.F("reviewers", result)},
		"Created reviewers map")
	return result
}

func (c *ghClient) Merge(issue int, mergeMethod string) {
	c.logger.DebugWith(log.MetaFields{log.F("issue.id", issue), log.F("strategy", mergeMethod)}, "Merging")
	options := &github.PullRequestOptions{MergeMethod: mergeMethod}
	_, _, err := c.client.PullRequests.Merge(context.Background(), c.owner, c.repo, issue, "", options)
	if err != nil {
		c.logger.ErrorWith(log.MetaFields{log.E(err), log.F("issue.id", issue)}, "Error trying to merge issue")
	}
}

func newClient(config client.ClientConfig, owner, repo string) *ghClient {
	source := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GetOAuthToken()},
	)
	c := oauth2.NewClient(context.Background(), source)

	return &ghClient{
		client: github.NewClient(c),
		secret: []byte(config.GetSecret()),
		owner:  owner,
		repo:   repo,
		logger: log.Get("GitHub.Client"),
	}
}

func newAppClient(config client.ClientConfig, owner, repo string, installation int) *ghClient {
	logger := log.Get("Github.Client")
	c, err := ghinstallation.NewKeyFromFile(
		http.DefaultTransport,
		config.GetApplicationID(),
		installation,
		"rivi.private-key.pem",
	)
	if err != nil {
		logger.ErrorWith(
			log.MetaFields{
				log.E(err),
				log.F("appid", config.GetApplicationID()),
				log.F("installation", installation),
				log.F("owner", owner),
				log.F("repo", repo),
			}, "Unable to create installation client",
		)
		return nil
	}
	return &ghClient{
		client: github.NewClient(&http.Client{Transport: c}),
		secret: []byte(config.GetSecret()),
		owner:  owner,
		repo:   repo,
		logger: logger,
	}
}
