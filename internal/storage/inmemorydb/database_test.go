package inmemorydb

import (
	"context"
	"ethereum-tracker-app/cmd/config"
	"fmt"
	"log"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestReturnByValue(t *testing.T) {
	// Test cases for int type
	int1, int2, int3 := 1, 2, 3
	intCases := []struct {
		input    []*int
		expected []int
	}{
		{
			input:    []*int{&int1, &int2, &int3},
			expected: []int{int1, int2, int3},
		},
		{
			input:    []*int{nil, &int1, nil, &int2},
			expected: []int{0, int1, 0, int2}, // 0 is the zero/default value of int
		},
	}

	for _, testCase := range intCases {
		t.Run("int", func(t *testing.T) {
			actual := returnByValue(testCase.input)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("input: %v, expected: %v, got: %v", testCase.input, testCase.expected, actual)
			}
		})
	}

	// Test cases for string type
	string1, string2, string3 := "a", "b", "c"
	stringCases := []struct {
		input    []*string
		expected []string
	}{
		{
			input:    []*string{&string1, &string2, &string3},
			expected: []string{string1, string2, string3},
		},
		{
			input:    []*string{nil, &string1, nil, &string2},
			expected: []string{"", string1, "", string2}, // "" is the zero/default value of string
		},
	}

	for _, testCase := range stringCases {
		t.Run("string", func(t *testing.T) {
			actual := returnByValue(testCase.input)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("input: %v, expected: %v, got: %v", testCase.input, testCase.expected, actual)
			}
		})
	}

	// Test cases for bool type
	bool1, bool2 := true, false
	boolCases := []struct {
		input    []*bool
		expected []bool
	}{
		{
			input:    []*bool{&bool1, &bool2, nil},
			expected: []bool{bool1, bool2, false}, // false is the zero/default value of bool
		},
		{
			input:    []*bool{nil, &bool1},
			expected: []bool{false, bool1},
		},
	}

	for _, testCase := range boolCases {
		t.Run("bool", func(t *testing.T) {
			actual := returnByValue(testCase.input)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("input: %v, expected: %v, got: %v", testCase.input, testCase.expected, actual)
			}
		})
	}
}

func TestSetBlock(t *testing.T) {
	type testCase struct {
		name        string
		blockNumber *big.Int
		isInserted  bool
	}

	testcases := []testCase{
		{
			name:        "insert corret block number",
			blockNumber: big.NewInt(123),
			isInserted:  true,
		},
		{
			name:        "insert another corret block number",
			blockNumber: big.NewInt(23),
			isInserted:  true,
		},
	}

	db := NewInmemortDBService(config.Config{}, log.New(os.Stdout, "app", log.LstdFlags)).(*inmemoryDB)
	for _, tt := range testcases {
		block := types.NewBlockWithHeader(&types.Header{Number: tt.blockNumber})
		db.SetBlock(context.Background(), block)
		db.mu.RLock()
		_, ok := db.blocks[block.NumberU64()]
		db.mu.RUnlock()
		fmt.Println(tt.isInserted, ok)
		assert.Equal(t, tt.isInserted, ok)
	}
}
