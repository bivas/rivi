package automerge

import "strings"

var (
	approvedPhrases = []string{
		"approve",
		"approved",
		"lgtm",
		"looks good to me",
	}
	approvedSearchPhrases map[string]bool
	strategies            = []string{
		"merge",
		"squash",
		"rebase",
	}
	searchStrategy map[string]bool
)

type rule struct {
	Require  int    `mapstructure:"require"`
	Strategy string `mapstructure:"strategy"`
}

func (r *rule) Defaults() {
	if r.Require == 0 {
		r.Require = 1
	}
	if r.Strategy == "" {
		r.Strategy = "merge"
	} else {
		search := strings.ToLower(r.Strategy)
		if _, ok := searchStrategy[search]; !ok {
			r.Strategy = "merge"
		}
	}
}

func init() {
	approvedSearchPhrases = make(map[string]bool)
	for _, a := range approvedPhrases {
		approvedSearchPhrases[a] = true
	}
	searchStrategy = make(map[string]bool)
	for _, s := range strategies {
		searchStrategy[s] = true
	}
}
