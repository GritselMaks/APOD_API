package testbinarstore

import "github.com/GritselMaks/BT_API/internal/store"

type BinarStore struct {
	handle map[string][]byte
}

func TestBinarStore() *BinarStore {
	return &BinarStore{handle: make(map[string][]byte)}
}

func (b *BinarStore) Set(key string, value []byte) error {
	b.handle[key] = value
	return nil
}

func (b *BinarStore) Get(key string) ([]byte, error) {
	v, ok := b.handle[key]
	if !ok {
		return nil, store.ErrNotFound
	}
	return v, nil
}

func (b *BinarStore) GetList() ([][]byte, error) {
	var pictures [][]byte
	for _, v := range b.handle {
		pictures = append(pictures, v)
	}
	return pictures, nil
}
