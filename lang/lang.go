package lang

import (
	"github.com/xuqingfeng/mailman/util"
)

type Lang struct {
	Type string `json:"type"`
}

func GetLang() (string, error) {

	boltStore, err := util.NewBoltStore(util.DBPath)
	if err != nil {
		return "", err
	}
	defer boltStore.Close()

	t, err := boltStore.Get([]byte(util.DefaultLang), util.KVBucketName)
	if err != nil {
		util.FileLog.Error(err.Error())
		return "", err
	}

	return string(t), nil
}

func SaveLang(lang Lang) error {

	boltStore, err := util.NewBoltStore(util.DBPath)
	if err != nil {
		return err
	}
	defer boltStore.Close()

	return boltStore.Set([]byte(util.DefaultLang), []byte(lang.Type), util.KVBucketName)
}
