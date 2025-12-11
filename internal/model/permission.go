package model

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限模型
type Permission struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"` // 权限名称，如：user:create
	DisplayName string `gorm:"type:varchar(100);not null" json:"display_name"`     // 显示名称，如：创建用户
	Description string `gorm:"type:varchar(255)" json:"description"`               // 权限描述
	Resource    string `gorm:"type:varchar(50);not null;index" json:"resource"`    // 资源类型，如：user, role, permission
	Action      string `gorm:"type:varchar(50);not null;index" json:"action"`      // 操作类型，如：create, read, update, delete
	Status      int    `gorm:"default:1" json:"status"`                            // 1: 启用, 0: 禁用

	// 关联关系
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}
