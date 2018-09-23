package github_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/yakshaving.art/propaganda/github"
)

func TestParsingPayloads(t *testing.T) {
	parser := github.NewParser("\\[announce\\]")
	tt := []struct {
		name           string
		jsonFilename   string
		expected       github.Payload
		shouldAnnounce bool
	}{
		{
			"MR Create",
			"fixtures/github-pr-create.json",
			github.Payload{
				PullRequest: github.PullRequest{
					State:  "open",
					Title:  "Just write something",
					URL:    "https://github.com/pcarranza/testing-prs/pull/1",
					Body:   "some payload message that comes from the description",
					Merged: false,
				},
				Repository: github.Repository{
					FullName: "pcarranza/testing-prs",
				},
			},
			false,
		},
		{
			"MR Merged",
			"fixtures/github-pr-merged.json",
			github.Payload{
				PullRequest: github.PullRequest{
					State:  "closed",
					Title:  "Just write something",
					URL:    "https://github.com/pcarranza/testing-prs/pull/1",
					Body:   "some payload message that comes from the description",
					Merged: true,
				},
				Repository: github.Repository{
					FullName: "pcarranza/testing-prs",
				},
			},
			true,
		},
		{
			"MR Closed without a merge",
			"fixtures/github-pr-close-no-merge.json",
			github.Payload{
				PullRequest: github.PullRequest{
					State:  "closed",
					Title:  "Second test",
					URL:    "https://github.com/pcarranza/testing-prs/pull/2",
					Body:   "some other payload message that comes from the description",
					Merged: false,
				},
				Repository: github.Repository{
					FullName: "pcarranza/testing-prs",
				},
			},
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)
			b, err := ioutil.ReadFile(tc.jsonFilename)
			a.Nilf(err, "could not read fixture file %s", tc.jsonFilename)
			a.NotNilf(b, "content should not be nil")

			mr, err := parser.Parse(b)
			a.NoErrorf(err, "could not unmarshal PR json")

			a.EqualValuesf(tc.expected, mr, "parsed merge request is not as expected")

			a.Equal(tc.expected.Text(), mr.Text())
			a.Equal(tc.expected.ProjectName(), mr.ProjectName())
			a.Equal(tc.shouldAnnounce, mr.ShouldAnnounce())
		})
	}
}

func TestHeadersMatcher(t *testing.T) {
	a := assert.New(t)
	p := github.NewParser(".*")

	a.Equal(false, p.MatchHeaders(map[string][]string{}))
	a.Equal(false, p.MatchHeaders(map[string][]string{"X-Gitlab-Event": {"pull_request"}}))
	a.Equal(false, p.MatchHeaders(map[string][]string{"X-Gitlab-Event": {"Merge Request Hook"}}))
	a.Equal(false, p.MatchHeaders(map[string][]string{"X-Github-Event": {"Merge Request Hook"}}))
	a.Equal(true, p.MatchHeaders(map[string][]string{"X-Github-Event": {"pull_request"}}))
}
