package service

import (
	"fmt"
	"leetcode/config"
	"leetcode/dao"
	"leetcode/model"
	"leetcode/utils"
	"sort"
	"sync"
	"time"
)

func CreateUser() (err error) {
	fmt.Println("请根据提示输入你要增加的用户信息")
	var ID, QQ, QQName string
	fmt.Print("ID: ")
	fmt.Scan(&ID)
	fmt.Print("QQ: ")
	fmt.Scan(&QQ)
	fmt.Print("QQName: ")
	fmt.Scan(&QQName)
	user := &model.User{
		ID:     ID,
		QQ:     QQ,
		QQName: QQName,
	}
	if err = dao.AddUser(user); err != nil {
		return err
	}
	return
}

func GetUser() (err error) {
	fmt.Println("请输入你要查询的用户ID")
	var ID string
	fmt.Scan(&ID)
	user, err := dao.GetUserByID(ID)
	fmt.Printf("ID: %s, QQ: %s, QQName: %s, Level: %d\n", user.ID, user.QQ, user.QQName, user.Level)
	return err
}

func GetAllUser() (err error) {
	db := dao.GetDB()
	keys := db.Iterator()
	fmt.Printf("考勤总人数: %d\n", len(keys))

	userinfos := []*model.User{}

	for _, ID := range db.Iterator() {
		var user *model.User
		user, err = dao.GetUserByID(string(ID))
		userinfos = append(userinfos, user)
	}

	sort.Slice(userinfos, func(i, j int) bool {
		return userinfos[i].Level < userinfos[j].Level
	})

	for i, user := range userinfos {
		nickName := config.Lazy[user.Level]
		fmt.Printf("%d --- ID: %s, QQ: %s, QQName: %s, Level: %s\n", i+1, user.ID, user.QQ, user.QQName, nickName)
	}
	return
}

func DeleteUser() (err error) {
	fmt.Println("请输入你要删除的用户ID")
	var ID string
	fmt.Scan(&ID)
	if err = dao.DeleteUser(ID); err != nil {
		return err
	}
	return
}

func StartAttendance() (err error) {
	fmt.Println("开始考勤")
	// 遍历 db 索引中的所有 key
	db := dao.GetDB()

	hard := []model.User{}
	lazy := []model.User{}

	keys := db.Iterator()
	mu := utils.NewSpinLock()
	wg := sync.WaitGroup{}
	for {
		for i := range keys {
			// 获取用户的 ID， 根据这个 ID 访问 leetcode
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				ID := keys[i]
				var lastSubmitTime *time.Time
				var userinfo *model.User
				userinfo, err = dao.GetUserByID(string(ID))
				lastSubmitTime, err = utils.FetchLastSubmitTime(string(ID))
				if err != nil {
					fmt.Println("Error fetching last submit time:", err)
					return
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
			}(i)
		}
		if err == nil {
			break
		}
	}

	wg.Wait()

	sort.Slice(hard, func(i, j int) bool {
		return hard[i].Level < hard[j].Level
	})

	sort.Slice(lazy, func(i, j int) bool {
		return lazy[i].Level < lazy[j].Level
	})

	fmt.Println("今日勤奋的同学是")
	for i := range hard {
		nickName := config.Lazy[hard[i].Level]
		fmt.Printf("%d --- ID: %s, QQ: %s, QQName: %s, Level: %s\n", i+1, hard[i].ID, hard[i].QQ, hard[i].QQName, nickName)
	}
	fmt.Println("今日懒惰的同学是")
	for i := range lazy {
		level := lazy[i].Level
		originName := config.Lazy[level]
		nowName := config.Lazy[min(9, level+1)]
		fmt.Printf("%d --- ID: %s, QQ: %s, QQName: %s, Level: %s --> %s\n", i+1, lazy[i].ID, lazy[i].QQ, lazy[i].QQName, originName, nowName)
		lazy[i].Level = min(9, level+1)
		dao.AddUser(&lazy[i])
	}
	return
}
