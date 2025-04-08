package storage

type StorageEngineInterface interface {
	Put(key, value string)
	Get(key string) (string, bool)
	Del(key string)
	Prefix(prefix string) []string
	Suffix(suffix string) []string
	Contains(sub string) []string
	// Import(keys,values []string)
	// Export()
}

const (
	GOMAP = "gomap"
	HASHMAP = "hashmap"
)

func NewStorageEngine(name string) StorageEngineInterface {
	if name == GOMAP {
		return NewGomap(16)
	} else if name == HASHMAP {
		hashFunc := func(key interface{}) int {
			return key.(int) // 简单示例，直接使用int值作为哈希
		}
		compareFunc := func(a, b interface{}) int {
			return a.(int) - b.(int)
		}
		hm := NewHashMap(16, hashFunc, compareFunc)
		return hm
	}
	return nil
}