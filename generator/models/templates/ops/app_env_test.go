package ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppEnv_MustParse(t *testing.T) {
	input := AppEnv{
		Data: Ops{
			Workdir:      "example_workdir",
			AppHTTPAlias: "example_alias",
			MainDatabase: Database{
				Host:     "example_host",
				Port:     "1234",
				Username: "example_user",
				Password: "example_password",
				Name:     "example_dbname",
				Debug:    "1",
			},
		},
	}
	expected := `
DB_HOST=example_host
DB_PORT=1234
DB_USERNAME=example_user
DB_PASSWORD=example_password
DB_DATABASE=example_dbname
DB_DEBUG=1
`
	result := input.MustParse()
	assert.Equal(t, expected, result)
}
