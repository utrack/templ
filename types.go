package templ

import (
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/a-h/lexical/parse"
)

// {% package templ %}
//
// {% import "strings" %}
// {% import strs "strings" %}
//
// {% templ RenderAddress(addr Address) %}
// 	<div>{%= addr.Address1 %}</div>
// 	<div>{%= addr.Address2 %}</div>
// 	<div>{%= addr.Address3 %}</div>
// 	<div>{%= addr.Address4 %}</div>
// {% endtempl %}
//
// {% templ Render(p Person) %}
//    <div>
//      <div>{%= p.Name() %}</div>
//      <a href={%= p.URL %}>{%= strings.ToUpper(p.Name()) %}</a>
//      <div>
//          {% if p.Type == "test" %}
//             <span>{%= "Test user" %}</span>
//          {% else %}
// 	    <span>{%= "Not test user" %}</span>
//          {% endif %}
//          {% for _, v := range p.Addresses %}
//             {% call RenderAddress(v) %}
//          {% endfor %}
//      </div>
//    </div>
// {% endtempl %}

// Source mapping to map from the source code of the template to the
// in-memory representation.
type Position struct {
	Index int64
	Line  int
	Col   int
}

func (p Position) String() string {
	return fmt.Sprintf("line %d, col %d (index %d)", p.Line, p.Col, p.Index)
}

// NewPosition initialises a position.
func NewPosition() Position {
	return Position{
		Index: 0,
		Line:  1,
		Col:   0,
	}
}

// NewPositionFromValues initialises a position.
func NewPositionFromValues(index int64, line, col int) Position {
	return Position{
		Index: index,
		Line:  line,
		Col:   col,
	}
}

// NewPositionFromInput creates a position from a parse input.
func NewPositionFromInput(pi parse.Input) Position {
	l, c := pi.Position()
	return Position{
		Index: pi.Index(),
		Line:  l,
		Col:   c,
	}
}

// NewExpression creates a Go expression.
func NewExpression(value string, from, to Position) Expression {
	return Expression{
		Value: value,
		Range: Range{
			From: from,
			To:   to,
		},
	}
}

// NewRange creates a range.
func NewRange(from, to Position) Range {
	return Range{
		From: from,
		To:   to,
	}
}

// Range of text within a file.
type Range struct {
	From Position
	To   Position
}

// Expression containing Go code.
type Expression struct {
	Value string
	Range Range
}

type TemplateFile struct {
	Package   Package
	Imports   []Import
	Templates []Template
}

func (tf TemplateFile) Write(w io.Writer) error {
	var indent int
	if err := tf.Package.Write(w, indent); err != nil {
		return err
	}
	if _, err := w.Write([]byte("\n\n")); err != nil {
		return err
	}
	for i := 0; i < len(tf.Imports); i++ {
		if err := tf.Imports[i].Write(w, indent); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}
	for i := 0; i < len(tf.Templates); i++ {
		if err := tf.Templates[i].Write(w, indent); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n\n")); err != nil {
			return err
		}
	}
	return nil
}

// {% package templ %}
type Package struct {
	Expression Expression
}

func (p Package) Write(w io.Writer, indent int) error {
	return writeIndent(w, indent, "{% package "+p.Expression.Value+" %}")
}

func writeIndent(w io.Writer, level int, s string) (err error) {
	if _, err = w.Write([]byte(strings.Repeat("\t", level))); err != nil {
		return
	}
	_, err = w.Write([]byte(s))
	return
}

// Whitespace.
type Whitespace struct {
	Value string
}

func (ws Whitespace) IsNode() bool { return true }
func (ws Whitespace) Write(w io.Writer, indent int) error {
	// Explicit whitespace nodes are removed from templates because they're auto-formatted.
	return nil
}

// {% import "strings" %}
// {% import strs "strings" %}
type Import struct {
	Expression Expression
}

func (imp Import) Write(w io.Writer, indent int) error {
	return writeIndent(w, indent, "{% import "+imp.Expression.Value+" %}\n")
}

// Template definition.
// {% templ Name(p Parameter) %}
//   {% if ... %}
//   <Element></Element>
// {% endtempl %}
type Template struct {
	Name       Expression
	Parameters Expression
	Children   []Node
}

func (t Template) Write(w io.Writer, indent int) error {
	if err := writeIndent(w, indent, "{% templ "+t.Name.Value+"("+t.Parameters.Value+") %}\n"); err != nil {
		return err
	}
	if err := writeNodes(w, indent+1, t.Children); err != nil {
		return err
	}
	if err := writeIndent(w, indent, "{% endtempl %}\n\n"); err != nil {
		return err
	}
	return nil
}

// A Node appears within a template, e.g. a StringExpression, Element, IfExpression etc.
type Node interface {
	IsNode() bool
	// Write out the string.
	Write(w io.Writer, indent int) error
}

// <a .../> or <div ...>...</div>
type Element struct {
	Name       string
	Attributes []Attribute
	Children   []Node
}

var blockElements = map[string]struct{}{
	"br":  {},
	"hr":  {},
	"div": {},
}

func (e Element) isBlockElement() bool {
	_, ok := blockElements[e.Name]
	return ok
}

