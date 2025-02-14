/**
 * ImagicMaker SDK Integration
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package tools

import (
  "fmt"

  "github.com/frog-engine/frog-go/pkg/logger"
  frogsdk "github.com/frog-engine/frog-sdk"
  "github.com/gofiber/fiber/v3"
)

type ImageTools struct{}

func NewImageTools() *ImageTools {
  return &ImageTools{}
}

// Process 处理图片格式转换和尺寸调整
// 1. 调用ImagicMaker SDK
// 2. 进行格式转换
// 3. 进行尺寸调整
// 4. 返回处理后的图片数据
func (t *ImageTools) Process(c fiber.Ctx, imageData []byte, width int, height int, format string) ([]byte, error) {
  // 初始化 sdk
  frogApi := frogsdk.GetAPI()
  logger.Println("frogApi", frogApi)

  // 从 Fiber context 获取原生 context
  // ctx := c.Context()

  // 读取图片数据
  // if err := frogApi.ReadImageBlob(imageData); err != nil {
  //   return nil, fmt.Errorf("failed to read image data: %w", err)
  // }

  // // 调整图片大小
  // if width > 0 && height > 0 {
  //   if err := frogApi.ResizeImage(uint(width), uint(height), 1); err != nil {
  //     return nil, fmt.Errorf("failed to resize image: %w", err)
  //   }
  // }

  // // 根据目标格式进行转换
  // switch strings.ToLower(format) {
  // case "jpg", "jpeg":
  //   // 转换为 JPEG 格式
  //   if err := frogApi.SetImageFormat("jpg"); err != nil {
  //     return nil, fmt.Errorf("failed to convert image to jpg: %w", err)
  //   }
  // case "png":
  //   // 转换为 PNG 格式
  //   if err := frogApi.SetImageFormat("png"); err != nil {
  //     return nil, fmt.Errorf("failed to convert image to png: %w", err)
  //   }
  // default:
  //   return nil, fmt.Errorf("unsupported image format: %s", format)
  // }

  // 获取处理后的图片数据
  processedImage, err := frogApi.ReadImageBlob(imageData)
  if err != nil {
    return nil, fmt.Errorf("failed to get image blob: %w", err)
  }

  return processedImage, nil
}
