/**
 * log middleware
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package middleware

import (
  "bytes"
  "io"

  "github.com/frog-engine/frog-go/pkg/logger"
  "github.com/gofiber/fiber/v3"
)

// RequestLoggerMiddleware 是一个中间件，用于记录每个请求的详细信息
func RequestLoggerMiddleware(next fiber.Handler) fiber.Handler {
  return func(c fiber.Ctx) error {
    logger.Println("RequestLoggerMiddleware:")
    PrintHandlerInfo(next) // 打印传入的 handler 信息
    LogRequestDetails(c)   // 打印请求详情
    // 调用传入的 next 处理函数
    return next(c)
  }
}

// RequestLoggerMiddleware 是一个中间件，用于记录每个请求的详细信息
// func RequestLoggerMiddleware(c fiber.Ctx) error {
//   logger.Println("RequestLoggerMiddleware:")
//   LogRequestDetails(c) // 打印request详情
//   return c.Next()      // 调用下一个处理程序
// }

// logRequestDetails logs the details of the incoming request
func LogRequestDetails(c fiber.Ctx) {
  logger.Printf("Request Method: %s, Request URL: %s, Request Headers: %v", c.Method(), c.OriginalURL(), c.GetReqHeaders())

  // Log the request body if needed (be cautious with large bodies)
  if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut {
    bodyBytes := c.Body()                    // 获取请求体
    bodyReader := bytes.NewReader(bodyBytes) // 将[]byte转为io.Reader

    // 使用io.ReadAll读取请求体内容
    bodyContent, err := io.ReadAll(bodyReader)
    if err == nil {
      logger.Printf("Request Body: %s", bodyContent)
      // Rewind the body so it can be used later
      c.Request().SetBody(bodyBytes)
    }
  }
}
