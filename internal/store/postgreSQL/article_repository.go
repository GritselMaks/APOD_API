package postgresql

import (
	"github.com/GritselMaks/BT_API/internal/store"
	"github.com/GritselMaks/BT_API/internal/store/models"
)

type ArticlesRepository struct {
	store *Store
}

func (r *ArticlesRepository) Create(a *models.Article) error {
	res := r.store.db.QueryRow(
		"INSERT INTO article (date_at, title, explanation) values ($1,$2,$3)",
		a.Date,
		a.Title,
		a.Explanation,
	)
	if res.Err() != nil {
		return store.ErrNotFound
	}
	return nil
}

func (r *ArticlesRepository) ShowArticlebByDate(date string) (*models.Article, error) {
	article := models.Article{}
	if err := r.store.db.QueryRow("SELECT * FROM article WHERE date_at=$1", date).Scan(&article.Date, &article.Title, &article.Explanation); err != nil {
		return nil, store.ErrNotFound
	}
	return &article, nil
}

func (r *ArticlesRepository) ShowArticles() ([]models.Article, error) {
	articles := []models.Article{}
	rows, err := r.store.db.Query("SELECT * FROM article")
	if err != nil {
		return nil, store.ErrNotContent
	}
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.Date, &article.Title, &article.Explanation); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}
