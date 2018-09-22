package gitlab_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/yakshaving.art/propaganda/gitlab"
)

func TestParsingPayloads(t *testing.T) {
	parser := gitlab.Parser{}
	tt := []struct {
		name           string
		jsonFilename   string
		expected       gitlab.MergeRequest
		shouldAnnounce bool
	}{
		{
			"MR Create",
			"fixtures/gitlab-mr-create.json",
			gitlab.MergeRequest{
				Kind: "merge_request",
				Project: gitlab.Project{
					PathWithNamespace: "pablo/testing-webhooks",
				},
				Attributes: gitlab.Attributes{
					State:       "opened",
					Title:       "[announce] Update README.md",
					Description: "Something in the description",
					URL:         "https://git.yakshaving.art/pablo/testing-webhooks/merge_requests/1",
					Action:      "open",
				},
			},
			false,
		},
		{
			"MR Merged",
			"fixtures/gitlab-mr-merged.json",
			gitlab.MergeRequest{
				Kind: "merge_request",
				Project: gitlab.Project{
					PathWithNamespace: "pablo/testing-webhooks",
				},
				Attributes: gitlab.Attributes{
					State:       "merged",
					Title:       "[announce] Update README.md",
					Description: "Something in the description",
					URL:         "https://git.yakshaving.art/pablo/testing-webhooks/merge_requests/1",
					Action:      "merge",
				},
			},
			true,
		},
		{
			"MR Closed without a merge",
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

			mr, err := parser.Parse(b)
			a.NoErrorf(err, "could not unmarshal MR json")

			a.EqualValuesf(tc.expected, mr, "parsed merge request is not as expected")

			a.Equal(tc.expected.Title(), mr.Title())
			a.Equal(tc.expected.Text(), mr.Text())
			a.Equal(tc.expected.URL(), mr.URL())
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
	mr, err := parser.Parse(b)
	a.Errorf(err, "json payload is not a merge request but a push")
	a.Equalf(gitlab.MergeRequest{}, mr, "merge request should be empty")
}
