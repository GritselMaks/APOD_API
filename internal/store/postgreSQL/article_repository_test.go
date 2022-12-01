package postgresql_test

import (
	"testing"

	"github.com/GritselMaks/BT_API/internal/store/models"
	store "github.com/GritselMaks/BT_API/internal/store/postgreSQL"
	"github.com/stretchr/testify/assert"
)

func TestArticleRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t)
	defer teardown("article")
	article, err := s.Articles().Create(&models.Articles{
		Date: "2022-11-29",
	})
	assert.NoError(t, err)
	assert.NotNil(t, article)
}

func TestArticleRepositoryFindByDate(t *testing.T) {
	s, teardown := store.TestStore(t)
	defer teardown("article")
	date := "2022-11-29"

	_, err := s.Articles().ShowArticlebByDate(date)
	assert.Error(t, err)

	s.Articles().Create(&models.Articles{Date: date, Title: "Title", Explanation: "Explanation"})
	article, err := s.Articles().ShowArticlebByDate(date)
	assert.NoError(t, err)
	assert.Equal(t, date, article.Date)
}

func TestArticleRepositoryShowArticles(t *testing.T) {
	s, teardown := store.TestStore(t)
	defer teardown("article")
	a1 := models.Articles{Date: "2022-11-29", Title: "Title", Explanation: "Explanation"}
	a2 := models.Articles{Date: "2022-11-30", Title: "Title", Explanation: "Explanation"}

	res, err := s.Articles().ShowArticles()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res))
	s.Articles().Create(&a1)
	s.Articles().Create(&a2)

	res, err = s.Articles().ShowArticles()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}
