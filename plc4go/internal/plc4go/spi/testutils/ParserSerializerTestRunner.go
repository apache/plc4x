//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package testutils

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	adsModel "github.com/apache/plc4x/plc4go/internal/plc4go/ads/readwrite"
	knxModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite"
	modbusModel "github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite"
	s7Model "github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/rs/zerolog/log"
	"github.com/subchen/go-xmldom"
	"os"
	"strconv"
	"strings"
	"testing"
)

func RunParserSerializerTestsuite(t *testing.T, testPath string, skippedTestCases ...string) {
	skippedTestCasesMap := map[string]bool{}
	for _, skippedTestCase := range skippedTestCases {
		skippedTestCasesMap[skippedTestCase] = true
	}
	// Get the current working directory
	path, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	// Check if the test-file is available
	info, err := os.Stat(path + "/../../../../" + testPath)
	if os.IsNotExist(err) {
		t.Error("Test-File doesn't exist")
	}
	if info.IsDir() {
		t.Error("Test-File refers to a directory")
	}

	// Open a reader for this file
	dat, err := os.Open(path + "/../../../../" + testPath)
	if err != nil {
		t.Error("Error opening file")
	}

	// Read the xml
	node := xmldom.Must(xmldom.Parse(dat)).Root

	if node.Name != "testsuite" {
		t.Error("Invalid document structure")
	}
	littleEndian := node.GetAttributeValue("bigEndian") != "true"
	var testsuiteName string
	for _, childPtr := range node.Children {
		child := *childPtr
		if child.Name == "name" {
			testsuiteName = child.Text
		} else if child.Name != "testcase" {
			t.Error("Invalid document structure")
			return
		} else {
			testCaseName := child.FindOneByName("name").Text
			t.Run(testCaseName, func(t *testing.T) {
				if skippedTestCasesMap[testCaseName] {
					log.Warn().Msgf("Testcase %s skipped", testCaseName)
					t.Skipf("Testcase %s skipped", testCaseName)
					return
				}
				t.Logf("running testsuite: %s test: %s", testsuiteName, testCaseName)
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
				referenceSerialized := referenceXml.FirstChild().XML()

				// Get the raw input by decoding the hex-encoded binary input
				rawInput, err := hex.DecodeString(rawInputText)
				if err != nil {
					t.Errorf("Error decoding test input")
					return
				}
				var readBuffer *utils.ReadBuffer
				if littleEndian {
					readBuffer = utils.NewLittleEndianReadBuffer(rawInput)
				} else {
					readBuffer = utils.NewReadBuffer(rawInput)
				}

				// Parse the input according to the settings of the testcase
				var helper interface {
					Parse(typeName string, arguments []string, io *utils.ReadBuffer) (interface{}, error)
				}
				switch testsuiteName {
				case "Modbus":
					helper = new(modbusModel.ModbusParserHelper)
				case "Beckhoff ADS/AMS":
					helper = new(adsModel.AdsParserHelper)
				case "S7":
					helper = new(s7Model.S7ParserHelper)
				case "KNXNet/IP":
					helper = new(knxModel.KnxnetipParserHelper)
				default:
					t.Errorf("Testsuite %s has not mapped parser", testsuiteName)
					return
				}
				msg, err := helper.Parse(rootType, parserArguments, readBuffer)
				if err != nil {
					t.Error("Error parsing input data: " + err.Error())
					return
				}

				// Serialize the parsed object to XML
				actualSerialized, err := xml.Marshal(msg)
				if err != nil {
					t.Error("Error serializing the actual message: " + err.Error())
					return
				}

				// Compare the actual and the expected xml
				err = CompareResults(actualSerialized, []byte(referenceSerialized))
				if err != nil {
					t.Error("Error comparing the results: " + err.Error())
					return
				}

				// If all was ok, serialize the object again
				s, ok := msg.(utils.Serializable)
				if !ok {
					t.Error("Couldn't cast message to Serializable")
					return
				}
				var writeBuffer *utils.WriteBuffer
				if littleEndian {
					writeBuffer = utils.NewLittleEndianWriteBuffer()
				} else {
					writeBuffer = utils.NewWriteBuffer()
				}
				err = s.Serialize(*writeBuffer)
				if !ok {
					t.Error("Couldn't serialize message back to byte array")
					return
				}

				// Check if the output matches in size and content
				rawOutput := writeBuffer.GetBytes()
				if len(rawInput) != len(rawOutput) {
					t.Errorf("Missmatched number of bytes expected ->%d != %d<-actual\nexpected:\t%x\nactual:\t\t%x", len(rawInput), len(rawOutput), rawInput, rawOutput)
					t.Errorf("Hexdumps\nexpected:\n%s\nactual:\n%s\n", utils.Dump(rawInput), utils.Dump(rawOutput))
					return
				}
				for i, val := range rawInput {
					if rawOutput[i] != val {
						t.Error("Raw output doesn't match input at position: " + strconv.Itoa(i))
						t.Errorf("Hexdumps\nexpected:\n%s\nactual:\n%s\n", utils.Dump(rawInput), utils.Dump(rawOutput))
						return
					}
				}
			})
		}
	}
	fmt.Printf("name = %v\n", node.Name)
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
