package commenter

import (
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/bot/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialization(t *testing.T) {
	input := map[string]interface{}{
		"comment": "comment1",
	}

	var f factory
	result := f.BuildAction(input)
	assert.NotNil(t, result, "should create action")
	s, ok := result.(*action)
	assert.True(t, ok, "should be of this package")
	assert.Equal(t, "comment1", s.rule.Comment, "comment")
}

func TestCommentNoComments(t *testing.T) {
	action := action{rule: &rule{Comment: "comment1"}}
	meta := &mock.MockEventData{Comments: []bot.Comment{}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedComments, 1, "added comments")
	assert.Len(t, meta.Comments, 1, "comments")
}

func TestNewCommentWithExisting(t *testing.T) {
	action := action{rule: &rule{Comment: "comment1"}}
	meta := &mock.MockEventData{Comments: []bot.Comment{bot.Comment{Comment: "comment2"}}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedComments, 1, "added Comments")
	assert.Len(t, meta.Comments, 2, "Comments")
}
