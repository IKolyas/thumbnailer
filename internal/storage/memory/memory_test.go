package memory

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRU(t *testing.T) {
	// Create temp directory for tests
	tempDir, err := os.MkdirTemp("", "lru_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("basic add and get", func(t *testing.T) {
		cache, err := NewLRUStorage(2, tempDir)
		assert.NoError(t, err)

		// Add first item
		err = cache.addToCache("key1", []byte("value1"))
		assert.NoError(t, err)
		filePath, ok := cache.cache["key1"]
		assert.True(t, ok)
		assert.FileExists(t, filePath)
		assert.Equal(t, []string{"key1"}, cache.order)

		// Add second item
		err = cache.addToCache("key2", []byte("value2"))
		assert.NoError(t, err)
		assert.Equal(t, 2, len(cache.cache))
		assert.Equal(t, []string{"key2", "key1"}, cache.order)
	})

	t.Run("eviction when capacity exceeded", func(t *testing.T) {
		cache, err := NewLRUStorage(2, tempDir)
		assert.NoError(t, err)

		err = cache.addToCache("key1", []byte("value1"))
		assert.NoError(t, err)
		err = cache.addToCache("key2", []byte("value2"))
		assert.NoError(t, err)
		err = cache.addToCache("key3", []byte("value3")) // Should evict key1
		assert.NoError(t, err)

		assert.Equal(t, 2, len(cache.cache))
		_, ok := cache.cache["key1"]
		assert.False(t, ok)
		assert.Equal(t, []string{"key3", "key2"}, cache.order)

		// Check file was deleted
		_, err = os.Stat(filepath.Join(tempDir, "key1"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("move to front on access", func(t *testing.T) {
		cache, err := NewLRUStorage(3, tempDir)
		assert.NoError(t, err)

		err = cache.addToCache("key1", []byte("value1"))
		assert.NoError(t, err)
		err = cache.addToCache("key2", []byte("value2"))
		assert.NoError(t, err)
		err = cache.addToCache("key3", []byte("value3"))
		assert.NoError(t, err)

		// Access key2 should move it to front
		cache.moveToFront("key2")
		assert.Equal(t, []string{"key2", "key3", "key1"}, cache.order)

		// Access key1 should move it to front
		cache.moveToFront("key1")
		assert.Equal(t, []string{"key1", "key2", "key3"}, cache.order)
	})

	t.Run("clear cache", func(t *testing.T) {
		cache, err := NewLRUStorage(2, tempDir)
		assert.NoError(t, err)

		err = cache.addToCache("key1", []byte("value1"))
		assert.NoError(t, err)
		err = cache.addToCache("key2", []byte("value2"))
		assert.NoError(t, err)

		err = cache.Clear()
		assert.NoError(t, err)

		assert.Equal(t, 0, len(cache.cache))
		assert.Equal(t, 0, len(cache.order))

		// Check files were deleted
		_, err = os.Stat(filepath.Join(tempDir, "key1"))
		assert.True(t, os.IsNotExist(err))
		_, err = os.Stat(filepath.Join(tempDir, "key2"))
		assert.True(t, os.IsNotExist(err))
	})
}
