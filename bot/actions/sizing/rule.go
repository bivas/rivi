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
