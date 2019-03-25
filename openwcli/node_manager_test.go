package openwcli

import (
	"fmt"
	"github.com/blocktree/openwallet/log"
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
	if cli == nil {
		return
	}

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