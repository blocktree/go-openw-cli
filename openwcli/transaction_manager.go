package openwcli

import (
	"encoding/json"
	"github.com/blocktree/go-openw-sdk/v2/openwsdk"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/blocktree/openwallet/v2/owtp"
)

// GetTokenBalance 获取代币余额
func (cli *CLI) GetTokenBalance(account *openwsdk.Account, contractID string) string {
	getBalance := "0"
	cli.api.GetBalanceByAccount(account.Symbol, account.AccountID, contractID, true,
		func(status uint64, msg string, balance *openwsdk.BalanceResult) {
			if status == owtp.StatusSuccess {
				getBalance = balance.Balance
			}
		})
	return getBalance
}

// GetTokenBalanceByContractAddress 通过代币合约的地址获取代币余额
func (cli *CLI) GetTokenBalanceByContractAddress(account *openwsdk.Account, symbol, address string) (*openwsdk.BalanceResult, error) {
	var (
		getBalance *openwsdk.BalanceResult
		callErr    error
	)

	token, findErr := cli.GetTokenContractList("Symbol", symbol, "Address", address)
	if findErr != nil {
		return nil, findErr
	}
	contractID := token[0].ContractID
	err := cli.api.GetBalanceByAccount(symbol, account.AccountID, contractID, true,
		func(status uint64, msg string, balance *openwsdk.BalanceResult) {
			if status == owtp.StatusSuccess {
				getBalance = balance
				getBalance.ContractToken = token[0].Token
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

// Transfer 转账交易
func (cli *CLI) Transfer(wallet *openwsdk.Wallet, account *openwsdk.Account, symbol, contractAddress, to, amount, sid, feeRate, memo, password string) ([]*openwsdk.Transaction, []*openwsdk.FailedRawTransaction, *openwallet.Error) {
	return cli.TransferExt(wallet, account, symbol, contractAddress, to, amount, sid, feeRate, memo, "", password)
}

// TransferExt 转账交易 + 扩展参数
func (cli *CLI) TransferExt(wallet *openwsdk.Wallet, account *openwsdk.Account, symbol, contractAddress, to, amount, sid, feeRate, memo, extParam, password string) ([]*openwsdk.Transaction, []*openwsdk.FailedRawTransaction, *openwallet.Error) {

	var (
		isContract  bool
		retTx       []*openwsdk.Transaction
		retFailed   []*openwsdk.FailedRawTransaction
		retRawTx    *openwsdk.RawTransaction
		err         error
		createErr   *openwallet.Error
		contractID  string
		tokenSymbol string
	)

	//:检查目标地址是否信任名单
	if !cli.IsTrustAddress(to, symbol) {
		return nil, nil, openwallet.Errorf(openwallet.ErrUnknownException, "%s is not in trust address list", to)
	}

	if len(password) == 0 {
		return nil, nil, openwallet.Errorf(openwallet.ErrCreateRawTransactionFailed, "unlock wallet password is empty. ")
	}

	//获取种子文件
	key, err := cli.getLocalKeyByWallet(wallet, password)
	if err != nil {
		return nil, nil, openwallet.Errorf(openwallet.ErrCreateRawTransactionFailed, err.Error())
	}

	if len(contractAddress) > 0 {
		isContract = true
		token, findErr := cli.GetTokenContractList("Symbol", symbol, "Address", contractAddress)
		if findErr != nil {
			return nil, nil, openwallet.ConvertError(findErr)
		}
		if len(token) == 0 {
			return nil, nil, openwallet.Errorf(openwallet.ErrSystemException, "can not find contract address")
		}
		contractID = token[0].ContractID
		tokenSymbol = token[0].Token
	}
	coin := openwsdk.Coin{
		Symbol:     symbol,
		IsContract: isContract,
		ContractID: contractID,
	}

	api := cli.api
	err = api.CreateTrade(account.AccountID, sid, coin, map[string]string{to: amount}, feeRate, memo, extParam, true,
		func(status uint64, msg string, rawTx *openwsdk.RawTransaction) {
			if status != owtp.StatusSuccess {
				createErr = openwallet.Errorf(status, msg)
				return
			}
			retRawTx = rawTx
		})
	if err != nil {
		return nil, nil, openwallet.ConvertError(err)
	}
	if createErr != nil {
		return nil, nil, createErr
	}

	//:打印交易单明细
	log.Infof("-----------------------------------------------")
	log.Infof("[%s %s Transfer]", symbol, tokenSymbol)
	log.Infof("SID: %s", retRawTx.Sid)
	log.Infof("From Account: %s", account.AccountID)
	log.Infof("To Address: %s", to)
	log.Infof("Send Amount: %s", amount)
	log.Infof("Fees: %v", retRawTx.Fees)
	log.Infof("FeeRate: %v", retRawTx.FeeRate)
	log.Infof("Memo: %v", memo)
	log.Infof("-----------------------------------------------")

	//签名交易单
	signatures, sigErr := cli.txSigner(retRawTx.Signatures, key)
	if sigErr != nil {
		return nil, nil, openwallet.Errorf(openwallet.ErrSignRawTransactionFailed, sigErr.Error())
	}
	retRawTx.Signatures = signatures

	//广播交易单
	err = api.SubmitTrade([]*openwsdk.RawTransaction{retRawTx}, true,
		func(status uint64, msg string, successTx []*openwsdk.Transaction, failedRawTxs []*openwsdk.FailedRawTransaction) {
			if status != owtp.StatusSuccess {
				createErr = openwallet.Errorf(status, msg)
				return
			}

			retTx = successTx
			retFailed = failedRawTxs
		})
	if err != nil {
		return nil, nil, openwallet.ConvertError(err)
	}
	if createErr != nil {
		return nil, nil, createErr
	}

	if len(retTx) > 0 {
		//打印交易单
		log.Info("send transaction successfully.")
		log.Info("transaction id:", retTx[0].TxID)
	} else if len(retFailed) > 0 {
		//打印交易单
		log.Errorf("send transaction failed.")
		tx := retFailed[0]
		log.Warningf("[Failed] reason: %s", tx.Reason)
		if tx.RawTx != nil {
			log.Warningf("[Failed] rawHex: %s", tx.RawTx.RawHex)
			for accountID, signatures := range tx.RawTx.Signatures {
				log.Warningf("[Failed] signature accountID: %s", accountID)
				for _, keySignature := range signatures {
					signaturesJSON, jsonErr := json.Marshal(keySignature)
					if jsonErr == nil {
						log.Warningf("[Failed] keySignature: %s", string(signaturesJSON))
					}
				}
			}
		}

		return retTx, retFailed, openwallet.Errorf(openwallet.ErrSubmitRawTransactionFailed, tx.Reason)
	}

	return retTx, retFailed, nil
}
