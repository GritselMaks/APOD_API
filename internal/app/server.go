package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/GritselMaks/BT_API/internal/apod"
	"github.com/GritselMaks/BT_API/internal/store/models"
	store "github.com/GritselMaks/BT_API/internal/store/postgreSQL"
	"github.com/GritselMaks/BT_API/internal/store/pudgestore"
	"github.com/GritselMaks/BT_API/internal/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const kMaxLogfileSize = 10 * 1024 * 1024

type Server struct {
	config     Config
	router     *mux.Router
	logger     *logrus.Logger
	store      *store.Store
	pudgeStore *pudgestore.Pudge

	apodClient *apod.APODClient
	apodChan   chan apod.ApodOutput
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
	s.apodChan = make(chan apod.ApodOutput)
}

func (s *Server) configRouter() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/articles", s.GetArticles()).Methods("GET")
	router.Handle("/article/{date}", s.GetArticleWithDate()).Methods("GET")
	router.Handle("/picture/{date}", s.GetPicture()).Methods("GET")
	return router
}

func (s *Server) configLoger() {
	logger := logrus.New()
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	logPath, err := utils.StableFilePath(s.config.LogPath)
	if err != nil {
		logger.Error(err)
	}
	loggerFile, err := utils.InitFile(*logPath, kMaxLogfileSize)
	if err != nil {
		logger.Error(err)
	} else {
		logger.SetOutput(loggerFile)
	}
	s.logger = logger
}

func (s *Server) configStore(conf *store.DBConfig) error {
	store, err := store.OpenStore(s.config.Store)
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
	s.apodClient = apod.NewApod(*pdg)
	return nil
}

// ServeHTTP starts the server and blocks until the provided context is closed.
func (s *Server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	s.logger.Info("server starting.....")
	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		s.logger.Info("server.Serve: context closed")
		s.apodClient.Stop()
		time.Sleep(500 * time.Millisecond)
		s.logger.Info("server.Serve: gracefully shutting down")
		shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		errCh <- srv.Shutdown(shutDownCtx)
	}()

	// Run APOD
	go s.apodClient.Run(s.apodChan, s.logger)
	go s.AddNewArticle()

	// Run the server. This will block until the provided context is closed.
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server.Server error: %s", err.Error())
	}

	// Return any errors that happened during shutdown.
	if err := <-errCh; err != nil {
		s.logger.Info("failed to shutdown: %s", err.Error())
		return fmt.Errorf("failed to shutdown: %s", err.Error())
	}

	return nil
}

// Creates an HTTP server using the provided handler,
func (s *Server) ServeHTTPHandler(ctx context.Context) error {
	addr := fmt.Sprintf(s.config.Http.Host + ":" + s.config.Http.Port)
	return s.ServeHTTP(ctx, &http.Server{
		Addr:    addr,
		Handler: s.router,
	})
}

func (s *Server) AddNewArticle() {
	for {
		newPicture := <-s.apodChan
		err := s.AddArticle(newPicture)
		if err != nil {
			s.logger.Errorf("Error.Server: store to database:%s", err.Error())
		}
	}
}

func (s *Server) SavePicture(ulr string) (*string, error) {
	_, err := s.apodClient.GetPicture(ulr)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Server) AddArticle(a apod.ApodOutput) error {
	article := models.MakeArticle(a)
	id, err := s.store.Articles().Create(article)
	if err != nil {
		return err
	}
	fmt.Println(id)
	return nil
}
