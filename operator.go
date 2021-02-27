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

func newOperator(opTypeStr, numStr string) (*operator, error) {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return nil, err
	}
	op, ok := opTypeToStr[opTypeStr]
	if !ok {
		return nil, fmt.Errorf("Invalid operorator type %q", opTypeStr)
	}
	return &operator{op, num}, nil
}
