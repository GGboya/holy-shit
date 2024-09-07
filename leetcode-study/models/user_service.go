package models

import (
	"fmt"
	"leetcode/config"
	"leetcode/dao"
	"leetcode/entities"
	"leetcode/utils"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type UserService struct {
	dbClient dao.DBClientInterface
}

// NewUserService 创建新的 UserService 实例
func NewUserService(dbClient dao.DBClientInterface) *UserService {
	return &UserService{dbClient: dbClient}
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(user *entities.User) error {
	if err := s.dbClient.AddUser(user); err != nil {
		return err
	}
	return nil
}

// GetAllUsers 获取所有用户
func (s *UserService) GetAllUsers() ([]*entities.User, error) {
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

// Reset 重置用户等级
func (s *UserService) Reset() error {
	keys := s.dbClient.GetAllKeys()
	for _, ID := range keys {
		user, err := s.dbClient.GetUserByID(string(ID))
		if err != nil {
			return err
		}
		temp, err := utils.ConvrtUserFormatByteToNormal(user)
		if err != nil {
			return err
		}
		temp.Level = 0
		if err := s.dbClient.AddUser(temp); err != nil {
			return err
		}
	}
	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ID string) error {
	if err := s.dbClient.DeleteUser(ID); err != nil {
		return err
	}
	return nil
}
func (s *UserService) StartAttendance() (err error) {
	fmt.Println("开始考勤")
	// 遍历 db 索引中的所有 key
	hard := []entities.User{}
	lazy := []entities.User{}

	keys := s.dbClient.GetAllKeys()
	mu := utils.NewSpinLock()
	wg := sync.WaitGroup{}
	for {
		var temerr error
		for i := range keys {
			// 获取用户的 ID， 根据这个 ID 访问 leetcode
			wg.Add(1)
			go func(i int) {
				if err := func() error {
					defer wg.Done()
					ID := keys[i]
					var lastSubmitTime *time.Time
					var userinfo *entities.User
					temp, err := s.dbClient.GetUserByID(string(ID))
					if err != nil {
						log.Println("Error getting user by ID:", err)
						return err
					}
					userinfo, err = utils.ConvrtUserFormatByteToNormal(temp)
					if err != nil {
						log.Println("Error converting user format byte to normal:", err)
						return err
					}
					lastSubmitTime, err = utils.FetchLastSubmitTime(string(ID))
					if err != nil {
						log.Println("Error fetching last submit time:", err)
						return err
					}
					currentTime := time.Now()
					duration := currentTime.Sub(*lastSubmitTime)
					mu.Lock()
					defer mu.Unlock()
					if duration < 24*time.Hour {
						hard = append(hard, *userinfo)
					} else {
						lazy = append(lazy, *userinfo)
					}
					return nil
				}(); err != nil {
					temerr = err
				}

			}(i)
		}
		if temerr == nil {
			break
		}
	}

	wg.Wait()

	sort.Slice(hard, func(i, j int) bool {
		if hard[i].Level == hard[j].Level {
			return i < j
		}
		return hard[i].Level < hard[j].Level
	})

	sort.Slice(lazy, func(i, j int) bool {
		if lazy[i].Level == lazy[j].Level {
			return i < j
		}
		return lazy[i].Level < lazy[j].Level
	})

	type msg struct {
		qq, title, content string
	}
	msgs := []msg{}
	god := []string{}
	log.Println("今日勤奋的同学是")
	for i := range hard {
		nickName := config.Lazy[hard[i].Level]
		fmt.Printf("%d --- ID: %s, QQ: %s, QQName: %s, Level: %s\n", i+1, hard[i].ID, hard[i].QQ, hard[i].QQName, nickName)
		content := fmt.Sprintf("你的称号没有变化， 仍旧是 %s, 请继续努力", nickName)
		title := "恭喜你，完成考勤！！！"
		msgs = append(msgs, msg{
			qq:      hard[i].QQ,
			title:   title,
			content: content,
		})
	}
	log.Println("今日懒惰的同学是")
	for i := range lazy {
		level := lazy[i].Level
		originName := config.Lazy[level]
		nowName := config.Lazy[min(9, level+1)]
		log.Printf("%d --- ID: %s, QQ: %s, QQName: %s, Level: %s --> %s\n", i+1, lazy[i].ID, lazy[i].QQ, lazy[i].QQName, originName, nowName)
		lazy[i].Level = min(9, level+1)
		s.dbClient.AddUser(&lazy[i])

		content := fmt.Sprintf("你的称号由 %s 变为 %s, 请继续努力", originName, nowName)
		title := "你今天没有刷题哦， 请再接再厉"
		msgs = append(msgs, msg{
			qq:      lazy[i].QQ,
			title:   title,
			content: content,
		})
		if level == 9 {
			god = append(god, lazy[i].QQ)
		}
	}

	wg = sync.WaitGroup{}
	for i := range msgs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			utils.SendEmail(msgs[i].qq, msgs[i].title, msgs[i].content)
		}(i)
	}

	wg.Wait()

	if len(god) != 0 {
		content := strings.Join(god, "\n")
		title := "有同学要飞升了，请手动处理邮件内容中的同学"
		utils.SendEmail(os.Getenv("QQ"), title, content)
	}
	return
}
