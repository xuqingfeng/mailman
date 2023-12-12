package util

import (
	"testing"

	_ "mailman/statik"
)

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

func TestGetContentFromStatik(t *testing.T) {

	_, err := GetContentFromStatik("/responsive.html")
	if err != nil {
		t.Errorf("E! GetContentFromStatik() failed: %v", err)
	}
}
