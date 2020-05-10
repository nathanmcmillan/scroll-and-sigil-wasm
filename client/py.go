package main

import (
	"strconv"
)

// ParseInt func
func ParseInt(value string) int {
	num, _ := strconv.ParseInt(value, 10, 64)
	return int(num)
}

// ParseFloat func
func ParseFloat(value string) float32 {
	num, _ := strconv.ParseFloat(value, 64)
	return float32(num)
}

// Array struct
type Array struct {
	data []interface{}
}

// NewArray func
func NewArray() *Array {
	v := &Array{}
	v.data = make([]interface{}, 0)
	return v
}
