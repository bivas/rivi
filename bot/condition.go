package bot

import (
	"regexp"
	"strings"

	"strconv"

	"github.com/bivas/rivi/util/log"
	"github.com/spf13/viper"
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
				log.WarningWith(
					log.MetaFields{log.F("condition", "FileCondition"), log.F("issue", meta.GetShortName()), log.E(err)},
					"Unable to compile regex '%s'", pattern)
				continue
			}
			compiled = append(compiled, re)
		}
		if len(compiled) == 0 {
			log.ErrorWith(
				log.MetaFields{log.F("condition", "FileCondition"), log.F("issue", meta.GetShortName())},
				"All configured patterns have failed to compile")
			return false
		}
		for _, check := range meta.GetFileNames() {
			for _, reg := range compiled {
				if reg.MatchString(check) {
					log.DebugWith(
						log.MetaFields{log.F("condition", "FileCondition"), log.F("issue", meta.GetShortName())},
						"Matched with regex '%s' on file '%s'", reg.String(), check)
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
					log.DebugWith(
						log.MetaFields{log.F("condition", "FileCondition"), log.F("issue", meta.GetShortName())},
						"Matched with extension '%s' on file '%s'", ext, check)
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
	title := meta.GetTitle()
	if c.StartsWith != "" && strings.HasPrefix(title, c.StartsWith) {
		log.DebugWith(
			log.MetaFields{log.F("condition", "TitleCondition"), log.F("issue", meta.GetShortName())},
			"Matched with prefix '%s' on title '%s'", c.StartsWith, title)
		return true
	}
	if c.EndsWith != "" && strings.HasSuffix(title, c.EndsWith) {
		log.DebugWith(
			log.MetaFields{log.F("condition", "TitleCondition"), log.F("issue", meta.GetShortName())},
			"Matched with suffix '%s' on title '%s'", c.EndsWith, title)
		return true
	}
	if len(c.Patterns) > 0 {
		compiled := make([]*regexp.Regexp, 0)
		for _, pattern := range c.Patterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				log.WarningWith(
					log.MetaFields{log.F("condition", "TitleCondition"), log.F("issue", meta.GetShortName()), log.E(err)},
					"Unable to compile regex '%s'", pattern)
				continue
			}
			compiled = append(compiled, re)
		}
		if len(compiled) == 0 {
			log.ErrorWith(
				log.MetaFields{log.F("condition", "TitleCondition"), log.F("issue", meta.GetShortName())},
				"All configured patterns have failed to compile")
			return false
		}
		for _, reg := range compiled {
			if reg.MatchString(title) {
				log.DebugWith(
					log.MetaFields{log.F("condition", "TitleCondition"), log.F("issue", meta.GetShortName())},
					"Matched with pattern '%s' on title '%s'", reg.String(), title)
				return true
			}
		}
	}
	return false
}

type DescriptionCondition struct {
	StartsWith string   `mapstructure:"starts-with,omitempty"`
	EndsWith   string   `mapstructure:"ends-with,omitempty"`
	Patterns   []string `mapstructure:"patterns,omitempty"`
}

func (c *DescriptionCondition) IsEmpty() bool {
	return c.StartsWith == "" && c.EndsWith == "" && len(c.Patterns) == 0
}

func (c *DescriptionCondition) Match(meta EventData) bool {
	description := meta.GetDescription()
	if c.StartsWith != "" && strings.HasPrefix(description, c.StartsWith) {
		log.DebugWith(
			log.MetaFields{log.F("condition", "DescriptionCondition"), log.F("issue", meta.GetShortName())},
			"Matched with prefix '%s' on description '%s'",
			c.StartsWith,
			description)
		return true
	}
	if c.EndsWith != "" && strings.HasSuffix(description, c.EndsWith) {
		log.DebugWith(
			log.MetaFields{log.F("condition", "DescriptionCondition"), log.F("issue", meta.GetShortName())},
			"Matched with suffix '%s' on description '%s'",
			c.EndsWith,
			description)
		return true
	}
	if len(c.Patterns) > 0 {
		compiled := make([]*regexp.Regexp, 0)
		for _, pattern := range c.Patterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				log.WarningWith(log.MetaFields{log.F("condition", "DescriptionCondition"), log.F("issue", meta.GetShortName()), log.E(err)},
					"Unable to compile regex '%s'", pattern)
				continue
			}
			compiled = append(compiled, re)
		}
		if len(compiled) == 0 {
			log.ErrorWith(
				log.MetaFields{log.F("condition", "DescriptionCondition"), log.F("issue", meta.GetShortName())},
				"All configured patterns have failed to compile")
			return false
		}
		for _, reg := range compiled {
			if reg.MatchString(description) {
				log.DebugWith(
					log.MetaFields{log.F("condition", "DescriptionCondition"), log.F("issue", meta.GetShortName())},
					"Matched with pattern '%s' on description '%s'",
					reg.String(),
					description)
				return true
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
		log.DebugWith(
			log.MetaFields{log.F("condition", "RefCondition"), log.F("issue", meta.GetShortName())},
			"Matched RefCondition with match on ref '%s'", ref)
		return true
	}
	if len(c.Patterns) > 0 {
		compiled := make([]*regexp.Regexp, 0)
		for _, pattern := range c.Patterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				log.WarningWith(
					log.MetaFields{log.F("condition", "RefCondition"), log.F("issue", meta.GetShortName()), log.E(err)},
					"Unable to compile regex '%s'", pattern)
				continue
			}
			compiled = append(compiled, re)
		}
		if len(compiled) == 0 {
			log.ErrorWith(
				log.MetaFields{log.F("condition", "RefCondition"), log.F("issue", meta.GetShortName())},
				"All configured patterns have failed to compile")
			return false
		}
		for _, reg := range compiled {
			if reg.MatchString(ref) {
				log.DebugWith(
					log.MetaFields{log.F("condition", "RefCondition"), log.F("issue", meta.GetShortName())},
					"Matched with regex '%s' on ref '%s'", reg.String(), ref)
				return true
			}
		}
	}
	return false
}

