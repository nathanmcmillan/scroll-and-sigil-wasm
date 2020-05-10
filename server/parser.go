package main

import (
	"container/list"
	"strings"
)

// ParserRead func
func ParserRead(str []byte) map[string]interface{} {
	data := make(map[string]interface{})
	stack := list.New()
	stack.PushFront(data)
	var key strings.Builder
	var value strings.Builder
	state := "key"
	num := len(str)
	for i := 0; i < num; i++ {
		c := str[i]
		if c == ':' {
			state = "value"
		} else if c == ',' {
			pc := str[i-1]
			if pc != '}' && pc != ']' {
				front := stack.Front().Value
				switch front.(type) {
				case *Array:
					array := front.(*Array)
					array.data = append(array.data, value.String())
				default:
					front.(map[string]interface{})[key.String()] = value.String()
					key.Reset()
					state = "key"
				}
				value.Reset()
			}
		} else if c == '{' {
			dict := make(map[string]interface{})
			front := stack.Front().Value
			switch front.(type) {
			case *Array:
				array := front.(*Array)
				array.data = append(array.data, dict)
				state = "key"
			default:
				front.(map[string]interface{})[key.String()] = dict
				key.Reset()
			}
			stack.PushFront(dict)
		} else if c == '[' {
			list := NewArray()
			front := stack.Front().Value
			switch front.(type) {
			case *Array:
				array := front.(*Array)
				array.data = append(array.data, list)
			default:
				front.(map[string]interface{})[key.String()] = list
				key.Reset()
			}
			stack.PushFront(list)
			state = "value"
		} else if c == '}' {
			pc := str[i-1]
			if pc != ',' && pc != '{' && pc != ']' && pc != '}' {
				front := stack.Front().Value
				front.(map[string]interface{})[key.String()] = value.String()
				key.Reset()
				value.Reset()
			}
			stack.Remove(stack.Front())
			front := stack.Front().Value
			switch front.(type) {
			case *Array:
				state = "value"
			default:
				state = "key"
			}
		} else if c == ']' {
			pc := str[i-1]
			if pc != ',' && pc != '[' && pc != ']' && pc != '}' {
				front := stack.Front().Value
				array := front.(*Array)
				array.data = append(array.data, value.String())
				value.Reset()
			}
			stack.Remove(stack.Front())
			front := stack.Front().Value
			switch front.(type) {
			case *Array:
				state = "value"
			default:
				state = "key"
			}
		} else if state == "key" {
			key.WriteByte(c)
		} else {
			value.WriteByte(c)
		}
	}
	pc := str[num-1]
	if pc != ',' && pc != ']' && pc != '}' {
		front := stack.Front().Value
		front.(map[string]interface{})[key.String()] = value.String()
	}
	return data
}
