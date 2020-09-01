package spec

func GetFixtureSpec() Spec {
	return Spec{
		Name:    "Test Example",
		Module:  "github.com/68696c6c/test-example",
		License: "none",
		Author: Author{
			Name:         "Aaron Hill",
			Email:        "68696c6c@gmail.com",
			Organization: "GOAT",
		},
		Commands: []Command{
			{
				Name: "cmd",
				Args: []string{},
				Use:  "This is an example command",
			},
			{
				Name: "cmd:one-arg",
				Args: []string{"id"},
				Use:  "This is an example command",
			},
			{
				Name: "cmd:two-args",
				Args: []string{"id", "name"},
				Use:  "This is an example command",
			},
		},
		Resources: []Resource{
			{
				Name:    "organization",
				HasMany: []string{"users"},
				Fields: []ResourceField{
					{
						Name:     "name",
						Type:     "string",
						Required: true,
					},
				},
			},
			{
				Name:      "user",
				BelongsTo: []string{"organization"},
				HasMany:   []string{"tokens"},
				Fields: []ResourceField{
					{
						Name:     "name",
						Type:     "string",
						Required: true,
						Unique:   false,
					},
					{
						Name:     "email",
						Type:     "string",
						Required: true,
						Unique:   true,
					},
				},
			},
			{
				Name:      "token",
				BelongsTo: []string{"user"},
				Fields: []ResourceField{
					{
						Name:     "key",
						Type:     "string",
						Required: true,
						Unique:   true,
					},
					{
						Name:     "expires",
						Type:     "time.Time",
						Required: true,
						Unique:   false,
					},
				},
				Actions: []string{"create", "delete"},
				Custom:  []string{"refresh"},
			},
		},
	}
}

