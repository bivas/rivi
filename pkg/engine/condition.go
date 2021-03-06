package engine

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bivas/rivi/pkg/types"
	"github.com/bivas/rivi/pkg/util/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

var lc = log.Get("engine.condition")

type sectionCondition interface {
	IsEmpty() bool
	Match(meta types.Data) bool
}

func matchNotPatterns(section string, patterns []string, meta types.Data, dataAccessor func(types.Data) []string) bool {
	if len(patterns) == 0 {
		return false
	} else {
		compiled := make([]*regexp.Regexp, 0)
		for _, pattern := range patterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				lc.WarningWith(
					log.MetaFields{
						log.F("condition", section),
						log.F("issue", meta.GetShortName()),
						log.E(err),
						log.F("regex", pattern),
					},
					"Unable to compile a not regex")
				continue
			}
			compiled = append(compiled, re)
		}
		if len(compiled) == 0 {
			lc.ErrorWith(
				log.MetaFields{log.F("condition", section), log.F("issue", meta.GetShortName())},
				"All configured not-patterns have failed to compile")
			return false
		}
		for _, check := range dataAccessor(meta) {
			for _, reg := range compiled {
				if !reg.MatchString(check) {
					lc.DebugWith(
						log.MetaFields{
							log.F("condition", section),
							log.F("issue", meta.GetShortName()),
							log.F("pattern", reg.String()),
							log.F("check", check),
						},
						"Matched with not-pattern")
					return true
				}
			}
		}
	}
	return false
}

func matchPatterns(section string, patterns []string, meta types.Data, dataAccessor func(types.Data) []string) bool {
	if len(patterns) == 0 {
		return false
	} else {
		compiled := make([]*regexp.Regexp, 0)
		for _, pattern := range patterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				lc.WarningWith(
					log.MetaFields{
						log.F("condition", section),
						log.F("issue", meta.GetShortName()),
						log.E(err),
						log.F("regex", pattern),
					},
					"Unable to compile regex")
				continue
			}
			compiled = append(compiled, re)
		}
		if len(compiled) == 0 {
			lc.ErrorWith(
				log.MetaFields{log.F("condition", section), log.F("issue", meta.GetShortName())},
				"All configured patterns have failed to compile")
			return false
		}
		for _, check := range dataAccessor(meta) {
			for _, reg := range compiled {
				if reg.MatchString(check) {
					lc.DebugWith(
						log.MetaFields{
							log.F("condition", section),
							log.F("issue", meta.GetShortName()),
							log.F("pattern", reg.String()),
							log.F("check", check),
						},
						"Matched with pattern")
					return true
				}
			}
		}
	}
	return false
}

type FilesCondition struct {
	Patterns   []string `mapstructure:"patterns,omitempty"`
	Extensions []string `mapstructure:"extensions,omitempty"`
}

func (c *FilesCondition) IsEmpty() bool {
	return len(c.Patterns) == 0 && len(c.Extensions) == 0
}

func (c *FilesCondition) checkPattern(meta types.Data) bool {
	return matchPatterns("FileCondition", c.Patterns, meta, func(types.Data) []string {
		return meta.GetFileNames()
	})
}

func (c *FilesCondition) checkExt(meta types.Data) bool {
	if len(c.Extensions) == 0 {
		return false
	} else {
		for _, check := range meta.GetFileExtensions() {
			for _, ext := range c.Extensions {
				if ext == check {
					lc.DebugWith(
						log.MetaFields{log.F("condition", "FileCondition"), log.F("issue", meta.GetShortName())},
						"Matched with extension '%s' on file '%s'", ext, check)
					return true
				}
			}
		}
	}
	return false
}

