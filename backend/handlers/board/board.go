package board

import (
	"net/http"
	"strconv"

	"progress-wall-backend/models"
	"progress-wall-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BoardHandler 看板处理器
type BoardHandler struct {
	boardService *services.BoardService
}

// NewBoardHandler 创建看板处理器
func NewBoardHandler(db *gorm.DB) *BoardHandler {
	return &BoardHandler{
		boardService: services.NewBoardService(db),
	}
}

// GetBoard 获取单个看板（包含嵌套的列和任务）
func (h *BoardHandler) GetBoard(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("boardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的看板ID"})
		return
	}

	board, err := h.boardService.GetBoardByID(uint(boardID))
	if err != nil {
		if err == services.ErrBoardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, board)
}

// GetBoards 获取用户的所有看板
func (h *BoardHandler) GetBoards(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	boards, err := h.boardService.GetBoardsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"boards": boards})
}

// CreateBoard 创建看板
func (h *BoardHandler) CreateBoard(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
		ProjectID   uint   `json:"project_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	board := &models.Board{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		ProjectID:   req.ProjectID,
		OwnerID:     userID,
		Status:      models.BoardStatusActive,
	}

	if err := h.boardService.CreateBoard(board); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, board)
}

// UpdateBoard 更新看板
func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("boardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的看板ID"})
		return
	}

	var req struct {
		Name        *string             `json:"name"`
		Description *string             `json:"description"`
		Color       *string             `json:"color"`
		Status      *models.BoardStatus `json:"status"`
		Position    *int                `json:"position"`
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

	if err := h.boardService.UpdateBoard(uint(boardID), updates); err != nil {
		if err == services.ErrBoardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteBoard 删除看板
func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("boardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的看板ID"})
		return
	}

	if err := h.boardService.DeleteBoard(uint(boardID)); err != nil {
		if err == services.ErrBoardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
