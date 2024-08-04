package models

import (
	"errors"
	"leetcode/config"
	"leetcode/dao"
	"leetcode/entities"
	"leetcode/utils"
	"sort"
)

type UserSecretService struct {
	dbClient dao.DBClientInterface
}

// NewUserService 创建新的 UserService 实例
func NewUserSecretService(dbClient dao.DBClientInterface) *UserService {
	return &UserService{dbClient: dbClient}
}

// CreateUser 创建新用户
func (s *UserSecretService) CreateUser(user *entities.User) error {
	if err := s.dbClient.AddUser(user); err != nil {
		return err
	}
	return nil
}

// GetAllUsers 获取所有用户
func (s *UserSecretService) GetAllUsers() ([]*entities.User, error) {
	keys := s.dbClient.GetAllKeys()
	var users []*entities.User

	for _, ID := range keys {
		user, err := s.dbClient.GetUserByID(string(ID))
		if err != nil {
			return nil, err
		}
		temp, err := utils.ConvrtUserFormatByteToNormal(user)
		if err != nil {
			return nil, err
		}
		users = append(users, temp)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Level < users[j].Level
	})

	for _, user := range users {
		user.NickName = config.Lazy[user.Level]
	}
	return users, nil
}

// DeleteUser 删除用户
func (s *UserSecretService) DeleteUser(username string) error {
	if err := s.dbClient.DeleteUser(username); err != nil {
		return err
	}
	return nil
}

// Authenticate 用户认证
func (s *UserSecretService) Authenticate(username, password string) (*entities.UserSecret, error) {
	user, err := s.dbClient.GetUserByID(username)
	if err != nil {
		return nil, err
	}
	temp, err := utils.ConvrtUserFormatByteToSecret(user)
	if err != nil {
		return nil, err
	}

	if temp.Password != password {
		return nil, errors.New("invalid credentials")
	}

	return temp, err
}
