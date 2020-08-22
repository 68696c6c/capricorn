package module

import "github.com/68696c6c/capricorn/generator/models/spec"

func GetFixtureModule() Module {
	f := spec.GetFixtureSpec()
	return NewModuleFromSpec(f)
}

const FixtureModuleYAML = ``
