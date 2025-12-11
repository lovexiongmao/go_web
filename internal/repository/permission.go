package repository

import (
	"errors"

	"go_test/internal/model"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(permission *model.Permission) error
	GetByID(id uint) (*model.Permission, error)
	GetByName(name string) (*model.Permission, error)
	Update(permission *model.Permission) error
	Delete(id uint) error
	List(offset, limit int) ([]*model.Permission, int64, error)
	GetByResourceAndAction(resource, action string) (*model.Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(permission *model.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) GetByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Preload("Roles").First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) GetByName(name string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Update(permission *model.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Permission{}, id).Error
}

func (r *permissionRepository) List(offset, limit int) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	err := r.db.Model(&model.Permission{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Roles").Offset(offset).Limit(limit).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

func (r *permissionRepository) GetByResourceAndAction(resource, action string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("resource = ? AND action = ?", resource, action).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}
