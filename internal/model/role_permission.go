package model

import (
	"time"

	"gorm.io/gorm"
)

// RolePermission 角色-权限关联表
type RolePermission struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	RoleID       uint `gorm:"not null;index;uniqueIndex:idx_role_permission" json:"role_id"`       // 角色ID
	PermissionID uint `gorm:"not null;index;uniqueIndex:idx_role_permission" json:"permission_id"` // 权限ID

	// 关联关系
	Role       Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
