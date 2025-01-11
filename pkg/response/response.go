/**
 * HTTP Response Helper
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package response

import (
	"encoding/json"
	"frog-go/internal/models"
	"frog-go/pkg/code"
	"net/http"
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

func Success(w http.ResponseWriter, data interface{}) {
  JSON(w, http.StatusOK, Response{
    Code:    code.Success,
    Message: "success",
    Data:    data,
  })
}

func Error(w http.ResponseWriter, httpCode int, message string) {
  JSON(w, httpCode, Response{
    Code:    code.Error,
    Message: message,
  })
}

func JSON(w http.ResponseWriter, httpCode int, data interface{}) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(httpCode)
  json.NewEncoder(w).Encode(data)
}
