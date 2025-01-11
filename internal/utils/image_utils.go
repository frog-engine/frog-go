/**
 * Image Processing Utilities
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package utils

import (
  "context"
  "io"
  "net/http"
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
