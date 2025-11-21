package services

import (
	"errors"
	"progress-wall-backend/models"

	"gorm.io/gorm"
)

// BoardService 看板服务
type BoardService struct {
	db *gorm.DB
}

// NewBoardService 创建看板服务
func NewBoardService(db *gorm.DB) *BoardService {
	return &BoardService{
		db: db,
	}
}

// GetBoardByID 根据ID获取看板（包含嵌套的列和任务）
func (s *BoardService) GetBoardByID(boardID uint) (*models.Board, error) {
	var board models.Board
	result := s.db.
		Preload("Columns", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Columns.Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Columns.Tasks.Assignee").
		Preload("Columns.Tasks.Creator").
		Preload("Owner").
		First(&board, boardID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrBoardNotFound
		}
		return nil, errors.New("查询看板失败")
	}

	return &board, nil
}

// GetBoardsByUserID 获取用户的所有看板
func (s *BoardService) GetBoardsByUserID(userID uint) ([]models.Board, error) {
	var boards []models.Board
	result := s.db.
		Where("owner_id = ?", userID).
		Order("position ASC").
		Find(&boards)

	if result.Error != nil {
		return nil, errors.New("查询看板列表失败")
	}

	return boards, nil
}

// CreateBoard 创建看板
func (s *BoardService) CreateBoard(board *models.Board) error {
	if err := s.db.Create(board).Error; err != nil {
		return errors.New("创建看板失败")
	}
	return nil
}

// UpdateBoard 更新看板
func (s *BoardService) UpdateBoard(boardID uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.Board{}).Where("id = ?", boardID).Updates(updates)
	if result.Error != nil {
		return errors.New("更新看板失败")
	}
	if result.RowsAffected == 0 {
		return ErrBoardNotFound
	}
	return nil
}

// DeleteBoard 删除看板（软删除）
func (s *BoardService) DeleteBoard(boardID uint) error {
	result := s.db.Delete(&models.Board{}, boardID)
	if result.Error != nil {
		return errors.New("删除看板失败")
	}
	if result.RowsAffected == 0 {
		return ErrBoardNotFound
	}
	return nil
}
