package pudgestore

import (
	"reflect"
	"testing"
)

func TestPudge_Set(t *testing.T) {
	pudge, teardown := TestPudgeStore(t)
	defer teardown()

	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name    string
		p       *Pudge
		args    args
		wantErr bool
	}{
		{
			name:    "NoError",
			p:       pudge,
			args:    args{key: "2018-22-10", value: []byte("test")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Pudge.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPudge_Get(t *testing.T) {
	pudge, teardown := TestPudgeStore(t)
	defer teardown()
	pudge.Set("key", []byte("value"))

	tests := []struct {
		name    string
		p       *Pudge
		key     string
		want    []byte
		wantErr bool
	}{
		{
			name:    "NoError",
			p:       pudge,
			key:     "key",
			want:    []byte("value"),
			wantErr: false,
		},
		{
			name:    "Error",
			p:       pudge,
			key:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Get(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pudge.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Pudge.Get() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}

func TestPudge_GetList(t *testing.T) {
	pudge, teardown := TestPudgeStore(t)
	defer teardown()
	value := []byte("value")
	value2 := []byte("value2")
	pudge.Set("key", value)
	pudge.Set("key2", value2)
	tests := []struct {
		name    string
		p       *Pudge
		want    [][]byte
		wantErr bool
	}{
		{
			name:    "NoError",
			p:       pudge,
			want:    [][]byte{value, value2},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetList()
			if (err != nil) != tt.wantErr {
				t.Errorf("Pudge.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Pudge.GetList() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
