package models

import (
	"time"

	"gorm.io/gorm"
)

// ActivityLog 活动日志表
type ActivityLog struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint           `json:"user_id" gorm:"not null;index;comment:'执行操作的用户ID'"`
	Username    string         `json:"username" gorm:"size:50;not null;comment:'执行操作的用户名（冗余字段，优化查询）'"`
	ActionType  string         `json:"action_type" gorm:"size:50;not null;index;comment:'操作类型：create/update/delete/move/comment等'"`
	EntityType  string         `json:"entity_type" gorm:"size:50;not null;index;comment:'实体类型：board/task/comment/attachment等'"`
	EntityID    uint           `json:"entity_id" gorm:"not null;index;comment:'实体ID'"`
	BoardID     *uint          `json:"board_id" gorm:"index;comment:'关联的看板ID（用于看板级别查询）'"`
	TaskID      *uint          `json:"task_id" gorm:"index;comment:'关联的任务ID（用于任务级别查询）'"`
	ProjectID   *uint          `json:"project_id" gorm:"index;comment:'关联的项目ID'"`
	Description string         `json:"description" gorm:"type:text;not null;comment:'操作描述文本'"`
	Metadata    string         `json:"metadata" gorm:"type:json;comment:'额外的元数据（JSON格式）'"`
	IPAddress   string         `json:"ip_address" gorm:"size:45;comment:'操作者IP地址'"`
	UserAgent   string         `json:"user_agent" gorm:"size:255;comment:'用户代理字符串'"`
	CreatedAt   time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	User    User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Board   *Board   `json:"board,omitempty" gorm:"foreignKey:BoardID"`
	Task    *Task    `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// TableName 指定表名
func (ActivityLog) TableName() string {
	return "activity_logs"
}

// ActivityActionType 定义常用的操作类型
const (
	ActionCreate  = "create"
	ActionUpdate  = "update"
	ActionDelete  = "delete"
	ActionMove    = "move"
	ActionComment = "comment"
	ActionAttach  = "attach"
	ActionAssign  = "assign"
	ActionArchive = "archive"
	ActionRestore = "restore"
)

// ActivityEntityType 定义常用的实体类型
const (
	EntityBoard      = "board"
	EntityColumn     = "column"
	EntityTask       = "task"
	EntityComment    = "comment"
	EntityAttachment = "attachment"
	EntityLabel      = "label"
	EntityProject    = "project"
)
