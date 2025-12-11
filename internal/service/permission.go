package service

import (
	"errors"

	"go_test/internal/model"
	"go_test/internal/repository"

	"gorm.io/gorm"
)

type PermissionService interface {
	CreatePermission(name, displayName, description, resource, action string) (*model.Permission, error)
	GetPermissionByID(id uint) (*model.Permission, error)
	GetPermissionByName(name string) (*model.Permission, error)
	UpdatePermission(id uint, displayName, description string, status int) (*model.Permission, error)
	DeletePermission(id uint) error
	ListPermissions(page, pageSize int) ([]*model.Permission, int64, error)
	GetPermissionByResourceAndAction(resource, action string) (*model.Permission, error)
}

type permissionService struct {
	permissionRepo repository.PermissionRepository
}

func NewPermissionService(permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionService{permissionRepo: permissionRepo}
}

func (s *permissionService) CreatePermission(name, displayName, description, resource, action string) (*model.Permission, error) {
	// 检查权限名称是否已存在
	existingPermission, err := s.permissionRepo.GetByName(name)
	if err == nil && existingPermission != nil {
		return nil, errors.New("权限名称已存在")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	permission := &model.Permission{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Resource:    resource,
		Action:      action,
		Status:      1,
	}

	err = s.permissionRepo.Create(permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *permissionService) GetPermissionByID(id uint) (*model.Permission, error) {
	return s.permissionRepo.GetByID(id)
}

func (s *permissionService) GetPermissionByName(name string) (*model.Permission, error) {
	return s.permissionRepo.GetByName(name)
}

func (s *permissionService) UpdatePermission(id uint, displayName, description string, status int) (*model.Permission, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if displayName != "" {
		permission.DisplayName = displayName
	}
	if description != "" {
		permission.Description = description
	}
	if status >= 0 {
		permission.Status = status
	}

	err = s.permissionRepo.Update(permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *permissionService) DeletePermission(id uint) error {
	return s.permissionRepo.Delete(id)
}

func (s *permissionService) ListPermissions(page, pageSize int) ([]*model.Permission, int64, error) {
	offset := (page - 1) * pageSize
	return s.permissionRepo.List(offset, pageSize)
}

func (s *permissionService) GetPermissionByResourceAndAction(resource, action string) (*model.Permission, error) {
	return s.permissionRepo.GetByResourceAndAction(resource, action)
}
