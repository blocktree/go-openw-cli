package openwcli

import (
	"encoding/json"
	"fmt"
	"github.com/blocktree/go-openw-sdk/v2/openwsdk"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/blocktree/openwallet/v2/owtp"
	"github.com/blocktree/openwallet/v2/timer"
	"time"
)

const (
	trustHostID = "TrustServer"
)

func (cli *CLI) ServeTransmitNode(autoReconnect bool) error {

	keychain, err := cli.GetKeychain()
	if err != nil {
		return err
	}
	cert, _ := keychain.Certificate()
	node := owtp.NewNode(owtp.NodeConfig{
		Cert:       cert,
		TimeoutSEC: cli.config.requesttimeout,
	})

	cli.transmitNode = node

	//绑定本地路由方法
	cli.transmitNode.HandleFunc("getTrustNodeInfo", cli.getTrustNodeInfo)
	cli.transmitNode.HandleFunc("createWalletViaTrustNode", cli.createWalletViaTrustNode)
	cli.transmitNode.HandleFunc("createAccountViaTrustNode", cli.createAccountViaTrustNode)
	cli.transmitNode.HandleFunc("sendTransactionViaTrustNode", cli.sendTransactionViaTrustNode)
	cli.transmitNode.HandleFunc("setSummaryInfoViaTrustNode", cli.setSummaryInfoViaTrustNode)
	cli.transmitNode.HandleFunc("findSummaryInfoByWalletIDViaTrustNode", cli.findSummaryInfoByWalletIDViaTrustNode)
	cli.transmitNode.HandleFunc("startSummaryTaskViaTrustNode", cli.startSummaryTaskViaTrustNode)
	cli.transmitNode.HandleFunc("stopSummaryTaskViaTrustNode", cli.stopSummaryTaskViaTrustNode)
	cli.transmitNode.HandleFunc("updateInfoViaTrustNode", cli.updateInfoViaTrustNode)
	cli.transmitNode.HandleFunc("appendSummaryTaskViaTrustNode", cli.appendSummaryTaskViaTrustNode)
	cli.transmitNode.HandleFunc("removeSummaryTaskViaTrustNode", cli.removeSummaryTaskViaTrustNode)
	cli.transmitNode.HandleFunc("getCurrentSummaryTaskViaTrustNode", cli.getCurrentSummaryTaskViaTrustNode)
	cli.transmitNode.HandleFunc("getSummaryTaskLogViaTrustNode", cli.getSummaryTaskLogViaTrustNode)
	cli.transmitNode.HandleFunc("getLocalWalletListViaTrustNode", cli.getLocalWalletListViaTrustNode)
	cli.transmitNode.HandleFunc("getTrustAddressListViaTrustNode", cli.getTrustAddressListViaTrustNode)
	cli.transmitNode.HandleFunc("signTransactionViaTrustNode", cli.signTransactionViaTrustNode)
	cli.transmitNode.HandleFunc("triggerABIViaTrustNode", cli.triggerABIViaTrustNode)
	cli.transmitNode.HandleFunc("signHashViaTrustNode", cli.signHashViaTrustNode)

	//自动连接
	if autoReconnect {
		go cli.autoReconnectTransmitNode()
		return nil
	}

	//单独连接
	err = cli.connectTransmitNode()
	if err != nil {
		return err
	}

	return nil
}

//connectTransmitNode
func (cli *CLI) connectTransmitNode() error {

	connectCfg := owtp.ConnectConfig{}
	connectCfg.Address = cli.config.trustedserver
	connectCfg.ConnectType = owtp.Websocket
	connectCfg.EnableSSL = cli.config.enabletrustserverssl
	connectCfg.EnableSignature = false
	connectCfg.EnableKeyAgreement = cli.config.enablekeyagreement

	//建立连接
	_, err := cli.transmitNode.Connect(trustHostID, connectCfg)
	if err != nil {
		return err
	}

	//开启协商密码
	//if cli.config.enablekeyagreement {
	//	if err = cli.transmitNode.KeyAgreement(trustHostID, "aes"); err != nil {
	//		return err
	//	}
	//}

	//向服务器发送连接成功
	err = cli.nodeDidConnectedServer()
	if err != nil {
		return err
	}

	return nil
}

