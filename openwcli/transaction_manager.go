package openwcli

import (
	"fmt"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/OpenWallet/owtp"
	"github.com/blocktree/go-openw-sdk/openwsdk"
	"github.com/google/uuid"
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

