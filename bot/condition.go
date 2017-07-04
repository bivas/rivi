package bot

import (
	"github.com/bivas/rivi/util"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

type sectionCondition interface {
	IsEmpty() bool
	Match(meta EventData) bool
}

type FilesCondition struct {
	Patterns   []string `mapstructure:"patterns,omitempty"`
	Extensions []string `mapstructure:"extensions,omitempty"`
}

func (c *FilesCondition) IsEmpty() bool {
	return len(c.Patterns) == 0 && len(c.Extensions) == 0
}

func (c *FilesCondition) checkPattern(meta EventData) bool {
	if len(c.Patterns) == 0 {
		return false
	} else {
		compiled := make([]*regexp.Regexp, 0)
		for _, pattern := range c.Patterns {
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
					util.Logger.Debug("Matched FileCondition with regex '%s' on file '%s'", reg.String(), check)
					return true
				}
			}
		}
	}
	return false
}

func (c *FilesCondition) checkExt(meta EventData) bool {
	if len(c.Extensions) == 0 {
		return false
	} else {
		for _, check := range meta.GetFileExtensions() {
			for _, ext := range c.Extensions {
				if ext == check {
					util.Logger.Debug("Matched FileCondition with extension '%s' on file '%s'", ext, check)
					return true
				}
			}
		}
	}
	return false
}

func (c *FilesCondition) Match(meta EventData) bool {
	return c.checkPattern(meta) || c.checkExt(meta)
}

type TitleCondition struct {
	StartsWith string   `mapstructure:"starts-with,omitempty"`
	EndsWith   string   `mapstructure:"ends-with,omitempty"`
	Patterns   []string `mapstructure:"patterns,omitempty"`
}

func (c *TitleCondition) IsEmpty() bool {
	return c.StartsWith == "" && c.EndsWith == "" && len(c.Patterns) == 0
}

func (c *TitleCondition) Match(meta EventData) bool {
	if c.StartsWith == "" && c.EndsWith == "" && len(c.Patterns) == 0 {
		return false
	} else {
		title := meta.GetTitle()
		if c.StartsWith != "" && strings.HasPrefix(title, c.StartsWith) {
			util.Logger.Debug("Matched TitleCondition with prefix '%s' on title '%s'", c.StartsWith, title)
			return true
		}
		if c.EndsWith != "" && strings.HasSuffix(title, c.EndsWith) {
			util.Logger.Debug("Matched TitleCondition with suffix '%s' on title '%s'", c.EndsWith, title)
			return true
		}
		if len(c.Patterns) > 0 {
			compiled := make([]*regexp.Regexp, 0)
			for _, pattern := range c.Patterns {
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
					util.Logger.Debug("Matched TitleCondition with pattern '%s' on title '%s'", reg.String(), title)
					return true
				}
			}
		}
	}
	return false
}

type RefCondition struct {
	Equals   string   `mapstructure:"match,omitempty"`
	Patterns []string `mapstructure:"patterns,omitempty"`
}

func (c *RefCondition) IsEmpty() bool {
	return c.Equals == "" && len(c.Patterns) == 0
}

func (c *RefCondition) Match(meta EventData) bool {
	ref := meta.GetRef()
	if c.Equals != "" && ref == c.Equals {
		util.Logger.Debug("Matched RefCondition with match on ref '%s'", ref)
		return true
	}
	if len(c.Patterns) > 0 {
		compiled := make([]*regexp.Regexp, 0)
		for _, pattern := range c.Patterns {
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
			if reg.MatchString(ref) {
				util.Logger.Debug("Matched RefCondition with regex '%s' on ref '%s'", reg.String(), ref)
				return true
			}
		}
	}
	return false
}

type Condition struct {
	Order         int            `mapstructure:"order,omitempty"`
	IfLabeled     []string       `mapstructure:"if-labeled,omitempty"`
	SkipIfLabeled []string       `mapstructure:"skip-if-labeled,omitempty"`
	Files         FilesCondition `mapstructure:"files,omitempty"`
	Title         TitleCondition `mapstructure:"title,omitempty"`
	Ref           RefCondition   `mapstructure:"ref,omitempty"`
}

func (c *Condition) checkIfLabeled(meta EventData) bool {
	if len(c.IfLabeled) == 0 {
		return false
	} else {
		for _, check := range c.IfLabeled {
			for _, label := range meta.GetLabels() {
				if check == label {
					util.Logger.Debug("Matched Condition with label '%s'", check)
					return true
				}
			}
		}
	}
	return false
}

func (c *Condition) checkAllEmpty(meta EventData) bool {
	empty := len(c.IfLabeled) == 0 &&
		c.Files.IsEmpty() &&
		c.Title.IsEmpty() &&
		c.Ref.IsEmpty()
	util.Logger.Debug("Condition is empty = %s", empty)
	return empty
}

func (c *Condition) Match(meta EventData) bool {
	match := c.checkAllEmpty(meta) ||
		c.checkIfLabeled(meta) ||
		c.Title.Match(meta) ||
		c.Files.Match(meta) ||
		c.Ref.Match(meta)

	if match {
		for _, check := range c.SkipIfLabeled {
			for _, label := range meta.GetLabels() {
				if check == label {
					util.Logger.Debug("Skipping Condition with label '%s'", check)
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
