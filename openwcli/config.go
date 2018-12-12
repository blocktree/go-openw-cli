package openwcli

import (
	"github.com/astaxie/beego/config"
)

// 默认配置
var (
	defaultConfig = `

# Remote Server
remoteserver = "www.openwallet.site"

# API Version
version = "1.0.0"

# App Key
appkey = "1234qwer"

# App Secret
appsecret = "qwer1234"

# Log file path
logdir = "/usr/logs/"

# Data directory, store keys, databases, backups
datadir = "/usr/data/"

# Wallet Summary Period
summaryperiod = "1h"

`
)

//配置
type Config struct {

	// 远程服务
	remoteserver string
	//版本号
	version string
	//应用key
	appkey string
	//应用密钥
	appsecret string
	//日期路径
	logdir string
	//数据缓存路径
	datadir string
	//汇总时间定时
	summaryperiod string
}

//初始化一个配置对象
func NewConfig(c config.Configer) *Config {
	conf := &Config{}
	conf.remoteserver = c.String("remoteserver")
	conf.version = c.String("version")
	conf.appkey = c.String("appkey")
	conf.appsecret = c.String("appsecret")
	conf.logdir = c.String("logdir")
	conf.datadir = c.String("datadir")
	conf.summaryperiod = c.String("summaryperiod")
	return conf
}


// 加载工具配置
func LoadConfig(path string) (*Config, error) {

	c, err := config.NewConfig("ini", path)
	if err != nil {
		return nil, err
	}

	conf := NewConfig(c)
	return conf, nil
}
