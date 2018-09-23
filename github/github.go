package github

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/yakshaving.art/propaganda/core"
	"regexp"
	"strings"
)

// Headers
// "X-Hub-Signature":[]string{"sha1=be490c94029284a1074f6ed7d6f551affcfa6e8b"},
// "User-Agent":[]string{"GitHub-Hookshot/32d792e"},
// "Content-Length":[]string{"20774"},
// "X-Github-Delivery":[]string{"97757e90-bb8e-11e8-8017-464837e8ed07"},
// "X-Github-Event":[]string{"pull_request"},
// "Accept-Encoding":[]string{"gzip"},
// "Accept":[]string{"*/*"},
// "Content-Type":[]string{"application/json"}

// Parser implements the core.Parser type for GitHub Pull Request
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
	if event, ok := headers["X-Github-Event"]; ok {
		if len(event) != 1 {
			return false
		}
		return event[0] == "pull_request"
	}
	return false
}

// Parse parses a payload and returns a a valid one if everything is in place for it to be announced
func (p Parser) Parse(payload []byte) (core.Announcement, error) {
	var pl Payload
	if err := json.Unmarshal(payload, &pl); err != nil {
		return pl, fmt.Errorf("could not parse json payload: %s", err)
	}

	if !p.matcher.MatchString(pl.PullRequest.Title) {
		return pl, fmt.Errorf("MR title '%s' is not annouceable", pl.PullRequest.Title)
	}

	pl.PullRequest.Title = strings.TrimSpace(p.matcher.ReplaceAllString(pl.PullRequest.Title, ""))

	return pl, nil
}

// Payload wraps a Github pull request
type Payload struct {
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}

// PullRequest implements a pull request payload object
type PullRequest struct {
	URL    string `json:"html_url"`
	State  string `json:"state"`
	Title  string `json:"title"`
	Merged bool   `json:"merged"`
	Body   string `json:"body"`
}

// ProjectName implements Annoucement
func (pl Payload) ProjectName() string {
	return pl.Repository.FullName
}

// ShouldAnnounce implements Annoucement
func (pl Payload) ShouldAnnounce() bool {
	return pl.PullRequest.Merged && pl.PullRequest.State == "closed"
}

// Text implements Annoucement
func (pl Payload) Text() string {
	return fmt.Sprintf("*%s*\n\n%s\n\n*URL:* %s",
		pl.PullRequest.Title,
		pl.PullRequest.Body,
		pl.PullRequest.URL)
}

// Repository holds the repository information
type Repository struct {
	FullName string `json:"full_name"`
}
