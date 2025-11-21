package services

import (
	"errors"
	"progress-wall-backend/models"

	"gorm.io/gorm"
)

// ColumnService 列服务
type ColumnService struct {
	db *gorm.DB
}

// NewColumnService 创建列服务
func NewColumnService(db *gorm.DB) *ColumnService {
	return &ColumnService{
		db: db,
	}
}

// GetColumnByID 根据ID获取列
func (s *ColumnService) GetColumnByID(columnID uint) (*models.Column, error) {
	var column models.Column
	result := s.db.
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		First(&column, columnID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrColumnNotFound
		}
		return nil, errors.New("查询列失败")
	}

	return &column, nil
}

// GetColumnsByBoardID 获取看板的所有列
func (s *ColumnService) GetColumnsByBoardID(boardID uint) ([]models.Column, error) {
	var columns []models.Column
	result := s.db.
		Where("board_id = ?", boardID).
		Order("position ASC").
		Find(&columns)

	if result.Error != nil {
		return nil, errors.New("查询列列表失败")
	}

	return columns, nil
}

// CreateColumn 创建列
func (s *ColumnService) CreateColumn(column *models.Column) error {
	// 获取当前看板的最大position
	var maxPosition int
	s.db.Model(&models.Column{}).
		Where("board_id = ?", column.BoardID).
		Select("COALESCE(MAX(position), -1)").
		Scan(&maxPosition)
	column.Position = maxPosition + 1

	if err := s.db.Create(column).Error; err != nil {
		return errors.New("创建列失败")
	}
	return nil
}

// UpdateColumn 更新列
func (s *ColumnService) UpdateColumn(columnID uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.Column{}).Where("id = ?", columnID).Updates(updates)
	if result.Error != nil {
		return errors.New("更新列失败")
	}
	if result.RowsAffected == 0 {
		return ErrColumnNotFound
	}
	return nil
}

// DeleteColumn 删除列（软删除）
func (s *ColumnService) DeleteColumn(columnID uint) error {
	result := s.db.Delete(&models.Column{}, columnID)
	if result.Error != nil {
		return errors.New("删除列失败")
	}
	if result.RowsAffected == 0 {
		return ErrColumnNotFound
	}
	return nil
}

// ReorderColumns 重新排序列
func (s *ColumnService) ReorderColumns(boardID uint, columnIDs []uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, columnID := range columnIDs {
		if err := tx.Model(&models.Column{}).
			Where("id = ? AND board_id = ?", columnID, boardID).
			Update("position", i).Error; err != nil {
			tx.Rollback()
			return errors.New("重新排序列失败")
		}
	}

	return tx.Commit().Error
}
