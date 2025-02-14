/**
 * Image Model Definition
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package models

type Image struct {
  URL     string `json:"url"`
  Width   int    `json:"width"`
  Height  int    `json:"height"`
  Format  string `json:"format"`
  Size    int64  `json:"size"`
  Quality int    `json:"quality"`
}

// ImageRequest 请求图片参数
type ImageRequest struct {
  URL     string  `json:"url" query:"url"`         // 远程图片 URL
  Quality int32   `json:"quality" query:"quality"` // 质量（1-100）
  Crop    string  `json:"crop" query:"crop"`       // 裁剪格式 x,y,w,h
  Scale   float32 `json:"scale" query:"scale"`     // 缩放比例
  Rotate  string  `json:"rotate" query:"rotate"`   // 旋转角度
  Overlay string  `json:"overlay" query:"overlay"` // 合成图片路径
  Format  string  `json:"format" query:"format"`   // 目标格式（jpg/png/webp）
}

// ImageProcessing 图片处理配置
type ImageProcessing struct {
  Width   int    `json:"width"`
  Height  int    `json:"height"`
  Format  string `json:"format"`
  Quality int    `json:"quality"`
}
