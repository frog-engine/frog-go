# 开发规范

## 1. 项目概述
图片同步和异步转码系统使用Go语言开发，采用 Fiber 或 Gin 框架进行 Web 开发。

## 2. 项目结构

```plaintext
/project-root
│
├── /cmd/                  # 主程序入口
│   └── main.go
│
├── /internal/             # 内部模块，包含核心业务逻辑
│   ├── /handlers/         # 请求处理
│   ├── /models/           # 数据模型
│   ├── /services/         # 服务层，业务逻辑处理
│   ├── /repositories/     # 数据存储层，数据库操作
│   ├── /utils/            # 公共工具函数
│   ├── /middlewares/      # 中间件
│   └── /config/           # 配置文件
│
├── /pkg/                  # 公共库，其他项目也可以使用
│
├── /scripts/              # 脚本文件
│
├── /migrations/           # 数据库迁移文件
│
└── /docs/                 # 文档
```

## 3. 依赖管理
使用 Go Modules 管理项目依赖。
项目初始化时运行 go mod init <module-name> 来初始化模块。
使用 go get 安装依赖，定期运行 go mod tidy 清理不再使用的依赖。

## 4. 包和模块命名
包名：包名应尽量简洁并符合其功能，例如：
- handlers：处理请求
- models：定义数据结构
- services：实现业务逻辑
- repositories：数据库操作
- middlewares：中间件
- utils：公用工具类
文件名：文件名应简洁，能明确描述其功能，如 user_handler.go, product_service.go 等。

## 5. 命名规范
- 变量命名：使用驼峰命名法。例：userName, userAge。
- 常量命名：常量应全大写，单词之间使用下划线分隔。例如：MAX_VALUE。
- 结构体命名：结构体名应为名词，并采用大驼峰命名法。例：User, Product。
- 函数命名：函数名应为动词，并且首字母大写（若导出）或小写（若非导出）。例：GetUser(), createProduct()。

## 6. 代码风格
- 缩进：使用 Tab 进行缩进，避免使用空格。
- 行长度：每行代码的长度不超过 80 个字符，长行可分行处理。
- 空行：函数之间应保持空行，以提高可读性。
- 代码格式化：使用 go fmt 自动格式化代码，保持一致性。

## 7. 错误处理
错误处理是 Go 语言中的重要部分。每个函数或方法的返回值中应包括 error 类型，调用者需要显式检查错误。

示例：

```go
func GetUser(id int) (*User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %v", err)
    }
    return user, nil
}
```

### Gin 例子：
```go
func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := h.service.GetUserByID(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, user)
}
```

### Fiber 例子：
```go
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    id := c.Params("id")
    user, err := h.service.GetUserByID(id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(user)
}
```

### 8. 接口设计
接口设计要符合 RESTful API 标准，尽量简洁和一致。所有 API 路由应遵循以下格式：

```go
GET /users - 获取用户列表
GET /users/{id} - 获取单个用户
POST /users - 创建新用户
PUT /users/{id} - 更新用户信息
DELETE /users/{id} - 删除用户
示例：用户接口
```

### Gin 框架实现：
```go
func SetupRouter() *gin.Engine {
    r := gin.Default()
    userHandler := NewUserHandler()

    // 用户相关接口
    users := r.Group("/users")
    {
        users.GET("/", userHandler.GetAllUsers)
        users.GET("/:id", userHandler.GetUserByID)
        users.POST("/", userHandler.CreateUser)
        users.PUT("/:id", userHandler.UpdateUser)
        users.DELETE("/:id", userHandler.DeleteUser)
    }

    return r
}
```

### Fiber 框架实现：
```go
func SetupRouter(app *fiber.App) {
    userHandler := NewUserHandler()

    // 用户相关接口
    users := app.Group("/users")
    {
        users.Get("/", userHandler.GetAllUsers)
        users.Get("/:id", userHandler.GetUserByID)
        users.Post("/", userHandler.CreateUser)
        users.Put("/:id", userHandler.UpdateUser)
        users.Delete("/:id", userHandler.DeleteUser)
    }
}
```

## 9. 中间件
中间件在请求处理流程中起着重要作用。常见的中间件包括日志记录、认证、错误处理等。

### Gin 例子：
```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        latency := time.Since(start)
        log.Printf("Request: %s %s, Latency: %v", c.Request.Method, c.Request.URL, latency)
    }
}
```

### Fiber 例子：
```go
func Logger(c *fiber.Ctx) error {
    start := time.Now()
    err := c.Next()
    latency := time.Since(start)
    log.Printf("Request: %s %s, Latency: %v", c.Method(), c.Path(), latency)
    return err
}
```

## 10. 日志记录
日志是调试和监控的重要工具。建议使用 logrus 或 zap 等库进行日志记录。

- 日志级别：使用不同的日志级别（如 Info, Warn, Error）来记录不同的重要性级别的信息。
- 日志格式：统一使用结构化日志格式，方便后续的分析和搜索。

## 11. 性能优化
- 连接池：使用数据库连接池来提高数据库操作性能。
- 缓存：对于高频访问的数据，使用内存缓存（如 Redis）进行缓存。
- 异步处理：对于耗时的操作，可以使用 goroutine 异步执行，避免阻塞请求处理。

## 12. 单元测试
每个功能模块都应编写单元测试。 使用 testing 包进行测试，测试文件名应以 _test.go 结尾。

```go
// 这里使用 Gin 或 Fiber 提供的测试工具
func TestGetUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    handler := NewUserHandler(mockRepo)
    mockRepo.On("FindByID", 1).Return(&User{ID: 1, Name: "John"}, nil)
    req, err := http.NewRequest("GET", "/users/1", nil)
    if err != nil {
        t.Fatal(err)
    }
}
```

# 参考
- uber-go/guide 的中文翻译： https://github.com/xxjwxc/uber_go_guide_cn
- 官方Effective Go：https://go.dev/doc/effective_go
- Google Go 官方规范：https://google.github.io/styleguide/go/guide
- Google Go 编程规范：https://gocn.github.io/styleguide/
- Google 开发Wiki：https://go.dev/wiki/