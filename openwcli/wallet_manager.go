package openwcli

import (
	"encoding/hex"
	"fmt"
	"github.com/asdine/storm/q"
	"github.com/blocktree/OpenWallet/common"
	"github.com/blocktree/OpenWallet/common/file"
	"github.com/blocktree/OpenWallet/hdkeystore"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/OpenWallet/openwallet"
	"github.com/blocktree/OpenWallet/owtp"
	"github.com/blocktree/go-openw-sdk/openwsdk"
	"github.com/bndr/gotabulate"
	"path/filepath"
	"strings"
	"time"
)

//CreateWalletOnServer
func (cli *CLI) CreateWalletOnServer(name, password string) (*openwsdk.Wallet, error) {

	var (
		key       *hdkeystore.HDKey
		retWallet *openwsdk.Wallet
	)

	if len(name) == 0 {
		return nil, fmt.Errorf("wallet name is empty. ")
	}

	if len(password) == 0 {
		return nil, fmt.Errorf("wallet password is empty. ")
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
		return nil, err
	}

	//登记钱包的openw-server
	cli.api.CreateWallet(name, key.KeyID, true,
		func(status uint64, msg string, wallet *openwsdk.Wallet) {
			if status == owtp.StatusSuccess {
				log.Info("Wallet create successfully, key path:", filePath)
				retWallet = wallet
			} else {
				log.Error("create wallet on server failed, unexpected error:", msg)

				//创建失败，删除key文件
				file.Delete(filePath)
			}
		})

	return retWallet, nil
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

//GetWalletByWalletID 查找本地且线上有的钱包对象
func (cli *CLI) GetWalletByWalletID(walletID string) (*openwsdk.Wallet, error) {

	var (
		findErr error
	)

	localWallet, err := cli.GetWalletByWalletIDOnLocal(walletID)
	if err != nil {
		return nil, err
	}

	err = cli.api.FindWalletByWalletID(walletID, true,
		func(status uint64, msg string, wallet *openwsdk.Wallet) {
			if status == owtp.StatusSuccess && wallet != nil {
				localWallet = wallet
			} else {
				findErr = fmt.Errorf(msg)
			}
		})
	if err != nil {
		return nil, err
	}

	if findErr != nil {
		return nil, findErr
	}

	return localWallet, nil
}

//GetWalletByWalletIDOnLocal 查找本地种子目录的钱包对象
func (cli *CLI) GetWalletByWalletIDOnLocal(walletID string) (*openwsdk.Wallet, error) {
	localWallets, err := openwallet.GetWalletsByKeyDir(cli.config.keydir)
	if err != nil {
		return nil, err
	}
	for _, w := range localWallets {
		if w.WalletID == walletID {
			selectedWallet := &openwsdk.Wallet{
				WalletID: w.WalletID,
				Alias:    w.Alias,
			}
			return selectedWallet, nil
		}

	}

	return nil, fmt.Errorf("can not find local wallet by walletID: %s", walletID)
}

//printWalletList 打印钱包列表
func (cli *CLI) printWalletList(list []*openwsdk.Wallet) {

	if list != nil && len(list) > 0 {
		tableInfo := make([][]interface{}, 0)

		for i, w := range list {
			tableInfo = append(tableInfo, []interface{}{
				i, w.Alias, w.WalletID, w.AccountIndex + 1,
			})
		}

		t := gotabulate.Create(tableInfo)
		// Set Headers
		t.SetHeaders([]string{"No.", "Name", "WalletID", "Accounts"})

		//打印信息
		fmt.Println(t.Render("simple"))
	} else {
		fmt.Println("No wallet was created locally. ")
	}
}

//CreateAccountOnServer
func (cli *CLI) CreateAccountOnServer(name, password, symbol string, wallet *openwsdk.Wallet) (*openwsdk.Account, []*openwsdk.Address, error) {

	var (
		key            *hdkeystore.HDKey
		selectedSymbol *openwsdk.Symbol
		retAccount     *openwsdk.Account
		retAddresses   []*openwsdk.Address
		err            error
		retErr         error
	)

	if len(name) == 0 {
		return nil, nil, fmt.Errorf("acount name is empty. ")
	}

	if len(password) == 0 {
		return nil, nil, fmt.Errorf("wallet password is empty. ")
	}

	selectedSymbol, err = cli.GetSymbolInfo(symbol)
	if err != nil {
		return nil, nil, err
	}

	keystore := hdkeystore.NewHDKeystore(
		cli.config.keydir,
		hdkeystore.StandardScryptN,
		hdkeystore.StandardScryptP,
	)

	fileName := fmt.Sprintf("%s-%s.key", wallet.Alias, wallet.WalletID)

	key, err = keystore.GetKey(
		wallet.WalletID,
		fileName,
		password,
	)
	if err != nil {
		return nil, nil, err
	}

	newaccount, err := wallet.CreateAccount(name, selectedSymbol, key)
	if err != nil {
		return nil, nil, err
	}

	//登记钱包的openw-server
	err = cli.api.CreateNormalAccount(newaccount, true,
		func(status uint64, msg string, account *openwsdk.Account, addresses []*openwsdk.Address) {
			if status == owtp.StatusSuccess {
				log.Infof("create [%s] account successfully", selectedSymbol.Coin)
				log.Infof("new accountID: %s", account.AccountID)
				if len(addresses) > 0 {
					log.Infof("new address: %s", addresses[0].Address)
				}

				retAccount = account
				retAddresses = addresses
			} else {
				log.Error("create account on server failed, unexpected error:", msg)
				retErr = fmt.Errorf(msg)
			}
		})

	if err != nil {
		return nil, nil, err
	}
	if retErr != nil {
		return nil, nil, retErr
	}

	return retAccount, retAddresses, nil
}

//GetAccountOnServerByAccountID 从服务器获取账户
func (cli *CLI) GetAccountByAccountID(accountID string) (*openwsdk.Account, error) {

	var (
		getAccount *openwsdk.Account
		err        error
		retErr     error
	)

	err = cli.api.FindAccountByAccountID(accountID, true,
		func(status uint64, msg string, account *openwsdk.Account) {
			if status == owtp.StatusSuccess {
				getAccount = account
			} else {
				retErr = fmt.Errorf(msg)
			}
		})
	if err != nil {
		return nil, err
	}
	if retErr != nil {
		return nil, retErr
	}

	return getAccount, nil
}

//GetAccountsOnServer 从服务器获取账户列表
func (cli *CLI) GetAccountsOnServer(walletID string) ([]*openwsdk.Account, error) {

	var (
		list   = make([]*openwsdk.Account, 0)
		err    error
		retErr error
	)

	err = cli.api.FindAccountByWalletID(walletID, true,
		func(status uint64, msg string, accounts []*openwsdk.Account) {
			if status == owtp.StatusSuccess && len(accounts) > 0 {
				list = append(list, accounts...)
			} else {
				retErr = fmt.Errorf(msg)
			}
		})

	if err != nil {
		return nil, err
	}
	if retErr != nil {
		return nil, retErr
	}

	return list, nil
}

//printAccountList 打印账户列表
func (cli *CLI) printAccountList(list []*openwsdk.Account) {

	if list != nil && len(list) > 0 {
		tableInfo := make([][]interface{}, 0)

		for i, w := range list {

			//读取汇总信息
			var sum openwsdk.SummarySetting
			cli.db.One("AccountID", w.AccountID, &sum)

			tableInfo = append(tableInfo, []interface{}{
				i, w.Alias, w.AccountID, w.Symbol, w.Balance, w.AddressIndex + 1,
				sum.SumAddress, sum.Threshold, sum.MinTransfer, sum.RetainedBalance, sum.Confirms,
			})
		}

		t := gotabulate.Create(tableInfo)
		// Set Headers
		t.SetHeaders([]string{"No.", "Name", "AccountID", "Symbol", "Balance", "Addresses",
			"Summary Address", "Summary Threshold", "Min Transfer", "Retained Balance", "Confirms"})

		//打印信息
		fmt.Println(t.Render("simple"))
	} else {
		fmt.Println("No account was created locally. ")
	}
}

//CreateAddressOnServer
func (cli *CLI) CreateAddressOnServer(walletID, accountID string, count uint64) error {

	if len(accountID) == 0 {
		return fmt.Errorf("accountID is empty. ")
	}

	if len(walletID) == 0 {
		return fmt.Errorf("walleID is empty. ")
	}

	if count == 0 {
		return fmt.Errorf("create address count can not 0. ")
	}

	err := cli.api.CreateBatchAddress(walletID, accountID, count, true,
		func(status uint64, msg string, addresses []string) {
			if status == owtp.StatusSuccess {
				log.Infof("create [%d] addresses successfully", len(addresses))
				//:保存到本地数据库，导出到文件夹
				timestamp := time.Now()
				filename := "["+accountID+"]-" + common.TimeFormat("20060102150405", timestamp) + ".txt"
				filePath := filepath.Join(cli.config.exportaddressdir, filename)
				if flag := cli.exportAddressToFile(addresses, filePath); flag {
					log.Infof("addresses has been exported into: %s", filePath)
				} else {
					log.Infof("addresses export failed")
				}
			} else {
				log.Error("create account on server failed, unexpected error:", msg)
			}
		})

	return err
}

//exportAddressToFile 导出地址到文件中
func (cli *CLI) exportAddressToFile(addresses []string, filePath string) bool {

	var (
		content string
	)

	for _, a := range addresses {
		content = content + a + "\n"
	}

	return file.WriteFile(filePath, []byte(content), true)
}

//SearchAddressOnServer
func (cli *CLI) SearchAddressOnServer(address string) (*openwsdk.Address, error) {

	if len(address) == 0 {
		return nil, fmt.Errorf("address is empty. ")
	}

	var addr *openwsdk.Address

	cli.api.FindAddressByAddress(address, true,
		func(status uint64, msg string, address *openwsdk.Address) {
			if status == owtp.StatusSuccess {
				addr = address
			} else {
				log.Error("search address on server failed, unexpected error:", msg)
			}
		})

	return addr, nil
}

//GetAddressesOnServer
func (cli *CLI) GetAddressesOnServer(walletID, accountID string, offset, limit int) ([]*openwsdk.Address, error) {

	list := make([]*openwsdk.Address, 0)

	if len(accountID) == 0 {
		return nil, fmt.Errorf("accountID is empty. ")
	}

	if len(walletID) == 0 {
		return nil, fmt.Errorf("walleID is empty. ")
	}

	cli.api.FindAddressByAccountID(accountID, offset, limit, true,
		func(status uint64, msg string, addresses []*openwsdk.Address) {
			if status == owtp.StatusSuccess {
				list = addresses
			} else {
				log.Error("get address on server failed, unexpected error:", msg)
			}
		})

	return list, nil
}

//printAddressList 打印地址列表
func (cli *CLI) printAddressList(walletID string, list []*openwsdk.Address, password string) error {

	var (
		isShowPrivateKey bool
		privatekey       = ""
		key              *hdkeystore.HDKey
	)

	if len(password) != 0 {
		isShowPrivateKey = true
	}

	if isShowPrivateKey {

		keystore := hdkeystore.NewHDKeystore(
			cli.config.keydir,
			hdkeystore.StandardScryptN,
			hdkeystore.StandardScryptP,
		)

		w, err := cli.GetWalletByWalletIDOnLocal(walletID)
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("%s-%s.key", w.Alias, w.WalletID)

		key, err = keystore.GetKey(
			w.WalletID,
			fileName,
			password,
		)
		if err != nil {
			return err
		}
	}

	if list != nil && len(list) > 0 {
		tableInfo := make([][]interface{}, 0)

		for i, a := range list {

			if isShowPrivateKey && key != nil {

				selectedSymbol, err := cli.GetSymbolInfo(a.Symbol)
				if err != nil {
					return err
				}

				extKey, err := key.DerivedKeyWithPath(a.HdPath, uint32(selectedSymbol.Curve))
				if err != nil {
					return err
				}

				privateKeyBytes, err := extKey.GetPrivateKeyBytes()
				if err != nil {
					return err
				}

				privatekey = hex.EncodeToString(privateKeyBytes)
			}

			tableInfo = append(tableInfo, []interface{}{
				i, a.Address, a.WalletID, a.AccountID, a.Symbol, a.Balance, a.PublicKey, privatekey,
			})

		}
		t := gotabulate.Create(tableInfo)
		// Set Headers
		t.SetHeaders([]string{"No.", "Address", "WalletID", "AccounttID", "Symbol", "Balance", "publicKey", "privateKey"})

		//打印信息
		fmt.Println(t.Render("simple"))
	} else {
		fmt.Println("No address was created locally. ")
	}

	return nil
}

//UpdateSymbols 更新主链
func (cli *CLI) UpdateSymbols() error {

	const (
		limit = 500
	)

	var getSymbols []*openwsdk.Symbol
	err := cli.api.GetSymbolList(0, limit, true,
		func(status uint64, msg string, symbols []*openwsdk.Symbol) {
			getSymbols = symbols
		})
	if err != nil {
		return err
	}

	tx, err := cli.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, s := range getSymbols {

		i := 0
		for {

			var getTokenContract []*openwsdk.TokenContract
			err = cli.api.GetContracts(s.Coin, i, limit, true,
				func(status uint64, msg string, tokenContract []*openwsdk.TokenContract) {
					getTokenContract = tokenContract
				})
			if err != nil || len(getTokenContract) == 0 {
				break
			}

			//保存主链上的合约信息
			for _, c := range getTokenContract {
				err = tx.Save(c)
				if err != nil {
					return err
				}
			}

			i = i + limit

		}

		//保存主链信息
		err = tx.Save(s)
		if err != nil {

			return err
		}

	}

	return tx.Commit()
}

//UpdateSymbols 更新主链
func (cli *CLI) UpdateTokenContracts(symbol string) error {
	var getTokenContract []*openwsdk.TokenContract
	err := cli.api.GetContracts(symbol, 0, 5000, true,
		func(status uint64, msg string, tokenContract []*openwsdk.TokenContract) {
			getTokenContract = tokenContract
		})
	if err != nil {
		return err
	}

	tx, err := cli.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, s := range getTokenContract {

		err = tx.Save(s)
		if err != nil {

			return err
		}

	}

	return tx.Commit()
}

//GetLocalSymbolList 查询本地保存主链
func (cli *CLI) GetSymbolList() ([]*openwsdk.Symbol, error) {
	var getSymbols []*openwsdk.Symbol
	err := cli.db.All(&getSymbols)

	//没有数据，更新数据
	if getSymbols == nil || len(getSymbols) == 0 {
		err = cli.UpdateSymbols()
		if err != nil {
			return nil, err
		}

		err = cli.db.All(&getSymbols)
		if err != nil {
			return nil, err
		}

	}
	return getSymbols, nil
}

//printSymbolList 打印主链列表
func (cli *CLI) printSymbolList(list []*openwsdk.Symbol) {

	if list != nil && len(list) > 0 {
		tableInfo := make([][]interface{}, 0)

		for _, w := range list {
			tableInfo = append(tableInfo, []interface{}{
				w.Name, w.Coin, w.Curve, w.Decimals,
			})
		}

		t := gotabulate.Create(tableInfo)
		// Set Headers
		t.SetHeaders([]string{"Name", "Symbol", "ECC Type", "Decimals"})

		//打印信息
		fmt.Println(t.Render("simple"))
	} else {
		fmt.Println("No Symbol. ")
	}
}

//GetLocalSymbolInfo 查询本地主链信息
func (cli *CLI) GetSymbolInfo(symbol string) (*openwsdk.Symbol, error) {

	getSymbols, err := cli.GetSymbolList()
	if err != nil {
		return nil, err
	}

	for _, s := range getSymbols {
		if s.Coin == strings.ToUpper(symbol) {
			return s, nil
		}
	}

	return nil, fmt.Errorf("can not find symbol info")
}

//GetContractList 查询本地保存代币合约信息
func (cli *CLI) GetTokenContractList(cols ...interface{}) ([]*openwsdk.TokenContract, error) {
	query := make([]q.Matcher, 0)
	var getTokenContracts []*openwsdk.TokenContract

	if len(cols)%2 != 0 {
		return nil, fmt.Errorf("condition param is not pair")
	}

	for i := 0; i < len(cols); i = i + 2 {
		field := common.NewString(cols[i])
		val := cols[i+1]
		query = append(query, q.Eq(field.String(), val))
	}

	err := cli.db.Select(q.And(query...)).Find(&getTokenContracts)

	//没有数据，更新数据
	if getTokenContracts == nil || len(getTokenContracts) == 0 {
		err = cli.UpdateSymbols()
		if err != nil {
			return nil, err
		}

		err = cli.db.Select(q.And(query...)).Find(&getTokenContracts)
		if err != nil {
			return nil, err
		}

	}
	return getTokenContracts, nil
}

//printTokenContractList 打印代币合约列表
func (cli *CLI) printTokenContractList(list []*openwsdk.TokenContract) {

	if list != nil && len(list) > 0 {
		tableInfo := make([][]interface{}, 0)

		for _, w := range list {
			tableInfo = append(tableInfo, []interface{}{
				w.ContractID, w.Symbol, w.Name, w.Token, w.Address, w.Protocol, w.Decimals,
			})
		}

		t := gotabulate.Create(tableInfo)
		// Set Headers
		t.SetHeaders([]string{"ContractID", "Symbol", "Name", "Token", "Address", "Protocol", "Decimals"})

		//打印信息
		fmt.Println(t.Render("simple"))
	} else {
		fmt.Println("No TokenContract. ")
	}
}

//GetTokenContractInfo 查询单个合约信息
func (cli *CLI) GetTokenContractInfo(contractID string) (*openwsdk.TokenContract, error) {

	getTokenContracts, err := cli.GetTokenContractList()
	if err != nil {
		return nil, err
	}

	for _, c := range getTokenContracts {
		if c.ContractID == contractID {
			return c, nil
		}
	}

	return nil, fmt.Errorf("can not find symbol info")
}

//SetSummaryInfo 设置账户的汇总设置
func (cli *CLI) SetSummaryInfo(obj *openwsdk.SummarySetting) error {
	//TODO:查询钱包和账户是否存在
	return cli.db.Save(obj)
}

//getLocalKeyByWallet
func (cli *CLI) getLocalKeyByWallet(wallet *openwsdk.Wallet, password string) (*hdkeystore.HDKey, error) {
	keystore := hdkeystore.NewHDKeystore(
		cli.config.keydir,
		hdkeystore.StandardScryptN,
		hdkeystore.StandardScryptP,
	)

	fileName := fmt.Sprintf("%s-%s.key", wallet.Alias, wallet.WalletID)

	key, err := keystore.GetKey(
		wallet.WalletID,
		fileName,
		password,
	)
	if err != nil {
		return nil, err
	}
	return key, nil
}
