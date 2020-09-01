package src

import (
	"testing"

	"github.com/68696c6c/capricorn/generator/models/module"
	"github.com/68696c6c/capricorn/generator/models/spec"

	"github.com/stretchr/testify/assert"
)

func TestSRC_NewSRCDDD(t *testing.T) {
	f := spec.GetFixtureSpec()
	m := module.NewModuleFromSpec(f)

	// _ = NewSRCDDD(m, "/base/path/test-example")
	result := NewSRCDDD(m, "/base/path/test-example")
	resultYAML := result.String()
	println(resultYAML)

	assert.Equal(t, FixtureSRCYAML, resultYAML)
	// assert.True(t, false)
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
      structs: []
      functions:
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
          type: ""
        body: "\n\tq := query.NewQueryBuilder(c)\n\n\tresult, errs := c.repo.List(q)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c, errs, \"failed to list organizations\",
          goat.RespondServerError)\n\t\treturn\n\t}\n\n\terrs = c.repo.SetQueryTotal(q)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c, errs, \"failed to count organizations\",
          goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondData(c, ListResponse{result,
          q.Pagination})\n"
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
          type: ""
        body: "\n\ti := c.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(c, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\tm,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(c, \"organization does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(c, errs, \"failed to get organization\",
          goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\tgoat.RespondData(c,
          ListResponse{m})\n"
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
          type: ""
        body: "\n\treq, ok := goat.GetRequest(c).(*CreateRequest)\n\tif !ok {\n\t\tc.errors.HandleMessage(c,
          \"failed to get request\", goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\tm
          := req.Organization\n\terrs := c.repo.Save(&m)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c,
          errs, \"failed to save organization\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondCreated(c,
          ListResponse{m})\n"
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
          type: ""
        body: "\n\ti := c.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(c, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t_,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(c, \"organization does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(c, errs, \"failed to get organization\",
          goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\treq, ok := goat.GetRequest(c).(*UpdateRequest)\n\tif
          !ok {\n\t\tc.errors.HandleMessage(c, \"failed to get request\", goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\terrs
          = c.repo.Save(&req.Organization)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c,
          errs, \"failed to save organization\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondCreated(c,
          ListResponse{req.Organization})\n"
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
          type: ""
        body: "\n\ti := c.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(c, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\tm,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(c, \"organization does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(c, errs, \"failed to get organization\",
          goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\t// @TODO generate
          model factories.\n\t// @TODO generate model validators.\n\terrs = c.repo.Delete(&m)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c, errs, \"failed to delete
          organization\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondValid(c)\n"
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
      structs: []
      functions:
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
          type: ""
        body: "\n\tq := query.NewQueryBuilder(c)\n\n\tresult, errs := c.repo.List(q)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c, errs, \"failed to list users\",
          goat.RespondServerError)\n\t\treturn\n\t}\n\n\terrs = c.repo.SetQueryTotal(q)\n\tif
          len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c, errs, \"failed to count users\",
          goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondData(c, ListResponse{result,
          q.Pagination})\n"
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
          type: ""
        body: "\n\ti := c.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(c, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\tm,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(c, \"user does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(c, errs, \"failed to get user\", goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\tgoat.RespondData(c,
          ListResponse{m})\n"
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
          type: ""
        body: "\n\treq, ok := goat.GetRequest(c).(*CreateRequest)\n\tif !ok {\n\t\tc.errors.HandleMessage(c,
          \"failed to get request\", goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\tm
          := req.User\n\terrs := c.repo.Save(&m)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c,
          errs, \"failed to save user\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondCreated(c,
          ListResponse{m})\n"
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
          type: ""
        body: "\n\ti := c.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(c, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t_,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(c, \"user does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(c, errs, \"failed to get user\", goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\treq,
          ok := goat.GetRequest(c).(*UpdateRequest)\n\tif !ok {\n\t\tc.errors.HandleMessage(c,
          \"failed to get request\", goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\terrs
          = c.repo.Save(&req.User)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c,
          errs, \"failed to save user\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondCreated(c,
          ListResponse{req.User})\n"
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
          type: ""
        body: "\n\ti := c.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(c, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\tm,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(c, \"user does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(c, errs, \"failed to get user\", goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\terrs
          = c.repo.Delete(&m)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c,
          errs, \"failed to delete user\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondValid(c)\n"
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
      structs: []
      functions:
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
          type: ""
        body: "\n\treq, ok := goat.GetRequest(c).(*CreateRequest)\n\tif !ok {\n\t\tc.errors.HandleMessage(c,
          \"failed to get request\", goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\tm
          := req.Token\n\terrs := c.repo.Save(&m)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c,
          errs, \"failed to save token\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondCreated(c,
          ListResponse{m})\n"
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
          type: ""
        body: "\n\ti := c.Param(\"id\")\n\tid, err := goat.ParseID(i)\n\tif err !=
          nil {\n\t\tc.errors.HandleErrorM(c, err, \"failed to parse id: \"+i, goat.RespondBadRequestError)\n\t\treturn\n\t}\n\n\tm,
          errs := c.repo.GetByID(id)\n\tif len(errs) > 0 {\n\t\tif goat.RecordNotFound(errs)
          {\n\t\t\tc.errors.HandleMessage(c, \"token does not exist\", goat.RespondNotFoundError)\n\t\t\treturn\n\t\t}
          else {\n\t\t\tc.errors.HandleErrorsM(c, errs, \"failed to get token\", goat.RespondServerError)\n\t\t\treturn\n\t\t}\n\t}\n\n\t//
          @TODO generate model factories.\n\t// @TODO generate model validators.\n\terrs
          = c.repo.Delete(&m)\n\tif len(errs) > 0 {\n\t\tc.errors.HandleErrorsM(c,
          errs, \"failed to delete token\", goat.RespondServerError)\n\t\treturn\n\t}\n\n\tgoat.RespondValid(c)\n"
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
