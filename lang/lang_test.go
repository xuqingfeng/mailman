package lang

import "testing"

func TestSaveLang(t *testing.T) {

	lang := Lang{"zh"}
	err := SaveLang(lang)
	if err != nil {
		t.Errorf("SaveLang() fail %v", err)
	}
}

func TestGetLang(t *testing.T) {

	lng, err := GetLang()
	if err != nil {
		t.Errorf("GetLang() fail %v", err)
	}
	t.Log(lng)
}
