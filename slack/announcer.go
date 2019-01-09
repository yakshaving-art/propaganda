package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	conf "gitlab.com/yakshaving.art/propaganda/configuration"
	"gitlab.com/yakshaving.art/propaganda/core"
	"gitlab.com/yakshaving.art/propaganda/metrics"

	"github.com/sirupsen/logrus"
)

// Announcer announces things to slack
type Announcer struct {
	WebhookURL string
	Proxy      string
}

// transforms markdown links to slack links
var re = regexp.MustCompile("\\[([^\\]])\\]\\(([^\\)]+)\\)")

// Announce implements core.Announcer interface
func (a Announcer) Announce(announcement core.Announcement) error {
	body, err := json.Marshal(payload{
		Markdown:    true,
		Text:        re.ReplaceAllString(announcement.Text(), "<$2|$1>"),
		Channel:     conf.GetConfiguration().GetChannel(announcement.ProjectName()),
		UnfurlLinks: false,
		UnfurlMedia: false,
	})
	if err != nil {
		metrics.AnnouncementErrors.WithLabelValues("encoding").Inc()
		return fmt.Errorf("failed to encode payload as json: %s", err)
	}

	logrus.Debugf("posting payload: %s", string(body))

	req, err := http.NewRequest(http.MethodPost, a.WebhookURL, bytes.NewReader(body))
	if err != nil {
		metrics.AnnouncementErrors.WithLabelValues("request").Inc()
		return fmt.Errorf("failed to create POST request: %s", err)
	}
	req.Header.Add("Content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		metrics.AnnouncementErrors.WithLabelValues("unknown").Inc()
		return fmt.Errorf("failed to call slack webhook: %s", err)
	}

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		logrus.Debugf("payload pushed to slack with response: %#v", resp)
		metrics.AnnouncementSuccesses.WithLabelValues(announcement.ProjectName()).Inc()

	default:
		metrics.AnnouncementErrors.WithLabelValues(resp.Status).Inc()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to push payload with code %d, and also to read the response body: %s", resp.StatusCode, err)
		}
		defer resp.Body.Close()
		return fmt.Errorf("failed to push payload to slack with code %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// Payload is the slack payload object used to send data
type payload struct {
	Text        string `json:"text,omitempty"`
	Markdown    bool   `json:"mrkdwn,omitempty"`
	Channel     string `json:"channel,omitempty"`
	UnfurlLinks bool   `json:"unfurl_links,omitempty"`
	UnfurlMedia bool   `json:"unfurl_media,omitempty"`
}
