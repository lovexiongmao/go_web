package dig

import (
	"go_test/internal/config"
	"go_test/internal/database"
	"go_test/internal/handler"
	"go_test/internal/logger"
	"go_test/internal/middleware"
	"go_test/internal/model"
	"go_test/internal/repository"
	"go_test/internal/router"
	"go_test/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

// Container 依赖注入容器
type Container struct {
	*dig.Container
}

// NewContainer 创建新的依赖注入容器
func NewContainer() *Container {
	c := dig.New()

	// 提供配置
	c.Provide(config.LoadConfig)

	// 提供请求日志Logger
	c.Provide(logger.NewLogger)

	// 提供审计日志Logger（使用命名参数区分）
	c.Provide(func(cfg *config.Config) *logger.Logger {
		return logger.NewAuditLogger(cfg)
	}, dig.Name("auditLogger"))

	// 提供数据库
	c.Provide(database.NewDatabase)

	// 提供Repository
	c.Provide(repository.NewUserRepository)

	// 提供Service
	c.Provide(service.NewUserService)

	// 提供Handler
	c.Provide(handler.NewUserHandler)

	// 提供中间件（使用命名参数区分）
	c.Provide(func(log *logger.Logger) gin.HandlerFunc {
		return middleware.LoggerMiddleware(log)
	}, dig.Name("logger"))

	// 审计中间件参数结构体
	type AuditMiddlewareParams struct {
		dig.In
		AuditLogger *logger.Logger `name:"auditLogger"`
	}
	c.Provide(func(params AuditMiddlewareParams) gin.HandlerFunc {
		return middleware.AuditMiddleware(params.AuditLogger)
	}, dig.Name("audit"))

	// 提供路由
	c.Provide(router.SetupRouter)

	return &Container{Container: c}
}

// InitializeDatabase 初始化数据库表
func InitializeDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		// 如果需要数据库层面的审计日志表，可以取消注释
		// &database.AuditLog{},
	)
}
