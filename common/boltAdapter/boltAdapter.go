package boltAdapter

import (
	"encoding/json"
	"log"
	"bytes"
	"github.com/boltdb/bolt"
	"fmt"
)

type BoltConnection struct{
	Database *bolt.DB
	LastError error
}

var (
	Connection BoltConnection
)

// Connect to the database
func Connect(path string) (BoltConnection){
  db, err := bolt.Open("gowebapp.db", 0600, nil)
	if err != nil {
	  log.Fatalln(err)
	}
  Connection.Database = db
	Connection.LastError = err
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
		fmt.Println(bucketName)
		encodedRecord, e := json.Marshal(dataStruct)
		if e != nil {
			return e
		}
		e = bucket.Put([]byte(key), encodedRecord)
		fmt.Println("commited")
		if e != nil {
			return e
		}
		return nil
	})
	return err
}

func (db BoltConnection) View(bucketName string, key string, dataStruct interface{}) error {
	err := db.Database.View(func(tx *bolt.Tx) error {
		// Get the bucket
		var e error
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		v := b.Get([]byte(key))
		if len(v) < 1 {
			return bolt.ErrInvalid
		}
		e = json.Unmarshal(v, &dataStruct)
		return e
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

func (db BoltConnection) Seek(bucketName string, term string) ([]string, error){
	result := []string{}
	var err error
	err = db.Database.View(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		c := b.Cursor()
		prefix := []byte(term)
		for k,_ := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			result = append(result,string(k[:]))
		}
		return nil
	})
	return result,err
}
