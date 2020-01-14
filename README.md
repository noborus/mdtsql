# mdtsql

Execute SQL in Markdown table

## Usage

Executes SQL for markdown containing table.
The table name is the file name without the extension (.md).

```sh
mdtsql query "SELECT * FROM file" file.md
| c1 | a  | b  | c  |
|----|----|----|----|
|  1 | a1 | b1 | c1 |
|  2 | a2 | b2 | c2 |
|  3 | a3 | b3 | c3 |
```

If the markdown includes multiple tables,
the second and subsequent tables are marked with `_number`.

```sh
mdtsql query "SELECT * FROM file_2" file.md
```

Specify the output format with option -o.
-o csv, -o ltsv, -ojson ...

```sh
mdtsql -o csv query "SELECT * FROM file" file.md
1,a1,b1,c1
2,a2,b2,c2
3,a3,b3,c3
```

