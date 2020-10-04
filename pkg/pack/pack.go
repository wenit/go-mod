package pack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("get abs path of module path: %w", err)
	}
	packageInfo := NewPackageInfo()
	packageInfo.Excludes = excludes
	packageInfo.Version = version
	packageInfo.SourceDir = path
	packageInfo.TargetDir = outputDirectory

	err = module(packageInfo)
	if err != nil {
		return fmt.Errorf("get module file: %w", err)
	}
	log.Printf("项目打包完成，输出目录：%s", outputDirectory)
	return nil
}

func pack(packageInfo *PackageInfo) error {
	return module(packageInfo)
}
func module(packageInfo *PackageInfo) error {
	outputDirectory := packageInfo.TargetDir
	path := packageInfo.SourceDir
	version := packageInfo.Version
	if !common.PathExists(outputDirectory) {
		err := common.MkDirs(outputDirectory)
		if err != nil {
			return fmt.Errorf("create output directory: %s,error %w", outputDirectory, err)
		}
	}

	moduleFile, err := getModuleFile(path, version)
	if err != nil {
		return fmt.Errorf("get module file: %w", err)
	}
	packageInfo.ModObj = moduleFile
	packageInfo.ModName = moduleFile.Module.Mod.Path

	if err := createZipArchiveCommon(packageInfo); err != nil {
		return fmt.Errorf("create zip archive: %w", err)
	}

	if err := createInfoFile(packageInfo); err != nil {
		return fmt.Errorf("create info file: %w", err)
	}

	if err := copyModuleFile(packageInfo); err != nil {
		return fmt.Errorf("copy module file: %w", err)
	}

	if err := createZiphash(packageInfo); err != nil {
		return fmt.Errorf("createZiphash file: %w", err)
	}
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

func createZipArchiveCommon(packageInfo *PackageInfo) error {
	outputDirectory := packageInfo.TargetDir
	path := packageInfo.SourceDir
	version := packageInfo.Version
	excludes := packageInfo.Excludes
	outputPath := filepath.Join(outputDirectory, version+".zip")
	packageInfo.ZipFilePath = outputPath
	prefix := fmt.Sprintf("%s@%s/", packageInfo.ModObj.Module.Mod.Path, version)
	filter := strings.Split(excludes, ",")
	return common.ZipFilter(path, outputPath, prefix, filter)
}

// func createZipArchive(path string, moduleFile *modfile.File, outputDirectory string) error {
func createZipArchive(packageInfo *PackageInfo) error {
	outputDirectory := packageInfo.TargetDir
	version := packageInfo.Version
	outputPath := filepath.Join(outputDirectory, version+".zip")

	var zipContents bytes.Buffer
	if err := zip.CreateFromDir(&zipContents, packageInfo.ModObj.Module.Mod, packageInfo.SourceDir); err != nil {
		return fmt.Errorf("create zip from dir: %w", err)
	}

	if err := ioutil.WriteFile(outputPath, zipContents.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing zip file: %w", err)
	}

	return nil
}

func createInfoFile(packageInfo *PackageInfo) error {
	outputDirectory := packageInfo.TargetDir
	version := packageInfo.Version
	infoFilePath := filepath.Join(outputDirectory, version+".info")
	packageInfo.InfoFilePath = infoFilePath
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
		Version: version,
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

func createZiphash(packageInfo *PackageInfo) error {
	outputDirectory := packageInfo.TargetDir
	version := packageInfo.Version
	zipFilePath := filepath.Join(outputDirectory, version+".zip")
	hashFilePath := filepath.Join(outputDirectory, version+".ziphash")
	packageInfo.ZipFilePath = zipFilePath
	packageInfo.ZipHashFilePath = hashFilePath

	file, err := os.Create(hashFilePath)
	if err != nil {
		return fmt.Errorf("create hash file: %w", err)
	}
	defer file.Close()

	hashValue := common.ZipHash(zipFilePath)
	if hashValue == "" {
		return fmt.Errorf("generate ziphash file: %w", err)
	}

	if _, err := file.Write([]byte(hashValue)); err != nil {
		return fmt.Errorf("write ziphash file: %w", err)
	}

	return nil
}

func getInfoFileFormattedTime(currentTime time.Time) string {
	const infoFileTimeFormat = "2006-01-02T15:04:05Z"
	return currentTime.Format(infoFileTimeFormat)
}

func copyModuleFile(packageInfo *PackageInfo) error {
	outputDirectory := packageInfo.TargetDir
	version := packageInfo.Version
	path := packageInfo.SourceDir
	if outputDirectory == "." {
		return nil
	}

	sourcePath := filepath.Join(path, "go.mod")
	destinationPath := filepath.Join(outputDirectory, version+".mod")

	packageInfo.ModFilePath = destinationPath

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
