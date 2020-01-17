package mdtsql

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/noborus/trdsql"
)

const TestData = "testdata"

func TestMarkdownQuery(t *testing.T) {
	type args struct {
		fileName string
		query    string
		caption  bool
	}
	tests := []struct {
		name      string
		outStream *bytes.Buffer
		errStream *bytes.Buffer
		args      args
		want      string
		wantErr   bool
	}{
		{
			name:      "testErr",
			outStream: new(bytes.Buffer),
			errStream: new(bytes.Buffer),
			args: args{
				fileName: "test",
				query:    "SELECT 1",
			},
			wantErr: true,
		},
		{
			name:      "test1",
			outStream: new(bytes.Buffer),
			errStream: new(bytes.Buffer),
			args: args{
				fileName: filepath.Join(TestData, "test.md"),
				query:    "SELECT 1",
			},
			want:    "1\n",
			wantErr: false,
		},
		{
			name:      "testCSV",
			outStream: new(bytes.Buffer),
			errStream: new(bytes.Buffer),
			args: args{
				fileName: filepath.Join(TestData, "test.md"),
				query:    "SELECT * FROM test",
			},
			want:    "1,a1,b1,c1\n",
			wantErr: false,
		},
		{
			name:      "testABC1",
			outStream: new(bytes.Buffer),
			errStream: new(bytes.Buffer),
			args: args{
				fileName: filepath.Join(TestData, "abc.md"),
				query:    "SELECT * FROM abc",
			},
			want:    "1,a1,b1,c1\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := trdsql.NewWriter(
				trdsql.OutFormat(trdsql.CSV),
				trdsql.OutStream(tt.outStream),
				trdsql.ErrStream(tt.errStream),
			)
			if err := MarkdownQuery(tt.args.fileName, tt.args.query, tt.args.caption, w); (err != nil) != tt.wantErr {
				t.Errorf("MarkdownQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				res := tt.outStream.String()
				if res != tt.want {
					t.Errorf("markdownQuery() result = %v, want %v", res, tt.want)
				}
			}
		})
	}
}
