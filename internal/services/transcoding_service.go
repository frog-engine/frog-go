/**
 * Image Transcoding Service
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package services

import (
	"context"
	"fmt"
	"time"

	"frog-go/internal/cache"
	"frog-go/internal/tools"
	"frog-go/internal/utils"
	"frog-go/pkg/logger"
)

type TranscodingService struct {
  cache      *cache.ImageCache
  imageTools *tools.ImageTools
}

func NewTranscodingService(cache *cache.ImageCache, tools *tools.ImageTools) *TranscodingService {
  return &TranscodingService{
    cache:      cache,
    imageTools: tools,
  }
}

// ProcessImage 处理图片转码
func (s *TranscodingService) ProcessImage(ctx context.Context, imageURL, width, height, format string) ([]byte, error) {
  // 生成缓存key
  cacheKey := fmt.Sprintf("%s_%s_%s_%s", imageURL, width, height, format)

  // 设置处理超时
  ctx, cancel := context.WithTimeout(ctx, 450*time.Millisecond)
  defer cancel()

  // 先查询缓存
  if cachedImage, found := s.cache.Get(cacheKey); found {
    logger.Debug("Cache hit for image: " + imageURL)
    return cachedImage, nil
  }

  // 下载原图
  originalImage, err := utils.DownloadImage(ctx, imageURL)
  if err != nil {
    return nil, fmt.Errorf("download failed: %w", err)
  }

  // 图片处理
  processedImage, err := s.imageTools.Process(ctx, originalImage, width, height, format)
  if err != nil {
    return nil, fmt.Errorf("processing failed: %w", err)
  }

  // 存入缓存
  s.cache.Set(cacheKey, processedImage)

  return processedImage, nil
}
