package database

import (
	"log"
	"progress-wall-backend/models"

	"gorm.io/gorm"
)

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) error {
	log.Println("开始执行数据库迁移...")

	// 迁移所有模型
	err := db.AutoMigrate(
		// 用户和权限相关
		&models.User{},
		&models.UserPermission{},
		&models.UserPermissionGroup{},
		&models.UserPermissionGroupUser{},
		&models.UserPermissionGroup{},

		// 项目和看板相关
		&models.Project{},
		&models.ProjectMember{},
		&models.Board{},

		// 任务相关
		&models.Column{},
		&models.Task{},
		&models.Comment{},
		&models.Attachment{},
		&models.Label{},
		&models.TaskLabel{},

		// 活动日志
		&models.ActivityLog{},
	)

	if err != nil {
		log.Printf("数据库迁移失败: %v", err)
		return err
	}

	log.Println("数据库迁移完成")
	return nil
}
