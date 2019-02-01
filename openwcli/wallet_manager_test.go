package openwcli

import (
	"fmt"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/go-openw-sdk/openwsdk"
	"testing"
)

func TestCLI_CreateWalletOnServer(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	_, err := cli.CreateWalletOnServer("testwallet", "12345678")
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
		_, _, err = cli.CreateAccountOnServer("helleo", "12345678", "BTC", wallets[0])
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
			err = cli.CreateAddressOnServer(accounts[0].WalletID, accounts[0].AccountID, 20)
			if err != nil {
				log.Error("CreateAddressOnServer error:", err)
				return
			}
		}
	}

}

func TestCLI_GetAddressesOnServer(t *testing.T) {

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
			addresses, err := cli.GetAddressesOnServer(accounts[0].WalletID, accounts[0].AccountID, 0, 50)
			if err != nil {
				log.Error("GetAddressesOnServer error:", err)
				return
			}

			for i, w := range addresses {
				log.Info("address[", i, "]:", w)
			}
		}
	}

}

func TestCLI_SearchAddressOnServer(t *testing.T) {

	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	addr := "n4LX1HnPnM4Xwy61abUFFpzqoXctHsmmeJ"

	address, err := cli.SearchAddressOnServer(addr)
	if err != nil {
		log.Error("SearchAddressOnServer error:", err)
		return
	}
	log.Info("address:", address)
	cli.printAddressList(address.WalletID, []*openwsdk.Address{address}, "12345678")
}

func TestCLI_GetSymbolList(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	symbols, err := cli.GetSymbolList()
	if err != nil {
		log.Error("GetSymbolList error:", err)
		return
	}

	for _, s := range symbols {
		fmt.Printf("symbol: %+v\n", s)
	}
}

func TestCLI_GetSymbolInfo(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	symbol, err := cli.GetSymbolInfo("BTC")
	if err != nil {
		log.Error("GetSymbolInfo error:", err)
		return
	}

	fmt.Printf("symbol: %+v\n", symbol)
}

func TestCLI_GetTokenContractList(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	tokens, err := cli.GetTokenContractList()
	if err != nil {
		log.Error("GetTokenContractList error:", err)
		return
	}

	for _, s := range tokens {
		fmt.Printf("token: %+v\n", s)
	}
}

func TestCLI_GetTokenContractInfo(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	token, err := cli.GetTokenContractInfo("zMRvEwlnUmAvPGCoaAOxINcDFI6gQDIeDlLjg19GGCg=")
	if err != nil {
		log.Error("GetTokenContractInfo error:", err)
		return
	}

	fmt.Printf("token: %+v\n", token)
}

func TestCLI_SetSummaryInfo(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	err := cli.SetSummaryInfo(&openwsdk.SummarySetting{
		"VzRF939isEwpz7wLUwqULpmhct2wsApdm4",
		"J3wiDj2jMGdp9aqmALhQtEkJQch4YN9e38TEXzRgZyKY",
		"mp1JDsi7Dr2PkcWu1j4SUSTXJqXjFMaeVx",
		"1",
		"0.1",
		"0",
		1,
	})
	if err != nil {
		log.Error("SetSummaryInfo error:", err)
		return
	}

	var sumSets openwsdk.SummarySetting
	//读取汇总信息
	err = cli.db.One("AccountID", "J3wiDj2jMGdp9aqmALhQtEkJQch4YN9e38TEXzRgZyKY", &sumSets)
	if err != nil {
		log.Error("GetSummaryInfo error:", err)
		return
	}
	fmt.Printf("SummaryInfo: %+v\n", sumSets)
}