func (e Element) IsNode() bool { return true }
func (e Element) Write(w io.Writer, indent int) error {
	// If it's a block element, start a newline and indent.
	if e.isBlockElement() {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
		if err := writeIndent(w, indent, "<"+e.Name); err != nil {
			return err
		}
	} else {
		if err := writeIndent(w, 0, "<"+e.Name); err != nil {
			return err
		}
	}
	for i := 0; i < len(e.Attributes); i++ {
		if _, err := w.Write([]byte(" ")); err != nil {
			return err
		}
		a := e.Attributes[i]
		if _, err := w.Write([]byte(a.String())); err != nil {
			return err
		}
	}
	// Doesn't have children, close up and exit.
	if len(e.Children) == 0 {
		if err := writeIndent(w, indent, "/>"); err != nil {
			return err
		}
		if e.isBlockElement() {
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		}
		return nil
	}
	// Has children.
	if e.isBlockElement() {
		if _, err := w.Write([]byte(">\n")); err != nil {
			return err
		}
		indent++
	} else {
		if err := writeIndent(w, indent, ">"); err != nil {
			return err
		}
		indent = 0
	}
	if err := writeNodes(w, indent, e.Children); err != nil {
		return err
	}
	if e.isBlockElement() {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
		indent--
	}
	return writeIndent(w, indent, "</"+e.Name+">")
}

func writeNodes(w io.Writer, indent int, nodes []Node) error {
	for i := 0; i < len(nodes); i++ {
		c := nodes[i]
		//TODO: Add newlines depending on what type of node it is.
		if err := c.Write(w, indent); err != nil {
			return err
		}
	}
	return nil
}

type Attribute interface {
	IsAttribute() bool
	String() string
}

// href=""
type ConstantAttribute struct {
	Name  string
	Value string
}

func (ca ConstantAttribute) IsAttribute() bool { return true }
func (ca ConstantAttribute) String() string {
	return ca.Name + `"` + html.EscapeString(ca.Value) + `"`
}

// href={%= ... }
type ExpressionAttribute struct {
	Name  string
	Value StringExpression
}

func (ea ExpressionAttribute) IsAttribute() bool { return true }
func (ea ExpressionAttribute) String() string {
	return ea.Name + `{%= ` + ea.Value.Expression.Value + ` %}`
}

// Nodes.

// {%= ... %}
type StringExpression struct {
	Expression Expression
}

func (se StringExpression) IsNode() bool { return true }
func (se StringExpression) Write(w io.Writer, indent int) error {
	return writeIndent(w, indent, `{%= `+se.Expression.Value+` %}`)
}

// {% call Other(p.First, p.Last) %}
type CallTemplateExpression struct {
	// Name of the template to execute.
	Name Expression
	// Arguments to pass.
	Arguments Expression
}

func (cte CallTemplateExpression) IsNode() bool { return true }
func (cte CallTemplateExpression) Write(w io.Writer, indent int) error {
	return writeIndent(w, indent, `{% call `+cte.Name.Value+`(`+cte.Arguments.Value+`) %}`)
}

// {% if p.Type == "test" && p.thing %}
// {% endif %}
type IfExpression struct {
	Expression Expression
	Then       []Node
	Else       []Node
}

func (n IfExpression) IsNode() bool { return true }
func (n IfExpression) Write(w io.Writer, indent int) error {
	if err := writeIndent(w, indent, "{% if "+n.Expression.Value+" %}\n"); err != nil {
		return err
	}
	indent++
	if err := writeNodes(w, indent, n.Then); err != nil {
		return err
	}
	indent--
	if len(n.Else) > 0 {
		if err := writeIndent(w, indent, "{% else %}\n"); err != nil {
			return err
		}
		indent++
		if err := writeNodes(w, indent, n.Else); err != nil {
			return err
		}
		indent--
	}
	if err := writeIndent(w, indent, "{% endif %}\n"); err != nil {
		return err
	}
	return nil
}

// {% switch p.Type %}
//  {% case "Something" %}
//  {% endcase %}
// {% endswitch %}
type SwitchExpression struct {
	Expression Expression
	Cases      []CaseExpression
	Default    []Node
}

func (se SwitchExpression) IsNode() bool { return true }
func (se SwitchExpression) Write(w io.Writer, indent int) error {
	if err := writeIndent(w, indent, "{% switch "+se.Expression.Value+" %}\n"); err != nil {
		return err
	}
	indent++
	for i := 0; i < len(se.Cases); i++ {
		c := se.Cases[i]
		if err := writeIndent(w, indent, "{% case "+c.Expression.Value+" %}\n"); err != nil {
			return err
		}
		indent++
		if err := writeNodes(w, indent, c.Children); err != nil {
			return err
		}
		indent--
		if err := writeIndent(w, indent, "{% endcase %}\n"); err != nil {
			return err
		}
	}
	if len(se.Default) > 0 {
		if err := writeIndent(w, indent, "{% default %}\n"); err != nil {
			return err
		}
		indent++
		if err := writeNodes(w, indent, se.Default); err != nil {
			return err
		}
		indent--
		if err := writeIndent(w, indent, "{% enddefault %}\n"); err != nil {
			return err
		}
	}
	indent--
	if err := writeIndent(w, indent, "{% endswitch %}\n"); err != nil {
		return err
	}
	return nil
}

// {% case "Something" %}
// ...
// {% endcase %}
type CaseExpression struct {
	Expression Expression
	Children   []Node
}

// {% for i, v := range p.Addresses %}
//   {% call Address(v) %}
// {% endfor %}
type ForExpression struct {
	Expression Expression
	Children   []Node
}

func (fe ForExpression) IsNode() bool { return true }
func (fe ForExpression) Write(w io.Writer, indent int) error {
	if err := writeIndent(w, indent, "{% for "+fe.Expression.Value+" %}\n"); err != nil {
		return err
	}
	indent++
	if err := writeNodes(w, indent, fe.Children); err != nil {
		return err
	}
	indent--
	if err := writeIndent(w, indent, "{% endfor %}\n"); err != nil {
		return err
	}
	return nil
}
