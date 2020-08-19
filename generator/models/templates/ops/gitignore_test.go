package ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitignore_MustParse(t *testing.T) {
	input := Gitignore{
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
.DS_Store
.idea
vendor
.app.env
}
`
	result := input.MustParse()
	assert.Equal(t, expected, result)
}
