package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"` // 角色名称，如：admin, editor, viewer
	DisplayName string `gorm:"type:varchar(100);not null" json:"display_name"`     // 显示名称，如：管理员
	Description string `gorm:"type:varchar(255)" json:"description"`               // 角色描述
	Status      int    `gorm:"default:1" json:"status"`                            // 1: 启用, 0: 禁用

	// 关联关系
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}
