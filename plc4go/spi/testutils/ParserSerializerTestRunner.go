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
	"encoding/binary"
	"encoding/hex"
	"fmt"
	abethModel "github.com/apache/plc4x/plc4go/protocols/abeth/readwrite"
	adsModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite"
	bacnetModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite"
	df1Model "github.com/apache/plc4x/plc4go/protocols/df1/readwrite"
	eipModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite"
	firmataModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite"
	knxModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite"
	modbusModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite"
	s7Model "github.com/apache/plc4x/plc4go/protocols/s7/readwrite"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
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
	dirOffset := "/../../../"
	testFile := path + dirOffset + testPath
	info, err := os.Stat(testFile)
	if os.IsNotExist(err) {
		t.Errorf("Test-File %s doesn't exist", testFile)
	}
	if info.IsDir() {
		t.Errorf("Test-File %s refers to a directory", testFile)
	}

	// Open a reader for this file
	dat, err := os.Open(testFile)
	if err != nil {
		t.Error("Error opening file")
	}

	// Read the xml
	node := xmldom.Must(xmldom.Parse(dat)).Root

	if node.Name != "testsuite" {
		t.Error("Invalid document structure")
	}
	var byteOrder binary.ByteOrder
	if node.GetAttributeValue("byteOrder") != "LITTLE_ENDIAN" {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
	var (
		testsuiteName string
		protocolName  string
		outputFlavor  string
	)
	for _, childPtr := range node.Children {
		child := *childPtr
		if child.Name == "name" {
			testsuiteName = child.Text
		} else if child.Name == "protocolName" {
			protocolName = child.Text
		} else if child.Name == "outputFlavor" {
			outputFlavor = child.Text
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
				referenceSerialized := referenceXml.FirstChild().XMLPretty()

				// Get the raw input by decoding the hex-encoded binary input
				rawInput, err := hex.DecodeString(rawInputText)
				if err != nil {
					t.Errorf("Error decoding test input")
					return
				}
				var readBuffer utils.ReadBuffer
				if byteOrder == binary.LittleEndian {
					readBuffer = utils.NewReadBufferByteBased(rawInput, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
				} else {
					readBuffer = utils.NewReadBufferByteBased(rawInput)
				}

				// Parse the input according to the settings of the testcase
				var helper interface {
					Parse(typeName string, arguments []string, io utils.ReadBuffer) (interface{}, error)
				}
				switch protocolName {
				case "abeth":
					helper = new(abethModel.AbethParserHelper)
				case "ads":
					helper = new(adsModel.AdsParserHelper)
				case "bacnetip":
					helper = new(bacnetModel.BacnetipParserHelper)
				case "df1":
					helper = new(df1Model.Df1ParserHelper)
				case "eip":
					helper = new(eipModel.EipParserHelper)
				case "firmata":
					helper = new(firmataModel.FirmataParserHelper)
				case "modbus":
					helper = new(modbusModel.ModbusParserHelper)
				case "s7":
					helper = new(s7Model.S7ParserHelper)
				case "knxnetip":
					helper = new(knxModel.KnxnetipParserHelper)
				default:
					t.Errorf("Testsuite %s has not mapped parser for %s", testsuiteName, protocolName)
					return
				}
				_ = outputFlavor
				msg, err := helper.Parse(rootType, parserArguments, readBuffer)
				if err != nil {
					t.Error("Error parsing input data: ", err)
					return
				}

				{
					// First try to use the native xml writer
					serializable := msg.(utils.Serializable)
					buffer := utils.NewXmlWriteBuffer()
					if err := serializable.SerializeWithWriteBuffer(buffer); err == nil {
						actualXml := buffer.GetXmlString()
						if err := CompareResults([]byte(actualXml), []byte(referenceSerialized)); err != nil {
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
								testCaseName)
							assert.Equal(t, referenceSerialized, actualXml)
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
				if byteOrder == binary.LittleEndian {
					writeBuffer = utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
				} else {
					writeBuffer = utils.NewWriteBufferByteBased()
				}
				err = s.SerializeWithWriteBuffer(writeBuffer)
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
