package dao

import "leetcode/entities"

// DBClientInterface 定义数据库客户端接口
type DBClientInterface interface {
	GetUserByID(userID string) ([]byte, error)
	AddUser(user *entities.User) error
	DeleteUser(userID string) error
	GetAllKeys() [][]byte
}
