package string

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromString_MustParse(t *testing.T) {
	expected := `
	values := []BuildStatus{
		"in queue",
		"started",
		"building",
		"archiving",
		"uploading",
		"distributing",
		"failed",
		"succeeded",
		"complete",
		"interrupted",
	}
	for _, v := range values {
		if string(v) == s {
			return BuildStatus(s), nil
		}
	}
	return "", errors.Errorf("'%s' is not a valid build status", s)
`
	input := NewFromString("BuildStatus", "build status", []string{
		"in queue",
		"started",
		"building",
		"archiving",
		"uploading",
		"distributing",
		"failed",
		"succeeded",
		"complete",
		"interrupted",
	})
	result := input.MustParse()
	assert.Equal(t, expected, result)
}
