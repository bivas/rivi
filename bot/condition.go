package bot

import (
	"github.com/bivas/rivi/util"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

type Condition struct {
	Order         int      `mapstructure:"order,omitempty"`
	IfLabeled     []string `mapstructure:"if-labeled,omitempty"`
	SkipIfLabeled []string `mapstructure:"skip-if-labeled,omitempty"`
	Files         struct {
		Patterns   []string `mapstructure:"patterns,omitempty"`
		Extensions []string `mapstructure:"extensions,omitempty"`
	} `mapstructure:"files,omitempty"`
	Title struct {
		StartsWith string   `mapstructure:"starts-with,omitempty"`
		EndsWith   string   `mapstructure:"ends-with,omitempty"`
		Patterns   []string `mapstructure:"patterns,omitempty"`
	} `mapstructure:"title,omitempty"`
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
	if len(c.Files.Patterns) == 0 {
		return false
	} else {
		compiled := make([]*regexp.Regexp, 0)
		for _, pattern := range c.Files.Patterns {
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
	if len(c.Files.Extensions) == 0 {
		return false
	} else {
		for _, check := range meta.GetFileExtensions() {
			for _, ext := range c.Files.Extensions {
				if ext == check {
					return true
				}
			}
		}
	}
	return false
}

func (c *Condition) checkTitle(meta EventData) bool {
	titleCondition := c.Title
	if titleCondition.StartsWith == "" && titleCondition.EndsWith == "" && len(titleCondition.Patterns) == 0 {
		return false
	} else {
		title := meta.GetTitle()
		if titleCondition.StartsWith != "" && strings.HasPrefix(title, titleCondition.StartsWith) {
			return true
		}
		if titleCondition.EndsWith != "" && strings.HasSuffix(title, titleCondition.EndsWith) {
			return true
		}
		if len(titleCondition.Patterns) > 0 {
			compiled := make([]*regexp.Regexp, 0)
			for _, pattern := range titleCondition.Patterns {
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
			for _, reg := range compiled {
				if reg.MatchString(title) {
					return true
				}
			}
		}
	}
	return false
}

func (c *Condition) checkFiles(meta EventData) bool {
	return c.checkPattern(meta) || c.checkExt(meta)
}

func (c *Condition) checkAllEmpty(meta EventData) bool {
	return len(c.IfLabeled) == 0 &&
		len(c.Files.Patterns) == 0 && len(c.Files.Extensions) == 0 &&
		c.Title.StartsWith == "" && c.Title.EndsWith == "" && len(c.Title.Patterns) == 0
}

func (c *Condition) Match(meta EventData) bool {
	match := c.checkAllEmpty(meta) || c.checkIfLabeled(meta) || c.checkTitle(meta) || c.checkFiles(meta)

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
