package repository

import (
	"errors"

	"go_test/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*model.User, int64, error)
	// 用户角色管理
	AssignRoles(userID uint, roleIDs []uint) error
	RemoveRoles(userID uint, roleIDs []uint) error
	GetRoles(userID uint) ([]*model.Role, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").Preload("Roles.Permissions").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		// 如果是记录不存在，返回错误以便上层判断
		// 如果是其他错误，也返回错误
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) List(offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	err := r.db.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Roles").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) AssignRoles(userID uint, roleIDs []uint) error {
	var roles []model.Role
	if err := r.db.Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	return r.db.Model(&user).Association("Roles").Append(roles)
}

func (r *userRepository) RemoveRoles(userID uint, roleIDs []uint) error {
	var roles []model.Role
	if err := r.db.Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	return r.db.Model(&user).Association("Roles").Delete(roles)
}

func (r *userRepository) GetRoles(userID uint) ([]*model.Role, error) {
	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	var roles []*model.Role
	err := r.db.Model(&user).Association("Roles").Find(&roles)
	return roles, err
}
