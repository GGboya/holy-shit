package main

import (
	"fmt"
	"leetcode/dao"
	"leetcode/service"
	"log"
)

func main() {

	// Initialize the database
	if err := dao.Init(); err != nil {
		log.Fatalf("error initializing database: %v", err)
	}

	// 看用户想执行什么操作？
	for {
		var action int
		fmt.Println("============================欢迎使用考勤系统============================")
		fmt.Println("请输入你想要执行的操作：")
		fmt.Println("1. 添加用户")
		fmt.Println("2. 获取用户")
		fmt.Println("3. 删除用户")
		fmt.Println("4. 获取所有用户")
		fmt.Println("5. 开始考勤")
		fmt.Println("6. 重置")
		fmt.Println("7. 退出")

		fmt.Scan(&action)
		switch action {
		case 1:
			err := service.CreateUser()
			if err != nil {
				log.Fatalf("error creating user: %v", err)
			}
		case 2:
			err := service.GetUser()
			if err != nil {
				log.Fatalf("error getting user: %v", err)
			}
		case 3:
			err := service.DeleteUser()
			if err != nil {
				log.Fatalf("error deleting user: %v", err)
			}
		case 4:
			err := service.GetAllUser()
			if err != nil {
				log.Fatalf("error deleting user: %v", err)
			}
		case 5:
			err := service.StartAttendance()
			if err != nil {
				log.Fatalf("error starting attendance: %v", err)
			}
		case 6:
			// 重置
			err := service.Reset()
			if err != nil {
				log.Fatal("error resetting database: %v", err)
			}
		case 7:
			return
		}
	}

}
