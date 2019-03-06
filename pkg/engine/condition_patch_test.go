package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatchFirstLine(t *testing.T) {
	var rule rule
	rule.condition.Patch.Hunk.StartsAt = 1
	rule.condition.Patch.Hunk.Patterns = []string{"first"}
	patch :=
		`@@ -0,0 +1,1 @@
+This is the Copyright line`
	meta := &mockData{Patch: map[string]*string{
		"file1.txt": &patch,
	}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
	rule.condition.Patch.Hunk.Patterns = []string{"Copyright"}
	assert.True(t, rule.Accept(meta), "should match")
}

func TestPatchMatchingHunkStartAt(t *testing.T) {
	var rule rule
	rule.condition.Patch.Hunk.StartsAt = 5
	rule.condition.Patch.Hunk.Patterns = []string{"Copyright"}
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

func TestPatchMatchingHunkStartAtWithExtension(t *testing.T) {
	var rule rule
	rule.condition.Patch.Hunk.StartsAt = 1
	rule.condition.Patch.Hunk.Patterns = []string{"Copyright"}
	rule.condition.MatchKind = "any"
	rule.condition.Files.Extensions = []string{".scala"}
	patch :=
		`@@ -0,0 +1,1 @@
+This is the Copyright line`
	meta := &mockData{
		Patch: map[string]*string{
			"file1.txt": &patch,
		},
		FileExtensions: []string{".go"}}
	assert.True(t, rule.Accept(meta), "any should match")
	rule.condition.MatchKind = "all"
	assert.False(t, rule.Accept(meta), "all shouldn't match")
	rule.condition.Files.Extensions = []string{".scala", ".go"}
	assert.True(t, rule.Accept(meta), "ext should match")
}

func TestPatchNotMatchingHunkStartAtWithExtension(t *testing.T) {
	var rule rule
	rule.condition.Patch.Hunk.StartsAt = 1
	rule.condition.Patch.Hunk.NotPatterns = []string{"foofoo"}
	rule.condition.MatchKind = "any"
	rule.condition.Files.Extensions = []string{".scala"}
	patch :=
		`@@ -0,0 +1,1 @@
+This is the Copyright line`
	meta := &mockData{
		Patch: map[string]*string{
			"file1.txt": &patch,
		},
		FileExtensions: []string{".go"}}
	assert.True(t, rule.Accept(meta), "any should match")
	rule.condition.MatchKind = "all"
	assert.False(t, rule.Accept(meta), "all shouldn't match")
	rule.condition.Files.Extensions = []string{".scala", ".go"}
	assert.True(t, rule.Accept(meta), "ext should match")
	rule.condition.Patch.Hunk.NotPatterns = []string{"Copyright"}
	assert.False(t, rule.Accept(meta), "not shouldn't match")
}

func TestPatchAnyHunk(t *testing.T) {
	var rule rule
	rule.condition.Patch.Hunk.Patterns = []string{"foofoo"}
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
	rule.condition.Patch.Hunk.Patterns = []string{"blanks"}
	assert.True(t, rule.Accept(meta), "should match")
}

func TestPatchNoPatch(t *testing.T) {
	var rule rule
	rule.condition.Patch.Hunk.StartsAt = 5
	rule.condition.Patch.Hunk.Patterns = []string{"Copyright"}
	meta := &mockData{Patch: map[string]*string{
		"file1.txt": nil,
	}}
	assert.False(t, rule.Accept(meta), "shouldn't match")
}
