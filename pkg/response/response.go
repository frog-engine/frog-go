/**
 * HTTP Response Helper (Fiber version)
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package response

import (
	"github.com/frog-engine/frog-go/internal/models"
	"github.com/frog-engine/frog-go/pkg/code"

	"github.com/gofiber/fiber/v3"
)
 
 type Response struct {
   Code    int         `json:"code"`
   Message string      `json:"message"`
   Data    interface{} `json:"data,omitempty"`
 }
 
 func PaginationWrapper[T any](dataKey string, dataList []T, pagination models.Pagination) map[string]interface{} {
   result := map[string]interface{}{
     "page":       pagination.Page,
     "size":       pagination.Size,
     "total":      pagination.Total,
     "totalPages": pagination.TotalPages(),
     dataKey:      dataList,
   }
   return result
 }

 func Text(ctx fiber.Ctx, data interface{}) error {
  ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
  return ctx.Status(fiber.StatusOK).SendString(data.(string))
}
 
 func Success(ctx fiber.Ctx, data interface{}) error {
   return ctx.Status(fiber.StatusOK).JSON(Response{
     Code:    code.Success,
     Message: "success",
     Data:    data,
   })
 }
 
 func Error(ctx fiber.Ctx, httpCode int, message string) error {
   // Using Fiber's chainable method to set status code and JSON response
   return ctx.Status(httpCode).JSON(Response{
     Code:    code.Error,
     Message: message,
   })
 }
 