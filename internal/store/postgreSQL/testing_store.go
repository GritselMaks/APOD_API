package postgresql

import (
	"fmt"
	"strings"
	"testing"
)

const (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "postgres"
	dbPassword = "Root"
	dbName     = "bt_api_test"
)

// TestStore
func TestStore(t *testing.T) (*Store, func(...string)) {
	t.Helper()
	config := DBConfig{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
	}
	s, err := OpenStore(&config)
	if err != nil {
		t.Fatal(err)
	}
	return s, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := s.db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}
		s.CloseStore()
	}
}
