package golang

import (
	"fmt"

	"github.com/68696c6c/girraph"

	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

type Package interface {
	SetBasePath(string) Package
	GetBasePath() string
	SetName(string) Package
	GetName() string
	GetFullPath() string
	SetImport(string) Package
	GetImport() string
	SetReference(string) Package
	GetReference() string
	AddFile(*File) Package
	SetFiles([]*File) Package
	GetFiles() []*File
}

type pkg struct {
	BasePath  string
	Name      string
	Import    string
	Reference string
	Files     []*File
}

func (p *pkg) SetBasePath(basePath string) Package {
	p.BasePath = basePath
	return p
}

func (p *pkg) GetBasePath() string {
	return p.BasePath
}

func (p *pkg) SetName(name string) Package {
	p.Name = name
	return p
}

func (p *pkg) GetName() string {
	return p.Name
}

func (p *pkg) GetFullPath() string {
	return fmt.Sprintf("%s/%s", p.GetBasePath(), p.GetName())
}

func (p *pkg) SetImport(path string) Package {
	p.Import = path
	return p
}

func (p *pkg) GetImport() string {
	return p.Import
}

func (p *pkg) GetReference() string {
	return p.Reference
}

func (p *pkg) SetReference(reference string) Package {
	p.Reference = reference
	return p
}

func (p *pkg) AddFile(file *File) Package {
	file.SetPackage(p.Reference)
	p.Files = append(p.Files, file)
	return p
}

func (p *pkg) SetFiles(files []*File) Package {
	p.Files = files
	for _, file := range p.Files {
		file.SetPackage(p.Reference)
	}
	return p
}

func (p *pkg) GetFiles() []*File {
	return p.Files
}

func MakePackageNode(name string) girraph.Tree[Package] {
	return girraph.MakeTree[Package]().SetMeta(&pkg{
		Name:      name,
		Reference: name,
		Files:     []*File{},
	})
}

func MakePackage(importPath string) Package {
	return &pkg{
		Import: importPath,
		Files:  []*File{},
	}
}

func SetPaths(basePath, baseImport string, node girraph.Tree[Package]) {
	meta := node.GetMeta()
	metaName := baseImport
	if meta.GetName() != "" {
		metaName = utils.JoinPath(baseImport, meta.GetName())
	}
	meta.SetImport(metaName).SetBasePath(basePath)
	imp := meta.GetImport()
	path := meta.GetFullPath()
	for _, child := range node.GetChildren() {
		SetPaths(path, imp, child)
	}
}

func mustGenerate(p Package) {
	dirPath := fmt.Sprintf("%s/%s", p.GetBasePath(), p.GetName())
	utils.PanicError(utils.CreateDir(dirPath))
	for _, file := range p.GetFiles() {
		err := utils.WriteFile(dirPath, file.GetFullName(), file.GetContents())
		if err != nil {
			panic(err)
		}
	}
}

func MustGenerate(node girraph.Tree[Package]) {
	p := node.GetMeta()
	mustGenerate(p)
	for _, child := range node.GetChildren() {
		MustGenerate(child)
	}
}
