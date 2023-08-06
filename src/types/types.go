package types // Change this to the actual package name

type MRPluginInterface interface {
	Map(input []byte) []MapOutputRecord
	Reduce(input ReduceInputRecord) ReduceOutputRecord
}

type MapInputRecord struct {
	Key   []byte
	Value []byte
}

type MapOutputRecord struct {
	Key   []byte
	Value []byte
}

type ReduceInputRecord struct {
	Key    []byte
	Values [][]byte
}

type ReduceOutputRecord struct {
	Key   []byte
	Value []byte
}
