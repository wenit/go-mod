package common

import (
	"testing"
)

func TestGetGoPath(t *testing.T) {

	t.Log(GetGoPath())
}

func TestGetGoModulePath(t *testing.T) {
	t.Log(GetGoModulePath())
}

func TestPathExists(t *testing.T) {
	t.Log(PathExists("c:/1"))
	t.Log(PathExists("."))
}

func TestZip(t *testing.T) {
	fileDir := "F:/github/go-mod"
	// zipFile := "d:/1.zip"
	zipFile := "F:/github/go-mod/target/1.zip"

	err := Zip(fileDir, zipFile)
	// err := Zip(fileDir, zipFile)

	if err != nil {
		t.Log("压缩文件错误：", err)
	}
}

func TestZipFilter(t *testing.T) {
	fileDir := "F:/github/go-mod"
	// zipFile := "d:/1.zip"
	zipFile := "F:/github/go-mod/target/1.zip"

	filter := []string{"releases", ".git", ".vscode"}

	err := ZipFilter(fileDir, zipFile, "github.com/wenit/go-mod@v1.0.1", filter)
	// err := Zip(fileDir, zipFile)

	if err != nil {
		t.Log("压缩文件错误：", err)
	}
}

func TestUnzip(t *testing.T) {
	zipFile := "F:/github/go-mod/target/1.zip"

	distDir := "F:/github/go-mod/target/temp"

	err := Unzip(zipFile, distDir)

	if err != nil {
		t.Log("解压文件错误：", err)
	}
}