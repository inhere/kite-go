package lcstorage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type FileStorage struct {
	mu       sync.RWMutex
	data     map[string]any
	filePath string
}

func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		data:     make(map[string]any),
		filePath: filePath,
	}
}

func (s *FileStorage) Get(key string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, exists := s.data[key]
	if !exists {
		return nil, fmt.Errorf("key %q not found", key)
	}
	return val, nil
}

func (s *FileStorage) Set(key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return s.doStore()
}

func (s *FileStorage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return s.doStore()
}

func (s *FileStorage) Clear() error {
	s.mu.Lock()
	s.data = make(map[string]any)
	s.mu.Unlock()

	return s.Store()
}

func (s *FileStorage) Store() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.doStore()
}

func (s *FileStorage) doStore() error {
	dataBytes, err := json.Marshal(s.data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err1 := os.WriteFile(s.filePath, dataBytes, 0644); err1 != nil {
		return fmt.Errorf("failed to write storage file: %w", err)
	}
	return nil
}

func (s *FileStorage) Restore() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dataBytes, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read storage file: %w", err)
	}

	if len(dataBytes) > 0 {
		if err := json.Unmarshal(dataBytes, &s.data); err != nil {
			return fmt.Errorf("failed to unmarshal storage data: %w", err)
		}
	}
	return nil
}

func (s *FileStorage) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *FileStorage) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

func (s *FileStorage) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data) == 0
}

func (s *FileStorage) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.data[key]
	return exists
}
