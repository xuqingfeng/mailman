package util

import (
	"github.com/boltdb/bolt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

func TestNewBoltStore(t *testing.T) {

	boltStore, err := NewBoltStore(DBPath)
	defer boltStore.Close()
	if err != nil {
		t.Errorf("NewBoltStore() fail %v", err)
	}
	if boltStore == nil {
		t.Error("NewBoltStore() fail")
	}

	fakePath := "/tmp/fakeDir/fake.db"
	_, err = NewBoltStore(fakePath)
	if err == nil {
		t.Error("db path doesn't exist NewBoltStore() should fail")
	}
}

func TestSetGet(t *testing.T) {

	k, v := "testKey", "testValue"
	boltStore, _ := NewBoltStore(DBPath)
	defer boltStore.Close()
	err := boltStore.Set([]byte(k), []byte(v), testBucketName)
	if err != nil {
		t.Errorf("Set() fail %v", err)
	}
	valueInBolt, err := boltStore.Get([]byte(k), testBucketName)
	if err != nil {
		t.Errorf("Get() fail %v", err)
	}
	if valueInBolt == nil {
		t.Error("Get() fail")
	}
}

func TestDelete(t *testing.T) {

	k := "testKey"
	boltStore, _ := NewBoltStore(DBPath)
	defer boltStore.Close()
	err := boltStore.Delete([]byte(k), testBucketName)
	if err != nil {
		t.Errorf("Delete() fail %v", err)
	}
	valueToBeDeleted, _ := boltStore.Get([]byte(k), testBucketName)
	if valueToBeDeleted != nil {
		t.Error("Delete() fail")
	}
}

func TestGetRange(t *testing.T) {

	boltStore, _ := NewBoltStore(DBPath)
	defer boltStore.Close()
	kvGroup := map[string]string{
		"testKey0": "testValue0",
		"testKey1": "testValue1",
	}
	for k, v := range kvGroup {
		boltStore.Set([]byte(k), []byte(v), testBucketName)
	}
	kvGroupInBolt, _, err := boltStore.GetRange(testBucketName)
	if err != nil {
		t.Errorf("GetRange() fail %v", err)
	}
	if !reflect.DeepEqual(kvGroup, kvGroupInBolt) {
		t.Error("GetRange() fail")
	}
}

func TestDeleteBucket(t *testing.T) {

	boltStore, _ := NewBoltStore(DBPath)
	defer boltStore.Close()
	err := boltStore.DeleteBucket(testBucketName)
	if err != nil {
		t.Errorf("DeleteBucket() fail %v", err)
	}
}

func TestMultiOpenBolt(t *testing.T) {
	_, err := bolt.Open("testdata/test.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Errorf("bolt open fail %v", err)
	}

	_, err = bolt.Open("testdata/test.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err == nil {
		t.Errorf("bolt reopen should fail  %v", err)
	}
}
