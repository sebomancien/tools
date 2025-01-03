package utils

import (
	"log"
	"slices"
	"testing"
)

func TestSortedList(t *testing.T) {
	data := []int{9, 3, 7, 8, 3, 2, 1}
	capacity := 5
	compare := func(a int, b int) int {
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	}

	list := NewSortedList(capacity, compare)
	for _, v := range data {
		list.Insert(v)
	}

	slices.Sort(data)

	if slices.Compare(list.array, data[len(data)-capacity:]) != 0 {
		log.Fatalf("slices do not match. returned:%v. expected:%v", list.array, data[:capacity])
	}
}
