package utils

import (
	"reflect"
	"testing"
)

func TestPreparePictureUrl(t *testing.T) {
	type args struct {
		host string
		port string
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				host: "127.0.0.1",
				port: "8080",
				path: "url",
			},
			want: "http://127.0.0.1:8080/picture/url",
		}, {
			name: "case 2",
			args: args{
				host: "",
				port: "8080",
				path: "url",
			},
			want: "http://localhost:8080/picture/url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PreparePictureUrl(tt.args.host, tt.args.port, tt.args.path); got != tt.want {
				t.Errorf("PreparePictureUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_preparePath(t *testing.T) {
	type args struct {
		path string
	}
	res := [3]string{"~/1.txt", "", "/home/sasha/1.txt"}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{
			name:    "case_1",
			args:    args{path: "~/1.txt"},
			want:    &res[0],
			wantErr: false,
		},
		{
			name:    "case_2",
			args:    args{path: ""},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "case_3",
			args:    args{path: "/home/maks/1.txt"},
			want:    &res[2],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := preparePath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("preparePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !reflect.DeepEqual(&got, &tt.want) {
					t.Errorf("preparePath() = %v, want %v", &got, &tt.want)
				}
			}
		})
	}
}
