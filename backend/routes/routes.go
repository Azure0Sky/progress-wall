package routes

import (
	"strings"

	"progress-wall-backend/config"
	"progress-wall-backend/handlers/auth"
	"progress-wall-backend/handlers/board"
	"progress-wall-backend/handlers/column"
	"progress-wall-backend/handlers/task"
	"progress-wall-backend/handlers/user"
	"progress-wall-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 设置路由
func SetupRoutes(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// 根据配置设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 配置CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = strings.Split(cfg.CORS.AllowOrigins, ",")
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// 初始化处理器
	loginHandler := auth.NewLoginHandler(db, cfg)
	registerHandler := auth.NewRegisterHandler(db, cfg)
	profileHandler := user.NewProfileHandler(db)
	boardHandler := board.NewBoardHandler(db)
	columnHandler := column.NewColumnHandler(db)
	taskHandler := task.NewTaskHandler(db)

	// 公开路由（不需要认证）
	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", loginHandler.Login)
			authGroup.POST("/register", registerHandler.Register)
		}
	}

	// 受保护的路由（需要认证）
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// 用户相关
		protected.GET("/user/profile", profileHandler.GetProfile)

		// 看板相关
		protected.GET("/boards", boardHandler.GetBoards)
		protected.POST("/boards", boardHandler.CreateBoard)
		protected.GET("/boards/:boardId", boardHandler.GetBoard)
		protected.PUT("/boards/:boardId", boardHandler.UpdateBoard)
		protected.DELETE("/boards/:boardId", boardHandler.DeleteBoard)

		// 列相关
		protected.GET("/boards/:boardId/columns", columnHandler.GetColumns)
		protected.POST("/boards/:boardId/columns", columnHandler.CreateColumn)
		protected.GET("/columns/:columnId", columnHandler.GetColumn)
		protected.PUT("/columns/:columnId", columnHandler.UpdateColumn)
		protected.DELETE("/columns/:columnId", columnHandler.DeleteColumn)

		// 任务相关
		protected.GET("/columns/:columnId/tasks", taskHandler.GetTasks)
		protected.POST("/columns/:columnId/tasks", taskHandler.CreateTask)
		protected.GET("/tasks/:taskId", taskHandler.GetTask)
		protected.PUT("/tasks/:taskId", taskHandler.UpdateTask)
		protected.DELETE("/tasks/:taskId", taskHandler.DeleteTask)
		protected.PATCH("/tasks/:taskId/move", taskHandler.MoveTask)
	}

	return r
}
