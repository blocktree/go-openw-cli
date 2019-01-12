package openwcli

import "github.com/blocktree/OpenWallet/owtp"

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
