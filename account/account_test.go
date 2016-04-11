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
    test1Account = Account{
        "test1@example.com",
        "password1",
    }
)

func TestSaveAccount(t *testing.T) {

    err := SaveAccount(testAccount)
    if err != nil {
        t.Errorf("SaveAccount() fail %v", err)
    }
    SaveAccount(test1Account)
}
func TestGetAccountEmail(t *testing.T) {

    accounts, err := GetAccountEmail()
    t.Logf("%#v", accounts)
    if err != nil {
        t.Errorf("GetAccountEmail() fail %v", err)
    }
    if len(accounts) < 1 {
        t.Error("GetAccountEmail() fail")
    }
}

func TestGetAccountInfo(t *testing.T) {

    accountToReturn, err := GetAccountInfo(testAccount.Email)
    if err != nil {
        t.Errorf("GetAccountInfo() fail %v", err)
    }
    if !reflect.DeepEqual(testAccount, accountToReturn) {
        t.Error("GetAccountInfo() fail")
    }
}

func TestDeleteAccount(t *testing.T) {

    err := DeleteAccount("test@example.com")
    if err != nil {
        t.Errorf("DeleteAccount() fail %v", err)
    }
    DeleteAccount("test1@example.com")
}

func TestEncryptPassword(t *testing.T) {

    password := []byte("fakePassword")
    encryptedPassword := encryptPassword(password)
    if string(password) == string(encryptedPassword) {
        t.Error("encryptPassword fail")
    }
    decryptedPassword := decryptPassword(encryptedPassword)
    if string(password) != string(decryptedPassword) {
        t.Error("decryptPassword fail")
    }
}
