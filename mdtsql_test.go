package mdtsql

import (
	"bytes"
	"path/filepath"
	"reflect"
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
		{
			name:      "testLink",
			outStream: new(bytes.Buffer),
			errStream: new(bytes.Buffer),
			args: args{
				fileName: filepath.Join(TestData, "link.md"),
				query:    "SELECT * FROM link",
			},
			want:    "[github](https://github.com/),TD1\n[google](https://google.com/),TD2\n",
			wantErr: false,
		},
		{
			name:      "testDecoration",
			outStream: new(bytes.Buffer),
			errStream: new(bytes.Buffer),
			args: args{
				fileName: filepath.Join(TestData, "decoration.md"),
				query:    "SELECT * FROM decoration",
			},
			want:    "`code`,TD1\n*emh*,TD2\n**strong**,TD3\n~~del~~,TD4\n",
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

func TestAnalyze(t *testing.T) {
	type args struct {
		fileName string
		caption  bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "testErr",
			args: args{
				fileName: "testErr",
				caption:  false,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test1",
			args: args{
				fileName: filepath.Join(TestData, "test.md"),
				caption:  false,
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				fileName: filepath.Join(TestData, "test.md"),
				caption:  true,
			},
			want:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Analyze(tt.args.fileName, tt.args.caption)
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.tableName, tt.want) {
					t.Errorf("Analyze() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestImporter_Dump(t *testing.T) {
	type fields struct {
		fileName string
		caption  bool
	}
	tests := []struct {
		name   string
		fields fields
		wantW  string
	}{
		{
			name: "test1",
			fields: fields{
				fileName: filepath.Join(TestData, "test.md"),
				caption:  false,
			},
			wantW: `Table Name: [test]
+-------------+------+
| column name | type |
+-------------+------+
| c1          | text |
| a           | text |
| b           | text |
| c           | text |
+-------------+------+

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im, err := Analyze(tt.fields.fileName, tt.fields.caption)
			if err != nil {
				t.Errorf("Importer.Dump() %s", err)
				t.Skip()
			}
			w := &bytes.Buffer{}
			im.Dump(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Importer.Dump() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
