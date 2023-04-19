package db

import (
	"github.com/boltdb/bolt"
	"log"
	"time"
)

/* ---- BOLT DB FUNCTIONS ---- */
var db *bolt.DB

func GetDB() *bolt.DB {
	var err error
	db, err = bolt.Open("UNKNOWN", 0600, nil)
	if err != nil {
		log.Panic("Cannot create DB")
	}
	return db
}

func PutDB(db *bolt.DB, bucketName string, key string, value []byte) {
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

func ViewDB(bucketName string, key string) []byte {
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
	return value
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
    log.Println(tempUserMap)
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
