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
remoteserver = "www.openwallet.link"

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

# The custom name of local node
localname = "blocktree"

# Be trusted client server
trustedserver = "client.blocktree.top"

# Enable client server request local transfer
enablerequesttransfer = false

# Enable client server execute summary task
enableexecutesummarytask = false

# Enable key agreement on local node communicate with client server
enablekeyagreement = true

# Enable https or wss
enablessl = false

# Network request timeout, unit: second
requesttimeout = 60

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
	//被托管节点服务地址
	trustedserver string
	//是否接受被托管节点发起的转账请求
	enablerequesttransfer bool
	//是否接受被托管节点执行汇总任务
	enableexecutesummarytask bool
	//是否开启协商密码通信
	enablekeyagreement bool
	//是否支持ssl：https，wss等
	enablessl bool
	//网络请求超时，单位：秒
	requesttimeout int64
	//本地节点自定义的名字
	localname string
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
	conf.trustedserver = c.String("trustedserver")
	conf.localname = c.String("localname")
	conf.enablerequesttransfer, _ = c.Bool("enablerequesttransfer")
	conf.enableexecutesummarytask, _ = c.Bool("enableexecutesummarytask")
	conf.enablekeyagreement, _ = c.Bool("enablekeyagreement")
	conf.enablessl, _ = c.Bool("enablessl")
	conf.requesttimeout, _ = c.Int64("requesttimeout")

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
