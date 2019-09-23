package openwcli

import (
	"fmt"
	"testing"

	"github.com/blocktree/go-openw-sdk/openwsdk"
	"github.com/blocktree/openwallet/log"
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
	walletID := "W3LxqTNAcXFqW7HGcTuERRLXKdNWu17Ccx"
	wallet, err := cli.GetWalletByWalletID(walletID)
	if err != nil {
		log.Error("GetWalletByWalletID error:", err)
		return
	}

	if wallet != nil {
		_, _, err = cli.CreateAccountOnServer("mainnetLTC", "12345678", "LTC", wallet)
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

	walletID := "W3LxqTNAcXFqW7HGcTuERRLXKdNWu17Ccx"
	accounts, err := cli.GetAccountsOnServer(walletID)
	if err != nil {
		log.Error("GetAccountsOnServer error:", err)
		return
	}
	for i, w := range accounts {
		log.Info("account[", i, "]:", w)
	}
}

func TestCLI_GetAccountOnServer(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	accountID := "PgHCcfMbcw1zXRNZo23NFjRdBmcN5tzrb1j5McRLJbG"
	account, err := cli.GetAccountByAccountID(accountID)
	if err != nil {
		log.Error("GetAccountByAccountID error:", err)
		return
	}
	log.Infof("account: %+v", account)
}

func TestCLI_CreateAddressOnServer(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	walletID := "WN84dVZXpgVixsvXnU8jkFWD1qWHp15LpA"
	accountID := "7ww2Gpfy8pN6HTngbMFBTEMAaVRGEpkmsiNkgAgqGQGf"
	err := cli.CreateAddressOnServer(walletID, accountID, 1000)
	if err != nil {
		log.Error("CreateAddressOnServer error:", err)
		return
	}

}

func TestCLI_GetAddressesOnServer(t *testing.T) {

	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	walletID := "W3LxqTNAcXFqW7HGcTuERRLXKdNWu17Ccx"
	accountID := "PgHCcfMbcw1zXRNZo23NFjRdBmcN5tzrb1j5McRLJbG"

	addresses, err := cli.GetAddressesOnServer(walletID, accountID, 0, 50)
	if err != nil {
		log.Error("GetAddressesOnServer error:", err)
		return
	}

	for i, w := range addresses {
		log.Info("address[", i, "]:", w)
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

//func TestCLI_SetSummaryInfo(t *testing.T) {
//	cli := getTestOpenwCLI()
//	if cli == nil {
//		return
//	}
//	err := cli.SetSummaryInfo(&openwsdk.SummarySetting{
//		"VzRF939isEwpz7wLUwqULpmhct2wsApdm4",
//		"J3wiDj2jMGdp9aqmALhQtEkJQch4YN9e38TEXzRgZyKY",
//		"mp1JDsi7Dr2PkcWu1j4SUSTXJqXjFMaeVx",
//		"1",
//		"0.1",
//		"0",
//		1,
//	})
//	if err != nil {
//		log.Error("SetSummaryInfo error:", err)
//		return
//	}
//
//	var sumSets openwsdk.SummarySetting
//	//读取汇总信息
//	err = cli.db.One("AccountID", "J3wiDj2jMGdp9aqmALhQtEkJQch4YN9e38TEXzRgZyKY", &sumSets)
//	if err != nil {
//		log.Error("GetSummaryInfo error:", err)
//		return
//	}
//	fmt.Printf("SummaryInfo: %+v\n", sumSets)
//}

func TestCLI_GetAllTokenContractBalance(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	accountID := "DCgKWqyefttTqWbyS4ihFsyyvL4jHcF4XBTa3KAGwEmF"
	list, err := cli.GetAllTokenContractBalance(accountID, "")
	if err != nil {
		log.Error("GetAllTokenContractBalance error:", err)
		return
	}
	cli.printTokenContractBalanceList(list, "TRX")
}

func TestCLI_GetAllTokenContractBalanceByAddress(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	accountID := "BBxgBEn7AoRhNqsS7vjD625B5SafFFdY1QMX7Zq8M9jn"
	address := "WhWV4XcD7UJzt2bAcVe48PN1Cxwh8HAyoi"
	list, err := cli.GetAllTokenContractBalanceByAddress(accountID, address, "WICC")
	if err != nil {
		log.Error("GetAllTokenContractBalance error:", err)
		return
	}
	cli.printTokenContractBalanceList(list, "WICC")
}

func TestCLI_printAccountSummaryInfo(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	cli.printAccountSummaryInfo()
}

func TestCLI_GetTokenBalanceByContractAddresss(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	accountID := "DCgKWqyefttTqWbyS4ihFsyyvL4jHcF4XBTa3KAGwEmF"
	address := "THvZvKPLHKLJhEFYKiyqj6j8G8nGgfg7ur"
	account := &openwsdk.Account{Symbol: "TRX", AccountID: accountID}
	balance, err := cli.GetTokenBalanceByContractAddress(account, address)
	if err != nil {
		t.Errorf("GetAllTokenContractBalance error: %v", err)
		return
	}
	fmt.Printf("balance: %+v\n", balance)
}

func TestCLI_UpdateSymbols(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	err := cli.UpdateSymbols()
	if err != nil {
		t.Errorf("UpdateSymbols error: %v", err)
		return
	}
	log.Infof("update info success")
}

func TestCLI_AddTrustAddress(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	trustAddr := openwsdk.NewTrustAddress("WhWV4XcD7UJzt2bAcVe48PN1Cxwh8HAyoi", "WICC2", "testwicc")
	err := cli.AddTrustAddress(trustAddr)
	if err != nil {
		t.Errorf("AddTrustAddress error: %v", err)
		return
	}
	addrs, err := cli.ListTrustAddress("")
	if err != nil {
		t.Errorf("ListTrustAddress error: %v", err)
		return
	}
	cli.printListTrustAddress(addrs)
	cli.printTrustAddressStatus()
}

func TestCLI_ListTrustAddress(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	addrs, err := cli.ListTrustAddress("")
	if err != nil {
		t.Errorf("ListTrustAddress error: %v", err)
		return
	}
	cli.printListTrustAddress(addrs)
	cli.printTrustAddressStatus()
}

func TestCLI_EnableTrustAddress(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	cli.printTrustAddressStatus()
	cli.EnableTrustAddress()
	cli.printTrustAddressStatus()
}


func TestCLI_DisableTrustAddress(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	cli.printTrustAddressStatus()
	cli.DisableTrustAddress()
	cli.printTrustAddressStatus()
}

func TestCLI_IsTrustAddress(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	flag := cli.IsTrustAddress("WhWV4XcD7UJzt2bAcVe48PN1Cxwh8HAyoi22", "WICC")
	log.Infof("WhWV4XcD7UJzt2bAcVe48PN1Cxwh8HAyoi22: %v", flag)

	flag = cli.IsTrustAddress("WhWV4XcD7UJzt2bAcVe48PN1Cxwh8HAyoi", "WICC")
	log.Infof("WhWV4XcD7UJzt2bAcVe48PN1Cxwh8HAyoi: %v", flag)
}