package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRole 用户-角色关联表
type UserRole struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID uint `gorm:"not null;index;uniqueIndex:idx_user_role" json:"user_id"` // 用户ID
	RoleID uint `gorm:"not null;index;uniqueIndex:idx_user_role" json:"role_id"` // 角色ID

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}
