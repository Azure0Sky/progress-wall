package routes

import (
	"strings"

	"progress-wall-backend/config"
	"progress-wall-backend/handlers/activity"
	"progress-wall-backend/handlers/auth"
	"progress-wall-backend/handlers/board"
	"progress-wall-backend/handlers/column"
	"progress-wall-backend/handlers/notification"
	"progress-wall-backend/handlers/project"
	"progress-wall-backend/handlers/task"
	"progress-wall-backend/handlers/team"
	"progress-wall-backend/handlers/user"
	"progress-wall-backend/middleware"
	"progress-wall-backend/services"

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
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true

	// 在开发环境下允许所有Origin，避免跨域问题
	if cfg.Server.Mode == "release" {
		corsConfig.AllowOrigins = strings.Split(cfg.CORS.AllowOrigins, ",")
	} else {
		corsConfig.AllowOriginFunc = func(origin string) bool {
			return true
		}
	}

	r.Use(cors.New(corsConfig))

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	permService := services.NewPermissionService(db)
	rbac := middleware.NewRBACMiddleware(permService, db)

	// 初始化处理器
	loginHandler := auth.NewLoginHandler(db, cfg)
	registerHandler := auth.NewRegisterHandler(db, cfg)
	profileHandler := user.NewProfileHandler(db)
	projectHandler := project.NewProjectHandler(db)
	boardHandler := board.NewBoardHandler(db)
	columnHandler := column.NewColumnHandler(db)
	taskHandler := task.NewTaskHandler(db)
	teamHandler := team.NewTeamHandler(db)
	boardActivitiesHandler := activity.NewBoardActivitiesHandler(db)
	taskActivitiesHandler := activity.NewTaskActivitiesHandler(db)
	// 添加通知处理器初始化
	notificationHandler := notification.NewNotificationHandler(db)

	// 公开路由（不需要认证）
	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", loginHandler.Login)
			authGroup.POST("/register", registerHandler.Register)
		}

		// 通知相关
		api.POST("/notifications", notificationHandler.ReceiveTaskNotification)
	}

	// 受保护的路由（需要认证）
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// 用户相关
		protected.GET("/user/profile", profileHandler.GetProfile)
		protected.PUT("/user/profile", profileHandler.UpdateProfile)
		protected.POST("/user/avatar", profileHandler.UploadAvatar)

		// Team Routes
		protected.POST("/teams", teamHandler.CreateTeam)
		protected.GET("/teams", teamHandler.GetMyTeams)
		protected.GET("/teams/:teamId",
			rbac.RequireTeamAccess("view", "teamId"),
			teamHandler.GetTeam,
		)
		protected.GET("/teams/:teamId/members",
			rbac.RequireTeamAccess("view", "teamId"),
			teamHandler.GetTeamMembers,
		)
		protected.POST("/teams/:teamId/members",
			rbac.RequireTeamAccess("manage", "teamId"),
			teamHandler.AddMember,
		)

		// Project Routes
		protected.POST("/teams/:teamId/projects",
			rbac.RequireTeamAccess("manage", "teamId"),
			projectHandler.CreateProject,
		)
		protected.GET("/teams/:teamId/projects",
			rbac.RequireTeamAccess("view", "teamId"),
			projectHandler.GetTeamProjects,
		)
		protected.GET("/projects", projectHandler.GetProjects)
		protected.GET("/projects/:projectId",
			rbac.RequireProjectAccess("view", "projectId", "project"),
			projectHandler.GetProject,
		)
		protected.PUT("/projects/:projectId",
			rbac.RequireProjectAccess("manage", "projectId", "project"),
			projectHandler.UpdateProject,
		)
		protected.DELETE("/projects/:projectId",
			rbac.RequireProjectAccess("manage", "projectId", "project"),
			projectHandler.DeleteProject,
		)

		// 看板相关
		protected.GET("/boards", boardHandler.GetBoards)
		protected.GET("/projects/:projectId/boards",
			rbac.RequireProjectAccess("view", "projectId", "project"),
			boardHandler.GetBoardsByProject,
		)
		protected.POST("/projects/:projectId/boards",
			rbac.RequireProjectAccess("manage", "projectId", "project"),
			boardHandler.CreateBoard,
		)
		protected.GET("/boards/:boardId",
			rbac.RequireProjectAccess("view", "boardId", "board"),
			boardHandler.GetBoard,
		)
		protected.PUT("/boards/:boardId",
			rbac.RequireProjectAccess("manage", "boardId", "board"),
			boardHandler.UpdateBoard,
		)
		protected.DELETE("/boards/:boardId",
			rbac.RequireProjectAccess("manage", "boardId", "board"),
			boardHandler.DeleteBoard,
		)

		// 列相关
		protected.GET("/boards/:boardId/columns",
			rbac.RequireProjectAccess("view", "boardId", "board"),
			columnHandler.GetColumns,
		)
		protected.POST("/boards/:boardId/columns",
			// Only admins can create columns
			rbac.RequireProjectAccess("manage", "boardId", "board"),
			columnHandler.CreateColumn,
		)
		protected.GET("/columns/:columnId",
			rbac.RequireProjectAccess("view", "columnId", "column"),
			columnHandler.GetColumn,
		)
		protected.PUT("/columns/:columnId",
			rbac.RequireProjectAccess("manage", "columnId", "column"),
			columnHandler.UpdateColumn,
		)
		protected.DELETE("/columns/:columnId",
			rbac.RequireProjectAccess("manage", "columnId", "column"),
			columnHandler.DeleteColumn,
		)

		// 任务相关
		protected.GET("/columns/:columnId/tasks",
			rbac.RequireProjectAccess("view", "columnId", "column"),
			taskHandler.GetTasks,
		)
		protected.POST("/columns/:columnId/tasks",
			rbac.RequireProjectAccess("view", "columnId", "column"),
			taskHandler.CreateTask,
		)
		protected.GET("/tasks/:taskId",
			rbac.RequireProjectAccess("view", "taskId", "task"),
			taskHandler.GetTask,
		)
		protected.PUT("/tasks/:taskId",
			rbac.RequireProjectAccess("view", "taskId", "task"),
			taskHandler.UpdateTask,
		)
		protected.DELETE("/tasks/:taskId",
			rbac.RequireProjectAccess("view", "taskId", "task"),
			taskHandler.DeleteTask,
		)
		protected.PATCH("/tasks/:taskId/move",
			rbac.RequireProjectAccess("view", "taskId", "task"),
			taskHandler.MoveTask,
		)

		// 看板活动日志
		protected.GET("/boards/:boardId/activities",
			rbac.RequireProjectAccess("view", "boardId", "board"),
			boardActivitiesHandler.GetBoardActivities,
		)

		// 任务活动日志
		protected.GET("/tasks/:taskId/activities",
			rbac.RequireProjectAccess("view", "taskId", "task"),
			taskActivitiesHandler.GetTaskActivities,
		)
	}

	return r
}
