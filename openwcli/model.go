package openwcli

import (
	"github.com/blocktree/OpenWallet/owtp"
	"github.com/blocktree/go-openw-sdk/openwsdk"
)

const (
	CLIBucket          = "CLIBucket"
	CurrentKeychainKey = "current_keychain"
)

//密钥对
type Keychain struct {
	NodeID     string `json:"nodeID" storm:"id"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`

	privateKeyBytes []byte
	publicKeyBytes  []byte
}

//初始化密钥对
func NewKeychain(cert owtp.Certificate) *Keychain {
	priv, pub := cert.KeyPair()
	keychain := &Keychain{
		NodeID:          cert.ID(),
		PrivateKey:      priv,
		PublicKey:       pub,
		privateKeyBytes: cert.PrivateKeyBytes(),
		publicKeyBytes:  cert.PublicKeyBytes(),
	}
	return keychain

}

func (keychain *Keychain) Certificate() (owtp.Certificate, error) {
	return owtp.NewCertificate(keychain.PrivateKey, "")
}

//SummarySetting 汇总设置信息
type SummarySetting struct {
	WalletID        string `json:"walletID"`
	AccountID       string `json:"accountID" storm:"id"`
	SumAddress      string `json:"sumAddress"`
	Threshold       string `json:"threshold"`
	MinTransfer     string `json:"minTransfer"`
	RetainedBalance string `json:"retainedBalance"`
	Confirms        uint64 `json:"confirms"`
}

type SummaryTask struct {
	Wallets []SummaryWalletTask `json:"wallets"`
}

type SummaryAccountTask struct {
	AccountID string   `json:"accountID"`
	Contracts []string `json:"contracts"`
	//account   *openwsdk.Account
}

type SummaryWalletTask struct {
	WalletID string               `json:"walletID"`
	Password string               `json:"password"`
	Accounts []SummaryAccountTask `json:"accounts"`
	wallet   *openwsdk.Wallet
}

/*
{
	"wallets": [
		{
			"walletID": "1234qwer",
			"password": "12345678",
			"accounts": [
				{
					"accountID": "123",
					"contracts":[
						"all", //全部合约
						"3qoe2ll2=", //指定的合约ID
					]
				},
			],
		},
	]
}
*/
