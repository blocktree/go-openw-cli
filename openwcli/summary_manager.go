package openwcli

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/blocktree/go-openw-sdk/openwsdk"
	"github.com/blocktree/openwallet/common"
	"github.com/blocktree/openwallet/hdkeystore"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/owtp"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

func (cli *CLI) TransferAll(wallet *openwsdk.Wallet, account *openwsdk.Account, contractAddress, to, sid, feeRate, memo, password string) error {

	var (
		isContract  bool
		contractID  string
		tokenSymbol string
		balance     string
	)

	key, err := cli.getLocalKeyByWallet(wallet, password)
	if err != nil {
		return err
	}

	accountTask := &openwsdk.SummaryAccountTask{
		AccountID: account.AccountID,
		FeeRate:   feeRate,
		Memo:      memo,
		SummarySetting: &openwsdk.SummarySetting{
			SumAddress:      to,
			RetainedBalance: "0",
			Confirms:        0,
			MinTransfer:     "0",
		},
	}

	if len(contractAddress) > 0 {
		isContract = true
		tokenBalance, findErr := cli.GetTokenBalanceByContractAddress(account, contractAddress)
		if findErr != nil {
			return findErr
		}
		contractID = tokenBalance.ContractID
		tokenSymbol = tokenBalance.Token
		balance = tokenBalance.Balance.Balance
	} else {
		balance = account.Balance
	}
	coin := openwsdk.Coin{
		Symbol:     account.Symbol,
		IsContract: isContract,
		ContractID: contractID,
	}

	log.Infof("Summary account[%s] Symbol: %s, token: %s ", account.AccountID, account.Symbol, tokenSymbol)

	//汇总账户
	err = cli.summaryAccount(account, accountTask, key, balance, *accountTask.SummarySetting, coin, "", decimal.Zero)
	if err != nil {
		return fmt.Errorf("Summary wallet[%s] account[%s] main coin unexpected error: %v ", wallet.WalletID, account.AccountID, err)
	}

	return nil
}

//SummaryWallets 执行汇总流程
func (cli *CLI) SummaryTask() {

	log.Infof("[Summary Task Start]------%s", common.TimeFormat("2006-01-02 15:04:05"))

	cli.mu.RLock()
	defer cli.mu.RUnlock()
	//读取参与汇总的钱包
	for _, task := range cli.summaryTask.Wallets {

		if task.Wallet == nil {
			w, err := cli.GetWalletByWalletIDOnLocal(task.WalletID)
			if err != nil {
				log.Errorf("Summary wallet[%s] unexpected error: %v", task.WalletID, err)
				continue
			}
			task.Wallet = w
		}

		key, err := cli.getLocalKeyByWallet(task.Wallet, task.Password)
		if err != nil {
			log.Errorf("Summary wallet[%s] unexpected error: %v", task.WalletID, err)
			continue
		}

		for _, accountTask := range task.Accounts {

			account, err := cli.GetAccountByAccountID(accountTask.AccountID)
			if err != nil {
				continue
			}

			if !accountTask.OnlyContracts {
				//汇总账户主币
				err = cli.SummaryAccountMainCoin(accountTask, account, key)
				if err != nil {
					log.Errorf("Summary wallet[%s] account[%s] main coin unexpected error: %v", task.WalletID, account.AccountID, err)
				}
			}

			//汇总账户主币
			err = cli.SummaryAccountTokenContracts(accountTask, account, key)
			if err != nil {
				log.Errorf("Summary wallet[%s] account[%s] token contracts unexpected error: %v", task.WalletID, account.AccountID, err)
			}

		}

	}

	log.Infof("[Summary Task End]------%s", common.TimeFormat("2006-01-02 15:04:05"))
}

