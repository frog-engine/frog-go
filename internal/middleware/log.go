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
  "log"
  "net/http"
)

// RequestLoggerMiddleware 是一个中间件，用于记录每个请求的详细信息
func RequestLoggerMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Println("RequestLoggerMiddleware:")
    PrintHanderInfo(next)
    LogRequestDetails(r) // 打印request详情
    next.ServeHTTP(w, r) // 调用下一个处理程序
  })
}

// logRequestDetails logs the details of the incoming request
func LogRequestDetails(r *http.Request) {
  log.Printf("Request Method: %s, Request URL: %s, Request Headers: %v", r.Method, r.URL, r.Header)

  // Log the request body if needed (be cautious with large bodies)
  if r.Method == http.MethodPost || r.Method == http.MethodPut {
    bodyBytes, err := io.ReadAll(r.Body)
    if err == nil {
      log.Printf("Request Body: %s", bodyBytes)
      r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
    }
  }
}
