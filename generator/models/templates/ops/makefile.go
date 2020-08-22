package ops

import (
	utils2 "github.com/68696c6c/capricorn/generator/models/utils"
	"github.com/68696c6c/capricorn/generator/utils"
)

const MakefileTemplate = `DCR = docker-compose run --rm

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
	$(DCR) $(APP_NAME) golangci-lint run --fix
`

type Makefile struct {
	Name utils2.FileData `yaml:"name"`
	Path utils2.PathData `yaml:"path"`

	Data Ops `yaml:"data"`
}

// This is only used for testing.
func (m Makefile) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_makefile", MakefileTemplate, m.Data)
	if err != nil {
		panic(err)
	}
	return result
}

func (m Makefile) Generate() error {
	err := utils.GenerateFile(m.Path.Base, m.Name.Full, MakefileTemplate, m.Data)
	if err != nil {
		return err
	}
	return nil
}