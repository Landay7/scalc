package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ReadNumbersInFile(t *testing.T) {
	set, err := readNumbersInFile("testdata/a.txt")
	assert.Nil(t, err)
	assert.Equal(t, set, []int{1, 2, 3})

	set, err = readNumbersInFile("testdata/invalid_file.txt")
	assert.Nil(t, set)
	assert.NotNil(t, err)
}
