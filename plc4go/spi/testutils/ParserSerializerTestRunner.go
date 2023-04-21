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

package testutils

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"

	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/stretchr/testify/assert"
	"github.com/subchen/go-xmldom"
)

type ParserSerializerTestsuite struct {
	name             string
	protocolName     string
	outputFlavor     string
	driverParameters map[string]string
	byteOrder        binary.ByteOrder
	parser           Parser
	rootTypeParser   func(utils.ReadBufferByteBased) (any, error)
	testcases        []ParserSerializerTestcase
}

type ParserSerializerTestcase struct {
	name            string
	rawInputText    string
	rootType        string
	referenceXml    string
	parserArguments []string
}

func (p *ParserSerializerTestsuite) Run(t *testing.T, testcase ParserSerializerTestcase) error {
	t.Run(testcase.name, func(t *testing.T) {
		t.Logf("running testsuite: %s test: %s", p.name, testcase.name)

		// Get the raw input by decoding the hex-encoded binary input
		rawInput, err := hex.DecodeString(testcase.rawInputText)
		if err != nil {
			t.Errorf("Error decoding test input")
			return
		}
		var readBuffer utils.ReadBuffer
		if p.byteOrder == binary.LittleEndian {
			readBuffer = utils.NewReadBufferByteBased(rawInput, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
		} else {
			readBuffer = utils.NewReadBufferByteBased(rawInput)
		}

		// Parse the input according to the settings of the testcase
		msg, err := p.parser.Parse(testcase.rootType, testcase.parserArguments, readBuffer)
		if err != nil {
			t.Error("Error parsing input data: ", err)
			return
		}

		{
			// First try to use the native xml writer
			serializable := msg.(utils.Serializable)
			buffer := utils.NewXmlWriteBuffer()
			if err := serializable.SerializeWithWriteBuffer(context.Background(), buffer); err == nil {
				actualXml := buffer.GetXmlString()
				if err := CompareResults(t, []byte(actualXml), []byte(testcase.referenceXml)); err != nil {
					border := strings.Repeat("=", 100)
					fmt.Printf(
						"\n"+
							// Border
							"%[1]s\n"+
							// Testcase name
							"%[4]s\n"+
							// diff detected message
							"Diff detected\n"+
							// Border
							"%[1]s\n"+
							// xml
							"%[2]s\n"+
							// Border
							"%[1]s\n%[1]s\n"+
							// Text
							"Differences were found after parsing (Use the above xml in the testsuite to disable this warning).\n"+
							// Diff
							"%[3]s\n"+
							// Double Border
							"%[1]s\n%[1]s\n",
						border,
						actualXml,
						err,
						testcase.name)
					assert.Equal(t, testcase.referenceXml, actualXml)
					t.Error("Error comparing the results: " + err.Error())
					return
				}
			}
		}

		// If all was ok, serialize the object again
		s, ok := msg.(utils.Serializable)
		if !ok {
			t.Error("Couldn't cast message to Serializable")
			return
		}
		var writeBuffer utils.WriteBufferByteBased
		if p.byteOrder == binary.LittleEndian {
			writeBuffer = utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
		} else {
			writeBuffer = utils.NewWriteBufferByteBased()
		}
		err = s.SerializeWithWriteBuffer(context.Background(), writeBuffer)
		if !ok {
			t.Error("Couldn't serialize message back to byte array")
			return
		}

		// Check if the output matches in size and content
		rawOutput := writeBuffer.GetBytes()
		if len(rawInput) != len(rawOutput) {
			t.Errorf("Missmatched number of bytes expected ->%d != %d<-actual\nexpected:\t%x\nactual:\t\t%x", len(rawInput), len(rawOutput), rawInput, rawOutput)
			t.Errorf("Hexdumps:\n%s", utils.DiffHex(rawInput, rawOutput))
			return
		}
		for i, val := range rawInput {
			if rawOutput[i] != val {
				t.Error("Raw output doesn't match input at position: " + strconv.Itoa(i))
				t.Errorf("Hexdumps:\n%s", utils.DiffHex(rawInput, rawOutput))
				return
			}
		}
	})
	return nil
}

type Parser interface {
	Parse(typeName string, arguments []string, io utils.ReadBuffer) (any, error)
}

func RunParserSerializerTestsuite(t *testing.T, testPath string, parser Parser, options ...WithOption) {
	t.Log("Extract testsuite options")
	var rootTypeParser func(utils.ReadBufferByteBased) (any, error)
	skippedTestCasesMap := map[string]bool{}
	for _, withOption := range options {
		switch option := withOption.(type) {
		case withRootTypeParser:
			t.Logf("Using root type parser for better output")
			rootTypeParser = option.rootTypeParser
		case withSkippedTestCases:
			t.Log("Skipping test cases:")
			for _, skippedTestCase := range option.skippedTestCases {
				t.Logf("Skipping %s", skippedTestCase)
				skippedTestCasesMap[skippedTestCase] = true
			}
		}
	}

	// Read the test-specification as XML file
	rootNode := ParseParserSerializerTestSuiteXml(t, testPath)

	// Parse the contents of the test-specification
	testsuite := ParseParserSerializerTestSuite(t, *rootNode, parser, rootTypeParser)

	t.Logf("Running %d testcases", len(testsuite.testcases))
	for _, testcase := range testsuite.testcases {
		t.Run(testcase.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Fatalf("\n-------------------------------------------------------\nPanic Failure\n%+v\n%s\n-------------------------------------------------------\n\n", err, debug.Stack())
				}
			}()
			if skippedTestCasesMap[testcase.name] {
				t.Logf("Testcase %s skipped", testcase.name)
				t.Skipf("Testcase %s skipped", testcase.name)
				return
			}
			t.Logf("Running testcase %s", testcase.name)
			if err := testsuite.Run(t, testcase); err != nil {
				t.Fatalf("\n-------------------------------------------------------\nFailure\n%+v\n-------------------------------------------------------\n\n", err)
			}
		})
	}
	t.Log("Done running testcases")
	// Execute the tests in the testsuite
	t.Logf(testsuite.name)
}

