package util

import (
	"sort"
)

type SortRange[T any] struct {
	Map map[string]T
}

func NewSortRange[T any](_map map[string]T) SortRange[T] {
	return SortRange[T]{Map: _map}
}

func (r SortRange[T]) Iter() <-chan struct {
	Key   string
	Value T
} {
	ch := make(chan struct {
		Key   string
		Value T
	})

	// Step 1: Convert the map to a slice of key-value pairs
	pairs := make([]struct {
		Key   string
		Value T
	}, len(r.Map))

	i := 0
	for key, value := range r.Map {
		pairs[i] = struct {
			Key   string
			Value T
		}{key, value}
		i++
	}

	// Step 2: Sort the slice by key
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key < pairs[j].Key
	})

	go func() {
		defer close(ch)

		for _, v := range pairs {
			ch <- v
		}
	}()

	return ch
}
