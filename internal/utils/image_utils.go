/**
 * Image Processing Utilities
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package utils

import (
  "bytes"
  "context"
  "errors"
  "image"
  "image/gif"
  "image/jpeg"
  "image/png"
  "io"
  "net/http"
  // "github.com/chai2010/webp"
  // "github.com/Kagami/go-avif"
)

// DownloadImage 从URL下载图片
func DownloadImage(ctx context.Context, url string) ([]byte, error) {
  req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
  if err != nil {
    return nil, err
  }

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  return io.ReadAll(resp.Body)
}

// ImageToBytes converts an image.Image to a byte slice in the specified format.
func ImageToBytes(img image.Image, format string) ([]byte, error) {
  var buf bytes.Buffer

  switch format {
  case "jpeg", "jpg":
    if err := jpeg.Encode(&buf, img, nil); err != nil {
      return nil, err
    }
  case "png":
    if err := png.Encode(&buf, img); err != nil {
      return nil, err
    }
  case "gif":
    if err := gif.Encode(&buf, img, nil); err != nil {
      return nil, err
    }
    /*
       case "webp":
         // arch -arm64 brew install webp
         // sudo yum install libwebp-devel
         if err := webp.Encode(&buf, img, nil); err != nil {
           return nil, err
         }
       case "avif":
         // arch -arm64 brew install aom
         // sudo yum install libaom-devel
         if err := avif.Encode(&buf, img, nil); err != nil {
           return nil, err
         }
    */
  default:
    return nil, errors.New("unsupported format: " + format)
  }

  return buf.Bytes(), nil
}
