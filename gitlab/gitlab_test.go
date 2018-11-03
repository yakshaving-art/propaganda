package gitlab_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/yakshaving.art/propaganda/gitlab"
)

func TestParsingPayloads(t *testing.T) {
	parser := gitlab.NewParser("\\[announce\\]", "mytoken")
	tt := []struct {
		name           string
		token          string
		jsonFilename   string
		expected       gitlab.MergeRequest
		shouldAnnounce bool
	}{
		{
			"MR Create",
			"mytoken",
			"fixtures/gitlab-mr-create.json",
			gitlab.MergeRequest{
				Kind: "merge_request",
				Project: gitlab.Project{
					PathWithNamespace: "pablo/testing-webhooks",
				},
				Attributes: gitlab.Attributes{
					State:       "opened",
					Title:       "Update README.md",
					Description: "Something in the description",
					URL:         "https://git.yakshaving.art/pablo/testing-webhooks/merge_requests/1",
					Action:      "open",
				},
			},
			false,
		},
		{
			"MR Merged",
			"mytoken",
			"fixtures/gitlab-mr-merged.json",
			gitlab.MergeRequest{
				Kind: "merge_request",
				Project: gitlab.Project{
					PathWithNamespace: "pablo/testing-webhooks",
				},
				Attributes: gitlab.Attributes{
					State:       "merged",
					Title:       "Update README.md",
					Description: "Something in the description",
					URL:         "https://git.yakshaving.art/pablo/testing-webhooks/merge_requests/1",
					Action:      "merge",
				},
			},
			true,
		},
		{
			"MR Closed without a merge",
			"mytoken",
			"fixtures/gitlab-mr-close-no-merge.json",
			gitlab.MergeRequest{
				Kind: "merge_request",
				Project: gitlab.Project{
					PathWithNamespace: "pablo/testing-webhooks",
				},
				Attributes: gitlab.Attributes{
					State:       "closed",
					Title:       "Update README.md",
					Description: "other description",
					URL:         "https://git.yakshaving.art/pablo/testing-webhooks/merge_requests/2",
					Action:      "close",
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

			mr, err := parser.Parse(map[string][]string{
				"X-Gitlab-Token": {tc.token},
			}, b)
			a.NoErrorf(err, "could not unmarshal MR json")

			a.EqualValuesf(tc.expected, mr, "parsed merge request is not as expected")

			a.Equal(tc.expected.Text(), mr.Text())
			a.Equal(tc.expected.ProjectName(), mr.ProjectName())
			a.Equal(tc.shouldAnnounce, mr.ShouldAnnounce())
		})
	}
}

func TestInvalidPayloadErrs(t *testing.T) {
	a := assert.New(t)

	b, err := ioutil.ReadFile("fixtures/gitlab-push-event.json")
	a.Nil(err, "could not read fixture file")
	a.NotNilf(b, "content should not be nil")

	parser := gitlab.Parser{}
	mr, err := parser.Parse(map[string][]string{
		"X-Gitlab-Token": {""},
	}, b)
	a.Errorf(err, "json payload is not a merge request but a push")
	a.Equalf(gitlab.MergeRequest{}, mr, "merge request should be empty")
}

func TestHeadersMatcher(t *testing.T) {
	a := assert.New(t)
	p := gitlab.NewParser(".*", "token")

	a.Equal(false, p.MatchHeaders(map[string][]string{}))
	a.Equal(false, p.MatchHeaders(map[string][]string{"X-Github-Event": {"Merge Request Hook"}}))
	a.Equal(false, p.MatchHeaders(map[string][]string{"X-Github-Event": {"pull_request"}}))
	a.Equal(false, p.MatchHeaders(map[string][]string{"X-Gitlab-Event": {"pull_request"}}))
	a.Equal(true, p.MatchHeaders(map[string][]string{"X-Gitlab-Event": {"Merge Request Hook"}}))
}
