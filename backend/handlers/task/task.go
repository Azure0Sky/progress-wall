package task

import (
	"net/http"
	"strconv"
	"time"

	"progress-wall-backend/models"
	"progress-wall-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	taskService *services.TaskService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{
		taskService: services.NewTaskService(db),
	}
}

// GetTask 获取单个任务
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	task, err := h.taskService.GetTaskByID(uint(taskID))
	if err != nil {
		if err == services.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetTasks 获取列的所有任务
func (h *TaskHandler) GetTasks(c *gin.Context) {
	columnID, err := strconv.ParseUint(c.Param("columnId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的列ID"})
		return
	}

	tasks, err := h.taskService.GetTasksByColumnID(uint(columnID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// CreateTask 创建任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无法获取用户信息"})
		return
	}

	columnID, err := strconv.ParseUint(c.Param("columnId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的列ID"})
		return
	}

	var req struct {
		Title          string               `json:"title" binding:"required"`
		Description    string               `json:"description"`
		Priority       *models.TaskPriority `json:"priority"`
		DueDate        *time.Time           `json:"due_date"`
		StartDate      *time.Time           `json:"start_date"`
		EstimatedHours *float64             `json:"estimated_hours"`
		AssigneeID     *uint                `json:"assignee_id"`
		ProjectID      uint                 `json:"project_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	priority := models.TaskPriorityMedium
	if req.Priority != nil {
		priority = *req.Priority
	}

	task := &models.Task{
		Title:          req.Title,
		Description:    req.Description,
		Priority:       priority,
		Status:         models.TaskStatusTodo,
		ColumnID:       uint(columnID),
		CreatorID:      userID,
		AssigneeID:     req.AssigneeID,
		ProjectID:      req.ProjectID,
		DueDate:        req.DueDate,
		StartDate:      req.StartDate,
		EstimatedHours: req.EstimatedHours,
	}

	if err := h.taskService.CreateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTask 更新任务
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var req struct {
		Title          *string              `json:"title"`
		Description    *string              `json:"description"`
		Priority       *models.TaskPriority `json:"priority"`
		Status         *models.TaskStatus   `json:"status"`
		DueDate        *time.Time           `json:"due_date"`
		StartDate      *time.Time           `json:"start_date"`
		EndDate        *time.Time           `json:"end_date"`
		EstimatedHours *float64             `json:"estimated_hours"`
		ActualHours    *float64             `json:"actual_hours"`
		AssigneeID     *uint                `json:"assignee_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.DueDate != nil {
		updates["due_date"] = *req.DueDate
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = *req.EndDate
	}
	if req.EstimatedHours != nil {
		updates["estimated_hours"] = *req.EstimatedHours
	}
	if req.ActualHours != nil {
		updates["actual_hours"] = *req.ActualHours
	}
	if req.AssigneeID != nil {
		updates["assignee_id"] = *req.AssigneeID
	}

	if err := h.taskService.UpdateTask(uint(taskID), updates); err != nil {
		if err == services.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteTask 删除任务
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	if err := h.taskService.DeleteTask(uint(taskID)); err != nil {
		if err == services.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// MoveTask 移动任务（拖拽排序）
func (h *TaskHandler) MoveTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var req struct {
		NewColumnID uint `json:"newColumnId" binding:"required"`
		NewOrder    int  `json:"newOrder" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：需要 newColumnId 和 newOrder"})
		return
	}

	if err := h.taskService.MoveTask(uint(taskID), req.NewColumnID, req.NewOrder); err != nil {
		if err == services.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "移动成功"})
}
