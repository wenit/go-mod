package pack

import "testing"

func TestInstall(t *testing.T) {
	path := "f:/github/go-mod"
	version := "v1.1.1"
	outputDirectory := "f:/github/go-mod/target"
	excludes := ".svn,.git,.vscode,target,releases"
	err := Install(path, version, outputDirectory, excludes)
	if err != nil {
		t.Log("安装失败：", err)
	}
}
