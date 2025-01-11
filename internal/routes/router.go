/**
 * Router Configuration
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package routes

import (
  "log"
  "net/http"

  "frog-go/config"
  "frog-go/internal/cache"
  "frog-go/internal/handlers"
  "frog-go/internal/middleware"
  "frog-go/internal/repositories"
  "frog-go/internal/services"
  "frog-go/internal/tools"
  "frog-go/pkg/response"

  "github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
  router := mux.NewRouter()

  // Load configuration and connect to the database
  config := config.GetConfig()
  databaseConfig := (*repositories.DatabaseConfig)(&config.Database)
  db, err := repositories.ConnectDatabase(databaseConfig)
  if err != nil {
    log.Fatal("Error connecting to the database:", err)
  }

  // Initialize dependencies
  imageCache := cache.NewImageCache()
  imageTools := tools.NewImageTools()
  transcodingService := services.NewTranscodingService(imageCache, imageTools)
  imageHandler := handlers.NewImageHandler(transcodingService)

  // Initialize user service and handler
  userRepo := repositories.NewSQLUserRepository(db)
  userService := services.NewUserService(userRepo)
  userHandler := handlers.NewUserHandler(userService)

  router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    response.Success(w, map[string]string{
      "name":        "欢迎来到 frog-go Image Transcoding Service",
      "version":     "1.0",
      "description": "High Performance Image Transcoding API",
      "author":      "jarryli@gmail.com",
    })
  }).Methods("GET", "OPTIONS")

  router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
  }).Methods("GET", "OPTIONS")

  router.Handle("/hello",
    middleware.Chain(
      http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        response.Success(w, map[string]string{
          "name": "Hello World！你好！",
        })
      }),
      middleware.RequestLoggerMiddleware,
    ),
  ).Methods("GET", "OPTIONS")

  /* 针对 user 的增删改查接口 */

  // 创建用户
  // 参数: user对象，包含用户的基本信息，如 name, email, phone 等，参见User Model
  // 请求示例：POST /api/user { "name": "张三", "email": "zhangsan@example.com", "phone": 13912345678 }
  router.Handle(
    "/api/user",
    middleware.APIChain(userHandler.CreateUser, "/api/user"),
  ).Methods("POST", "OPTIONS")

  // 查询用户 by ID
  // 参数: id，用户的唯一标识符，整型值
  // 请求示例：GET /api/user/12
  router.Handle(
    "/api/user/{id:[0-9]+}",
    middleware.APIChain(userHandler.GetUserById, "/api/user"),
  ).Methods("GET", "OPTIONS")

  // 分页查询用户列表
  // 参数: page (页码), size (每页数量), fields (返回的字段), condition (查询条件)
  // 请求示例：GET /api/user/list?page=1&size=10&fields=name,email&condition=id>8
  router.Handle(
    "/api/user/list",
    middleware.APIChain(userHandler.FindPagedUsers, "/api/user"),
  ).Methods("GET", "OPTIONS")

  // 查询所有用户
  // 无需额外参数，返回所有用户信息
  // 请求示例：GET /api/users
  router.Handle(
    "/api/users",
    middleware.APIChain(userHandler.GetAllUsers, "/api/user"),
  ).Methods("GET", "OPTIONS")

  // 更新用户信息
  // 参数: user对象，包含需要更新的用户信息，如 id, name, email 等
  // 请求示例：PUT /api/user { "id": 123, "name": "John Updated", "email": "updated.john.doe@example.com" }
  router.Handle(
    "/api/user",
    middleware.APIChain(userHandler.UpdateUser, "/api/user"),
  ).Methods("PUT", "OPTIONS")

  // 删除用户
  // 参数: id，待删除用户的唯一标识符，整型值
  // 请求示例：DELETE /api/user/123
  router.Handle(
    "/api/user/{id:[0-9]+}",
    middleware.APIChain(userHandler.DeleteUser, "/api/user"),
  ).Methods("DELETE", "OPTIONS")

  // 根据字段分组统计用户数据
  // 参数: field，指定字段进行分组统计（如：age, gender等）
  // 请求示例：GET /api/user/group?field=age
  router.Handle(
    "/api/user/group",
    middleware.APIChain(userHandler.GroupByHandler, "/api/user"),
  ).Methods("GET", "OPTIONS")

  /* 针对image的处理接口 */

  // 图像处理
  // 参数: image (图像文件ID或URL)，action (处理操作，如 resize, rotate 等)
  // 可选参数: width, height (针对resize操作)
  // 请求示例：GET /api/image/process?image=img123&action=resize&width=200&height=300
  router.Handle(
    "/api/image/process",
    middleware.APIChain(imageHandler.ProcessImage, "/api/image"),
  ).Methods("GET", "OPTIONS")

  return router
}
