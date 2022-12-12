package apod

import (
	"testing"
	"time"

	"github.com/GritselMaks/BT_API/internal/store/testbinarstore"
)

func TestMakeApodQueryInput(t *testing.T) {
	type args struct {
		date  string
		start string
		end   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "noError_case1",
			args: args{
				date:  "2012-11-10",
				start: "",
				end:   ""},
			wantErr: false,
		},
		{
			name: "noError_case2",
			args: args{
				date:  "",
				start: "2012-11-10",
				end:   "2012-11-11"},
			wantErr: false,
		},
		{
			name: "noError_case3",
			args: args{
				date:  "",
				start: "",
				end:   ""},
			wantErr: false,
		},
		{
			name: "Error_case1",
			args: args{
				date:  "2012-11-10",
				start: "2012-11-9",
				end:   "2012-11-11"},
			wantErr: true,
		},
		{
			name: "Error_case3",
			args: args{
				date:  "2012-11-10",
				start: "2012-11-9",
				end:   ""},
			wantErr: true,
		},
		{
			name: "Error_case3",
			args: args{
				date:  "2012-11-10",
				start: "",
				end:   "2012-11-9"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MakeApodQueryInput(tt.args.date, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test %s. MakeApodQueryInput() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}

		})
	}
}

func TestAPODClient_Query(t *testing.T) {
	st := testbinarstore.TestBinarStore()
	apodCLient := NewAPOD("", st)
	tests := []struct {
		name        string
		a           APODClient
		queryParams *ApodQueryInput
		wantErr     bool
	}{
		{
			name:        "noError_case1",
			a:           apodCLient,
			queryParams: &ApodQueryInput{Date: time.Now(), StartDate: time.Time{}, EndDate: time.Time{}},
			wantErr:     false,
		},
		{
			name:        "noError_case2",
			a:           apodCLient,
			queryParams: &ApodQueryInput{Date: time.Time{}, StartDate: time.Date(2022, 12, 10, 0, 0, 0, 0, time.Local), EndDate: time.Now()},
			wantErr:     false,
		},
		{
			name:        "Error_case1",
			a:           apodCLient,
			queryParams: &ApodQueryInput{Date: time.Now(), StartDate: time.Now(), EndDate: time.Now()},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.a.Query(tt.queryParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("APODClient.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAPODClient_GetContent(t *testing.T) {
	st := testbinarstore.TestBinarStore()
	apodCLient := NewAPOD("", st)
	type args struct {
		date  string
		start string
		end   string
	}
	tests := []struct {
		name    string
		a       APODClient
		args    args
		wantErr bool
	}{
		{
			name:    "noError_case1",
			a:       apodCLient,
			args:    args{date: "2016-11-12"},
			wantErr: false,
		},
		{
			name:    "noError_case2",
			a:       apodCLient,
			args:    args{start: "2022-11-12", end: "2022-12-11"},
			wantErr: false,
		},
		{
			name:    "Error_case1",
			a:       apodCLient,
			args:    args{date: "2016-11-11", start: "2016-11-12", end: "2017-11-11"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.a.GetContent(tt.args.date, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("APODClient.GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
