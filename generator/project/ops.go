package project

import (
	"github.com/68696c6c/capricorn_rnd/generator/golang"
	"github.com/68696c6c/capricorn_rnd/generator/spec"
	"github.com/68696c6c/capricorn_rnd/generator/utils"
	"github.com/68696c6c/girraph"
)

const appEnv = `
DB_HOST={{ .MainDatabase.Host }}
DB_PORT={{ .MainDatabase.Port }}
DB_USERNAME={{ .MainDatabase.Username }}
DB_PASSWORD={{ .MainDatabase.Password }}
DB_DATABASE={{ .MainDatabase.Name }}
DB_DEBUG={{ .MainDatabase.Debug }}`

func makeAppEnv(ops spec.Ops) *golang.File {
	body, err := utils.ParseTemplateToString("tmp_template_appEnv", appEnv, ops)
	if err != nil {
		panic(err)
	}
	return golang.MakeFile(".app", "env").SetContents(body)
}

const appTemplateEnv = `
DB_HOST={{ .MainDatabase.Host }}
DB_PORT={{ .MainDatabase.Port }}
DB_USERNAME={{ .MainDatabase.Username }}
DB_PASSWORD=
DB_DATABASE={{ .MainDatabase.Name }}
DB_DEBUG={{ .MainDatabase.Debug }}`

func makeAppTemplateEnv(ops spec.Ops) *golang.File {
	body, err := utils.ParseTemplateToString("tmp_template_appTemplateEnv", appTemplateEnv, ops)
	if err != nil {
		panic(err)
	}
	return golang.MakeFile(".app.template", "env").SetContents(body)
}

const dockerCompose = `version: "3.6"

networks:
  default:
    external:
      name: docker-dev

volumes:
  pkg:
  db-volume:

services:

  app:
    image: app:dev
    command: ./app server 80
    depends_on:
      - db
    volumes:
      - pkg:/go/pkg
      - ./:/{{ .Workdir }}
    working_dir: /{{ .Workdir }}
    ports:
      - "80"
    env_file:
      - .app.env
    environment:
      VIRTUAL_HOST: {{ .AppHTTPAlias }}.local
      ENV: local
      HTTP_PORT: 80
      MIGRATION_PATH: /{{ .Workdir }}/src/database
    networks:
      default:
        aliases:
          - {{ .AppHTTPAlias }}.local

  db:
    platform: linux/amd64
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: {{ .MainDatabase.Password }}
      MYSQL_DATABASE: {{ .MainDatabase.Name }}
    ports:
      - "${HOST_DB_PORT:-3310}:{{ .MainDatabase.Port }}"
    volumes:
      - db-volume:/var/lib/mysql`

func makeDockerCompose(ops spec.Ops) *golang.File {
	body, err := utils.ParseTemplateToString("tmp_template_dockerCompose", dockerCompose, ops)
	if err != nil {
		panic(err)
	}
	return golang.MakeFile("docker-compose", "yml").SetContents(body)
}

const dockerfile = `version: "3.6"

networks:
  default:
    external:
      name: docker-dev

volumes:
  pkg:
  db-volume:

services:

  app:
    image: app:dev
    command: ./app server 80
    depends_on:
      - db
    volumes:
      - pkg:/go/pkg
      - ./:/{{ .Workdir }}
    working_dir: /{{ .Workdir }}
    ports:
      - "80"
    env_file:
      - .app.env
    environment:
      VIRTUAL_HOST: {{ .AppHTTPAlias }}.local
      ENV: local
      HTTP_PORT: 80
      MIGRATION_PATH: /{{ .Workdir }}/src/database
    networks:
      default:
        aliases:
          - {{ .AppHTTPAlias }}.local

  db:
    platform: linux/amd64
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: {{ .MainDatabase.Password }}
      MYSQL_DATABASE: {{ .MainDatabase.Name }}
    ports:
      - "${HOST_DB_PORT:-3310}:{{ .MainDatabase.Port }}"
    volumes:
      - db-volume:/var/lib/mysql`

func makeDockerfile(ops spec.Ops) *golang.File {
	body, err := utils.ParseTemplateToString("tmp_template_dockerfile", dockerfile, ops)
	if err != nil {
		panic(err)
	}
	return golang.MakeFile("Dockerfile", "").SetContents(body)
}