//SummaryAccountMainCoin 汇总账户主币
func (cli *CLI) SummaryAccountMainCoin(accountTask *openwsdk.SummaryAccountTask, account *openwsdk.Account, key *hdkeystore.HDKey) error {

	var (
		err     error
		sumSets openwsdk.SummarySetting
		symbol  string
	)

	//读取汇总信息
	err = cli.db.One("AccountID", account.AccountID, &sumSets)
	if err != nil {
		return fmt.Errorf("Summary account[%s] can not find account summary setting ", account.AccountID)
	}

	if sumSets.SumAddress == "" {
		log.Errorf("Summary account[%s] summary address is empty!", account.AccountID)
		return err
	}

	//balance, _ := decimal.NewFromString(account.Balance)
	//threshold, _ := decimal.NewFromString(sumSets.Threshold)

	//检查是否需要切换symbol
	if len(accountTask.SwitchSymbol) > 0 {
		symbol = accountTask.SwitchSymbol
	} else {
		symbol = account.Symbol
	}

	coin := openwsdk.Coin{
		Symbol:     symbol,
		IsContract: false,
	}

	log.Infof("Summary account[%s] Symbol: %s start", account.AccountID, symbol)

	err = cli.summaryAccountProcess(account, accountTask, key, account.Balance, *accountTask.SummarySetting, coin)

	log.Infof("Summary account[%s] Symbol: %s end", account.AccountID, symbol)
	log.Infof("----------------------------------------------------------------------------------------")

	if err != nil {
		return err
	}

	return nil
}

//SummaryAccountTokenContracts 汇总账户代币合约
func (cli *CLI) SummaryAccountTokenContracts(accountTask *openwsdk.SummaryAccountTask, account *openwsdk.Account, key *hdkeystore.HDKey) error {

	var (
		err     error
		sumSets openwsdk.SummarySetting
		symbol  string
	)

	if len(accountTask.Contracts) == 0 {
		return nil
	}

	//检查是否需要切换symbol
	if len(accountTask.SwitchSymbol) > 0 {
		symbol = accountTask.SwitchSymbol
	} else {
		symbol = account.Symbol
	}

	tokenBalances, err := cli.GetAllTokenContractBalance(account.AccountID, symbol)
	if err != nil {
		return err
	}

	//查询已选token过程，查地址
	findSelectedTokensFunc := func(t string) (bool, *openwsdk.SummaryContractTask) {

		if setting, ok := accountTask.Contracts["all"]; ok {
			return true, setting
		}

		for c, s := range accountTask.Contracts {
			if c == t {
				return true, s
			}
		}
		return false, nil
	}

	//读取汇总信息
	err = cli.db.One("AccountID", account.AccountID, &sumSets)
	if err != nil {
		return err
	}

	if sumSets.SumAddress == "" {
		log.Errorf("Summary account[%s] summary address is empty!")
		return err
	}

	for _, token := range tokenBalances {

		//找不到已选合约跳到下一个
		find, contrackTask := findSelectedTokensFunc(token.Address)
		if !find {
			continue
		}

		if contrackTask.SummarySetting == nil {
			contrackTask.SummarySetting = &sumSets
		} else {
			contrackTask.SummarySetting.SumAddress = sumSets.SumAddress
		}

		//查询合约余额
		//tokenBalance := cli.GetTokenBalance(account, token.ContractID)

		coin := openwsdk.Coin{
			Symbol:     symbol,
			IsContract: true,
			ContractID: token.ContractID,
		}

		log.Infof("Summary account[%s] Symbol: %s, token: %s start", account.AccountID, symbol, token.Token)

		err = cli.summaryAccountProcess(account, accountTask, key, token.Balance.Balance, *contrackTask.SummarySetting, coin)

		log.Infof("Summary account[%s] Symbol: %s, token: %s end", account.AccountID, symbol, token.Token)

		if err != nil {
			continue
		}

	}
	return nil
}

