package src

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const modelTemplate = `
package models

import "github.com/68696c6c/goat"

{{ $tick := "` + "`" + `" }}
// swagger:model {{ .StructName }}
type {{ .StructName }} struct {
	goat.Model
	{{- range $key, $value := .Fields }}
	{{ $value.FieldName }} {{ $value.Type }} {{ $tick }}{{ $value.Tag }}{{ $tick }}
	{{- end }}
}

`

func CreateModels(spec utils.Spec, logger *logrus.Logger) error {
	logPrefix := "CreateModels | "

	err := utils.CreateDir(spec.Paths.Models)
	if err != nil {
		return errors.Wrapf(err, "failed to create models directory '%s'", spec.Paths.Models)
	}

	// Create models.
	for _, m := range spec.Models {
		m.StructName = utils.SeparatedToExported(m.Name)

		// Build relations.
		if len(m.BelongsTo) > 0 {
			logger.Debug(logPrefix, "model belongs to")

			for _, r := range m.BelongsTo {
				logger.Debug(logPrefix, "relation: ", r)

				f := &utils.Field{
					Name:      fmt.Sprintf("%s_id", r),
					FieldName: fmt.Sprintf("%sID", utils.SeparatedToExported(r)),
					Type:      "goat.ID",
				}
				m.Fields = append([]*utils.Field{f}, m.Fields...)
			}
		}
		if len(m.HasMany) > 0 {
			logger.Debug(logPrefix, "model has many")

			for _, r := range m.HasMany {
				logger.Debug(logPrefix, "relation: ", r)

				t := inflection.Singular(r)
				f := &utils.Field{
					Name:      r,
					FieldName: utils.SeparatedToExported(r),
					Type:      fmt.Sprintf("[]*%s", utils.SeparatedToExported(t)),
				}
				m.Fields = append(m.Fields, f)
			}
		}

		// Set field names and annotations.
		for _, f := range m.Fields {
			if f.FieldName == "" {
				f.FieldName = utils.SeparatedToExported(f.Name)
			}
			var extra string
			if f.Required {
				extra = ` binding:"required"`
			}
			f.Tag = fmt.Sprintf(`json:"%s"%s`, f.Name, extra)
		}

		err = utils.GenerateGoFile(spec.Paths.Models, m.Name, modelTemplate, *m)
		if err != nil {
			return errors.Wrap(err, "failed to generate model")
		}
	}

	return nil
}
