/**
 * Image Handler
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package handlers

import (
  "github.com/frog-engine/frog-go/internal/models"
  "github.com/frog-engine/frog-go/internal/services"
  "github.com/frog-engine/frog-go/pkg/logger"
  "github.com/frog-engine/frog-go/pkg/response"

  "github.com/gofiber/fiber/v3"
)

type ImageHandler struct {
  imageService *services.ImageService
}

func NewImageHandler(ts *services.ImageService) *ImageHandler {
  return &ImageHandler{
    imageService: ts,
  }
}

// ProcessImage 处理图片转码请求
func (h *ImageHandler) ProcessImage(c fiber.Ctx) error {
  // 获取并验证参数
  var imageRequest models.ImageRequest
  if err := c.Bind().Query(&imageRequest); err != nil {
    return response.Error(c, fiber.StatusBadRequest, "invalid request")
  }
  logger.Printf("ImageRequest: %+v\n", imageRequest)
  if imageRequest.URL == "" {
    response.Error(c, fiber.StatusBadRequest, "url is required")
  }

  if imageRequest.URL == "" {
    return response.Error(c, fiber.StatusBadRequest, "missing image url")
  }

  // 调用图片转码服务
  processedImage, err := h.imageService.ProcessImage(c, imageRequest)
  if err != nil {
    logger.Errorf("Failed to process image: %v", err)
    return response.Error(c, fiber.StatusInternalServerError, "image processing failed")
  }

  // 设置响应头
  c.Set("Content-Type", "image/"+imageRequest.Format)
  return c.Send(processedImage)
  // return response.Success(c, imageRequest)
}

// 处理图片 API
// func processImage1(c fiber.Ctx) error {
//   // 获取参数
//   imageURL := c.FormValue("image_url") // 远程 URL
//   format := c.FormValue("format")      // 目标格式（jpg/png/webp）
//   quality := c.FormValue("quality")    // 质量（1-100）
//   crop := c.FormValue("crop")          // 裁剪格式 x,y,w,h
//   scale := c.FormValue("scale")        // 缩放 w,h
//   rotate := c.FormValue("rotate")      // 旋转角度
//   overlay := c.FormValue("overlay")    // 合成图片路径

//   var imgData []byte
//   var err error

//   if imageURL != "" {
//     // 下载远程图片
//     imgData, err = downloadImage(imageURL)
//     if err != nil {
//       return c.Status(400).SendString("下载图片失败: " + err.Error())
//     }
//   }

//   // 保存到临时文件
//   inputFile := "input.png"
//   outputFile := "output." + format
//   err = os.WriteFile(inputFile, imgData, 0644)
//   if err != nil {
//     return c.Status(500).SendString("无法保存输入图片")
//   }

//   // 处理图片
//   err = transformImage(inputFile, outputFile, format, quality, crop, scale, rotate, overlay)
//   if err != nil {
//     return c.Status(500).SendString("图片处理失败: " + err.Error())
//   }

//   // 返回处理后的图片
//   return c.SendFile(outputFile)
// }
