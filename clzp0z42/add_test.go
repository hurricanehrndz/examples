package add

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type inputsType struct {
	a int
	b int
}

var testCases = []struct {
	inputs   inputsType
	expected int
}{
	{inputsType{2, 3}, 5},
	{inputsType{3, 3}, 6},
}

func TestAddTable(t *testing.T) {
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test add case: %d", i), func(t *testing.T) {
			assert.Equal(t, tc.expected, add(tc.inputs.a, tc.inputs.b))
		})
	}
}
