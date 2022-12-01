package store

import "github.com/GritselMaks/BT_API/internal/store/models"

type ArticlesRepository struct {
	store *Store
}

func (r *ArticlesRepository) Create(a *models.Articles) (*int, error) {
	if err := r.store.db.QueryRow(
		"INSERT INTO article (date_at, title, explanation) values ($1,$2,$3) RETURNING id",
		a.Date,
		a.Title,
		a.Explanation,
	).Scan(&a.ID); err != nil {
		return nil, err
	}
	return &a.ID, nil
}

func (r *ArticlesRepository) ShowArticlebByDate(date string) (*models.Articles, error) {
	article := models.Articles{}
	if err := r.store.db.QueryRow("SELECT * FROM article WHERE date_at=$1", date).Scan(&article.ID, &article.Date, &article.Title, &article.Explanation); err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticlesRepository) ShowArticles() ([]models.Articles, error) {
	articles := []models.Articles{}
	rows, err := r.store.db.Query("SELECT * FROM article")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var article models.Articles
		if err := rows.Scan(&article.ID, &article.Date, &article.Title, &article.Explanation); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}
