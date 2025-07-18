package set

import (
	"cmp"
	"slices"
)

type OrderedSet[T cmp.Ordered] struct {
	items []T
	set   map[T]struct{}
}

func NewOrderedSet[T cmp.Ordered]() *OrderedSet[T] {
	return &OrderedSet[T]{set: map[T]struct{}{}}
}

func FromSlice[T cmp.Ordered](slice []T) *OrderedSet[T] {
	s := NewOrderedSet[T]()
	s.AddSlice(slice)
	return s
}

func (s *OrderedSet[T]) Items() []T {
	return s.items
}

func (s *OrderedSet[T]) SortedItems() []T {
	slices.Sort(s.items)
	return s.items
}

func (s *OrderedSet[T]) Contains(item T) bool {
	_, ok := s.set[item]
	return ok
}

func (s *OrderedSet[T]) AddMaps(sets ...map[T]struct{}) {
	for _, set := range sets {
		s.AddMap(set)
	}
}

func (s *OrderedSet[T]) AddMap(set map[T]struct{}) {
	for item := range set {
		if _, ok := s.set[item]; !ok {
			s.set[item] = struct{}{}
			s.items = append(s.items, item)
		}
	}
}

func (s *OrderedSet[T]) AddSlices(slices ...[]T) {
	for _, items := range slices {
		s.AddSlice(items)
	}
}

func (s *OrderedSet[T]) AddSlice(items []T) {
	for _, item := range items {
		if _, ok := s.set[item]; !ok {
			s.set[item] = struct{}{}
			s.items = append(s.items, item)
		}
	}
}

func (s *OrderedSet[T]) AddItems(items ...T) {
	s.AddSlice(items)
}
