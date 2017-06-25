package trigger

import (
	"bytes"
	"github.com/bivas/rivi/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"text/template"
	"time"
)

func TestTemplate(t *testing.T) {
	temp, e := template.New("message").Parse(defaultTemplate)
	if e != nil {
		t.Fatalf("Template has failed. %s", e)
	}
	stamp, err := time.Parse("2006-01-02T15:04:05", "2017-10-28T10:13:44")
	if err != nil {
		t.Fatalf("Unable to get time. %s", err)
	}
	message := message{Time: stamp, Number: 1, Title: "title1", State: "open", Owner: "my", Repo: "repo", Origin: "self"}
	var buffer bytes.Buffer
	temp.Execute(&buffer, message)
	expected :=
		util.StripNonSpaceWhitespaces(`{
			"time":"2017-10-28 10:13:44 +0000 UTC",
			"message":"Triggered by Rivi Bot",
			"item":{
				"repository":"my/repo",
				"state":"open",
				"id":1,
				"title":"title1"
			}
		}`)
	assert.Equal(t, expected, buffer.String(), "json")
}
