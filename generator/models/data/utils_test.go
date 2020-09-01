package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakePathData(t *testing.T) {
	result := MakePathData("/base/path", "dir")

	assert.Equal(t, "/base/path", result.Base)
	assert.Equal(t, "/base/path/dir", result.Full)
}

func TestMakeGoFileData(t *testing.T) {
	fileResult, pathResult := MakeGoFileData("/base/path", "filename")

	assert.Equal(t, "filename", fileResult.Base)
	assert.Equal(t, "filename.go", fileResult.Full)

	assert.Equal(t, "/base/path", pathResult.Base)
	assert.Equal(t, "/base/path/filename.go", pathResult.Full)
}

func TestMakePackageData(t *testing.T) {
	pkgName := "pkgname"
	result := MakePackageData("github.com/user/example", pkgName)
	expectedName := MakeName(pkgName)

	assert.Equal(t, expectedName, result.Name)
	assert.Equal(t, "pkgname", result.Reference)
	assert.Equal(t, "github.com/user/example", result.Path.Base)
	assert.Equal(t, "github.com/user/example/pkgname", result.Path.Full)
}

func TestMakePackageData_SnakeName(t *testing.T) {
	pkgName := "pkg_name"
	result := MakePackageData("github.com/user/example", pkgName)
	expectedName := MakeName(pkgName)

	assert.Equal(t, expectedName, result.Name)
	assert.Equal(t, "pkg_name", result.Reference)
	assert.Equal(t, "github.com/user/example", result.Path.Base)
	assert.Equal(t, "github.com/user/example/pkg_name", result.Path.Full)
}
