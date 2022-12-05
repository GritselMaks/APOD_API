package models

import "github.com/GritselMaks/BT_API/internal/apod"

type Article struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Url         string `json:"url"`
}

func MakeArticle(a apod.ApodOutput) *Article {
	return &Article{
		Title:       a.Title,
		Date:        a.Date,
		Explanation: a.Explanation,
	}
}
