package project

import "github.com/68696c6c/capricorn_rnd/generator/golang"

var (
	PkgStdFmt       = golang.MakePackage("fmt")
	PkgStdOs        = golang.MakePackage("os")
	PkgStdSqlDriver = golang.MakePackage("database/sql/driver")

	PkgGoat       = golang.MakePackage("github.com/68696c6c/goat")
	PkgQuery      = golang.MakePackage("github.com/68696c6c/goat/query")
	PkgGin        = golang.MakePackage("github.com/gin-gonic/gin")
	PkgErrors     = golang.MakePackage("github.com/pkg/errors")
	PkgGorm       = golang.MakePackage("github.com/jinzhu/gorm")
	PkgValidation = golang.MakePackage("github.com/go-ozzo/ozzo-validation")
	PkgLogrus     = golang.MakePackage("github.com/sirupsen/logrus")
	PkgGoose      = golang.MakePackage("github.com/pressly/goose")
	PkgCobra      = golang.MakePackage("github.com/spf13/cobra")
	PkgViper      = golang.MakePackage("github.com/spf13/viper")
	PkgSqlDriver  = golang.MakePackage("_ \"github.com/go-sql-driver/mysql\"")
)
