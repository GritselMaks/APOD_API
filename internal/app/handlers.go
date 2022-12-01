package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) GetArticles() http.Handler {
	return http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugf("starting")
		defer s.logger.Debugf("finishing")
		articles, err := s.store.Articles().ShowArticles()
		if err != nil {
			s.logger.Debugf("failed getting data from store: %s", err.Error())
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		for i, article := range articles {
			url := fmt.Sprintf("http://%s:%s/picture/%s", s.config.Http.Host, s.config.Http.Port, article.Date)
			articles[i].Url = url
		}
		ResponseWithJSON(w, http.StatusOK, articles)
	}))
}

func (s *Server) GetArticleWithDate() http.Handler {
	return http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugf("starting")
		defer s.logger.Debugf("finishing")
		date, ok := mux.Vars(r)["date"]
		if !ok {
			s.logger.Debugf("failed getting request data")
			RespondWithError(w, http.StatusBadRequest, "Bad request data")
			return
		}
		article, err := s.store.Articles().ShowArticlebByDate(date)
		if err != nil {
			s.logger.Debugf("failed getting data from store: %s", err.Error())
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		url := fmt.Sprintf("http://%s:%s/picture/%s", s.config.Http.Host, s.config.Http.Port, article.Date)
		article.Url = url
		ResponseWithJSON(w, http.StatusOK, article)
	}))
}

func (s *Server) GetPicture() http.Handler {
	return http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugf("starting")
		defer s.logger.Debugf("finishing")
		date, ok := mux.Vars(r)["date"]
		if !ok {
			s.logger.Debugf("failed getting request data")
			RespondWithError(w, http.StatusBadRequest, "Bad request data")
			return
		}
		picture, err := s.pudgeStore.Get(date)
		if err != nil {
			s.logger.Debugf("failed getting data from store: %s", err.Error())
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithPicture(w, http.StatusOK, picture)
	}))
}
