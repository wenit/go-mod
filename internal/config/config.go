package config

// Config 配置信息
type Config struct {
	FromDir      string // 源码目录
	TargetDir    string // 本地打包输出目录
	ProxyAddress string // 代理地址，go-mod-server的api服务地址，端口一般为go-mod-server的代理端口+1
	Version      string // 版本号
	Excludes     string // 排除目录，多个目录使用逗号分割
}

// NewConfig 创建配置信息对象
func NewConfig() *Config {
	return &Config{}
}
