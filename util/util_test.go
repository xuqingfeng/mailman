package util

import "testing"

func TestGetHomeDir(t *testing.T) {

	homeDir, err := GetHomeDir()
	if err != nil {
		t.Errorf("GetHomeDir() fail %v", err)
	}
	if "" == homeDir {
		t.Error("GetHomeDir() fail")
	}
}

func TestGetUserName(t *testing.T) {

	username := GetUserName()
	if 0 == len(username) {
		t.Error("GetUserName() fail")
	}
}
