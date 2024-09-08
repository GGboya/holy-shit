// dao/db.go

package dao

import (
	"encoding/json"
	"fmt"
	"os"

	"leetcode/ggdb"
	"leetcode/model"
	// Import other necessary packages
)

var db *ggdb.DB

// Init initializes the database
func Init() error {
	op := ggdb.Options{
		DataFileSize: 1 << 20,
		DirPath:      os.Getenv("DIRPATH"),
		IndexType:    ggdb.BTree,
		SyncWrites:   true,
	}
	var err error
	fmt.Println(os.Getenv("DIRPATH"))
	db, err = ggdb.Open(op)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	return nil
}

func GetDB() *ggdb.DB {
	return db
}

// GetUserByID retrieves a user from the database by ID
func GetUserByID(userID string) (*model.User, error) {
	value, err := db.Get([]byte(userID))
	if err != nil {
		return nil, fmt.Errorf("error getting user from database: %v", err)
	}

	var user model.User
	if err := json.Unmarshal(value, &user); err != nil {
		return nil, fmt.Errorf("error unmarshaling user data: %v", err)
	}
	return &user, nil
}

// AddUser adds a new user to the database
func AddUser(user *model.User) error {
	value, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshaling user data: %v", err)
	}
	if err := db.Put([]byte(user.ID), value); err != nil {
		return fmt.Errorf("error putting user data into database: %v", err)
	}
	return nil
}

func DeleteUser(userID string) error {
	if err := db.Delete([]byte(userID)); err != nil {
		return fmt.Errorf("error deleting user from database: %v", err)
	}
	return nil
}