//summaryAccountProcess 汇总账户过程
func (cli *CLI) summaryAccountProcess(account *openwsdk.Account, task *openwsdk.SummaryAccountTask, key *hdkeystore.HDKey, balance string, sumSets openwsdk.SummarySetting, coin openwsdk.Coin) error {

	var (
		feesSupportAccountID string
		feesSupportBalance   = decimal.Zero
	)

	balanceDec, _ := decimal.NewFromString(balance)
	threshold, _ := decimal.NewFromString(sumSets.Threshold)

	log.Infof("Summary account[%s] Current Balance: %v, threshold: %v", account.AccountID, balance, threshold)

	// 查询手续费账户是否存在，是否在当前钱包下，相同的symbol，并且检查手续费账户余额是否报警
	if task.FeesSupportAccount != nil && coin.IsContract {
		//代币汇总才需要手续费账户
		feesSupportAccountID = task.FeesSupportAccount.AccountID
		feesSupportAccounInfo, err := cli.GetAccountByAccountID(feesSupportAccountID)
		if err != nil {
			return fmt.Errorf("fees support account: %s can not find", feesSupportAccountID)
		}

		//手续费是否合约代币
		if task.FeesSupportAccount.IsTokenContract {
			//代币作为手续费
			contractAddress := task.FeesSupportAccount.ContractAddress
			if len(contractAddress) == 0 {
				return fmt.Errorf("fees support account use token contract for fees, contract address is empty")
			}
			tokenBalance, err := cli.GetTokenBalanceByContractAddress(feesSupportAccounInfo, task.FeesSupportAccount.ContractAddress)
			if err == nil {
				feesSupportBalance, _ = decimal.NewFromString(tokenBalance.Balance.Balance)
			}

		} else {
			//主币作为手续费
			feesSupportBalance, _ = decimal.NewFromString(feesSupportAccounInfo.Balance)
		}

		lowBalanceWarning, _ := decimal.NewFromString(task.FeesSupportAccount.LowBalanceWarning)
		lowBalanceStop, _ := decimal.NewFromString(task.FeesSupportAccount.LowBalanceStop)
		if feesSupportBalance.LessThan(lowBalanceWarning) {
			log.Warningf("fees support account balance: %s is less then %s", feesSupportBalance.String(), lowBalanceWarning.String())
		}
		if feesSupportBalance.LessThan(lowBalanceStop) {
			return fmt.Errorf("fees support account: %s stop work", feesSupportBalance.String())
		}
	}

	//如果余额大于阀值，汇总的地址
	if balanceDec.GreaterThan(threshold) {
		return cli.summaryAccount(account, task, key, balance, sumSets, coin, feesSupportAccountID, feesSupportBalance)
	}

	return nil
}

