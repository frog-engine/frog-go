/**
 * chain middleware
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package middleware

import (
  "github.com/frog-engine/frog-go/pkg/logger"
  "github.com/gofiber/fiber/v3"
)

func PermissionMiddleware(requiredPermission string) func(fiber.Handler) fiber.Handler {
  return func(next fiber.Handler) fiber.Handler {
    return func(c fiber.Ctx) error {
      logger.Println("PermissionMiddleware:")
      PrintHandlerInfo(next) // 打印下一步中间件的信息

      // 获取用户权限（示例逻辑）
      userPermissions := getUserPermissions(c)

      logger.Println("requiredPermission:", requiredPermission, "userPermissions:", userPermissions[0])
      // TODO: 检查权限
      // if !hasPermission(userPermissions, requiredPermission) {
      // return c.Status(fiber.StatusForbidden).SendString("Forbidden: insufficient permissions")
      // }

      // 调用下一个中间件或处理程序
      return next(c)
    }
  }
}

// 模拟从请求上下文中提取用户权限的函数
func getUserPermissions(c fiber.Ctx) []string {
  // 从请求上下文中提取用户ID（可以通过JWT、Session 或其他方式获取）
  userID := c.Locals("userID")

  // 模拟根据用户ID从缓存或数据库获取权限
  if userID == nil {
    // 如果找不到用户ID，可以返回一个空权限或做一些错误处理
    return []string{"read", "write"}
  }

  // 示例：如果用户是 "admin"，返回所有权限
  if userID == "admin" {
    return []string{"read", "write", "delete"}
  }

  // 示例：普通用户返回有限的权限
  return []string{"read", "write"}
}

func hasPermission(userPermissions []string, requiredPermission string) bool {
  for _, perm := range userPermissions {
    if perm == requiredPermission {
      return true
    }
  }
  return false
}

// Token 验证中间件
func TokenValidationMiddleware(next fiber.Handler) fiber.Handler {
  return func(c fiber.Ctx) error {
    logger.Println("TokenValidationMiddleware:")
    // 验证 Token（示例逻辑）
    token := c.Get("Authorization")
    // TODO: 跳过验证 Token
    if token == "" {
      token = "valid-token"
    }
    if token == "" || token != "valid-token" {
      return c.Status(fiber.StatusForbidden).SendString("Forbidden: Invalid token")
    }
    return next(c) // 继续调用下一个中间件
  }
}

// SSO 登录中间件
func SSOLoginMiddleware(next fiber.Handler) fiber.Handler {
  return func(c fiber.Ctx) error {
    logger.Println("SSOLoginMiddleware:")
    // 检查 SSO 登录状态（示例逻辑）
    if !isLoggedIn(c) {
      return c.Status(fiber.StatusUnauthorized).SendString("SSO Required")
    }
    return next(c) // 继续调用下一个中间件
  }
}

func isLoggedIn(c fiber.Ctx) bool {
  // 示例 SSO 检查逻辑
  user := c.Get("X-SSO-User")
  logger.Println("isLoggedIn user:", user)
  return true
}

// CORSMiddleware 设置 CORS 头
func CORSMiddleware(next fiber.Handler) fiber.Handler {
  return func(c fiber.Ctx) error {
    logger.Println("CORSMiddleware received:", c.Method())

    // 处理预检请求（OPTIONS）
    if c.Method() == fiber.MethodOptions {
      c.Set("Access-Control-Allow-Origin", "*") // 允许所有来源，或根据需要修改为特定域名
      c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
      c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
      return c.SendStatus(fiber.StatusOK)
    }

    // 对于其他请求，继续设置 CORS 头部
    c.Set("Access-Control-Allow-Origin", "*")
    c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    return next(c) // 继续调用下一个中间件
  }
}
