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

package xmldom

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleXML = `<?xml version="1.0" encoding="UTF-8"?>
<testsuite byteOrder="LITTLE_ENDIAN">
  <name>example</name>
  <driver-parameters>
    <parameter><name>host</name><value>1.2.3.4</value></parameter>
    <parameter><name>port</name><value>1234</value></parameter>
  </driver-parameters>
  <testcase>
    <name>case1</name>
    <raw>DEADBEEF</raw>
  </testcase>
</testsuite>`

func parseSample(t *testing.T) *Document {
	t.Helper()
	doc, err := Parse(strings.NewReader(sampleXML))
	require.NoError(t, err)
	require.NotNil(t, doc.Root)
	return doc
}

func TestParse_BuildsRoot(t *testing.T) {
	doc := parseSample(t)
	assert.Equal(t, "testsuite", doc.Root.Name)
	assert.Equal(t, "LITTLE_ENDIAN", doc.Root.GetAttributeValue("byteOrder"))
}

func TestParse_PreservesParentLinks(t *testing.T) {
	doc := parseSample(t)
	name := doc.Root.GetChild("name")
	require.NotNil(t, name)
	assert.Equal(t, doc.Root, name.Parent)
}

func TestParse_ProcInst(t *testing.T) {
	doc := parseSample(t)
	assert.Contains(t, doc.ProcInst, `xml`)
	assert.Contains(t, doc.ProcInst, `version="1.0"`)
}

func TestParse_Error(t *testing.T) {
	_, err := Parse(strings.NewReader("<unclosed>"))
	assert.Error(t, err)
}

func TestMust_PanicsOnError(t *testing.T) {
	assert.Panics(t, func() {
		Must(nil, assert.AnError)
	})
}

func TestMust_ReturnsDocOnSuccess(t *testing.T) {
	doc := &Document{}
	got := Must(doc, nil)
	assert.Same(t, doc, got)
}

func TestGetAttribute_Missing(t *testing.T) {
	doc := parseSample(t)
	assert.Nil(t, doc.Root.GetAttribute("missing"))
	assert.Equal(t, "", doc.Root.GetAttributeValue("missing"))
}

func TestGetChild(t *testing.T) {
	doc := parseSample(t)
	name := doc.Root.GetChild("name")
	require.NotNil(t, name)
	assert.Equal(t, "example", name.Text)
	assert.Nil(t, doc.Root.GetChild("does-not-exist"))
}

func TestGetChildren(t *testing.T) {
	doc := parseSample(t)
	params := doc.Root.GetChild("driver-parameters").GetChildren("parameter")
	assert.Len(t, params, 2)

	// GetChildren is direct-children-only — does NOT descend.
	directNames := doc.Root.GetChildren("name")
	assert.Len(t, directNames, 1, "GetChildren should only match direct children")
}

func TestFirstChild(t *testing.T) {
	doc := parseSample(t)
	assert.Equal(t, "name", doc.Root.FirstChild().Name)

	leaf := &Node{}
	assert.Nil(t, leaf.FirstChild())
}

func TestFindOneByName_DepthFirst(t *testing.T) {
	doc := parseSample(t)
	// `name` appears under <testsuite>, under each <parameter>, and under <testcase>.
	// Depth-first should hit the top-level <name> first.
	found := doc.Root.FindOneByName("name")
	require.NotNil(t, found)
	assert.Equal(t, "example", found.Text)

	assert.Nil(t, doc.Root.FindOneByName("does-not-exist"))
}

func TestFindByName_RecursiveAll(t *testing.T) {
	doc := parseSample(t)
	// <name> elements: 1 top-level, 2 in parameters, 1 in testcase = 4 total.
	all := doc.Root.FindByName("name")
	assert.Len(t, all, 4)
}

func TestXMLPretty_SelfClosing(t *testing.T) {
	n := &Node{Name: "leaf"}
	assert.Equal(t, "<leaf />\n", n.XMLPretty())
}

func TestXMLPretty_WithText(t *testing.T) {
	n := &Node{Name: "v", Text: "hello"}
	assert.Equal(t, "<v>hello</v>\n", n.XMLPretty())
}

func TestXMLPretty_EscapesText(t *testing.T) {
	n := &Node{Name: "v", Text: "a<b&c"}
	out := n.XMLPretty()
	assert.Contains(t, out, "&lt;")
	assert.Contains(t, out, "&amp;")
}

func TestXMLPretty_AttributesAndChildren(t *testing.T) {
	doc := parseSample(t)
	tc := doc.Root.GetChild("testcase")
	require.NotNil(t, tc)
	out := tc.XMLPretty()
	// Should contain nested elements with 2-space indentation.
	assert.Contains(t, out, "<testcase>")
	assert.Contains(t, out, "  <name>case1</name>")
	assert.Contains(t, out, "  <raw>DEADBEEF</raw>")
	assert.Contains(t, out, "</testcase>")
}

func TestXMLPretty_EscapesAttribute(t *testing.T) {
	n := &Node{
		Name:       "v",
		Attributes: []*Attribute{{Name: "a", Value: `"<>&`}},
	}
	out := n.XMLPretty()
	assert.Contains(t, out, `&#34;`)
	assert.Contains(t, out, `&lt;`)
	assert.Contains(t, out, `&amp;`)
}

func TestRoundTrip_TextMutation(t *testing.T) {
	// The harness's normalizeXml mutates Node.Text in place; make sure that's possible.
	doc := parseSample(t)
	name := doc.Root.GetChild("name")
	name.Text = "changed"
	assert.Equal(t, "changed", doc.Root.GetChild("name").Text)
}
