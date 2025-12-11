package handler

import (
	"strconv"
	"time"

	"go_test/internal/service"
	"go_test/internal/util"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required" example:"admin"`       // 角色名称
	DisplayName string `json:"display_name" binding:"required" example:"管理员"` // 显示名称
	Description string `json:"description" example:"系统管理员，拥有所有权限"`            // 角色描述
}

type UpdateRoleRequest struct {
	DisplayName *string `json:"display_name" example:"管理员"`  // 显示名称（可选）
	Description *string `json:"description" example:"系统管理员"` // 角色描述（可选）
	Status      *int    `json:"status" example:"1"`          // 状态：1-启用，0-禁用（可选）
}

type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required" example:"[1,2,3]"` // 权限ID列表
}

type AssignUsersRequest struct {
	UserIDs []uint `json:"user_ids" binding:"required" example:"[1,2,3]"` // 用户ID列表
}

// RoleResponse 角色响应结构体（用于 Swagger 文档）
type RoleResponse struct {
	ID          uint      `json:"id" example:"1"`                            // 角色ID
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"` // 创建时间
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"` // 更新时间
	Name        string    `json:"name" example:"admin"`                      // 角色名称
	DisplayName string    `json:"display_name" example:"管理员"`                // 显示名称
	Description string    `json:"description" example:"系统管理员"`               // 角色描述
	Status      int       `json:"status" example:"1"`                        // 状态：1-启用，0-禁用
}

// CreateRole 创建角色
// @Summary      创建角色
// @Description  创建一个新角色
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        role  body      CreateRoleRequest  true  "角色信息"
// @Success      201   {object}  util.Response{data=RoleResponse}
// @Failure      400   {object}  util.Response
// @Failure      500   {object}  util.Response
// @Router       /roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestWithError(c, "请求参数错误", err)
		return
	}

	role, err := h.roleService.CreateRole(req.Name, req.DisplayName, req.Description)
	if err != nil {
		util.InternalServerErrorWithError(c, "创建角色失败", err)
		return
	}

	util.CreatedWithMessage(c, "角色创建成功", role)
}

// GetRole 获取角色详情
// @Summary      获取角色详情
// @Description  根据角色ID获取角色详细信息
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "角色ID"
// @Success      200  {object}  util.Response{data=RoleResponse}
// @Failure      400  {object}  util.Response
// @Failure      404  {object}  util.Response
// @Router       /roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的角色ID")
		return
	}

	role, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		util.NotFound(c, "角色不存在")
		return
	}

	util.Success(c, role)
}

// UpdateRole 更新角色
// @Summary      更新角色
// @Description  根据角色ID更新角色信息（支持部分更新）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "角色ID"
// @Param        role  body      UpdateRoleRequest true  "角色信息"
// @Success      200   {object}  util.Response{data=RoleResponse}
// @Failure      400   {object}  util.Response
// @Failure      500   {object}  util.Response
// @Router       /roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的角色ID")
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestWithError(c, "请求参数错误", err)
		return
	}

	displayName := ""
	if req.DisplayName != nil {
		displayName = *req.DisplayName
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	status := -1
	if req.Status != nil {
		status = *req.Status
	}

	role, err := h.roleService.UpdateRole(uint(id), displayName, description, status)
	if err != nil {
		util.InternalServerErrorWithError(c, "更新角色失败", err)
		return
	}

	util.SuccessWithMessage(c, "角色更新成功", role)
}

// DeleteRole 删除角色
// @Summary      删除角色
// @Description  根据角色ID删除角色
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "角色ID"
// @Success      200  {object}  util.Response
// @Failure      400  {object}  util.Response
// @Failure      500  {object}  util.Response
// @Router       /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的角色ID")
		return
	}

	err = h.roleService.DeleteRole(uint(id))
	if err != nil {
		util.InternalServerErrorWithError(c, "删除角色失败", err)
		return
	}

	util.SuccessWithMessage(c, "角色删除成功", nil)
}

// ListRoles 获取角色列表
// @Summary      获取角色列表
// @Description  分页获取角色列表
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "页码"  default(1)
// @Param        page_size query     int  false  "每页数量"  default(10)
// @Success      200       {object}  util.Response{data=[]RoleResponse}
// @Failure      400       {object}  util.Response
// @Router       /roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	roles, total, err := h.roleService.ListRoles(page, pageSize)
	if err != nil {
		util.InternalServerErrorWithError(c, "获取角色列表失败", err)
		return
	}

	util.SuccessWithPagination(c, roles, total, page, pageSize)
}

