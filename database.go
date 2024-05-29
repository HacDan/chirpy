package database

import (
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	file, err := os.OpenFile(path)
	if errors.Is(err, os.ErrNotExist) {
		os.Create(path)
		file, _ = os.OpenFile(path)
	}
}
