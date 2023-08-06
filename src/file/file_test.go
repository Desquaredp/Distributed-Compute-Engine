package file

import (
	"reflect"
	"testing"
)

func TestFileHandler_ChecksumOnDisk(t *testing.T) {
	type fields struct {
		fileName    string
		fileSize    int64
		dataStream  []byte
		checksum    [16]byte
		location    []string
		fragments   []*Fragment
		fragmentMap map[string]*Fragment
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{
			name: "Test ChecksumOnDisk",
			fields: fields{
				fileName:    "file1",
				fileSize:    100,
				dataStream:  []byte("data"),
				checksum:    [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				location:    []string{"localhost:8080"},
				fragments:   []*Fragment{},
				fragmentMap: map[string]*Fragment{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileHandler{
				fileName:    tt.fields.fileName,
				fileSize:    tt.fields.fileSize,
				dataStream:  tt.fields.dataStream,
				checksum:    tt.fields.checksum,
				location:    tt.fields.location,
				fragments:   tt.fields.fragments,
				fragmentMap: tt.fields.fragmentMap,
			}
			f.ChecksumOnDisk()
		})
	}
}

func TestFileHandler_ValidateChecksumFromFile(t *testing.T) {
	type fields struct {
		fileName    string
		fileSize    int64
		dataStream  []byte
		checksum    [16]byte
		location    []string
		fragments   []*Fragment
		fragmentMap map[string]*Fragment
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValid bool
	}{
		// TODO: Add test cases.
		{
			name: "Test ValidateChecksumFromFile",
			fields: fields{
				fileName:    "file1",
				fileSize:    100,
				dataStream:  []byte("data"),
				checksum:    [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				location:    []string{"localhost:8080"},
				fragments:   []*Fragment{},
				fragmentMap: map[string]*Fragment{},
			},
			args: args{
				fileName: "file1",
			},
			wantValid: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileHandler{
				fileName:    tt.fields.fileName,
				fileSize:    tt.fields.fileSize,
				dataStream:  tt.fields.dataStream,
				checksum:    tt.fields.checksum,
				location:    tt.fields.location,
				fragments:   tt.fields.fragments,
				fragmentMap: tt.fields.fragmentMap,
			}
			if gotValid := f.ValidateChecksumFromFile(tt.args.fileName); gotValid != tt.wantValid {
				t.Errorf("ValidateChecksumFromFile() = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}

func TestFileHandler_ReadCheckSum(t *testing.T) {
	type fields struct {
		fileName    string
		fileSize    int64
		dataStream  []byte
		checksum    [16]byte
		location    []string
		fragments   []*Fragment
		fragmentMap map[string]*Fragment
	}
	tests := []struct {
		name      string
		fields    fields
		wantFiles []string
	}{
		// TODO: Add test cases.

		{
			name: "Test ReadCheckSum",
			fields: fields{
				fileName:   "file1",
				fileSize:   100,
				dataStream: []byte("data"),
				checksum:   [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			},
			wantFiles: []string{"file1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileHandler{
				fileName:    tt.fields.fileName,
				fileSize:    tt.fields.fileSize,
				dataStream:  tt.fields.dataStream,
				checksum:    tt.fields.checksum,
				location:    tt.fields.location,
				fragments:   tt.fields.fragments,
				fragmentMap: tt.fields.fragmentMap,
			}
			if gotFiles, _ := f.ReadCheckSum(); !reflect.DeepEqual(gotFiles, tt.wantFiles) {
				t.Errorf("ReadCheckSum() = %v, want %v", gotFiles, tt.wantFiles)
			}
		})
	}
}

func TestFileHandler_ReadFile(t *testing.T) {
	type fields struct {
		fileName    string
		fileSize    int64
		dataStream  []byte
		checksum    [16]byte
		location    []string
		fragments   []*Fragment
		fragmentMap map[string]*Fragment
	}
	tests := []struct {
		name     string
		fields   fields
		wantFile []byte
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "Test ReadFile",
			fields: fields{
				fileName: "file1",
			},

			wantFile: nil,
			wantErr:  true,
		},

		{
			name: "Test ReadFile",
			fields: fields{
				fileName: "config.yaml",
			},
			wantFile: []byte("wabba labba dub dub"),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileHandler{
				fileName:    tt.fields.fileName,
				fileSize:    tt.fields.fileSize,
				dataStream:  tt.fields.dataStream,
				checksum:    tt.fields.checksum,
				location:    tt.fields.location,
				fragments:   tt.fields.fragments,
				fragmentMap: tt.fields.fragmentMap,
			}
			gotFile, err := f.ReadFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFile, tt.wantFile) {
				t.Errorf("ReadFile() gotFile = %v, want %v", gotFile, tt.wantFile)
			}
		})
	}
}

func TestFileHandler_ChecksumOnDisk1(t *testing.T) {
	type fields struct {
		fileName    string
		fileSize    int64
		dataStream  []byte
		checksum    [16]byte
		location    []string
		fragments   []*Fragment
		fragmentMap map[string]*Fragment
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.

		{
			name: "Test ChecksumOnDisk",
			fields: fields{
				fileName: "file1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileHandler{
				fileName:    tt.fields.fileName,
				fileSize:    tt.fields.fileSize,
				dataStream:  tt.fields.dataStream,
				checksum:    tt.fields.checksum,
				location:    tt.fields.location,
				fragments:   tt.fields.fragments,
				fragmentMap: tt.fields.fragmentMap,
			}
			if err := f.ChecksumOnDisk(); (err != nil) != tt.wantErr {
				t.Errorf("ChecksumOnDisk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileHandler_FillFragmentData(t *testing.T) {
	type fields struct {
		fileName    string
		fileSize    int64
		dataStream  []byte
		checksum    [16]byte
		location    []string
		fragments   []*Fragment
		fragmentMap map[string]*Fragment
	}
	type args struct {
		fragId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.

		{
			name: "Test FillFragmentData",
			fields: fields{
				fileName: "file1",
			},
			args: args{
				fragId: "file1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileHandler{
				fileName:    tt.fields.fileName,
				fileSize:    tt.fields.fileSize,
				dataStream:  tt.fields.dataStream,
				checksum:    tt.fields.checksum,
				location:    tt.fields.location,
				fragments:   tt.fields.fragments,
				fragmentMap: tt.fields.fragmentMap,
			}
			f.FillFragmentData(tt.args.fragId)
		})
	}
}
