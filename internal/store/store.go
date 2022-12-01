package store

type Store interface {
	Articles() IArticlesRepository
}
