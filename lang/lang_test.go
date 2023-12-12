package lang

import (
	"testing"

	"mailman/util"
)

func TestSaveLang(t *testing.T) {

	lang := Lang{"zh"}
	err := SaveLang(lang)
	if err != nil {
		t.Errorf("SaveLang() fail %v", err)
	}
	// fail
	fakePath := "/tmp/fakeDir/fake.db"
	util.DBPath, fakePath = fakePath, util.DBPath
	err = SaveLang(lang)
	if err == nil {
		t.Error("db file doesn't exist SaveLog() should fail")
	}
	util.DBPath, fakePath = fakePath, util.DBPath
}

func TestGetLang(t *testing.T) {

	lng, err := GetLang()
	if err != nil {
		t.Errorf("GetLang() fail %v", err)
	}
	t.Log(lng)

	// fail
	fakePath := "/tmp/fakeDir/fake.db"
	util.DBPath, fakePath = fakePath, util.DBPath
	_, err = GetLang()
	if err == nil {
		t.Error("db file doesn't exist GetLog() should fail")
	}
	util.DBPath, fakePath = fakePath, util.DBPath
	fakeLang := "fakeLang"
	util.DefaultLang, fakeLang = fakeLang, util.DefaultLang
	_, err = GetLang()
	if err == nil {
		t.Error("lang doesn't exist GetLog() should fail")
	}
	util.DefaultLang, fakeLang = fakeLang, util.DefaultLang
}
