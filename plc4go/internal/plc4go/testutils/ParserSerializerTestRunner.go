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
    "github.com/subchen/go-xmldom"
    "os"
    model2 "plc4x.apache.org/plc4go/v0/internal/plc4go/modbus/readwrite"
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
    "strconv"
    "strings"
    "testing"
)

func RunParserSerializerTestsuite(t *testing.T, testPath string) {
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
	var testsuiteName string
	for _, childPtr := range node.Children {
		curFailed := false
		child := *childPtr
		if child.Name == "name" {
			testsuiteName = child.Text
		} else if child.Name != "testcase" {
			t.Error("Invalid document structure")
			curFailed = true
		} else {
			t.Logf("running testsuite: %s test: %s", testsuiteName, (*(child.FindOneByName("name"))).Text)
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
				t.Fail()
				curFailed = true
			}
			readBuffer := utils.NewReadBuffer(rawInput)

			// Parse the input according to the settings of the testcase
			helper := new(model2.ModbusParserHelper)
			msg, err := helper.Parse(rootType, parserArguments, readBuffer)
			if err != nil {
				t.Error("Error parsing input data: " + err.Error())
				t.Fail()
				curFailed = true
			}

			// Serialize the parsed object to XML
			actualSerialized, err := xml.Marshal(msg)
			if err != nil {
				t.Error("Error serializing the actual message: " + err.Error())
				t.Fail()
				curFailed = true
			}

			// Compare the actual and the expected xml
			err = CompareResults(actualSerialized, []byte(referenceSerialized))
			if err != nil {
				t.Error("Error comparing the results: " + err.Error())
				t.Fail()
				curFailed = true
			}

			// If all was ok, serialize the object again
			s, ok := msg.(utils.Serializable)
			if !ok {
				t.Error("Couldn't cast message to Serializable")
				t.Fail()
				curFailed = true
			}
			writeBuffer := utils.NewWriteBuffer()
			err = s.Serialize(*writeBuffer)
			if !ok {
				t.Error("Couldn't serialize message back to byte array")
				t.Fail()
				curFailed = true
			}

			// Check if the output matches in size and content
			rawOutput := writeBuffer.GetBytes()
			if len(rawInput) != len(rawOutput) {
				t.Error("Couldn't serialize message back to byte array")
				t.Fail()
				curFailed = true
			}
			for i, val := range rawInput {
				if rawOutput[i] != val {
					t.Error("Raw output doesn't match input at position: " + strconv.Itoa(i))
					t.Fail()
					curFailed = true
				}
			}

			if curFailed {
				// All worked
				t.Logf("FAILED")
			} else {
				// All worked
				t.Logf("SUCCESS")
			}
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
