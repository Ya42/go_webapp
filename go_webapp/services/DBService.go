package service

import(
  "fmt"
  "time"
  "github.com/bolt"
  "encoding/json"
)

var db *bolt.DB
var err error

func OpenDB(dbpath string, timeout int) *bolt.DB{
  //Set to -1 to optout timeout
  if int = -1{
    db, err = bolt.Open(dbpath, 0600, nil)
  }else{
    db, err = bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 1 * timeout.second})
  }
  if err != nil{
    fmt.Println(err)
    return nil
  }
  return db
}

func Update(db *bolt.DB, bucketName string, key string, dataStruct interface{}) (byte[], error) {
	err := db.Update(func(tx *bolt.Tx) error {
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
	return encodedRecord, err
}

func View(db *bolt.DB, bucketName string, key string, dataStruct interface{}) (byte[], error) {
	err := db.View(func(tx *bolt.Tx) error {
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
	return v, err
}

func Delete(db *bolt.DB, bucketName string, key string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		return b.Delete([]byte(key))
	})
	return err
}

func CloseDB(db *bolt.DB) error{
  err = db.Close()
  if err != nil{
    fmt.Println("DB failed to close")
  }
}
