package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Dummy tests, we should show the Table test
func TestIntIsInSlice(t *testing.T) {
	valueToFind := 69
	elements := []int{1, 2, 3, 4, 85, 69}
	result := Contains(valueToFind, elements)
	assert.True(t, result)
}

func TestIntIsNotInSlice(t *testing.T) {
	valueToFind := 69
	elements := []int{1, 2, 3, 4, 85}
	result := Contains(valueToFind, elements)
	assert.False(t, result)
}

func TestContainsInt(t *testing.T) {
	testCases := []struct {
		Name           string
		ValueToFind    int
		Elements       []int
		ExpectedResult bool
	}{
		{
			Name:           "The slice contains the element",
			ValueToFind:    69,
			Elements:       []int{1, 2, 3, 4, 85, 69},
			ExpectedResult: true,
		},
		{
			Name:           "The slice does not contains the element",
			ValueToFind:    69,
			Elements:       []int{1, 2, 3, 4, 85},
			ExpectedResult: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			result := Contains(testCase.ValueToFind, testCase.Elements)
			assert.Equal(t, testCase.ExpectedResult, result)
		})
	}
}
