package binary_test

import (
	"github.com/Sergii-Kirichok/pr/pkg/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBits(t *testing.T) {
	const (
		n = 6
		r = binary.Bits(0b0100_0000)
	)

	assert.Equal(t, r, binary.GetBits(n))
}

func TestSetBits(t *testing.T) {
	const (
		b = binary.Bits(0b10101011)
		f = binary.Bits(0b01110001)
		r = binary.Bits(0b11111011)
	)

	assert.Equal(t, r, binary.SetBits(b, f))
}

func TestClearBits(t *testing.T) {
	const (
		b = binary.Bits(0b10101011)
		f = binary.Bits(0b01110010)
		r = binary.Bits(0b10001001)
	)

	assert.Equal(t, r, binary.ClearBits(b, f))
}

func TestToggleBits(t *testing.T) {
	const (
		b = binary.Bits(0b10101011)
		f = binary.Bits(0b01011000)
		r = binary.Bits(0b11110011)
	)

	assert.Equal(t, r, binary.ToggleBits(b, f))
}

func TestHasBits(t *testing.T) {
	testCases := []struct {
		name     string
		b, f     binary.Bits
		expected bool
	}{
		{
			name:     "Yes",
			b:        binary.Bits(0b10101011),
			f:        binary.Bits(0b00101000),
			expected: true,
		},
		{
			name:     "No",
			b:        binary.Bits(0b10101011),
			f:        binary.Bits(0b00011000),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, binary.HasBits(tc.b, tc.f))
		})
	}
}

func TestClearFlag(t *testing.T) {
	type input struct {
		bitPos  uint
		numBits uint
	}

	testCases := []struct {
		name   string
		input  input
		output binary.Bits
	}{
		{
			name:   "First",
			input:  input{bitPos: 2, numBits: 0},
			output: 0b0000_0000,
		},
		{
			name:   "Second",
			input:  input{bitPos: 3, numBits: 2},
			output: 0b0001_1000,
		},
		{
			name:   "Third",
			input:  input{bitPos: 6, numBits: 3},
			output: 0b0001_1100_0000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := binary.ClearFlag(tc.input.bitPos, tc.input.numBits)
			assert.Equal(t, tc.output, actual)
		})
	}
}

func TestSetFromTo(t *testing.T) {
	b := binary.Bits(0b1000_0000)
	flag := binary.Bits(0b01)

	expected := binary.Bits(0b0100_0000)
	actual := binary.SetFromTo(6, 7, b, flag)
	//fmt.Printf("%08b\n", actual)
	assert.Equal(t, expected, actual)
}
