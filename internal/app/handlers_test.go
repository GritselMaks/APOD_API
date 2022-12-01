package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GritselMaks/BT_API/internal/store/models"
	"github.com/GritselMaks/BT_API/internal/store/testbinarstore"
	"github.com/GritselMaks/BT_API/internal/store/testsqlstore"
	"github.com/sirupsen/logrus"
)

func testServer() *Server {
	cfg := Config{}
	s := NewServer(cfg)
	s.store = testsqlstore.TestStore()
	s.pudgeStore = testbinarstore.TestBinarStore()
	s.logger = logrus.New()
	s.router = s.configRouter()
	return s
}

func Test_GetArticles(t *testing.T) {
	t.Parallel()
	srv := testServer()
	a := models.Articles{
		Date: "2022-12-1",
	}
	srv.store.Articles().Create(&a)
	type args struct {
		req *http.Request
	}

	tests := []struct {
		name     string
		args     func(t *testing.T) args
		wantCode int
		wantBody bool
	}{
		{
			name: "/articles http.StatusOK",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/articles", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			resp := httptest.NewRecorder()
			srv.router.ServeHTTP(resp, tArgs.req)

			if resp.Result().StatusCode != tt.wantCode {
				t.Fatalf("Test:%s| the status code should be [%d] but received [%d]", tt.name, resp.Result().StatusCode, tt.wantCode)
			}
		})
	}
}

func Test_GetArticleWithDate(t *testing.T) {
	t.Parallel()
	srv := testServer()
	a := models.Articles{
		Date: "2022-12-1",
	}
	srv.store.Articles().Create(&a)
	type args struct {
		req *http.Request
	}

	tests := []struct {
		name     string
		args     func(t *testing.T) args
		wantCode int
	}{
		{
			name: "/articles/date http.StatusOK ",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/article/2022-12-1", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "/articles/data http.StatusNotFound ",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/article/", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "/articles/data http.StatusBadRequest ",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/article/2022-22-22", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			resp := httptest.NewRecorder()
			srv.router.ServeHTTP(resp, tArgs.req)

			if resp.Result().StatusCode != tt.wantCode {
				t.Fatalf("Test:%s| the status code should be [%d] but received [%d]", tt.name, resp.Result().StatusCode, tt.wantCode)
			}
		})
	}
}

func Test_GetPicture(t *testing.T) {
	t.Parallel()
	srv := testServer()
	a := models.Articles{
		Date: "2022-12-1",
	}
	srv.store.Articles().Create(&a)
	srv.pudgeStore.Set("2022-12-1", []byte(a.Date))

	type args struct {
		req *http.Request
	}

	tests := []struct {
		name     string
		args     func(t *testing.T) args
		wantCode int
	}{
		{
			name: "/picture/date http.StatusOK ",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/picture/2022-12-1", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "/picture/data http.StatusNotFound ",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/picture/", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "/picture/data http.StatusBadRequest ",
			args: func(*testing.T) args {
				req, err := http.NewRequest("GET", "/picture/2022-22-22", nil)
				if err != nil {
					t.Fatalf("fail to create request: %s", err.Error())
				}
				return args{
					req: req,
				}
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			resp := httptest.NewRecorder()
			srv.router.ServeHTTP(resp, tArgs.req)
			if resp.Result().StatusCode != tt.wantCode {
				t.Fatalf("Test:%s| the status code should be [%d] but received [%d]", tt.name, resp.Result().StatusCode, tt.wantCode)
			}
		})
	}
}
