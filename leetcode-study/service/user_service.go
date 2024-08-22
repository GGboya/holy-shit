package service

import (
	"fmt"
	"leetcode/config"
	"leetcode/dao"
	"leetcode/model"
	"leetcode/utils"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"gopkg.in/gomail.v2"
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
	qq := strings.TrimSpace(QQ)
	user := &model.User{
		ID:     ID,
		QQ:     qq,
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
	fmt.Printf("ID: %s, QQ: %s, QQName: %s, Level: %s\n", user.ID, user.QQ, user.QQName, config.Lazy[user.Level])
	return err
}

func GetAllUser() (err error) {
	db := dao.GetDB()
	keys := db.Iterator()
	fmt.Printf("考勤总人数: %d\n", len(keys))

	userinfos := []*model.User{}

	for _, ID := range keys {
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

func Reset() (err error) {
	db := dao.GetDB()
	keys := db.Iterator()
	for _, ID := range keys {
		var user *model.User
		user, err = dao.GetUserByID(string(ID))
		if err != nil {
			return err
		}
		user.Level = 0
		dao.AddUser(user)
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
	god := []string{}
	type msg struct {
		qq, title, content string
	}
	msgs := []msg{}

	fmt.Println("今日勤奋的同学是")
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

	fmt.Println("今日懒惰的同学是")
	for i := range lazy {
		level := lazy[i].Level
		originName := config.Lazy[level]
		nowName := config.Lazy[min(9, level+1)]
		fmt.Printf("%d --- ID: %s, QQ: %s, QQName: %s, Level: %s --> %s\n", i+1, lazy[i].ID, lazy[i].QQ, lazy[i].QQName, originName, nowName)
		lazy[i].Level = min(9, level+1)
		dao.AddUser(&lazy[i])
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
			SendEmail(msgs[i].qq, msgs[i].title, msgs[i].content)
		}(i)
	}

	wg.Wait()

	if len(god) != 0 {
		content := strings.Join(god, "\n")
		title := "有同学要飞升了，请手动处理邮件内容中的同学"
		SendEmail(os.Getenv("QQ"), title, content)
	}
	return
}

func SendEmail(qq, title, content string) error {
	// 从环境变量中获取邮箱信息
	qqEmail := os.Getenv("QQ_EMAIL")
	qqAuthCode := os.Getenv("QQ_AUTH_CODE")

	// 验证环境变量是否加载成功
	if qqEmail == "" || qqAuthCode == "" {
		log.Fatal("QQ_EMAIL or QQ_AUTH_CODE is not set in .env file")
	}

	m := gomail.NewMessage()

	// 发件人
	m.SetHeader("From", qqEmail)

	// 收件人
	m.SetHeader("To", qq+"@qq.com")

	// 邮件标题
	m.SetHeader("Subject", title)

	// 邮件内容
	m.SetBody("text/plain", content)

	// QQ邮箱SMTP服务器信息
	d := gomail.NewDialer("smtp.qq.com", 587, qqEmail, qqAuthCode)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Email sent successfully!")
	return nil
}