func (c *FilesCondition) Match(meta types.Data) bool {
	fileConditionCounter.Inc()
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

func (c *TitleCondition) Match(meta types.Data) bool {
	titleConditionCounter.Inc()
	title := meta.GetTitle()
	if c.StartsWith != "" && strings.HasPrefix(title, c.StartsWith) {
		lc.DebugWith(
			log.MetaFields{log.F("condition", "TitleCondition"), log.F("issue", meta.GetShortName())},
			"Matched with prefix '%s' on title '%s'", c.StartsWith, title)
		return true
	}
	if c.EndsWith != "" && strings.HasSuffix(title, c.EndsWith) {
		lc.DebugWith(
			log.MetaFields{log.F("condition", "TitleCondition"), log.F("issue", meta.GetShortName())},
			"Matched with suffix '%s' on title '%s'", c.EndsWith, title)
		return true
	}
	return matchPatterns("TitleCondition", c.Patterns, meta, func(types.Data) []string {
		return []string{title}
	})
}

type DescriptionCondition struct {
	StartsWith string   `mapstructure:"starts-with,omitempty"`
	EndsWith   string   `mapstructure:"ends-with,omitempty"`
	Patterns   []string `mapstructure:"patterns,omitempty"`
}

func (c *DescriptionCondition) IsEmpty() bool {
	return c.StartsWith == "" && c.EndsWith == "" && len(c.Patterns) == 0
}

func (c *DescriptionCondition) Match(meta types.Data) bool {
	descriptionConditionCounter.Inc()
	description := meta.GetDescription()
	if c.StartsWith != "" && strings.HasPrefix(description, c.StartsWith) {
		lc.DebugWith(
			log.MetaFields{log.F("condition", "DescriptionCondition"), log.F("issue", meta.GetShortName())},
			"Matched with prefix '%s' on description '%s'",
			c.StartsWith,
			description)
		return true
	}
	if c.EndsWith != "" && strings.HasSuffix(description, c.EndsWith) {
		lc.DebugWith(
			log.MetaFields{log.F("condition", "DescriptionCondition"), log.F("issue", meta.GetShortName())},
			"Matched with suffix '%s' on description '%s'",
			c.EndsWith,
			description)
		return true
	}
	return matchPatterns("DescriptionCondition", c.Patterns, meta, func(types.Data) []string {
		return []string{description}
	})
}

type RefCondition struct {
	Equals   string   `mapstructure:"match,omitempty"`
	Patterns []string `mapstructure:"patterns,omitempty"`
}

func (c *RefCondition) IsEmpty() bool {
	return c.Equals == "" && len(c.Patterns) == 0
}

func (c *RefCondition) Match(meta types.Data) bool {
	refConditionCounter.Inc()
	ref := meta.GetRef()
	if c.Equals != "" && ref == c.Equals {
		lc.DebugWith(
			log.MetaFields{log.F("condition", "RefCondition"), log.F("issue", meta.GetShortName())},
			"Matched RefCondition with match on ref '%s'", ref)
		return true
	}
	return matchPatterns("RefCondition", c.Patterns, meta, func(types.Data) []string {
		return []string{ref}
	})
}

var commentsRegex = silentCompile("^([><]{1}[=]?|[=]{2}|)[ ]*([0-9]+)[ ]*$")

func silentCompile(pattern string) *regexp.Regexp {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		lc.ErrorWith(
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

func (c *CommentsCondition) Match(meta types.Data) bool {
	if c.IsEmpty() {
		return false
	}
	commentsConditionCounter.Inc()
	if commentsRegex == nil {
		lc.WarningWith(
			log.MetaFields{log.F("condition", "CommentsCondition"),
				log.F("issue", meta.GetShortName())},
			"comments regex is nil")
		return false
	}
	count := int64(len(meta.GetComments()))
	countGroups := commentsRegex.FindStringSubmatch(c.Count)
	if len(countGroups) != 3 {
		lc.WarningWith(
			log.MetaFields{log.F("condition", "CommentsCondition"),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups)},
			"No groups matched'")
		return false
	}
	query, err := strconv.ParseInt(countGroups[2], 10, 64)
	if err != nil {
		lc.WarningWith(
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
		lc.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.F("matcher", op),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups),
				log.F("matched", matched), log.F("comments.count", count)},
			"matching with '=='")
		break
	case op == ">":
		// gt
		matched = count > query
		lc.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.F("matcher", op),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups),
				log.F("matched", matched), log.F("comments.count", count)},
			"matching with '>'")
	case op == ">=":
		// gte
		matched = count >= query
		lc.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.F("matcher", op),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups),
				log.F("matched", matched), log.F("comments.count", count)},
			"matching with '>='")
	case op == "<":
		// lt
		matched = count < query
		lc.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.F("matcher", op),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups),
				log.F("matched", matched), log.F("comments.count", count)},
			"matching with '<'")
	case op == "<=":
		// lte
		matched = count <= query
		lc.DebugWith(
			log.MetaFields{log.F("condition", "CommentsCondition"), log.F("matcher", op),
				log.F("issue", meta.GetShortName()), log.F("groups", countGroups),
				log.F("matched", matched), log.F("comments.count", count)},
			"matching with '<='")
	}
	return matched
}

