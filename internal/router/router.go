package router

import (
	"go_test/docs/swagger" // Swagger 文档
	"go_test/internal/config"
	"go_test/internal/handler"
	"go_test/internal/util"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/dig"
)

// RouterParams 路由参数结构体，用于接收命名参数
type RouterParams struct {
	dig.In

	Config            *config.Config
	LoggerMiddleware  gin.HandlerFunc `name:"logger"`
	AuditMiddleware   gin.HandlerFunc `name:"audit"`
	UserHandler       *handler.UserHandler
	RoleHandler       *handler.RoleHandler
	PermissionHandler *handler.PermissionHandler
}

func SetupRouter(params RouterParams) *gin.Engine {
	cfg := params.Config
	loggerMiddleware := params.LoggerMiddleware
	auditMiddleware := params.AuditMiddleware
	userHandler := params.UserHandler
	roleHandler := params.RoleHandler
	permissionHandler := params.PermissionHandler
	// 在创建路由之前设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()

	// 全局中间件
	r.Use(loggerMiddleware)
	r.Use(auditMiddleware)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		util.SuccessWithMessage(c, "服务运行正常", gin.H{
			"status": "ok",
		})
	})

	// Swagger UI 文档
	// 确保 Swagger 文档被注册（导入 docs/swagger 包会自动执行 init 函数）
	_ = swagger.SwaggerInfo
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由组
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// 角色相关路由
		roles := api.Group("/roles")
		{
			roles.POST("", roleHandler.CreateRole)
			roles.GET("", roleHandler.ListRoles)
			roles.GET("/:id", roleHandler.GetRole)
			roles.PUT("/:id", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
			// 角色权限管理
			roles.POST("/:id/permissions", roleHandler.AssignPermissions)
			roles.DELETE("/:id/permissions", roleHandler.RemovePermissions)
			roles.GET("/:id/permissions", roleHandler.GetRolePermissions)
			// 角色用户管理
			roles.POST("/:id/users", roleHandler.AssignUsers)
			roles.DELETE("/:id/users", roleHandler.RemoveUsers)
			roles.GET("/:id/users", roleHandler.GetRoleUsers)
		}

		// 权限相关路由
		permissions := api.Group("/permissions")
		{
			permissions.POST("", permissionHandler.CreatePermission)
			permissions.GET("", permissionHandler.ListPermissions)
			permissions.GET("/:id", permissionHandler.GetPermission)
			permissions.PUT("/:id", permissionHandler.UpdatePermission)
			permissions.DELETE("/:id", permissionHandler.DeletePermission)
		}
	}

	return r
}
