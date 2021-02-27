package main

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

func Test_calculate(t *testing.T) {
	command1 := "[ LE 2 testdata/a.txt [ GR 1 testdata/b.txt testdata/c.txt ] ]"
	set1, err1 := calculate(strings.Fields(command1))
	assert.Nil(t, err1)
	assert.Equal(t, set1, []int{1, 4})

	command2 := "[ EQ 1 [ GR 1 testdata/a.txt testdata/c.txt ] [ GR 1 testdata/b.txt testdata/c.txt ] ]"
	set2, err2 := calculate(strings.Fields(command2))
	assert.Nil(t, err2)
	assert.Equal(t, set2, []int{1, 4})

	command3 := "[ GR 1 testdata/c.txt [ EQ 3 testdata/a.txt testdata/a.txt testdata/b.txt ] ]"
	set3, err3 := calculate(strings.Fields(command3))
	assert.Nil(t, err3)
	assert.Equal(t, set3, []int{2, 3})
}

func Test_calculate_invalid_exprission(t *testing.T) {
	command1 := "[ LE 2 testdata/a.txt [ GR 1 testdata/b.txt testdata/c.txt ] "
	set, err := calculate(strings.Fields(command1))
	assert.NotNil(t, err)
	assert.Nil(t, set)
}

func Test_calculate_no_file(t *testing.T) {
	command1 := "[ LE 2 testdata/not_such_file.txt [ GR 1 testdata/b.txt testdata/c.txt ] ]"
	set, err := calculate(strings.Fields(command1))
	assert.NotNil(t, err)
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
	result := makeCountMap(sets)
	expectedResult := map[int]int{
		1: 1,
		2: 1,
		3: 3,
		4: 2,
		5: 1,
		6: 1,
	}
	assert.Equal(t, expectedResult, result)
}

func Test_makeCountMap2(t *testing.T) {
	sets := []sortedSet{
		{3, 4},
		{5, 6},
		{7, 8},
	}
	result := makeCountMap(sets)
	expectedResult := map[int]int{
		3: 1,
		4: 1,
		5: 1,
		6: 1,
		7: 1,
		8: 1,
	}
	assert.Equal(t, expectedResult, result)
}

func Test_сalculateLeaf_GR(t *testing.T) {
	sets := []sortedSet{
		{1, 2, 3},
		{3, 4, 4, 4},
		{3, 4, 5, 6},
	}
	op := &operator{GR, 2}
	result1 := сalculateLeaf(op, sets)
	assert.Equal(t, result1, []int{3})
}

func Test_сalculateLeaf_EQ(t *testing.T) {
	sets := []sortedSet{
		{3},
		{3, 4, 4, 4},
		{},
		{3, 4, 5, 6},
	}
	op := &operator{EQ, 2}
	result2 := сalculateLeaf(op, sets)
	assert.Equal(t, result2, []int{4})
}

func Test_сalculateLeaf_LE(t *testing.T) {
	sets := []sortedSet{
		{1, 2, 3},
		{3, 4, 4, 4},
		{3, 4, 5, 6},
	}
	op := &operator{LE, 2}
	result := сalculateLeaf(op, sets)
	sort.Ints(result)
	assert.Equal(t, result, []int{1, 2, 5, 6})
}
