package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"leetcode/dao"
	"leetcode/entities"
	"leetcode/models"
	"leetcode/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *models.UserService
}

func NewUserControllers(dbClient dao.DBClientInterface) *UserController {
	return &UserController{UserService: models.NewUserService(dbClient)}
}
func (ctl *UserController) GetAllUsersHandler(c *gin.Context) {
	users, err := ctl.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// AddUserHandler 添加用户的处理程序
func (ctl *UserController) AddUserHandler(c *gin.Context) {
	var newUser entities.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// 首先需要鉴定一下 id 是否存在

	_, err := utils.FetchLastSubmitTime(newUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid leetcode id"})
		return
	}

	// 还得确认下 QQ 是否存在
	err = utils.SendEmail(newUser.QQ, "welcome to ggboy coding", "this email show that you join the leetcode study")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid QQ number"})
		return
	}

	err = ctl.UserService.CreateUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}

// DeleteUserHandler 删除用户的处理程序
func (ctl *UserController) DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	log.Println("id:", id)
	err := ctl.UserService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (ctl *UserController) StartDailyJob() {
	for {
		// 计算下一个 24:00 的时间
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		duration := time.Until(next)

		fmt.Printf("Next job will run in %v\n", duration)

		// 等待直到下一个 24:00
		time.Sleep(duration)

		// 获取当前日期
		currentTime := time.Now()

		// 判断当前日期是否是月末
		if isLastDayOfMonth(currentTime) {
			// 如果是月末，执行 ResetHandler
			log.Println("It's the last day of the month, resetting users...")
			err := ctl.UserService.Reset() // 调用重置功能
			if err != nil {
				fmt.Println("Failed to reset users:", err)
			}
		} else {
			// 否则执行考勤任务
			log.Println("Starting daily attendance...")
			err := ctl.UserService.StartAttendance() // 调用考勤功能
			if err != nil {
				fmt.Println("Failed to start attendance:", err)
			}
		}
	}
}

// 判断是否是月末
func isLastDayOfMonth(t time.Time) bool {
	nextDay := t.AddDate(0, 0, 1) // 获取下一天
	return nextDay.Day() == 1     // 如果下一天是1号，则今天是月末
}
