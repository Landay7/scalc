package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_fileCache_get(t *testing.T) {
	var cache = newFileCache()
	filename := "testdata/a.txt"

	value1, err := cache.get(filename)
	require.NoError(t, err)
	assert.Equal(t, len(value1), 3)

	values2, _ := cache.get(filename)
	assert.Equal(t, value1, values2)
}

func Test_fileCache_get_NoSuchFile(t *testing.T) {
	cache := newFileCache()
	_, err := cache.get("no-such-file")
	require.Error(t, err)
}
