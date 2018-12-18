package openwcli

import (
	"fmt"
	"github.com/blocktree/OpenWallet/common/file"
	"github.com/blocktree/OpenWallet/hdkeystore"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/OpenWallet/openwallet"
	"github.com/blocktree/OpenWallet/owtp"
	"github.com/blocktree/go-openw-sdk/openwsdk"
	"github.com/bndr/gotabulate"
)

//CreateWalletOnServer
func (cli *CLI) CreateWalletOnServer(name, password string) error {

	var (
		key *hdkeystore.HDKey
	)

	if len(name) == 0 {
		return fmt.Errorf("wallet name is empty. ")
	}

	if len(password) == 0 {
		return fmt.Errorf("wallet password is empty. ")
	}

	//随机生成keystore
	key, filePath, err := hdkeystore.StoreHDKey(
		cli.config.keydir,
		name,
		password,
		hdkeystore.StandardScryptN,
		hdkeystore.StandardScryptP,
	)

	if err != nil {
		return err
	}

	//登记钱包的openw-server
	cli.api.CreateWallet(name, key.KeyID, key.RootPath, true,
		func(status uint64, msg string, wallet *openwsdk.Wallet) {
		if status == owtp.StatusSuccess {
			log.Notice("Wallet create successfully, key path:", filePath)
		} else {
			log.Error("create wallet on server failed, unexpected error:", msg)

			//创建失败，删除key文件
			file.Delete(filePath)
		}
	})

	return nil
}

//GetWalletsByKeyDir 通过给定的文件路径加载keystore文件得到钱包列表
func (cli *CLI) GetWalletsOnServer() ([]*openwsdk.Wallet, error) {
	localWallets, err := openwallet.GetWalletsByKeyDir(cli.config.keydir)
	if err != nil {
		return nil, err
	}
	serverWallets := make([]*openwsdk.Wallet, 0)

	for _, w := range localWallets {
		cli.api.FindWalletByWalletID(w.WalletID, true,
			func(status uint64, msg string, wallet *openwsdk.Wallet) {
			if status == owtp.StatusSuccess && wallet != nil {
				serverWallets = append(serverWallets, wallet)
			}
		})
	}

	return serverWallets, nil
}

//printWalletList 打印钱包列表
func (cli *CLI) printWalletList(list []*openwsdk.Wallet) {

	if list != nil && len(list) > 0 {
		tableInfo := make([][]interface{}, 0)

		for i, w := range list {
			tableInfo = append(tableInfo, []interface{}{
				i, w.WalletID, w.Alias, w.AccountIndex + 1,
			})
		}

		t := gotabulate.Create(tableInfo)
		// Set Headers
		t.SetHeaders([]string{"No.", "ID", "Name", "Accounts"})

		//打印信息
		fmt.Println(t.Render("simple"))
	} else {
		fmt.Println("No wallet was created locally. ")
	}
}