type Condition struct {
	Order         int                  `mapstructure:"order,omitempty"`
	MatchKind     string               `mapstructure:"match-kind,omitempty"`
	IfLabeled     []string             `mapstructure:"if-labeled,omitempty"`
	SkipIfLabeled []string             `mapstructure:"skip-if-labeled,omitempty"`
	Files         FilesCondition       `mapstructure:"files,omitempty"`
	Title         TitleCondition       `mapstructure:"title,omitempty"`
	Description   DescriptionCondition `mapstructure:"description,omitempty"`
	Ref           RefCondition         `mapstructure:"ref,omitempty"`
	Comments      CommentsCondition    `mapstructure:"comments,omitempty"`
	Patch         PatchCondition       `mapstructure:"patch,omitempty"`
}

func (c *Condition) checkIfLabeled(meta types.Data) bool {
	if len(c.IfLabeled) == 0 {
		return false
	} else {
		labelConditionCounter.Inc()
		for _, check := range c.IfLabeled {
			for _, label := range meta.GetLabels() {
				if check == label {
					lc.DebugWith(
						log.MetaFields{log.F("issue", meta.GetShortName())},
						"Matched Condition with label '%s'", check)
					return true
				}
			}
		}
	}
	return false
}

func (c *Condition) checkAllEmpty(meta types.Data) bool {
	empty := len(c.IfLabeled) == 0 &&
		c.Files.IsEmpty() &&
		c.Title.IsEmpty() &&
		c.Description.IsEmpty() &&
		c.Ref.IsEmpty() &&
		c.Comments.IsEmpty() &&
		c.Patch.IsEmpty()
	if empty {
		lc.DebugWith(
			log.MetaFields{log.F("issue", meta.GetShortName())},
			"Condition is empty")
		emptyConditionCounter.Inc()
	}
	return empty
}

func (c *Condition) Match(meta types.Data) bool {
	sections := []sectionCondition{
		&c.Title,
		&c.Description,
		&c.Files,
		&c.Ref,
		&c.Comments,
		&c.Patch,
	}
	match := false
	switch strings.ToLower(c.MatchKind) {
	case "any", "":
		match = c.checkAllEmpty(meta) || c.checkIfLabeled(meta)
		if !match {
			for _, section := range sections {
				match = match || section.Match(meta)
			}
		}
	case "all":
		match = c.checkAllEmpty(meta)
		if !match {
			match = len(c.IfLabeled) == 0 || c.checkIfLabeled(meta)
			for _, section := range sections {
				match = match && (section.IsEmpty() || section.Match(meta))
			}
		}
	default:
		lc.WarningWith(
			log.MetaFields{log.F("match-kind", c.MatchKind), log.F("issue", meta.GetShortName())},
			"Unknown match kind")
	}

	if match && len(c.SkipIfLabeled) > 0 {
		blockByLabelConditionCounter.Inc()
		for _, check := range c.SkipIfLabeled {
			for _, label := range meta.GetLabels() {
				if check == label {
					lc.DebugWith(
						log.MetaFields{log.F("issue", meta.GetShortName())},
						"Skipping Condition with label '%s'", check)
					return false
				}
			}
		}
	}
	return match
}

func createCounter(name string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "rivi",
		Subsystem: "condition",
		Name:      name,
		Help:      fmt.Sprintf("Condition counter for %s", name),
	})
}

var (
	emptyConditionCounter        = createCounter("empty")
	labelConditionCounter        = createCounter("label")
	blockByLabelConditionCounter = createCounter("not_labeled")
	commentsConditionCounter     = createCounter("comments")
	descriptionConditionCounter  = createCounter("description")
	fileConditionCounter         = createCounter("files")
	refConditionCounter          = createCounter("ref")
	titleConditionCounter        = createCounter("title")
	patchConditionCounter        = createCounter("patch")
)

func buildConditionFromConfiguration(config *viper.Viper) Condition {
	var result Condition
	exists := config.Get("condition")
	if exists != nil {
		condition := config.Sub("condition")
		if e := condition.Unmarshal(&result); e != nil {
			lc.ErrorWith(log.MetaFields{log.E(e)}, "Unable to unmarshal rule condition")
		}
	}
	return result
}

func init() {
	prometheus.Register(emptyConditionCounter)
	prometheus.Register(labelConditionCounter)
	prometheus.Register(blockByLabelConditionCounter)
	prometheus.Register(commentsConditionCounter)
	prometheus.Register(descriptionConditionCounter)
	prometheus.Register(fileConditionCounter)
	prometheus.Register(refConditionCounter)
	prometheus.Register(titleConditionCounter)
	prometheus.Register(patchConditionCounter)
}
