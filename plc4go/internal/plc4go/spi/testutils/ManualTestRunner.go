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
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/rs/zerolog/log"
	"math/rand"
	"strings"
	"time"
)

type ManualTestCase struct {
	Address           string
	ExpectedReadValue interface{}
	WriteValue        interface{}
}

type ManualTestSuite struct {
	ConnectionString string
	DriverManager    plc4go.PlcDriverManager
	TestCases        []ManualTestCase
}

func NewManualTestSuite(connectionString string, driverManager plc4go.PlcDriverManager) *ManualTestSuite {
	return &ManualTestSuite{
		ConnectionString: connectionString,
		DriverManager:    driverManager,
	}
}

func (m ManualTestSuite) AddTestCase(address string, expectedReadValue interface{}) {
	m.TestCases = append(m.TestCases, ManualTestCase{address, expectedReadValue, nil})
}

func (m ManualTestSuite) Run() {
	connectionResult := <-m.DriverManager.GetConnection(m.ConnectionString)
	if connectionResult.Err != nil {
		panic(connectionResult.Err)
	}
	connection := connectionResult
	log.Info().Msg("Reading all types in separate requests")
	// Run all entries separately:
	for _, testCase := range m.TestCases {
		fieldName := testCase.Address

		{
			// Prepare the read-request
			readRequestBuilder := connection.Connection.ReadRequestBuilder()
			readRequestBuilder.AddQuery(fieldName, testCase.Address)
			readRequest, err := readRequestBuilder.Build()
			if err != nil {
				panic(err)
			}

			// Execute the read request
			readResponseResult := <-readRequest.Execute()
			if readResponseResult.Err != nil {
				panic(err)
			}
			readResponse := readResponseResult.Response

			// Check the result
			assertEquals(1, len(readResponse.GetFieldNames()), fieldName)
			assertEquals(fieldName, readResponse.GetFieldNames()[0], fieldName)
			assertEquals(model.PlcResponseCode_OK, readResponse.GetResponseCode(fieldName), fieldName)
			assertNotNil(readResponse.GetValue(fieldName), fieldName)
			if readResponse.GetValue(fieldName).IsList() {
				plcList := readResponse.GetValue(fieldName).GetList()
				expectedValues := testCase.ExpectedReadValue.([]interface{})
				for j := 0; j < len(expectedValues); j++ {
					assertEquals(expectedValues[j], plcList[j].GetString(), fmt.Sprintf("%s[%d]", fieldName, j))
				}
			} else {
				assertEquals(fmt.Sprint(testCase.ExpectedReadValue), readResponse.GetValue(fieldName).GetString(), fieldName)
			}

			log.Info().Msg("Success")
		}

		{
			// Read all items in one big request.
			// Shuffle the list of test cases and run the test 10 times.
			log.Info().Msg("Reading all items together in random order")
			for i := 0; i < 100; i++ {
				log.Info().Msgf(" - run number %d of %d", i, 100)
				shuffledTestcases := append(make([]ManualTestCase, 0), m.TestCases...)
				rand.Seed(time.Now().UnixNano())
				rand.Shuffle(len(shuffledTestcases), func(i, j int) {
					shuffledTestcases[i], shuffledTestcases[j] = shuffledTestcases[j], shuffledTestcases[i]
				})

				sb := strings.Builder{}
				for _, testCase := range shuffledTestcases {
					sb.WriteString(testCase.Address)
					sb.WriteString(", ")
				}
				log.Info().Msgf("       using order: %s", sb.String())

				builder := connection.Connection.ReadRequestBuilder()
				for _, testCase := range shuffledTestcases {
					fieldName := testCase.Address
					builder.AddQuery(fieldName, testCase.Address)
				}
				readRequest, err := builder.Build()
				if err != nil {
					panic(err)
				}

				// Execute the read request
				readResponseResult := <-readRequest.Execute()
				if readResponseResult.Err != nil {
					panic(err)
				}
				readResponse := readResponseResult.Response

				// Check the result
				assertEquals(len(shuffledTestcases), len(readResponse.GetFieldNames()))
				for _, testCase := range shuffledTestcases {
					fieldName := testCase.Address
					assertEquals(model.PlcResponseCode_OK, readResponse.GetResponseCode(fieldName))
					assertNotNil(readResponse.GetValue(fieldName))
					if readResponse.GetValue(fieldName).IsList() {
						plcList := readResponse.GetValue(fieldName).GetList()
						expectedValues := testCase.ExpectedReadValue.([]interface{})
						for j := 0; j < len(expectedValues); j++ {
							assertEquals(expectedValues[j], plcList[j].GetString(), fmt.Sprintf("%s[%d]", fieldName, j))
						}
					} else {
						assertEquals(fmt.Sprint(testCase.ExpectedReadValue), readResponse.GetValue(fieldName).GetString(), fieldName)
					}
				}
			}
			log.Info().Msg("Success")
		}
	}
}

func assertEquals(expected interface{}, actual interface{}, message ...string) {
	if expected != actual {
		log.Error().Msgf("actual %s doesn't match expected %s\nmessage: %s", actual, expected, message)
		panic("Assertion failed")
	}
}

func assertNotNil(actual interface{}, message ...string) {
	if nil != actual {
		log.Error().Msgf("actual %s is nil\nmessage: %s", actual, message)
		panic("Assertion failed")
	}
}
