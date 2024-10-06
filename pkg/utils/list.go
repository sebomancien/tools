package utils

import (
	"iter"
	"slices"
)

type SortedList[T any] struct {
	capacity int
	array    []T
	compare  func(a T, b T) int
}

func NewSortedList[T any](capacity int, compare func(a T, b T) int) *SortedList[T] {
	return &SortedList[T]{
		capacity: capacity,
		array:    []T{},
		compare:  compare,
	}
}

func (s *SortedList[T]) Insert(value T) {
	index, _ := slices.BinarySearchFunc(s.array, value, s.compare)
	s.array = slices.Insert(s.array, index, value)

	// Make sure the list does not exceed the maximum capacity
	if len(s.array) > s.capacity {
		s.array = s.array[len(s.array)-s.capacity:]
	}
}

func (s *SortedList[T]) Values() iter.Seq[T] {
	return slices.Values(s.array)
}

func (s *SortedList[T]) All() iter.Seq2[int, T] {
	return slices.All(s.array)
}
