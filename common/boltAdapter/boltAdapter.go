package boltAdapter

import (
	"encoding/json"
	"log"
	"github.com/boltdb/bolt"
)

type BoltConnection struct{
	Database *bolt.DB
	LastError error
}

var (
	Connection BoltConnection
)

// Connect to the database
func Connect(path string, timeout int) BoltConnection{
  db, err := bolt.Open(path, 0600, nil)
	if err != nil {
	  log.Fatalln("Bolt Driver Error", err)
		Connection.Database = nil
		Connection.LastError = err
		return Connection
	}
	Connection.Database = db
	return Connection
}

func (db BoltConnection) Close(){
  db.Database.Close()
}

func (db BoltConnection) Update(bucketName string, key string, dataStruct interface{}) error {
	err := db.Database.Update(func(tx *bolt.Tx) error {
		// Create the bucket
		bucket, e := tx.CreateBucketIfNotExists([]byte(bucketName))
		if e != nil {
			return e
		}
		encodedRecord, e := json.Marshal(dataStruct)
		if e != nil {
			return e
		}
		if e = bucket.Put([]byte(key), encodedRecord); e != nil {
			return e
		}
		return nil
	})
	return err
}

func (db BoltConnection) View(bucketName string, key string, dataStruct interface{}) error {
	err := db.Database.View(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		v := b.Get([]byte(key))
		if len(v) < 1 {
			return bolt.ErrInvalid
		}
		e := json.Unmarshal(v, &dataStruct)
		if e != nil {
			return e
		}
		return nil
	})
	return err
}

func (db BoltConnection) Delete(bucketName string, key string) error {
	err := db.Database.Update(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		return b.Delete([]byte(key))
	})
	return err
}
