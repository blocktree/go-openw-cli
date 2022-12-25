package openwcli

import (
	"encoding/json"
	"github.com/blocktree/go-openw-sdk/v2/openwsdk"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/google/uuid"
	"testing"
	"time"
)

func testFindAccountByID(accountID string, list []*openwsdk.Account) *openwsdk.Account {
	for _, a := range list {
		if a.AccountID == accountID {
			return a
		}
	}
	return nil
}

func testCLITransfer(walletID, accountID, symbol, contractAddress, amount, to, memo, password string) *openwallet.Error {
	cli := getTestOpenwCLI()
	if cli == nil {
		return openwallet.Errorf(openwallet.ErrUnknownException, "init cli error")
	}

	wallet, err := cli.GetWalletByWalletID(walletID)
	if err != nil {
		//log.Error("GetWalletByWalletID error:", err)
		return openwallet.ConvertError(err)
	}

	account, err := cli.GetAccountByAccountID(symbol, accountID)
	if err != nil {
		//log.Error("GetAccountByAccountID error:", err)
		return openwallet.ConvertError(err)
	}

	if account != nil {
		sid := uuid.New().String()
		_, _, exErr := cli.Transfer(wallet, account, symbol, contractAddress, to, amount, sid, "", memo, password)
		if exErr != nil {
			//log.Error("Transfer error code: %d, msg: %s", exErr.Code(), exErr.Error())
			return exErr
		}
	}

	return nil
}

func testCLITransferAll(walletID, accountID, symbol, to, password string) *openwallet.Error {
	cli := getTestOpenwCLI()
	if cli == nil {
		return openwallet.Errorf(openwallet.ErrUnknownException, "init cli error")
	}

	wallet, err := cli.GetWalletByWalletID(walletID)
	if err != nil {
		//log.Error("GetWalletByWalletID error:", err)
		return openwallet.ConvertError(err)
	}

	account, err := cli.GetAccountByAccountID(symbol, accountID)
	if err != nil {
		//log.Error("GetAccountByAccountID error:", err)
		return openwallet.ConvertError(err)
	}

	if account != nil {
		sid := uuid.New().String()
		exErr := cli.TransferAll(wallet, account, "MATIC", "", to, sid, "", "", password)
		if exErr != nil {
			//log.Error("Transfer error code: %d, msg: %s", exErr.Code(), exErr.Error())
			return openwallet.ConvertError(err)
		}
	}

	return nil
}

func TestCLI_Transfer_LTC(t *testing.T) {
	walletID := "W3LxqTNAcXFqW7HGcTuERRLXKdNWu17Ccx"
	accountID := "PgHCcfMbcw1zXRNZo23NFjRdBmcN5tzrb1j5McRLJbG"
	amount := "0.001"
	to := "LcaFc1pmJBsS7MQyMvZaboppuuvGFubD49"
	password := "12345678"
	err := testCLITransfer(walletID, accountID, "ETH", "", amount, to, "", password)
	if err != nil {
		t.Errorf("Transfer error code: %d, msg: %s", err.Code(), err.Error())
		return
	}
}

func TestCLI_TransferAll_Coin(t *testing.T) {
	walletID := "W4W8i2F27c1YB63rMZswSquVSnS3265MAF"
	accountID := "3BSJAseva4A2oZgmEuMEVtdSgnj5UXCunYzyWtK7dj4b"
	to := "0x8C178b782fab1d0686D88bC16B31F80431098fa1"
	password := "1234qwer"
	err := testCLITransferAll(walletID, accountID, "MATIC", to, password)
	if err != nil {
		t.Errorf("Transfer All error code: %d, msg: %s", err.Code(), err.Error())
		return
	}
}

