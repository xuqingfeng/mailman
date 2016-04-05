package account

import (
    "github.com/xuqingfeng/mailman/util"
)

type Account struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

//****Account START****
func GetAccountEmail() ([]string, error) {

    boltStore, err := util.NewBoltStore(util.DBPath)
    if err != nil {
        return nil, err
    }
    defer boltStore.Close()

    accounts, err := boltStore.GetRange(util.AccountBucketName)
    if err != nil {
        util.FileLog.Error(err.Error())
        return nil, err
    }
    var ret []string
    for k, _ := range accounts {
        ret = append(ret, k)
    }

    return ret, nil
}
func SaveAccount(account Account) error {

    boltStore, err := util.NewBoltStore(util.DBPath)
    if err != nil {
        return err
    }
    defer boltStore.Close()

    return boltStore.Set([]byte(account.Email), []byte(account.Password), util.AccountBucketName)
}
func GetAccountInfo(email string) (Account, error) {

    boltStore, err := util.NewBoltStore(util.DBPath)
    if err != nil {
        return Account{}, err
    }
    defer boltStore.Close()

    password, err := boltStore.Get([]byte(email), util.AccountBucketName)
    if err != nil {
        util.FileLog.Error(err.Error())
        return Account{}, err
    }

    return Account{email, string(password)}, nil
}
func DeleteAccount(email string) error {

    boltStore, err := util.NewBoltStore(util.DBPath)
    if err != nil {
        return err
    }
    defer boltStore.Close()

    return boltStore.Delete([]byte(email), util.AccountBucketName)
}

//****Account END****
