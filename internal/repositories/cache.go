package repositories

type Cache interface {
  Get(key string) ([]byte, bool)
  Set(key string, value []byte)
}
