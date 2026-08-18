package task

import (
	"sync"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	tasks  map[int64]Task
	nextID int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks:  make(map[int64]Task),
		nextID: 1,
	}
}
