package gitlab

// Headers:
// "Connection":[]string{"close"},
// "X-Forwarded-For":[]string{"51.15.62.67"},
// "X-Gitlab-Event":[]string{"Merge Request Hook"},
// "X-Gitlab-Token":[]string{"my-secret-token"},
// "Accept-Encoding":[]string{"gzip"},
// "User-Agent":[]string{"Go-http-client/1.1"},
// "Content-Length":[]string{"3963"},
// "Content-Type":[]string{"application/json"},

import (
	"fmt"
	"regexp"
	"strings"

	"encoding/json"

	"gitlab.com/yakshaving.art/propaganda/core"

	"github.com/sirupsen/logrus"
)

// Parser implements the core.Parser type for GitLab merge webhooks
type Parser struct {
	matcher *regexp.Regexp
}

// NewParser creates a new parser using the pattern provided
func NewParser(pattern string) Parser {
	re, err := regexp.Compile(pattern)
	if err != nil {
		logrus.Fatalf("could not compile regexp pattern for announcements: %s", err)
	}

	return Parser{
		matcher: re,
	}
}

// MatchHeaders indicates that the headers match with the kind of request
func (Parser) MatchHeaders(headers map[string][]string) bool {
	if event, ok := headers["X-Gitlab-Event"]; ok {
		if len(event) != 1 {
			return false
		}
		return event[0] == "Merge Request Hook"
	}
	return false
}

// Parse parses a payload and returns a a valid one if everything is in place for it to be announced
func (p Parser) Parse(payload []byte) (core.Announcement, error) {
	var mr MergeRequest
	if err := json.Unmarshal(payload, &mr); err != nil {
		return mr, fmt.Errorf("could not parse json payload: %s", err)
	}
	if mr.Kind != "merge_request" {
		return MergeRequest{}, fmt.Errorf("json payload is not a merge request but a %s", mr.Kind)
	}

	if !p.matcher.MatchString(mr.Attributes.Title) {
		return MergeRequest{}, fmt.Errorf("MR title '%s' is not annouceable", mr.Attributes.Title)
	}

	mr.Attributes.Title = strings.TrimSpace(p.matcher.ReplaceAllString(mr.Attributes.Title, ""))

	return mr, nil
}

// MergeRequest is the MR object
type MergeRequest struct {
	Kind       string     `json:"object_kind"`
	Project    Project    `json:"project"`
	Attributes Attributes `json:"object_attributes"`
}

// Text implements Annoucement
func (m MergeRequest) Text() string {
	return fmt.Sprintf("*%s*\n\n%s\n\n*URL:* %s",
		m.Attributes.Title,
		m.Attributes.Description,
		m.Attributes.URL)
}

// ShouldAnnounce implements Announcement
func (m MergeRequest) ShouldAnnounce() bool {
	return m.Attributes.State == "merged"
}

// ProjectName implements Announcement
func (m MergeRequest) ProjectName() string {
	return m.Project.PathWithNamespace
}

// Project is used to identify which project it is including the namespace
type Project struct {
	PathWithNamespace string `json:"path_with_namespace"`
}

// Attributes represent things like state, title, url or action
type Attributes struct {
	State       string `json:"state"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Action      string `json:"action"`
}
