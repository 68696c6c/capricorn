package golang

import (
	"fmt"
	"strings"
)

type SubTemplate interface {
	MustParse() string
}

func getJoinedValueString(values []Value) string {
	var builtValues []string
	for _, v := range values {
		builtValues = append(builtValues, fmt.Sprintf("%s %s", v.Name, v.Type))
	}
	joinedValues := strings.Join(builtValues, ", ")
	return strings.TrimSpace(joinedValues)
}
