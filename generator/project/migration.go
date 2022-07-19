package project

import (
	"fmt"
	"strings"

	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

type MigrationTemplateData struct {
	UpModels   string
	DownModels string
	TxArgName  string
}

const migrationInit = `
	goose.AddMigration(upInitialMigration, downInitialMigration)`

func makeMigrationInit(data MigrationTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_migrationInit", migrationInit, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "init",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgGoose},
		},
		Body: body,
	}
}

const migrationUp = `
	goat.Init()

	db, err := goat.GetMigrationDB()
	if err != nil {
		return errors.Wrap(err, "failed to initialize migration connection")
	}
	db.AutoMigrate(&users.User{})

	return nil`

func makeMigrationUp(data MigrationTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_migrationUp", migrationUp, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "upInitialMigration",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgGoose},
		},
		Arguments: []golang.Value{
			{
				Name: data.TxArgName,
				Type: "*sql.Tx",
			},
		},
		ReturnValues: []golang.Value{
			{
				Type: "error",
			},
		},
		Body: body,
	}
}

const migrationDown = `
	goat.Init()

	db, err := goat.GetMigrationDB()
	if err != nil {
		return errors.Wrap(err, "failed to initialize migration connection")
	}
	db.DropTable(&users.User{})

	return nil`

func makeMigrationDown(data MigrationTemplateData) *golang.Function {
	body, err := utils.ParseTemplateToString("tmp_template_migrationDown", migrationDown, data)
	if err != nil {
		panic(err)
	}
	return &golang.Function{
		Name: "downInitialMigration",
		Imports: golang.Imports{
			Vendor: []golang.Package{PkgGoat, PkgGoose},
		},
		Arguments: []golang.Value{
			{
				Name: data.TxArgName,
				Type: "*sql.Tx",
			},
		},
		ReturnValues: []golang.Value{
			{
				Type: "error",
			},
		},
		Body: body,
	}
}

func MakeInitialMigration(modelReferences []string, modelsPkg golang.Package) *golang.File {
	var upReferences []string
	var downReferences []string
	for _, modelRef := range modelReferences {
		upReferences = append(upReferences, fmt.Sprintf("db.AutoMigrate(&%s{})", modelRef))
		downReferences = append(downReferences, fmt.Sprintf("db.DropTable(&%s{})", modelRef))
	}
	data := MigrationTemplateData{
		UpModels:   strings.Join(upReferences, "\n"),
		DownModels: strings.Join(downReferences, "\n"),
		TxArgName:  "_",
	}

	migrationName := utils.MakeGooseMigrationName("", "initial_migration")
	return golang.MakeFile(migrationName).SetImports(golang.Imports{
		Standard: []golang.Package{PkgStdSqlDriver},
		Vendor:   []golang.Package{PkgGoat, PkgGoose, modelsPkg},
		App:      []golang.Package{modelsPkg},
	}).SetFunctions([]*golang.Function{
		makeMigrationInit(data),
		makeMigrationUp(data),
		makeMigrationDown(data),
	})
}
