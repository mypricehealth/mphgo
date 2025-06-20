package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderedSetCreator(t *testing.T) {
	t.Parallel()

	s := NewOrderedSet[int]()
	assert.NotNil(t, s)

	s.AddItems(1, 3, 2, 1, 4)
	assert.Equal(t, []int{1, 3, 2, 4}, s.Items())

	assert.True(t, s.Contains(1))

	s.AddMap(map[int]struct{}{1: {}, 5: {}})
	assert.Equal(t, []int{1, 3, 2, 4, 5}, s.Items())

	s.AddSlice([]int{1, 6, 7})
	assert.Equal(t, []int{1, 3, 2, 4, 5, 6, 7}, s.Items())

	s.AddSlices([]int{1, 10, 11}, []int{1, 8, 9})
	assert.Equal(t, []int{1, 3, 2, 4, 5, 6, 7, 10, 11, 8, 9}, s.Items())

	s.AddMaps(map[int]struct{}{1: {}, 13: {}}, map[int]struct{}{1: {}, 12: {}})
	assert.Equal(t, []int{1, 3, 2, 4, 5, 6, 7, 10, 11, 8, 9, 13, 12}, s.Items())

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, s.SortedItems())
}

func TestFromSlice(t *testing.T) {
	t.Parallel()

	s := FromSlice([]int{1, 3, 2, 1, 4})
	assert.Equal(t, []int{1, 3, 2, 4}, s.Items())
}
