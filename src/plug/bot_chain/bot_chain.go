package main

import (
	"src/types"
	"strings"
)

type MRBotsPluginImpl struct{}

func (mr *MRBotsPluginImpl) Map(input []byte) []types.MapOutputRecord {
	words := strings.Fields(string(input))
	if len(words) < 2 {
		return []types.MapOutputRecord{}
	}
	output := make([]types.MapOutputRecord, 1)

	output = append(output, types.MapOutputRecord{
		Key:   []byte(words[1]),
		Value: []byte(words[0]),
	})
	return output
}

func (mr *MRBotsPluginImpl) Reduce(input types.ReduceInputRecord) types.ReduceOutputRecord {

	values := input.Values

	//make values a vector of strings and convert to bytes
	values_str := make([]string, len(values))
	for i, _ := range values {
		values_str[i] = string(values[i])
	}
	valuesString := strings.Join(values_str, ",")

	output := types.ReduceOutputRecord{
		Key:   input.Key,
		Value: []byte(valuesString),
	}
	return output
}

var MRPlugin MRBotsPluginImpl
