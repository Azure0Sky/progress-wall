package column

import (
	"net/http"
	"strconv"

	"progress-wall-backend/models"
	"progress-wall-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ColumnHandler 列处理器
type ColumnHandler struct {
	columnService *services.ColumnService
}

// NewColumnHandler 创建列处理器
func NewColumnHandler(db *gorm.DB) *ColumnHandler {
	return &ColumnHandler{
		columnService: services.NewColumnService(db),
	}
}

// GetColumn 获取单个列
func (h *ColumnHandler) GetColumn(c *gin.Context) {
	columnID, err := strconv.ParseUint(c.Param("columnId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的列ID"})
		return
	}

	column, err := h.columnService.GetColumnByID(uint(columnID))
	if err != nil {
		if err == services.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, column)
}

// GetColumns 获取看板的所有列
func (h *ColumnHandler) GetColumns(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("boardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的看板ID"})
		return
	}

	columns, err := h.columnService.GetColumnsByBoardID(uint(boardID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"columns": columns})
}

// CreateColumn 创建列
func (h *ColumnHandler) CreateColumn(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("boardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的看板ID"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	column := &models.Column{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		BoardID:     uint(boardID),
		Status:      models.ColumnStatusActive,
	}

	if err := h.columnService.CreateColumn(column); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, column)
}

// UpdateColumn 更新列
func (h *ColumnHandler) UpdateColumn(c *gin.Context) {
	columnID, err := strconv.ParseUint(c.Param("columnId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的列ID"})
		return
	}

	var req struct {
		Name        *string              `json:"name"`
		Description *string              `json:"description"`
		Color       *string              `json:"color"`
		Status      *models.ColumnStatus `json:"status"`
		Position    *int                 `json:"position"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Position != nil {
		updates["position"] = *req.Position
	}

	if err := h.columnService.UpdateColumn(uint(columnID), updates); err != nil {
		if err == services.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteColumn 删除列
func (h *ColumnHandler) DeleteColumn(c *gin.Context) {
	columnID, err := strconv.ParseUint(c.Param("columnId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的列ID"})
		return
	}

	if err := h.columnService.DeleteColumn(uint(columnID)); err != nil {
		if err == services.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
