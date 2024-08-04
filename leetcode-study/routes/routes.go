package routes

import (
	"leetcode/controllers"
	"leetcode/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置路由
func SetupRoutes() *gin.Engine {
	dbClient := dao.NewDB()
	// dbSecretClinet := dao.NewSecretDB()
	// 初始化控制器
	userController := controllers.NewUserControllers(dbClient)
	// userSecretController := controllers.NewUserSecretControllers(dbSecretClinet)
	router := gin.Default()

	// Serve static files
	router.Static("/css", "./frontend/css")
	router.Static("/js", "./frontend/js")

	// 路由根路径返回 HTML 页面
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	api := router.Group("/api")
	{
		// 用户相关路由
		api.GET("/users", userController.GetAllUsersHandler)
		api.POST("/users", userController.AddUserHandler)
		api.DELETE("/users/:id", userController.DeleteUserHandler)

		// 管理员相关 API
		api.POST("/reset", userController.ResetHandler)
		api.POST("/attendance", userController.StartAttendanceHandler)
	}

	// 加载 HTML 模板
	router.LoadHTMLGlob("frontend/*.html")

	return router
}
