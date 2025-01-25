/**
 * Frog - Fast Image Processing Service
 *
 * @author jarryli@gmail.com
 * @date 2024-03-20
 */

package main

import (
  "flag"
  "strconv"
  "time"

  "github.com/frog-engine/frog-go/config"
  "github.com/frog-engine/frog-go/internal/routes"
  "github.com/frog-engine/frog-go/pkg/logger"
  frogsdk "github.com/frog-engine/frog-sdk"

  "github.com/gofiber/fiber/v3"
  "github.com/gofiber/fiber/v3/middleware/cors"
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
  logger.Info("Config loaded successfully. " + cfg.Server.Addr) // 使用 logger 记录信息
  logger.Println("config:\r\n", cfg)

  // 初始化日志
  logger.Init() // 初始化日志系统

  // 初始化 frogsdk
  frogsdk.Init(nil, "")
  defer frogsdk.Terminate()
  logger.Printf("%s initialized successfully. ", frogsdk.Name())

  // 将配置中的超时值转换为整型并检查错误
  readTimeout, err := strconv.Atoi(cfg.Server.ReadTimeout)
  if err != nil {
    logger.Errorf("Invalid ReadTimeout value: %v", err)
  }

  writeTimeout, err := strconv.Atoi(cfg.Server.WriteTimeout)
  if err != nil {
    logger.Errorf("Invalid WriteTimeout value: %v", err)
  }

  // 创建 Fiber 实例，去掉 MaxHeaderSize 和 Prefork
  app := fiber.New(fiber.Config{
    StrictRouting: true, // 启用严格路由匹配
    // CaseSensitive: true, // 路由区分大小写
    ReadTimeout:  time.Duration(readTimeout) * time.Second,
    WriteTimeout: time.Duration(writeTimeout) * time.Second,
  })

  // 使用 Fiber 的日志中间件，记录每个请求
  app.Use(func(c fiber.Ctx) error {
    // 这里的 c.Method() 和 c.Path() 都可以调用
    logger.Info("Request received: " + c.Method() + " " + c.Path()) // 使用自定义日志
    return c.Next()
  })

  // 配置 CORS 中间件
  app.Use(cors.New(cors.Config{
    AllowOrigins: []string{"*"},                                                        // 使用 []string
    AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"}, // 使用 []string
    AllowHeaders: []string{"Origin", "Content-Type", "Accept"},                         // 使用 []string
  }))

  // 初始化路由
  routes.SetupRouter(app)

  addr := cfg.Server.Addr + ":" + cfg.Server.Port
  logger.Info("Server starting on :" + addr)

  // 启动 HTTP 服务
  if err := app.Listen(addr); err != nil {
    logger.Errorf("Server failed to start: %v", err) // 使用 logger 记录错误
    logger.Fatalf("Server failed to start: %v", err)
  }
}
