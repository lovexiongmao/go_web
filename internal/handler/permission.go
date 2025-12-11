package handler

import (
	"strconv"
	"time"

	"go_test/internal/service"
	"go_test/internal/util"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService service.PermissionService
}

func NewPermissionHandler(permissionService service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService}
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required" example:"user:create"`  // 权限名称
	DisplayName string `json:"display_name" binding:"required" example:"创建用户"` // 显示名称
	Description string `json:"description" example:"允许创建用户"`                   // 权限描述
	Resource    string `json:"resource" binding:"required" example:"user"`     // 资源类型
	Action      string `json:"action" binding:"required" example:"create"`     // 操作类型
}

type UpdatePermissionRequest struct {
	DisplayName *string `json:"display_name" example:"创建用户"`  // 显示名称（可选）
	Description *string `json:"description" example:"允许创建用户"` // 权限描述（可选）
	Status      *int    `json:"status" example:"1"`           // 状态：1-启用，0-禁用（可选）
}

// PermissionResponse 权限响应结构体（用于 Swagger 文档）
type PermissionResponse struct {
	ID          uint      `json:"id" example:"1"`                            // 权限ID
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"` // 创建时间
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"` // 更新时间
	Name        string    `json:"name" example:"user:create"`                // 权限名称
	DisplayName string    `json:"display_name" example:"创建用户"`               // 显示名称
	Description string    `json:"description" example:"允许创建用户"`              // 权限描述
	Resource    string    `json:"resource" example:"user"`                   // 资源类型
	Action      string    `json:"action" example:"create"`                   // 操作类型
	Status      int       `json:"status" example:"1"`                        // 状态：1-启用，0-禁用
}

// CreatePermission 创建权限
// @Summary      创建权限
// @Description  创建一个新权限
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Param        permission  body      CreatePermissionRequest  true  "权限信息"
// @Success      201          {object}  util.Response{data=PermissionResponse}
// @Failure      400          {object}  util.Response
// @Failure      500          {object}  util.Response
// @Router       /permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequestWithError(c, "请求参数错误", err)
		return
	}

	permission, err := h.permissionService.CreatePermission(req.Name, req.DisplayName, req.Description, req.Resource, req.Action)
	if err != nil {
		util.InternalServerErrorWithError(c, "创建权限失败", err)
		return
	}

	util.CreatedWithMessage(c, "权限创建成功", permission)
}

// GetPermission 获取权限详情
// @Summary      获取权限详情
// @Description  根据权限ID获取权限详细信息
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "权限ID"
// @Success      200  {object}  util.Response{data=PermissionResponse}
// @Failure      400  {object}  util.Response
// @Failure      404  {object}  util.Response
// @Router       /permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的权限ID")
		return
	}

	permission, err := h.permissionService.GetPermissionByID(uint(id))
	if err != nil {
		util.NotFound(c, "权限不存在")
		return
	}

	util.Success(c, permission)
}

// UpdatePermission 更新权限
// @Summary      更新权限
// @Description  根据权限ID更新权限信息（支持部分更新）
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Param        id          path      int                    true  "权限ID"
// @Param        permission  body      UpdatePermissionRequest true  "权限信息"
// @Success      200         {object}  util.Response{data=PermissionResponse}
// @Failure      400         {object}  util.Response
// @Failure      500         {object}  util.Response
// @Router       /permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的权限ID")
		return
	}

	var req UpdatePermissionRequest
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

	permission, err := h.permissionService.UpdatePermission(uint(id), displayName, description, status)
	if err != nil {
		util.InternalServerErrorWithError(c, "更新权限失败", err)
		return
	}

	util.SuccessWithMessage(c, "权限更新成功", permission)
}

// DeletePermission 删除权限
// @Summary      删除权限
// @Description  根据权限ID删除权限
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "权限ID"
// @Success      200  {object}  util.Response
// @Failure      400  {object}  util.Response
// @Failure      500  {object}  util.Response
// @Router       /permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的权限ID")
		return
	}

	err = h.permissionService.DeletePermission(uint(id))
	if err != nil {
		util.InternalServerErrorWithError(c, "删除权限失败", err)
		return
	}

	util.SuccessWithMessage(c, "权限删除成功", nil)
}

// ListPermissions 获取权限列表
// @Summary      获取权限列表
// @Description  分页获取权限列表
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "页码"  default(1)
// @Param        page_size query     int  false  "每页数量"  default(10)
// @Success      200       {object}  util.Response{data=[]PermissionResponse}
// @Failure      400       {object}  util.Response
// @Router       /permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	permissions, total, err := h.permissionService.ListPermissions(page, pageSize)
	if err != nil {
		util.InternalServerErrorWithError(c, "获取权限列表失败", err)
		return
	}

	util.SuccessWithPagination(c, permissions, total, page, pageSize)
}
