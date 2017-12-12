package engine

import (
	"testing"

	"github.com/bivas/rivi/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMatchLabel(t *testing.T) {
	var rule1, rule2 rule
	rule1.condition.IfLabeled = []string{"label1"}
	rule1.name = "rule1"
	rule2.condition.SkipIfLabeled = []string{"pending-approval"}
	rule2.name = "rule2"
	meta := &mockData{Labels: []string{"label1", "pending-approval"}}
	matched := make([]string, 0)
	for _, rule := range []rule{rule1, rule2} {
		if rule.Accept(meta) {
			matched = append(matched, rule.Name())
		}
	}
	assert.Contains(t, matched, "rule1", "matched")
	assert.NotContains(t, matched, "rule2", "matched")
}

func TestMatchNoLabel(t *testing.T) {
	var rule1 rule
	rule1.condition.IfLabeled = []string{"label1"}
	meta := &mockData{Labels: []string{"pending-approval"}}
	result := rule1.condition.Match(meta)
	assert.False(t, result, "no label to match")
}

func TestSkipLabel(t *testing.T) {
	var tested rule
	tested.condition.SkipIfLabeled = []string{"pending-approval"}
	tested.name = "rule"
	meta := &mockData{Labels: []string{"pending-approval"}}
	matched := make([]string, 0)
	for _, r := range []rule{tested} {
		if r.Accept(meta) {
			matched = append(matched, r.Name())
		}
	}
	assert.Len(t, matched, 0, "matched")
}

func TestMatchPattern(t *testing.T) {
	var tested rule
	tested.condition.Files.Patterns = []string{".*/foo.txt"}
	tested.name = "rule4"
	meta := &mockData{
		Labels: []string{"pending-approval"},
		FileNames: []string{
			"foo.txt",
			"path/to/docs/foo.txt",
		}}
	matched := make([]string, 0)
	for _, r := range []rule{tested} {
		if r.Accept(meta) {
			matched = append(matched, r.Name())
		}
	}
	assert.Len(t, matched, 1, "matched")
	assert.Contains(t, matched, "rule4", "matched")
}

func TestMatchExt(t *testing.T) {
	var tested rule
	tested.condition.Files.Extensions = []string{".go"}
	tested.name = "rule4"
	meta := &mockData{FileExtensions: []string{".scala", ".go"}}
	matched := make([]string, 0)
	for _, r := range []rule{tested} {
		if r.Accept(meta) {
			matched = append(matched, r.Name())
		}
	}
	assert.Len(t, matched, 1, "matched")
	assert.Contains(t, matched, "rule4", "matched")
}

func TestMatchNoExt(t *testing.T) {
	var tested rule
	tested.condition.Files.Extensions = []string{".go"}
	meta := &mockData{FileExtensions: []string{".scala"}}
	result := tested.condition.Files.Match(meta)
	assert.False(t, result, "nothing to match")
}

