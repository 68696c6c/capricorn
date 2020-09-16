package golang

import (
	"fmt"
	"strings"

	"github.com/68696c6c/capricorn/generator/utils"
)

type Imports struct {
	Standard []string `yaml:"standard,omitempty"`
	App      []string `yaml:"app,omitempty"`
	Vendor   []string `yaml:"vendor,omitempty"`
}

func (m Imports) HasImports() bool {
	return len(m.Standard) > 0 || len(m.App) > 0 || len(m.Vendor) > 0
}

func (m Imports) MustParse() string {
	if !m.HasImports() {
		return ""
	}

	appendSection := func(heap []string, section []string) []string {
		if len(section) > 0 {
			var sRes []string
			for _, i := range section {
				sRes = append(sRes, fmt.Sprintf(`	"%s"`, i))
			}
			heap = append(heap, strings.Join(sRes, "\n"))
		}
		return heap
	}

	var sectionImports []string
	sectionImports = appendSection(sectionImports, m.Standard)
	sectionImports = appendSection(sectionImports, m.App)
	sectionImports = appendSection(sectionImports, m.Vendor)

	result := []string{"import ("}

	// Separate each section with an additional line break.
	result = append(result, strings.Join(sectionImports, "\n\n"))

	result = append(result, ")")

	return strings.Join(result, "\n")
}

func MergeImports(target, additional Imports) Imports {
	target.Standard = append(target.Standard, additional.Standard...)
	target.App = append(target.App, additional.App...)
	target.Vendor = append(target.Vendor, additional.Vendor...)
	return Imports{
		Standard: utils.RemoveDuplicateStrings(target.Standard),
		App:      utils.RemoveDuplicateStrings(target.App),
		Vendor:   utils.RemoveDuplicateStrings(target.Vendor),
	}
}
