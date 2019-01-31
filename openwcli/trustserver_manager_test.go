package openwcli

import "testing"

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
