package ports

type Store interface {
	Put(key, value string) error
	Append(key, arg string) (string, error)
	Get(key string) (string, error)
}
