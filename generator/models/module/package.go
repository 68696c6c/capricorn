package module

import (
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/models/templates"
	"github.com/68696c6c/capricorn/generator/utils"
)

const (
	pkgSRC        = "src"
	pkgOPS        = "ops"
	pkgDocker     = "docker"
	pkgApp        = "app"
	pkgCMD        = "cmd"
	pkgDB         = "db"
	pkgHTTP       = "http"
	pkgRepos      = "repos"
	pkgModels     = "models"
	pkgMigrations = "migrations"
	pkgSeeders    = "seeders"
)

type Package struct {
	Name models.Name        `yaml:"name"`
	Path templates.PathData `yaml:"path"` // e.g. full: module/path/src/app/domain, base: domain
}

type Packages struct {
	SRC        Package `yaml:"src"`
	OPS        Package `yaml:"ops"`
	Docker     Package `yaml:"docker"`
	App        Package `yaml:"app"`
	CMD        Package `yaml:"cmd"`
	DB         Package `yaml:"database"`
	HTTP       Package `yaml:"http"`
	Repos      Package `yaml:"repos"`
	Models     Package `yaml:"models"`
	Migrations Package `yaml:"migrations"`
	Seeders    Package `yaml:"seeders"`
	Domains    Package `yaml:"domains"`
}

func makePackages(root string, resources []Resource) Packages {
	pSRC := MakePackage(root, pkgSRC)
	srcPath := pSRC.Path.Full

	pApp := MakePackage(srcPath, pkgApp)

	pDB := MakePackage(srcPath, pkgDB)
	dbPath := pDB.Path.Full

	result := Packages{
		Docker:     MakePackage(root, pkgDocker),
		OPS:        MakePackage(root, pkgOPS),
		SRC:        pSRC,
		App:        pApp,
		CMD:        MakePackage(srcPath, pkgCMD),
		HTTP:       MakePackage(srcPath, pkgHTTP),
		Repos:      MakePackage(srcPath, pkgRepos),
		Models:     MakePackage(srcPath, pkgModels),
		DB:         pDB,
		Migrations: MakePackage(dbPath, pkgMigrations),
		Seeders:    MakePackage(dbPath, pkgSeeders),
		Domains:    pApp,
	}

	return result
}

func MakePackage(pkgBase, pkgName string) Package {
	return Package{
		Name: models.MakeName(pkgName),
		Path: templates.PathData{
			Base: pkgName,
			Full: utils.JoinPath(pkgBase, pkgName),
		},
	}
}

func (m Package) GetImport() string {
	return m.Path.Full
}

func (m Package) GetReference() string {
	return m.Path.Base
}
