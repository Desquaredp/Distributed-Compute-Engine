package main

import (
	"fmt"
	"net/url"
	"src/types"
	"strings"
)

type MRDomainPluginImpl struct{}

func (mr *MRDomainPluginImpl) Map(input []byte) []types.MapOutputRecord {
	words := strings.Fields(string(input))
	output := make([]types.MapOutputRecord, 1)

	url, _ := url.Parse(words[3])
	domain := url.Hostname()

	output[0] = types.MapOutputRecord{
		Key:   []byte(domain),
		Value: []byte("1"),
	}
	return output
}

func (mr *MRDomainPluginImpl) Reduce(input types.ReduceInputRecord) types.ReduceOutputRecord {
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

var MRPlugin MRDomainPluginImpl
