package pudgestore

type Repository interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	GetList() ([][]byte, error)
}
