package module

import "github.com/68696c6c/capricorn/generator/models/utils"

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
	SRC        utils.PackageData `yaml:"src"`
	OPS        utils.PackageData `yaml:"ops"`
	Docker     utils.PackageData `yaml:"docker"`
	App        utils.PackageData `yaml:"app"`
	CMD        utils.PackageData `yaml:"cmd"`
	DB         utils.PackageData `yaml:"database"`
	HTTP       utils.PackageData `yaml:"http"`
	Repos      utils.PackageData `yaml:"repos"`
	Models     utils.PackageData `yaml:"models"`
	Migrations utils.PackageData `yaml:"migrations"`
	Seeders    utils.PackageData `yaml:"seeders"`
	Domains    utils.PackageData `yaml:"domains"`
}

func makePackages(root string) Packages {
	pSRC := utils.MakePackageData(root, pkgSRC)
	srcPath := pSRC.Path.Full

	pApp := utils.MakePackageData(srcPath, pkgApp)

	pDB := utils.MakePackageData(srcPath, pkgDB)
	dbPath := pDB.Path.Full

	result := Packages{
		Docker:     utils.MakePackageData(root, pkgDocker),
		OPS:        utils.MakePackageData(root, pkgOPS),
		SRC:        pSRC,
		App:        pApp,
		CMD:        utils.MakePackageData(srcPath, pkgCMD),
		HTTP:       utils.MakePackageData(srcPath, pkgHTTP),
		Repos:      utils.MakePackageData(srcPath, pkgRepos),
		Models:     utils.MakePackageData(srcPath, pkgModels),
		DB:         pDB,
		Migrations: utils.MakePackageData(dbPath, pkgMigrations),
		Seeders:    utils.MakePackageData(dbPath, pkgSeeders),
		Domains:    pApp,
	}

	return result
}
