package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type sortedSet []int

func addToCountMap(countMap map[int]int, set []int) map[int]int {
	// maps a int to the number of times that it is found in the sets
	if countMap == nil {
		countMap = make(map[int]int)
	}

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

func сalculateLeaf(op *operator, countMap map[int]int) []int {
	res := filterByOp(op, countMap)
	sort.Ints(res)
	return res
}

// Used for error messages
func tokToExpr(tokens []string) string {
	return strings.Join(tokens, " ")
}

func parseBrackets(tokens []string) ([]string, error) {
	l := len(tokens)
	if l < 2 || tokens[0] != "[" || tokens[l-1] != "]" {
		return nil, fmt.Errorf("expression must start and end with square brackets: (%s)", tokToExpr(tokens))
	}
	return tokens[1 : l-2], nil
}

func parseOperator(tokens []string) (*operator, []string, error) {
	if len(tokens) < 2 {
		return nil, nil, errors.New("expression too short to be valid")
	}
	op, err := newOperator(tokens[0], tokens[1])
	if err != nil {
		return nil, nil, err
	}
	return op, tokens[2:], nil
}

func parseSets(countMap map[int]int, tokens []string, depth int) (map[int]int, []string, error) {
	for len(tokens) > 0 && tokens[0] != "]" {
		var set []int
		var err error
		tok := tokens[0]
		if tok == "[" {
			set, tokens, err = calculate(tokens, depth+1)
		} else {
			set, err = readNumbersInFile(tok)
			tokens = tokens[1:]
		}
		if err != nil {
			return nil, nil, err
		}
		addToCountMap(countMap, set)
	}
	return countMap, tokens, nil
}

// calculate takes an expression that is split into tokens
// Example:
//     [ LE 2 a.txt [ GR 1 b.txt c.txt ] ]
// Is passed as:
//     []string{"[" "LE" "2" "a.txt" "[" "GR" "1" "b.txt" "c.txt" "]" "]"}
func calculate(tokens []string, depth int) ([]int, []string, error) {
	if len(tokens) < 5 {
		return nil, nil, fmt.Errorf("expression too short to be valid: %s", tokToExpr(tokens))
	}
	if tokens[0] != "[" {
		return nil, nil, fmt.Errorf("expression must start with '[' (found '%s')", tokens[0])
	}
	tokens = tokens[1:]
	op, tokens, err := parseOperator(tokens)
	if err != nil {
		return nil, nil, err
	}
	countMap := make(map[int]int)
	countMap, tokens, err = parseSets(countMap, tokens, depth)
	if err != nil {
		return nil, nil, err
	}
	if len(tokens) < 1 || tokens[0] != "]" {
		return nil, nil, errors.New("expression must end with ']'")
	}
	tokens = tokens[1:]
	if depth == 0 && len(tokens) > 0 {
		return nil, nil, fmt.Errorf("unexpected values at end of expression: %s", tokToExpr(tokens))
	}
	return сalculateLeaf(op, countMap), tokens, nil
}
