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
  snake: test_example
  kebob: test-example
  exported: TestExample
  unexported: testExample
path:
  full: github.com/68696c6c/test-example
  base: test-example
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
    name:
      snake: src
      kebob: src
      exported: Src
      unexported: src
    path:
      full: github.com/68696c6c/test-example/src
      base: src
  ops:
    name:
      snake: ops
      kebob: ops
      exported: Ops
      unexported: ops
    path:
      full: github.com/68696c6c/test-example/ops
      base: ops
  docker:
    name:
      snake: docker
      kebob: docker
      exported: Docker
      unexported: docker
    path:
      full: github.com/68696c6c/test-example/docker
      base: docker
  app:
    name:
      snake: app
      kebob: app
      exported: App
      unexported: app
    path:
      full: github.com/68696c6c/test-example/src/app
      base: app
  cmd:
    name:
      snake: cmd
      kebob: cmd
      exported: Cmd
      unexported: cmd
    path:
      full: github.com/68696c6c/test-example/src/cmd
      base: cmd
  database:
    name:
      snake: db
      kebob: db
      exported: Db
      unexported: db
    path:
      full: github.com/68696c6c/test-example/src/db
      base: db
  http:
    name:
      snake: http
      kebob: http
      exported: Http
      unexported: http
    path:
      full: github.com/68696c6c/test-example/src/http
      base: http
  repos:
    name:
      snake: repos
      kebob: repos
      exported: Repos
      unexported: repos
    path:
      full: github.com/68696c6c/test-example/src/repos
      base: repos
  models:
    name:
      snake: models
      kebob: models
      exported: Models
      unexported: models
    path:
      full: github.com/68696c6c/test-example/src/models
      base: models
  migrations:
    name:
      snake: migrations
      kebob: migrations
      exported: Migrations
      unexported: migrations
    path:
      full: github.com/68696c6c/test-example/src/db/migrations
      base: migrations
  seeders:
    name:
      snake: seeders
      kebob: seeders
      exported: Seeders
      unexported: seeders
    path:
      full: github.com/68696c6c/test-example/src/db/seeders
      base: seeders
  domains:
    name:
      snake: app
      kebob: app
      exported: App
      unexported: app
    path:
      full: github.com/68696c6c/test-example/src/app
      base: app
commands:
- name:
    snake: cmd
    kebob: cmd
    exported: Cmd
    unexported: cmd
- name:
    snake: cmd:one_arg
    kebob: cmd:one-arg
    exported: Cmd:oneArg
    unexported: cmd:oneArg
- name:
    snake: cmd:two_args
    kebob: cmd:two-args
    exported: Cmd:twoArgs
    unexported: cmd:twoArgs
resources:
- key:
    resource: organization
    field: ""
  name:
    snake: organization
    kebob: organization
    exported: Organization
    unexported: organization
  fields:
  - key:
      resource: organization
      field: id
    name:
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
      snake: organizations
      kebob: organizations
      exported: Organizations
      unexported: organizations
    actions: []
  repo:
    name:
      snake: organizations
      kebob: organizations
      exported: Organizations
      unexported: organizations
    actions: []
  service:
    name:
      snake: organization_service
      kebob: organization-service
      exported: OrganizationService
      unexported: organizationService
    actions: []
- key:
    resource: user
    field: ""
  name:
    snake: user
    kebob: user
    exported: User
    unexported: user
  fields:
  - key:
      resource: user
      field: id
    name:
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
      snake: users
      kebob: users
      exported: Users
      unexported: users
    actions: []
  repo:
    name:
      snake: users
      kebob: users
      exported: Users
      unexported: users
    actions: []
  service:
    name:
      snake: user_service
      kebob: user-service
      exported: UserService
      unexported: userService
    actions: []
- key:
    resource: token
    field: ""
  name:
    snake: token
    kebob: token
    exported: Token
    unexported: token
  fields:
  - key:
      resource: token
      field: id
    name:
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
      snake: tokens
      kebob: tokens
      exported: Tokens
      unexported: tokens
    actions:
    - snake: create
      kebob: create
      exported: Create
      unexported: create
    - snake: delete
      kebob: delete
      exported: Delete
      unexported: delete
  repo:
    name:
      snake: tokens
      kebob: tokens
      exported: Tokens
      unexported: tokens
    actions:
    - snake: create
      kebob: create
      exported: Create
      unexported: create
    - snake: delete
      kebob: delete
      exported: Delete
      unexported: delete
  service:
    name:
      snake: token_service
      kebob: token-service
      exported: TokenService
      unexported: tokenService
    actions:
    - snake: refresh
      kebob: refresh
      exported: Refresh
      unexported: refresh
`
