package util

import (
	"errors"
	"github.com/boltdb/bolt"
)

const (
	dbFileMode = 0600
)

var (
	KVBucketName       = []byte("kv")
	AccountBucketName  = []byte("account")
	ContactsBucketName = []byte("contacts")
	SmtpBucketName     = []byte("smtp")
	testBucketName     = []byte("test")
	KeyNotFoundErr     = errors.New("Key Not Found")
)

type BoltStore struct {
	conn *bolt.DB
	path string
}

func NewBoltStore(path string) (*BoltStore, error) {

	//handle, err := bolt.Open(path, dbFileMode, &bolt.Options{Timeout: 1 * time.Second})
	handle, err := bolt.Open(path, dbFileMode, nil)
	if err != nil {
		FileLog.Error(err.Error())
		return nil, err
	}

	store := &BoltStore{
		handle,
		path,
	}

	if err := store.initialize(); err != nil {
		FileLog.Error(err.Error())
		return nil, err
	}

	return store, nil
}

func (b *BoltStore) initialize() error {

	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	// ?
	defer tx.Rollback()

	// create all buckets
	if _, err = tx.CreateBucketIfNotExists(KVBucketName); err != nil {
		return err
	}
	if _, err = tx.CreateBucketIfNotExists(AccountBucketName); err != nil {
		return err
	}
	if _, err = tx.CreateBucketIfNotExists(ContactsBucketName); err != nil {
		return err
	}
	if _, err = tx.CreateBucketIfNotExists(SmtpBucketName); err != nil {
		return err
	}
	if _, err = tx.CreateBucketIfNotExists(testBucketName); err != nil {
		return err
	}

	return tx.Commit()
}

// defer close
func (b *BoltStore) Close() error {

	return b.conn.Close()
}

func (b *BoltStore) Set(k, v, bucketName []byte) error {

	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(bucketName)
	if err := bucket.Put(k, v); err != nil {
		return err
	}

	return tx.Commit()
}

func (b *BoltStore) Get(k, bucketName []byte) ([]byte, error) {

	tx, err := b.conn.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(bucketName)
	// The returned value is only valid for the life of the transaction.
	val := bucket.Get(k)

	if val == nil {
		return nil, KeyNotFoundErr
	}

	return append([]byte{}, val...), nil
}

func (b *BoltStore) Delete(k, bucketName []byte) error {

	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(bucketName)

	err = bucket.Delete(k)
	if err != nil {
		return err
	}

	// commit needed !
	return tx.Commit()
}

func (b *BoltStore) DeleteBucket(bucketName []byte) error {

	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.DeleteBucket(bucketName)

	return tx.Commit()
}

// range map
// the iteration order is not specified and is not guaranteed to be the same from one iteration to the next.
// https://blog.golang.org/go-maps-in-action
func (b *BoltStore) GetRange(bucketName []byte) (map[string]string, []string, error) {

	tx, err := b.conn.Begin(false)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	ret := make(map[string]string)
	var order []string
	curs := tx.Bucket(bucketName).Cursor()
	// reverse
	for k, v := curs.First(); k != nil; k, v = curs.Next() {
		ret[string(k)] = string(v)
		order = append(order, string(k))
		//FileLog.Warn("k: " + string(k) + " " + "v: " + string(v))
	}
	return ret, order, nil
}