//Run 运行商户节点管理
func (cli *CLI) autoReconnectTransmitNode() error {

	var (
		err error
		//连接状态通道
		reconnect = make(chan bool, 1)
		//断开状态通道
		disconnected = make(chan struct{}, 1)
		//重连时的等待时间
		reconnectWait = 5
	)

	defer func() {
		close(reconnect)
		close(disconnected)
	}()

	//断开连接通知
	cli.transmitNode.SetCloseHandler(func(n *owtp.OWTPNode, peer owtp.PeerInfo) {
		disconnected <- struct{}{}
	})

	//启动连接
	reconnect <- true

	//节点运行时
	for {
		select {
		case <-reconnect:
			//重新连接
			log.Info("Connecting to", cli.config.trustedserver)
			err = cli.connectTransmitNode()
			if err != nil {
				log.Errorf("Connect %s node failed unexpected error: %v", trustHostID, err)
				disconnected <- struct{}{}
			} else {
				log.Infof("Connect %s node successfully.", trustHostID)
			}

		case <-disconnected:
			//重新连接，前等待
			log.Info("Auto reconnect after", reconnectWait, "seconds...")
			time.Sleep(time.Duration(reconnectWait) * time.Second)
			reconnect <- true
		}
	}

	return nil
}

/*********** 客户服务平台业务方法调用 ***********/

func (cli *CLI) nodeDidConnectedServer() error {

	params := map[string]interface{}{
		"appID": cli.config.appid,
		"nodeInfo": openwsdk.TrustNodeInfo{
			NodeID:      cli.transmitNode.NodeID(),
			NodeName:    cli.config.localname,
			ConnectType: owtp.Websocket,
		},
	}

	err := cli.transmitNode.Call(trustHostID, "newNodeJoin", params,
		true, func(resp owtp.Response) {
			if resp.Status != owtp.StatusSuccess {
				log.Error(resp.Msg)
			}
		})

	return err
}

/*********** 本地路由方法实现 ***********/

func (cli *CLI) getTrustNodeInfo(ctx *owtp.Context) {
	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	info := openwsdk.TrustNodeInfo{
		NodeID:      cli.transmitNode.NodeID(),
		NodeName:    cli.config.localname,
		ConnectType: owtp.Websocket,
		Version:     Version,
		GitRev:      GitRev,
		BuildTime:   BuildTime,
	}

	ctx.Response(info, owtp.StatusSuccess, "success")
}

func (cli *CLI) createWalletViaTrustNode(ctx *owtp.Context) {
	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	alias := ctx.Params().Get("alias").String()
	password := ctx.Params().Get("password").String()

	wallet, err := cli.CreateWalletOnServer(alias, password)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}
	ctx.Response(wallet, owtp.StatusSuccess, "success")
}

func (cli *CLI) createAccountViaTrustNode(ctx *owtp.Context) {

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	alias := ctx.Params().Get("alias").String()
	walletID := ctx.Params().Get("walletID").String()
	symbol := ctx.Params().Get("symbol").String()
	password := ctx.Params().Get("password").String()

	if len(password) == 0 {
		//钱包是否已经解锁
		if p, exist := cli.unlockWallets[walletID]; exist {
			password = p
		}
	}

	wallet, err := cli.GetWalletByWalletID(walletID)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	account, addresses, err := cli.CreateAccountOnServer(alias, password, symbol, wallet)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}
	ctx.Response(map[string]interface{}{
		"account": account,
		"address": addresses,
	}, owtp.StatusSuccess, "success")

}

func (cli *CLI) sendTransactionViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enablerequesttransfer {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [transfer] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	accountID := ctx.Params().Get("accountID").String()
	sid := ctx.Params().Get("sid").String()
	contractAddress := ctx.Params().Get("contractAddress").String()
	password := ctx.Params().Get("password").String()
	amount := ctx.Params().Get("amount").String()
	address := ctx.Params().Get("address").String()
	feeRate := ctx.Params().Get("feeRate").String()
	memo := ctx.Params().Get("memo").String()
	extParam := ctx.Params().Get("extParam").String()

	account, err := cli.GetAccountByAccountID(accountID)
	if err != nil {
		ctx.Response(nil, openwallet.ErrAccountNotFound, err.Error())
		return
	}

	wallet, err := cli.GetWalletByWalletID(account.WalletID)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	if len(password) == 0 {
		//钱包是否已经解锁
		if p, exist := cli.unlockWallets[wallet.WalletID]; exist {
			password = p
		}
	}

	retTx, retFailed, exErr := cli.TransferExt(wallet, account, contractAddress, address, amount, sid, feeRate, memo, extParam, password)
	if exErr != nil {
		ctx.Response(nil, exErr.Code(), exErr.Error())
		return
	}

	ctx.Response(map[string]interface{}{
		"failure": retFailed,
		"success": retTx,
	}, owtp.StatusSuccess, "success")
}

