package boltAdapter

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/boltdb/bolt"
)

// Connect to the database
func Connect(path string, timeout int) *bolt.DB{
  if db, err = bolt.Open(path, 0600, nil); err != nil {
		log.Println("Bolt Driver Error", err)
		return nil, err
	}
	return db, err
}

func (db *bolt.DB) Close(){

}

func (db *bolt.DB) Update(bucketName string, key string, dataStruct interface{}) error {
	err := db.Update(func(tx *bolt.Tx) error {
		// Create the bucket
		bucket, e := tx.CreateBucketIfNotExists([]byte(bucketName))
		if e != nil {
			return e
		}
		// Encode the record
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

func (db *bolt.DB) View(bucketName string, key string, dataStruct interface{}) error {
	err := BoltDB.View(func(tx *bolt.Tx) error {
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

func (db *bolt.DB) Delete(bucketName string, key string) error {
	err := BoltDB.Update(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		return b.Delete([]byte(key))
	})
	return err
}
