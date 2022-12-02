package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	app "github.com/GritselMaks/BT_API/internal/app"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Parse config
	loadConfig := func() *app.Config {
		cfg, err := app.LoadConfig()
		if err != nil {
			log.Fatal("App::Initialize load config error: ", err)
		}
		return cfg
	}
	conf := loadConfig()
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", conf.Store.User, conf.Store.Password, conf.Store.Host, conf.Store.DBName)

	// Update migrations
	m, err := migrate.New(
		"file://../../internal/store/migrations",
		databaseUrl)
	if err != nil {
		log.Fatal("App::load migrate error: ", err)
	}
	m.Up()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer func() {
		stop()
		if r := recover(); r != nil {
			log.Fatal("application panic", "panic", r)
		}
	}()

	// Create server and start
	s := app.NewServer(*conf)
	s.Initialize()

	// If you want to add pictures from the previous month, you should uncommitted.
	// New picture will add every day.

	// s.AddContent()

	err = s.ServeHTTPHandler(ctx)
	defer stop()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("successful shutdown")
}
