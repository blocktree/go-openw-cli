/*
 * Copyright 2019 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package openwcli

import (
	"encoding/json"
	"github.com/blocktree/go-openw-sdk/v2/openwsdk"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/blocktree/openwallet/v2/owtp"
	"strings"
)

// CallABI 直接调用ABI方法
func (cli *CLI) CallABI(account *openwsdk.Account, contractAddress string, abiParam []string) (*openwsdk.SmartContractCallResult, *openwallet.Error) {

	var (
		isContract    bool
		retCallResult *openwsdk.SmartContractCallResult
		err           error
		createErr     *openwallet.Error
		contractID    string
		tokenSymbol   string
	)

	if len(contractAddress) > 0 {
		isContract = true
		token, findErr := cli.GetTokenContractList("Symbol", account.Symbol, "Address", contractAddress)
		if findErr != nil {
			return nil, openwallet.ConvertError(findErr)
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
	err = api.CallSmartContractABI(account.AccountID, coin, abiParam, "", 0,
		true, func(status uint64, msg string, callResult *openwsdk.SmartContractCallResult) {
			if status != owtp.StatusSuccess {
				createErr = openwallet.Errorf(status, msg)
				return
			}
			retCallResult = callResult
		})
	if err != nil {
		return nil, openwallet.ConvertError(err)
	}
	if createErr != nil {
		return nil, createErr
	}

	//:打印交易单明细
	log.Infof("-----------------------------------------------")
	log.Infof("[%s %s TriggerABI]", account.Symbol, tokenSymbol)
	log.Infof("From Account: %s", account.AccountID)
	log.Infof("Contract Address: %s", contractAddress)
	log.Infof("ABI Param: %s", strings.Join(abiParam, ","))
	if retCallResult.Status == openwallet.SmartContractCallResultStatusSuccess {
		log.Infof("ABI Status: success")
		log.Infof("ABI Result: %v", retCallResult.Value)
	} else {
		log.Infof("ABI Status: fail")
		log.Infof("ABI Exception: %v", retCallResult.Exception)
	}
	log.Infof("-----------------------------------------------")

	return retCallResult, nil
}

// TriggerABI 触发合约ABI接口
func (cli *CLI) TriggerABI(wallet *openwsdk.Wallet, account *openwsdk.Account, symbol, contractAddress, contractABI, amount, sid, feeRate, password string, abiParam []string, raw string, rawType uint64, awaitResult bool) (*openwsdk.SmartContractReceipt, *openwallet.Error) {

	var (
		isContract  bool
		retReceipt  *openwsdk.SmartContractReceipt
		retTx       []*openwsdk.SmartContractReceipt
		retFailed   []*openwsdk.FailureSmartContractLog
		retRawTx    *openwsdk.SmartContractRawTransaction
		err         error
		createErr   *openwallet.Error
		contractID  string
		tokenSymbol string
	)

	if len(password) == 0 {
		return nil, openwallet.Errorf(openwallet.ErrCreateRawTransactionFailed, "unlock wallet password is empty. ")
	}

	//获取种子文件
	key, err := cli.getLocalKeyByWallet(wallet, password)
	if err != nil {
		return nil, openwallet.Errorf(openwallet.ErrCreateRawTransactionFailed, err.Error())
	}

	if len(contractAddress) > 0 {
		isContract = true
		token, findErr := cli.GetTokenContractList("Symbol", account.Symbol, "Address", contractAddress)
		if findErr == nil && len(token) > 0 {
			contractID = token[0].ContractID
			tokenSymbol = token[0].Token
		}
	}
	coin := openwsdk.Coin{
		Symbol:          symbol,
		IsContract:      isContract,
		ContractID:      contractID,
		ContractAddress: contractAddress,
		ContractABI:     contractABI,
	}

	api := cli.api
	err = api.CreateSmartContractTrade(sid, account.AccountID, coin, abiParam, raw, rawType, feeRate, amount,
		true, func(status uint64, msg string, rawTx *openwsdk.SmartContractRawTransaction) {
			if status != owtp.StatusSuccess {
				createErr = openwallet.Errorf(status, msg)
				return
			}
			retRawTx = rawTx
		})
	if err != nil {
		return nil, openwallet.ConvertError(err)
	}
	if createErr != nil {
		return nil, createErr
	}

	//:打印交易单明细
	log.Infof("-----------------------------------------------")
	log.Infof("[%s %s TriggerABI]", account.Symbol, tokenSymbol)
	log.Infof("SID: %s", retRawTx.Sid)
	log.Infof("From Account: %s", account.AccountID)
	log.Infof("Contract Address: %s", contractAddress)
	log.Infof("ABI Param: %s", strings.Join(abiParam, ","))
	log.Infof("Fees: %v", retRawTx.Fees)
	log.Infof("-----------------------------------------------")

	//签名交易单
	signatures, sigErr := cli.txSigner(retRawTx.Signatures, key)
	if sigErr != nil {
		return nil, openwallet.Errorf(openwallet.ErrSignRawTransactionFailed, sigErr.Error())
	}
	retRawTx.Signatures = signatures
	retRawTx.AwaitResult = awaitResult
	retRawTx.AwaitTimeout = uint64(cli.config.requesttimeout)
	//广播交易单
	err = api.SubmitSmartContractTrade([]*openwsdk.SmartContractRawTransaction{retRawTx}, true,
		func(status uint64, msg string, successTx []*openwsdk.SmartContractReceipt, failedRawTxs []*openwsdk.FailureSmartContractLog) {
			if status != owtp.StatusSuccess {
				createErr = openwallet.Errorf(status, msg)
				return
			}

			retTx = successTx
			retFailed = failedRawTxs
		})
	if err != nil {
		return nil, openwallet.ConvertError(err)
	}
	if createErr != nil {
		return nil, createErr
	}

	if len(retTx) > 0 {
		//打印交易单
		log.Info("send transaction successfully.")
		log.Info("transaction id:", retTx[0].TxID)
		retReceipt = retTx[0]
	} else if len(retFailed) > 0 {
		//打印交易单
		log.Errorf("send transaction failed.")
		tx := retFailed[0]
		log.Warningf("[Failed] reason: %s", tx.Reason)
		if tx.RawTx != nil {
			log.Warningf("[Failed] rawHex: %s", tx.RawTx.Raw)
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

		return nil, openwallet.Errorf(openwallet.ErrSubmitRawSmartContractTransactionFailed, tx.Reason)
	}

	if retReceipt != nil && retReceipt.BlockHeight > 0 {

		if retReceipt.Success == openwallet.TxStatusSuccess {
			log.Infof("get receipt status: success")
			log.Info("show receipt events:")
			for _, event := range retReceipt.Events {
				log.Infof("contract address: %s, contract name: %s", event.ContractAddr, event.ContractName)
				log.Infof("[%s]%s", event.Event, event.Value)
			}
		} else {
			log.Infof("get receipt status: fail")
		}
	}

	return retReceipt, nil
}
