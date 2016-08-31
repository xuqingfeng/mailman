package util

import "testing"

func TestGetHomeDir(t *testing.T) {

	homeDir := GetHomeDir()
	if "" == homeDir {
		t.Error("GetHomeDir() fail")
	}
}

func TestGetUser(t *testing.T) {

	username := GetUser()
	if 0 == len(username) {
		t.Error("GetUser() fail")
	}
}
