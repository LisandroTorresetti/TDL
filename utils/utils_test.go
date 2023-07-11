package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
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
