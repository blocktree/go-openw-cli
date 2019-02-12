package openwcli

import (
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/OpenWallet/owtp"
	"github.com/blocktree/OpenWallet/timer"
	"github.com/blocktree/go-openw-sdk/openwsdk"
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
	connectCfg.EnableSSL = cli.config.enablessl
	connectCfg.EnableSignature = false

	//建立连接
	err := cli.transmitNode.Connect(trustHostID, connectCfg)
	if err != nil {
		return err
	}

	//开启协商密码
	if cli.config.enablekeyagreement {
		if err = cli.transmitNode.KeyAgreement(trustHostID, "aes"); err != nil {
			return err
		}
	}

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
		reconnectWait = 10
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
				log.Infof("Connect %s node successfully. \n", trustHostID)
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
		ctx.Response(nil, owtp.ErrCustomError, "appID is incorrect")
		return
	}

	info := openwsdk.TrustNodeInfo{
		NodeID:      cli.transmitNode.NodeID(),
		NodeName:    cli.config.localname,
		ConnectType: owtp.Websocket,
	}

	ctx.Response(info, owtp.StatusSuccess, "success")
}

func (cli *CLI) createWalletViaTrustNode(ctx *owtp.Context) {
	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, owtp.ErrCustomError, "appID is incorrect")
		return
	}

	alias := ctx.Params().Get("alias").String()
	password := ctx.Params().Get("password").String()

	wallet, err := cli.CreateWalletOnServer(alias, password)
	if err != nil {
		ctx.Response(nil, owtp.ErrCustomError, err.Error())
		return
	}
	ctx.Response(wallet, owtp.StatusSuccess, "success")
}

func (cli *CLI) createAccountViaTrustNode(ctx *owtp.Context) {

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, owtp.ErrCustomError, "appID is incorrect")
		return
	}

	alias := ctx.Params().Get("alias").String()
	walletID := ctx.Params().Get("walletID").String()
	symbol := ctx.Params().Get("symbol").String()
	password := ctx.Params().Get("password").String()

	wallet, err := cli.GetWalletByWalletID(walletID)
	if err != nil {
		ctx.Response(nil, owtp.ErrCustomError, err.Error())
		return
	}

	account, addresses, err := cli.CreateAccountOnServer(alias, password, symbol, wallet)
	if err != nil {
		ctx.Response(nil, owtp.ErrCustomError, err.Error())
		return
	}
	ctx.Response(map[string]interface{}{
		"account": account,
		"address": addresses,
	}, owtp.StatusSuccess, "success")

}

func (cli *CLI) sendTransactionViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enablerequesttransfer {
		ctx.Response(nil, owtp.ErrCustomError, "the node has disabled [transfer] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, owtp.ErrCustomError, "appID is incorrect")
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

	account, err := cli.GetAccountByAccountID(accountID)
	if err != nil {
		ctx.Response(nil, owtp.ErrCustomError, err.Error())
		return
	}

	wallet, err := cli.GetWalletByWalletID(account.WalletID)
	if err != nil {
		ctx.Response(nil, owtp.ErrCustomError, err.Error())
		return
	}

	retTx, retFailed, err := cli.Transfer(wallet, account, contractAddress, address, amount, sid, feeRate, memo, password)
	ctx.Response(map[string]interface{}{
		"failure": retFailed,
		"success": retTx,
	}, owtp.StatusSuccess, "success")
}

func (cli *CLI) setSummaryInfoViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enableexecutesummarytask {
		ctx.Response(nil, owtp.ErrCustomError, "the node has disabled [transfer] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, owtp.ErrCustomError, "appID is incorrect")
		return
	}

	summarySetting := openwsdk.NewSummarySetting(ctx.Params().Get("summarySetting"))
	err := cli.SetSummaryInfo(summarySetting)
	if err != nil {
		ctx.Response(nil, owtp.ErrCustomError, "summary info save failed")
		return
	}

	ctx.Response(nil, owtp.StatusSuccess, "success")

}

func (cli *CLI) findSummaryInfoByWalletIDViaTrustNode(ctx *owtp.Context) {

	var (
		err     error
		sumSets []*openwsdk.SummarySetting
	)

	if !cli.config.enableexecutesummarytask {
		ctx.Response(nil, owtp.ErrCustomError, "the node has disabled [transfer] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, owtp.ErrCustomError, "appID is incorrect")
		return
	}
	walletID := ctx.Params().Get("walletID").String()

	//读取汇总配置
	err = cli.db.Find("WalletID", walletID, &sumSets)
	if err != nil {
		ctx.Response(nil, owtp.ErrCustomError, "can not find summary info")
		return
	}

	ctx.Response(sumSets, owtp.StatusSuccess, "success")
}

func (cli *CLI) startSummaryTaskViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enableexecutesummarytask {
		ctx.Response(nil, owtp.ErrCustomError, "the node has disabled [transfer] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, owtp.ErrCustomError, "appID is incorrect")
		return
	}

	if cli.summaryTaskTimer != nil && cli.summaryTaskTimer.Running() {
		ctx.Response(nil, owtp.ErrCustomError, "summary task timer is running")
		return
	}

	summaryTask := openwsdk.NewSummaryTask(ctx.Params().Get("summaryTask"))
	cycleSec := ctx.Params().Get("cycleSec").Int()
	cli.summaryTask = summaryTask

	log.Infof("The timer for summary task start now. Execute by every %v seconds.", cycleSec)

	//启动钱包汇总程序
	sumTimer := timer.NewTask(time.Duration(cycleSec)*time.Second, cli.SummaryTask)
	sumTimer.Start()

	cli.summaryTaskTimer = sumTimer

	ctx.Response(nil, owtp.StatusSuccess, "success")

}

func (cli *CLI) stopSummaryTaskViaTrustNode(ctx *owtp.Context) {

	if !cli.config.enableexecutesummarytask {
		ctx.Response(nil, owtp.ErrCustomError, "the node has disabled [transfer] ability")
		return
	}

	appID := ctx.Params().Get("appID").String()

	if appID != cli.config.appid {
		ctx.Response(nil, owtp.ErrCustomError, "appID is incorrect")
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
		ctx.Response(nil, owtp.ErrCustomError, "appID is incorrect")
		return
	}

	err := cli.UpdateSymbols()
	if err != nil {
		ctx.Response(nil, owtp.ErrCustomError, err.Error())
		return
	}

	ctx.Response(nil, owtp.StatusSuccess, "success")
}
