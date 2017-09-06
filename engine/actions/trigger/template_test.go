package trigger

import (
	"bytes"
	"testing"
	"text/template"
	"time"

	"github.com/bivas/rivi/types"
	"github.com/bivas/rivi/util"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	temp, e := template.New("message").Parse(defaultTemplateBody)
	if e != nil {
		t.Fatalf("Template has failed. %s", e)
	}
	stamp, err := time.Parse("2006-01-02T15:04:05", "2017-10-28T10:13:44")
	if err != nil {
		t.Fatalf("Unable to get time. %s", err)
	}
	message := message{
		Time:   stamp,
		Number: 1,
		Title:  "title1",
		State:  "open",
		Owner:  "my",
		Repo:   "repo",
		Origin: types.Origin{User: "self"}}
	var buf bytes.Buffer
	temp.Execute(&buf, message)
	expected :=
		util.StripNonSpaceWhitespaces(`{
			"time":"2017-10-28 10:13:44 +0000 UTC",
			"message":"Triggered by Rivi",
			"item":{
				"repository":"my/repo",
				"state":"open",
				"id":1,
				"title":"title1"
			}
		}`)
	assert.Equal(t, expected, buf.String(), "json")
}
