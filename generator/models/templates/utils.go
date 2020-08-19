package templates

type Template interface {
	Generate() error
}

// MustParse is for rendering a template inside of another template.  Since it is called inside of a template, there is
// no way to handle an error, leaving no option but to panic, hence "must" in the name.
type SubTemplate interface {
	MustParse() string
}

type FileData struct {
	// e.g. path: src/app/domain/example.go
	// e.g. file: example.go
	Full string `yaml:"full"`

	// e.g. path: src/app/domain/
	// e.g. file: example
	Base string `yaml:"base"`
}
