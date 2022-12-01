package testsqlstore

import (
	"fmt"

	"github.com/GritselMaks/BT_API/internal/store/models"
)

// Test ArticleRepository
type ArticlesRepository struct {
	store *Store
}

func (r *ArticlesRepository) Create(a *models.Articles) error {
	if a != nil {
		r.store.db[a.Date] = *a
	}
	return nil
}

func (r *ArticlesRepository) ShowArticlebByDate(date string) (*models.Articles, error) {
	article, ok := r.store.db[date]
	if !ok {
		return nil, fmt.Errorf("User not found")
	}
	return &article, nil
}

func (r *ArticlesRepository) ShowArticles() ([]models.Articles, error) {
	var artisles []models.Articles
	for k, v := range r.store.db {
		if k != v.Date {
			return nil, fmt.Errorf("Error")
		}
		artisles = append(artisles, v)
	}
	return artisles, nil
}
