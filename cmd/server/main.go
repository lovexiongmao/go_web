// @title           Go Test API
// @version         1.0
// @description     这是一个 Go Web项目的 RESTful API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes   http https

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_test/docs/swagger" // Swagger 文档
	"go_test/internal/config"
	"go_test/internal/logger"
	"go_test/pkg/dig"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 初始化 Swagger 文档（导入 docs/swagger 包会自动执行 init 函数）
	_ = swagger.SwaggerInfo

	// 创建依赖注入容器
	container := dig.NewContainer()

	// 启动服务
	if err := container.Invoke(startServer); err != nil {
		panic(fmt.Sprintf("启动服务失败: %v", err))
	}
}

func startServer(
	cfg *config.Config,
	log *logger.Logger,
	db *gorm.DB,
	r *gin.Engine,
) error {
	// Gin模式已在router.SetupRouter中设置

	// 初始化数据库表（只在配置允许时执行）
	// 生产环境应设置 DB_AUTO_MIGRATE=false，使用专门的迁移工具
	if cfg.Database.AutoMigrate {
		if err := dig.InitializeDatabase(db, cfg); err != nil {
			return fmt.Errorf("初始化数据库失败: %v", err)
		}
		log.Info("数据库表初始化成功！")
	} else {
		log.Info("跳过数据库自动迁移（生产环境模式，请使用专门的迁移工具）")
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: r,
	}

	// 启动服务器（在goroutine中）
	go func() {
		log.Infof("服务器启动在 http://%s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭服务器...")

	// 设置5秒的超时时间用于关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}

	log.Info("服务器已退出！")

	return nil
}
