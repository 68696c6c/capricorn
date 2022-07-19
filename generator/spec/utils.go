package spec

// import "github.com/68696c6c/capricorn/generator/models/templates/ops"

func GetFixtureInput() []byte {
	input := `name: Test Example
module: github.com/68696c6c/test-example
license: none
author:
  name: Aaron Hill
  email: 68696c6c@gmail.com
  organization: GOAT
ops:
  workdir: test-example
  app_http_alias: test-example.local
  database:
    host: db
    port: 3306
    username: root
    password: secret
    name: test_example
    debug: 1
commands:
- name: cmd
  use: This is an example command
- name: cmd:one-arg
  args:
  - id
  use: This is an example command
- name: cmd:two-args
  args:
  - id
  - name
  use: This is an example command
enums:
- name: user_type
  type: string
  values:
  - user
  - admin
  - super
resources:
- name: organization
  has_many:
  - users
  fields:
  - name: name
    type: string
    required: true
- name: user
  belongs_to:
  - organization
  has_many:
  - tokens
  fields:
  - name: type
    enum: user_type
    required: true
  - name: name
    type: string
    required: true
  - name: email
    type: string
    required: true
    unique: true
- name: token
  belongs_to:
  - user
  fields:
  - name: key
    type: string
    required: true
    unique: true
  - name: expires
    type: time.Time
    required: true
  actions:
  - create
  - delete
  custom:
  - refresh
`
	return []byte(input)
}

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
		Ops: Ops{
			Workdir:      "test-example",
			AppHTTPAlias: "test-example.local",
			MainDatabase: Database{
				Host:     "db",
				Port:     "3306",
				Username: "root",
				Password: "secret",
				Name:     "test_example",
				Debug:    "1",
			},
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
		Enums: []Enum{
			{
				Name: "user_type",
				Type: "string",
				Values: []string{
					"user",
					"admin",
					"super",
				},
			},
		},
		Resources: []Resource{
			{
				Name:    "organization",
				HasMany: []string{"users"},
				Fields: []*ResourceField{
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
				Fields: []*ResourceField{
					{
						Name:     "type",
						Type:     "user_type",
						Enum:     "user_type",
						Required: true,
						Unique:   false,
					},
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
				Fields: []*ResourceField{
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
  app_http_alias: test-example.local
  database:
    host: db
    port: "3306"
    username: root
    password: secret
    name: test_example
    debug: "1"
packages:
  src:
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example
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
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/app
  cmd:
    reference: cmd
    name:
      space: cmd
      snake: cmd
      kebob: cmd
      exported: Cmd
      unexported: cmd
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/cmd
  database:
    reference: db
    name:
      space: db
      snake: db
      kebob: db
      exported: Db
      unexported: db
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/db
  http:
    reference: http
    name:
      space: http
      snake: http
      kebob: http
      exported: Http
      unexported: http
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/http
  repos:
    reference: repos
    name:
      space: repos
      snake: repos
      kebob: repos
      exported: Repos
      unexported: repos
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/repos
  models:
    reference: models
    name:
      space: models
      snake: models
      kebob: models
      exported: Models
      unexported: models
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/models
  migrations:
    reference: migrations
    name:
      space: migrations
      snake: migrations
      kebob: migrations
      exported: Migrations
      unexported: migrations
    path:
      base: github.com/68696c6c/test-example/db
      full: github.com/68696c6c/test-example/db/migrations
  seeders:
    reference: seeders
    name:
      space: seeders
      snake: seeders
      kebob: seeders
      exported: Seeders
      unexported: seeders
    path:
      base: github.com/68696c6c/test-example/db
      full: github.com/68696c6c/test-example/db/seeders
  domains:
    reference: app
    name:
      space: app
      snake: app
      kebob: app
      exported: App
      unexported: app
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/app
  enums:
    reference: enums
    name:
      space: enums
      snake: enums
      kebob: enums
      exported: Enums
      unexported: enums
    path:
      base: github.com/68696c6c/test-example/app
      full: github.com/68696c6c/test-example/app/enums
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
enums:
  user_type:
    enum_type: string
    inflection:
      single:
        space: user type
        snake: user_type
        kebob: user-type
        exported: UserType
        unexported: userType
      plural:
        space: user types
        snake: user_types
        kebob: user-types
        exported: UserTypes
        unexported: userTypes
    type_data:
      reference: enums.UserType
      package: enums
      name: UserType
      receiver_name: u
      data_type: string
    values:
    - user
    - admin
    - super
resources:
- key:
    resource: organization
  inflection:
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
      field: name
    name:
      space: name
      snake: name
      kebob: name
      exported: Name
      unexported: name
    type_data:
      reference: string
      name: string
      data_type: VARCHAR
    is_required: true
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
  fields_meta:
    goat:
    - key:
        resource: organization
        field: id
      name:
        space: id
        snake: id
        kebob: id
        exported: Id
        unexported: id
      type_data:
        reference: goat.ID
        package: goat
        name: ID
        receiver_name: i
        data_type: BINARY(16) NOT NULL
      is_primary: true
    - key:
        resource: organization
        field: created_at
      name:
        space: created at
        snake: created_at
        kebob: created-at
        exported: CreatedAt
        unexported: createdAt
      type_data:
        reference: time.Time
        package: time
        name: Time
        receiver_name: t
        data_type: NOT NULL DEFAULT CURRENT_TIMESTAMP
    - key:
        resource: organization
        field: updated_at
      name:
        space: updated at
        snake: updated_at
        kebob: updated-at
        exported: UpdatedAt
        unexported: updatedAt
      type_data:
        reference: time.Time
        package: time
        name: Time
        receiver_name: t
        data_type: NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
        is_pointer: true
    - key:
        resource: organization
        field: deleted_at
      name:
        space: deleted at
        snake: deleted_at
        kebob: deleted-at
        exported: DeletedAt
        unexported: deletedAt
      type_data:
        reference: time.Time
        package: time
        name: Time
        receiver_name: t
        data_type: NULL DEFAULT NULL
        is_pointer: true
    - key:
        resource: organization
        field: name
      name:
        space: name
        snake: name
        kebob: name
        exported: Name
        unexported: name
      type_data:
        reference: string
        name: string
        data_type: VARCHAR
      is_required: true
    model:
    - key:
        resource: organization
        field: name
      name:
        space: name
        snake: name
        kebob: name
        exported: Name
        unexported: name
      type_data:
        reference: string
        name: string
        data_type: VARCHAR
      is_required: true
    has_many:
    - key:
        resource: organization
        field: Users
      relation:
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
      name:
        space: users
        snake: users
        kebob: users
        exported: Users
        unexported: users
      type_data:
        reference: users.User
        package: users
        name: User
        receiver_name: u
        is_pointer: true
        is_slice: true
- key:
    resource: user
  inflection:
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
      field: organization_id
    name:
      space: organization id
      snake: organization_id
      kebob: organization-id
      exported: OrganizationId
      unexported: organizationId
    type_data:
      reference: goat.ID
      package: goat
      name: ID
      receiver_name: i
      data_type: BINARY(16) NOT NULL
  - key:
      resource: user
      field: type
    name:
      space: type
      snake: type
      kebob: type
      exported: Type
      unexported: type
    type_data:
      reference: enums.UserType
      package: enums
      name: UserType
      receiver_name: u
      data_type: string
    is_required: true
  - key:
      resource: user
      field: name
    name:
      space: name
      snake: name
      kebob: name
      exported: Name
      unexported: name
    type_data:
      reference: string
      name: string
      data_type: VARCHAR
    is_required: true
  - key:
      resource: user
      field: email
    name:
      space: email
      snake: email
      kebob: email
      exported: Email
      unexported: email
    type_data:
      reference: string
      name: string
      data_type: VARCHAR
    is_required: true
    is_unique: true
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
  fields_meta:
    goat:
    - key:
        resource: user
        field: id
      name:
        space: id
        snake: id
        kebob: id
        exported: Id
        unexported: id
      type_data:
        reference: goat.ID
        package: goat
        name: ID
        receiver_name: i
        data_type: BINARY(16) NOT NULL
      is_primary: true
    - key:
        resource: user
        field: created_at
      name:
        space: created at
        snake: created_at
        kebob: created-at
        exported: CreatedAt
        unexported: createdAt
      type_data:
        reference: time.Time
        package: time
        name: Time
        receiver_name: t
        data_type: NOT NULL DEFAULT CURRENT_TIMESTAMP
    - key:
        resource: user
        field: updated_at
      name:
        space: updated at
        snake: updated_at
        kebob: updated-at
        exported: UpdatedAt
        unexported: updatedAt
      type_data:
        reference: time.Time
        package: time
        name: Time
        receiver_name: t
        data_type: NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
        is_pointer: true
    - key:
        resource: user
        field: deleted_at
      name:
        space: deleted at
        snake: deleted_at
        kebob: deleted-at
        exported: DeletedAt
        unexported: deletedAt
      type_data:
        reference: time.Time
        package: time
        name: Time
        receiver_name: t
        data_type: NULL DEFAULT NULL
        is_pointer: true
    - key:
        resource: user
        field: organization_id
      name:
        space: organization id
        snake: organization_id
        kebob: organization-id
        exported: OrganizationId
        unexported: organizationId
      type_data:
        reference: goat.ID
        package: goat
        name: ID
        receiver_name: i
        data_type: BINARY(16) NOT NULL
    - key:
        resource: user
        field: type
      name:
        space: type
        snake: type
        kebob: type
        exported: Type
        unexported: type
      type_data:
        reference: enums.UserType
        package: enums
        name: UserType
        receiver_name: u
        data_type: string
      is_required: true
    - key:
        resource: user
        field: name
      name:
        space: name
        snake: name
        kebob: name
        exported: Name
        unexported: name
      type_data:
        reference: string
        name: string
        data_type: VARCHAR
      is_required: true
    - key:
        resource: user
        field: email
      name:
        space: email
        snake: email
        kebob: email
        exported: Email
        unexported: email
      type_data:
        reference: string
        name: string
        data_type: VARCHAR
      is_required: true
      is_unique: true
    model:
    - key:
        resource: user
        field: organization_id
      name:
        space: organization id
        snake: organization_id
        kebob: organization-id
        exported: OrganizationId
        unexported: organizationId
      type_data:
        reference: goat.ID
        package: goat
        name: ID
        receiver_name: i
        data_type: BINARY(16) NOT NULL
    - key:
        resource: user
        field: type
      name:
        space: type
        snake: type
        kebob: type
        exported: Type
        unexported: type
      type_data:
        reference: enums.UserType
        package: enums
        name: UserType
        receiver_name: u
        data_type: string
      is_required: true
    - key:
        resource: user
        field: name
      name:
        space: name
        snake: name
        kebob: name
        exported: Name
        unexported: name
      type_data:
        reference: string
        name: string
        data_type: VARCHAR
      is_required: true
    - key:
        resource: user
        field: email
      name:
        space: email
        snake: email
        kebob: email
        exported: Email
        unexported: email
      type_data:
        reference: string
        name: string
        data_type: VARCHAR
      is_required: true
      is_unique: true
    belongs_to:
    - key:
        resource: user
        field: organization
      name:
        space: organization
        snake: organization
        kebob: organization
        exported: Organization
        unexported: organization
      type_data:
        reference: organizations.Organization
        package: organizations
        name: Organization
        receiver_name: o
        is_pointer: true
    has_many:
    - key:
        resource: user
        field: Tokens
      relation:
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
      name:
        space: tokens
        snake: tokens
        kebob: tokens
        exported: Tokens
        unexported: tokens
      type_data:
        reference: tokens.Token
        package: tokens
        name: Token
        receiver_name: t
        is_pointer: true
        is_slice: true
    unique:
    - key:
        resource: user
        field: email
      name:
        space: email
        snake: email
        kebob: email
        exported: Email
        unexported: email
      type_data:
        reference: string
        name: string
        data_type: VARCHAR
      is_required: true
      is_unique: true
- key:
    resource: token
  inflection:
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
      field: user_id
    name:
      space: user id
      snake: user_id
      kebob: user-id
      exported: UserId
      unexported: userId
    type_data:
      reference: goat.ID
      package: goat
      name: ID
      receiver_name: i
      data_type: BINARY(16) NOT NULL
  - key:
      resource: token
      field: key
    name:
      space: key
      snake: key
      kebob: key
      exported: Key
      unexported: key
    type_data:
      reference: string
      name: string
      data_type: VARCHAR
    is_required: true
    is_unique: true
  - key:
      resource: token
      field: expires
    name:
      space: expires
      snake: expires
      kebob: expires
      exported: Expires
      unexported: expires
    type_data:
      reference: time.Time
      package: time
      name: Time
    is_required: true
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
  fields_meta:
    goat:
    - key:
        resource: token
        field: id
      name:
        space: id
        snake: id
        kebob: id
        exported: Id
        unexported: id
      type_data:
        reference: goat.ID
        package: goat
        name: ID
        receiver_name: i
        data_type: BINARY(16) NOT NULL
      is_primary: true
    - key:
        resource: token
        field: created_at
      name:
        space: created at
        snake: created_at
        kebob: created-at
        exported: CreatedAt
        unexported: createdAt
      type_data:
        reference: time.Time
        package: time
        name: Time
        receiver_name: t
        data_type: NOT NULL DEFAULT CURRENT_TIMESTAMP
    - key:
        resource: token
        field: updated_at
      name:
        space: updated at
        snake: updated_at
        kebob: updated-at
        exported: UpdatedAt
        unexported: updatedAt
      type_data:
        reference: time.Time
        package: time
        name: Time
        receiver_name: t
        data_type: NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
        is_pointer: true
    - key:
        resource: token
        field: deleted_at
      name:
        space: deleted at
        snake: deleted_at
        kebob: deleted-at
        exported: DeletedAt
        unexported: deletedAt
      type_data:
        reference: time.Time
        package: time
        name: Time
        receiver_name: t
        data_type: NULL DEFAULT NULL
        is_pointer: true
    - key:
        resource: token
        field: user_id
      name:
        space: user id
        snake: user_id
        kebob: user-id
        exported: UserId
        unexported: userId
      type_data:
        reference: goat.ID
        package: goat
        name: ID
        receiver_name: i
        data_type: BINARY(16) NOT NULL
    - key:
        resource: token
        field: key
      name:
        space: key
        snake: key
        kebob: key
        exported: Key
        unexported: key
      type_data:
        reference: string
        name: string
        data_type: VARCHAR
      is_required: true
      is_unique: true
    - key:
        resource: token
        field: expires
      name:
        space: expires
        snake: expires
        kebob: expires
        exported: Expires
        unexported: expires
      type_data:
        reference: time.Time
        package: time
        name: Time
      is_required: true
    model:
    - key:
        resource: token
        field: user_id
      name:
        space: user id
        snake: user_id
        kebob: user-id
        exported: UserId
        unexported: userId
      type_data:
        reference: goat.ID
        package: goat
        name: ID
        receiver_name: i
        data_type: BINARY(16) NOT NULL
    - key:
        resource: token
        field: key
      name:
        space: key
        snake: key
        kebob: key
        exported: Key
        unexported: key
      type_data:
        reference: string
        name: string
        data_type: VARCHAR
      is_required: true
      is_unique: true
    - key:
        resource: token
        field: expires
      name:
        space: expires
        snake: expires
        kebob: expires
        exported: Expires
        unexported: expires
      type_data:
        reference: time.Time
        package: time
        name: Time
      is_required: true
    belongs_to:
    - key:
        resource: token
        field: user
      name:
        space: user
        snake: user
        kebob: user
        exported: User
        unexported: user
      type_data:
        reference: users.User
        package: users
        name: User
        receiver_name: u
        is_pointer: true
    unique:
    - key:
        resource: token
        field: key
      name:
        space: key
        snake: key
        kebob: key
        exported: Key
        unexported: key
      type_data:
        reference: string
        name: string
        data_type: VARCHAR
      is_required: true
      is_unique: true
`
