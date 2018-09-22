package slack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gitlab.com/yakshaving.art/propaganda/core"
	"gitlab.com/yakshaving.art/propaganda/metrics"

	"github.com/sirupsen/logrus"
)

// Announcer announces things to slack
type Announcer struct {
	WebhookURL string
	Proxy      string
	// Channel    string
}

// Announce implements core.Announcer interface
func (a Announcer) Announce(announcement core.Announcement) {

	body, err := json.Marshal(payload{
		Markdown: true,
		Text:     announcement.Text(),
	})
	if err != nil {
		metrics.AnnouncementErrors.WithLabelValues("encoding").Inc()
		logrus.Errorf("failed to encode payload as json: %s", err)
		return
	}

	logrus.Debugf("posting payload: %s", string(body))

	req, err := http.NewRequest(http.MethodPost, a.WebhookURL, bytes.NewReader(body))
	if err != nil {
		metrics.AnnouncementErrors.WithLabelValues("request").Inc()
		logrus.Errorf("failed to create POST request: %s", err)
		return
	}
	req.Header.Add("Content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		metrics.AnnouncementErrors.WithLabelValues("unknown").Inc()
		logrus.Errorf("failed to call slack webhook: %s", err)
		return
	}

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		logrus.Debugf("payload pushed to slack with response: %#v", resp)
		metrics.AnnouncementSuccesses.WithLabelValues(announcement.ProjectName()).Inc()

	default:
		logrus.Debugf("payload failed to slack with response: %#v", resp)
		metrics.AnnouncementErrors.WithLabelValues(resp.Status).Inc()
	}
}

type payload struct {
	// Username string `json:"username,omitempty"`
	Text     string `json:"text,omitempty"`
	Markdown bool   `json:"mrkdwn,omitempty"`
	// Parse       string       `json:"parse,omitempty"`
	// IconUrl     string       `json:"icon_url,omitempty"`
	// IconEmoji   string       `json:"icon_emoji,omitempty"`
	// Channel string `json:"channel,omitempty"`
	// LinkNames   string       `json:"link_names,omitempty"`
	// Attachments []Attachment `json:"attachments,omitempty"`
	// UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	// UnfurlMedia bool         `json:"unfurl_media,omitempty"`
}

// type Field struct {
// 	Title string `json:"title"`
// 	Value string `json:"value"`
// 	Short bool   `json:"short"`
// }

// type Attachment struct {
// 	Fallback   *string   `json:"fallback"`
// 	Color      *string   `json:"color"`
// 	PreText    *string   `json:"pretext"`
// 	AuthorName *string   `json:"author_name"`
// 	AuthorLink *string   `json:"author_link"`
// 	AuthorIcon *string   `json:"author_icon"`
// 	Title      *string   `json:"title"`
// 	TitleLink  *string   `json:"title_link"`
// 	Text       *string   `json:"text"`
// 	ImageUrl   *string   `json:"image_url"`
// 	Fields     []*Field  `json:"fields"`
// 	Footer     *string   `json:"footer"`
// 	FooterIcon *string   `json:"footer_icon"`
// 	Timestamp  *int64    `json:"ts"`
// 	MarkdownIn *[]string `json:"mrkdwn_in"`
// }
