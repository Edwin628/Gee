package geeorm

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestEngine(t *testing.T) {
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Error("open db failed: ", err)
	}
	engine.Close()
}
