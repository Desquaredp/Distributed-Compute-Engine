package main

import (
	"go.uber.org/zap"
	"src/types"
	"testing"
)

func TestNodeManager_sortKeys(t *testing.T) {
	type fields struct {
		port   string
		logger *zap.Logger
	}
	type args struct {
		output []types.MapOutputRecord
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test sortKeys",
			fields: fields{
				port:   "8080",
				logger: zap.NewNop(),
			},
			args: args{
				output: []types.MapOutputRecord{

					{[]byte("b"), []byte("2")},
					{[]byte("c"), []byte("3")},
					{[]byte("a"), []byte("1")},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nm := &NodeManager{
				NodeManagerPort: tt.fields.port,
				logger:          tt.fields.logger,
			}
			if err := nm.sortKeys(tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("sortKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
