package routes

import (
	"progress-wall-backend/config"
	"progress-wall-backend/handlers/activity"
	"progress-wall-backend/handlers/auth"
	"progress-wall-backend/handlers/user"
	"progress-wall-backend/middleware"
	"strings"

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
	boardActivitiesHandler := activity.NewBoardActivitiesHandler(db)
	taskActivitiesHandler := activity.NewTaskActivitiesHandler(db)

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

		// 看板活动日志
		protected.GET("/boards/:board_id/activities", boardActivitiesHandler.GetBoardActivities)

		// 任务活动日志
		protected.GET("/tasks/:task_id/activities", taskActivitiesHandler.GetTaskActivities)
	}

	return r
}
