package src

import (
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"testing"

	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/spec"

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
path:
  base: /base/path/test-example
  full: /base/path/test-example/src
app:
  container:
    name:
      base: ""
      full: ""
    path:
      base: ""
      full: ""
    package:
      reference: ""
      name:
        space: ""
        snake: ""
        kebob: ""
        exported: ""
        unexported: ""
      path:
        base: ""
        full: ""
    imports:
      standard: []
      app: []
      vendor: []
    init_function:
      name: ""
      imports:
        standard: []
        app: []
        vendor: []
      arguments: []
      return_values: []
      receiver:
        name: ""
        type: ""
      body: ""
    consts: []
    vars: []
    interfaces: []
    structs: []
    functions: []
  domains:
  - controller:
      name:
        base: controller
        full: controller.go
      path:
        base: github.com/68696c6c/test-example/src/app/organizations
        full: github.com/68696c6c/test-example/src/app/organizations/controller.go
      package:
        reference: organizations
        name:
          space: organizations
          snake: organizations
          kebob: organizations
          exported: Organizations
          unexported: organizations
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/organizations
      imports:
        standard: []
        app: []
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat/query
        - github.com/gin-gonic/gin
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs:
      - name: Organizations
        fields:
        - name: repo
          type: Repo
          tags: []
        - name: errors
          type: goat.ErrorHandler
          tags: []
      - name: CreateRequest
        fields:
        - name: Organization
          type: ""
          tags: []
      - name: UpdateRequest
        fields:
        - name: Organization
          type: ""
          tags: []
      - name: Response
        fields:
        - name: Organization
          type: ""
          tags: []
      - name: ListResponse
        fields:
        - name: Data
          type: '[]*Organization'
          tags:
          - key: json
            values:
            - data
        - name: ""
          type: query.Pagination
          tags:
          - key: json
            values:
            - pagination
      functions:
      - name: NewController
        imports:
          standard: []
          app: []
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: repo
          type: Repo
        - name: errors
          type: goat.ErrorHandler
        return_values:
        - name: ""
          type: Controller
        receiver:
          name: ""
          type: ""
        body: "\n\treturn Controller{\n\t\trepo: repo,\n\t\terrors: errors,\n\t}\n"
      - name: List
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/68696c6c/goat/query
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
        receiver:
          name: c
          type: Controller
        body: "\n\tq := query.NewQueryBuilder(cx)\n\n\tresult, errs := c.repo.Filter(q)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx, errs, \"failed to list organizations\",
          goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondData(cx, ListResponse{\n\t\tData:
          result,\n\t\tPagination: q.Pagination,\n\t})\n"
      - name: View
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
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
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
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
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
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
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
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
    controller_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    repo:
      name:
        base: repo
        full: repo.go
      path:
        base: github.com/68696c6c/test-example/src/app/organizations
        full: github.com/68696c6c/test-example/src/app/organizations/repo.go
      package:
        reference: organizations
        name:
          space: organizations
          snake: organizations
          kebob: organizations
          exported: Organizations
          unexported: organizations
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/organizations
      imports:
        standard: []
        app: []
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces:
      - name: Repo
        functions:
        - name: NewRepo
          imports:
            standard: []
            app: []
            vendor:
            - github.com/jinzhu/gorm
          arguments:
          - name: d
            type: '*gorm.DB'
          return_values:
          - name: ""
            type: Repo
          receiver:
            name: ""
            type: ""
          body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
        - name: getBaseQuery
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments: []
          return_values:
          - name: ""
            type: '*gorm.DB'
          receiver:
            name: r
            type: RepoGorm
          body: return r.db.Model(&Organization{})
        - name: getFilteredQuery
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - name: ""
            type: '*gorm.DB'
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
            {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
        - name: applyPaginationToQuery
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif
            err != nil {\n\t\treturn errors.Wrap(err, \"failed to set organization
            query pagination\")\n\t}\n\treturn nil\n"
        - name: Filter
          imports:
            standard: []
            app: []
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
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: id
            type: goat.ID
          return_values:
          - name: ""
            type: Organization
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tm := Organization{\n\t\tModel: goat.Model{\n\t\t\tID: id,\n\t\t},\n\t}\n\terrs
            := r.db.First(&m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn m, goat.ErrorsToError(errs)\n\t}\n\treturn
            m, nil\n"
        - name: Save
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*Organization'
          return_values:
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
            else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0
            {\n\t\treturn goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
        - name: Save
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*Organization'
          return_values:
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
            else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0
            {\n\t\treturn goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
        - name: Delete
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*Organization'
          return_values:
          - name: ""
            type: error
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
          tags: []
      functions:
      - name: NewRepo
        imports:
          standard: []
          app: []
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: d
          type: '*gorm.DB'
        return_values:
        - name: ""
          type: Repo
        receiver:
          name: ""
          type: ""
        body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
      - name: getBaseQuery
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments: []
        return_values:
        - name: ""
          type: '*gorm.DB'
        receiver:
          name: r
          type: RepoGorm
        body: return r.db.Model(&Organization{})
      - name: getFilteredQuery
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - name: ""
          type: '*gorm.DB'
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
          {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
      - name: applyPaginationToQuery
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif err
          != nil {\n\t\treturn errors.Wrap(err, \"failed to set organization query
          pagination\")\n\t}\n\treturn nil\n"
      - name: Filter
        imports:
          standard: []
          app: []
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
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: id
          type: goat.ID
        return_values:
        - name: ""
          type: Organization
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tm := Organization{\n\t\tModel: goat.Model{\n\t\t\tID: id,\n\t\t},\n\t}\n\terrs
          := r.db.First(&m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn m, goat.ErrorsToError(errs)\n\t}\n\treturn
          m, nil\n"
      - name: Save
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*Organization'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
          else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      - name: Save
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*Organization'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
          else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      - name: Delete
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*Organization'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terrs :=  r.db.Delete(m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
    repo_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    model:
      name:
        base: model
        full: model.go
      path:
        base: github.com/68696c6c/test-example/src/app/organizations
        full: github.com/68696c6c/test-example/src/app/organizations/model.go
      package:
        reference: organizations
        name:
          space: organizations
          snake: organizations
          kebob: organizations
          exported: Organizations
          unexported: organizations
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/organizations
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs:
      - name: Organization
        fields: []
      functions:
      - name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
    model_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    service:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    service_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    validator:
      name:
        base: validator
        full: validator.go
      path:
        base: github.com/68696c6c/test-example/src/app/organizations
        full: github.com/68696c6c/test-example/src/app/organizations/validator.go
      package:
        reference: organizations
        name:
          space: organizations
          snake: organizations
          kebob: organizations
          exported: Organizations
          unexported: organizations
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/organizations
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    validator_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
  - controller:
      name:
        base: controller
        full: controller.go
      path:
        base: github.com/68696c6c/test-example/src/app/users
        full: github.com/68696c6c/test-example/src/app/users/controller.go
      package:
        reference: users
        name:
          space: users
          snake: users
          kebob: users
          exported: Users
          unexported: users
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/users
      imports:
        standard: []
        app: []
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat/query
        - github.com/gin-gonic/gin
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs:
      - name: Users
        fields:
        - name: repo
          type: Repo
          tags: []
        - name: errors
          type: goat.ErrorHandler
          tags: []
      - name: CreateRequest
        fields:
        - name: User
          type: ""
          tags: []
      - name: UpdateRequest
        fields:
        - name: User
          type: ""
          tags: []
      - name: Response
        fields:
        - name: User
          type: ""
          tags: []
      - name: ListResponse
        fields:
        - name: Data
          type: '[]*User'
          tags:
          - key: json
            values:
            - data
        - name: ""
          type: query.Pagination
          tags:
          - key: json
            values:
            - pagination
      functions:
      - name: NewController
        imports:
          standard: []
          app: []
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: repo
          type: Repo
        - name: errors
          type: goat.ErrorHandler
        return_values:
        - name: ""
          type: Controller
        receiver:
          name: ""
          type: ""
        body: "\n\treturn Controller{\n\t\trepo: repo,\n\t\terrors: errors,\n\t}\n"
      - name: List
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/68696c6c/goat/query
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
        receiver:
          name: c
          type: Controller
        body: "\n\tq := query.NewQueryBuilder(cx)\n\n\tresult, errs := c.repo.Filter(q)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(cx, errs, \"failed to list users\",
          goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondData(cx, ListResponse{\n\t\tData:
          result,\n\t\tPagination: q.Pagination,\n\t})\n"
      - name: View
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
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
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
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
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
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
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
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
    controller_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    repo:
      name:
        base: repo
        full: repo.go
      path:
        base: github.com/68696c6c/test-example/src/app/users
        full: github.com/68696c6c/test-example/src/app/users/repo.go
      package:
        reference: users
        name:
          space: users
          snake: users
          kebob: users
          exported: Users
          unexported: users
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/users
      imports:
        standard: []
        app: []
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces:
      - name: Repo
        functions:
        - name: NewRepo
          imports:
            standard: []
            app: []
            vendor:
            - github.com/jinzhu/gorm
          arguments:
          - name: d
            type: '*gorm.DB'
          return_values:
          - name: ""
            type: Repo
          receiver:
            name: ""
            type: ""
          body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
        - name: getBaseQuery
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments: []
          return_values:
          - name: ""
            type: '*gorm.DB'
          receiver:
            name: r
            type: RepoGorm
          body: return r.db.Model(&User{})
        - name: getFilteredQuery
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - name: ""
            type: '*gorm.DB'
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
            {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
        - name: applyPaginationToQuery
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif
            err != nil {\n\t\treturn errors.Wrap(err, \"failed to set user query pagination\")\n\t}\n\treturn
            nil\n"
        - name: Filter
          imports:
            standard: []
            app: []
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
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: id
            type: goat.ID
          return_values:
          - name: ""
            type: User
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tm := User{\n\t\tModel: goat.Model{\n\t\t\tID: id,\n\t\t},\n\t}\n\terrs
            := r.db.First(&m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn m, goat.ErrorsToError(errs)\n\t}\n\treturn
            m, nil\n"
        - name: Save
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*User'
          return_values:
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
            else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0
            {\n\t\treturn goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
        - name: Save
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*User'
          return_values:
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
            else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0
            {\n\t\treturn goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
        - name: Delete
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*User'
          return_values:
          - name: ""
            type: error
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
          tags: []
      functions:
      - name: NewRepo
        imports:
          standard: []
          app: []
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: d
          type: '*gorm.DB'
        return_values:
        - name: ""
          type: Repo
        receiver:
          name: ""
          type: ""
        body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
      - name: getBaseQuery
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments: []
        return_values:
        - name: ""
          type: '*gorm.DB'
        receiver:
          name: r
          type: RepoGorm
        body: return r.db.Model(&User{})
      - name: getFilteredQuery
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - name: ""
          type: '*gorm.DB'
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
          {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
      - name: applyPaginationToQuery
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif err
          != nil {\n\t\treturn errors.Wrap(err, \"failed to set user query pagination\")\n\t}\n\treturn
          nil\n"
      - name: Filter
        imports:
          standard: []
          app: []
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
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: id
          type: goat.ID
        return_values:
        - name: ""
          type: User
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tm := User{\n\t\tModel: goat.Model{\n\t\t\tID: id,\n\t\t},\n\t}\n\terrs
          := r.db.First(&m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn m, goat.ErrorsToError(errs)\n\t}\n\treturn
          m, nil\n"
      - name: Save
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*User'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
          else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      - name: Save
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*User'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
          else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      - name: Delete
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*User'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terrs :=  r.db.Delete(m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
    repo_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    model:
      name:
        base: model
        full: model.go
      path:
        base: github.com/68696c6c/test-example/src/app/users
        full: github.com/68696c6c/test-example/src/app/users/model.go
      package:
        reference: users
        name:
          space: users
          snake: users
          kebob: users
          exported: Users
          unexported: users
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/users
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs:
      - name: User
        fields: []
      functions:
      - name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
    model_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    service:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    service_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    validator:
      name:
        base: validator
        full: validator.go
      path:
        base: github.com/68696c6c/test-example/src/app/users
        full: github.com/68696c6c/test-example/src/app/users/validator.go
      package:
        reference: users
        name:
          space: users
          snake: users
          kebob: users
          exported: Users
          unexported: users
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/users
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    validator_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
  - controller:
      name:
        base: controller
        full: controller.go
      path:
        base: github.com/68696c6c/test-example/src/app/tokens
        full: github.com/68696c6c/test-example/src/app/tokens/controller.go
      package:
        reference: tokens
        name:
          space: tokens
          snake: tokens
          kebob: tokens
          exported: Tokens
          unexported: tokens
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/tokens
      imports:
        standard: []
        app: []
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
        - github.com/68696c6c/goat
        - github.com/gin-gonic/gin
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs:
      - name: Tokens
        fields:
        - name: repo
          type: Repo
          tags: []
        - name: errors
          type: goat.ErrorHandler
          tags: []
      - name: CreateRequest
        fields:
        - name: Token
          type: ""
          tags: []
      - name: Response
        fields:
        - name: Token
          type: ""
          tags: []
      functions:
      - name: NewController
        imports:
          standard: []
          app: []
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: repo
          type: Repo
        - name: errors
          type: goat.ErrorHandler
        return_values:
        - name: ""
          type: Controller
        receiver:
          name: ""
          type: ""
        body: "\n\treturn Controller{\n\t\trepo: repo,\n\t\terrors: errors,\n\t}\n"
      - name: Create
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
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
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/gin-gonic/gin
        arguments:
        - name: cx
          type: '*gin.Context'
        return_values: []
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
    controller_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    repo:
      name:
        base: repo
        full: repo.go
      path:
        base: github.com/68696c6c/test-example/src/app/tokens
        full: github.com/68696c6c/test-example/src/app/tokens/repo.go
      package:
        reference: tokens
        name:
          space: tokens
          snake: tokens
          kebob: tokens
          exported: Tokens
          unexported: tokens
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/tokens
      imports:
        standard: []
        app: []
        vendor:
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/pkg/errors
        - github.com/68696c6c/goat/query
        - github.com/jinzhu/gorm
        - github.com/68696c6c/goat
        - github.com/68696c6c/goat
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces:
      - name: Repo
        functions:
        - name: NewRepo
          imports:
            standard: []
            app: []
            vendor:
            - github.com/jinzhu/gorm
          arguments:
          - name: d
            type: '*gorm.DB'
          return_values:
          - name: ""
            type: Repo
          receiver:
            name: ""
            type: ""
          body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
        - name: getBaseQuery
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments: []
          return_values:
          - name: ""
            type: '*gorm.DB'
          receiver:
            name: r
            type: RepoGorm
          body: return r.db.Model(&Token{})
        - name: getFilteredQuery
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - name: ""
            type: '*gorm.DB'
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
            {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
        - name: applyPaginationToQuery
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
            - github.com/pkg/errors
            - github.com/68696c6c/goat/query
            - github.com/jinzhu/gorm
          arguments:
          - name: q
            type: '*query.Query'
          return_values:
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif
            err != nil {\n\t\treturn errors.Wrap(err, \"failed to set token query
            pagination\")\n\t}\n\treturn nil\n"
        - name: Save
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*Token'
          return_values:
          - name: ""
            type: error
          receiver:
            name: r
            type: RepoGorm
          body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
            else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0
            {\n\t\treturn goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
        - name: Delete
          imports:
            standard: []
            app: []
            vendor:
            - github.com/68696c6c/goat
          arguments:
          - name: m
            type: '*Token'
          return_values:
          - name: ""
            type: error
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
          tags: []
      functions:
      - name: NewRepo
        imports:
          standard: []
          app: []
          vendor:
          - github.com/jinzhu/gorm
        arguments:
        - name: d
          type: '*gorm.DB'
        return_values:
        - name: ""
          type: Repo
        receiver:
          name: ""
          type: ""
        body: "\n\treturn {\n\t\tdb: d,\n\t}\n"
      - name: getBaseQuery
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments: []
        return_values:
        - name: ""
          type: '*gorm.DB'
        receiver:
          name: r
          type: RepoGorm
        body: return r.db.Model(&Token{})
      - name: getFilteredQuery
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - name: ""
          type: '*gorm.DB'
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tresult, err := q.ApplyToGorm(r.getBaseQuery())\n\tif err != nil
          {\n\t\treturn result, err\n\t}\n\treturn result, nil\n"
      - name: applyPaginationToQuery
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
          - github.com/pkg/errors
          - github.com/68696c6c/goat/query
          - github.com/jinzhu/gorm
        arguments:
        - name: q
          type: '*query.Query'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terr := goat.ApplyPaginationToQuery(q, r.getBaseQuery())\n\tif err
          != nil {\n\t\treturn errors.Wrap(err, \"failed to set token query pagination\")\n\t}\n\treturn
          nil\n"
      - name: Save
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*Token'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\tvar errs []error\n\tif m.Model.ID.Valid() {\n\t\terrs = r.db.Save(m).GetErrors()\n\t}
          else {\n\t\terrs = r.db.Create(m).GetErrors()\n\t}\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
      - name: Delete
        imports:
          standard: []
          app: []
          vendor:
          - github.com/68696c6c/goat
        arguments:
        - name: m
          type: '*Token'
        return_values:
        - name: ""
          type: error
        receiver:
          name: r
          type: RepoGorm
        body: "\n\terrs :=  r.db.Delete(m).GetErrors()\n\tif len(errs) > 0 {\n\t\treturn
          goat.ErrorsToError(errs)\n\t}\n\treturn nil\n"
    repo_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    model:
      name:
        base: model
        full: model.go
      path:
        base: github.com/68696c6c/test-example/src/app/tokens
        full: github.com/68696c6c/test-example/src/app/tokens/model.go
      package:
        reference: tokens
        name:
          space: tokens
          snake: tokens
          kebob: tokens
          exported: Tokens
          unexported: tokens
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/tokens
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs:
      - name: Token
        fields: []
      functions:
      - name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
    model_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    service:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    service_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    validator:
      name:
        base: validator
        full: validator.go
      path:
        base: github.com/68696c6c/test-example/src/app/tokens
        full: github.com/68696c6c/test-example/src/app/tokens/validator.go
      package:
        reference: tokens
        name:
          space: tokens
          snake: tokens
          kebob: tokens
          exported: Tokens
          unexported: tokens
        path:
          base: github.com/68696c6c/test-example/src/app
          full: github.com/68696c6c/test-example/src/app/tokens
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
    validator_test:
      name:
        base: ""
        full: ""
      path:
        base: ""
        full: ""
      package:
        reference: ""
        name:
          space: ""
          snake: ""
          kebob: ""
          exported: ""
          unexported: ""
        path:
          base: ""
          full: ""
      imports:
        standard: []
        app: []
        vendor: []
      init_function:
        name: ""
        imports:
          standard: []
          app: []
          vendor: []
        arguments: []
        return_values: []
        receiver:
          name: ""
          type: ""
        body: ""
      consts: []
      vars: []
      interfaces: []
      structs: []
      functions: []
cmd:
  root:
    name:
      base: ""
      full: ""
    path:
      base: ""
      full: ""
    package:
      reference: ""
      name:
        space: ""
        snake: ""
        kebob: ""
        exported: ""
        unexported: ""
      path:
        base: ""
        full: ""
    imports:
      standard: []
      app: []
      vendor: []
    init_function:
      name: ""
      imports:
        standard: []
        app: []
        vendor: []
      arguments: []
      return_values: []
      receiver:
        name: ""
        type: ""
      body: ""
    consts: []
    vars: []
    interfaces: []
    structs: []
    functions: []
  server:
    name:
      base: ""
      full: ""
    path:
      base: ""
      full: ""
    package:
      reference: ""
      name:
        space: ""
        snake: ""
        kebob: ""
        exported: ""
        unexported: ""
      path:
        base: ""
        full: ""
    imports:
      standard: []
      app: []
      vendor: []
    init_function:
      name: ""
      imports:
        standard: []
        app: []
        vendor: []
      arguments: []
      return_values: []
      receiver:
        name: ""
        type: ""
      body: ""
    consts: []
    vars: []
    interfaces: []
    structs: []
    functions: []
  migrate:
    name:
      base: ""
      full: ""
    path:
      base: ""
      full: ""
    package:
      reference: ""
      name:
        space: ""
        snake: ""
        kebob: ""
        exported: ""
        unexported: ""
      path:
        base: ""
        full: ""
    imports:
      standard: []
      app: []
      vendor: []
    init_function:
      name: ""
      imports:
        standard: []
        app: []
        vendor: []
      arguments: []
      return_values: []
      receiver:
        name: ""
        type: ""
      body: ""
    consts: []
    vars: []
    interfaces: []
    structs: []
    functions: []
  seed:
    name:
      base: ""
      full: ""
    path:
      base: ""
      full: ""
    package:
      reference: ""
      name:
        space: ""
        snake: ""
        kebob: ""
        exported: ""
        unexported: ""
      path:
        base: ""
        full: ""
    imports:
      standard: []
      app: []
      vendor: []
    init_function:
      name: ""
      imports:
        standard: []
        app: []
        vendor: []
      arguments: []
      return_values: []
      receiver:
        name: ""
        type: ""
      body: ""
    consts: []
    vars: []
    interfaces: []
    structs: []
    functions: []
  custom: []
http:
  routes:
    name:
      base: ""
      full: ""
    path:
      base: ""
      full: ""
    package:
      reference: ""
      name:
        space: ""
        snake: ""
        kebob: ""
        exported: ""
        unexported: ""
      path:
        base: ""
        full: ""
    imports:
      standard: []
      app: []
      vendor: []
    init_function:
      name: ""
      imports:
        standard: []
        app: []
        vendor: []
      arguments: []
      return_values: []
      receiver:
        name: ""
        type: ""
      body: ""
    consts: []
    vars: []
    interfaces: []
    structs: []
    functions: []
`