//summaryAccount 汇总单个账户
func (cli *CLI) summaryAccount(account *openwsdk.Account, task *openwsdk.SummaryAccountTask,
	key *hdkeystore.HDKey, balance string, sumSets openwsdk.SummarySetting, coin openwsdk.Coin,
	feesSupportAccountID string, feesSupportBalance decimal.Decimal) error {

	const (
		limit = 200
	)

	var (
		err                  error
		createErr            error
		retTx                []*openwsdk.Transaction
		retFailed            []*openwsdk.FailedRawTransaction
		retRawTxs            []*openwsdk.RawTransaction
		retRawFeesSupportTxs []*openwsdk.RawTransaction
	)

	log.Infof("Summary account[%s] Current Balance = %v ", account.AccountID, balance)
	log.Infof("Summary account[%s] Summary Address = %v ", account.AccountID, sumSets.SumAddress)
	log.Infof("Summary account[%s] Start Create Summary Transaction", account.AccountID)

	//分页汇总交易
	for i := 0; i < int(account.AddressIndex)+1; i = i + limit {
		err = nil
		retRawTxs = make([]*openwsdk.RawTransaction, 0)
		retRawFeesSupportTxs = make([]*openwsdk.RawTransaction, 0)
		retTx = nil
		retFailed = nil

		log.Infof("Create Summary Transaction in address range [%d...%d]", i, i+limit)

		//:记录汇总批次号
		sid := uuid.New().String()
		//log.Debugf("sid: %+v", sid)
		err = cli.api.CreateSummaryTx(account.AccountID, sumSets.SumAddress, coin,
			task.FeeRate, sumSets.MinTransfer, sumSets.RetainedBalance,
			i, limit, sumSets.Confirms, sid, task.FeesSupportAccount, task.Memo, true,
			func(status uint64, msg string, rawTxs []*openwsdk.RawTransaction) {
				log.Debugf("status: %d, msg: %s", status, msg)
				for _, rawTx := range rawTxs {
					//log.Debugf("rawTx: %+v", rawTx)
					if rawTx.ErrorMsg != nil && rawTx.ErrorMsg.Code != 0 {
						log.Warning(rawTx.ErrorMsg.Err)
					} else {

						switch rawTx.AccountID {
						case account.AccountID:
							retRawTxs = append(retRawTxs, rawTx)
						case feesSupportAccountID:
							retRawFeesSupportTxs = append(retRawFeesSupportTxs, rawTx)
						}
						//if rawTx.AccountID == feesSupportAccountID {
						//	retRawFeesSupportTxs = append(retRawFeesSupportTxs, rawTx)
						//	log.Notice("create fees support account transaction for summary task")
						//}

						//retRawTxs = append(retRawTxs, rawTx)
					}
				}

				if status != owtp.StatusSuccess {
					createErr = fmt.Errorf(msg)
				}
			})

		if err != nil {
			log.Warn("CreateSummaryTransaction unexpected error: %v", err)
			continue
		}

		if createErr != nil {
			log.Warn("CreateSummaryTransaction unexpected error:", createErr)
			continue
		}

		//正常的汇总交易
		if len(retRawTxs) > 0 {

			//发送汇总交易
			signedRawTxs := make([]*openwsdk.RawTransaction, 0)
			txIDs := make([]string, 0)
			sids := make([]string, 0)
			for _, rawTx := range retRawTxs {

				//签名交易
				err = openwsdk.SignRawTransaction(rawTx, key)
				if err != nil {
					log.Warn("SignRawTransaction unexpected error: %v", err)
					continue
				}

				signedRawTxs = append(signedRawTxs, rawTx)

				//log.Debugf("retRawTxs: %+v", rawTx)
			}

			if len(signedRawTxs) == 0 {
				continue
			}

			//	广播交易单
			err = cli.api.SubmitTrade(signedRawTxs, true,
				func(status uint64, msg string, successTx []*openwsdk.Transaction, failedRawTxs []*openwsdk.FailedRawTransaction) {

					//log.Debugf("status: %d, msg: %s", status, msg)
					if status != owtp.StatusSuccess {
						createErr = fmt.Errorf(msg)
						return
					}

					retTx = successTx
					retFailed = failedRawTxs
				})
			if err != nil {
				log.Warningf("SubmitRawTransaction unexpected error: %v", err)
				continue
			}
			if createErr != nil {
				log.Warningf("SubmitRawTransaction unexpected error: %v", createErr)
				continue
			}

			//打印汇总交易结果
			totalSumAmount := decimal.Zero
			totalCostFees := decimal.Zero

			for _, tx := range retTx {

				//只计算汇总账户的 总的汇总数量，手续费
				log.Infof("[Success] txid: %s", tx.Txid)

				fees, _ := decimal.NewFromString(tx.Fees)

				totalCostFees = totalCostFees.Add(fees)
				txIDs = append(txIDs, tx.Txid)
				sids = append(sids, tx.Sid)
				//统计汇总总数
				for i, a := range tx.ToAddress {
					if a == sumSets.SumAddress {
						amount, _ := decimal.NewFromString(tx.ToAddressV[i])
						totalSumAmount = totalSumAmount.Add(amount)
					}
				}
			}

			for _, tx := range retFailed {
				log.Warn("[Failed] reason:", tx.Reason)
			}

			//:记录汇总情况
			totalSumAmount = totalSumAmount.Sub(totalCostFees)
			summaryTaskLog := openwsdk.SummaryTaskLog{
				Sid:            sid,
				WalletID:       account.WalletID,
				AccountID:      account.AccountID,
				StartAddrIndex: i,
				EndAddrIndex:   i + limit,
				Coin:           coin,
				SuccessCount:   len(retTx),
				FailCount:      len(retFailed),
				TxIDs:          txIDs,
				Sids:           sids,
				TotalSumAmount: totalSumAmount.String(),
				TotalCostFees:  totalCostFees.String(),
				CreateTime:     time.Now().Unix(),
			}
			err = cli.db.Save(&summaryTaskLog)
			if err != nil {
				log.Infof("Save summary task log failed: %s", err.Error())
			} else {
				log.Infof("Save summary task log successfully")
			}
		}

		//发送手续费账户交易单
		if len(retRawFeesSupportTxs) > 0 {

			log.Std.Notice("create fees support account transaction for summary task")
			log.Std.Notice("current fees support account balance: %s %s ", feesSupportBalance.String(), coin.Symbol)

			signedRawTxs := make([]*openwsdk.RawTransaction, 0)
			for _, rawTx := range retRawFeesSupportTxs {
				//签名交易
				err = openwsdk.SignRawTransaction(rawTx, key)
				if err != nil {
					log.Warn("SignRawTransaction unexpected error: %v", err)
					continue
				}

				signedRawTxs = append(signedRawTxs, rawTx)
			}

			if len(signedRawTxs) == 0 {
				continue
			}

			//	广播交易单
			err = cli.api.SubmitTrade(signedRawTxs, true,
				func(status uint64, msg string, successTx []*openwsdk.Transaction, failedRawTxs []*openwsdk.FailedRawTransaction) {
					if status != owtp.StatusSuccess {
						createErr = fmt.Errorf(msg)
						return
					}

					retTx = successTx
					retFailed = failedRawTxs
				})
			if err != nil {
				log.Warningf("SubmitRawTransaction unexpected error: %v", err)
				continue
			}
			if createErr != nil {
				log.Warningf("SubmitRawTransaction unexpected error: %v", createErr)
				continue
			}

			//打印手续费交易结果
			totalSupportCostFees := decimal.Zero

			for _, tx := range retTx {
				amount, _ := decimal.NewFromString(tx.Amount)
				totalSupportCostFees = totalSupportCostFees.Add(amount)
				//:手续费账户消息处理
				log.Std.Notice(" [fees support account transfer Success] txid: %s", tx.Txid)
			}

			for _, tx := range retFailed {
				log.Warn("[fees support account transfer Failed] reason:", tx.Reason)
			}

			log.Std.Notice("fees support account total cost: %s %s", totalSupportCostFees.String(), coin.Symbol)
		}
	}

	return nil
}

