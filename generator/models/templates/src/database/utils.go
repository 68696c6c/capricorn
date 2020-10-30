package database

import (
	"fmt"
	"time"
)

func makeGooseMigrationName(name string) string {
	// This is copied from github.com/pressly/goose/create.go CreateWithTemplate function and should match what that function does.
	version := time.Now().Format("20060102150405")
	return fmt.Sprintf("%v_%v.%v", version, name, "go")
}
