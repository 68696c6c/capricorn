package project

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
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

func CreateModels(spec Spec, logger *logrus.Logger) error {
	logPrefix := "CreateModels | "

	err := createDir(spec.Paths.Models)
	if err != nil {
		return errors.Wrapf(err, "failed to create models directory '%s'", spec.Paths.Models)
	}

	// Create models.
	for _, m := range spec.Models {
		m.StructName = snakeToCamel(m.Name)

		// Build relations.
		if len(m.BelongsTo) > 0 {
			logger.Debug(logPrefix, "model belongs to")

			for _, r := range m.BelongsTo {
				logger.Debug(logPrefix, "relation: ", r)

				f := &field{
					Name:      fmt.Sprintf("%s_id", r),
					FieldName: fmt.Sprintf("%sID", snakeToCamel(r)),
					Type:      "goat.ID",
				}
				m.Fields = append([]*field{f}, m.Fields...)
			}
		}
		if len(m.HasMany) > 0 {
			logger.Debug(logPrefix, "model has many")

			for _, r := range m.HasMany {
				logger.Debug(logPrefix, "relation: ", r)

				t := inflection.Singular(r)
				f := &field{
					Name:      r,
					FieldName: snakeToCamel(r),
					Type:      fmt.Sprintf("[]*%s", snakeToCamel(t)),
				}
				m.Fields = append(m.Fields, f)
			}
		}

		// Set field names and annotations.
		for _, f := range m.Fields {
			if f.FieldName == "" {
				f.FieldName = snakeToCamel(f.Name)
			}
			var extra string
			if f.Required {
				extra = ` binding:"required"`
			}
			f.Tag = fmt.Sprintf(`json:"%s"%s`, f.Name, extra)
		}

		err = generateFile(spec.Paths.Models, m.Name, modelTemplate, *m)
		if err != nil {
			return errors.Wrap(err, "failed to generate model")
		}
	}

	return nil
}
