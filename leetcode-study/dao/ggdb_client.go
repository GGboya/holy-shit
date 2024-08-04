// dao/ggdb_client.go

package dao

import (
	"encoding/json"
	"fmt"
	"leetcode/entities"
	"leetcode/ggdb"
)

// DBClient 结构体用于普通数据库
type DBClient struct {
	db *ggdb.DB
}

// NewDBClient 创建新的数据库客户端实例
func NewDBClient(db *ggdb.DB) *DBClient {
	return &DBClient{db: db}
}

// GetUserByID 从数据库中获取用户
func (c *DBClient) GetUserByID(userID string) ([]byte, error) {
	value, err := c.db.Get([]byte(userID))
	if err != nil {
		return nil, fmt.Errorf("error getting user from database: %v", err)
	}
	return value, nil
}

// AddUser 将新用户添加到数据库中
func (c *DBClient) AddUser(user *entities.User) error {
	value, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshaling user data: %v", err)
	}
	if err := c.db.Put([]byte(user.ID), value); err != nil {
		return fmt.Errorf("error putting user data into database: %v", err)
	}
	return nil
}

// DeleteUser 从数据库中删除用户
func (c *DBClient) DeleteUser(userID string) error {
	if err := c.db.Delete([]byte(userID)); err != nil {
		return fmt.Errorf("error deleting user from database: %v", err)
	}
	return nil
}

func (c *DBClient) GetAllKeys() [][]byte {
	keys := c.db.Iterator()
	return keys
}
