package services

import (
	"errors"
	"progress-wall-backend/models"
	"progress-wall-backend/repository"

	"gorm.io/gorm"
)

// ActivityLogService 活动日志服务
type ActivityLogService struct {
	repo repository.ActivityLogRepository
	db   *gorm.DB
}

// NewActivityLogService 创建活动日志服务实例
func NewActivityLogService(
	repo repository.ActivityLogRepository,
	db *gorm.DB,
) *ActivityLogService {
	return &ActivityLogService{
		repo: repo,
		db:   db,
	}
}

// CreateLog 创建活动日志
func (s *ActivityLogService) CreateLog(log *models.ActivityLog) error {
	return s.repo.Create(log)
}

// GetBoardActivities 获取看板活动日志
func (s *ActivityLogService) GetBoardActivities(boardID uint, page, pageSize int) ([]models.ActivityLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // 默认每页20条
	}

	offset := (page - 1) * pageSize
	return s.repo.GetByBoardID(boardID, offset, pageSize)
}

// GetTaskActivities 获取任务活动日志
func (s *ActivityLogService) GetTaskActivities(taskID uint, page, pageSize int) ([]models.ActivityLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // 默认每页20条
	}

	offset := (page - 1) * pageSize
	return s.repo.GetByTaskID(taskID, offset, pageSize)
}

// VerifyBoardAccess 验证用户是否有权访问看板
func (s *ActivityLogService) VerifyBoardAccess(userID, boardID uint) error {
	// 首先检查看板是否存在
	var board models.Board
	if err := s.db.First(&board, boardID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBoardNotFound
		}
		return err
	}

	// 检查用户是否是看板所有者或项目成员
	var count int64
	err := s.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", board.ProjectID, userID).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 && board.OwnerID != userID {
		return ErrAccessDenied
	}

	return nil
}

// VerifyTaskAccess 验证用户是否有权访问任务
func (s *ActivityLogService) VerifyTaskAccess(userID, taskID uint) error {
	// 首先检查任务是否存在
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	// 检查用户是否是项目成员
	var count int64
	err := s.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", task.ProjectID, userID).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 && task.CreatorID != userID {
		return ErrAccessDenied
	}

	return nil
}
