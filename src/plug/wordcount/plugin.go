package main

import (
	"fmt"
	"src/types"
	"strings"
)

type MRPlugImpl struct{}

func (mr *MRPlugImpl) Map(input []byte) []types.MapOutputRecord {
	words := strings.Fields(string(input))
	output := make([]types.MapOutputRecord, len(words))
	for i, word := range words {
		output[i] = types.MapOutputRecord{
			Key:   []byte(word),
			Value: []byte("1"),
		}
	}
	return output
}

func (mr *MRPlugImpl) Reduce(input types.ReduceInputRecord) types.ReduceOutputRecord {
	count := 0
	for _, _ = range input.Values {
		count++
	}
	output := types.ReduceOutputRecord{
		Key:   input.Key,
		Value: []byte(fmt.Sprintf("%d", count)),
	}
	return output
}

var MRPlugin MRPlugImpl
