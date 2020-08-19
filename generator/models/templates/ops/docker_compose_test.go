package ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerCompose_MustParse(t *testing.T) {
	input := DockerCompose{
		Data: Ops{
			Workdir:      "example_workdir",
			AppHTTPAlias: "example_alias",
			MainDatabase: Database{
				Host:     "example_host",
				Port:     "1234",
				Username: "example_user",
				Password: "example_password",
				Name:     "example_dbname",
				Debug:    "1",
			},
		},
	}
	expected := `version: "3.6"

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
      - ./:/example_workdir
    working_dir: /example_workdir
    ports:
      - "80"
    env_file:
      - .app.env
    environment:
      VIRTUAL_HOST: example_alias.local
      ENV: local
      HTTP_PORT: 80
      MIGRATION_PATH: /example_workdir/src/database
    networks:
      default:
        aliases:
          - example_alias.local

  db:
    image: mysql:5.7
    environment:
      MYSQL_ROOT_PASSWORD: example_password
      MYSQL_DATABASE: example_dbname
    ports:
      - "${HOST_DB_PORT:-3310}:1234"
    volumes:
      - db-volume:/var/lib/mysql
`
	result := input.MustParse()
	assert.Equal(t, expected, result)
}
