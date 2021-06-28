package openwcli

import (
	"bytes"
	"github.com/blocktree/openwallet/v2/log"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestCLI_ConnectTransmitNode(t *testing.T) {

	var (
		endRunning = make(chan bool, 1)
	)

	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	err := cli.ServeTransmitNode(true)
	if err != nil {
		t.Logf("ConnectTransmitNode error: %v\n", err)
		return
	}

	<-endRunning
}

func TestLocalIPAddress(t *testing.T) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		t.Logf("LocalIPAddress error: %v\n", err)
		return
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Infof("addr: %s", ipnet.IP.String())
			}
		}
	}
}

func TestExternalIPAddress(t *testing.T) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		t.Logf("ExternalIPAddress error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	log.Infof("addr: %s", string(content))
}

func TestCLI_TestRunTimeTask(t *testing.T) {

	testRunTimeTask := func (cli *CLI) {
		for {
			walletID := "W3LxqTNAcXFqW7HGcTuERRLXKdNWu17Ccx"
			accounts, cErr := cli.GetAccountsOnServer(walletID)
			if cErr != nil {
				log.Errorf("GetAccountsOnServer failed, err: %v", cErr)
			}
			log.Infof("accounts length: %d", len(accounts))
			time.Sleep(1 * time.Second)
		}
	}

	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	err := cli.ServeTransmitNode(true)
	if err != nil {
		t.Logf("ConnectTransmitNode error: %v\n", err)
		return
	}

	testRunTimeTask(cli)

}

