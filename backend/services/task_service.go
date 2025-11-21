package services

import (
	"errors"
	"progress-wall-backend/models"

	"gorm.io/gorm"
)

// TaskService 任务服务
type TaskService struct {
	db *gorm.DB
}

// NewTaskService 创建任务服务
func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		db: db,
	}
}

// GetTaskByID 根据ID获取任务
func (s *TaskService) GetTaskByID(taskID uint) (*models.Task, error) {
	var task models.Task
	result := s.db.
		Preload("Assignee").
		Preload("Creator").
		Preload("Column").
		Preload("Labels").
		First(&task, taskID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, errors.New("查询任务失败")
	}

	return &task, nil
}

// GetTasksByColumnID 获取列的所有任务
func (s *TaskService) GetTasksByColumnID(columnID uint) ([]models.Task, error) {
	var tasks []models.Task
	result := s.db.
		Where("column_id = ?", columnID).
		Order("position ASC").
		Find(&tasks)

	if result.Error != nil {
		return nil, errors.New("查询任务列表失败")
	}

	return tasks, nil
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(task *models.Task) error {
	// 获取当前列的最大position
	var maxPosition int
	s.db.Model(&models.Task{}).
		Where("column_id = ?", task.ColumnID).
		Select("COALESCE(MAX(position), -1)").
		Scan(&maxPosition)
	task.Position = maxPosition + 1

	if err := s.db.Create(task).Error; err != nil {
		return errors.New("创建任务失败")
	}
	return nil
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(taskID uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.Task{}).Where("id = ?", taskID).Updates(updates)
	if result.Error != nil {
		return errors.New("更新任务失败")
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// DeleteTask 删除任务（软删除）
func (s *TaskService) DeleteTask(taskID uint) error {
	result := s.db.Delete(&models.Task{}, taskID)
	if result.Error != nil {
		return errors.New("删除任务失败")
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// MoveTask 移动任务到新列和新位置
func (s *TaskService) MoveTask(taskID uint, newColumnID uint, newOrder int) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取任务
	var task models.Task
	if err := tx.First(&task, taskID).Error; err != nil {
		tx.Rollback()
		return ErrTaskNotFound
	}

	oldColumnID := task.ColumnID

	// 如果移动到不同列，需要更新两个列中的任务位置
	if oldColumnID != newColumnID {
		// 从旧列中移除：将旧列中位置大于当前任务位置的所有任务位置减1
		if err := tx.Model(&models.Task{}).
			Where("column_id = ? AND position > ?", oldColumnID, task.Position).
			Update("position", gorm.Expr("position - 1")).Error; err != nil {
			tx.Rollback()
			return errors.New("更新旧列任务位置失败")
		}

		// 在新列中插入：将新列中位置大于等于newOrder的所有任务位置加1
		if err := tx.Model(&models.Task{}).
			Where("column_id = ? AND position >= ?", newColumnID, newOrder).
			Update("position", gorm.Expr("position + 1")).Error; err != nil {
			tx.Rollback()
			return errors.New("更新新列任务位置失败")
		}

		// 更新任务的列ID和位置
		if err := tx.Model(&task).
			Updates(map[string]interface{}{
				"column_id": newColumnID,
				"position":  newOrder,
			}).Error; err != nil {
			tx.Rollback()
			return errors.New("更新任务位置失败")
		}
	} else {
		// 同一列内移动
		oldPosition := task.Position
		if oldPosition < newOrder {
			// 向后移动：将位置在 (oldPosition, newOrder] 之间的任务位置减1
			if err := tx.Model(&models.Task{}).
				Where("column_id = ? AND position > ? AND position <= ?", newColumnID, oldPosition, newOrder).
				Update("position", gorm.Expr("position - 1")).Error; err != nil {
				tx.Rollback()
				return errors.New("更新任务位置失败")
			}
		} else if oldPosition > newOrder {
			// 向前移动：将位置在 [newOrder, oldPosition) 之间的任务位置加1
			if err := tx.Model(&models.Task{}).
				Where("column_id = ? AND position >= ? AND position < ?", newColumnID, newOrder, oldPosition).
				Update("position", gorm.Expr("position + 1")).Error; err != nil {
				tx.Rollback()
				return errors.New("更新任务位置失败")
			}
		}

		// 更新任务位置
		if err := tx.Model(&task).Update("position", newOrder).Error; err != nil {
			tx.Rollback()
			return errors.New("更新任务位置失败")
		}
	}

	return tx.Commit().Error
}
