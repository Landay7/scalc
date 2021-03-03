package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
)

func Test_calculate_1(t *testing.T) {
	command1 := "[ LE 2 testdata/a.txt [ GR 1 testdata/b.txt testdata/c.txt ] ]"
	set1, tokens, err := calculate(strings.Fields(command1), 0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(tokens))
	assert.Equal(t, set1, []int{1, 4})
}

func Test_calculate_2(t *testing.T) {
	command2 := "[ EQ 1 [ GR 1 testdata/a.txt testdata/c.txt ] [ GR 1 testdata/b.txt testdata/c.txt ] ]"
	set2, tokens, err := calculate(strings.Fields(command2), 0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(tokens))
	assert.Equal(t, set2, []int{1, 4})
}

func Test_calculate_3(t *testing.T) {
	command3 := "[ GR 1 testdata/c.txt [ EQ 3 testdata/a.txt testdata/a.txt testdata/b.txt ] ]"
	set3, tokens, err := calculate(strings.Fields(command3), 0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(tokens))
	assert.Equal(t, set3, []int{2, 3})
}

func Test_calculate_mismatched_brackets(t *testing.T) {
	commandWithNoClosedBracket := "[ LE 2 testdata/a.txt [ GR 1 testdata/b.txt testdata/c.txt ] "
	set, tokens, err := calculate(strings.Fields(commandWithNoClosedBracket), 0)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(tokens))
	assert.Nil(t, set)
}

func Test_calculate_no_file(t *testing.T) {
	command1 := "[ LE 2 testdata/not_such_file.txt [ GR 1 testdata/b.txt testdata/c.txt ] ]"
	set, tokens, err := calculate(strings.Fields(command1), 0)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(tokens))
	assert.Nil(t, set)
}

func Test_filterByOp(t *testing.T) {
	countMap := map[int]int{
		1: 1,
		2: 2,
		3: 1,
		4: 2,
		5: 3,
		7: 3,
	}

	operation1, _ := newOperator("GR", "2")
	result1 := filterByOp(operation1, countMap)
	sort.Ints(result1)
	assert.Equal(t, result1, []int{5, 7})

	operation2, _ := newOperator("EQ", "2")
	result2 := filterByOp(operation2, countMap)
	sort.Ints(result2)
	assert.Equal(t, result2, []int{2, 4})

	operation3, _ := newOperator("LE", "2")
	result3 := filterByOp(operation3, countMap)
	sort.Ints(result3)
	assert.Equal(t, result3, []int{1, 3})
}

func Test_makeCountMap1(t *testing.T) {
	sets := []sortedSet{
		{1, 2, 3},
		{3, 4, 4, 4},
		{3, 4, 5, 6},
	}
	var countMap map[int]int
	for _, s := range sets {
		countMap = addToCountMap(countMap, s)
	}
	expectedResult := map[int]int{
		1: 1,
		2: 1,
		3: 3,
		4: 2,
		5: 1,
		6: 1,
	}
	assert.Equal(t, expectedResult, countMap)
}

func Test_makeCountMap2(t *testing.T) {
	sets := []sortedSet{
		{3, 4},
		{5, 6},
		{7, 8},
	}
	var countMap map[int]int
	for _, s := range sets {
		countMap = addToCountMap(countMap, s)
	}
	expectedResult := map[int]int{
		3: 1,
		4: 1,
		5: 1,
		6: 1,
		7: 1,
		8: 1,
	}
	assert.Equal(t, expectedResult, countMap)
}

func Test_сalculateLeaf_GR(t *testing.T) {
	sets := []sortedSet{
		{1, 2, 3},
		{3, 4, 4, 4},
		{3, 4, 5, 6},
	}
	var countMap map[int]int
	for _, s := range sets {
		countMap = addToCountMap(countMap, s)
	}
	op := &operator{GR, 2}
	result1 := сalculateLeaf(op, countMap)
	assert.Equal(t, result1, []int{3})
}

func Test_сalculateLeaf_EQ(t *testing.T) {
	sets := []sortedSet{
		{3},
		{3, 4, 4, 4},
		{},
		{3, 4, 5, 6},
	}
	var countMap map[int]int
	for _, s := range sets {
		countMap = addToCountMap(countMap, s)
	}
	op := &operator{EQ, 2}
	result2 := сalculateLeaf(op, countMap)
	assert.Equal(t, result2, []int{4})
}

func Test_сalculateLeaf_LE(t *testing.T) {
	sets := []sortedSet{
		{1, 2, 3},
		{3, 4, 4, 4},
		{3, 4, 5, 6},
	}
	var countMap map[int]int
	for _, s := range sets {
		countMap = addToCountMap(countMap, s)
	}
	op := &operator{LE, 2}
	result := сalculateLeaf(op, countMap)
	sort.Ints(result)
	assert.Equal(t, result, []int{1, 2, 5, 6})
}

func Test_parseBrackets(t *testing.T) {
	badTestData := [][]string {
		{},
		{""},
		{"["},
		{"]", "["},
		{"[", "["},
	}
	for _, testData := range badTestData {
		_, err := parseBrackets(testData)
		assert.Error(t, err)
	}
}
