package module

import (
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/models/templates"
	"github.com/68696c6c/capricorn/generator/utils"
)

const (
	dirSRC        = "src"
	dirOPS        = "ops"
	dirDocker     = "docker"
	dirApp        = "app"
	dirCMD        = "cmd"
	dirDB         = "db"
	dirHTTP       = "http"
	dirRepos      = "repos"
	dirModels     = "models"
	dirMigrations = "migrations"
	dirSeeders    = "seeders"
)

type Package struct {
	Name models.Name        `yaml:"name"`
	Path templates.FileData `yaml:"path"` // e.g. full: module/path/src/app/domain, base: domain
}

type Packages struct {
	SRC        Package   `yaml:"src"`
	OPS        Package   `yaml:"ops"`
	Docker     Package   `yaml:"docker"`
	App        Package   `yaml:"app"`
	CMD        Package   `yaml:"cmd"`
	DB         Package   `yaml:"database"`
	HTTP       Package   `yaml:"http"`
	Repos      Package   `yaml:"repos"`
	Models     Package   `yaml:"models"`
	Migrations Package   `yaml:"migrations"`
	Seeders    Package   `yaml:"seeders"`
	Domains    []Package `yaml:"domains"`
}

func makePackages(root string, resources []Resource) Packages {
	pkgSRC := makePackage(root, dirSRC)
	srcPath := pkgSRC.Path.Full

	pkgApp := makePackage(srcPath, dirApp)
	appPath := pkgApp.Path.Full

	pkgDB := makePackage(srcPath, dirDB)
	dbPath := pkgDB.Path.Full

	result := Packages{
		Docker:     makePackage(root, dirDocker),
		OPS:        makePackage(root, dirOPS),
		SRC:        pkgSRC,
		App:        pkgApp,
		CMD:        makePackage(srcPath, dirCMD),
		HTTP:       makePackage(srcPath, dirHTTP),
		Repos:      makePackage(srcPath, dirRepos),
		Models:     makePackage(srcPath, dirModels),
		DB:         pkgDB,
		Migrations: makePackage(dbPath, dirMigrations),
		Seeders:    makePackage(dbPath, dirSeeders),
	}

	for _, r := range resources {
		result.Domains = append(result.Domains, makePackage(appPath, r.Name.Kebob))
	}

	return result
}

func makePackage(pkgBase, pkgName string) Package {
	return Package{
		Name: makeName(pkgName),
		Path: templates.FileData{
			Base: pkgName,
			Full: utils.JoinPath(pkgBase, pkgName),
		},
	}
}