const makefile = `DCR = docker-compose run --rm

NETWORK_NAME ?= docker-dev
APP_NAME = app
DB_NAME = db

DOC_PATH_BASE = docs/swagger.json
DOC_PATH_FINAL = docs/api-spec.json

.PHONY: docs build

.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    image-local        build the $(APP_NAME):dev Docker image for local development'
	@echo '    image-built        build the $(APP_NAME):built Docker image for task running'
	@echo '    build              compile the app for use in Docker'
	@echo '    init               initialize the Go module'
	@echo '    deps               install dependencies'
	@echo '    setup-network      create local Docker network'
	@echo '    setup              set up local databases'
	@echo '    local              spin up local dev environment'
	@echo '    local-down         tear down local dev environment'
	@echo '    migrate            migrate the local database'
	@echo '    migration          create a new migration'
	@echo '    docs               build the Swagger docs'
	@echo '    docs-server        build and serve the Swagger docs'
	@echo '    test               run unit tests'
	@echo '    lint               run the linter'
	@echo '    lint-fix           run the linter and fix any problems'
	@echo


default: .DEFAULT

image-local:
	docker build . -f docker/Dockerfile --target dev -t $(APP_NAME):dev

image-built:
	docker build . -f docker/Dockerfile --target dev -t $(APP_NAME):built

build:
	$(DCR) $(APP_NAME) go build -i -o $(APP_NAME)

deps:
	$(DCR) $(APP_NAME) go mod tidy
	$(DCR) $(APP_NAME) go mod vendor

setup-network:
	docker network create docker-dev || exit 0

setup: setup-network image-local deps build
	$(DCR) $(DB_NAME) mysql -u root -p{{ .MainDatabase.Password }} -h $(DB_NAME) -e "CREATE DATABASE IF NOT EXISTS {{ .MainDatabase.Name }}"
	$(DCR) $(APP_NAME) bash -c "./$(APP_NAME) migrate up && ./$(APP_NAME) seed"

local: local-down build
	NETWORK_NAME="$(NETWORK_NAME)" docker-compose up

local-down:
	NETWORK_NAME="$(NETWORK_NAME)" docker-compose down

test:
	$(DCR) $(APP_NAME) go test ./... -cover

migrate: build
	$(DCR) $(APP_NAME) ./app migrate up

migration: build
	$(DCR) $(APP_NAME) goose -dir src/db/migrations create $(name)

docs: build
	$(DCR) $(APP_NAME) bash -c "GO111MODULE=off swagger generate spec -mo '$(DOC_PATH_BASE)'"
	$(DCR) $(APP_NAME) ./$(APP_NAME) gen:docs

docs-server: docs
	swagger serve "$(DOC_PATH_FINAL)"

lint:
	$(DCR) $(APP_NAME) golangci-lint run

lint-fix:
	$(DCR) $(APP_NAME) golangci-lint run --fix`

func makeMakefile(ops spec.Ops) *golang.File {
	body, err := utils.ParseTemplateToString("tmp_template_makefile", makefile, ops)
	if err != nil {
		panic(err)
	}
	return golang.MakeFile("Makefile", "").SetContents(body)
}

const gitignore = `
.DS_Store
.idea
vendor
.app.env`

const entrypoint = `#!/bin/sh
set -e

if [ -n "$CHAMBER_ENV" ]; then
    echo "importing chamber secrets from ${APP_NAME}"
    chamber exec "${APP_NAME}" -- "$@"
else
    eval "$@"
fi`

const preDeploy = `#!/bin/sh
set -eu

echo "Running migrations."
./app migrate up

echo "Done."`

const gitkeep = `ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
 ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ                        ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ 
ВВВВВВВВВВВВВВВВВВВВВВВ                          ВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ                            ВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ         ВВВВВВВВ            ВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ         ВВВВВВВВВВ          ВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ         ВВВВВВВВВВВ         ВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ         ВВВВВВВВВВ         ВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ         ВВВВВВВВВ         ВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ                         ВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ                         ВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ                             ВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ                               ВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ         ВВВВВВВВВВВВВ          ВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ         ВВВВВВВВВВВВВВ          ВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ         ВВВВВВВВВВВВВВ          ВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ         ВВВВВВВВВВВВВ           ВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ         ВВВВВВВВВВВВ           ВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ                               ВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВ                             ВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВ                          ВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
 ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
   ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ
      ВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВВ`

func makeOps() girraph.Tree[golang.Package] {
	build := golang.MakePackageNode("_build")
	build.GetMeta().SetFiles([]*golang.File{
		golang.MakeFile(".gitkeep", "").SetContents(gitkeep),
	})

	deploy := golang.MakePackageNode("_deploy")
	deploy.GetMeta().SetFiles([]*golang.File{
		golang.MakeFile(".gitkeep", "").SetContents(gitkeep),
	})

	env := golang.MakePackageNode("env")
	env.GetMeta().SetFiles([]*golang.File{
		golang.MakeFile("development", "yml").SetContents("hello"),
		golang.MakeFile("development.secrets", "yml").SetContents("hello"),
	})

	inputs := golang.MakePackageNode("inputs")
	inputs.GetMeta().SetFiles([]*golang.File{
		golang.MakeFile("development", "yml").SetContents("hello"),
		golang.MakeFile("development.secrets", "yml").SetContents("hello"),
	})

	templates := golang.MakePackageNode("templates").SetChildren([]girraph.Tree[golang.Package]{
		env,
		inputs,
	})

	scripts := golang.MakePackageNode("scripts")
	scripts.GetMeta().SetFiles([]*golang.File{
		golang.MakeFile("entrypoint", "sh").SetContents(entrypoint),
		golang.MakeFile("pre-deploy", "sh").SetContents(preDeploy),
	})

	return golang.MakePackageNode("ops").SetChildren([]girraph.Tree[golang.Package]{
		build,
		deploy,
		scripts,
		templates,
	})
}
