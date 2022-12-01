package models

import "github.com/GritselMaks/BT_API/internal/apod"

type Articles struct {
	ID          int    `json:"-"`
	Date        string `json:"data"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Url         string `json:"url"`
}

func MakeArticle(a apod.ApodOutput) *Articles {
	return &Articles{
		Title:       a.Title,
		Date:        a.Date,
		Explanation: a.Explanation,
	}
}
