package golang

import (
	"fmt"
	"strings"
)

type Imports struct {
	// Standard []string `yaml:"standard,omitempty"`
	// App      []string `yaml:"app,omitempty"`
	// Vendor   []string `yaml:"vendor,omitempty"`

	Standard []Package
	Vendor   []Package
	App      []Package
}

func (m Imports) HasImports() bool {
	return len(m.Standard) > 0 || len(m.App) > 0 || len(m.Vendor) > 0
}

func (m Imports) MustString() string {
	if !m.HasImports() {
		return ""
	}

	appendSection := func(heap []string, section []Package) []string {
		if len(section) > 0 {
			var sRes []string
			for _, p := range section {
				i := p.GetImport()
				if strings.Contains(i, `"`) {
					sRes = append(sRes, fmt.Sprintf(`	%s`, i))
				} else {
					sRes = append(sRes, fmt.Sprintf(`	"%s"`, i))
				}
			}
			heap = append(heap, strings.Join(sRes, "\n"))
		}
		return heap
	}

	var sectionImports []string
	sectionImports = appendSection(sectionImports, m.Standard)
	sectionImports = appendSection(sectionImports, m.Vendor)
	sectionImports = appendSection(sectionImports, m.App)

	result := []string{"import ("}

	// Separate each section with an additional line break.
	result = append(result, strings.Join(sectionImports, "\n\n"))

	result = append(result, ")")

	return strings.Join(result, "\n")
}

func MergeImports(target, additional Imports) Imports {
	target.Standard = append(target.Standard, additional.Standard...)
	target.Vendor = append(target.Vendor, additional.Vendor...)
	target.App = append(target.App, additional.App...)
	return Imports{
		Standard: removeDuplicateImports(target.Standard),
		Vendor:   removeDuplicateImports(target.Vendor),
		App:      removeDuplicateImports(target.App),
	}
}

func removeDuplicateImports(items []Package) []Package {
	keys := make(map[string]bool)
	var result []Package
	for _, i := range items {
		imp := i.GetImport()
		if _, ok := keys[imp]; !ok {
			keys[imp] = true
			result = append(result, i)
		}
	}
	return result
}

// func (m Imports) MustString() string {
// 	if !m.HasImports() {
// 		return ""
// 	}
//
// 	appendSection := func(heap []string, section []string) []string {
// 		if len(section) > 0 {
// 			var sRes []string
// 			for _, i := range section {
// 				if strings.Contains(i, `"`) {
// 					sRes = append(sRes, fmt.Sprintf(`	%s`, i))
// 				} else {
// 					sRes = append(sRes, fmt.Sprintf(`	"%s"`, i))
// 				}
// 			}
// 			heap = append(heap, strings.Join(sRes, "\n"))
// 		}
// 		return heap
// 	}
//
// 	var sectionImports []string
// 	sectionImports = appendSection(sectionImports, m.Standard)
// 	sectionImports = appendSection(sectionImports, m.App)
// 	sectionImports = appendSection(sectionImports, m.Vendor)
//
// 	result := []string{"import ("}
//
// 	// Separate each section with an additional line break.
// 	result = append(result, strings.Join(sectionImports, "\n\n"))
//
// 	result = append(result, ")")
//
// 	return strings.Join(result, "\n")
// }

// func MergeImports(target, additional Imports) Imports {
// 	target.Standard = append(target.Standard, additional.Standard...)
// 	target.App = append(target.App, additional.App...)
// 	target.Vendor = append(target.Vendor, additional.Vendor...)
// 	return Imports{
// 		Standard: utils.RemoveDuplicateStrings(target.Standard),
// 		App:      utils.RemoveDuplicateStrings(target.App),
// 		Vendor:   utils.RemoveDuplicateStrings(target.Vendor),
// 	}
// }
