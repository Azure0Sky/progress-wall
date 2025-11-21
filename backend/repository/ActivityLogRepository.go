package repository

import (
	"progress-wall-backend/models"

	"gorm.io/gorm"
)

// ActivityLogRepository 活动日志仓储接口
type ActivityLogRepository interface {
	// Create 创建活动日志
	Create(log *models.ActivityLog) error

	// GetByBoardID 根据看板ID获取活动日志（支持分页）
	GetByBoardID(boardID uint, offset, limit int) ([]models.ActivityLog, int64, error)

	// GetByTaskID 根据任务ID获取活动日志（支持分页）
	GetByTaskID(taskID uint, offset, limit int) ([]models.ActivityLog, int64, error)

	// GetByProjectID 根据项目ID获取活动日志（支持分页）
	GetByProjectID(projectID uint, offset, limit int) ([]models.ActivityLog, int64, error)
}

type activityLogRepository struct {
	db *gorm.DB
}

// NewActivityLogRepository 创建活动日志仓储实例
func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

// Create 创建活动日志
func (r *activityLogRepository) Create(log *models.ActivityLog) error {
	return r.db.Create(log).Error
}

// GetByBoardID 根据看板ID获取活动日志（支持分页）
func (r *activityLogRepository) GetByBoardID(boardID uint, offset, limit int) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var total int64

	// 查询总数
	if err := r.db.Model(&models.ActivityLog{}).
		Where("board_id = ?", boardID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询分页数据，预加载用户信息
	err := r.db.Where("board_id = ?", boardID).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, total, err
}

// GetByTaskID 根据任务ID获取活动日志（支持分页）
func (r *activityLogRepository) GetByTaskID(taskID uint, offset, limit int) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var total int64

	// 查询总数
	if err := r.db.Model(&models.ActivityLog{}).
		Where("task_id = ?", taskID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询分页数据，预加载用户信息
	err := r.db.Where("task_id = ?", taskID).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, total, err
}

// GetByProjectID 根据项目ID获取活动日志（支持分页）
func (r *activityLogRepository) GetByProjectID(projectID uint, offset, limit int) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var total int64

	// 查询总数
	if err := r.db.Model(&models.ActivityLog{}).
		Where("project_id = ?", projectID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询分页数据，预加载用户信息
	err := r.db.Where("project_id = ?", projectID).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, total, err
}