func (cli *CLI) signSummaryRawTransaction(retRawTxs []*openwsdk.RawTransaction, key *hdkeystore.HDKey) ([]*openwsdk.RawTransaction, error) {
	signedRawTxs := make([]*openwsdk.RawTransaction, 0)
	for _, rawTx := range retRawTxs {
		//签名交易
		err := openwsdk.SignRawTransaction(rawTx, key)
		if err != nil {
			log.Warn("SignRawTransaction unexpected error: %v", err)
			return nil, err
		}

		signedRawTxs = append(signedRawTxs, rawTx)
	}

	if len(signedRawTxs) == 0 {
		return nil, fmt.Errorf("not transactions have been signed")
	}
	return signedRawTxs, nil
}

func FindExistedSummaryWalletTask(walletID string, tasks []*openwsdk.SummaryWalletTask) (int, *openwsdk.SummaryWalletTask) {
	for i, w := range tasks {
		if w.WalletID == walletID {
			return i, w
		}
	}
	return -1, nil
}

func FindExistedSummaryAccountTask(accountID string, tasks []*openwsdk.SummaryAccountTask) (int, *openwsdk.SummaryAccountTask) {
	for i, w := range tasks {
		if w.AccountID == accountID {
			return i, w
		}
	}
	return -1, nil
}

//checkSummaryTaskIsHaveSettings 检查汇总任务中的账户是否已配置
func (cli *CLI) checkSummaryTaskIsHaveSettings(task *openwsdk.SummaryTask) error {

	for _, w := range task.Wallets {
		for _, account := range w.Accounts {

			accounInfo, err := cli.GetAccountByAccountID(account.AccountID)
			if err != nil {
				return fmt.Errorf("summary task account: %s can not find", account.AccountID)
			}

			var sumSets openwsdk.SummarySetting
			//读取汇总信息
			err = cli.db.One("AccountID", account.AccountID, &sumSets)
			if err != nil {
				return fmt.Errorf("Summary account[%s] can not find account summary setting ", account.AccountID)
			}

			if sumSets.SumAddress == "" {
				log.Errorf("Summary account[%s] summary address is empty!", account.AccountID)
				return err
			}

			if account.SummarySetting == nil {
				account.SummarySetting = &sumSets
			} else {
				account.SummarySetting.SumAddress = sumSets.SumAddress
			}

			//:查询手续费账户是否存在，是否在当前钱包下，相同的symbol，并且检查手续费账户余额是否报警
			if account.FeesSupportAccount != nil && account.FeesSupportAccount.AccountID != "" {

				feesSupportAccountInfo, err := cli.GetAccountByAccountID(account.FeesSupportAccount.AccountID)
				if err != nil {
					return fmt.Errorf("fees support account: %s can not find", account.FeesSupportAccount.AccountID)
				}
				if feesSupportAccountInfo.WalletID != accounInfo.WalletID {
					return fmt.Errorf("fees support account: %s walletID is not equal: %s", account.FeesSupportAccount.AccountID, accounInfo.WalletID)
				}

				if len(account.SwitchSymbol) > 0 {

					//SwitchSymbol是否存在
					_, findErr := cli.GetSymbolInfo(account.SwitchSymbol)
					if findErr != nil {
						return fmt.Errorf("can not find switch symbol")
					}

					//允许切换账户的symbol
					if feesSupportAccountInfo.Symbol != account.SwitchSymbol {
						return fmt.Errorf("fees support account: %s symbol is not equal summary task account switch symbol: %s", account.FeesSupportAccount.AccountID, accounInfo.Symbol)
					}
				} else {
					if feesSupportAccountInfo.Symbol != accounInfo.Symbol {
						return fmt.Errorf("fees support account: %s symbol is not equal summary task account: %s", account.FeesSupportAccount.AccountID, accounInfo.Symbol)
					}
				}

			}
		}
	}
	return nil
}

