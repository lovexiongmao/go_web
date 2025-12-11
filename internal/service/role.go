package service

import (
	"errors"

	"go_test/internal/model"
	"go_test/internal/repository"

	"gorm.io/gorm"
)

type RoleService interface {
	CreateRole(name, displayName, description string) (*model.Role, error)
	GetRoleByID(id uint) (*model.Role, error)
	GetRoleByName(name string) (*model.Role, error)
	UpdateRole(id uint, displayName, description string, status int) (*model.Role, error)
	DeleteRole(id uint) error
	ListRoles(page, pageSize int) ([]*model.Role, int64, error)
	// 角色权限管理
	AssignPermissions(roleID uint, permissionIDs []uint) error
	RemovePermissions(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) ([]*model.Permission, error)
	// 用户角色管理
	AssignUsers(roleID uint, userIDs []uint) error
	RemoveUsers(roleID uint, userIDs []uint) error
	GetRoleUsers(roleID uint) ([]*model.User, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

func (s *roleService) CreateRole(name, displayName, description string) (*model.Role, error) {
	// 检查角色名称是否已存在
	existingRole, err := s.roleRepo.GetByName(name)
	if err == nil && existingRole != nil {
		return nil, errors.New("角色名称已存在")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role := &model.Role{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Status:      1,
	}

	err = s.roleRepo.Create(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) GetRoleByID(id uint) (*model.Role, error) {
	return s.roleRepo.GetByID(id)
}

func (s *roleService) GetRoleByName(name string) (*model.Role, error) {
	return s.roleRepo.GetByName(name)
}

func (s *roleService) UpdateRole(id uint, displayName, description string, status int) (*model.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if displayName != "" {
		role.DisplayName = displayName
	}
	if description != "" {
		role.Description = description
	}
	if status >= 0 {
		role.Status = status
	}

	err = s.roleRepo.Update(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) DeleteRole(id uint) error {
	return s.roleRepo.Delete(id)
}

func (s *roleService) ListRoles(page, pageSize int) ([]*model.Role, int64, error) {
	offset := (page - 1) * pageSize
	return s.roleRepo.List(offset, pageSize)
}

func (s *roleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}

func (s *roleService) RemovePermissions(roleID uint, permissionIDs []uint) error {
	return s.roleRepo.RemovePermissions(roleID, permissionIDs)
}

func (s *roleService) GetRolePermissions(roleID uint) ([]*model.Permission, error) {
	return s.roleRepo.GetPermissions(roleID)
}

func (s *roleService) AssignUsers(roleID uint, userIDs []uint) error {
	return s.roleRepo.AssignUsers(roleID, userIDs)
}

func (s *roleService) RemoveUsers(roleID uint, userIDs []uint) error {
	return s.roleRepo.RemoveUsers(roleID, userIDs)
}

func (s *roleService) GetRoleUsers(roleID uint) ([]*model.User, error) {
	return s.roleRepo.GetUsers(roleID)
}
