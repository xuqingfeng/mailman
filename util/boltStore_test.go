package util

import (
	"os"
	"reflect"
	"testing"
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
	if boltStore != nil {
		t.Log("NewBoltStore() success")
	}
}

func TestSetGet(t *testing.T) {

	k, v := "testKey", "testValue"
	boltStore, _ := NewBoltStore(DBPath)
	defer boltStore.Close()
	err := boltStore.Set([]byte(k), []byte(v), kvBucketName)
	if err != nil {
		t.Errorf("Set() fail %v", err)
	}
	valueInBolt, err := boltStore.Get([]byte(k), kvBucketName)
	if err != nil {
		t.Errorf("Get() fail %v", err)
	}
	if valueInBolt != nil {
		t.Logf("Get() success %s", valueInBolt)
	}
}

func TestDelete(t *testing.T) {

	k := "testKey"
	boltStore, _ := NewBoltStore(DBPath)
	defer boltStore.Close()
	err := boltStore.Delete([]byte(k), kvBucketName)
	if err != nil {
		t.Errorf("Delete() fail %v", err)
	}
	valueToBeDeleted, _ := boltStore.Get([]byte(k), kvBucketName)
	if valueToBeDeleted == nil {
		t.Log("Delete() success")
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
		boltStore.Set([]byte(k), []byte(v), kvBucketName)
	}
	kvGroupInBolt, err := boltStore.GetRange(kvBucketName)
	if err != nil {
		t.Errorf("GetRange() fail %v", err)
	}
	if reflect.DeepEqual(kvGroup, kvGroupInBolt) {
		t.Log("GetRange() success")
	}
}

func TestDeleteBucket(t *testing.T) {

	boltStore, _ := NewBoltStore(DBPath)
	defer boltStore.Close()
	err := boltStore.DeleteBucket(kvBucketName)
	if err != nil {
		t.Errorf("DeleteBucket() fail %v", err)
	}
}
