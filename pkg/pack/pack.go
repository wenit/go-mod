package pack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wenit/go-mod/pkg/common"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/zip"
)

// Package 打包
func Package(path string, version string, outputDirectory string, excludes string) error {
	_, err := module(path, version, outputDirectory, excludes)
	if err != nil {
		return fmt.Errorf("get module file: %w", err)
	}
	return nil
}

func module(path string, version string, outputDirectory string, excludes string) (*modfile.File, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("get abs path of module path: %w", err)
	}

	if !common.PathExists(outputDirectory) {
		err := common.MkDirs(outputDirectory)
		if err != nil {
			return nil, fmt.Errorf("create output directory: %s,error %w", outputDirectory, err)
		}
	}

	moduleFile, err := getModuleFile(path, version)
	if err != nil {
		return nil, fmt.Errorf("get module file: %w", err)
	}

	// common.ZipFilter()

	if err := createZipArchiveCommon(path, moduleFile, outputDirectory, excludes); err != nil {
		return nil, fmt.Errorf("create zip archive: %w", err)
	}

	if err := createInfoFile(moduleFile, outputDirectory); err != nil {
		return nil, fmt.Errorf("create info file: %w", err)
	}

	if err := copyModuleFile(path, moduleFile, outputDirectory); err != nil {
		return nil, fmt.Errorf("copy module file: %w", err)
	}

	return moduleFile, nil
}

// Install 本地安装
func Install(path string, version string, outputDirectory string, excludes string) error {
	moduleFile, err := module(path, version, outputDirectory, excludes)
	if err != nil {
		return fmt.Errorf("get module file: %w", err)
	}

	zipFile := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".zip")

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

	err = common.Unzip(zipFile, modulePath)
	if err != nil {
		return err
	}

	// github.com/wenit/go-mod\@v
	prefix := fmt.Sprintf("%s/@v", moduleFile.Module.Mod.Path)
	downloadPath := common.GetGoModuleCacheDownloadPath()

	srcInfoFile := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".info")
	srcModFile := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".mod")
	srcZipFile := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".zip")

	dstInfoFile := filepath.Join(downloadPath, prefix, moduleFile.Module.Mod.Version+".info")
	dstModFile := filepath.Join(downloadPath, prefix, moduleFile.Module.Mod.Version+".mod")
	dstZipFile := filepath.Join(downloadPath, prefix, moduleFile.Module.Mod.Version+".zip")
	// 1、copy info
	common.CopyFile(srcInfoFile, dstInfoFile)
	// 2、copy mod
	common.CopyFile(srcModFile, dstModFile)
	// 3、copy zip
	common.CopyFile(srcZipFile, dstZipFile)

	return nil
}

func getModuleFile(path string, version string) (*modfile.File, error) {
	path = filepath.Join(path, "go.mod")
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open module file: %w", err)
	}
	defer file.Close()

	moduleBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read module file: %w", err)
	}

	moduleFile, err := modfile.Parse(path, moduleBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("parse module file: %w", err)
	}

	if moduleFile.Module == nil {
		return nil, fmt.Errorf("parsing module returned nil module")
	}

	moduleFile.Module.Mod.Version = version

	return moduleFile, nil
}

func createZipArchiveCommon(path string, moduleFile *modfile.File, outputDirectory string, excludes string) error {
	outputPath := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".zip")
	prefix := fmt.Sprintf("%s@%s/", moduleFile.Module.Mod.Path, moduleFile.Module.Mod.Version)
	filter := strings.Split(excludes, ",")
	return common.ZipFilter(path, outputPath, prefix, filter)
}
func createZipArchive(path string, moduleFile *modfile.File, outputDirectory string) error {
	outputPath := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".zip")

	var zipContents bytes.Buffer
	if err := zip.CreateFromDir(&zipContents, moduleFile.Module.Mod, path); err != nil {
		return fmt.Errorf("create zip from dir: %w", err)
	}

	if err := ioutil.WriteFile(outputPath, zipContents.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing zip file: %w", err)
	}

	return nil
}

func createInfoFile(moduleFile *modfile.File, outputDirectory string) error {
	infoFilePath := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".info")
	file, err := os.Create(infoFilePath)
	if err != nil {
		return fmt.Errorf("create info file: %w", err)
	}
	defer file.Close()

	type infoFile struct {
		Version string
		Time    string
	}

	currentTime := getInfoFileFormattedTime(time.Now())
	info := infoFile{
		Version: moduleFile.Module.Mod.Version,
		Time:    currentTime,
	}

	infoBytes, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal info file: %w", err)
	}

	if _, err := file.Write(infoBytes); err != nil {
		return fmt.Errorf("write info file: %w", err)
	}

	return nil
}

func getInfoFileFormattedTime(currentTime time.Time) string {
	const infoFileTimeFormat = "2006-01-02T15:04:05Z"
	return currentTime.Format(infoFileTimeFormat)
}

func copyModuleFile(path string, moduleFile *modfile.File, outputDirectory string) error {
	if outputDirectory == "." {
		return nil
	}

	sourcePath := filepath.Join(path, "go.mod")
	destinationPath := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".mod")

	if sourcePath == destinationPath {
		return nil
	}

	moduleContents, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("read module file: %w", err)
	}

	if err := ioutil.WriteFile(destinationPath, moduleContents, 0644); err != nil {
		return fmt.Errorf("write module file: %w", err)
	}

	return nil
}
