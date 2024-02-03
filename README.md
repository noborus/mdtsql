# mdtsql

A CLI tool that executes SQL queries and converts the results into a Markdown table.

## install

### Go install

```console
go install github.com/noborus/mdtsql/cmd/mdtsql@latest
```

### Homebrew

```console
brew install noborus/tap/mdtsql
```

## Usage

Executes SQL for markdown containing table.
The result can be output to CSV, JSON, LTSV, YAML, Markdown, etc.

```sh
mdtsql query "SELECT * FROM file.md"
```

```sh
mdtsql table file.md
```

### option

```console
mdtsql --help
Execute SQL for table in markdown.
The result can be output to CSV, JSON, LTSV, YAML, Markdown, etc.

Usage:
  mdtsql [flags]
  mdtsql [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List and analyze SQL dumps
  query       Execute SQL queries on markdown table and tabular data
  table       SQL(SELECT * FROM table) for markdown table and tabular data

Flags:
  -d, --Delimiter string   output delimiter (CSV only) (default ",")
  -O, --Header             output header (CSV only)
  -o, --OutFormat string   output format=at|csv|ltsv|json|jsonl|tbln|raw|md|vf|yaml (default "md")
  -c, --caption            caption table name
      --config string      config file (default is $HOME/.mdtsql.yaml)
      --debug              debug print
  -h, --help               help for mdtsql
  -v, --version            display version information

Use "mdtsql [command] --help" for more information about a command.
```

### Example

```sh
mdtsql query "SELECT * FROM file.md"
```

| c1 | a  | b  | c  |
|----|----|----|----|
|  1 | a1 | b1 | c1 |
|  2 | a2 | b2 | c2 |
|  3 | a3 | b3 | c3 |

If the markdown includes multiple tables,
the second and subsequent tables are marked with `::number`.

```sh
mdtsql query "SELECT * FROM file.md::1"
```

Specify the output format with option -o.
-o csv, -o ltsv, -ojson ...

```sh
mdtsql -o csv query "SELECT * FROM file.md"
```

```CSV
1,a1,b1,c1
2,a2,b2,c2
3,a3,b3,c3
```

### List Command

The `list` command displays all the tables in the specified markdown file.

```sh
mdtsql list file.md
```

```sh
mdtsql list abc.md
Table Name: [0]
+-------------+------+
| column name | type |
+-------------+------+
| c1          | text |
| a           | text |
| b           | text |
| c           | text |
+-------------+------+

Table Name: [1]
+-------------+------+
| column name | type |
+-------------+------+
| c1          | text |
| a           | text |
| b           | text |
| c           | text |
+-------------+------+

Table Name: [2]
+-------------+------+
| column name | type |
+-------------+------+
| c1          | text |
| a           | text |
| b           | text |
| c           | text |
+-------------+------+
```

## Table Command

The `table` command executes SQL(SELECT * FROM table) for markdown table and tabular data.

```sh
mdtsql table file.md
```

```sh
mdtsql table file.md::1
```

## Caption option

The  `--caption` or `-c` option specifies a caption name, not a sequential number.
This allows you to specify the same table even if the order changes.

```sh
mdtsql --caption list testdata/abc.md
Table Name: [header]
+-------------+------+
| column name | type |
+-------------+------+
| c1          | text |
| a           | text |
| b           | text |
| c           | text |
+-------------+------+

Table Name: [caption]
+-------------+------+
| column name | type |
+-------------+------+
| c1          | text |
| a           | text |
| b           | text |
| c           | text |
+-------------+------+

Table Name: [caption_1]
+-------------+------+
| column name | type |
+-------------+------+
| c1          | text |
| a           | text |
| b           | text |
| c           | text |
+-------------+------+
```

```sh
mdtsql --caption query "SELECT * FROM testdata/abc.md::caption_1"
| c1 | a  | b  | c  |
|----|----|----|----|
|  1 | a1 | b1 | c1 |
|  2 | a2 | b2 | c2 |
|  3 | a3 | b3 | c3 |
```

## Multiple queries

You can specify multiple queries with the `;` separator.

```console
mdtsql query "INSERT INTO abc.md::2 (c1, a, b, c) VALUES ('4', 'a4', 'b4', 'c4');SELECT * FROM abc.md::2"
```

| c1 | a  | b  | c  |
|----|----|----|----|
|  1 | a1 | b1 | c1 |
|  2 | a2 | b2 | c2 |
|  3 | a3 | b3 | c3 |
| **4** | **a4** | **b4** | **c4** |

```console
mdtsql query "UPDATE abc.md::2 SET c='u4' WHERE c1=3;SELECT * FROM abc.md::2"
```

| c1 | a  | b  | c  |
|----|----|----|----|
|  1 | a1 | b1 | c1 |
|  2 | a2 | b2 | c2 |
|  3 | a3 | b3 | **u4** |
