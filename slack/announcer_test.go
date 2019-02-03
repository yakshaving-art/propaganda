package slack_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/yakshaving.art/propaganda/configuration"
	"gitlab.com/yakshaving.art/propaganda/slack"

	"github.com/stretchr/testify/assert"
)

type announcement struct {
	project string
	text    string
}

func (a announcement) ProjectName() string {
	return a.project
}

func (a announcement) Text() string {
	return a.text
}

func (a announcement) ShouldAnnounce() bool {
	return true
}

func TestSlackAnnouncerCanSucceed(t *testing.T) {
	configuration.Load([]byte("default_channel: general"))

	a := announcement{
		text: `test [text](http://endpoint)
		[one](https://secondlink.md) [link](./thirdlink.md)
		[]()
		[](#anchor)
		[no link]()
		[link](http://something/])
		[test/><javascript alert("test";) /><!----](blah)`,
		project: "some/project",
	}
	ass := assert.New(t)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		ass.NoError(err)
		defer r.Body.Close()

		data := make(map[string]interface{}, 0)
		if err = json.Unmarshal(b, &data); err != nil {
			http.Error(w, err.Error(), 500)
		}
		if ass.Equal(`test <http://endpoint|text>
		<https://secondlink.md|one> <./thirdlink.md|link>
		[]()
		[](#anchor)
		[no link]()
		<http://something/]|link>
		<blah|test/><javascript alert("test";) /><!---->`,
			data["text"]) {
			w.WriteHeader(200)
		} else {
			http.Error(w, fmt.Sprintf("text is not as expected: %s", data["text"]), 400)
		}
	}))
	defer s.Close()

	announcer := slack.Announcer{
		WebhookURL: s.URL,
	}

	ass.NoError(announcer.Announce(a))
}

func TestSlackAnnouncerCanFail(t *testing.T) {
	ass := assert.New(t)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid payload", 400)
	}))
	defer s.Close()

	announcer := slack.Announcer{
		WebhookURL: s.URL,
	}
	ass.Errorf(announcer.Announce(announcement{
		text:    "invalid test text",
		project: "some/project",
	}), "failed to push payload to slack with code 400: invalid payload")
}
