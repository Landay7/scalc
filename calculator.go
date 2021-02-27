package main

import (
	"sort"
	"errors"
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
func calcDeepestSubexpression(operations *[]*operator, depth int, operands *[]numberSet) error {
	if len(*operations) == 0 {
		return errors.New("No operator")
	}
	operation := (*operations)[len(*operations) - 1]
	*operations = (*operations)[:len(*operations) - 1]
	var sets []sortedSet
	for len(*operands) > 0 {
		nums := (*operands)[len(*operands) - 1]
		if depth < nums.depth {
			panic("not processing deepest expression")
		}
		if nums.depth != depth {
			break
		}
		sets = append(sets, nums.values)
		*operands = (*operands)[:len(*operands) - 1]
	}
	result := сalculateLeaf(operation, sets)
	*operands = append(*operands, numberSet{result, depth - 1})
	return nil
}

// calculate takes an expression that is split into tokens
// Example:
//     [ LE 2 a.txt [ GR 1 b.txt c.txt ] ]
// Is passed as:
//     []string{"[" "LE" "2" "a.txt" "[" "GR" "1" "b.txt" "c.txt" "]" "]"}
func calculate(tokens []string) ([]int, error) {
	var operations []*operator
	var sets []numberSet
	var fCache = newFileCache()
	depth := 0
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		switch token {
		case "[":
			depth++
		case "]":
			if err := calcDeepestSubexpression(&operations, depth, &sets); err != nil{
				return nil, err
			}
			depth--
		case "GR", "EQ", "LE":
			i++
			n := tokens[i]
			operator, err := newOperator(token, n)
			if err != nil {
				return nil, err
			}
			operations = append(operations, operator)
		default:
			setData, err := fCache.get(token)
			if err != nil {
				return nil, err
			}
			sets = append(sets, numberSet{setData, depth})
		}
	}
	if depth != 0 {
		return nil, errors.New("Expression not valid")
	}
	if len(sets) != 1 {
		return nil, errors.New("Something went wrong")
	}
	vals := sets[0].values
	sort.Ints(vals)
	return vals, nil
}
