package src

import (
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	repoMethodCreate = "create"
	repoMethodUpdate = "update"
	repoMethodSave   = "save"
	repoMethodGet    = "get"
	repoMethodList   = "list"
	repoMethodDelete = "delete"
)

const repoTemplate = `
package repos

import (
	"time"

	"{{.ModelsImportPath}}"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/jinzhu/gorm"
)

type {{.InterfaceName}} interface {
	{{- range $key, $value := .InterfaceTemplateMethods }}
	{{ $value }}
	{{- end }}
}

type {{.StructName}} struct {
	db *gorm.DB
}

func {{.ConstructorName}}(d *gorm.DB) {{.StructName}} {
	return {{.StructName}}{
		db: d,
	}
}


{{- range $key, $value := .MethodTemplates }}
{{ $value }}
{{- end }}

`

const repoInterfaceSaveTemplate = `Save(model *models.{{.ModelStructName}}) (errs []error)`
const repoSaveTemplate = `
func (r {{.StructName}}) Save(model *models.{{.ModelStructName}}) (errs []error) {
	if model.Model.ID.Valid() {
		errs = r.db.Save(model).GetErrors()
	} else {
		errs = r.db.Create(model).GetErrors()
	}
	return
}
`

const repoInterfaceGetTemplate = `GetByID(id goat.ID) (*models.{{.ModelStructName}}, []error)`
const repoGetByIDTemplate = `
func (r {{.StructName}}) GetByID(id goat.ID) (*models.{{.ModelStructName}}, []error) {
	m := &models.{{.ModelStructName}}{}
	errs := r.db.First(m, "id = ?", id).GetErrors()
	return m, errs
}
`

const repoInterfaceListTemplate = `List(q *query.Query) (m []*models.{{.ModelStructName}}, errs []error)
	SetQueryTotal(q *query.Query) (errs []error)`
const repoListTemplate = `
func (r {{.StructName}}) List(q *query.Query) (m []*models.{{.ModelStructName}}, errs []error) {
	base := r.db.Model(&models.{{.ModelStructName}}{})

	qr, err := q.ApplyToGorm(base)
	if err != nil {
		return m, []error{err}
	}

	errs = qr.Find(&m).GetErrors()
	return
}

func (r {{.StructName}}) SetQueryTotal(q *query.Query) (errs []error) {
	base := r.db.Model(&models.{{.ModelStructName}}{})

	qr, err := q.ApplyToGormCount(base)
	if err != nil {
		return []error{err}
	}

	var count uint
	errs = qr.Count(&count).GetErrors()
	if len(errs) > 0 {
		return errs
	}

	q.Pagination.Total = count

	return []error{}
}`

const repoInterfaceDeleteTemplate = `Delete(model *models.{{.ModelStructName}}) []error`
const repoDeleteTemplate = `
func (r {{.StructName}}) Delete(model *models.{{.ModelStructName}}) []error {
	n := time.Now()
	model.DeletedAt = &n
	return r.db.Save(model).GetErrors()
}
`

func CreateRepos(spec *utils.Spec, logger *logrus.Logger) error {
	logPrefix := "CreateRepos | "

	err := utils.CreateDir(spec.Paths.Repos)
	if err != nil {
		return errors.Wrapf(err, "failed to create repos directory '%s'", spec.Paths.Repos)
	}

	for _, r := range spec.Repos {
		logger.Debug(logPrefix, "repo ", r)
		logger.Debug(logPrefix, "repo model ", r.Model)

		model := utils.SnakeToExportedName(r.Model)
		plural := inflection.Plural(model)
		r.Name = inflection.Plural(r.Model) + "_repo"
		r.InterfaceName = plural + "Repo"
		r.StructName = plural + "RepoGORM"
		r.ModelsImportPath = spec.Imports.Models
		r.ModelStructName = model
		r.ConstructorName = "New" + r.InterfaceName

		// If no methods were specified, default to all.
		if len(r.Methods) == 0 {
			r.Methods = []string{
				repoMethodSave,
				repoMethodGet,
				repoMethodList,
				repoMethodDelete,
			}
		}

		for _, m := range r.Methods {
			logger.Debug(logPrefix, "method ", m)

			switch m {
			case repoMethodCreate:
				fallthrough
			case repoMethodUpdate:
				fallthrough
			case repoMethodSave:
				method, err := utils.ParseTemplateToString("repo_save", repoSaveTemplate, r)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo method 'save'")
				}
				r.MethodTemplates = append(r.MethodTemplates, method)
				intMethod, err := utils.ParseTemplateToString("repo_interface_save", repoInterfaceSaveTemplate, r)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo interface method 'save'")
				}
				r.InterfaceTemplateMethods = append(r.InterfaceTemplateMethods, intMethod)

			case repoMethodGet:
				method, err := utils.ParseTemplateToString("repo_get", repoGetByIDTemplate, r)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo method 'get'")
				}
				r.MethodTemplates = append(r.MethodTemplates, method)
				intMethod, err := utils.ParseTemplateToString("repo_interface_get", repoInterfaceGetTemplate, r)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo interface method 'get'")
				}
				r.InterfaceTemplateMethods = append(r.InterfaceTemplateMethods, intMethod)

			case repoMethodList:
				method, err := utils.ParseTemplateToString("repo_list", repoListTemplate, r)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo method 'list'")
				}
				r.MethodTemplates = append(r.MethodTemplates, method)
				intMethod, err := utils.ParseTemplateToString("repo_interface_list", repoInterfaceListTemplate, r)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo interface method 'list'")
				}
				r.InterfaceTemplateMethods = append(r.InterfaceTemplateMethods, intMethod)

			case repoMethodDelete:
				method, err := utils.ParseTemplateToString("repo_delete", repoDeleteTemplate, r)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo method 'delete'")
				}
				r.MethodTemplates = append(r.MethodTemplates, method)
				intMethod, err := utils.ParseTemplateToString("repo_interface_delete", repoInterfaceDeleteTemplate, r)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo interface method 'delete'")
				}
				r.InterfaceTemplateMethods = append(r.InterfaceTemplateMethods, intMethod)
			}
		}

		err = utils.GenerateGoFile(spec.Paths.Repos, r.Name, repoTemplate, *r)
		if err != nil {
			return errors.Wrap(err, "failed to generate repo")
		}
	}

	return nil
}
