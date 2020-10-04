package pack

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Deploy 远程发布
func Deploy(path string, version string, outputDirectory string, excludes string, apiProxyAddr string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("get abs path of module path: %w", err)
	}
	packageInfo := NewPackageInfo()
	packageInfo.Excludes = excludes
	packageInfo.Version = version
	packageInfo.SourceDir = path
	packageInfo.TargetDir = outputDirectory
	packageInfo.APIProxyAddr = apiProxyAddr

	err = install(packageInfo)
	if err != nil {
		return err
	}
	err = upload(packageInfo)
	if err != nil {
		return err
	}
	return nil
}

func upload(packageInfo *PackageInfo) error {
	method := "POST"
	url := packageInfo.APIProxyAddr
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	infofile, infoErrFile := os.Open(packageInfo.InfoFilePath)
	defer infofile.Close()
	infoPart, infoErrFile := writer.CreateFormFile("infoFile", packageInfo.InfoFilePath)
	_, infoErrFile = io.Copy(infoPart, infofile)
	if infoErrFile != nil {
		return infoErrFile
	}

	modFile, modErrFile := os.Open(packageInfo.ModFilePath)
	defer modFile.Close()
	modPart, modErrFile := writer.CreateFormFile("modFile", packageInfo.ModFilePath)
	_, modErrFile = io.Copy(modPart, modFile)
	if modErrFile != nil {
		return modErrFile
	}

	zipFile, zipErrFile := os.Open(packageInfo.ZipFilePath)
	defer zipFile.Close()
	zipPart, zipErrFile := writer.CreateFormFile("zipFile", packageInfo.ZipFilePath)
	_, zipErrFile = io.Copy(zipPart, zipFile)
	if zipErrFile != nil {
		return zipErrFile
	}

	err := writer.Close()
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	status := res.StatusCode
	if status == 200 {
		// log.Printf("上传成功[%s-%s]", component.ArtifactID, component.Version)
		errMsg := fmt.Sprintf("上传成功[%s-%s]", packageInfo.ModName, packageInfo.Version)
		log.Println(errMsg)
	} else {
		errMsg := fmt.Sprintf("上传失败[%s-%s]，错误码：%d", packageInfo.ModName, packageInfo.Version, status)
		log.Println(errMsg)
		body, _ := ioutil.ReadAll(res.Body)
		errMsg = fmt.Sprintf("错误内容\t:%s", string(body))
		log.Println(errMsg)
	}
	return nil
}
