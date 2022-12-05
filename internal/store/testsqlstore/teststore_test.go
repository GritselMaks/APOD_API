package testsqlstore_test

import (
	"testing"

	"github.com/GritselMaks/BT_API/internal/store/models"
	"github.com/GritselMaks/BT_API/internal/store/testsqlstore"
	"github.com/stretchr/testify/assert"
)

func TestArticleRepository_Create(t *testing.T) {
	s := testsqlstore.TestStore()
	err := s.Articles().Create(&models.Article{
		Date: "2022-11-29",
	})
	assert.NoError(t, err)
}
