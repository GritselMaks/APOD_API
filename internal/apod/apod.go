package apod

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/GritselMaks/BT_API/internal/store"
	"github.com/GritselMaks/BT_API/internal/store/models"
)

const (
	duration   int64  = 12
	url        string = "https://api.nasa.gov/planetary/apod"
	pictureUrl string = "https://api.nasa.gov/planetary/apod/static/"
	defaultKey string = "DEMO_KEY"
)

var ErrorBadQueryParams = errors.New("bad query params")

// APOD is the struct with url for getting content
type APODClient struct {
	store   store.BinarStorage
	apodUrl string
}

func NewAPOD(key string, store store.BinarStorage) APODClient {
	apiKey := defaultKey
	if len(key) != 0 {
		apiKey = key
	}
	return APODClient{
		apodUrl: fmt.Sprintf("%s?api_key=%s", url, apiKey),
		store:   store,
	}
}

// ApodQueryInput is the input for an Apod Query
type ApodQueryInput struct {
	Date      time.Time `json:"date"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Count     int       `json:"count"`
	Thumbs    bool      `json:"thumbs"`
}

func MakeApodQueryInput(date, start, end string) (*ApodQueryInput, error) {
	// "Start" and "end" of a date cannot be used with "date"
	if len(date) != 0 && (len(start) != 0 || len(end) != 0) {
		return nil, ErrorBadQueryParams
	}
	query := ApodQueryInput{}
	if len(date) != 0 {
		t, error := time.Parse("2006-01-02", date)
		if error != nil {
			return nil, ErrorBadQueryParams
		}
		query.Date = t
	}

	// "Start" and "end" of a date should be used together
	if len(start) != 0 && len(end) != 0 {
		startDate, error := time.Parse("2006-01-02", start)
		if error != nil {
			return nil, ErrorBadQueryParams
		}
		endDate, error := time.Parse("2006-01-02", end)
		if error != nil {
			return nil, ErrorBadQueryParams
		}
		query.StartDate = startDate
		query.EndDate = endDate
	}

	if len(date) == 0 && len(start) == 0 && len(end) == 0 {
		query.Date = time.Now()
	}

	return &query, nil
}

// ApodQueryOutput is the output from an Apod Query
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

func (a APODClient) GetContent(date, start, end string) ([]models.Article, error) {
	// Parse queryParams
	queryParams, err := MakeApodQueryInput(date, start, end)
	if err != nil {
		return nil, err
	}

	// Do request to APOD service
	rows, err := a.Query(queryParams)
	if err != nil {
		return nil, err
	}

	// Make store models
	var result []models.Article
	for _, row := range rows {
		a := MakeArticle(row)
		result = append(result, *a)
	}
	return result, nil
}

func (a APODClient) SavePicture(url, key string) error {
	value, err := a.GetPicture(url)
	if err != nil {
		return err
	}
	return a.store.Set(key, value)
}

func MakeArticle(a ApodOutput) *models.Article {
	return &models.Article{
		Title:       a.Title,
		Date:        a.Date,
		Explanation: a.Explanation,
		Url:         a.Url,
	}
}

// Make query to APOD service
func (a APODClient) Query(queryParams *ApodQueryInput) ([]ApodOutput, error) {
	var queryUrl string
	if !queryParams.Date.IsZero() {
		if !queryParams.StartDate.IsZero() || !queryParams.EndDate.IsZero() {
			return nil, ErrorBadQueryParams
		}
		queryUrl += fmt.Sprintf("&date=%s", queryParams.Date.Format("2006-01-02"))
	}

	if !queryParams.StartDate.IsZero() {
		queryUrl += fmt.Sprintf("&start_date=%s", queryParams.StartDate.Format("2006-01-02"))
	}

	if !queryParams.EndDate.IsZero() {
		queryUrl += fmt.Sprintf("&end_date=%s", queryParams.EndDate.Format("2006-01-02"))
	}

	resp, err := http.Get(a.apodUrl + queryUrl)
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
	var queryRows []ApodOutput

	if !queryParams.Date.IsZero() {
		var queryRow ApodOutput
		err = json.Unmarshal(body, &queryRow)
		queryRows = append(queryRows, queryRow)
	} else {
		err = json.Unmarshal(body, &queryRows)
	}
	return queryRows, err
}

// Func Download picture from APOD
func (a APODClient) GetPicture(url string) ([]byte, error) {
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
