package bot

import (
	"github.com/bivas/rivi/util"
	"github.com/spf13/viper"
	"regexp"
)

type Condition struct {
	Order         int      `mapstructure:"order,omitempty"`
	IfLabeled     []string `mapstructure:"if-labeled,omitempty"`
	SkipIfLabeled []string `mapstructure:"skip-if-labeled,omitempty"`
	Filter        struct {
		Patterns   []string `mapstructure:"patterns,omitempty"`
		Extensions []string `mapstructure:"extensions,omitempty"`
	} `mapstructure:"filter,omitempty"`
}

func (c *Condition) checkIfLabeled(meta EventData) bool {
	if len(c.IfLabeled) == 0 {
		return false
	} else {
		for _, check := range c.IfLabeled {
			for _, label := range meta.GetLabels() {
				if check == label {
					return true
				}
			}
		}
	}
	return false
}

func (c *Condition) checkPattern(meta EventData) bool {
	if len(c.Filter.Patterns) == 0 {
		return false
	} else {
		compiled := make([]*regexp.Regexp, 0)
		for _, pattern := range c.Filter.Patterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				util.Logger.Warning("Unable to compile regex '%s'. %s", pattern, err)
				continue
			}
			compiled = append(compiled, re)
		}
		if len(compiled) == 0 {
			util.Logger.Error("All configured patterns have failed to compile")
			return false
		}
		for _, check := range meta.GetFileNames() {
			for _, reg := range compiled {
				if reg.MatchString(check) {
					return true
				}
			}
		}
	}
	return false
}

func (c *Condition) checkExt(meta EventData) bool {
	if len(c.Filter.Extensions) == 0 {
		return false
	} else {
		for _, check := range meta.GetFileExtensions() {
			for _, ext := range c.Filter.Extensions {
				if ext == check {
					return true
				}
			}
		}
	}
	return false
}

func (c *Condition) checkAllEmpty(meta EventData) bool {
	return len(c.IfLabeled) == 0 && len(c.Filter.Patterns) == 0 && len(c.Filter.Extensions) == 0
}

func (c *Condition) Match(meta EventData) bool {
	match := c.checkAllEmpty(meta) || c.checkIfLabeled(meta) || c.checkPattern(meta) || c.checkExt(meta)

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
