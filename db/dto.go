package db

type Storable interface {
	GetPrimaryKey() int
}
type DB[T Storable] interface {
	Insert(obj T)
	Delete(key int) T
	Get(key int) (T, error)
}
