package testsqlstore

import (
	"errors"
	"fmt"

	"github.com/GritselMaks/BT_API/internal/store/models"
)

// Test ArticleRepository
type ArticlesRepository struct {
	store *Store
}

func (r *ArticlesRepository) Create(a *models.Article) error {
	if a != nil {
		r.store.db[a.Date] = *a
	}
	return nil
}

func (r *ArticlesRepository) ShowArticlebByDate(date string) (*models.Article, error) {
	article, ok := r.store.db[date]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &article, nil
}

func (r *ArticlesRepository) ShowArticles() ([]models.Article, error) {
	var artisles []models.Article
	for k, v := range r.store.db {
		if k != v.Date {
			return nil, fmt.Errorf("Error")
		}
		artisles = append(artisles, v)
	}
	return artisles, nil
}
