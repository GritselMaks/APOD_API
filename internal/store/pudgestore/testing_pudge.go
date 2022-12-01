package pudgestore

import (
	"testing"

	"github.com/recoilme/pudge"
)

func TestPudgeStore(t *testing.T) (*Pudge, func()) {
	path := "~/test_pudge.db"
	cfg := &pudge.Config{
		SyncInterval: 0}
	db, err := pudge.Open(path, cfg)
	if err != nil {
		t.Fatal(err)

	}
	return &Pudge{handle: db}, func() {
		db.DeleteFile()
	}
}
