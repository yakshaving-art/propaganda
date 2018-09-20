package gitlab_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/yakshaving.art/propaganda/gitlab"
)

func TestParsingPayloads(t *testing.T) {
	tt := []struct {
		name         string
		jsonFilename string
		expected     gitlab.MergeRequest
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
					State:  "opened",
					Title:  "[announce] Update README.md",
					URL:    "https://git.yakshaving.art/pablo/testing-webhooks/merge_requests/1",
					Action: "open",
				},
			},
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
					State:  "merged",
					Title:  "[announce] Update README.md",
					URL:    "https://git.yakshaving.art/pablo/testing-webhooks/merge_requests/1",
					Action: "merge",
				},
			},
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
					State:  "closed",
					Title:  "Update README.md",
					URL:    "https://git.yakshaving.art/pablo/testing-webhooks/merge_requests/2",
					Action: "close",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)
			b, err := ioutil.ReadFile(tc.jsonFilename)
			a.Nilf(err, "could not read fixture file %s", tc.jsonFilename)
			a.NotNilf(b, "content should not be nil")

			mr, err := gitlab.ParseMergeRequest(b)
			a.NoErrorf(err, "could not unmarshal MR json")

			a.EqualValuesf(tc.expected, mr, "parsed merge request is not as expected")
		})
	}
}

func TestInvalidPayloadErrs(t *testing.T) {
	a := assert.New(t)

	b, err := ioutil.ReadFile("fixtures/gitlab-push-event.json")
	a.Nil(err, "could not read fixture file")
	a.NotNilf(b, "content should not be nil")

	mr, err := gitlab.ParseMergeRequest(b)
	a.Errorf(err, "json payload is not a merge request but a push")
	a.Equalf(gitlab.MergeRequest{}, mr, "merge request should be empty")
}