func TestTitleStartsWith(t *testing.T) {
	var rule rule
	rule.condition.Title.StartsWith = "BUGFIX"
	meta := &mockData{Title: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Title = "BUGFIX it"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestTitleEndsWith(t *testing.T) {
	var rule rule
	rule.condition.Title.EndsWith = "WIP"
	meta := &mockData{Title: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Title = "BUGFIX WIP"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestTitlePattern(t *testing.T) {
	var rule rule
	rule.condition.Title.Patterns = []string{".* Bug( )?[0-9]{5} .*"}
	meta := &mockData{Title: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Title = "This PR for Bug1 with comment"
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Title = "This PR for Bug 45456 with comment"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestDescriptionStartsWith(t *testing.T) {
	var rule rule
	rule.condition.Description.StartsWith = "BUGFIX"
	meta := &mockData{Description: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Description = "BUGFIX it"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestDescriptionEndsWith(t *testing.T) {
	var rule rule
	rule.condition.Description.EndsWith = "WIP"
	meta := &mockData{Description: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Description = "BUGFIX WIP"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestDescriptionPattern(t *testing.T) {
	var rule rule
	rule.condition.Description.Patterns = []string{"(?s)~~~.*deps:"}
	meta := &mockData{Description: "NOT A BUGFIX"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Description = "~~~\n     test_priorities"
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Description = "~~~\n    deps:\nplenty of dependencies"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestMatchEmptyCondition(t *testing.T) {
	meta := &mockData{}
	rule := rule{condition: Condition{}}
	assert.True(t, rule.Accept(meta), "none")
}

func TestMatchRef(t *testing.T) {
	var rule rule
	rule.condition.Ref.Equals = "master"
	meta := &mockData{Ref: "development"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Ref = "master"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestRefPatters(t *testing.T) {
	var rule rule
	rule.condition.Ref.Patterns = []string{"integration-v[0-9]{2}$"}
	meta := &mockData{Ref: "development"}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	meta.Ref = "integration-v11"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestCommentsCountNoOp(t *testing.T) {
	var rule rule
	rule.condition.Comments.Count = "5"
	meta := &mockData{Comments: []types.Comment{{Commenter: "user1", Comment: "comment1"}}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	rule.condition.Comments.Count = "1"
	assert.True(t, rule.Accept(meta), "shouldn't match")
}

func TestCommentsCountEquals(t *testing.T) {
	var rule rule
	rule.condition.Comments.Count = "==5"
	meta := &mockData{Comments: []types.Comment{{Commenter: "user1", Comment: "comment1"}}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	rule.condition.Comments.Count = "==1"
	assert.True(t, rule.Accept(meta), "shouldn't match")
}

func TestCommentsCountLessThan(t *testing.T) {
	var rule rule
	rule.condition.Comments.Count = "<1"
	meta := &mockData{Comments: []types.Comment{{Commenter: "user1", Comment: "comment1"}}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	rule.condition.Comments.Count = "<5"
	assert.True(t, rule.Accept(meta), "shouldn't match")
}

func TestCommentsCountLessThanEquals(t *testing.T) {
	var rule rule
	rule.condition.Comments.Count = "<=1"
	meta := &mockData{Comments: []types.Comment{
		{Commenter: "user1", Comment: "comment1"}, {Commenter: "user2", Comment: "comment2"}}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	rule.condition.Comments.Count = "<=5"
	assert.True(t, rule.Accept(meta), "shouldn't match")
}

func TestCommentsCountGreaterThan(t *testing.T) {
	var rule rule
	rule.condition.Comments.Count = ">5"
	meta := &mockData{Comments: []types.Comment{
		{Commenter: "user1", Comment: "comment1"}, {Commenter: "user2", Comment: "comment2"}}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	rule.condition.Comments.Count = ">1"
	assert.True(t, rule.Accept(meta), "shouldn't match")
}

func TestCommentsCountGreaterThanEquals(t *testing.T) {
	var rule rule
	rule.condition.Comments.Count = ">=5"
	meta := &mockData{Comments: []types.Comment{{Commenter: "user1", Comment: "comment1"}}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	rule.condition.Comments.Count = ">=1"
	assert.True(t, rule.Accept(meta), "shouldn't match")
}

func TestPatchFirstLine(t *testing.T) {
	var rule rule
	rule.condition.Patch.Hunk.StartsAt = 1
	rule.condition.Patch.Hunk.Pattern = "first"
	patch :=
		`@@ -0,0 +1,1 @@
+This is the Copyright line`
	meta := &mockData{Patch: map[string]*string{
		"file1.txt": &patch,
	}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	rule.condition.Patch.Hunk.Pattern = "Copyright"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestPatchNoMatchingHunk(t *testing.T) {
	var rule rule
	rule.condition.Patch.Hunk.StartsAt = 5
	rule.condition.Patch.Hunk.Pattern = "Copyright"
	patch :=
		`@@ -0,0 +1,1 @@
+This is the Copyright line`
	meta := &mockData{Patch: map[string]*string{
		"file1.txt": &patch,
	}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	rule.condition.Patch.Hunk.StartsAt = 1
	assert.True(t, rule.Accept(meta), "should match")
}

func TestPatchAnyHunk(t *testing.T) {
	var rule rule
	rule.condition.Patch.Hunk.Pattern = "foofoo"
	patch :=
		`@@ -0,0 +1,10 @@
+This is the first line
+This is the second line
+This is the third line
+This is the forth line
+This is the fifth line
+
+
+
+
+Last line after blanks`
	meta := &mockData{Patch: map[string]*string{
		"file1.txt": &patch,
	}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	rule.condition.Patch.Hunk.Pattern = "blanks"
	assert.True(t, rule.Accept(meta), "should match")
}

func TestPatchNoPatch(t *testing.T) {
	var rule rule
	rule.condition.Patch.Hunk.StartsAt = 5
	rule.condition.Patch.Hunk.Pattern = "Copyright"
	meta := &mockData{Patch: map[string]*string{
		"file1.txt": nil,
	}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
}

func TestCheckPatternNoPatterns(t *testing.T) {
	result := patternCheck("test", []string{}, &mockData{}, func(types.Data) []string {
		return nil
	})
	assert.False(t, result, "no pattern")
}

func TestCheckPatternNothingCompiles(t *testing.T) {
	result := patternCheck("test", []string{"[a"}, &mockData{}, func(types.Data) []string {
		return nil
	})
	assert.False(t, result, "no pattern compiled")
}

func TestBuildConditionFromEmptyConfiguration(t *testing.T) {
	result := buildConditionFromConfiguration(viper.New())
	assert.NotNil(t, result, "default condition")
	assert.True(t, result.checkAllEmpty(&mockData{}), "empty condition")
}
