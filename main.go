package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tailings-dam/internal/handler"
	"tailings-dam/internal/service"
	"tailings-dam/internal/store"
)

func main() {
	dbPath := "tailings-dam.db"
	if env := os.Getenv("DB_PATH"); env != "" {
		dbPath = env
	}

	st, err := store.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to init store: %v", err)
	}
	defer st.Close()

	// 检查数据库连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := st.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 初始化服务层
	damService := service.NewDamService(st)
	monitoringService := service.NewMonitoringService(st)
	readingService := service.NewReadingService(st)
	alertService := service.NewAlertService(st)
	inspectionService := service.NewInspectionService(st)
	drainageService := service.NewDrainageService(st)

	// 初始化 HTTP 处理器
	h := handler.NewHandler(
		damService,
		monitoringService,
		readingService,
		alertService,
		inspectionService,
		drainageService,
	)

	// 配置 HTTP 服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      h.Routes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器（非阻塞）
	go func() {
		log.Printf("尾矿库安全监测系统启动，监听端口 :%s", port)
		log.Printf("数据库路径: %s", dbPath)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}

	log.Println("服务器已关闭")
}
