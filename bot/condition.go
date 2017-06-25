package bot

import (
	"github.com/bivas/rivi/util"
	"github.com/spf13/viper"
	"regexp"
)

type Condition struct {
	IfLabeled     []string `mapstructure:"if-labeled,omitempty"`
	SkipIfLabeled []string `mapstructure:"skip-if-labeled,omitempty"`
	Filter        struct {
		Pattern   string `mapstructure:"pattern,omitempty"`
		Extension string `mapstructure:"extension,omitempty"`
	} `mapstructure:"filter,omitempty"`
}

func (c *Condition) checkIfLabeled(meta EventData) bool {
	accept := false
	if len(c.IfLabeled) == 0 {
		accept = true
	} else {
		for _, check := range c.IfLabeled {
			for _, label := range meta.GetLabels() {
				accept = accept || check == label
			}
		}
	}
	return accept
}

func (c *Condition) checkPattern(meta EventData) bool {
	if c.Filter.Pattern == "" {
		return true
	} else {
		for _, check := range meta.GetFileNames() {
			matched, e := regexp.MatchString(c.Filter.Pattern, check)
			if e != nil {
				util.Logger.Debug("Error checking filter %s", e)
			} else if matched {
				return true
			}
		}
	}
	return false
}

func (c *Condition) checkExt(meta EventData) bool {
	if c.Filter.Extension == "" {
		return true
	} else {
		for _, check := range meta.GetFileExtensions() {
			accept := c.Filter.Extension == check
			if accept {
				return true
			}
		}
	}
	return false
}

func (c *Condition) Match(meta EventData) bool {
	match := c.checkIfLabeled(meta) && c.checkPattern(meta) && c.checkExt(meta)

	if match {
		for _, check := range c.SkipIfLabeled {
			for _, label := range meta.GetLabels() {
				if check == label {
					return false
				}
			}
		}
	}
	return match
}

func buildConditionFromConfiguration(config *viper.Viper) Condition {
	var result Condition
	exists := config.Get("condition")
	if exists != nil {
		condition := config.Sub("condition")
		if e := condition.Unmarshal(&result); e != nil {
			util.Logger.Error("Unable to unmarshal rule condition. %s", e)
		}
	}
	return result
}