// AssignPermissions 分配权限给角色
// @Summary      分配权限给角色
// @Description  为角色分配权限
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      int                      true  "角色ID"
// @Param        body  body      AssignPermissionsRequest true  "权限ID列表"
// @Success      200   {object}  util.Response
// @Failure      400   {object}  util.Response
// @Failure      500   {object}  util.Response
// @Router       /roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的角色ID")
		return
	}

	var req AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestWithError(c, "请求参数错误", err)
		return
	}

	err = h.roleService.AssignPermissions(uint(id), req.PermissionIDs)
	if err != nil {
		util.InternalServerErrorWithError(c, "分配权限失败", err)
		return
	}

	util.SuccessWithMessage(c, "权限分配成功", nil)
}

// RemovePermissions 移除角色的权限
// @Summary      移除角色的权限
// @Description  从角色中移除权限
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      int                      true  "角色ID"
// @Param        body  body      AssignPermissionsRequest true  "权限ID列表"
// @Success      200   {object}  util.Response
// @Failure      400   {object}  util.Response
// @Failure      500   {object}  util.Response
// @Router       /roles/{id}/permissions [delete]
func (h *RoleHandler) RemovePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的角色ID")
		return
	}

	var req AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestWithError(c, "请求参数错误", err)
		return
	}

	err = h.roleService.RemovePermissions(uint(id), req.PermissionIDs)
	if err != nil {
		util.InternalServerErrorWithError(c, "移除权限失败", err)
		return
	}

	util.SuccessWithMessage(c, "权限移除成功", nil)
}

// GetRolePermissions 获取角色的权限列表
// @Summary      获取角色的权限列表
// @Description  获取角色拥有的所有权限
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "角色ID"
// @Success      200  {object}  util.Response{data=[]PermissionResponse}
// @Failure      400  {object}  util.Response
// @Router       /roles/{id}/permissions [get]
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的角色ID")
		return
	}

	permissions, err := h.roleService.GetRolePermissions(uint(id))
	if err != nil {
		util.InternalServerErrorWithError(c, "获取权限列表失败", err)
		return
	}

	util.Success(c, permissions)
}

// AssignUsers 分配用户给角色
// @Summary      分配用户给角色
// @Description  为角色分配用户
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "角色ID"
// @Param        body  body      AssignUsersRequest true  "用户ID列表"
// @Success      200   {object}  util.Response
// @Failure      400   {object}  util.Response
// @Failure      500   {object}  util.Response
// @Router       /roles/{id}/users [post]
func (h *RoleHandler) AssignUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的角色ID")
		return
	}

	var req AssignUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestWithError(c, "请求参数错误", err)
		return
	}

	err = h.roleService.AssignUsers(uint(id), req.UserIDs)
	if err != nil {
		util.InternalServerErrorWithError(c, "分配用户失败", err)
		return
	}

	util.SuccessWithMessage(c, "用户分配成功", nil)
}

// RemoveUsers 移除角色的用户
// @Summary      移除角色的用户
// @Description  从角色中移除用户
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "角色ID"
// @Param        body  body      AssignUsersRequest true  "用户ID列表"
// @Success      200   {object}  util.Response
// @Failure      400   {object}  util.Response
// @Failure      500   {object}  util.Response
// @Router       /roles/{id}/users [delete]
func (h *RoleHandler) RemoveUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的角色ID")
		return
	}

	var req AssignUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestWithError(c, "请求参数错误", err)
		return
	}

	err = h.roleService.RemoveUsers(uint(id), req.UserIDs)
	if err != nil {
		util.InternalServerErrorWithError(c, "移除用户失败", err)
		return
	}

	util.SuccessWithMessage(c, "用户移除成功", nil)
}

// GetRoleUsers 获取角色的用户列表
// @Summary      获取角色的用户列表
// @Description  获取拥有该角色的所有用户
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "角色ID"
// @Success      200  {object}  util.Response{data=[]UserResponse}
// @Failure      400  {object}  util.Response
// @Router       /roles/{id}/users [get]
func (h *RoleHandler) GetRoleUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的角色ID")
		return
	}

	users, err := h.roleService.GetRoleUsers(uint(id))
	if err != nil {
		util.InternalServerErrorWithError(c, "获取用户列表失败", err)
		return
	}

	util.Success(c, users)
}