func (cli *CLI) setSummaryInfoViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enableeditsummarysettings {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [edit summary settings] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	summarySetting := openwsdk.NewSummarySetting(ctx.Params().Get("summarySetting"))

	//汇总配置是否已初始化，若初始化后不能再有信任节点设置
	setup, err := cli.getSummarySettingByAccount(summarySetting.AccountID)
	if setup != nil {
		ctx.Response(nil, ErrorSummarySettingFailed, "summary setting has been initialized")
		return
	}

	err = cli.SetSummaryInfo(summarySetting)
	if err != nil {
		ctx.Response(nil, ErrorSummarySettingFailed, "summary info save failed")
		return
	}

	ctx.Response(nil, owtp.StatusSuccess, "success")

}

func (cli *CLI) findSummaryInfoByWalletIDViaTrustNode(ctx *owtp.Context) {

	var (
		//err     error
		sumSets []*openwsdk.SummarySetting
	)

	_, err := cli.getDB()
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, "cli database open failed")
		return
	}
	defer cli.closeDB()

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}
	walletID := ctx.Params().Get("walletID").String()

	//读取汇总配置
	cli.db.Find("WalletID", walletID, &sumSets)
	//if err != nil {
	//	ctx.Response(nil, openwallet.ErrUnknownException, "can not find summary info")
	//	return
	//}

	ctx.Response(sumSets, owtp.StatusSuccess, "success")
}

func (cli *CLI) startSummaryTaskViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enableexecutesummarytask {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [execute summary task] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()
	operateType := ctx.Params().Get("operateType").Int()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	summaryTask := openwsdk.NewSummaryTask(ctx.Params().Get("summaryTask"))
	cycleSec := ctx.Params().Get("cycleSec").Int()

	if cycleSec <= 0 {
		ctx.Response(nil, openwallet.ErrUnknownException, "cycleSec must be greater than 0")
		return
	}

	//检查汇总任务的参数是否传入密码
	for _, summaryWalletTask := range summaryTask.Wallets {
		if len(summaryWalletTask.Password) == 0 {
			//钱包是否已经解锁
			if p, exist := cli.unlockWallets[summaryWalletTask.WalletID]; exist {
				summaryWalletTask.Password = p
			}
		}
	}

	//:先检查汇总任务是否有汇总配置
	err := cli.checkSummaryTaskIsHaveSettings(summaryTask)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	switch operateType {
	case openwsdk.SummaryTaskOperateTypeReset:

		cli.mu.Lock()
		cli.summaryTask = summaryTask
		cli.mu.Unlock()

	case openwsdk.SummaryTaskOperateTypeAdd:
		cli.appendSummaryTasks(summaryTask)
	}

	if cli.summaryTaskTimer != nil && cli.summaryTaskTimer.Running() {
		log.Warning("summary task timer is running")
		//ctx.Response(nil, ErrorSummaryTaskTimerIsRunning, "summary task timer is running")
		//return
	} else {

		log.Infof("The timer for summary task start now. Execute by every %v seconds.", cycleSec)

		//启动钱包汇总程序
		sumTimer := timer.NewTask(time.Duration(cycleSec)*time.Second, cli.SummaryTask)
		sumTimer.Start()
		cli.summaryTaskTimer = sumTimer
		//马上执行一次汇总
		go cli.SummaryTask()

	}

	ctx.Response(nil, owtp.StatusSuccess, "The timer for summary task start running")

}

func (cli *CLI) stopSummaryTaskViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enableexecutesummarytask {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [execute summary task] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	if cli.summaryTaskTimer != nil && cli.summaryTaskTimer.Running() {
		cli.summaryTaskTimer.Stop()
		cli.summaryTaskTimer = nil
	}

	log.Infof("The timer for summary task has been stopped.")

	ctx.Response(nil, owtp.StatusSuccess, "success")
}

func (cli *CLI) updateInfoViaTrustNode(ctx *owtp.Context) {

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	err := cli.UpdateSymbols()
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	ctx.Response(nil, owtp.StatusSuccess, "success")
}

