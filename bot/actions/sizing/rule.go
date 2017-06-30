package sizing

import (
	"math"
)

type sizingRule struct {
	Name                  string
	Label                 string `mapstructure:"label"`
	Comment               string `mapstructure:"comment"`
	ChangedFilesThreshold int    `mapstructure:"changed-files-threshold"`
	ChangesThreshold      int    `mapstructure:"changes-threshold"`
}

func (rule *sizingRule) Defaults() {
	if rule.ChangesThreshold == 0 {
		rule.ChangesThreshold = math.MaxInt32
	}
	if rule.ChangedFilesThreshold == 0 {
		rule.ChangedFilesThreshold = math.MaxInt32
	}
}

type rules []sizingRule

func (r rules) Len() int {
	return len(r)
}

func (r rules) Less(i, j int) bool {
	return r[i].ChangedFilesThreshold < r[j].ChangedFilesThreshold && r[i].ChangesThreshold < r[j].ChangesThreshold
}

func (r rules) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
