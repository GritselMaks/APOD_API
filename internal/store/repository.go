package store

import "github.com/GritselMaks/BT_API/internal/store/models"

type IArticlesRepository interface {
	Create(a *models.Articles) (*int, error)
	ShowArticlebByDate(date string) (*models.Articles, error)
	ShowArticles() ([]models.Articles, error)
}
