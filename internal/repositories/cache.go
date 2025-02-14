package repositories

import "image"

type Cache interface {
  Set(key string, value []byte)
  Get(key string) ([]byte, bool)
  SetImage(key string, value image.Image)
  GetImage(key string) (image.Image, bool)
}
