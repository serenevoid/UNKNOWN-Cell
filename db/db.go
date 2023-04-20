package db

import (
	"log"
	"math/rand"
	"time"

	"github.com/boltdb/bolt"
)

/* ---- BOLT DB FUNCTIONS ---- */
var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open("UNKNOWN.db", 0600, nil)
	if err != nil {
		log.Panic("Cannot create DB")
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Channels"))
		if err != nil {
			log.Panic("Cannot create bucket Channels")
		}
		_, err = tx.CreateBucketIfNotExists([]byte("Guilds"))
		if err != nil {
			log.Panic("Cannot create new bucket")
		}
		return err
	})
	if err != nil {
		log.Panic("Cannot update bucket to write")
	}
}

func CloseDB() {
	db.Close()
}

func InsertDataToBucket(bucketName string, key string, value []byte) {
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

func DeleteDataFromBucket(bucketName string, key string) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		err := bucket.Delete([]byte(key))
		if err != nil {
			log.Panic("Cannot delete key from DB")
		}
		return nil
	})
	if err != nil {
		log.Panic("Cannot open DB todelete key")
	}
}

func IsKeyPresentInBucket(bucketName string, key string) bool {
	isPresent := false
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket.Get([]byte(key)) != nil {
			isPresent = true
		}
		return nil
	})
	if err != nil {
		log.Panic("Cannot read from DB")
		return false
	}
	return isPresent
}

func GetRandomSubscribers(ring func(string)) {
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Channels"))
		if bucket == nil {
			log.Panic("Bucket not found")
		}

		// Count the number of keys in the bucket
		count := 0
		bucket.ForEach(func(_, _ []byte) error {
			count++
			return nil
		})

		// Select 9 random keys
		rand.Seed(time.Now().UnixNano())
		selectedKeys := make(map[string]bool)
		for len(selectedKeys) < 9 && count > 0 {
			// Generate a random key index
			index := rand.Intn(count)

			// Iterate over the keys and select the key at the random index
			i := 0
			bucket.ForEach(func(k, _ []byte) error {
				if i == index {
					selectedKeys[string(k)] = true
				}
				i++
				return nil
			})

			// Decrement the count
			count--
		}

		// Ring the selected keys
		for key := range selectedKeys {
			ring(key)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func GetKeyCount(bucketName string) int {
    count := 0
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			log.Panic("Bucket not found")
		}

		// Count the number of keys in the bucket
		bucket.ForEach(func(_, _ []byte) error {
			count++
			return nil
		})
		return nil
	})
	if err != nil {
		log.Panic("Cannot count Subscribers")
	}
    return count
}

/* ---- IN MEMORY DB FUNCTIONS ---- */
var (
	waitList       = make([]string, 0)
	connectionMap  = make(map[string]string)
	channelUserMap = make(map[string]string)
	tempUserMap    = make(map[string][]string)
	reportList     = make(map[string]int)
	banList        = make(map[string]int)
)

func PushWaitList(channelID string) {
	waitList = append(waitList, channelID)
}

func PopWaitList() string {
	if len(waitList) > 0 {
		element := waitList[0]
		waitList = waitList[1:]
		return element
	}
	return ""
}

func RemoveWaitList(index int) {
	waitList = append(waitList[0:index], waitList[index+1:]...)
}

func AddConnection(user1 string, user2 string) {
	connectionMap[user1] = user2
	connectionMap[user2] = user1
}

func ViewConnection(user string) string {
	return connectionMap[user]
}

func GetConnectionCount() int {
	return len(connectionMap)
}

func RemoveConnection(user1 string, user2 string) {
	delete(connectionMap, user1)
	delete(connectionMap, user2)
}

func AddChannelUser(channelID string, userID string) {
	channelUserMap[channelID] = userID
}

func ViewChannelUser(channelID string) string {
	return channelUserMap[channelID]
}

func RemoveChannelUser(channelID string) {
	delete(channelUserMap, channelID)
}

func addTempUsers(channelID string, userID string) {
	tempUserMap[channelID] = append(tempUserMap[channelID], userID)
}

func GetTempUserIndex(channelID string, userID string) int {
	for index, value := range tempUserMap[channelID] {
		if value == userID {
			return index + 1
		}
	}
	addTempUsers(channelID, userID)
	return len(tempUserMap[channelID])
}

func RemoveTempUsers(channelID string) {
	delete(tempUserMap, channelID)
}

func ReportUser(userID string) {
	currentCount := reportList[userID] / 1000
	today := time.Now().YearDay()
	reportList[userID] = (currentCount+1)*1000 + today
	if currentCount > 4 {
		BanUser(userID, today)
	}
}

func BanUser(userID string, today int) {
	banList[userID] = today
}

func IsBanned(userID string) bool {
	if banList[userID] < time.Now().YearDay() {
		delete(banList, userID)
		if gap := reportList[userID] % 1000; time.Now().YearDay()-gap > 5 {
			delete(reportList, userID)
		} else {
			// Windowed report count reduction
			reportList[userID] = reportList[userID] - gap*1000
		}
		return false
	} else {
		return true
	}
}

func IsWaiting(userID string) int {
	for index, value := range waitList {
		if value == userID {
			return index
		}
	}
	return -1
}
