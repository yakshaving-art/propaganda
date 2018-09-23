package configuration_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/yakshaving.art/propaganda/configuration"
	"io/ioutil"
	"testing"
)

func TestLoadingValidConfiguration(t *testing.T) {
	a := assert.New(t)

	b, err := ioutil.ReadFile("fixtures/valid.yml")
	a.NoError(err)

	a.NoError(configuration.Load(b))

	a.EqualValues(
		configuration.Configuration{
			DefaultChannel: "#propaganda",
			Repositories: map[string]string{
				"pcarranza/test-webhooks": "#private_channel",
			},
		},
		configuration.GetConfiguration(),
	)
}
