package dto

import "time"

// ActivityLogResponse 活动日志响应
type ActivityLogResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	ActionType  string    `json:"action_type"`
	EntityType  string    `json:"entity_type"`
	EntityID    uint      `json:"entity_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// ActivityLogListResponse 活动日志列表响应
type ActivityLogListResponse struct {
	Data       []ActivityLogResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// PaginationQuery 分页查询参数
type PaginationQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"limit" binding:"omitempty,min=1,max=100"`
}
