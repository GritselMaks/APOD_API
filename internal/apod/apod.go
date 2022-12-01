package apod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/GritselMaks/BT_API/internal/store/pudgestore"
	"github.com/sirupsen/logrus"
)

const (
	duration   int64  = 12
	url        string = "https://api.nasa.gov/planetary/apod"
	defaultKey string = "DEMO_KEY"
)

type APODClient struct {
	ApiKey    string
	Url       string
	Durration int64
	Quit      chan bool

	pudgeStore pudgestore.Pudge
}

func NewApod(pudge pudgestore.Pudge) *APODClient {
	return &APODClient{
		ApiKey:     defaultKey,
		Durration:  duration,
		Quit:       make(chan bool),
		Url:        url,
		pudgeStore: pudge,
	}
}

type ApodOutput struct {
	Title          string `json:"title,omitempty"`
	Explanation    string `json:"explanation,omitempty"`
	Date           string `json:"date,omitempty"`
	MediaType      string `json:"media_type,omitempty"`
	Url            string `json:"url,omitempty"`
	HdUrl          string `json:"hdurl,omitempty"`
	ThumbnailUrl   string `json:"thumbnail_url,omitempty"`
	Copyright      string `json:"Copyright,omitempty"`
	ServiceVersion string `json:"service_version,omitempty"`
}

func (a *APODClient) Query() (*ApodOutput, error) {
	resp, err := http.Get(fmt.Sprintf("%s?api_key=%s", a.Url, a.ApiKey))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.Header.Get("X-RateLimit-Remaining") == "0" {
		return nil, fmt.Errorf("you have exceeded your rate limit")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result ApodOutput
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *APODClient) GetPicture(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *APODClient) Run(ch chan ApodOutput, logger *logrus.Logger) {
	for {
		select {
		case <-a.Quit:
			return
		case <-time.After(time.Duration(a.Durration) * time.Hour):
			apod, err := a.Query()
			if err != nil {
				logger.Errorf("APOND.Error:%s", err.Error())
				break
			}
			data, err := a.GetPicture(apod.Url)
			if err != nil {
				logger.Errorf("APOND.Error:%s", err.Error())
				break
			}
			err = a.SavePicture(apod.Date, data)
			if err != nil {
				logger.Errorf("APOND.Error:%s", err.Error())
				break
			}
			ch <- *apod
		}
	}
}

func (a *APODClient) Stop() {
	a.Quit <- true
}

func (a *APODClient) SavePicture(key string, value []byte) error {
	return a.pudgeStore.Set(key, value)
}
