package module

import "github.com/68696c6c/capricorn/generator/models/data"

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

type Packages struct {
	SRC        data.PackageData `yaml:"src"`
	OPS        data.PackageData `yaml:"ops"`
	Docker     data.PackageData `yaml:"docker"`
	App        data.PackageData `yaml:"app"`
	CMD        data.PackageData `yaml:"cmd"`
	DB         data.PackageData `yaml:"database"`
	HTTP       data.PackageData `yaml:"http"`
	Repos      data.PackageData `yaml:"repos"`
	Models     data.PackageData `yaml:"models"`
	Migrations data.PackageData `yaml:"migrations"`
	Seeders    data.PackageData `yaml:"seeders"`
	Domains    data.PackageData `yaml:"domains"`
}

func makePackages(root string) Packages {
	pSRC := data.MakePackageData(root, pkgSRC)
	srcPath := pSRC.Path.Full

	pApp := data.MakePackageData(srcPath, pkgApp)

	pDB := data.MakePackageData(srcPath, pkgDB)
	dbPath := pDB.Path.Full

	result := Packages{
		Docker:     data.MakePackageData(root, pkgDocker),
		OPS:        data.MakePackageData(root, pkgOPS),
		SRC:        pSRC,
		App:        pApp,
		CMD:        data.MakePackageData(srcPath, pkgCMD),
		HTTP:       data.MakePackageData(srcPath, pkgHTTP),
		Repos:      data.MakePackageData(srcPath, pkgRepos),
		Models:     data.MakePackageData(srcPath, pkgModels),
		DB:         pDB,
		Migrations: data.MakePackageData(dbPath, pkgMigrations),
		Seeders:    data.MakePackageData(dbPath, pkgSeeders),
		Domains:    pApp,
	}

	return result
}
