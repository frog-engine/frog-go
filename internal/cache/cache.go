/**
 * Image Cache Implementation
 *
 * @author jarryli@gmail.com
 * @date 2024-12-20
 */

package cache

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type ImageCache struct {
  cache *cache.Cache
  mutex sync.RWMutex
}

func NewImageCache() *ImageCache {
  // 创建一个默认过期时间为1小时，每10分钟清理一次的缓存
  c := cache.New(1*time.Hour, 10*time.Minute)
  return &ImageCache{
    cache: c,
  }
}

func (ic *ImageCache) Get(key string) ([]byte, bool) {
  ic.mutex.RLock()
  defer ic.mutex.RUnlock()

  if data, found := ic.cache.Get(key); found {
    return data.([]byte), true
  }
  return nil, false
}

func (ic *ImageCache) Set(key string, data []byte) {
  ic.mutex.Lock()
  defer ic.mutex.Unlock()

  ic.cache.Set(key, data, cache.DefaultExpiration)
}
