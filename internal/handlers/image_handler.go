/**
 * Image Processing Handler
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package handlers

import (
	"net/http"
	"strconv"

	"frog-go/internal/services"
	"frog-go/pkg/logger"
	"frog-go/pkg/response"
)

type ImageHandler struct {
  transcodingService *services.TranscodingService
}

func NewImageHandler(ts *services.TranscodingService) *ImageHandler {
  return &ImageHandler{
    transcodingService: ts,
  }
}

// ProcessImage 处理图片转码请求
func (h *ImageHandler) ProcessImage(w http.ResponseWriter, r *http.Request) {
  ctx := r.Context()

  // 获取并验证参数
  imageURL := r.URL.Query().Get("url")
  width := r.URL.Query().Get("w")
  height := r.URL.Query().Get("h")
  format := r.URL.Query().Get("format")

  if imageURL == "" {
    response.Error(w, http.StatusBadRequest, "missing image url")
    return
  }

  // 参数转换
  widthInt, _ := strconv.Atoi(width)
  heightInt, _ := strconv.Atoi(height)
  if widthInt <= 0 || heightInt <= 0 {
    response.Error(w, http.StatusBadRequest, "invalid dimensions")
    return
  }

  // 调用转码服务
  processedImage, err := h.transcodingService.ProcessImage(ctx, imageURL, width, height, format)
  if err != nil {
    logger.Errorf("Failed to process image: %v", err)
    response.Error(w, http.StatusInternalServerError, "image processing failed")
    return
  }

  // 设置响应头
  w.Header().Set("Content-Type", "image/"+format)
  w.Write(processedImage)
}