func ParseParserSerializerTestSuiteXml(t *testing.T, testPath string) *xmldom.Node {
	// Get the current working directory
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Check if the test-file is available
	dirOffset := "/../../../"
	testFile := path + dirOffset + testPath
	info, err := os.Stat(testFile)
	if os.IsNotExist(err) {
		t.Logf("Test-File %s doesn't exist", testFile)
		t.Fatal(err)
	}
	if info.IsDir() {
		t.Logf("Test-File %s refers to a directory", testFile)
		t.Fatal(err)
	}

	// Open a reader for this file
	dat, err := os.Open(testFile)
	if err != nil {
		t.Error("Error opening file")
	}

	// Read the xml
	node := xmldom.Must(xmldom.Parse(dat)).Root
	return node
}

func ParseParserSerializerTestSuite(t *testing.T, node xmldom.Node, parser Parser, rootTypeParser func(utils.ReadBufferByteBased) (any, error)) *ParserSerializerTestsuite {
	if node.Name != "testsuite" {
		t.Fatal("invalid document structure")
	}
	var byteOrder binary.ByteOrder
	if node.GetAttributeValue("byteOrder") != "LITTLE_ENDIAN" {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
	var testsuiteName string
	var protocolName string
	var outputFlavor string
	driverParameters := make(map[string]string)
	var testcases []ParserSerializerTestcase
	for _, childPtr := range node.Children {
		child := *childPtr
		switch child.Name {
		case "name":
			testsuiteName = child.Text
		case "protocolName":
			protocolName = child.Text
		case "outputFlavor":
			outputFlavor = child.Text
		case "driver-parameters":
			parameterList := child.FindByName("parameter")
			for _, parameter := range parameterList {
				nameElement := parameter.FindOneByName("name")
				valueElement := parameter.FindOneByName("value")
				if nameElement == nil || valueElement == nil {
					t.Fatal("invalid parameter found: no present")
				}
				name := nameElement.Text
				value := valueElement.Text
				if name == "" || value == "" {
					t.Fatal("invalid parameter found: empty")
				}
				driverParameters[name] = value
			}
		case "testcase":
			testcaseName := child.FindOneByName("name").Text
			rawInputText := (*(child.FindOneByName("raw"))).Text
			rootType := (*(child.FindOneByName("root-type"))).Text
			parserArgumentsXml := child.FindOneByName("parser-arguments")
			var parserArguments []string
			if parserArgumentsXml != nil {
				for _, parserArgumentXml := range parserArgumentsXml.Children {
					parserArguments = append(parserArguments, parserArgumentXml.Text)
				}
			}
			referenceXml := child.FindOneByName("xml")
			normalizeXml(referenceXml)
			referenceSerialized := referenceXml.FirstChild().XMLPretty()
			testcase := ParserSerializerTestcase{
				name:            testcaseName,
				rawInputText:    rawInputText,
				rootType:        rootType,
				referenceXml:    referenceSerialized,
				parserArguments: parserArguments,
			}
			testcases = append(testcases, testcase)
		default:
			t.Fatalf("invalid document structure. Unhandled element %s", child.Name)
		}
	}
	t.Logf("Parsed testsuite name: %s", testsuiteName)

	return &ParserSerializerTestsuite{
		name:             testsuiteName,
		protocolName:     protocolName,
		outputFlavor:     outputFlavor,
		driverParameters: driverParameters,
		byteOrder:        byteOrder,
		parser:           parser,
		rootTypeParser:   rootTypeParser,
		testcases:        testcases,
	}
}

// Mainly remove linebreaks from text content.
func normalizeXml(input *xmldom.Node) {
	if len(input.Children) > 0 {
		for _, child := range input.Children {
			normalizeXml(child)
		}
	}
	if len(input.Text) > 0 {
		if strings.Contains(input.Text, "\n") {
			input.Text = strings.Replace(input.Text, "\n", "", -1)
		}
	}
}
