package jayson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverseSlice(t *testing.T) {

	for _, item := range []struct {
		input  []int
		output []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{[]int{1, 2, 3, 4, 5, 6}, []int{6, 5, 4, 3, 2, 1}},
		{[]int{1, 2, 3, 4, 5, 6, 7}, []int{7, 6, 5, 4, 3, 2, 1}},
		{[]int{}, []int{}},
		{nil, []int{}},
	} {
		result := reverseSlice(item.input)
		assert.Equal(t, item.output, result)
	}
}
