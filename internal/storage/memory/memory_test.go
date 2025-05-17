package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRU(t *testing.T) {
	t.Run("basic add and get", func(t *testing.T) {
		cache := NewLRUStorage(2)

		// Добавляем первый элемент
		cache.addToCache("key1", []byte("value1"))
		val, ok := cache.cache["key1"]
		assert.True(t, ok)
		assert.Equal(t, []byte("value1"), val)
		assert.Equal(t, []string{"key1"}, cache.order)

		// Добавляем второй элемент
		cache.addToCache("key2", []byte("value2"))
		assert.Equal(t, 2, len(cache.cache))
		assert.Equal(t, []string{"key2", "key1"}, cache.order)
	})

	t.Run("eviction when capacity exceeded", func(t *testing.T) {
		cache := NewLRUStorage(2)

		cache.addToCache("key1", []byte("value1"))
		cache.addToCache("key2", []byte("value2"))
		cache.addToCache("key3", []byte("value3")) // Должен вытеснить key1

		assert.Equal(t, 2, len(cache.cache))
		_, ok := cache.cache["key1"]
		assert.False(t, ok)
		assert.Equal(t, []string{"key3", "key2"}, cache.order)
	})

	t.Run("move to front on access", func(t *testing.T) {
		cache := NewLRUStorage(3)

		cache.addToCache("key1", []byte("value1"))
		cache.addToCache("key2", []byte("value2"))
		cache.addToCache("key3", []byte("value3"))

		// Доступ к key2 должен переместить его в начало
		cache.moveToFront("key2")
		assert.Equal(t, []string{"key2", "key3", "key1"}, cache.order)

		// Доступ к key1 должен переместить его в начало
		cache.moveToFront("key1")
		assert.Equal(t, []string{"key1", "key2", "key3"}, cache.order)
	})

	t.Run("clear cache", func(t *testing.T) {
		cache := NewLRUStorage(2)

		cache.addToCache("key1", []byte("value1"))
		cache.addToCache("key2", []byte("value2"))

		cache.Clear()

		assert.Equal(t, 0, len(cache.cache))
		assert.Equal(t, 0, len(cache.order))
	})
}
