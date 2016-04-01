package contacts

import (
	"github.com/xuqingfeng/mailman/util"
)

type Contacts struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

//****Contacts START****
func GetContacts() ([]Contacts, error) {

	boltStore, err := util.NewBoltStore(util.DBPath)
	if err != nil {
		return nil, err
	}
	defer boltStore.Close()

	contacts, err := boltStore.GetRange(util.ContactsBucketName)
	if err != nil {
		util.FileLog.Error(err.Error())
		return nil, err
	}
	var contactsList []Contacts
	for k, v := range contacts {
		contactsList = append(contactsList, Contacts{k, v})
	}

	return contactsList, nil
}
func SaveContacts(contacts Contacts) error {

	boltStore, err := util.NewBoltStore(util.DBPath)
	if err != nil {
		return err
	}
	defer boltStore.Close()

	return boltStore.Set([]byte(contacts.Email), []byte(contacts.Name), util.ContactsBucketName)
}
func DeleteContacts(email string) error {

	boltStore, err := util.NewBoltStore(util.DBPath)
	if err != nil {
		return err
	}
	defer boltStore.Close()

	return boltStore.Delete([]byte(email), util.ContactsBucketName)
}

//****Contacts END****
