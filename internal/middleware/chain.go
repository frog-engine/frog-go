/**
 * chain middleware
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package middleware

import (
  "reflect"
  "runtime"

  "github.com/frog-engine/frog-go/pkg/logger"
  "github.com/gofiber/fiber/v3"
)

// 事件链，传入handler以及中间件，依次执行
func Chain(handler fiber.Handler, middlewares ...func(fiber.Handler) fiber.Handler) fiber.Handler {
  // 按逆序添加，顺序执行中间件，更符合习惯
  for i := len(middlewares) - 1; i >= 0; i-- {
    handler = middlewares[i](handler)
  }

  return handler
}

// 通用 API 中间件链，适用于不同的资源处理
func APIChain(handler fiber.Handler, permission string) fiber.Handler {
  return Chain(
    handler,
    CORSMiddleware,
    RequestLoggerMiddleware,
    SSOLoginMiddleware,
    TokenValidationMiddleware,
    PermissionMiddleware(permission),
  )
}

func PrintHandlerInfo(handler fiber.Handler) {
  // 获取函数指针信息
  funcPointer := reflect.ValueOf(handler).Pointer() // 获取函数指针
  function := runtime.FuncForPC(funcPointer)        // 获取运行时函数信息

  if function != nil {
    logger.Println("Handler name:", function.Name()) // 打印函数名称
  } else {
    logger.Println("Could not retrieve function name.")
  }
}
