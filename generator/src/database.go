package src

import (
	"fmt"
	"time"

	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const migrationTemplate = `package migrations

import (
	"database/sql"


	{{- range $key, $value := .Domains }}
	"{{ $value.Import }}"
	{{- end }}

	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

func init() {
	goat.Init()
	goose.AddMigration(upInitialMigration, downInitialMigration)
}

func upInitialMigration(tx *sql.Tx) error {
	db, err := goat.GetMigrationDB()
	if err != nil {
		return errors.Wrap(err, "failed to initialize migration connection")
	}

	{{- range $key, $value := .Domains }}
	db.AutoMigrate(&{{ $value.Name }}.{{ $value.Model.Name }}{})
	{{- end }}

	return nil
}

func downInitialMigration(tx *sql.Tx) error {
	db, err := goat.GetMigrationDB()
	if err != nil {
		return errors.Wrap(err, "failed to initialize migration connection")
	}


	{{- range $key, $value := .Domains }}
	db.DropTable(&{{ $value.Name }}.{{ $value.Model.Name }}{})
	{{- end }}

	return nil
}
`

const logPrefixDB = "CreateDatabase"

func CreateDatabase(spec *models.Project, logger *logrus.Logger) error {
	logger.Infof("%s | creating migrations %s", logPrefixDB, spec.Paths.Migrations)
	err := utils.CreateDir(spec.Paths.Migrations)
	if err != nil {
		return errors.Wrapf(err, "failed to create database migrations directory '%s'", spec.Paths.Migrations)
	}

	// dt := goat.TimeToYMDHISString(time.Now())
	dt := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_initial_migration.go", dt)
	err = utils.GenerateFile(spec.Paths.Migrations, fileName, migrationTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to generate initial migration")
	}

	return nil
}
