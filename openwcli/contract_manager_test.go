/*
 * Copyright 2019 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package openwcli

import (
	"github.com/google/uuid"
	"testing"
)

func TestCLI_CallABI(t *testing.T) {

	accountID := "7KgNQFx35ijMA43NgY89uaiwi9Tm4MH1PH68Kpnaqstu"
	contractAddress := "0x550cdb1020046b3115a4f8ccebddfb28b66beb27"
	abiParam := []string{"balanceOf", "0xe6a9cc4fe66e7b726e3e8ef8e32c308ce74c0996"}

	cli := getTestOpenwCLI()
	if cli == nil {
		t.Errorf("cli is not initialized")
		return
	}

	account, err := cli.GetAccountByAccountID(accountID)
	if err != nil {
		t.Errorf("GetAccountByAccountID failed, err: %v", err)
		return
	}

	if account == nil {
		t.Errorf("account is nil")
		return
	}

	_, exErr := cli.CallABI(account, contractAddress, abiParam)
	if exErr != nil {
		t.Errorf("CallABI failed, err: %v", exErr)
		return
	}

}

func TestCLI_TriggerABI(t *testing.T) {

	walletID := "W3LxqTNAcXFqW7HGcTuERRLXKdNWu17Ccx"
	accountID := "7KgNQFx35ijMA43NgY89uaiwi9Tm4MH1PH68Kpnaqstu"
	contractAddress := "0x550cdb1020046b3115a4f8ccebddfb28b66beb27"
	abiParam := []string{"transfer", "0x19a4b5d6ea319a5d5ad1d4cc00a5e2e28cac5ec3", "3456"}

	cli := getTestOpenwCLI()
	if cli == nil {
		t.Errorf("cli is not initialized")
		return
	}

	wallet, err := cli.GetWalletByWalletID(walletID)
	if err != nil {
		t.Errorf("GetWalletByWalletID failed, err: %v", err)
		return
	}

	account, err := cli.GetAccountByAccountID(accountID)
	if err != nil {
		t.Errorf("GetAccountByAccountID failed, err: %v", err)
		return
	}

	sid := uuid.New().String()
	_, exErr := cli.TriggerABI(wallet, account, contractAddress, "", "0", sid, "", "12345678", abiParam, "", 0)
	if exErr != nil {
		t.Errorf("TriggerABI failed, err: %v", exErr)
		return
	}
}
