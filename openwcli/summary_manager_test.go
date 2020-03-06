package openwcli

import (
	"encoding/json"
	"github.com/blocktree/go-openw-sdk/v2/openwsdk"
	"github.com/blocktree/openwallet/v2/log"
	"testing"
)

func TestCLI_appendSummaryTasks(t *testing.T) {

	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	plain := `

{
    "wallets": [
        {
            "walletID": "W1ixmQVGWX78MnFacrZ38i8kAivVifa2d7",
            "password": "",
            "accounts": [
                {
                    "accountID": "DCgKWqyefttTqWbyS4ihFsyyvL4jHcF4XBTa3KAGwEmF",
                    "onlyContracts": true,
                    "threshold": "5",
                    "minTransfer": "2",
                    "retainedBalance": "0",
					"addressLimit": 20,
                    "contracts": {
                        "1002000": {
                            "threshold": "30",
                            "minTransfer": "20",
                            "retainedBalance": "0"
                        },
                        "THvZvKPLHKLJhEFYKiyqj6j8G8nGgfg7ur": {
                            "threshold": "10",
                            "minTransfer": "0",
                            "retainedBalance": "0"
                        }
                    },
                    "feesSupportAccount": {         
                        "accountID": "Bp38EW8An9DnNah7pYCYsRur24VMXXifgyBj9mHCtJ17",       
                        "lowBalanceWarning": "10",
                        "fixSupportAmount": "2",
                        "feesScale": "2"
                    }
                }
            ]
        },
		{
            "walletID": "W22222222",
            "password": "",
            "accounts": [
                {
                    "accountID": "234234234",
                    "threshold": "5",
                    "minTransfer": "0",
                    "retainedBalance": "0"
                }
            ]
        }
    ]
}

`
	var summaryTask openwsdk.SummaryTask
	err := json.Unmarshal([]byte(plain), &summaryTask)
	if err != nil {
		t.Errorf("json.Unmarshal error: %v", err)
		return
	}

	cli.summaryTask = &summaryTask

	log.Infof("--------------- before append task ---------------")

	for _, w := range cli.summaryTask.Wallets {
		log.Infof("walletID: %s", w.WalletID)
		for _, a := range w.Accounts {
			log.Infof("accountID: %s", a.AccountID)
			for addr, _ := range a.Contracts {
				log.Infof("contractAddr: %s", addr)
			}
			log.Infof("addressLimit: %d", a.AddressLimit)
		}
	}

	append := `

{
    "wallets": [
        {
            "walletID": "W1ixmQVGWX78MnFacrZ38i8kAivVifa2d7",
            "password": "",
            "accounts": [
                {
                    "accountID": "DCgKWqyefttTqWbyS4ihFsyyvL4jHcF4XBTa3KAGwEmF",
                    "onlyContracts": true,
                    "threshold": "5",
                    "minTransfer": "2",
                    "retainedBalance": "0",
                    "contracts": {
                        "1002000": {
                            "threshold": "30",
                            "minTransfer": "20",
                            "retainedBalance": "0"
                        },
                        "123023": {
                            "threshold": "10",
                            "minTransfer": "0",
                            "retainedBalance": "0"
                        }
                    },
                    "feesSupportAccount": {         
                        "accountID": "Bp38EW8An9DnNah7pYCYsRur24VMXXifgyBj9mHCtJ17",       
                        "lowBalanceWarning": "10",
                        "fixSupportAmount": "2",
                        "feesScale": "2"
                    }
                },
				{
                    "accountID": "1111122222",
                    "threshold": "5",
                    "minTransfer": "2",
                    "retainedBalance": "0"
                }
            ]
        }
    ]
}

`

	var apppendTask openwsdk.SummaryTask
	err = json.Unmarshal([]byte(append), &apppendTask)
	if err != nil {
		t.Errorf("json.Unmarshal error: %v", err)
		return
	}

	cli.appendSummaryTasks(&apppendTask)

	log.Infof("--------------- after append task ---------------")

	for _, w := range cli.summaryTask.Wallets {
		log.Infof("walletID: %s", w.WalletID)
		for _, a := range w.Accounts {
			log.Infof("accountID: %s", a.AccountID)
			for addr, _ := range a.Contracts {
				log.Infof("contractAddr: %s", addr)
			}
		}
	}

}
