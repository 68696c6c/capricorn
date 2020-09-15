package ops

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/utils"
)

const GitignoreTemplate = `
.DS_Store
.idea
vendor
.app.env
}
`

type Gitignore struct {
	Name data.FileData `yaml:"name,omitempty"`
	Path data.PathData `yaml:"path,omitempty"`

	Data Ops `yaml:"data,omitempty"`
}

// This is only used for testing.
func (m Gitignore) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_gitignore", GitignoreTemplate, m.Data)
	if err != nil {
		panic(err)
	}
	return result
}

func (m Gitignore) Generate() error {
	err := utils.GenerateFile(m.Path.Base, m.Name.Full, GitignoreTemplate, m.Data)
	if err != nil {
		return err
	}
	return nil
}
