/**
 * Image Transcoding Service
 * @description 负责转码调用
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package services

import (
  "context"
  "fmt"
  "time"

  "github.com/frog-engine/frog-go/internal/models"
  "github.com/frog-engine/frog-go/internal/repositories"
  "github.com/frog-engine/frog-go/internal/tools"
  "github.com/frog-engine/frog-go/internal/utils"
  "github.com/frog-engine/frog-go/pkg/logger"

  "github.com/gofiber/fiber/v3"
)

type TranscodingService struct {
  cache      repositories.Cache
  imageTools tools.ImageTools
}

func NewTranscodingService(cache repositories.Cache, imageTools *tools.ImageTools) *TranscodingService {
  return &TranscodingService{
    cache:      cache,
    imageTools: *imageTools,
  }
}

// ProcessImage 处理图片转码
func (s *TranscodingService) ProcessImage(c fiber.Ctx, imageRequest models.ImageRequest) ([]byte, error) {
  // 生成缓存key
  cacheKey := fmt.Sprintf("%s_%s_%s", imageRequest.URL, imageRequest.Crop, imageRequest.Format)

  // 从 Fiber context 获取原生 context
  ctx := c.Context()

  // 设置处理超时
  ctx, cancel := context.WithTimeout(ctx, 450*time.Millisecond)
  defer cancel()

  // 先查询缓存
  if cachedImage, found := s.cache.Get(cacheKey); found {
    logger.Debug("Cache hit for image: " + imageRequest.URL)
    return cachedImage, nil
  }

  // 下载原图
  originalImageData, err := utils.DownloadImage(ctx, imageRequest.URL)
  if err != nil {
    return nil, fmt.Errorf("download failed: %w", err)
  }

  // originalImage, _, err := image.Decode(bytes.NewReader(originalImageData))
  // if err != nil {
  //   return nil, fmt.Errorf("download failed: %w", err)
  // }

  // // 图片处理
  // imageBytes, err := utils.ImageToBytes(originalImage, format)
  // if err != nil {
  //   return nil, fmt.Errorf("image conversion failed: %w", err)
  // }
  // processedImage, err := s.imageTools.Process(c, imageBytes, width, height, format)

  processedImage, err := s.imageTools.Process(c, originalImageData, imageRequest)
  if err != nil {
    return nil, fmt.Errorf("processing failed: %w", err)
  }

  // 存入缓存
  s.cache.Set(cacheKey, processedImage)

  return processedImage, nil
}
