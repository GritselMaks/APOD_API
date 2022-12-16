package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/GritselMaks/BT_API/internal/apod"
	"github.com/GritselMaks/BT_API/internal/event_loop"
	"github.com/GritselMaks/BT_API/internal/store"
	"github.com/GritselMaks/BT_API/internal/store/models"
	"github.com/GritselMaks/BT_API/internal/store/postgresql"
	"github.com/GritselMaks/BT_API/internal/store/pudgestore"
	"github.com/GritselMaks/BT_API/internal/utils"
	"github.com/GritselMaks/BT_API/pkg/logger"
	"github.com/gorilla/mux"
)

const (
	period time.Duration = time.Hour * 24
	offset time.Duration = time.Hour * 12
)

type Server struct {
	config     Config
	router     *mux.Router
	logger     *logger.Logger
	store      store.Store
	pudgeStore store.BinarStorage

	apodClient apod.IAPOD
	loop       event_loop.EventLoop
}

func NewServer(conf Config) *Server {
	return &Server{config: conf}
}

func (s *Server) Initialize() {
	s.configLoger()
	s.router = s.configRouter()
	if err := s.configStore(s.config.Store); err != nil {
		s.logger.Fatalf("error initialize storages: %v", err.Error())
	}
	s.loop = *event_loop.NewEvenlLoop()
}

func (s *Server) configRouter() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/articles", s.GetArticles()).Methods("GET")
	router.Handle("/article/{date}", s.GetArticleWithDate()).Methods("GET")
	router.Handle("/picture/{date}", s.GetPicture()).Methods("GET")
	return router
}

func (s *Server) configLoger() {
	s.logger = logger.NewLogger(s.config.LogLevel)
	s.logger.Info("logger is created")
}

func (s *Server) configStore(conf *postgresql.DBConfig) error {
	store, err := postgresql.OpenStore(s.config.Store)
	if err != nil {
		return err
	}

	path, err := utils.StableFilePath(s.config.LocalStore)
	if err != nil {
		return err
	}

	pdg, err := pudgestore.Open(*path)
	if err != nil {
		return err
	}

	s.store = store
	s.pudgeStore = pdg
	s.apodClient = apod.NewAPOD("", pdg)
	return nil
}

// ServeHTTP starts the server and blocks until the provided context is closed.
// Run shaduler for getting content
func (s *Server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	s.logger.Info("server starting.....")
	defer s.logger.Info("server stoped.....")
	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		s.logger.Info("server.Serve: context closed")
		time.Sleep(500 * time.Millisecond)
		s.logger.Info("server.Serve: gracefully shutting down")
		shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		errCh <- srv.Shutdown(shutDownCtx)
	}()

	// Run the server. This will block until the provided context is closed.
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server.Server error: %s", err.Error())
	}

	// Return any errors that happened during shutdown.
	if err := <-errCh; err != nil {
		s.logger.Errorf("failed to shutdown: %s", err.Error())
		return fmt.Errorf("failed to shutdown: %s", err.Error())
	}

	return nil
}

// Creates an HTTP server using the provided handler,
func (s *Server) ServerRun(ctx context.Context) error {
	// If you want to add pictures from the previous month, you should uncommitted.
	// New picture will add every day.
	// s.AddContent()

	s.loop.AddEvent(s.CustomLoopFunc())
	go s.loop.Schedule(ctx, period, offset)

	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	err := s.ServeHTTP(ctx, &http.Server{
		Addr:    addr,
		Handler: s.router,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) AddArticle(article *models.Article) error {
	err := s.store.Articles().Create(article)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) CustomLoopFunc() func() {
	return func() {
		articles, err := s.apodClient.GetContent("", "", "")
		if err != nil {
			s.logger.Errorf("Error get article: %s", err.Error())
			return
		}
		for _, article := range articles {
			if err := s.apodClient.SavePicture(article.Url, article.Date); err != nil {
				s.logger.Errorf("Error save picture: %s", err.Error())
				continue
			}
			if err := s.store.Articles().Create(&article); err != nil {
				s.logger.Errorf("Error save article: %s", err.Error())
				continue
			}
		}
	}
}

// Func run for adddig date in storage since 02.11.2022
func (s *Server) AddContent() {
	rows, err := s.apodClient.GetContent("", "2022-11-02", time.Now().Format("2006-01-02"))
	if err != nil {
		return
	}
	for _, r := range rows {
		if err := s.apodClient.SavePicture(r.Url, r.Date); err != nil {
			s.logger.Debugf("Error save picture: %s", err.Error())
			continue
		}
		if err := s.store.Articles().Create(&r); err != nil {
			s.logger.Debugf("Error save article: %s", err.Error())
			continue
		}
	}
}
