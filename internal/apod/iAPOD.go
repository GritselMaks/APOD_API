package apod

import "github.com/GritselMaks/BT_API/internal/store/models"

type IAPOD interface {
	GetContent(date, start, end string) ([]models.Article, error)
	SavePicture(url, key string) error
}
