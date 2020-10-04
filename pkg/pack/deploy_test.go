package pack

import (
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestDeploy(t *testing.T) {
	path := "f:/github/go-mod"
	version := "v1.1.1"
	outputDirectory := "f:/github/go-mod/target"
	excludes := ".svn,.git,.vscode,target,releases"
	apiProxyAddr := "http://localhost:8082/upload"
	err := Deploy(path, version, outputDirectory, excludes, apiProxyAddr)
	if err != nil {
		t.Log("部署失败：", err)
	}
}

func TestGetProxyURL(t *testing.T) {
	goproxyURL := os.Getenv("GOPROXY")

	proxyAddrs := strings.Split(goproxyURL, ",")

	// firstProxy := proxyAddrs[0]
	firstProxy := "https://mirrors.aliyun.com:8081/goproxy/"

	u, err := url.Parse(firstProxy)
	if err != nil {
		t.Log(err)
	}

	t.Log(u.Scheme)
	t.Log(u.Host)
	t.Log(u.Port())
	t.Log(u.Path)

	t.Log(goproxyURL)
	t.Log(firstProxy)
	t.Log(proxyAddrs[0])

}
