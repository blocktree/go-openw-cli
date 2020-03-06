package openwcli

import (
	"encoding/hex"
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/log"
	"testing"
)

func TestCLI_GenKeychain(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	keychain, err := GenKeychain()
	if err != nil {
		log.Error("GenKeychain error:", err)
		return
	}
	printKeychain(keychain)
}

func TestCLI_RegisterOnServer(t *testing.T) {
	cli := getTestOpenwCLI()
	//if cli == nil {
	//	return
	//}

	//生成keychain
	keychain, err := GenKeychain()
	if err != nil {
		log.Error("RegisterOnServer error:", err)
		return
	}

	err = cli.SaveCurrentKeychain(keychain)
	if err != nil {
		log.Error("RegisterOnServer error:", err)
		return
	}

	//配置APISDK
	err = cli.setupAPISDK(keychain)
	if err != nil {
		log.Error("RegisterOnServer error:", err)
		return
	}

	err = cli.RegisterOnServer()
	if err != nil {
		log.Error("RegisterOnServer error:", err)
		return
	}
}

func TestDelArrayObj(t *testing.T) {
	a := []int{0, 1, 2, 3, 4}
	fmt.Println(a)
	//删除第i个元素
	i := 2
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a)
}

func TestSha256(t *testing.T) {
	msg, _ := hex.DecodeString("00e0008403d67925c8c7fda675b4bf8e3230d2fccafd9c32be6414059bc3aa4bbb87d885484e447439716e41486e4641755038543947627a51326f385561616351736341635532120000000000000068656c6c6f20626f79")
	hash := owcrypt.Hash(msg, 0, owcrypt.HASH_ALG_SHA256)
	log.Infof("msg: %s", hex.EncodeToString(msg))
	log.Infof("hash: %s", hex.EncodeToString(hash))
}
