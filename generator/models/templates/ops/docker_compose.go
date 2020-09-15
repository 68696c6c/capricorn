package ops

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/utils"
)

const DockerComposeTemplate = `version: "3.6"

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
    image: mysql:5.7
    environment:
      MYSQL_ROOT_PASSWORD: {{ .MainDatabase.Password }}
      MYSQL_DATABASE: {{ .MainDatabase.Name }}
    ports:
      - "${HOST_DB_PORT:-3310}:{{ .MainDatabase.Port }}"
    volumes:
      - db-volume:/var/lib/mysql
`

type DockerCompose struct {
	Name data.FileData `yaml:"name,omitempty"`
	Path data.PathData `yaml:"path,omitempty"`

	Data Ops `yaml:"data,omitempty"`
}

// This is only used for testing.
func (m DockerCompose) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_docker_compose", DockerComposeTemplate, m.Data)
	if err != nil {
		panic(err)
	}
	return result
}

func (m DockerCompose) Generate() error {
	err := utils.GenerateFile(m.Path.Base, m.Name.Full, DockerComposeTemplate, m.Data)
	if err != nil {
		return err
	}
	return nil
}
