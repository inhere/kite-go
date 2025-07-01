package lcstorage

type SimpleStorage interface {
	Get(key string) (any, error)
	Set(key string, value any) error
	Delete(key string) error
	Clear() error
	Store() error
	Restore() error

	Keys() []string
	Len() int
	IsEmpty() bool
	Exists(key string) bool
}

type LocalStorage struct {
}

type XFileStorage struct {
}

// MixinFileStorage 小key统一存储到一个默认文件；大于 512 的key，存储到单独的文件中
type MixinFileStorage struct {
}

type MultiFileStorage struct {
}
