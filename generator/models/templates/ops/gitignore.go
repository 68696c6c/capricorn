package ops

import (
	utils2 "github.com/68696c6c/capricorn/generator/models/utils"
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
	Name utils2.FileData `yaml:"name"`
	Path utils2.PathData `yaml:"path"`

	Data Ops `yaml:"data"`
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