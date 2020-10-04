package pack

import "golang.org/x/mod/modfile"

// PackageInfo 包信息
type PackageInfo struct {
	InfoFilePath    string
	ModFilePath     string
	ZipFilePath     string
	ZipHashFilePath string

	SourceDir string // 源文件目录
	TargetDir string // 打包输出目录
	Excludes  string // 排除目录，多个目录使用逗号分割

	CacheInfoFilePath    string
	CacheModFilePath     string
	CacheZipFilePath     string
	CacheZipHashFilePath string

	ModObj  *modfile.File // module对象信息
	ModName string        // module信息
	Version string        // 版本号

	APIProxyAddr string // 私服上传地址
}

// NewPackageInfo 创建包信息对象
func NewPackageInfo() *PackageInfo {
	return &PackageInfo{}
}
