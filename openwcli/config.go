package openwcli

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/OpenWallet/common/file"
	"path/filepath"
)

// 默认配置
const (
	defaultConfig = `

# Remote Server
remoteserver = "www.openwallet.site"

# API Version
version = "1.0.0"

# App Key
appkey = "1234qwer"

# App ID
appid = "qwer1234"

# Log file path
logdir = "/usr/logs/"

# Data directory, store keys, databases, backups
datadir = "/usr/data/"

# Wallet Summary Period
summaryperiod = "1h"

`

	keyDirName = "key"
	dbDirName = "db"
)

//配置
type Config struct {

	// 远程服务
	remoteserver string
	//版本号
	version string
	//应用key
	appkey string
	//应用ID
	appid string
	//日期路径
	logdir string
	//数据缓存路径
	datadir string
	//汇总时间定时
	summaryperiod string
	//密钥目录
	keydir string
	//数据库目录
	dbdir string
}

//初始化一个配置对象
func NewConfig(c config.Configer) *Config {
	conf := &Config{}
	conf.remoteserver = c.String("remoteserver")
	conf.version = c.String("version")
	conf.appkey = c.String("appkey")
	conf.appid = c.String("appid")
	conf.logdir = c.String("logdir")
	conf.datadir = c.String("datadir")
	conf.summaryperiod = c.String("summaryperiod")

	conf.keydir = filepath.Join(conf.datadir, keyDirName)
	conf.dbdir =filepath.Join(conf.datadir, dbDirName)

	//建立文件夹
	file.MkdirAll(conf.datadir)
	file.MkdirAll(conf.logdir)
	file.MkdirAll(conf.keydir)
	file.MkdirAll(conf.dbdir)

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
