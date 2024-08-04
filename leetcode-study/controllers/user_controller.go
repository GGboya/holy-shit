package controllers

import (
	"log"
	"net/http"

	"leetcode/dao"
	"leetcode/entities"
	"leetcode/models"

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

	err := ctl.UserService.CreateUser(&newUser)
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
func (ctl *UserController) ResetHandler(c *gin.Context) {
	err := ctl.UserService.Reset()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Users reset successfully"})
}

func (ctl *UserController) StartAttendanceHandler(c *gin.Context) {
	err := ctl.UserService.StartAttendance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start attendance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Attendance started successfully"})
}