func (cli *CLI) appendSummaryTaskViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enableexecutesummarytask {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [execute summary task] ability")
		return
	}

	if cli.summaryTaskTimer == nil || !cli.summaryTaskTimer.Running() {
		ctx.Response(nil, ErrorSummaryTaskTimerIsNotStart, "summary task timer is not start")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	summaryTask := openwsdk.NewSummaryTask(ctx.Params().Get("summaryTask"))

	//检查汇总任务的参数是否传入密码
	for _, summaryWalletTask := range summaryTask.Wallets {
		if len(summaryWalletTask.Password) == 0 {
			//钱包是否已经解锁
			if p, exist := cli.unlockWallets[summaryWalletTask.WalletID]; exist {
				summaryWalletTask.Password = p
			}
		}
	}

	//:先检查汇总任务是否有汇总配置
	err := cli.checkSummaryTaskIsHaveSettings(summaryTask)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	cli.appendSummaryTasks(summaryTask)

	ctx.Response(nil, owtp.StatusSuccess, "success")

}

func (cli *CLI) removeSummaryTaskViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enableexecutesummarytask {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [execute summary task] ability")
		return
	}

	if cli.summaryTaskTimer == nil || !cli.summaryTaskTimer.Running() {
		ctx.Response(nil, ErrorSummaryTaskTimerIsNotStart, "summary task timer is not start")
		return
	}

	appID := ctx.Params().Get("appID").String()
	walletID := ctx.Params().Get("walletID").String()
	accountID := ctx.Params().Get("accountID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	cli.removeSummaryWalletTasks(walletID, accountID)

	ctx.Response(nil, owtp.StatusSuccess, "success")

}

func (cli *CLI) getCurrentSummaryTaskViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enableexecutesummarytask {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [execute summary task] ability")
		return
	}

	if cli.summaryTaskTimer == nil || !cli.summaryTaskTimer.Running() {
		ctx.Response(nil, ErrorSummaryTaskTimerIsNotStart, "summary task timer is not start")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	retTask := openwsdk.SummaryTask{Wallets: make([]*openwsdk.SummaryWalletTask, 0)}
	for _, wt := range cli.summaryTask.Wallets {
		newWt := *wt
		newWt.Password = ""
		newWt.Wallet = nil
		retTask.Wallets = append(retTask.Wallets, &newWt)
	}

	ctx.Response(retTask, owtp.StatusSuccess, "success")
}

func (cli *CLI) getSummaryTaskLogViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enableexecutesummarytask {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [execute summary task] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()
	offset := ctx.Params().Get("offset").Int()
	limit := ctx.Params().Get("limit").Int()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	logs, err := cli.GetSummaryTaskLog(offset, limit)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	ctx.Response(logs, owtp.StatusSuccess, "success")
}

func (cli *CLI) getLocalWalletListViaTrustNode(ctx *owtp.Context) {

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	wallets, err := cli.GetWalletsOnServer()
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	ctx.Response(wallets, owtp.StatusSuccess, "success")
}

func (cli *CLI) getTrustAddressListViaTrustNode(ctx *owtp.Context) {

	appID := ctx.Params().Get("appID").String()
	symbol := ctx.Params().Get("symbol").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	list, err := cli.ListTrustAddress(symbol)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	status := cli.TrustAddressStatus()

	ctx.Response(map[string]interface{}{
		"trustAddressList":   list,
		"enableTrustAddress": status,
	}, owtp.StatusSuccess, "success")

}

func (cli *CLI) signTransactionViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enablerequesttransfer {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [transfer] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()
	walletID := ctx.Params().Get("walletID").String()
	password := ctx.Params().Get("password").String()
	jsonRawTx := ctx.Params().Get("rawTx")

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	var rawTx openwsdk.RawTransaction
	err := json.Unmarshal([]byte(jsonRawTx.Raw), &rawTx)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	destination := ""
	amount := ""
	for to, a := range rawTx.To {
		//:检查目标地址是否信任名单
		if !cli.IsTrustAddress(to, rawTx.Coin.Symbol) {
			msg := fmt.Sprintf("%s is not in trust address list", to)
			ctx.Response(nil, openwallet.ErrUnknownException, msg)
			return
		}
		destination = to
		amount = a
	}

	wallet, err := cli.GetWalletByWalletID(walletID)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	if len(password) == 0 {
		//钱包是否已经解锁
		if p, exist := cli.unlockWallets[wallet.WalletID]; exist {
			password = p
		}
	}

	//获取种子文件
	key, err := cli.getLocalKeyByWallet(wallet, password)
	if err != nil {
		ctx.Response(nil, openwallet.ErrSignRawTransactionFailed, err.Error())
		return
	}

	tokenSymbol := ""
	if rawTx.Coin.IsContract {
		tokenContract, _ := cli.GetTokenContractInfo(rawTx.Coin.ContractID)
		if tokenContract != nil {
			tokenSymbol = tokenContract.Token
		}
	}

	//:打印交易单明细
	log.Infof("-----------------------------------------------")
	log.Infof("[%s %s Sign Transaction]", rawTx.Coin.Symbol, tokenSymbol)
	log.Infof("SID: %s", rawTx.Sid)
	log.Infof("From Account: %s", rawTx.AccountID)
	log.Infof("To Address: %s", destination)
	log.Infof("Send Amount: %s", amount)
	log.Infof("Fees: %v", rawTx.Fees)
	log.Infof("FeeRate: %v", rawTx.FeeRate)
	log.Infof("-----------------------------------------------")

	//签名交易
	signature, sigErr := cli.txSigner(rawTx.Signatures, key)
	if sigErr != nil {
		ctx.Response(nil, openwallet.ErrSignRawTransactionFailed, sigErr.Error())
		return
	}

	rawTx.Signatures = signature

	ctx.Response(map[string]interface{}{
		"signedRawTx": rawTx,
	}, owtp.StatusSuccess, "success")
}

