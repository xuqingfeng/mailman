package contacts

import (
	"reflect"
	"testing"
)

var (
	testContacts = Contacts{
		"test@example.com",
		"test",
	}
)

func TestSaveContacts(t *testing.T) {

	err := SaveContacts(testContacts)
	if err != nil {
		t.Errorf("SaveContacts() fail %v", err)
	}
}

func TestGetContacts(t *testing.T) {

	c, err := GetContacts()
	if err != nil {
		t.Errorf("GetContacts() fail %v", err)
	}
	contactsList := []Contacts{{"test@example.com", "test"}}
	if !reflect.DeepEqual(contactsList, c) {
		t.Error("GetContacts() fail")
	}
}

func TestDeleteContacts(t *testing.T) {

	err := DeleteContacts(testContacts.Email)
	if err != nil {
		t.Errorf("DeleteContacts() fail %v", err)
	}
}
