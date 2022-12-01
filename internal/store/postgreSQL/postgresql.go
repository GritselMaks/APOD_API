package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/GritselMaks/BT_API/internal/store"
	"github.com/GritselMaks/BT_API/internal/utils"
	_ "github.com/lib/pq"
)

const maxAttempts = 5

type Store struct {
	db *sql.DB
}

func OpenStore(config *DBConfig) (*Store, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	var store Store
	err := utils.DoWithTries(func() error {
		con, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			return err
		}
		if err := con.Ping(); err != nil {
			return err
		}
		store.db = con
		return nil
	}, maxAttempts, 3*time.Second)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *Store) CloseStore() {
	s.db.Close()
}

func (s *Store) Articles() store.IArticlesRepository {
	articlesRepository := &ArticlesRepository{store: s}
	return articlesRepository
}
