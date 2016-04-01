package contacts

import "testing"

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

func TestDeleteContacts(t *testing.T) {

	err := DeleteContacts(testContacts.Email)
	if err != nil {
		t.Errorf("DeleteContacts() fail %v", err)
	}
}
