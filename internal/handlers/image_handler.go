/**
 * Image Handler
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package handlers

import (
  "strconv"

  "github.com/frog-engine/frog-go/internal/services"
  "github.com/frog-engine/frog-go/pkg/logger"
  "github.com/frog-engine/frog-go/pkg/response"

  "github.com/gofiber/fiber/v3"
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
func (h *ImageHandler) ProcessImage(c fiber.Ctx) error {
  // 获取并验证参数
  imageURL := c.Query("url")
  width := c.Query("w")
  height := c.Query("h")
  format := c.Query("format")

  if imageURL == "" {
    return response.Error(c, fiber.StatusBadRequest, "missing image url")
  }

  // 参数转换
  widthInt, _ := strconv.Atoi(width)
  heightInt, _ := strconv.Atoi(height)
  if widthInt <= 0 || heightInt <= 0 {
    return response.Error(c, fiber.StatusBadRequest, "invalid dimensions")
  }

  // 调用转码服务
  processedImage, err := h.transcodingService.ProcessImage(c, imageURL, widthInt, heightInt, format)
  if err != nil {
    logger.Errorf("Failed to process image: %v", err)
    return response.Error(c, fiber.StatusInternalServerError, "image processing failed")
  }

  // 设置响应头
  c.Set("Content-Type", "image/"+format)
  return c.Send(processedImage)
}
