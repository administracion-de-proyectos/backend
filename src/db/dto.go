package db

type Storable interface {
	GetPrimaryKey() string
}
type DB[T Storable] interface {
	Insert(obj T)
	Update(obj T)
	Delete(key string) T
	Get(key string) (T, error)
	GetAll() ([]T, error)
}
