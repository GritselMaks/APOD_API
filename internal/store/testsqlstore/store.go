package testsqlstore

import (
	"github.com/GritselMaks/BT_API/internal/store"
	"github.com/GritselMaks/BT_API/internal/store/models"
)

type Store struct {
	db map[string]models.Articles
}

func TestStore() *Store {
	return &Store{db: make(map[string]models.Articles)}
}

func (s *Store) Articles() store.IArticlesRepository {
	articlesRepository := &ArticlesRepository{store: s}
	return articlesRepository
}
