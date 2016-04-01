package account

import (
	"reflect"
	"testing"
)

var (
	testAccount = Account{
		"test@example.com",
		"password",
	}
)

func TestSaveAccount(t *testing.T) {

	err := SaveAccount(testAccount)
	if err != nil {
		t.Errorf("SaveAccount() fail %v", err)
	}
}

func TestGetAccountInfo(t *testing.T) {

	accountToReturn, err := GetAccountInfo(testAccount.Email)
	if err != nil {
		t.Errorf("GetAccountInfo() fail %v", err)
	}
	if reflect.DeepEqual(testAccount, accountToReturn) {
		t.Logf("GetAccountInfo() success")
	}
}

func TestDeleteAccount(t *testing.T) {

	err := DeleteAccount("test@example.com")
	if err != nil {
		t.Errorf("DeleteAccount() fail %v", err)
	}
}
