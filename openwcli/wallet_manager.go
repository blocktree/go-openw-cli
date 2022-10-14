package openwcli

import (
	"encoding/hex"
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"path/filepath"
	"strings"
	"time"

	"github.com/asdine/storm/q"
	"github.com/blocktree/go-openw-sdk/v2/openwsdk"
	"github.com/blocktree/openwallet/v2/common"
	"github.com/blocktree/openwallet/v2/common/file"
	"github.com/blocktree/openwallet/v2/hdkeystore"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/blocktree/openwallet/v2/owtp"
	"github.com/bndr/gotabulate"
)

// CreateWalletOnServer
func (cli *CLI) CreateWalletOnServer(name, password string) (*openwsdk.Wallet, error) {

	var (
		key       *hdkeystore.HDKey
		retWallet *openwsdk.Wallet
		retErr    error
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

	walletParam := &openwsdk.Wallet{
		Alias:    name,
		WalletID: key.KeyID,
	}

	//登记钱包的openw-server
	err = cli.api.CreateWallet(walletParam, true,
		func(status uint64, msg string, wallet *openwsdk.Wallet) {
			if status == owtp.StatusSuccess {
				log.Info("Wallet create successfully, key path:", filePath)
				retWallet = wallet
			} else {
				log.Error("create wallet on server failed, unexpected error:", msg)
				retErr = openwallet.Errorf(status, msg)
				//创建失败，删除key文件
				file.Delete(filePath)
			}
		})
	if err != nil {
		return nil, err
	}

	return retWallet, retErr
}

// GetWalletsByKeyDir 通过给定的文件路径加载keystore文件得到钱包列表
func (cli *CLI) GetWalletsOnServer() ([]*openwsdk.Wallet, error) {
	localWallets, err := openwallet.GetWalletsByKeyDir(cli.config.keydir)
	if err != nil {
		return nil, err
	}
	serverWallets := make([]*openwsdk.Wallet, 0)

	for _, w := range localWallets {
		callErr := cli.api.FindWalletByWalletID(w.WalletID, true,
			func(status uint64, msg string, wallet *openwsdk.Wallet) {
				if status == owtp.StatusSuccess && wallet != nil {
					serverWallets = append(serverWallets, wallet)
				}
			})
		if callErr != nil {
			return nil, callErr
		}
	}

	return serverWallets, nil
}

// GetWalletByWalletID 查找本地且线上有的钱包对象
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

// GetWalletByWalletIDOnLocal 查找本地种子目录的钱包对象
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

// printWalletList 打印钱包列表
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

// CreateAccountOnServer
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
				log.Infof("create [%s] account successfully", selectedSymbol.Symbol)
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

// GetAccountOnServerByAccountID 从服务器获取账户
func (cli *CLI) GetAccountByAccountID(symbol, accountID string) (*openwsdk.Account, error) {

	var (
		getAccount *openwsdk.Account
		err        error
		retErr     error
	)

	err = cli.api.FindAccountByAccountID(symbol, accountID, 0, true,
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

// GetAccountsOnServer 从服务器获取账户列表
func (cli *CLI) GetAccountsOnServer(walletID string) ([]*openwsdk.Account, error) {

	var (
		list   = make([]*openwsdk.Account, 0)
		err    error
		retErr error
		lastID = int64(0)
		limit  = int64(200)
	)

	err = cli.api.FindAccountByWalletID("", walletID, lastID, limit, true,
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

// printAccountList 打印账户列表
func (cli *CLI) printAccountList(list []*openwsdk.Account) {

	_, err := cli.getDB()
	if err != nil {
		return
	}
	defer cli.closeDB()

	if list != nil && len(list) > 0 {
		tableInfo := make([][]interface{}, 0)

		for i, w := range list {

			//读取汇总信息
			sumTips := ""
			var sum openwsdk.SummarySetting
			err := cli.db.One("AccountID", w.AccountID, &sum)
			if err != nil {
				sumTips = "X"
			} else {
				sumTips = "√"
			}
			balanceStr := "0"
			//查询账户余额
			cli.api.GetBalanceByAccount(w.Symbol, w.AccountID, "",
				true, func(status uint64, msg string, balance *openwsdk.BalanceResult) {
					if status == owtp.StatusSuccess {
						balanceStr = balance.Balance
					} else {
						balanceStr = "N/A"
					}
				})

			tableInfo = append(tableInfo, []interface{}{
				i, w.Alias, w.AccountID, w.Symbol, balanceStr, w.AddressIndex + 1, sumTips,
			})
		}

		t := gotabulate.Create(tableInfo)
		// Set Headers
		t.SetHeaders([]string{"No.", "Name", "AccountID", "Symbol", "Balance", "Addresses",
			"Setup summary info"})

		//打印信息
		fmt.Println(t.Render("simple"))
	} else {
		fmt.Println("No account was created locally. ")
	}
}

// printAccountList 打印账户列表
func (cli *CLI) printAccountSummaryInfo() {

	_, err := cli.getDB()
	if err != nil {
		return
	}
	defer cli.closeDB()

	//读取汇总信息
	var sum []*openwsdk.SummarySetting
	err = cli.db.All(&sum)
	if err != nil || len(sum) == 0 {
		fmt.Println("No account setup summary info. ")
		return
	}

	tableInfo := make([][]interface{}, 0)

	for _, s := range sum {
		tableInfo = append(tableInfo, []interface{}{
			s.AccountID, s.SumAddress, s.Threshold, s.MinTransfer, s.RetainedBalance, s.Confirms,
		})
	}

	t := gotabulate.Create(tableInfo)
	// Set Headers
	t.SetHeaders([]string{"AccountID", "Summary Address", "Summary Threshold", "Min Transfer", "Retained Balance", "Confirms"})

	//打印信息
	fmt.Println(t.Render("simple"))
}

// CreateAddressOnServer
func (cli *CLI) CreateAddressOnServer(walletID, accountID, symbol string, count uint64) error {

	var (
		retErr error
	)

	if len(accountID) == 0 {
		return fmt.Errorf("accountID is empty. ")
	}

	if len(walletID) == 0 {
		return fmt.Errorf("walleID is empty. ")
	}

	if count == 0 {
		return fmt.Errorf("create address count can not 0. ")
	}

	err := cli.api.CreateAddress(symbol, walletID, accountID, count, true,
		func(status uint64, msg string, addresses []*openwsdk.Address) {
			if status == owtp.StatusSuccess {
				log.Infof("create [%d] addresses successfully", len(addresses))
				//:保存到本地数据库，导出到文件夹
				timestamp := time.Now()
				filename := "[" + accountID + "]-" + common.TimeFormat("20060102150405", timestamp) + ".txt"
				filePath := filepath.Join(cli.config.exportaddressdir, filename)
				if flag := cli.exportAddressToFile(addresses, filePath); flag {
					log.Infof("addresses has been exported into: %s", filePath)
				} else {
					log.Infof("addresses export failed")
				}
			} else {
				log.Error("create account on server failed, unexpected error:", msg)
				retErr = openwallet.Errorf(status, msg)
			}
		})

	if err != nil {
		return err
	}

	return retErr
}

// exportAddressToFile 导出地址到文件中
func (cli *CLI) exportAddressToFile(addresses []*openwsdk.Address, filePath string) bool {

	var (
		content string
	)

	for _, a := range addresses {
		content = content + a.Address + "\n"
	}

	return file.WriteFile(filePath, []byte(content), true)
}

// SearchAddressOnServer
func (cli *CLI) SearchAddressOnServer(symbol, address string) (*openwsdk.Address, error) {

	var (
		retErr error
	)

	if len(address) == 0 {
		return nil, fmt.Errorf("address is empty. ")
	}

	var addr *openwsdk.Address

	err := cli.api.FindAddressByAddress(symbol, address, true,
		func(status uint64, msg string, address *openwsdk.Address) {
			if status == owtp.StatusSuccess {
				addr = address
			} else {
				log.Error("search address on server failed, unexpected error:", msg)
				retErr = openwallet.Errorf(status, msg)
			}
		})
	if err != nil {
		return nil, err
	}

	return addr, retErr
}

// GetAddressesOnServer
func (cli *CLI) GetAddressesOnServer(walletID, accountID, symbol string, lastId, limit int64) ([]*openwsdk.Address, error) {

	var (
		retErr error
	)

	list := make([]*openwsdk.Address, 0)

	if len(accountID) == 0 {
		return nil, fmt.Errorf("accountID is empty. ")
	}

	if len(walletID) == 0 {
		return nil, fmt.Errorf("walleID is empty. ")
	}

	err := cli.api.FindAddressByAccountID(symbol, accountID, lastId, limit, true,
		func(status uint64, msg string, addresses []*openwsdk.Address) {
			if status == owtp.StatusSuccess {
				list = addresses
			} else {
				log.Error("get address on server failed, unexpected error:", msg)
				retErr = openwallet.Errorf(status, msg)
			}
		})
	if err != nil {
		return nil, err
	}

	return list, retErr
}

// printAddressList 打印地址列表
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

		for _, a := range list {

			balanceStr := "0"
			cli.api.GetBalanceByAddress(a.Symbol, a.Address, "",
				true, func(status uint64, msg string, balance *openwsdk.BalanceResult) {
					if status == owtp.StatusSuccess {
						balanceStr = balance.Balance
					} else {
						balanceStr = "N/A"
					}
				})

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
				a.Id, a.Address, a.WalletID, a.AccountID, a.Symbol, balanceStr, a.PublicKey, privatekey,
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

// UpdateSymbols 更新主链
func (cli *CLI) UpdateSymbols() error {

	const (
		limit = 500
	)

	_, err := cli.getDB()
	if err != nil {
		return err
	}
	defer cli.closeDB()

	var getSymbols []*openwsdk.Symbol
	err = cli.api.GetSymbolList("", 0, limit, 0, true,
		func(status uint64, msg string, total int, symbols []*openwsdk.Symbol) {
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
			err = cli.api.GetContracts(s.Symbol, "", i, limit, true,
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

// UpdateSymbols 更新主链
func (cli *CLI) UpdateTokenContracts(symbol string) error {
	var getTokenContract []*openwsdk.TokenContract
	err := cli.api.GetContracts(symbol, "", 0, 5000, true,
		func(status uint64, msg string, tokenContract []*openwsdk.TokenContract) {
			getTokenContract = tokenContract
		})
	if err != nil {
		return err
	}

	_, err = cli.getDB()
	if err != nil {
		return err
	}
	defer cli.closeDB()

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

// GetLocalSymbolList 查询本地保存主链
func (cli *CLI) GetSymbolList() ([]*openwsdk.Symbol, error) {

	_, err := cli.getDB()
	if err != nil {
		return nil, err
	}

	var getSymbols []*openwsdk.Symbol
	err = cli.db.All(&getSymbols)
	cli.closeDB()

	//没有数据，更新数据
	if getSymbols == nil || len(getSymbols) == 0 {

		err = cli.UpdateSymbols()
		if err != nil {
			return nil, err
		}

		_, err = cli.getDB()
		if err != nil {
			return nil, err
		}

		err = cli.db.All(&getSymbols)
		cli.closeDB()

	}
	return getSymbols, nil
}

// printSymbolList 打印主链列表
func (cli *CLI) printSymbolList(list []*openwsdk.Symbol) {

	if list != nil && len(list) > 0 {
		tableInfo := make([][]interface{}, 0)

		for _, w := range list {
			tableInfo = append(tableInfo, []interface{}{
				w.Name, w.Symbol, w.Curve, w.Decimals,
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

// GetLocalSymbolInfo 查询本地主链信息
func (cli *CLI) GetSymbolInfo(symbol string) (*openwsdk.Symbol, error) {

	getSymbols, err := cli.GetSymbolList()
	if err != nil {
		return nil, err
	}

	for _, s := range getSymbols {
		if s.Symbol == strings.ToUpper(symbol) {
			return s, nil
		}
	}

	return nil, fmt.Errorf("can not find symbol info")
}

// GetContractList 查询本地保存代币合约信息
func (cli *CLI) GetTokenContractList(cols ...interface{}) ([]*openwsdk.TokenContract, error) {

	var (
		query             = make([]q.Matcher, 0)
		getTokenContracts []*openwsdk.TokenContract
		err               error
	)

	if len(cols)%2 != 0 {
		return nil, fmt.Errorf("condition param is not pair")
	}

	for i := 0; i < len(cols); i = i + 2 {
		field := common.NewString(cols[i])
		val := cols[i+1]
		query = append(query, q.Eq(field.String(), val))
	}

	getTokenContractListFunc := func(queryMatcher []q.Matcher) []*openwsdk.TokenContract {
		var tokenContracts []*openwsdk.TokenContract
		_, err := cli.getDB()
		if err != nil {
			return nil
		}
		defer cli.closeDB()

		err = cli.db.Select(q.And(queryMatcher...)).Find(&tokenContracts)
		return tokenContracts
	}

	getTokenContracts = getTokenContractListFunc(query)

	//没有数据，更新数据
	if getTokenContracts == nil || len(getTokenContracts) == 0 {
		err = cli.UpdateSymbols()
		if err != nil {
			return nil, err
		}

		getTokenContracts = getTokenContractListFunc(query)

	}
	return getTokenContracts, nil
}

// printTokenContractList 打印代币合约列表
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

// GetTokenContractInfo 查询单个合约信息
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

// SetSummaryInfo 设置账户的汇总设置
func (cli *CLI) SetSummaryInfo(obj *openwsdk.SummarySetting) error {

	//检查账户是否存在
	account, err := cli.GetAccountByAccountID(obj.Symbol, obj.AccountID)
	if err != nil {
		return err
	}

	//把汇总地址添加到信任名单
	trustAddr := openwsdk.NewTrustAddress(
		obj.SumAddress,
		account.Symbol,
		"summary address")
	err = cli.AddTrustAddress(trustAddr)
	if err != nil {
		return err
	}

	_, err = cli.getDB()
	if err != nil {
		return err
	}
	defer cli.closeDB()

	return cli.db.Save(obj)
}

// getLocalKeyByWallet
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

// GetAllTokenContractBalance 查询账户合约余额
func (cli *CLI) GetAllTokenContractBalance(accountID string, symbol string) ([]*openwsdk.TokenBalance, error) {

	var (
		getErr      error
		getBalances []*openwsdk.TokenBalance
	)
	err := cli.api.GetAllTokenBalanceByAccount(accountID, symbol, true,
		func(status uint64, msg string, balance []*openwsdk.TokenBalance) {
			if status == owtp.StatusSuccess {
				getBalances = balance
			} else {
				getErr = fmt.Errorf(msg)
			}
		})
	if err != nil {
		return nil, err
	}

	if getErr != nil {
		return nil, getErr
	}

	return getBalances, nil
}

// GetAllTokenContractBalanceByAddress 查询地址合约余额
func (cli *CLI) GetAllTokenContractBalanceByAddress(accountID, address, symbol string) ([]*openwsdk.TokenBalance, error) {

	var (
		getErr      error
		getBalances []*openwsdk.TokenBalance
	)
	err := cli.api.GetAllTokenBalanceByAddress(accountID, address, symbol, true,
		func(status uint64, msg string, balance []*openwsdk.TokenBalance) {
			if status == owtp.StatusSuccess {
				getBalances = balance
			} else {
				getErr = fmt.Errorf(msg)
			}
		})
	if err != nil {
		return nil, err
	}

	if getErr != nil {
		return nil, getErr
	}

	return getBalances, nil
}

func findTokenContractByID(tokenList []*openwsdk.TokenContract, contractID string) *openwsdk.TokenContract {
	for _, c := range tokenList {
		if c.ContractID == contractID {
			return c
		}
	}
	return nil
}

// printTokenContractBalanceList 打印账户代币合约余额列表
func (cli *CLI) printTokenContractBalanceList(list []*openwsdk.TokenBalance, symbol string) {

	if list != nil && len(list) > 0 {
		tableInfo := make([][]interface{}, 0)

		getTokenContracts, err := cli.GetTokenContractList("Symbol", strings.ToUpper(symbol))
		if err != nil {
			fmt.Println("Please execute command 'updateinfo' first. ")
			return
		}

		for _, w := range list {

			token := findTokenContractByID(getTokenContracts, w.ContractID)
			if token == nil {
				continue
			}

			tableInfo = append(tableInfo, []interface{}{
				w.ContractID, token.Symbol, token.Name, w.Token, token.Address, token.Protocol, w.Balance.Balance,
			})
		}

		if len(tableInfo) == 0 {
			fmt.Println("Please execute command 'updateinfo' first. ")
			return
		}

		t := gotabulate.Create(tableInfo)
		// Set Headers
		t.SetHeaders([]string{"ContractID", "Symbol", "Name", "Token", "Address", "Protocol", "Balance"})

		//打印信息
		fmt.Println(t.Render("simple"))
	} else {
		fmt.Println("No Token Contract Balance.")
	}
}

// AddTrustAddress 添加白名单地址
func (cli *CLI) AddTrustAddress(trustAddress *openwsdk.TrustAddress) error {

	//检查symbol是否存在
	s, err := cli.GetSymbolInfo(trustAddress.Symbol)
	if err != nil {
		return err
	}

	trustAddress.Symbol = s.Symbol

	_, err = cli.getDB()
	if err != nil {
		return err
	}
	defer cli.closeDB()

	err = cli.db.Save(trustAddress)
	if err != nil {
		return err
	}
	return nil
}

// ListTrustAddress 白名单地址列表
func (cli *CLI) ListTrustAddress(symbol string) ([]*openwsdk.TrustAddress, error) {

	var (
		list []*openwsdk.TrustAddress
		err  error
	)

	_, err = cli.getDB()
	if err != nil {
		return nil, err
	}
	defer cli.closeDB()

	if symbol == "" {
		err = cli.db.All(&list)
	} else {
		err = cli.db.Find("Symbol", symbol, &list)
	}

	if err != nil {
		return nil, nil
	}
	return list, nil
}

// printListTrustAddress 白名单地址列表
func (cli *CLI) printListTrustAddress(addrs []*openwsdk.TrustAddress) {

	if len(addrs) == 0 {
		fmt.Println("No Trust Address info. ")
		return
	}

	tableInfo := make([][]interface{}, 0)

	for _, s := range addrs {
		t := time.Unix(s.CreateTime, 0)
		strTime := common.TimeFormat("2006-01-02 15:04:05", t)
		tableInfo = append(tableInfo, []interface{}{
			s.Address, strings.ToUpper(s.Symbol), s.Memo, strTime,
		})
	}

	t := gotabulate.Create(tableInfo)
	// Set Headers
	t.SetHeaders([]string{"Address", "Symbol", "Memo", "CreateTime"})

	//打印信息
	fmt.Println(t.Render("simple"))
}

// importSummaryAddressToTrustAddress 导入汇总地址到信任地址列表
func (cli *CLI) importSummaryAddressToTrustAddress() error {

	_, err := cli.getDB()
	if err != nil {
		return err
	}
	defer cli.closeDB()

	//读取汇总信息
	var sum []*openwsdk.SummarySetting
	cli.db.All(&sum)
	if len(sum) > 0 {
		for _, s := range sum {

			//检查账户是否存在
			account, err := cli.GetAccountByAccountID(s.Symbol, s.AccountID)
			if err != nil {
				return err
			}

			//把汇总地址添加到信任名单
			trustAddr := openwsdk.NewTrustAddress(
				s.SumAddress,
				account.Symbol,
				"summary address")
			err = cli.AddTrustAddress(trustAddr)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

// EnableTrustAddress
func (cli *CLI) EnableTrustAddress() error {

	_, err := cli.getDB()
	if err != nil {
		return err
	}
	defer cli.closeDB()

	var inited bool
	cli.db.Get(CLIBucket, InitTrustAddress, &inited)

	//第一次初始化都把现有的汇总地址导入到信任地址名单
	if !inited {

		err := cli.importSummaryAddressToTrustAddress()
		if err != nil {
			return err
		}

		err = cli.db.Set(CLIBucket, InitTrustAddress, true)
		if err != nil {
			return fmt.Errorf("Init Trust Address, unexpected error: %v ", err)
		}
	}

	err = cli.db.Set(CLIBucket, EnableTrustAddress, true)
	if err != nil {
		return fmt.Errorf("Enable Trust Address, unexpected error: %v ", err)
	}
	return nil
}

// DisableTrustAddress
func (cli *CLI) DisableTrustAddress() error {

	_, err := cli.getDB()
	if err != nil {
		return err
	}
	defer cli.closeDB()

	err = cli.db.Set(CLIBucket, EnableTrustAddress, false)
	if err != nil {
		return fmt.Errorf("Enable Trust Address, unexpected error: %v ", err)
	}
	return nil
}

// TrustAddressStatus
func (cli *CLI) TrustAddressStatus() bool {

	_, err := cli.getDB()
	if err != nil {
		log.Errorf("cli database open failed")
		return false
	}
	defer cli.closeDB()

	var status bool
	cli.db.Get(CLIBucket, EnableTrustAddress, &status)
	return status
}

// printTrustAddressStatus
func (cli *CLI) printTrustAddressStatus() {

	if cli.TrustAddressStatus() {
		fmt.Printf("######## Trust address is enabled. ######## \n")
	} else {
		fmt.Printf("######## Trust address is disabled. ######## \n")
	}
}

// IsTrustAddress
func (cli *CLI) IsTrustAddress(address, symbol string) bool {
	var (
		list []*openwsdk.TrustAddress
		err  error
	)

	if cli.TrustAddressStatus() {

		_, err = cli.getDB()
		if err != nil {
			fmt.Errorf("cli database open failed")
			return false
		}
		defer cli.closeDB()

		err = cli.db.Select(
			q.And(
				q.Eq("Address", address),
				q.Eq("Symbol", symbol),
			)).Find(&list)
		if err != nil {
			return false
		}
	}
	return true
}

// SignHash 哈希消息签名
func (cli *CLI) SignHash(address *openwsdk.Address, message, password string, appendV bool) (string, error) {

	wallet, err := cli.GetWalletByWalletID(address.WalletID)
	if err != nil {
		return "", err
	}

	//获取种子文件
	key, err := cli.getLocalKeyByWallet(wallet, password)
	if err != nil {
		return "", err
	}

	symbolInfo, err := cli.GetSymbolInfo(address.Symbol)
	if err != nil {
		return "", err
	}

	childKey, err := key.DerivedKeyWithPath(address.HdPath, uint32(symbolInfo.Curve))
	if err != nil {
		return "", err
	}

	keyBytes, err := childKey.GetPrivateKeyBytes()
	if err != nil {
		return "", err
	}

	hash, err := hex.DecodeString(strings.TrimPrefix(message, "0x"))
	if err != nil {
		return "", err
	}

	signature, v, sigErr := owcrypt.Signature(keyBytes, nil, hash, uint32(symbolInfo.Curve))
	if sigErr != owcrypt.SUCCESS {
		return "", fmt.Errorf("sign hash message failed")
	}

	if appendV {
		signature = append(signature, v)
	}

	return hex.EncodeToString(signature), nil
}
