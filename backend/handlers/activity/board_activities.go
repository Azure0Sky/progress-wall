package activity

import (
	"math"
	"net/http"
	"progress-wall-backend/dto"
	"progress-wall-backend/models"
	"progress-wall-backend/repository"
	"progress-wall-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BoardActivitiesHandler 看板活动日志处理器
type BoardActivitiesHandler struct {
	activityService *services.ActivityLogService
}

// NewBoardActivitiesHandler 创建看板活动日志处理器
func NewBoardActivitiesHandler(db *gorm.DB) *BoardActivitiesHandler {
	activityRepo := repository.NewActivityLogRepository(db)
	activityService := services.NewActivityLogService(activityRepo, db)

	return &BoardActivitiesHandler{
		activityService: activityService,
	}
}

// GetBoardActivities 获取看板活动日志
// GET /api/boards/:board_id/activities
func (h *BoardActivitiesHandler) GetBoardActivities(c *gin.Context) {
	// 获取当前登录用户ID
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 获取看板ID
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的看板ID"})
		return
	}

	// 验证用户是否有权访问该看板
	if err := h.activityService.VerifyBoardAccess(userID, uint(boardID)); err != nil {
		switch err {
		case services.ErrBoardNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "看板不存在"})
		case services.ErrAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问该看板"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "验证权限失败"})
		}
		return
	}

	// 获取分页参数
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分页参数"})
		return
	}

	// 设置默认值
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	// 获取活动日志
	logs, total, err := h.activityService.GetBoardActivities(uint(boardID), query.Page, query.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取活动日志失败"})
		return
	}

	// 转换为响应格式
	activities := make([]dto.ActivityLogResponse, len(logs))
	for i, log := range logs {
		activities[i] = convertToActivityLogResponse(log)
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	// 返回响应
	c.JSON(http.StatusOK, dto.ActivityLogListResponse{
		Data:       activities,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	})
}

// convertToActivityLogResponse 转换模型为响应DTO
func convertToActivityLogResponse(log models.ActivityLog) dto.ActivityLogResponse {
	return dto.ActivityLogResponse{
		ID:          log.ID,
		UserID:      log.UserID,
		Username:    log.Username,
		ActionType:  log.ActionType,
		EntityType:  log.EntityType,
		EntityID:    log.EntityID,
		Description: log.Description,
		CreatedAt:   log.CreatedAt,
	}
}
