package pack

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/wenit/go-mod/pkg/common"
)

// Install 本地安装
func Install(path string, version string, outputDirectory string, excludes string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("get abs path of module path: %w", err)
	}
	packageInfo := NewPackageInfo()
	packageInfo.Excludes = excludes
	packageInfo.Version = version
	packageInfo.SourceDir = path
	packageInfo.TargetDir = outputDirectory

	err = install(packageInfo)
	if err != nil {
		return err
	}
	return nil
}

// install 本地安装
func install(packageInfo *PackageInfo) error {
	version := packageInfo.Version
	outputDirectory := packageInfo.TargetDir

	err := module(packageInfo)
	moduleFile := packageInfo.ModObj
	if err != nil {
		return fmt.Errorf("get module file: %w", err)
	}
	log.Printf("项目打包完成，输出目录：%s", outputDirectory)

	zipFile := filepath.Join(outputDirectory, version+".zip")

	modulePath := common.GetGoModulePath()
	if modulePath == "" {
		return fmt.Errorf("get module root dir error")
	}

	if !common.PathExists(modulePath) {
		err := common.MkDirs(modulePath)
		if err != nil {
			return err
		}
	}

	// 解压文件到 $GOPATH/pkg/mod
	err = common.Unzip(zipFile, modulePath)
	if err != nil {
		return err
	}
	log.Printf("项目解压至本地mod仓库，输出目录：%s", modulePath)

	prefix := fmt.Sprintf("%s/@v", moduleFile.Module.Mod.Path)
	downloadPath := common.GetGoModuleCacheDownloadPath()

	srcInfoFile := filepath.Join(outputDirectory, version+".info")
	srcModFile := filepath.Join(outputDirectory, version+".mod")
	srcZipFile := filepath.Join(outputDirectory, version+".zip")
	srcZiphashFile := filepath.Join(outputDirectory, version+".ziphash")

	dstInfoFile := filepath.Join(downloadPath, prefix, version+".info")
	dstModFile := filepath.Join(downloadPath, prefix, version+".mod")
	dstZipFile := filepath.Join(downloadPath, prefix, version+".zip")
	dstZiphashFile := filepath.Join(downloadPath, prefix, version+".ziphash")

	packageInfo.CacheInfoFilePath = dstInfoFile
	packageInfo.CacheModFilePath = dstModFile
	packageInfo.CacheZipFilePath = dstZipFile
	packageInfo.CacheZipHashFilePath = dstZiphashFile

	// copy文件至缓存目录 ： $GOPATH/pkg/mod/cache/download
	// 1、copy info
	err = common.CopyFile(srcInfoFile, dstInfoFile)
	if err != nil {
		return err
	}
	log.Printf("复制info文件至缓存目录[%s]完成", dstZipFile)
	// 2、copy mod
	err = common.CopyFile(srcModFile, dstModFile)
	if err != nil {
		return err
	}
	log.Printf("复制mod文件至缓存目录[%s]完成", dstZipFile)
	// 3、copy zip
	err = common.CopyFile(srcZipFile, dstZipFile)
	if err != nil {
		return err
	}
	log.Printf("复制zip文件至缓存目录[%s]完成", dstZipFile)

	// 4、copy ziphash
	err = common.CopyFile(srcZiphashFile, dstZiphashFile)
	if err != nil {
		return err
	}
	log.Printf("复制ziphash文件至缓存目录[%s]完成", dstZiphashFile)

	return nil
}