func (cli *CLI) appendSummaryTasks(sums *openwsdk.SummaryTask) {
	cli.mu.Lock()
	defer cli.mu.Unlock()

	if cli.summaryTask == nil {
		cli.summaryTask = sums
	}

	for _, newWalletTask := range sums.Wallets {

		//查找钱包是否汇总中
		_, executingWallet := FindExistedSummaryWalletTask(newWalletTask.WalletID, cli.summaryTask.Wallets)
		if executingWallet != nil {
			//钱包汇总中...
			for _, newAccountTask := range newWalletTask.Accounts {
				//查找账户是否汇总中
				_, executingAccount := FindExistedSummaryAccountTask(newAccountTask.AccountID, executingWallet.Accounts)
				if executingAccount != nil {

					if executingAccount.Contracts == nil {
						executingAccount.Contracts = make(map[string]*openwsdk.SummaryContractTask)
					}
					//executingAccount.Contracts = newAccountTask.Contracts
					//追加或替换合约
					for addr, newContractTask := range newAccountTask.Contracts {
						executingAccount.Contracts[addr] = newContractTask
					}
				} else {
					executingWallet.Accounts = append(executingWallet.Accounts, newAccountTask)
					log.Infof("Summary account[%s] task has been appended ", newAccountTask.AccountID)
				}
			}

		} else {
			cli.summaryTask.Wallets = append(cli.summaryTask.Wallets, newWalletTask)
			log.Infof("Summary wallet[%s] task has been appended ", newWalletTask.WalletID)
		}
	}

}

func (cli *CLI) removeSummaryWalletTasks(walletID string, accountID string) {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	indexWallet, executingWallet := FindExistedSummaryWalletTask(walletID, cli.summaryTask.Wallets)
	if executingWallet != nil {
		if len(accountID) > 0 {
			//查找账户是否汇总中
			indexAccount, executingAccount := FindExistedSummaryAccountTask(accountID, executingWallet.Accounts)
			if executingAccount != nil {
				//移除汇总账户任务
				executingWallet.Accounts = append(executingWallet.Accounts[:indexAccount], executingWallet.Accounts[indexAccount+1:]...)
				log.Infof("Summary account[%s] task has been removed ", accountID)
			}
		} else {
			//移除汇总钱包任务
			cli.summaryTask.Wallets = append(cli.summaryTask.Wallets[:indexWallet], cli.summaryTask.Wallets[indexWallet+1:]...)
			log.Infof("Summary wallet[%s] task has been removed ", walletID)
		}
	}
}

func (cli *CLI) GetSummaryTaskLog(offset, limit int64) ([]*openwsdk.SummaryTaskLog, error) {
	var summaryTaskLog []*openwsdk.SummaryTaskLog
	//err := cli.db.All(&summaryTaskLog)
	err := cli.db.AllByIndex("CreateTime", &summaryTaskLog,
		storm.Limit(int(limit)), storm.Skip(int(offset)), storm.Reverse())
	if err != nil {
		return nil, err
	}
	return summaryTaskLog, nil
}
