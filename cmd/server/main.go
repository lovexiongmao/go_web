package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_test/internal/config"
	"go_test/internal/logger"
	"go_test/pkg/dig"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
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
	// 设置Gin模式（必须在创建路由之前设置，但路由在依赖注入时已创建）
	// 因此需要在router中根据配置创建，或者在这里重新设置.
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库表
	if err := dig.InitializeDatabase(db); err != nil {
		return fmt.Errorf("初始化数据库失败: %v", err)
	}

	log.Info("数据库表初始化成功！")

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
