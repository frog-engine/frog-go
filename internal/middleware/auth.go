/**
 * chain middleware
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package middleware

import (
  "log"
  "net/http"
)

// 权限校验中间件
func PermissionMiddleware(requiredPermission string) func(http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      log.Println("PermissionMiddleware:")
      PrintHanderInfo(next)
      // 获取用户权限（示例逻辑）
      userPermissions := getUserPermissions(r)
      if !hasPermission(userPermissions, requiredPermission) {
        http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
        return
      }
      next.ServeHTTP(w, r)
    })
  }
}

// 示例权限检查逻辑
func getUserPermissions(r *http.Request) []string {
  // 模拟从请求上下文中提取用户权限
  return []string{"read", "write"} // 示例权限
}

func hasPermission(userPermissions []string, requiredPermission string) bool {
  for _, perm := range userPermissions {
    if perm == requiredPermission {
      return true
    }
  }
  // TODO: return false
  return true
}

// Token 验证中间件
func TokenValidationMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Println("TokenValidationMiddleware:")
    PrintHanderInfo(next)
    // token := r.Header.Get("Authorization")
    // if token == "" {
    //   http.Error(w, "Unauthorized token", http.StatusUnauthorized)
    //   return
    // }
    // // 验证 Token（示例逻辑）
    // if token != "valid-token" {
    //   http.Error(w, "Forbidden", http.StatusForbidden)
    //   return
    // }
    next.ServeHTTP(w, r)
  })
}

// SSO 登录中间件
func SSOLoginMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Println("SSOLoginMiddleware:")
    PrintHanderInfo(next)
    // 检查 SSO 登录状态（示例逻辑）
    if !isLoggedIn(r) {
      http.Error(w, "SSO Required", http.StatusUnauthorized)
      return
    }
    next.ServeHTTP(w, r)
  })
}

func isLoggedIn(r *http.Request) bool {
  // 示例 SSO 检查逻辑
  // user := r.Header.Get("X-SSO-User")
  // return user != ""
  return true
}

// CORSMiddleware 设置 CORS 头
func CORSMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Println("CORSMiddleware received:", r.Method)

    // 处理预检请求（OPTIONS）
    if r.Method == http.MethodOptions {
      w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源，或根据需要修改为特定域名
      w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
      w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

      w.WriteHeader(http.StatusOK)
      return
    }

    // 对于其他请求，继续设置 CORS 头部
    w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    next.ServeHTTP(w, r)
  })
}
