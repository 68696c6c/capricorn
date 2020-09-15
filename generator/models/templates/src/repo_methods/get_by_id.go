package repo_methods

import (
	"fmt"

	"github.com/68696c6c/capricorn/generator/models/data"
	"github.com/68696c6c/capricorn/generator/models/templates/golang"
	"github.com/68696c6c/capricorn/generator/utils"
)

var getByIdBodyTemplate = `
	m := {{ .Single.Exported }}{
		Model: goat.Model{
			ID: id,
		},
	}
	errs := {{ .GetDbReference }}.First(&m).GetErrors()
	if len(errs) > 0 {
		return m, goat.ErrorsToError(errs)
	}
	return m, nil
`

type GetByID struct {
	dbFieldName string
	receiver    golang.Value
	Single      data.Name
}

func NewGetByID(meta MethodMeta) GetByID {
	return GetByID{
		dbFieldName: meta.DBFieldName,
		receiver:    meta.Receiver,
		Single:      meta.Resource.Inflection.Single,
	}
}

func (m GetByID) GetDbReference() string {
	return fmt.Sprintf("%s.%s", m.receiver.Name, m.dbFieldName)
}

func (m GetByID) MustGetFunction() golang.Function {
	return golang.Function{
		Name:         m.GetName(),
		Imports:      m.GetImports(),
		Receiver:     m.GetReceiver(),
		Arguments:    m.GetArgs(),
		ReturnValues: m.GetReturns(),
		Body:         m.MustParse(),
	}
}

func (m GetByID) GetName() string {
	return "GetByID"
}

func (m GetByID) GetImports() golang.Imports {
	return golang.Imports{
		Standard: nil,
		App:      nil,
		Vendor:   []string{data.ImportGoat},
	}
}

func (m GetByID) GetReceiver() golang.Value {
	return m.receiver
}

func (m GetByID) GetArgs() []golang.Value {
	return []golang.Value{
		{
			Name: "id",
			Type: "goat.ID",
		},
	}
}

func (m GetByID) GetReturns() []golang.Value {
	return []golang.Value{
		{
			Type: m.Single.Exported,
		},
		{
			Type: "error",
		},
	}
}

func (m GetByID) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_repo_get_by_id", getByIdBodyTemplate, m)
	if err != nil {
		panic(err)
	}
	return result
}
