package main

import (
	"fmt"
	"src/types"
	"strings"
)

type MRBotsPluginImpl struct{}

func (mr *MRBotsPluginImpl) Map(input []byte) []types.MapOutputRecord {
	words := strings.Fields(string(input))
	output := make([]types.MapOutputRecord, 1)

	output[0] = types.MapOutputRecord{
		Key:   []byte(words[2]),
		Value: []byte("1"),
	}
	return output
}

func (mr *MRBotsPluginImpl) Reduce(input types.ReduceInputRecord) types.ReduceOutputRecord {
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

var MRPlugin MRBotsPluginImpl
