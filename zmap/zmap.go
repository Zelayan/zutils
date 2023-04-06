package zmap

type Zmap interface {
	Load(key string) (any, bool)
	Store(key string, value any)
	Delete(key string)
}
