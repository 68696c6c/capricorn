package utils

import "github.com/jinzhu/inflection"

type Name struct {
	Space      string `yaml:"space,omitempty"`
	Snake      string `yaml:"snake,omitempty"`
	Kebob      string `yaml:"kebob,omitempty"`
	Exported   string `yaml:"exported,omitempty"`
	Unexported string `yaml:"unexported,omitempty"`
}

func MakeName(base string) Name {
	return Name{
		Space:      SeparatedToSpace(base),
		Snake:      SeparatedToSnake(base),
		Kebob:      SeparatedToKebob(base),
		Exported:   SeparatedToExported(base),
		Unexported: SeparatedToUnexported(base),
	}
}

type Inflection struct {
	Single Name `yaml:"single,omitempty"`
	Plural Name `yaml:"plural,omitempty"`
}

// input can be either plural or singular
func MakeInflection(input string) Inflection {
	single := inflection.Singular(input)
	plural := inflection.Plural(input)
	return Inflection{
		Single: MakeName(single),
		Plural: MakeName(plural),
	}
}
