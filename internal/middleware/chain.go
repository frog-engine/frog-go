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
  "reflect"
  "runtime"
)

// 事件链，传入handler以及中间件，依次执行
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {

  // 按逆序添加，顺序执行中间件，更符合习惯
  for i := len(middlewares) - 1; i >= 0; i-- {
    handler = middlewares[i](handler)
  }

  return handler
}

// 通用 API 中间件链，适用于不同的资源处理
func APIChain(handlerFunc func(http.ResponseWriter, *http.Request), permission string) http.Handler {
  return Chain(
    http.HandlerFunc(handlerFunc),
    CORSMiddleware,
    RequestLoggerMiddleware,
    SSOLoginMiddleware,
    TokenValidationMiddleware,
    PermissionMiddleware(permission),
  )
}

func PrintHanderInfo(handler http.Handler) {
  if handlerFunc, ok := handler.(http.HandlerFunc); ok {
    funcPointer := reflect.ValueOf(handlerFunc).Pointer()
    function := runtime.FuncForPC(funcPointer)
    if function != nil {
      // log.Println("Handler name:", function.Name())
    } else {
      log.Println("Could not retrieve function name.")
    }
  } else {
    log.Println("Not an http.HandlerFunc")
  }
}
