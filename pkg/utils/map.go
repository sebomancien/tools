package utils

import (
	"cmp"
	"iter"
	"slices"
)

type IndexKey[K cmp.Ordered] struct {
	Index int
	Key   K
}

func IterateMapSorted[T any, K cmp.Ordered](m map[K]T) iter.Seq2[IndexKey[K], T] {
	return func(yield func(index IndexKey[K], value T) bool) {
		// Get the map keys
		var keys []K
		for key := range m {
			keys = append(keys, key)
		}

		// Order the keys to have deterministic order in slice
		slices.Sort(keys)

		// Create the slice, applying the convert function foreach element
		for i, key := range keys {
			if !yield(IndexKey[K]{Index: i, Key: key}, m[key]) {
				return
			}
		}
	}
}
