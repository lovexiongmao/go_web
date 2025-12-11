package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code    int         `json:"code"`            // HTTP 状态码
	Message string      `json:"message"`         // 响应消息
	Data    interface{} `json:"data,omitempty"`  // 响应数据（可选）
	Error   string      `json:"error,omitempty"` // 错误详情（可选）
}

// Success 成功响应（200 OK）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "操作成功",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（200 OK）带自定义消息
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Created 创建成功响应（201 Created）
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    data,
	})
}

// CreatedWithMessage 创建成功响应（201 Created）带自定义消息
func CreatedWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	})
}

// NoContent 无内容响应（204 No Content）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest 错误响应（400 Bad Request）
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
		Error:   message,
	})
}

// BadRequestWithError 错误响应（400 Bad Request）带详细错误信息
func BadRequestWithError(c *gin.Context, message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg = err.Error()
	}
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
		Error:   errorMsg,
	})
}

// Unauthorized 未授权响应（401 Unauthorized）
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "未授权，请先登录"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: message,
		Error:   message,
	})
}

// Forbidden 禁止访问响应（403 Forbidden）
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "权限不足，禁止访问"
	}
	c.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: message,
		Error:   message,
	})
}

// NotFound 资源不存在响应（404 Not Found）
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: message,
		Error:   message,
	})
}

// Conflict 冲突响应（409 Conflict）
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Code:    http.StatusConflict,
		Message: message,
		Error:   message,
	})
}

// UnprocessableEntity 无法处理响应（422 Unprocessable Entity）
func UnprocessableEntity(c *gin.Context, message string) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Error:   message,
	})
}

// InternalServerError 服务器错误响应（500 Internal Server Error）
func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
		Error:   message,
	})
}

// InternalServerErrorWithError 服务器错误响应（500 Internal Server Error）带详细错误信息
func InternalServerErrorWithError(c *gin.Context, message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg = err.Error()
	}
	if message == "" {
		message = "服务器内部错误"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
		Error:   errorMsg,
	})
}

// Error 自定义错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   message,
	})
}

// ErrorWithData 自定义错误响应带数据
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
		Error:   message,
	})
}

// SuccessWithPagination 成功响应（200 OK）带分页信息
func SuccessWithPagination(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "操作成功",
		Data: gin.H{
			"list":      data,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
