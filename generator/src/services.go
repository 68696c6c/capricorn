package src

const serviceTemplate = `
package {{ .Package }}

import (
	{{- range $key, $value := .Imports }}
	"{{ $value }}"
	{{- end }}
)

type {{ .Name.Exported }} struct {
	repo Repo
}

type Options struct {
	Repo Repo
}

func {{ .Constructor }}(o Options) {{ .Name.Exported }} {
	return {{ .Name.Exported }}{
		repo: o.Repo,
	}
}

{{- range $key, $value := .MethodTemplates }}
{{ $value }}
{{- end }}

`

const serviceMethodTemplate = `
func (r {{.Receiver}}) {{.Signature}} {
	return nil
}`
