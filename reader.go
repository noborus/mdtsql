package mdtsql

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/noborus/trdsql"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	gast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// Caption is a caption specification.
var Caption bool

// MDTReader is a reader for markdown table.
type MDTReader struct {
	tableNames []string
	caption    bool
	names      []string
	types      []string
	tables     []ast.Node
	body       [][]interface{}
	source     []byte
}

// NewMDTReader returns a new MDTReader.
func NewMDTReader(reader io.Reader, opts *trdsql.ReadOpts) (trdsql.Reader, error) {
	r := MDTReader{}
	r.caption = Caption
	target := 0
	capTitle := ""
	if r.caption {
		capTitle = opts.InJQuery
	} else {
		target = targetTable(opts.InJQuery)
	}
	if err := r.parse(reader); err != nil {
		return nil, err
	}

	for i, node := range r.tables {
		if r.caption {
			if r.tableNames[i] != capTitle {
				continue
			}
		} else {
			if i != target {
				continue
			}
		}
		table, err := tableNode(r.source, node)
		if err != nil {
			return nil, err
		}
		r.names = table.names
		r.types = table.types
		r.body = table.body
	}

	return &r, nil
}

func targetTable(optString string) int {
	target := 0
	if optString != "" {
		n, err := strconv.Atoi(optString)
		if err == nil {
			target = n
		}
	}
	return target
}

// parse reads the content from the given io.Reader and parses it using the goldmark library.
// It populates the MDTReader's names and source fields based on the parsed content.
// The parsed content is then passed to the parseNode method for further processing.
// Returns an error if there was an issue reading or parsing the content.
func (r *MDTReader) parse(reader io.Reader) error {
	gmd := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	source, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	parser := gmd.Parser()
	node := parser.Parse(text.NewReader(source))

	for n := node.FirstChild(); n != nil; n = n.NextSibling() {
		switch n.Kind() {
		case gast.KindTableHeader:
			i := 0
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				col := toText(c.Text(source))
				if col == "" {
					col = fmt.Sprintf("c%d", i+1)
				}
				r.names = append(r.names, col)
				i++
			}
		case gast.KindTableRow:
			row := []string{}
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				rawText := []byte{}
				for i := 0; i < c.Lines().Len(); i++ {
					line := c.Lines().At(i)
					rawText = append(rawText, line.Value(source)...)
				}
				row = append(row, string(rawText))
			}
			data := make([]interface{}, len(row))
			for i, col := range row {
				data[i] = col
			}
		default:
			// fmt.Fprintf(os.Stderr, "unknown node:")
			// fmt.Fprintf(os.Stderr, "%v:%v\n", n.Kind(), n.Type())
		}
	}
	r.source = source
	return r.parseNode(node)
}

func (r *MDTReader) parseNode(node ast.Node) error {
	switch node.Type() {
	case ast.TypeDocument:
		for n := node.FirstChild(); n != nil; n = n.NextSibling() {
			if err := r.parseNode(n); err != nil {
				return err
			}
		}
	case ast.TypeBlock:
		if node.Kind() == gast.KindTable {
			r.tables = append(r.tables, node)
		}

		switch node.Kind() {
		case ast.KindHeading, ast.KindParagraph:
			if r.caption {
				caption := string(node.Text(r.source))
				if r.existsTableName(caption) {
					caption = incrementName(caption)
				}
				r.tableNames = append(r.tableNames, caption)
			}
		}
	default:
		fmt.Printf("unknown node %v:%v", node.Kind(), node.Type())
	}

	return nil
}

func (r *MDTReader) existsTableName(name string) bool {
	for _, n := range r.tableNames {
		if n == name {
			return true
		}
	}
	return false
}

func incrementName(name string) string {
	names := strings.Split(name, "_")
	if len(names) == 1 {
		return fmt.Sprintf("%s_1", name)
	}
	n, err := strconv.Atoi(names[len(names)-1])
	if err != nil {
		return fmt.Sprintf("%s_1", name)
	}
	n++
	return fmt.Sprintf("%s_%d", name, n)
}

func tableNode(source []byte, node ast.Node) (table, error) {
	t := table{}
	for n := node.FirstChild(); n != nil; n = n.NextSibling() {
		switch n.Kind() {
		case gast.KindTableHeader:
			i := 0
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				col := toText(c.Text(source))
				if col == "" {
					col = fmt.Sprintf("c%d", i+1)
				}
				t.names = append(t.names, col)
				i++
			}
		case gast.KindTableRow:
			row := []string{}
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				rawText := []byte{}
				for i := 0; i < c.Lines().Len(); i++ {
					line := c.Lines().At(i)
					rawText = append(rawText, line.Value(source)...)
				}
				row = append(row, string(rawText))
			}
			data := make([]interface{}, len(row))
			for i, col := range row {
				data[i] = col
			}
			t.body = append(t.body, data)
		default:
			return t, fmt.Errorf("unknown node %v:%v", n.Kind(), n.Type())
		}
	}
	t.types = make([]string, len(t.names))
	for i := 0; i < len(t.names); i++ {
		t.types[i] = trdsql.DefaultDBType
	}
	return t, nil
}

func (t MDTReader) Names() ([]string, error) {
	return t.names, nil
}

func (t MDTReader) Types() ([]string, error) {
	return t.types, nil
}

func (t MDTReader) PreReadRow() [][]interface{} {
	return t.body
}

// ReadRow only returns EOF.
func (t MDTReader) ReadRow(row []interface{}) ([]interface{}, error) {
	return nil, io.EOF
}

func toText(buf []byte) string {
	if len(buf) > 0 {
		return string(buf)
	}
	return ""
}

func init() {
	trdsql.RegisterReaderFunc("MD", NewMDTReader)
}
