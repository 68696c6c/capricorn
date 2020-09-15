package src

import (
	"testing"

	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/spec"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"

	"github.com/stretchr/testify/assert"
)

func TestSRC_NewSRCDDD(t *testing.T) {
	f := spec.GetFixtureSpec()
	m := module.NewModuleFromSpec(f, true)

	result := NewSRCDDD(m, "/base/path/test-example")
	resultYAML := result.String()
	// println(resultYAML)

	assert.Equal(t, FixtureSRCYAML, resultYAML)
}

func Test_mergeImports(t *testing.T) {
	stack := golang.Imports{
		Standard: []string{"one"},
		App:      []string{"one", "two"},
		Vendor:   []string{"one", "two", "three"},
	}

	additional := golang.Imports{
		Standard: []string{},
		App:      []string{},
		Vendor:   []string{"one", "two", "three", "four"},
	}

	stack = mergeImports(stack, additional)

	assert.Len(t, stack.Standard, 1)
	assert.Len(t, stack.App, 2)
	assert.Len(t, stack.Vendor, 4)
}

const FixtureSRCYAML = `package:
  reference: ""
  name:
    space: ""
    snake: ""
    kebob: ""
    exported: ""
    unexported: ""
  path:
    base: github.com/68696c6c/test-example
    full: github.com/68696c6c/test-example
path:
  base: /base/path/test-example
  full: /base/path/test-example
app:
  domains:
  - controller:
      name:
        base: controller
        full: controller.go
      path:
        base: github.com/68696c6c/test-example/app/organizations
        full: github.com/68696c6c/test-example/app/organizations/controller.go
      package:
        reference: organizations
        name:
          space: organizations
          snake: organizations
          kebob: organizations
          exported: Organizations
          unexported: organizations
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/organizations
      imports:
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat/query
        - github.com/gin-gonic/gin
      structs:
      - name: Organizations
        fields:
        - name: repo
          type: Repo
        - name: errors
          type: goat.ErrorHandler
      - name: CreateRequest
        fields:
        - name: Organization
      - name: UpdateRequest
        fields:
        - name: Organization
      - name: Response
        fields:
        - name: Organization
      - name: ListResponse
        fields:
        - name: Data
          type: '[]*Organization'
          tags:
          - key: json
            values:
            - data
        - type: query.Pagination
          tags:
          - key: json
            values:
            - pagination
      functions:
      - name: NewController
        imports:
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: repo
          type: Repo
        - name: errors
          type: goat.ErrorHandler
        return_values:
        - type: Controller
        body: "\n\treturn Controller{\n\t\trepo: repo,\n\t\terrors: errors,\n\t}\n"
      - name: List
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/68696c6c/goat/query
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\tq := query.NewQueryBuilder(cx)\n\n\tresult, errs := c.repo.Filter(q)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx, errs, \"failed to list organizations\",
          goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondData(cx, ListResponse{\n\t\tData:
          result,\n\t\tPagination: q.Pagination,\n\t})\n"
      - name: View
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\ti := cx.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(cx, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\tm,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(cx, \"organization does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(cx, errs, \"failed to get organization\",
          goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\tgoat.RespondData(cx,
          Response{m})\n"
      - name: Create
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\treq, ok := goat.GetRequest(cx).(*CreateRequest)\n\tif !ok {\n\t\tc.errors.HandleMessage(c,
          \"failed to get request\", goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\tm
          := req.Organization\n\terrs := c.repo.Save(&m)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx,
          errs, \"failed to save organization\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondCreated(cx,
          Response{m})\n"
      - name: Update
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\ti := cx.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(cx, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO replace this block with an existence validator and build \"not found\"
          handling into the repo.\n\t_, errs := c.repo.GetByID(id)\n\tif len(errs)
          > 0 {\n\t\tif goat.RecordNotFound(errs) {\n\t\t\tc.errors.HandleMessage(cx,
          \"organization does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(cx, errs, \"failed to get organization\",
          goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\treq, ok := goat.GetRequest(cx).(*UpdateRequest)\n\tif
          !ok {\n\t\tc.errors.HandleMessage(cx, \"failed to get request\", goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\terrs
          = c.repo.Save(&req.Organization)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx,
          errs, \"failed to save organization\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondCreated(cx,
          Response{req.Organization})\n"
      - name: Delete
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\ti := cx.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(cx, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\tm,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(cx, \"organization does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(cx, errs, \"failed to get organization\",
          goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\t// @TODO generate
          model factories.\n\t// @TODO generate model validators.\n\terrs = c.repo.Delete(&m)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx, errs, \"failed to delete
          organization\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondValid(cx)\n"
    repo:
      name:
        base: repo
        full: repo.go
      path:
        base: github.com/68696c6c/test-example/app/organizations
        full: github.com/68696c6c/test-example/app/organizations/repo.go
      package:
        reference: organizations
        name:
          space: organizations
          snake: organizations
          kebob: organizations
          exported: Organizations
          unexported: organizations
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/organizations
      imports:
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
      interfaces:
      - name: Repo
        functions:
        - name: NewRepo
          imports:
            vendor:
            - github.com/jinzhu/gorm
          arguments:
          - name: d
            type: '*gorm.DB'
          return_values:
          - type: Repo
          body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
        - name: getBaseQuery
          imports:
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          return_values:
          - type: '*gorm.DB'
          receiver:
            name: r
            type: RepoGorm
          body: return r.db.Model(&Organization{})
        - name: getFilteredQuery
          imports:
            vendor:
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - type: '*gorm.DB'
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
            {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
        - name: applyPaginationToQuery
          imports:
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif
            err != nil {\n\t\treturn errors.Wrap(err, \"failed to set organization
            query pagination\")\n\t}\n\treturn nil\n"
        - name: Filter
          imports:
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: id
            type: goat.ID
          return_values:
          - name: result
            type: '[]*Organization'
          - name: err
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tdataQuery, err := r.getFilteredQuery(q)\n\tif err != nil {\n\t\treturn
            result, errors.Wrap(err, \"failed to build filter sites query\")\n\t}\n\n\terrs
            := dataQuery.Find(&result).GetErrors()\n\tif len(errs) > 0 && goat.ErrorsBesidesRecordNotFound(errs)
            {\n\t\terr := goat.ErrorsToError(errs)\n\t\treturn result, errors.Wrap(err,
            \"failed to execute filter sites data query\")\n\t}\n\n\tif err := r.applyPaginationToQuery(q);
            err != nil {\n\t\treturn result, err\n\t}\n\n\treturn result, nil\n"
        - name: GetByID
          imports:
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: id
            type: goat.ID
          return_values:
          - type: Organization
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tm := Organization{\n\t\tModel: goat.Model{\n\t\t\tID: id,\n\t\t},\n\t}\n\terrs
            := r.db.First(&m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn m, goat.ErrorsToError(errs)\n\t}\n\treturn
            m, nil\n"
        - name: Save
          imports:
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*Organization'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
            else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0
            {\n\t\treturn goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
        - name: Save
          imports:
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*Organization'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
            else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0
            {\n\t\treturn goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
        - name: Delete
          imports:
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*Organization'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\terrs :=  r.db.Delete(m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn
            goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      structs:
      - name: Organization
        fields:
        - name: db
          type: '*gorm.DB'
      functions:
      - name: NewRepo
        imports:
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: d
          type: '*gorm.DB'
        return_values:
        - type: Repo
        body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
      - name: getBaseQuery
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        return_values:
        - type: '*gorm.DB'
        receiver:
          name: r
          type: RepoGorm
        body: return r.db.Model(&Organization{})
      - name: getFilteredQuery
        imports:
          vendor:
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - type: '*gorm.DB'
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
          {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
      - name: applyPaginationToQuery
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif err
          != nil {\n\t\treturn errors.Wrap(err, \"failed to set organization query
          pagination\")\n\t}\n\treturn nil\n"
      - name: Filter
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: id
          type: goat.ID
        return_values:
        - name: result
          type: '[]*Organization'
        - name: err
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tdataQuery, err := r.getFilteredQuery(q)\n\tif err != nil {\n\t\treturn
          result, errors.Wrap(err, \"failed to build filter sites query\")\n\t}\n\n\terrs
          := dataQuery.Find(&result).GetErrors()\n\tif len(errs) > 0 && goat.ErrorsBesidesRecordNotFound(errs)
          {\n\t\terr := goat.ErrorsToError(errs)\n\t\treturn result, errors.Wrap(err,
          \"failed to execute filter sites data query\")\n\t}\n\n\tif err := r.applyPaginationToQuery(q);
          err != nil {\n\t\treturn result, err\n\t}\n\n\treturn result, nil\n"
      - name: GetByID
        imports:
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: id
          type: goat.ID
        return_values:
        - type: Organization
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tm := Organization{\n\t\tModel: goat.Model{\n\t\t\tID: id,\n\t\t},\n\t}\n\terrs
          := r.db.First(&m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn m, goat.ErrorsToError(errs)\n\t}\n\treturn
          m, nil\n"
      - name: Save
        imports:
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*Organization'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
          else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      - name: Save
        imports:
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*Organization'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
          else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      - name: Delete
        imports:
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*Organization'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terrs :=  r.db.Delete(m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
    model:
      name:
        base: model
        full: model.go
      path:
        base: github.com/68696c6c/test-example/app/organizations
        full: github.com/68696c6c/test-example/app/organizations/model.go
      package:
        reference: organizations
        name:
          space: organizations
          snake: organizations
          kebob: organizations
          exported: Organizations
          unexported: organizations
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/organizations
      structs:
      - name: Organization
      functions:
      - {}
    validator:
      name:
        base: validator
        full: validator.go
      path:
        base: github.com/68696c6c/test-example/app/organizations
        full: github.com/68696c6c/test-example/app/organizations/validator.go
      package:
        reference: organizations
        name:
          space: organizations
          snake: organizations
          kebob: organizations
          exported: Organizations
          unexported: organizations
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/organizations
  - controller:
      name:
        base: controller
        full: controller.go
      path:
        base: github.com/68696c6c/test-example/app/users
        full: github.com/68696c6c/test-example/app/users/controller.go
      package:
        reference: users
        name:
          space: users
          snake: users
          kebob: users
          exported: Users
          unexported: users
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/users
      imports:
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat/query
        - github.com/gin-gonic/gin
      structs:
      - name: Users
        fields:
        - name: repo
          type: Repo
        - name: errors
          type: goat.ErrorHandler
      - name: CreateRequest
        fields:
        - name: User
      - name: UpdateRequest
        fields:
        - name: User
      - name: Response
        fields:
        - name: User
      - name: ListResponse
        fields:
        - name: Data
          type: '[]*User'
          tags:
          - key: json
            values:
            - data
        - type: query.Pagination
          tags:
          - key: json
            values:
            - pagination
      functions:
      - name: NewController
        imports:
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: repo
          type: Repo
        - name: errors
          type: goat.ErrorHandler
        return_values:
        - type: Controller
        body: "\n\treturn Controller{\n\t\trepo: repo,\n\t\terrors: errors,\n\t}\n"
      - name: List
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/68696c6c/goat/query
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\tq := query.NewQueryBuilder(cx)\n\n\tresult, errs := c.repo.Filter(q)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx, errs, \"failed to list users\",
          goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondData(cx, ListResponse{\n\t\tData:
          result,\n\t\tPagination: q.Pagination,\n\t})\n"
      - name: View
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\ti := cx.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(cx, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\tm,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(cx, \"user does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(cx, errs, \"failed to get user\", goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\tgoat.RespondData(cx,
          Response{m})\n"
      - name: Create
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\treq, ok := goat.GetRequest(cx).(*CreateRequest)\n\tif !ok {\n\t\tc.errors.HandleMessage(c,
          \"failed to get request\", goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\tm
          := req.User\n\terrs := c.repo.Save(&m)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx,
          errs, \"failed to save user\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondCreated(cx,
          Response{m})\n"
      - name: Update
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\ti := cx.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(cx, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO replace this block with an existence validator and build \"not found\"
          handling into the repo.\n\t_, errs := c.repo.GetByID(id)\n\tif len(errs)
          > 0 {\n\t\tif goat.RecordNotFound(errs) {\n\t\t\tc.errors.HandleMessage(cx,
          \"user does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(cx, errs, \"failed to get user\", goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\treq,
          ok := goat.GetRequest(cx).(*UpdateRequest)\n\tif !ok {\n\t\tc.errors.HandleMessage(cx,
          \"failed to get request\", goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\terrs
          = c.repo.Save(&req.User)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx,
          errs, \"failed to save user\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondCreated(cx,
          Response{req.User})\n"
      - name: Delete
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\ti := cx.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(cx, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\tm,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(cx, \"user does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(cx, errs, \"failed to get user\", goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\terrs
          = c.repo.Delete(&m)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx,
          errs, \"failed to delete user\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondValid(cx)\n"
    repo:
      name:
        base: repo
        full: repo.go
      path:
        base: github.com/68696c6c/test-example/app/users
        full: github.com/68696c6c/test-example/app/users/repo.go
      package:
        reference: users
        name:
          space: users
          snake: users
          kebob: users
          exported: Users
          unexported: users
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/users
      imports:
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
      interfaces:
      - name: Repo
        functions:
        - name: NewRepo
          imports:
            vendor:
            - github.com/jinzhu/gorm
          arguments:
          - name: d
            type: '*gorm.DB'
          return_values:
          - type: Repo
          body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
        - name: getBaseQuery
          imports:
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          return_values:
          - type: '*gorm.DB'
          receiver:
            name: r
            type: RepoGorm
          body: return r.db.Model(&User{})
        - name: getFilteredQuery
          imports:
            vendor:
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - type: '*gorm.DB'
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
            {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
        - name: applyPaginationToQuery
          imports:
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif
            err != nil {\n\t\treturn errors.Wrap(err, \"failed to set user query pagination\")\n\t}\n\treturn
            nil\n"
        - name: Filter
          imports:
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: id
            type: goat.ID
          return_values:
          - name: result
            type: '[]*User'
          - name: err
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tdataQuery, err := r.getFilteredQuery(q)\n\tif err != nil {\n\t\treturn
            result, errors.Wrap(err, \"failed to build filter sites query\")\n\t}\n\n\terrs
            := dataQuery.Find(&result).GetErrors()\n\tif len(errs) > 0 && goat.ErrorsBesidesRecordNotFound(errs)
            {\n\t\terr := goat.ErrorsToError(errs)\n\t\treturn result, errors.Wrap(err,
            \"failed to execute filter sites data query\")\n\t}\n\n\tif err := r.applyPaginationToQuery(q);
            err != nil {\n\t\treturn result, err\n\t}\n\n\treturn result, nil\n"
        - name: GetByID
          imports:
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: id
            type: goat.ID
          return_values:
          - type: User
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tm := User{\n\t\tModel: goat.Model{\n\t\t\tID: id,\n\t\t},\n\t}\n\terrs
            := r.db.First(&m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn m, goat.ErrorsToError(errs)\n\t}\n\treturn
            m, nil\n"
        - name: Save
          imports:
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*User'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
            else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0
            {\n\t\treturn goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
        - name: Save
          imports:
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*User'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
            else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0
            {\n\t\treturn goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
        - name: Delete
          imports:
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*User'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\terrs :=  r.db.Delete(m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn
            goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      structs:
      - name: User
        fields:
        - name: db
          type: '*gorm.DB'
      functions:
      - name: NewRepo
        imports:
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: d
          type: '*gorm.DB'
        return_values:
        - type: Repo
        body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
      - name: getBaseQuery
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        return_values:
        - type: '*gorm.DB'
        receiver:
          name: r
          type: RepoGorm
        body: return r.db.Model(&User{})
      - name: getFilteredQuery
        imports:
          vendor:
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - type: '*gorm.DB'
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
          {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
      - name: applyPaginationToQuery
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif err
          != nil {\n\t\treturn errors.Wrap(err, \"failed to set user query pagination\")\n\t}\n\treturn
          nil\n"
      - name: Filter
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: id
          type: goat.ID
        return_values:
        - name: result
          type: '[]*User'
        - name: err
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tdataQuery, err := r.getFilteredQuery(q)\n\tif err != nil {\n\t\treturn
          result, errors.Wrap(err, \"failed to build filter sites query\")\n\t}\n\n\terrs
          := dataQuery.Find(&result).GetErrors()\n\tif len(errs) > 0 && goat.ErrorsBesidesRecordNotFound(errs)
          {\n\t\terr := goat.ErrorsToError(errs)\n\t\treturn result, errors.Wrap(err,
          \"failed to execute filter sites data query\")\n\t}\n\n\tif err := r.applyPaginationToQuery(q);
          err != nil {\n\t\treturn result, err\n\t}\n\n\treturn result, nil\n"
      - name: GetByID
        imports:
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: id
          type: goat.ID
        return_values:
        - type: User
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tm := User{\n\t\tModel: goat.Model{\n\t\t\tID: id,\n\t\t},\n\t}\n\terrs
          := r.db.First(&m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn m, goat.ErrorsToError(errs)\n\t}\n\treturn
          m, nil\n"
      - name: Save
        imports:
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*User'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
          else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      - name: Save
        imports:
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*User'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
          else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      - name: Delete
        imports:
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*User'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terrs :=  r.db.Delete(m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
    model:
      name:
        base: model
        full: model.go
      path:
        base: github.com/68696c6c/test-example/app/users
        full: github.com/68696c6c/test-example/app/users/model.go
      package:
        reference: users
        name:
          space: users
          snake: users
          kebob: users
          exported: Users
          unexported: users
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/users
      structs:
      - name: User
      functions:
      - {}
    validator:
      name:
        base: validator
        full: validator.go
      path:
        base: github.com/68696c6c/test-example/app/users
        full: github.com/68696c6c/test-example/app/users/validator.go
      package:
        reference: users
        name:
          space: users
          snake: users
          kebob: users
          exported: Users
          unexported: users
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/users
  - controller:
      name:
        base: controller
        full: controller.go
      path:
        base: github.com/68696c6c/test-example/app/tokens
        full: github.com/68696c6c/test-example/app/tokens/controller.go
      package:
        reference: tokens
        name:
          space: tokens
          snake: tokens
          kebob: tokens
          exported: Tokens
          unexported: tokens
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/tokens
      imports:
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
      structs:
      - name: Tokens
        fields:
        - name: repo
          type: Repo
        - name: errors
          type: goat.ErrorHandler
      - name: CreateRequest
        fields:
        - name: Token
      - name: Response
        fields:
        - name: Token
      functions:
      - name: NewController
        imports:
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: repo
          type: Repo
        - name: errors
          type: goat.ErrorHandler
        return_values:
        - type: Controller
        body: "\n\treturn Controller{\n\t\trepo: repo,\n\t\terrors: errors,\n\t}\n"
      - name: Create
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\treq, ok := goat.GetRequest(cx).(*CreateRequest)\n\tif !ok {\n\t\tc.errors.HandleMessage(c,
          \"failed to get request\", goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\tm
          := req.Token\n\terrs := c.repo.Save(&m)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx,
          errs, \"failed to save token\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondCreated(cx,
          Response{m})\n"
      - name: Delete
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        receiver:
          name: c
          type: Controller
        body: "\n\ti := cx.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(cx, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\tm,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(cx, \"token does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(cx, errs, \"failed to get token\",
          goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\t// @TODO generate
          model factories.\n\t// @TODO generate model validators.\n\terrs = c.repo.Delete(&m)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx, errs, \"failed to delete
          token\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondValid(cx)\n"
    repo:
      name:
        base: repo
        full: repo.go
      path:
        base: github.com/68696c6c/test-example/app/tokens
        full: github.com/68696c6c/test-example/app/tokens/repo.go
      package:
        reference: tokens
        name:
          space: tokens
          snake: tokens
          kebob: tokens
          exported: Tokens
          unexported: tokens
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/tokens
      imports:
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
      interfaces:
      - name: Repo
        functions:
        - name: NewRepo
          imports:
            vendor:
            - github.com/jinzhu/gorm
          arguments:
          - name: d
            type: '*gorm.DB'
          return_values:
          - type: Repo
          body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
        - name: getBaseQuery
          imports:
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          return_values:
          - type: '*gorm.DB'
          receiver:
            name: r
            type: RepoGorm
          body: return r.db.Model(&Token{})
        - name: getFilteredQuery
          imports:
            vendor:
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - type: '*gorm.DB'
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
            {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
        - name: applyPaginationToQuery
          imports:
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif
            err != nil {\n\t\treturn errors.Wrap(err, \"failed to set token query
            pagination\")\n\t}\n\treturn nil\n"
        - name: Save
          imports:
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*Token'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
            else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0
            {\n\t\treturn goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
        - name: Delete
          imports:
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*Token'
          return_values:
          - type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\terrs :=  r.db.Delete(m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn
            goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      structs:
      - name: Token
        fields:
        - name: db
          type: '*gorm.DB'
      functions:
      - name: NewRepo
        imports:
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: d
          type: '*gorm.DB'
        return_values:
        - type: Repo
        body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
      - name: getBaseQuery
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        return_values:
        - type: '*gorm.DB'
        receiver:
          name: r
          type: RepoGorm
        body: return r.db.Model(&Token{})
      - name: getFilteredQuery
        imports:
          vendor:
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - type: '*gorm.DB'
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
          {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
      - name: applyPaginationToQuery
        imports:
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif err
          != nil {\n\t\treturn errors.Wrap(err, \"failed to set token query pagination\")\n\t}\n\treturn
          nil\n"
      - name: Save
        imports:
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*Token'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
          else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      - name: Delete
        imports:
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*Token'
        return_values:
        - type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terrs :=  r.db.Delete(m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
    model:
      name:
        base: model
        full: model.go
      path:
        base: github.com/68696c6c/test-example/app/tokens
        full: github.com/68696c6c/test-example/app/tokens/model.go
      package:
        reference: tokens
        name:
          space: tokens
          snake: tokens
          kebob: tokens
          exported: Tokens
          unexported: tokens
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/tokens
      structs:
      - name: Token
      functions:
      - {}
    validator:
      name:
        base: validator
        full: validator.go
      path:
        base: github.com/68696c6c/test-example/app/tokens
        full: github.com/68696c6c/test-example/app/tokens/validator.go
      package:
        reference: tokens
        name:
          space: tokens
          snake: tokens
          kebob: tokens
          exported: Tokens
          unexported: tokens
        path:
          base: github.com/68696c6c/test-example/app
          full: github.com/68696c6c/test-example/app/tokens
main:
  name:
    base: main
    full: main.go
  path:
    base: /base/path/test-example
    full: /base/path/test-example/main.go
  package:
    reference: main
    name:
      space: main
      snake: main
      kebob: main
      exported: Main
      unexported: main
    path:
      base: github.com/68696c6c/test-example
      full: github.com/68696c6c/test-example/main
  imports:
    standard:
    - os
    app:
    - github.com/68696c6c/test-example/cmd
  functions:
  - name: main
    body: "\n\tif err := cmd.Root.Execute(); err != nil {\n\t\tprintln(err)\n\t\tos.Exit(1)\n\t}"
`