// TriggerABIViaTrustNode 触发ABI上链调用
// @param nodeID 必填 节点ID
// @param accountID 必填 账户ID
// @param password 可选 钱包解锁密码
// @param contractAddress 必填 合约地址
// @param contractABI 可选 ABI定义
// @param amount 可选 主币数量
// @param feeRate 可选 自定义手续费率
// @param abiParam 可选 ABI参数组
// @param raw 可选 原始交易单
// @param rawType 可选 原始交易单编码类型，0：hex字符串，1：json字符串，2：base64字符串
func (cli *CLI) triggerABIViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enablerequesttransfer {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [transfer] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	accountID := ctx.Params().Get("accountID").String()
	sid := ctx.Params().Get("sid").String()
	contractAddress := ctx.Params().Get("contractAddress").String()
	contractABI := ctx.Params().Get("contractABI").String()
	password := ctx.Params().Get("password").String()
	abiArr := ctx.Params().Get("abiParam")
	amount := ctx.Params().Get("amount").String()
	feeRate := ctx.Params().Get("feeRate").String()
	raw := ctx.Params().Get("raw").String()
	rawType := ctx.Params().Get("rawType").Uint()
	awaitResult := ctx.Params().Get("awaitResult").Bool()

	abiParam := make([]string, 0)
	for _, s := range abiArr.Array() {
		abiParam = append(abiParam, s.String())
	}

	account, err := cli.GetAccountByAccountID(accountID)
	if err != nil {
		ctx.Response(nil, openwallet.ErrAccountNotFound, err.Error())
		return
	}

	wallet, err := cli.GetWalletByWalletID(account.WalletID)
	if err != nil {
		ctx.Response(nil, openwallet.ErrUnknownException, err.Error())
		return
	}

	if len(password) == 0 {
		//钱包是否已经解锁
		if p, exist := cli.unlockWallets[wallet.WalletID]; exist {
			password = p
		}
	}

	retTx, exErr := cli.TriggerABI(wallet, account, contractAddress, contractABI, amount, sid, feeRate, password, abiParam, raw, rawType, awaitResult)
	if exErr != nil {
		ctx.Response(nil, exErr.Code(), exErr.Error())
		return
	}

	ctx.Response(retTx, owtp.StatusSuccess, "success")
}

// signHashViaTrustNode 通过节点签名哈希消息
func (cli *CLI) signHashViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enablerequesttransfer {
		ctx.Response(nil, ErrorNodeAbilityDisabled, "the node has disabled [transfer] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, ErrorAppIDIncorrect, "appID is incorrect")
		return
	}

	walletID := ctx.Params().Get("walletID").String()
	accountID := ctx.Params().Get("accountID").String()
	message := ctx.Params().Get("message").String()
	password := ctx.Params().Get("password").String()
	address := ctx.Params().Get("address").String()
	symbol := ctx.Params().Get("symbol").String()
	hdPath := ctx.Params().Get("hdPath").String()

	if len(password) == 0 {
		//钱包是否已经解锁
		if p, exist := cli.unlockWallets[walletID]; exist {
			password = p
		}
	}

	addr := &openwsdk.Address{
		AppID:            appID,
		WalletID:         walletID,
		AccountID:        accountID,
		Symbol:           symbol,
		Address:          address,
		HdPath:           hdPath,
	}

	signature, err := cli.SignHash(addr, message, password)
	if err != nil {
		ctx.Response(nil, openwallet.ErrSystemException, err.Error())
		return
	}

	ctx.Response(map[string]interface{}{
		"signature": signature,
	}, owtp.StatusSuccess, "success")
}
