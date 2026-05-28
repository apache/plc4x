/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package xmldom provides a tiny DOM tree over encoding/xml for the test harness.
// It replaces the subset of github.com/subchen/go-xmldom that plc4go's testutils used.
package xmldom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Attribute struct {
	Name  string
	Value string
}

type Node struct {
	Document   *Document
	Parent     *Node
	Name       string
	Attributes []*Attribute
	Children   []*Node
	Text       string
}

type Document struct {
	ProcInst   string
	Directives []string
	Root       *Node
}

// Must panics if err is non-nil; otherwise returns doc.
func Must(doc *Document, err error) *Document {
	if err != nil {
		panic(err)
	}
	return doc
}

// Parse reads an XML document from r and returns a DOM tree.
func Parse(r io.Reader) (*Document, error) {
	dec := xml.NewDecoder(r)
	doc := &Document{}
	var current *Node
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			el := &Node{Document: doc, Parent: current, Name: t.Name.Local}
			for _, a := range t.Attr {
				el.Attributes = append(el.Attributes, &Attribute{Name: a.Name.Local, Value: a.Value})
			}
			if current != nil {
				current.Children = append(current.Children, el)
			}
			current = el
			if doc.Root == nil {
				doc.Root = el
			}
		case xml.EndElement:
			if current != nil {
				current = current.Parent
			}
		case xml.CharData:
			if current != nil {
				current.Text = string(bytes.TrimSpace(t))
			}
		case xml.ProcInst:
			doc.ProcInst = fmt.Sprintf("<?%s %s?>", t.Target, string(t.Inst))
		case xml.Directive:
			doc.Directives = append(doc.Directives, fmt.Sprintf("<!%s>", string(t)))
		}
	}
	return doc, nil
}

func (n *Node) GetAttribute(name string) *Attribute {
	for _, a := range n.Attributes {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func (n *Node) GetAttributeValue(name string) string {
	if a := n.GetAttribute(name); a != nil {
		return a.Value
	}
	return ""
}

func (n *Node) GetChild(name string) *Node {
	for _, c := range n.Children {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (n *Node) GetChildren(name string) []*Node {
	var out []*Node
	for _, c := range n.Children {
		if c.Name == name {
			out = append(out, c)
		}
	}
	return out
}

func (n *Node) FirstChild() *Node {
	if len(n.Children) > 0 {
		return n.Children[0]
	}
	return nil
}

// FindOneByName performs a depth-first search and returns the first node whose name matches.
func (n *Node) FindOneByName(name string) *Node {
	if n.Name == name {
		return n
	}
	for _, c := range n.Children {
		if x := c.FindOneByName(name); x != nil {
			return x
		}
	}
	return nil
}

// FindByName performs a depth-first search and returns every node whose name matches.
func (n *Node) FindByName(name string) []*Node {
	var out []*Node
	if n.Name == name {
		out = append(out, n)
	}
	for _, c := range n.Children {
		out = append(out, c.FindByName(name)...)
	}
	return out
}

// XMLPretty renders the node tree with two-space indentation, matching the format produced by
// the previous xmldom library so golden-file test fixtures keep comparing equal.
func (n *Node) XMLPretty() string {
	buf := new(bytes.Buffer)
	printXML(buf, n, 0, "  ")
	return buf.String()
}

func printXML(buf *bytes.Buffer, n *Node, level int, indent string) {
	pretty := len(indent) > 0
	if pretty {
		buf.WriteString(strings.Repeat(indent, level))
	}
	buf.WriteByte('<')
	buf.WriteString(n.Name)
	for _, a := range n.Attributes {
		buf.WriteByte(' ')
		buf.WriteString(a.Name)
		buf.WriteString(`="`)
		xml.Escape(buf, []byte(a.Value))
		buf.WriteByte('"')
	}
	if len(n.Children) == 0 && len(n.Text) == 0 {
		buf.WriteString(" />")
		if pretty {
			buf.WriteByte('\n')
		}
		return
	}
	buf.WriteByte('>')
	if len(n.Children) > 0 {
		if pretty {
			buf.WriteByte('\n')
		}
		for _, c := range n.Children {
			printXML(buf, c, level+1, indent)
		}
	}
	if len(n.Text) > 0 {
		xml.EscapeText(buf, []byte(n.Text))
	}
	if len(n.Children) > 0 && pretty {
		buf.WriteString(strings.Repeat(indent, level))
	}
	buf.WriteString("</")
	buf.WriteString(n.Name)
	buf.WriteByte('>')
	if pretty {
		buf.WriteByte('\n')
	}
}
