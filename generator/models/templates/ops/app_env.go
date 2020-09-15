package ops

import (
	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/utils"
)

const AppEnvTemplate = `
DB_HOST={{ .MainDatabase.Host }}
DB_PORT={{ .MainDatabase.Port }}
DB_USERNAME={{ .MainDatabase.Username }}
DB_PASSWORD={{ .MainDatabase.Password }}
DB_DATABASE={{ .MainDatabase.Name }}
DB_DEBUG={{ .MainDatabase.Debug }}
`

type AppEnv struct {
	Name data.FileData `yaml:"name,omitempty"`
	Path data.PathData `yaml:"path,omitempty"`

	Data Ops `yaml:"data,omitempty"`
}

// This is only used for testing.
func (m AppEnv) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_app_env", AppEnvTemplate, m.Data)
	if err != nil {
		panic(err)
	}
	return result
}

func (m AppEnv) Generate() error {
	err := utils.GenerateFile(m.Path.Base, m.Name.Full, AppEnvTemplate, m.Data)
	if err != nil {
		return err
	}
	return nil
}
