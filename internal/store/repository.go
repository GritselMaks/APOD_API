package store

import "github.com/GritselMaks/BT_API/internal/store/models"

type IArticlesRepository interface {
	Create(a *models.Article) error
	ShowArticlebByDate(date string) (*models.Article, error)
	ShowArticles() ([]models.Article, error)
}
