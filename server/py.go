package main

import (
	"crypto/rand"
	"fmt"
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

// UUID func
func UUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
