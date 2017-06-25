package commenter

import (
	"github.com/bivas/rivi/bot/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommentNoComments(t *testing.T) {
	action := action{rule: &rule{Comment: "comment1"}}
	meta := &mock.MockEventData{Comments: []string{}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedComments, 1, "added comments")
	assert.Equal(t, meta.Comments, []string{"comment1"}, "comments")
}

func TestNewCommentWithExisting(t *testing.T) {
	action := action{rule: &rule{Comment: "comment1"}}
	meta := &mock.MockEventData{Comments: []string{"comment2"}}
	config := &mock.MockConfiguration{}
	action.Apply(config, meta)
	assert.Len(t, meta.AddedComments, 1, "added Comments")
	assert.Len(t, meta.Comments, 2, "Comments")
}
