package enums

import (
	"fmt"

	"database/sql/driver"
	"github.com/pkg/errors"
)

// Represents the status of a site build.
type BuildStatus string

const buildStatusName string = "build status"

func BuildStatusFromString(s string) (BuildStatus, error) {
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
	return "", errors.Errorf("%s not a valid %s", s, buildStatusName)
}

func (t BuildStatus) String() string {
	return string(t)
}

// Gorm calls Scan to convert a raw database value into our custom type.
func (t *BuildStatus) Scan(value interface{}) error {
	stringValue := fmt.Sprintf("%v", value)
	result, err := BuildStatusFromString(stringValue)
	if err != nil {
		return err
	}
	*t = result
	return nil
}

// Gorm calls Value to convert our custom type into something it can work with.
func (t BuildStatus) Value() (driver.Value, error) {
	return string(t), nil
}