func TestCLI_Summary_BTC(t *testing.T) {

	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	//accountID := "7u7CQNdkaJXVszoj528Bink88aWgfay3rDxb1rsmDywA"

	plain := `

{
    "wallets": [
        {
            "walletID": "W6iq5mLbyCdXKmUDMZzFnDx4rFw8eTYddS",
            "password": "1234qwer",
            "accounts": [
                {
                    "accountID": "7u7CQNdkaJXVszoj528Bink88aWgfay3rDxb1rsmDywA",
                    "onlyContracts": true,
                    "threshold": "1000",
                    "minTransfer": "0",
                    "retainedBalance": "0",
                    "switchSymbol": "ETH",
                    "contracts": {
                        "all": {
                            "threshold": "1000",
                            "minTransfer": "0",
                            "retainedBalance": "0"
                        }
                    }
                }
            ]
        }
    ]
}

`
	var summaryTask openwsdk.SummaryTask
	err := json.Unmarshal([]byte(plain), &summaryTask)
	if err != nil {
		log.Error("json.Unmarshal error:", err)
		return
	}

	cli.summaryTask = &summaryTask

	//err = cli.SetSummaryInfo(&openwsdk.SummarySetting{
	//	"W1ixmQVGWX78MnFacrZ38i8kAivVifa2d7",
	//	accountID,
	//	"TE7nDCrQFPQTfRzsXRBFCLmVHDbiiw6BW9",
	//	"10000",
	//	"0",
	//	"0",
	//	0,
	//})
	//if err != nil {
	//	log.Error("SetSummaryInfo error:", err)
	//	return
	//}

	err = cli.checkSummaryTaskIsHaveSettings(&summaryTask)
	if err != nil {
		log.Error("checkSummaryTaskIsHaveSettings error:", err)
		return
	}

	cli.SummaryTask()

	//wallets, err := cli.GetWalletsOnServer()
	//if err != nil {
	//	log.Error("GetWalletsOnServer error:", err)
	//	return
	//}
	//
	//if len(wallets) > 0 {
	//	accounts, err := cli.GetAccountsOnServer(wallets[0].WalletID)
	//	if err != nil {
	//		log.Error("GetAccountsOnServer error:", err)
	//		return
	//	}
	//
	//	account := testFindAccountByID(accountID, accounts)
	//
	//	if account != nil {
	//		cli.SummaryTask()
	//	}
	//}
}

func TestCLI_GetSummaryTaskLog(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}
	logs, err := cli.GetSummaryTaskLog(0, 20)
	if err != nil {
		log.Error("GetAccountsOnServer error:", err)
		return
	}
	for i, l := range logs {
		log.Infof("log[%d]: %+v", i, l)
	}
}

func TestCLI_Transfer_Token(t *testing.T) {
	count := 1000
	for i := 0; i < count; i++ {
		walletID := "WCpVSv7AsTpLpkkc5tHApnqCKzdsoNKr8P"
		accountID := "72xtVNtkkiJyEHoFCSL9XjnwRQNjVm1GBpLwLG6Rk98h"
		amount := "0.0001"
		contractAddress := "evsio.token:TGC"
		to := "tgcopenwtest"
		password := "12345678"
		memo := "N3CCXUQL"
		err := testCLITransfer(walletID, accountID, "ETH", contractAddress, amount, to, memo, password)
		if err != nil {
			t.Errorf("Transfer error code: %d, msg: %s \n", err.Code(), err.Error())
			return
		}
		time.Sleep(500 * time.Millisecond)
	}

}

func TestCLI_Transfer_MATIC(t *testing.T) {
	walletID := "W4W8i2F27c1YB63rMZswSquVSnS3265MAF"
	accountID := "3BSJAseva4A2oZgmEuMEVtdSgnj5UXCunYzyWtK7dj4b"
	amount := "0.01"
	to := "0x8C178b782fab1d0686D88bC16B31F80431098fa1"
	password := "1234qwer"
	err := testCLITransfer(walletID, accountID, "MATIC", "", amount, to, "", password)
	if err != nil {
		t.Errorf("Transfer error code: %d, msg: %s", err.Code(), err.Error())
		return
	}
}
