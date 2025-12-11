package repository

import (
	"errors"

	"go_test/internal/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *model.Role) error
	GetByID(id uint) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Update(role *model.Role) error
	Delete(id uint) error
	List(offset, limit int) ([]*model.Role, int64, error)
	// 角色权限管理
	AssignPermissions(roleID uint, permissionIDs []uint) error
	RemovePermissions(roleID uint, permissionIDs []uint) error
	GetPermissions(roleID uint) ([]*model.Permission, error)
	// 用户角色管理
	AssignUsers(roleID uint, userIDs []uint) error
	RemoveUsers(roleID uint, userIDs []uint) error
	GetUsers(roleID uint) ([]*model.User, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Preload("Permissions").Preload("Users").First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

func (r *roleRepository) List(offset, limit int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	err := r.db.Model(&model.Role{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Permissions").Offset(offset).Limit(limit).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *roleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	var permissions []model.Permission
	if err := r.db.Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Permissions").Append(permissions)
}

func (r *roleRepository) RemovePermissions(roleID uint, permissionIDs []uint) error {
	var permissions []model.Permission
	if err := r.db.Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Permissions").Delete(permissions)
}

func (r *roleRepository) GetPermissions(roleID uint) ([]*model.Permission, error) {
	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return nil, err
	}

	var permissions []*model.Permission
	err := r.db.Model(&role).Association("Permissions").Find(&permissions)
	return permissions, err
}

func (r *roleRepository) AssignUsers(roleID uint, userIDs []uint) error {
	var users []model.User
	if err := r.db.Find(&users, userIDs).Error; err != nil {
		return err
	}

	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Users").Append(users)
}

func (r *roleRepository) RemoveUsers(roleID uint, userIDs []uint) error {
	var users []model.User
	if err := r.db.Find(&users, userIDs).Error; err != nil {
		return err
	}

	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Users").Delete(users)
}

func (r *roleRepository) GetUsers(roleID uint) ([]*model.User, error) {
	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return nil, err
	}

	var users []*model.User
	err := r.db.Model(&role).Association("Users").Find(&users)
	return users, err
}
