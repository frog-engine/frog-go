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

type ProcessingOptions struct {
  Width   int    `json:"width"`
  Height  int    `json:"height"`
  Format  string `json:"format"`
  Quality int    `json:"quality"`
}
