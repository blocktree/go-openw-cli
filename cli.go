package openwcli

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/blocktree/OpenWallet/openwallet"
	"github.com/coreos/bbolt"
	"path/filepath"
	"time"
)

type CLI struct {
	//工具配置
	config *Config
	db     *openwallet.StormDB
}

// 初始化工具
func NewCLI(conf string) (*CLI, error) {

	//加载配置
	c, err := loadConfig(conf)
	if err != nil {
		return nil, err
	}

	if len(c.appkey) == 0 {
		return nil, fmt.Errorf("appkey is empty. ")
	}

	if len(c.appsecret) == 0 {
		return nil, fmt.Errorf("appsecret is empty. ")
	}

	if len(c.remoteserver) == 0 {
		return nil, fmt.Errorf("remoteserver is empty. ")
	}

	dbfile := filepath.Join(c.datadir, c.appkey+".db")

	//加载数据
	db, err := openwallet.OpenStormDB(
		dbfile,
		storm.BoltOptions(0600, &bolt.Options{Timeout: 3 * time.Second}),
	)


	if err != nil {
		return nil, err
	}

	cli := &CLI{
		config: c,
		db:     db,
	}
	return cli, nil
}

func (cli *CLI) GetKeyPair() {

}
