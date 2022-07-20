package filesystem

import (
	"fmt"

	"github.com/68696c6c/girraph"

	"github.com/68696c6c/capricorn_rnd/generator/utils"
)

type File struct {
	Name      string
	Extension string
	Contents  string
}

func (f *File) SetName(name string) *File {
	f.Name = name
	return f
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) SetExtension(ext string) *File {
	f.Extension = ext
	return f
}

func (f *File) GetExtension() string {
	return f.Extension
}

func (f *File) GetFullName() string {
	if f.Extension == "" {
		return f.Name
	}
	return fmt.Sprintf("%s.%s", f.Name, f.Extension)
}

func (f *File) SetContents(contents string) *File {
	f.Contents = contents
	return f
}

func (f *File) GetContents() string {
	return f.Contents
}

func MakeFile(name, extension string) *File {
	return &File{
		Name:      name,
		Extension: extension,
	}
}

type Directory interface {
	SetBasePath(string) Directory
	GetBasePath() string
	SetName(string) Directory
	GetName() string
	GetFullPath() string
	SetFiles([]*File) Directory
	GetFiles() []*File
}

type directory struct {
	BasePath string
	Name     string
	Files    []*File
}

func (d *directory) SetBasePath(basePath string) Directory {
	d.BasePath = basePath
	return d
}

func (d *directory) GetBasePath() string {
	return d.BasePath
}

func (d *directory) SetName(name string) Directory {
	d.Name = name
	return d
}

func (d *directory) GetName() string {
	return d.Name
}

func (d *directory) GetFullPath() string {
	return fmt.Sprintf("%s/%s", d.GetBasePath(), d.GetName())
}

func (d *directory) SetFiles(files []*File) Directory {
	d.Files = files
	return d
}

func (d *directory) GetFiles() []*File {
	return d.Files
}

func MakeDirectory(name string) girraph.Tree[Directory] {
	return girraph.MakeTree[Directory]().SetMeta(&directory{
		Name:  name,
		Files: []*File{},
	})
}

func SetPaths(basePath string, node girraph.Tree[Directory]) {
	meta := node.GetMeta()
	meta.SetBasePath(basePath)
	path := meta.GetFullPath()
	for _, child := range node.GetChildren() {
		SetPaths(path, child)
	}
}

func mustGenerate(dir Directory) {
	dirPath := fmt.Sprintf("%s/%s", dir.GetBasePath(), dir.GetName())
	utils.PanicError(utils.CreateDir(dirPath))
	for _, file := range dir.GetFiles() {
		err := utils.WriteFile(dirPath, file.GetFullName(), file.GetContents())
		if err != nil {
			panic(err)
		}
	}
}

func MustGenerate(node girraph.Tree[Directory]) {
	p := node.GetMeta()
	mustGenerate(p)
	for _, child := range node.GetChildren() {
		MustGenerate(child)
	}
}
