package main

import (
	"fmt"
	"strconv"
)

type opType int

const (
	GR opType = iota
	EQ
	LE
)

var opTypeToStr = map[string]opType{
	"GR": GR,
	"EQ": EQ,
	"LE": LE,
}

type operator struct {
	op opType
	n  int
}

func newOperator(opTypeStr, numStr string) *operator {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		panic(fmt.Sprintf("Can not parse: %s", numStr))
	}
	op, ok := opTypeToStr[opTypeStr]
	if !ok {
		panic(fmt.Sprintf("Invalid operorator type %q", opTypeStr))
	}
	return &operator{op, num}
}
