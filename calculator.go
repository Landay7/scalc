package main

import (
	"github.com/golang-collections/collections/stack"
	"sort"
)

type sortedSet []int

type numberSet struct {
	values []int
	depth  int
}

func makeCountMap(setSlice []sortedSet) map[int]int {
	// maps a int to the number of times that it is found in the sets
	var countMap = make(map[int]int)
	for _, set := range setSlice {
		for i, val := range set {
			// assume the sets are sorted, but may have repeated values
			if i > 0 && set[i-1] == val {
				continue
			}
			if _, ok := countMap[val]; ok {
				countMap[val]++
			} else {
				countMap[val] = 1
			}
		}
	}
	return countMap
}

// Returns non-sorted set of unique values
func filterByOp(operation *operator, countMap map[int]int) []int {
	var result []int
	for key, value := range countMap {
		if operation.op == GR {
			if value > operation.n {
				result = append(result, key)
			}
		} else if operation.op == EQ {
			if value == operation.n {
				result = append(result, key)
			}

		} else if operation.op == LE {
			if value < operation.n {
				result = append(result, key)
			}
		}
	}
	return result
}

func сalculateLeaf(op *operator, sets []sortedSet) []int {
	countMap := makeCountMap(sets)
	return filterByOp(op, countMap)
}

// depth is the deepest depth at the current time
func calcDeepestSubexpression(operations *stack.Stack, depth int, operands *stack.Stack) {
	operation := operations.Pop()
	if operation == nil {
		panic("No operator!")
	}
	var sets []sortedSet
	for operand := operands.Peek(); operand != nil; operand = operands.Peek() {
		nums := operand.(numberSet)
		if depth < nums.depth {
			panic("not processing deepest expression")
		}
		if nums.depth != depth {
			break
		}
		sets = append(sets, nums.values)
		operands.Pop()
	}
	result := сalculateLeaf(operation.(*operator), sets)
	operands.Push(numberSet{result, depth - 1})
}

// calculate takes an expression that is split into tokens
// Example:
//     [ LE 2 a.txt [ GR 1 b.txt c.txt ] ]
// Is passed as:
//     []string{"[" "LE" "2" "a.txt" "[" "GR" "1" "b.txt" "c.txt" "]" "]"}
func calculate(tokens []string) []int {
	operations := stack.New()
	sets := stack.New()
	var fCache = newFileCache()
	depth := 0
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		switch token {
		case "[":
			depth++
		case "]":
			calcDeepestSubexpression(operations, depth, sets)
			depth--
		case "GR", "EQ", "LE":
			i++
			n := tokens[i]
			operations.Push(newOperator(token, n))
		default:
			setData, err := fCache.get(token)
			if err != nil {
				panic(err)
			}
			sets.Push(numberSet{setData, depth})
		}
	}
	vals := sets.Pop().(numberSet).values
	sort.Ints(vals)
	return vals
}
