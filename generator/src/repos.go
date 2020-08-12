package src

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const repoTemplate = `
package {{ .Package }}

import (
	{{- range $key, $value := .Imports }}
	"{{ $value }}"
	{{- end }}

	{{- range $key, $value := .VendorImports }}
	"{{ $value }}"
	{{- end }}
)

type {{ .InterfaceName }} interface {
	{{- range $key, $value := .InterfaceTemplates }}
	{{ $value }}
	{{- end }}
}

type {{ .Name.Exported }} struct {
	db *gorm.DB
}

func {{ .Constructor }}(d *gorm.DB) {{ .InterfaceName }} {
	return {{ .Name.Exported }}{
		db: d,
	}
}


{{- range $key, $value := .MethodTemplates }}
{{ $value }}
{{- end }}

`

const repoSaveTemplate = `
func (r {{ .Receiver }}) {{ .Signature }} {
	if m.Model.ID.Valid() {
		errs = r.db.Save(m).GetErrors()
	} else {
		errs = r.db.Create(m).GetErrors()
	}
	return
}`

const repoGetByIDTemplate = `
func (r {{ .Receiver }}) {{ .Signature }} {
	m := {{ .Resource.Single.Exported }}{}
	errs := r.db.First(&m, "id = ?", id).GetErrors()
	return m, errs
}`

const repoListTemplate = `
func (r {{ .Receiver }}) {{ .Signature }} {
	base := r.db.Model(&{{ .Resource.Single.Exported }}{})

	qr, err := q.ApplyToGorm(base)
	if err != nil {
		return m, []error{err}
	}

	errs = qr.Find(&m).GetErrors()
	return
}`

const repoSetQueryTotalTemplate = `
func (r {{ .Receiver }}) {{ .Signature }} {
	base := r.db.Model(&{{ .Resource.Single.Exported }}{})

	qr, err := q.ApplyToGormCount(base)
	if err != nil {
		return []error{err}
	}

	var count uint
	errs := qr.Count(&count).GetErrors()
	if len(errs) > 0 {
		return errs
	}

	q.Pagination.Total = count

	return []error{}
}`

const repoDeleteTemplate = `
func (r {{ .Receiver }}) {{ .Signature }} {
	return r.db.Delete(m).GetErrors()
}`

func CreateRepos(spec models.Project, logger *logrus.Logger) error {
	logPrefix := "CreateRepos | "

	err := utils.CreateDir(spec.Paths.Repos)
	if err != nil {
		return errors.Wrapf(err, "failed to create repos directory '%s'", spec.Paths.Repos)
	}

	for _, r := range spec.Repos {
		logger.Infof(logPrefix, fmt.Sprintf("creating repo %s", r.Filename))

		for _, m := range r.Methods {
			logger.Infof(logPrefix, fmt.Sprintf("creating repo method %s", m.Name))

			switch m.Name {
			case "Save":
				mt, err := utils.ParseTemplateToString("repo_save", repoSaveTemplate, m)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo method 'Save'")
				}
				r.MethodTemplates = append(r.MethodTemplates, mt)
				r.InterfaceTemplates = append(r.InterfaceTemplates, m.Signature)

			case "GetByID":
				mt, err := utils.ParseTemplateToString("repo_get", repoGetByIDTemplate, m)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo method 'GetByID'")
				}
				r.MethodTemplates = append(r.MethodTemplates, mt)
				r.InterfaceTemplates = append(r.InterfaceTemplates, m.Signature)

			case "List":
				mt, err := utils.ParseTemplateToString("repo_list", repoListTemplate, m)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo method 'List'")
				}
				r.MethodTemplates = append(r.MethodTemplates, mt)
				r.InterfaceTemplates = append(r.InterfaceTemplates, m.Signature)

			case "SetQueryTotal":
				mt, err := utils.ParseTemplateToString("query_total", repoSetQueryTotalTemplate, m)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo method 'SetQueryTotal'")
				}
				r.MethodTemplates = append(r.MethodTemplates, mt)
				r.InterfaceTemplates = append(r.InterfaceTemplates, m.Signature)

			case "Delete":
				mt, err := utils.ParseTemplateToString("repo_delete", repoDeleteTemplate, m)
				if err != nil {
					return errors.Wrap(err, "failed to generate repo method 'Delete'")
				}
				r.MethodTemplates = append(r.MethodTemplates, mt)
				r.InterfaceTemplates = append(r.InterfaceTemplates, m.Signature)
			}
		}

		err = utils.GenerateFile(spec.Paths.Repos, r.Filename, repoTemplate, r)
		if err != nil {
			return errors.Wrap(err, "failed to generate repo")
		}
	}

	return nil
}