var commentsRegex = silentCompile("^([><]{1}[=]?|[=]{2}|)[ ]*([0-9]+)[ ]*$")

func silentCompile(pattern string) *regexp.Regexp {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		log.ErrorWith(
			log.MetaFields{log.F("pattern", pattern), log.E(err)},
			"Unable to compile comments regex",
		)
		return nil
	}
	return compiled
}

type CommentsCondition struct {
	Count string `mapstructure:"count,omitempty"`
}

func (c *CommentsCondition) IsEmpty() bool {
	return c.Count == ""
}

func (c *CommentsCondition) Match(meta EventData) bool {
	if commentsRegex == nil {
		log.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"),
				log.F("issue", meta.GetShortName())},
			"comments regex is nil'")
		return false
	}
	if c.Count == "" {
		return false
	}
	count := int64(len(meta.GetComments()))
	countGroups := commentsRegex.FindStringSubmatch(c.Count)
	if len(countGroups) != 3 {
		log.WarningWith(
			log.MetaFields{log.F("condition", "CommentsCondition"),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups)},
			"No groups matched'")
		return false
	}
	query, err := strconv.ParseInt(countGroups[2], 10, 64)
	if err != nil {
		log.WarningWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.E(err),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups)},
			"No groups matched")
		return false
	}
	matched := false
	op := countGroups[1]
	switch {
	case op == "" || op == "==":
		//exact number or equals
		matched = count == query
		log.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.F("matcher", op),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups),
				log.F("matched", matched), log.F("comments.count", count)},
			"matching with '=='")
		break
	case op == ">":
		// gt
		matched = count > query
		log.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.F("matcher", op),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups),
				log.F("matched", matched), log.F("comments.count", count)},
			"matching with '>'")
	case op == ">=":
		// gte
		matched = count >= query
		log.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.F("matcher", op),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups),
				log.F("matched", matched), log.F("comments.count", count)},
			"matching with '>='")
	case op == "<":
		// lt
		matched = count < query
		log.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.F("matcher", op),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups),
				log.F("matched", matched), log.F("comments.count", count)},
			"matching with '<'")
	case op == "<=":
		// lte
		matched = count <= query
		log.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.F("matcher", op),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups),
				log.F("matched", matched), log.F("comments.count", count)},
			"matching with '<='")
	}
	return matched
}

type Condition struct {
	Order         int                  `mapstructure:"order,omitempty"`
	IfLabeled     []string             `mapstructure:"if-labeled,omitempty"`
	SkipIfLabeled []string             `mapstructure:"skip-if-labeled,omitempty"`
	Files         FilesCondition       `mapstructure:"files,omitempty"`
	Title         TitleCondition       `mapstructure:"title,omitempty"`
	Description   DescriptionCondition `mapstructure:"description,omitempty"`
	Ref           RefCondition         `mapstructure:"ref,omitempty"`
	Comments      CommentsCondition    `mapstructure:"comments,omitempty"`
}

func (c *Condition) checkIfLabeled(meta EventData) bool {
	if len(c.IfLabeled) == 0 {
		return false
	} else {
		for _, check := range c.IfLabeled {
			for _, label := range meta.GetLabels() {
				if check == label {
					log.DebugWith(
						log.MetaFields{log.F("issue", meta.GetShortName())},
						"Matched Condition with label '%s'", check)
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
		c.Description.IsEmpty() &&
		c.Ref.IsEmpty() &&
		c.Comments.IsEmpty()
	if empty {
		log.DebugWith(
			log.MetaFields{log.F("issue", meta.GetShortName())},
			"Condition is empty")
	}
	return empty
}

func (c *Condition) Match(meta EventData) bool {
	match := c.checkAllEmpty(meta) ||
		c.checkIfLabeled(meta) ||
		c.Title.Match(meta) ||
		c.Description.Match(meta) ||
		c.Files.Match(meta) ||
		c.Ref.Match(meta) ||
		c.Comments.Match(meta)

	if match {
		for _, check := range c.SkipIfLabeled {
			for _, label := range meta.GetLabels() {
				if check == label {
					log.DebugWith(
						log.MetaFields{log.F("issue", meta.GetShortName())},
						"Skipping Condition with label '%s'", check)
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
			log.ErrorWith(log.MetaFields{log.E(e)}, "Unable to unmarshal rule condition")
		}
	}
	return result
}
