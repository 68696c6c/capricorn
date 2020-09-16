package data

import "github.com/jinzhu/inflection"

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
