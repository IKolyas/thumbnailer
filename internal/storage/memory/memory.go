package memory

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/IKolyas/image-previewer/internal/core/image"
	"github.com/IKolyas/image-previewer/internal/storage/source"
)

type LRUStorage struct {
	capacity   int
	cache      map[string]string // key -> filepath
	order      []string
	mu         sync.Mutex
	storageDir string
}

func NewLRUStorage(capacity int, storageDir string) (*LRUStorage, error) {
	if err := os.MkdirAll(storageDir, 0o755); err != nil {
		return nil, err
	}

	return &LRUStorage{
		capacity: capacity,

		cache:      make(map[string]string),
		order:      make([]string, 0, capacity),
		storageDir: storageDir,
	}, nil
}

func (s *LRUStorage) Get(ctx context.Context, imgData *image.ImgData) ([]byte, error) {
	key := imgData.String()
	s.mu.Lock()
	defer s.mu.Unlock()

	if filePath, ok := s.cache[key]; ok {
		s.moveToFront(key)
		return os.ReadFile(filePath)
	}

	data, err := source.Get(ctx, imgData)
	if err != nil {
		return nil, err
	}

	err = s.addToCache(key, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *LRUStorage) addToCache(key string, imgData []byte) error {
	if len(s.order) >= s.capacity {
		oldest := s.order[len(s.order)-1]
		if err := s.removeFile(oldest); err != nil {
			return err
		}
		delete(s.cache, oldest)
		s.order = s.order[:len(s.order)-1]
	}

	filePath := filepath.Join(s.storageDir, key)

	if err := os.WriteFile(filePath, imgData, 0o600); err != nil {
		return err
	}

	s.cache[key] = filePath
	s.order = append([]string{key}, s.order...)

	return nil
}

func (s *LRUStorage) moveToFront(key string) {
	for i, k := range s.order {
		if k == key {
			s.order = slices.Delete(s.order, i, i+1)
			s.order = append([]string{key}, s.order...)
			break
		}
	}
}

func (s *LRUStorage) removeFile(key string) error {
	filePath, ok := s.cache[key]
	if !ok {
		return nil
	}
	return os.Remove(filePath)
}

func (s *LRUStorage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, filePath := range s.cache {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	s.cache = make(map[string]string)
	s.order = make([]string, 0, s.capacity)

	return nil
}
