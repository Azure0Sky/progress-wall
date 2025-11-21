package activity

import (
	"math"
	"net/http"
	"progress-wall-backend/dto"
	"progress-wall-backend/repository"
	"progress-wall-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TaskActivitiesHandler 任务活动日志处理器
type TaskActivitiesHandler struct {
	activityService *services.ActivityLogService
}

// NewTaskActivitiesHandler 创建任务活动日志处理器
func NewTaskActivitiesHandler(db *gorm.DB) *TaskActivitiesHandler {
	activityRepo := repository.NewActivityLogRepository(db)
	activityService := services.NewActivityLogService(activityRepo, db)

	return &TaskActivitiesHandler{
		activityService: activityService,
	}
}

// GetTaskActivities 获取任务活动日志
// GET /api/tasks/:task_id/activities
func (h *TaskActivitiesHandler) GetTaskActivities(c *gin.Context) {
	// 获取当前登录用户ID
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 获取任务ID
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	// 验证用户是否有权访问该任务
	if err := h.activityService.VerifyTaskAccess(userID, uint(taskID)); err != nil {
		switch err {
		case services.ErrTaskNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		case services.ErrAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问该任务"})
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
	logs, total, err := h.activityService.GetTaskActivities(uint(taskID), query.Page, query.PageSize)
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
