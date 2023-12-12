package account

import (
	"crypto/aes"
	"crypto/cipher"

	"mailman/util"
)

const (
	defaultKeyLength = 16
)

var (
	commonIV   = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	defaultKey = []byte("bQKR01UFB33NECJy")
)

type Account struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
	encryptedPassword, err := encryptPassword([]byte(account.Password))
	if err != nil {
		util.FileLog.Error(err)
	}
	account.Password = string(encryptedPassword)

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
	password, err = decryptPassword(password)
	if err != nil {
		util.FileLog.Error(err)
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

func encryptPassword(password []byte) ([]byte, error) {

	key := generateKey()
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	encryptedPassword := make([]byte, len(password))
	cfb.XORKeyStream(encryptedPassword, password)

	return encryptedPassword, nil
}
func decryptPassword(encryptedPassword []byte) ([]byte, error) {

	key := generateKey()
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBDecrypter(c, commonIV)
	decryptedPassword := make([]byte, len(encryptedPassword))
	cfb.XORKeyStream(decryptedPassword, encryptedPassword)

	return decryptedPassword, nil
}

func generateKey() []byte {
	var k string
	username := util.GetUser()
	if len(username) > defaultKeyLength {
		k = username[:defaultKeyLength]
	} else {
		left := defaultKeyLength - len(username)
		k = username + string(defaultKey[:left])
	}

	return []byte(k)
}
