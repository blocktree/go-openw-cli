package openwcli

import (
	"github.com/blocktree/OpenWallet/log"
	"testing"
)

func TestCLI_CreateWalletOnServer(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	err := cli.CreateWalletOnServer("gogogo", "12345678")
	if err != nil {
		log.Error("CreateWalletOnServer error:", err)
		return
	}
}

func TestCLI_GetWalletsOnServer(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	wallets, err := cli.GetWalletsOnServer()
	if err != nil {
		log.Error("GetWalletsOnServer error:", err)
		return
	}
	for i, w := range wallets {
		log.Info("wallet[", i, "]:", w)
	}
}

func TestCLI_CreateAccountOnServer(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	wallets, err := cli.GetWalletsOnServer()
	if err != nil {
		log.Error("GetWalletsOnServer error:", err)
		return
	}

	if len(wallets) > 0 {
		err = cli.CreateAccountOnServer("helleo", "12345678", "BTC", wallets[0])
		if err != nil {
			log.Error("CreateAccountOnServer error:", err)
			return
		}
	}
}

func TestCLI_GetAccountsOnServer(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	wallets, err := cli.GetWalletsOnServer()
	if err != nil {
		log.Error("GetWalletsOnServer error:", err)
		return
	}

	if len(wallets) > 0 {
		accounts, err := cli.GetAccountsOnServer(wallets[0].WalletID)
		if err != nil {
			log.Error("GetAccountsOnServer error:", err)
			return
		}
		for i, w := range accounts {
			log.Info("account[", i, "]:", w)
		}
	}
}

func TestCLI_CreateAddressOnServer(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	wallets, err := cli.GetWalletsOnServer()
	if err != nil {
		log.Error("GetWalletsOnServer error:", err)
		return
	}

	if len(wallets) > 0 {
		accounts, err := cli.GetAccountsOnServer(wallets[0].WalletID)
		if err != nil {
			log.Error("GetAccountsOnServer error:", err)
			return
		}

		if len(accounts) > 0 {
			err = cli.CreateAddressOnServer(accounts[0].WalletID, accounts[0].AccountID, 2)
			if err != nil {
				log.Error("CreateAddressOnServer error:", err)
				return
			}
		}
	}


}