package main

import (
	"github.com/boltdb/bolt"
	"log"
)

var db *bolt.DB

func openDB() {
	var err error
	db, err = bolt.Open("UNKNOWN", 0600, nil)
	if err != nil {
		log.Panic("Cannot create DB")
	}
}

func closeDB() {
	db.Close()
}

func store(bucketName string, key string, value []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			log.Panic("Cannot create new bucket")
		}
		err = bucket.Put([]byte(key), value)
		return err
	})
	if err != nil {
		log.Panic("Cannot update bucket to write")
	}
}

func retrieve(bucketName string, key string) []byte {
    var value []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
        value = bucket.Get([]byte(key))
        return nil
	})
	if err != nil {
		log.Panic("Cannot read from DB")
        return []byte{}
	}
    return value;
}
