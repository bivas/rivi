package status

import (
	"strings"
)

const (
	defaultDescription = "Rivi completed processing rules"
	defaultState       = "failure"
	unknownState       = "error"
)

var (
	states = []string{
		"error",
		"failure",
		"pending",
		"success",
	}
	searchState map[string]bool
)

type rule struct {
	State       string `mapstructure:"state"`
	Description string `mapstructure:"description"`
}

func (r *rule) Defaults() {
	if r.State == "" {
		r.State = defaultState
	} else {
		search := strings.ToLower(r.State)
		if _, ok := searchState[search]; !ok {
			r.State = unknownState
		}
	}
	if r.Description == "" {
		r.Description = defaultDescription
	}
}

func init() {
	searchState = make(map[string]bool)
	for _, s := range states {
		searchState[s] = true
	}
}
