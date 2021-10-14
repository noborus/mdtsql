# mdtsql

Execute SQL to markdown table and convert to other format

## install

```console
go install github.com/noborus/mdtsql/cmd/mdtsql@latest
```

## Usage

Executes SQL for markdown containing table.
The table name is the file name without the extension (.md).

```sh
mdtsql [option] [markdown file]
```

### option

```sh
mdtsql -h
Execute SQL for table in markdown.
The result can be output to CSV, JSON, LTSV, Markdwon, etc.

Usage:
  mdtsql [flags]

Flags:
  -d, --Delimiter string   output delimiter (CSV only) (default ",")
  -O, --Header             output header (CSV only)
  -o, --OutFormat string   output format=at|csv|ltsv|json|jsonl|tbln|raw|md|vf (default "md")
  -c, --caption            caption table name
      --config string      config file (default is $HOME/.mdtsql.yaml)
      --debug              debug print
  -h, --help               help for mdtsql
  -q, --query string       SQL query
  -t, --toggle             Help message for toggle
  -v, --version            display version information
```

### Example

```sh
mdtsql -q "SELECT * FROM file" file.md
```

| c1 | a  | b  | c  |
|----|----|----|----|
|  1 | a1 | b1 | c1 |
|  2 | a2 | b2 | c2 |
|  3 | a3 | b3 | c3 |

If the markdown includes multiple tables,
the second and subsequent tables are marked with `_number`.

```sh
mdtsql -q "SELECT * FROM file_2" file.md
```

Specify the output format with option -o.
-o csv, -o ltsv, -ojson ...

```sh
mdtsql -o csv query "SELECT * FROM file" file.md
```

```CSV
1,a1,b1,c1
2,a2,b2,c2
3,a3,b3,c3
```

If there is no `--query` or `-q` option,
analyze the markdown file and output the table information.

```sh
mdsql abc.md
Table Name: [abc]
+-------------+------+
| column name | type |
+-------------+------+
| c1          | text |
| a           | text |
| b           | text |
| c           | text |
+-------------+------+

Table Name: [abc_2]
+-------------+------+
| column name | type |
+-------------+------+
| c1          | text |
| a           | text |
| b           | text |
| c           | text |
+-------------+------+

Table Name: [abc_3]
+-------------+------+
| column name | type |
+-------------+------+
| c1          | text |
| a           | text |
| b           | text |
| c           | text |
+-------------+------+
```
