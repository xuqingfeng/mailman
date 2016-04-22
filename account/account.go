package account

import (
	"encoding/base64"
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

	_, order, err := boltStore.GetRange(util.AccountBucketName)
	if err != nil {
		util.FileLog.Error(err.Error())
		return nil, err
	}
	var ret []string
	for _, v := range order {
		ret = append(ret, v)
	}

	return ret, nil
}
func SaveAccount(account Account) error {

	boltStore, err := util.NewBoltStore(util.DBPath)
	if err != nil {
		return err
	}
	defer boltStore.Close()
	account.Password = string(encryptPassword([]byte(account.Password)))

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
	password = decryptPassword(password)

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

func encryptPassword(password []byte) []byte {

	return []byte(base64.URLEncoding.EncodeToString(password))
}
func decryptPassword(encryptedPassword []byte) []byte {

	decodedStr, _ := base64.URLEncoding.DecodeString(string(encryptedPassword))
	return decodedStr
}

//****Account END****
