package chfs

type Value[K comparable] interface {
	Key() K
}

/**
 * A key-value store that allows for bulk inserts/gets
 **/
type Store[K comparable, V Value[K]] interface {
	Get(key K) (V, error)
	Gets(keys []K) []V
	Put(value V) (K, error)
	Puts(values []V) []K
}
