package engine

import (
	"strings"

	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util/log"
	"github.com/waigani/diffparser"
)

type PatchCondition struct {
	Hunk struct {
		StartsAt    int      `mapstructure:"starts-at,omitempty"`
		Patterns    []string `mapstructure:"patterns,omitempty"`
		NotPatterns []string `mapstructure:"not-patterns,omitempty"`
	} `mapstructure:"hunk,omitempty"`
}

func (c *PatchCondition) IsEmpty() bool {
	return c.Hunk.StartsAt == 0 && len(c.Hunk.Patterns) == 0
}

func (c *PatchCondition) getAllHunks(meta types.Data) []string {
	extended, ok := meta.(types.ExtendedData)
	if !ok {
		lc.DebugWith(
			log.MetaFields{log.F("condition", "PatchCondition"), log.F("issue", meta.GetShortName())},
			"Meta data is not extended data")
		return []string{}
	}
	mergedLines := make([]string, 0)
	for filename, patch := range extended.GetPatch() {
		if patch == nil || *patch == "" {
			lc.DebugWith(log.MetaFields{
				log.F("filename", filename),
				log.F("condition", "PatchCondition"),
				log.F("issue", meta.GetShortName()),
			}, "Empty patch, file moved or removed")
			continue
		}
		diff, err := diffparser.Parse(*patch)
		if err != nil {
			lc.ErrorWith(
				log.MetaFields{
					log.E(err),
					log.F("condition", "PatchCondition"),
					log.F("issue", meta.GetShortName()),
					log.F("filename", filename),
				}, "Failed to read patch")
			continue
		}
		for _, file := range diff.Files {
			lc.DebugWith(log.MetaFields{
				log.F("filename", filename),
				log.F("condition", "PatchCondition"),
				log.F("issue", meta.GetShortName()),
				log.F("original", file.OrigName),
				log.F("new", file.NewName),
				log.F("hunks", len(file.Hunks)),
			}, "Diff file")
			for _, hunk := range file.Hunks {
				if c.Hunk.StartsAt > 0 && hunk.NewRange.Start != c.Hunk.StartsAt {
					continue
				}
				mergedLines = append(mergedLines, mergeLines(hunk.NewRange.Lines))
			}
		}
	}
	return mergedLines
}

func (c *PatchCondition) Match(meta types.Data) bool {
	if c.IsEmpty() {
		return false
	}
	patchConditionCounter.Inc()
	hunks := c.getAllHunks(meta)
	return matchPatterns("PatchCondition", c.Hunk.Patterns, meta, func(types.Data) []string { return hunks }) ||
		matchNotPatterns("PatchCondition", c.Hunk.NotPatterns, meta, func(types.Data) []string { return hunks })
}
func mergeLines(lines []*diffparser.DiffLine) string {
	result := make([]string, len(lines), len(lines))
	for i, line := range lines {
		result[i] = line.Content

	}
	return strings.Join(result, "\n")
}
