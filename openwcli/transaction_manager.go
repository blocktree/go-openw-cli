package openwcli

import (
	"fmt"
	"github.com/blocktree/OpenWallet/common"
	"github.com/blocktree/OpenWallet/hdkeystore"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/OpenWallet/owtp"
	"github.com/blocktree/go-openw-sdk/openwsdk"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

//Transfer 转账交易
func (cli *CLI) Transfer(wallet *openwsdk.Wallet, account *openwsdk.Account, contractAddress, to, amount, password string) error {

	var (
		isContract bool
		retTx      []*openwsdk.Transaction
		retFailed  []*openwsdk.FailedRawTransaction
		retRawTx   *openwsdk.RawTransaction
		err        error
		contractID string
	)

	//获取种子文件
	key, err := cli.getLocalKeyByWallet(wallet, password)
	if err != nil {
		return err
	}

	if len(contractAddress) > 0 {
		isContract = true
		token, findErr := cli.GetTokenContractList("Symbol", account.Symbol, "Address", contractAddress)
		if findErr != nil {
			return findErr
		}
		contractID = token[0].ContractID
	}
	coin := openwsdk.Coin{
		Symbol:     account.Symbol,
		IsContract: isContract,
		ContractID: contractID,
	}

	//创建新交易单
	sid := uuid.New().String()
	api := cli.api
	err = api.CreateTrade(account.AccountID, sid, coin, amount, to, "", "", true,
		func(status uint64, msg string, rawTx *openwsdk.RawTransaction) {
			if status != owtp.StatusSuccess {
				err = fmt.Errorf(msg)
				return
			}
			retRawTx = rawTx
		})
	if err != nil {
		return err
	}

	//签名交易单
	err = openwsdk.SignRawTransaction(retRawTx, key)
	if err != nil {
		return err
	}

	//广播交易单
	err = api.SubmitTrade(retRawTx, true,
		func(status uint64, msg string, successTx []*openwsdk.Transaction, failedRawTxs []*openwsdk.FailedRawTransaction) {
			if status != owtp.StatusSuccess {
				err = fmt.Errorf(msg)
				return
			}

			retTx = successTx
			retFailed = failedRawTxs
		})
	if err != nil {
		return err
	}

	if len(retTx) > 0 {
		//打印交易单
		log.Info("send transaction successfully.")
		log.Info("transaction id:", retTx[0].Txid)
	} else if len(retFailed) > 0 {
		//打印交易单
		log.Error("send transaction failed. unexpected error:", retFailed[0].Reason)
	}

	return nil
}

//SummaryWallets 执行汇总流程
func (cli *CLI) SummaryTask() {

	log.Infof("[Summary Task Start]------%s", common.TimeFormat("2006-01-02 15:04:05"))

	//读取参与汇总的钱包
	for _, task := range cli.summaryTask.Wallets {

		if task.wallet == nil {
			w, err := cli.GetWalletByWalletIDOnLocal(task.WalletID)
			if err != nil {
				log.Errorf("Summary wallet[%s] unexpected error: %v", task.WalletID, err)
				continue
			}
			task.wallet = w
		}

		key, err := cli.getLocalKeyByWallet(task.wallet, task.Password)
		if err != nil {
			log.Errorf("Summary wallet[%s] unexpected error: %v", task.WalletID, err)
			continue
		}

		for _, accountTask := range task.Accounts {

			account, err := cli.GetAccountByAccountID(accountTask.AccountID)
			if err != nil {
				continue
			}

			//汇总账户主币
			cli.SummaryAccountMainCoin(account, key)

		}

	}

	log.Infof("[Summary Task End]------%s", common.TimeFormat("2006-01-02 15:04:05"))
}

//SummaryAccountMainCoin 汇总账户主币
func (cli *CLI) SummaryAccountMainCoin(account *openwsdk.Account, key *hdkeystore.HDKey) error {

	const (
		limit = 1000
	)

	var (
		err       error
		retTx     []*openwsdk.Transaction
		retFailed []*openwsdk.FailedRawTransaction
		retRawTxs []*openwsdk.RawTransaction
		sumSets   SummarySetting
	)

	//读取汇总信息
	err = cli.db.One("AccountID", account.AccountID, &sumSets)
	if err != nil {
		return err
	}

	if sumSets.SumAddress == "" {
		log.Errorf("Summary account[%s] summary address is empty!")
		return err
	}

	balance, _ := decimal.NewFromString(account.Balance)
	threshold, _ := decimal.NewFromString(sumSets.Threshold)
	//如果余额大于阀值，汇总的地址
	if balance.GreaterThan(threshold) {

		coin := openwsdk.Coin{
			Symbol: account.Symbol,
		}

		log.Infof("Summary account[%s] Current Balance = %v ", account.AccountID, balance)
		log.Infof("Summary account[%s] Summary Address = %v ", account.AccountID, sumSets.SumAddress)
		log.Infof("Summary account[%s] Start Create Summary Transaction", account.AccountID)

		//分页汇总交易
		for i := 0; i < int(account.AddressIndex)+1; i = i + limit {
			err = nil
			retRawTxs = nil
			retTx = nil
			retFailed = nil

			log.Infof("Create Summary Transaction in address range [%d...%d]", i, i+limit)

			err = cli.api.CreateSummaryTx(account.AccountID, sumSets.SumAddress, coin,
				"", sumSets.MinTransfer, sumSets.RetainedBalance,
				i, limit, sumSets.Confirms, true,
				func(status uint64, msg string, rawTxs []*openwsdk.RawTransaction) {
					retRawTxs = rawTxs
					if status != owtp.StatusSuccess {
						err = fmt.Errorf(msg)
					}
				})

			if err != nil {
				log.Errorf("CreateSummaryTransaction unexpected error: %v", err)
				continue
			}

			for _, rawTx := range retRawTxs {
				//签名交易
				err = openwsdk.SignRawTransaction(rawTx, key)
				if err != nil {
					log.Errorf("SignRawTransaction unexpected error: %v", err)
					continue
				}

				//	广播交易单
				err = cli.api.SubmitTrade(rawTx, true,
					func(status uint64, msg string, successTx []*openwsdk.Transaction, failedRawTxs []*openwsdk.FailedRawTransaction) {
						if status != owtp.StatusSuccess {
							err = fmt.Errorf(msg)
							return
						}

						retTx = successTx
						retFailed = failedRawTxs
					})
				if err != nil {
					log.Errorf("SubmitRawTransaction unexpected error: %v", err)
					continue
				}

				//打印汇总交易结果

				for _, tx := range retTx {
					log.Infof("[Success] txid:", tx.Txid)
				}

				for _, tx := range retFailed {
					log.Errorf("[Failed] reason:", tx.Reason)
				}
			}

		}
	} else {
		log.Infof("Summary account[%s] Current Balance: %v，below threshold: %v", account.AccountID, balance, threshold)
	}
	return nil
}

//SummaryAccountTokenContracts 汇总账户代币合约
func (cli *CLI) SummaryAccountTokenContracts(account *openwsdk.Account, key *hdkeystore.HDKey) error {

	const (
		limit = 1000
	)

	var (
		err       error
		retTx     []*openwsdk.Transaction
		retFailed []*openwsdk.FailedRawTransaction
		retRawTxs []*openwsdk.RawTransaction
		sumSets   SummarySetting
	)

	//读取汇总信息
	err = cli.db.One("AccountID", account.AccountID, &sumSets)
	if err != nil {
		return err
	}

	if sumSets.SumAddress == "" {
		log.Errorf("Summary account[%s] summary address is empty!")
		return err
	}

	balance, _ := decimal.NewFromString(account.Balance)
	threshold, _ := decimal.NewFromString(sumSets.Threshold)
	//如果余额大于阀值，汇总的地址
	if balance.GreaterThan(threshold) {

		coin := openwsdk.Coin{
			Symbol: account.Symbol,
		}

		log.Infof("Summary account[%s] Current Balance = %v ", account.AccountID, balance)
		log.Infof("Summary account[%s] Summary Address = %v ", account.AccountID, sumSets.SumAddress)
		log.Infof("Summary account[%s] Start Create Summary Transaction", account.AccountID)

		//分页汇总交易
		for i := 0; i < int(account.AddressIndex)+1; i = i + limit {
			err = nil
			retRawTxs = nil
			retTx = nil
			retFailed = nil

			log.Infof("Create Summary Transaction in address range [%d...%d]", i, i+limit)

			err = cli.api.CreateSummaryTx(account.AccountID, sumSets.SumAddress, coin,
				"", sumSets.MinTransfer, sumSets.RetainedBalance,
				i, limit, sumSets.Confirms, true,
				func(status uint64, msg string, rawTxs []*openwsdk.RawTransaction) {
					retRawTxs = rawTxs
					if status != owtp.StatusSuccess {
						err = fmt.Errorf(msg)
					}
				})

			if err != nil {
				log.Errorf("CreateSummaryTransaction unexpected error: %v", err)
				continue
			}

			for _, rawTx := range retRawTxs {
				//签名交易
				err = openwsdk.SignRawTransaction(rawTx, key)
				if err != nil {
					log.Errorf("SignRawTransaction unexpected error: %v", err)
					continue
				}

				//	广播交易单
				err = cli.api.SubmitTrade(rawTx, true,
					func(status uint64, msg string, successTx []*openwsdk.Transaction, failedRawTxs []*openwsdk.FailedRawTransaction) {
						if status != owtp.StatusSuccess {
							err = fmt.Errorf(msg)
							return
						}

						retTx = successTx
						retFailed = failedRawTxs
					})
				if err != nil {
					log.Errorf("SubmitRawTransaction unexpected error: %v", err)
					continue
				}

				//打印汇总交易结果

				for _, tx := range retTx {
					log.Infof("[Success] txid:", tx.Txid)
				}

				for _, tx := range retFailed {
					log.Errorf("[Failed] reason:", tx.Reason)
				}
			}

		}
	} else {
		log.Infof("Summary account[%s] Current Balance: %v，below threshold: %v", account.AccountID, balance, threshold)
	}
	return nil
}