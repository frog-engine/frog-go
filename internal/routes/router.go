/**
 * Router Configuration
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package routes

import (
  "database/sql"

  "github.com/frog-engine/frog-go/config"
  "github.com/frog-engine/frog-go/internal/cache"
  "github.com/frog-engine/frog-go/internal/handlers"
  "github.com/frog-engine/frog-go/internal/middleware"
  "github.com/frog-engine/frog-go/internal/repositories"
  "github.com/frog-engine/frog-go/internal/services"
  "github.com/frog-engine/frog-go/internal/tools"
  "github.com/frog-engine/frog-go/pkg/logger"
  "github.com/frog-engine/frog-go/pkg/response"

  "github.com/gofiber/fiber/v3"
)

func SetupRouter(app *fiber.App) {
  // 加载配置并连接数据库
  config := config.GetConfig()
  db := initSqlDB(config)

  // 主页路由
  app.Get("/", func(c fiber.Ctx) error {
    return response.Success(c, map[string]string{
      "name":        "欢迎来到 frog-go Image Transcoding Service",
      "version":     "1.0",
      "description": "High Performance Image Transcoding API",
      "author":      "jarryli@gmail.com",
    })
  })

  // 健康检查路由
  app.Get("/health", func(c fiber.Ctx) error {
    return c.SendStatus(fiber.StatusOK)
  })

  // hello 路由
  app.Get("/hello", middleware.Chain(
    func(c fiber.Ctx) error {
      return response.Success(c, map[string]string{
        "name": "Hello World！你好！",
      })
    },
    middleware.RequestLoggerMiddleware,
  ))

  initUserRouter(app, db)
  initImageRouter(app, db)

}

func initSqlDB(config *config.Config) *sql.DB {
  databaseConfig := (*repositories.DatabaseConfig)(&config.Database)
  db, err := repositories.ConnectDatabase(databaseConfig)
  logger.Println("router->SetupRouter:", "db:", db)
  if err != nil {
    logger.Fatal("Error connecting to the database:", err)
  }
  return db
}

func initImageRouter(app *fiber.App, _ *sql.DB) {
  // 初始化依赖项
  imageCache := cache.NewImageCache()
  imageTools := tools.NewImageTools()
  transcodingService := services.NewTranscodingService(imageCache, imageTools)
  imageHandler := handlers.NewImageHandler(transcodingService)
  // Image 路由处理
  app.Get("/api/image/process", middleware.APIChain(imageHandler.ProcessImage, "/api/image"))
}

func initUserRouter(app *fiber.App, db *sql.DB) {
  // 初始化用户服务和处理器
  userRepo := repositories.NewSQLUserRepository(db)
  userService := services.NewUserService(userRepo)
  userHandler := handlers.NewUserHandler(userService)
  // User 路由处理，注意fiber关注顺序，最长的路径放前面
  app.Get("/api/user/list", middleware.APIChain(userHandler.FindPagedUsers, "/api/user"))
  app.Get("/api/user/group", middleware.APIChain(userHandler.GroupByHandler, "/api/user"))
  app.Post("/api/user", middleware.APIChain(userHandler.CreateUser, "/api/user"))
  app.Put("/api/user", middleware.APIChain(userHandler.UpdateUser, "/api/user"))
  app.Get("/api/users", middleware.APIChain(userHandler.GetAllUsers, "/api/user"))
  app.Get("/api/user/:id", middleware.APIChain(userHandler.GetUserById, "/api/user"))
  app.Delete("/api/user/:id", middleware.APIChain(userHandler.DeleteUser, "/api/user"))
}
