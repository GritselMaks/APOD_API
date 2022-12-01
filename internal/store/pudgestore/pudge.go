package pudgestore

import (
	"github.com/recoilme/pudge"
)

type Pudge struct {
	handle *pudge.Db
}

func Open(path string) (*Pudge, error) {
	cfg := &pudge.Config{
		SyncInterval: 0}

	db, err := pudge.Open(path, cfg)
	if err != nil {
		return nil, err
	}
	return &Pudge{handle: db}, nil
}

func (p *Pudge) Set(key string, value []byte) error {
	return p.handle.Set(key, value)
}

func (p *Pudge) Get(key string) ([]byte, error) {
	var data []byte
	err := p.handle.Get(key, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *Pudge) GetList() ([][]byte, error) {
	k, err := p.handle.Keys(nil, 0, 0, false)
	if err != nil {
		return nil, err
	}
	var pictures [][]byte
	for _, key := range k {
		var picture []byte
		err := p.handle.Get(key, &picture)
		if err == nil {
			pictures = append(pictures, picture)
		}
	}

	for i, j := 0, len(pictures)-1; i < j; i, j = i+1, j-1 {
		pictures[i], pictures[j] = pictures[j], pictures[i]
	}
	return pictures, nil
}
