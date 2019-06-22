package openwcli

import (
	"fmt"
	"github.com/blocktree/go-openw-sdk/openwsdk"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/blocktree/openwallet/owtp"
)

//GetTokenBalance 获取代币余额
func (cli *CLI) GetTokenBalance(account *openwsdk.Account, contractID string) string {
	getBalance := "0"
	cli.api.GetTokenBalanceByAccount(account.AccountID, contractID, true,
		func(status uint64, msg string, balance *openwsdk.TokenBalance) {
			if status == owtp.StatusSuccess {
				getBalance = balance.Balance.Balance
			}
		})
	return getBalance
}

//GetTokenBalanceByContractAddress 通过代币合约的地址获取代币余额
func (cli *CLI) GetTokenBalanceByContractAddress(account *openwsdk.Account, address string) (*openwsdk.TokenBalance, error) {
	var (
		getBalance *openwsdk.TokenBalance
		callErr error
	)

	token, findErr := cli.GetTokenContractList("Symbol", account.Symbol, "Address", address)
	if findErr != nil {
		return nil, findErr
	}
	contractID := token[0].ContractID
	err := cli.api.GetTokenBalanceByAccount(account.AccountID, contractID, true,
		func(status uint64, msg string, balance *openwsdk.TokenBalance) {
			if status == owtp.StatusSuccess {
				getBalance = balance
			} else {
				callErr = openwallet.Errorf(status, msg)
			}
		})
	if err != nil {
		return nil, err
	}
	if callErr != nil {
		return nil, err
	}
	return getBalance, nil
}

//Transfer 转账交易
func (cli *CLI) Transfer(wallet *openwsdk.Wallet, account *openwsdk.Account, contractAddress, to, amount, sid, feeRate, memo, password string) ([]*openwsdk.Transaction, []*openwsdk.FailedRawTransaction, error) {

	var (
		isContract  bool
		retTx       []*openwsdk.Transaction
		retFailed   []*openwsdk.FailedRawTransaction
		retRawTx    *openwsdk.RawTransaction
		err         error
		createErr   error
		contractID  string
		tokenSymbol string
	)

	//获取种子文件
	key, err := cli.getLocalKeyByWallet(wallet, password)
	if err != nil {
		return nil, nil, err
	}

	if len(contractAddress) > 0 {
		isContract = true
		token, findErr := cli.GetTokenContractList("Symbol", account.Symbol, "Address", contractAddress)
		if findErr != nil {
			return nil, nil, findErr
		}
		contractID = token[0].ContractID
		tokenSymbol = token[0].Token
	}
	coin := openwsdk.Coin{
		Symbol:     account.Symbol,
		IsContract: isContract,
		ContractID: contractID,
	}

	api := cli.api
	err = api.CreateTrade(account.AccountID, sid, coin, amount, to, feeRate, memo, true,
		func(status uint64, msg string, rawTx *openwsdk.RawTransaction) {
			if status != owtp.StatusSuccess {
				createErr = fmt.Errorf(msg)
				return
			}
			retRawTx = rawTx
		})
	if err != nil {
		return nil, nil, err
	}
	if createErr != nil {
		return nil, nil, createErr
	}

	//:打印交易单明细
	log.Infof("-----------------------------------------------")
	log.Infof("[%s %s Transfer]", account.Symbol, tokenSymbol)
	log.Infof("From Account: %s", account.AccountID)
	log.Infof("To Address: %s", to)
	log.Infof("Send Amount: %s", amount)
	log.Infof("Fees: %v", retRawTx.Fees)
	log.Infof("FeeRate: %v", retRawTx.FeeRate)
	log.Infof("Memo: %v", memo)
	log.Infof("-----------------------------------------------")

	//签名交易单
	err = openwsdk.SignRawTransaction(retRawTx, key)
	if err != nil {
		return nil, nil, err
	}

	//广播交易单
	err = api.SubmitTrade([]*openwsdk.RawTransaction{retRawTx}, true,
		func(status uint64, msg string, successTx []*openwsdk.Transaction, failedRawTxs []*openwsdk.FailedRawTransaction) {
			if status != owtp.StatusSuccess {
				createErr = fmt.Errorf(msg)
				return
			}

			retTx = successTx
			retFailed = failedRawTxs
		})
	if err != nil {
		return nil, nil, err
	}
	if createErr != nil {
		return nil, nil, createErr
	}

	if len(retTx) > 0 {
		//打印交易单
		log.Info("send transaction successfully.")
		log.Info("transaction id:", retTx[0].Txid)
	} else if len(retFailed) > 0 {
		//打印交易单
		log.Error("send transaction failed. unexpected error:", retFailed[0].Reason)
	}

	return retTx, retFailed, nil
}