const FixtureSpecYAML = `name:
  space: test example
  snake: test_example
  kebob: test-example
  exported: TestExample
  unexported: testExample
package:
  reference: test-example
  name:
    space: test example
    snake: test_example
    kebob: test-example
    exported: TestExample
    unexported: testExample
  path:
    base: github.com/68696c6c
    full: github.com/68696c6c/test-example
ops:
  workdir: test-example
  app_http_alias: test-example
  database:
    host: db
    port: "3306"
    username: root
    password: secret
    name: test_example
    debug: "1"
packages:
  src:
    reference: src
    name:
      space: src
      snake: src
      kebob: src
      exported: Src
      unexported: src
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/src
  ops:
    reference: ops
    name:
      space: ops
      snake: ops
      kebob: ops
      exported: Ops
      unexported: ops
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/ops
  docker:
    reference: docker
    name:
      space: docker
      snake: docker
      kebob: docker
      exported: Docker
      unexported: docker
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/docker
  app:
    reference: app
    name:
      space: app
      snake: app
      kebob: app
      exported: App
      unexported: app
    path:
      base: github.com/68696c6c/test-example/src
      full: github.com/68696c6c/test-example/src/app
  cmd:
    reference: cmd
    name:
      space: cmd
      snake: cmd
      kebob: cmd
      exported: Cmd
      unexported: cmd
    path:
      base: github.com/68696c6c/test-example/src
      full: github.com/68696c6c/test-example/src/cmd
  database:
    reference: db
    name:
      space: db
      snake: db
      kebob: db
      exported: Db
      unexported: db
    path:
      base: github.com/68696c6c/test-example/src
      full: github.com/68696c6c/test-example/src/db
  http:
    reference: http
    name:
      space: http
      snake: http
      kebob: http
      exported: Http
      unexported: http
    path:
      base: github.com/68696c6c/test-example/src
      full: github.com/68696c6c/test-example/src/http
  repos:
    reference: repos
    name:
      space: repos
      snake: repos
      kebob: repos
      exported: Repos
      unexported: repos
    path:
      base: github.com/68696c6c/test-example/src
      full: github.com/68696c6c/test-example/src/repos
  models:
    reference: models
    name:
      space: models
      snake: models
      kebob: models
      exported: Models
      unexported: models
    path:
      base: github.com/68696c6c/test-example/src
      full: github.com/68696c6c/test-example/src/models
  migrations:
    reference: migrations
    name:
      space: migrations
      snake: migrations
      kebob: migrations
      exported: Migrations
      unexported: migrations
    path:
      base: github.com/68696c6c/test-example/src/db
      full: github.com/68696c6c/test-example/src/db/migrations
  seeders:
    reference: seeders
    name:
      space: seeders
      snake: seeders
      kebob: seeders
      exported: Seeders
      unexported: seeders
    path:
      base: github.com/68696c6c/test-example/src/db
      full: github.com/68696c6c/test-example/src/db/seeders
  domains:
    reference: app
    name:
      space: app
      snake: app
      kebob: app
      exported: App
      unexported: app
    path:
      base: github.com/68696c6c/test-example/src
      full: github.com/68696c6c/test-example/src/app
commands:
- name:
    space: cmd
    snake: cmd
    kebob: cmd
    exported: Cmd
    unexported: cmd
- name:
    space: cmd:one arg
    snake: cmd:one_arg
    kebob: cmd:one-arg
    exported: Cmd:oneArg
    unexported: cmd:oneArg
- name:
    space: cmd:two args
    snake: cmd:two_args
    kebob: cmd:two-args
    exported: Cmd:twoArgs
    unexported: cmd:twoArgs
resources:
- key:
    resource: organization
    field: ""
  name:
    single:
      space: organization
      snake: organization
      kebob: organization
      exported: Organization
      unexported: organization
    plural:
      space: organizations
      snake: organizations
      kebob: organizations
      exported: Organizations
      unexported: organizations
  fields:
  - key:
      resource: organization
      field: id
    name:
      space: id
      snake: id
      kebob: id
      exported: Id
      unexported: id
    type: goat.ID
    index: null
    is_required: false
    is_primary: true
    is_goat_field: true
  - key:
      resource: organization
      field: created_at
    name:
      space: created at
      snake: created_at
      kebob: created-at
      exported: CreatedAt
      unexported: createdAt
    type: time.Time
    index: null
    is_required: false
    is_primary: false
    is_goat_field: true
  - key:
      resource: organization
      field: updated_at
    name:
      space: updated at
      snake: updated_at
      kebob: updated-at
      exported: UpdatedAt
      unexported: updatedAt
    type: '*time.Time'
    index: null
    is_required: false
    is_primary: false
    is_goat_field: true
  - key:
      resource: organization
      field: deleted_at
    name:
      space: deleted at
      snake: deleted_at
      kebob: deleted-at
      exported: DeletedAt
      unexported: deletedAt
    type: '*time.Time'
    index: null
    is_required: false
    is_primary: false
    is_goat_field: true
  - key:
      resource: organization
      field: name
    name:
      space: name
      snake: name
      kebob: name
      exported: Name
      unexported: name
    type: string
    index: null
    is_required: true
    is_primary: false
    is_goat_field: false
  controller:
    name:
      space: organizations
      snake: organizations
      kebob: organizations
      exported: Organizations
      unexported: organizations
    actions:
    - list
    - view
    - create
    - update
    - delete
  repo:
    name:
      space: organizations
      snake: organizations
      kebob: organizations
      exported: Organizations
      unexported: organizations
    actions:
    - list
    - view
    - create
    - update
    - delete
  service:
    name:
      space: organizations service
      snake: organizations_service
      kebob: organizations-service
      exported: OrganizationsService
      unexported: organizationsService
    actions: []
- key:
    resource: user
    field: ""
  name:
    single:
      space: user
      snake: user
      kebob: user
      exported: User
      unexported: user
    plural:
      space: users
      snake: users
      kebob: users
      exported: Users
      unexported: users
  fields:
  - key:
      resource: user
      field: id
    name:
      space: id
      snake: id
      kebob: id
      exported: Id
      unexported: id
    type: goat.ID
    index: null
    is_required: false
    is_primary: true
    is_goat_field: true
  - key:
      resource: user
      field: created_at
    name:
      space: created at
      snake: created_at
      kebob: created-at
      exported: CreatedAt
      unexported: createdAt
    type: time.Time
    index: null
    is_required: false
    is_primary: false
    is_goat_field: true
  - key:
      resource: user
      field: updated_at
    name:
      space: updated at
      snake: updated_at
      kebob: updated-at
      exported: UpdatedAt
      unexported: updatedAt
    type: '*time.Time'
    index: null
    is_required: false
    is_primary: false
    is_goat_field: true
  - key:
      resource: user
      field: deleted_at
    name:
      space: deleted at
      snake: deleted_at
      kebob: deleted-at
      exported: DeletedAt
      unexported: deletedAt
    type: '*time.Time'
    index: null
    is_required: false
    is_primary: false
    is_goat_field: true
  - key:
      resource: user
      field: name
    name:
      space: name
      snake: name
      kebob: name
      exported: Name
      unexported: name
    type: string
    index: null
    is_required: true
    is_primary: false
    is_goat_field: false
  - key:
      resource: user
      field: email
    name:
      space: email
      snake: email
      kebob: email
      exported: Email
      unexported: email
    type: string
    index: null
    is_required: true
    is_primary: false
    is_goat_field: false
  controller:
    name:
      space: users
      snake: users
      kebob: users
      exported: Users
      unexported: users
    actions:
    - list
    - view
    - create
    - update
    - delete
  repo:
    name:
      space: users
      snake: users
      kebob: users
      exported: Users
      unexported: users
    actions:
    - list
    - view
    - create
    - update
    - delete
  service:
    name:
      space: users service
      snake: users_service
      kebob: users-service
      exported: UsersService
      unexported: usersService
    actions: []
- key:
    resource: token
    field: ""
  name:
    single:
      space: token
      snake: token
      kebob: token
      exported: Token
      unexported: token
    plural:
      space: tokens
      snake: tokens
      kebob: tokens
      exported: Tokens
      unexported: tokens
  fields:
  - key:
      resource: token
      field: id
    name:
      space: id
      snake: id
      kebob: id
      exported: Id
      unexported: id
    type: goat.ID
    index: null
    is_required: false
    is_primary: true
    is_goat_field: true
  - key:
      resource: token
      field: created_at
    name:
      space: created at
      snake: created_at
      kebob: created-at
      exported: CreatedAt
      unexported: createdAt
    type: time.Time
    index: null
    is_required: false
    is_primary: false
    is_goat_field: true
  - key:
      resource: token
      field: updated_at
    name:
      space: updated at
      snake: updated_at
      kebob: updated-at
      exported: UpdatedAt
      unexported: updatedAt
    type: '*time.Time'
    index: null
    is_required: false
    is_primary: false
    is_goat_field: true
  - key:
      resource: token
      field: deleted_at
    name:
      space: deleted at
      snake: deleted_at
      kebob: deleted-at
      exported: DeletedAt
      unexported: deletedAt
    type: '*time.Time'
    index: null
    is_required: false
    is_primary: false
    is_goat_field: true
  - key:
      resource: token
      field: key
    name:
      space: key
      snake: key
      kebob: key
      exported: Key
      unexported: key
    type: string
    index: null
    is_required: true
    is_primary: false
    is_goat_field: false
  - key:
      resource: token
      field: expires
    name:
      space: expires
      snake: expires
      kebob: expires
      exported: Expires
      unexported: expires
    type: time.Time
    index: null
    is_required: true
    is_primary: false
    is_goat_field: false
  controller:
    name:
      space: tokens
      snake: tokens
      kebob: tokens
      exported: Tokens
      unexported: tokens
    actions:
    - create
    - delete
  repo:
    name:
      space: tokens
      snake: tokens
      kebob: tokens
      exported: Tokens
      unexported: tokens
    actions:
    - create
    - delete
  service:
    name:
      space: tokens service
      snake: tokens_service
      kebob: tokens-service
      exported: TokensService
      unexported: tokensService
    actions:
    - refresh
`
