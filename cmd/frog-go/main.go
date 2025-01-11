/**
 * Frog - Fast Image Processing Service
 *
 * @author jarryli@gmail.com
 * @date 2024-03-20
 */

package main

import (
  "flag"
  "log"
  "net/http"
  "strconv"
  "time"

  "frog-go/config"
  "frog-go/internal/routes"
  "frog-go/pkg/logger"
)

func main() {
  defaultEnv := "test"
  // 解析命令行参数
  env := flag.String("env", defaultEnv, "Application environment (production, test)")
  flag.Parse()
  // 初始化配置
  config.Init(env)

  // 获取配置
  cfg := config.GetConfig()
  log.Println("config:\r\n", cfg)

  // 初始化日志
  logger.Init()

  // 初始化路由
  router := routes.SetupRouter()

  readTimeout, _ := strconv.Atoi(cfg.Server.ReadTimeout)
  writeTimeout, _ := strconv.Atoi(cfg.Server.WriteTimeout)

  // 启动HTTP服务
  server := &http.Server{
    Addr:           cfg.Server.Addr + ":" + cfg.Server.Port,
    Handler:        router,
    ReadTimeout:    time.Duration(readTimeout) * time.Second,
    WriteTimeout:   time.Duration(writeTimeout) * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  logger.Info("Server starting on :" + server.Addr)

  if err := server.ListenAndServe(); err != nil {
    log.Fatalf("Server failed to start: %v", err)
  }
}
