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
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/rs/zerolog/log"
	"math/rand"
	"reflect"
	"strings"
	"testing"
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
	t                *testing.T
}

func NewManualTestSuite(connectionString string, driverManager plc4go.PlcDriverManager, t *testing.T) *ManualTestSuite {
	return &ManualTestSuite{
		ConnectionString: connectionString,
		DriverManager:    driverManager,
		t:                t,
	}
}

func (m *ManualTestSuite) AddTestCase(address string, expectedReadValue interface{}) {
	m.TestCases = append(m.TestCases, ManualTestCase{address, expectedReadValue, nil})
}

func (m *ManualTestSuite) Run() {
	connectionResult := <-m.DriverManager.GetConnection(m.ConnectionString)
	if connectionResult.GetErr() != nil {
		panic(connectionResult.GetErr())
	}
	connection := connectionResult
	log.Info().Msg("Reading all types in separate requests")
	// Run all entries separately:
	for _, testCase := range m.TestCases {
		fieldName := testCase.Address
		m.t.Run(fieldName, func(t *testing.T) {
			m.runSingleTest(t, connection, fieldName, testCase)
		})
	}
	m.t.Run("combinedTest", func(t *testing.T) {
		m.runBurstTest(t, connection)
	})
}

func (m *ManualTestSuite) runSingleTest(t *testing.T, connection plc4go.PlcConnectionConnectResult, fieldName string, testCase ManualTestCase) {
	// Prepare the read-request
	readRequestBuilder := connection.GetConnection().ReadRequestBuilder()
	readRequestBuilder.AddQuery(fieldName, testCase.Address)
	readRequest, err := readRequestBuilder.Build()
	if err != nil {
		panic(err)
	}

	// Execute the read request
	readResponseResult := <-readRequest.Execute()
	if readResponseResult.GetErr() != nil {
		t.Errorf("Error getting response %v", readResponseResult.GetErr())
		t.FailNow()
		return
	}
	readResponse := readResponseResult.GetResponse()

	// Check the result
	assertEquals(t, 1, len(readResponse.GetFieldNames()), fieldName)
	assertEquals(t, fieldName, readResponse.GetFieldNames()[0], fieldName)
	assertEquals(t, model.PlcResponseCode_OK, readResponse.GetResponseCode(fieldName), fieldName)
	assertNotNil(t, readResponse.GetValue(fieldName), fieldName)
	expectation := reflect.ValueOf(testCase.ExpectedReadValue)
	if readResponse.GetValue(fieldName).IsList() && (expectation.Kind() == reflect.Slice || expectation.Kind() == reflect.Array) {
		plcList := readResponse.GetValue(fieldName).GetList()
		for j := 0; j < expectation.Len(); j++ {
			assertEquals(t, expectation.Index(j).Interface(), plcList[j], fmt.Sprintf("%s[%d]", fieldName, j))
		}
	} else {
		assertEquals(t, fmt.Sprint(testCase.ExpectedReadValue), readResponse.GetValue(fieldName).GetString(), fieldName)
	}
}

func (m *ManualTestSuite) runBurstTest(t *testing.T, connection plc4go.PlcConnectionConnectResult) {
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

		builder := connection.GetConnection().ReadRequestBuilder()
		for _, testCase := range shuffledTestcases {
			fieldName := testCase.Address
			builder.AddQuery(fieldName, testCase.Address)
		}
		readRequest, err := builder.Build()
		if err != nil {
			t.Errorf("Error building request %v", err)
			return
		}

		// Execute the read request
		readResponseResult := <-readRequest.Execute()
		if readResponseResult.GetErr() != nil {
			t.Errorf("Error getting response %v", err)
			return
		}
		readResponse := readResponseResult.GetResponse()

		// Check the result
		assertEquals(t, len(shuffledTestcases), len(readResponse.GetFieldNames()))
		for _, testCase := range shuffledTestcases {
			fieldName := testCase.Address
			assertEquals(t, model.PlcResponseCode_OK, readResponse.GetResponseCode(fieldName))
			assertNotNil(t, readResponse.GetValue(fieldName))
			expectation := reflect.ValueOf(testCase.ExpectedReadValue)
			if readResponse.GetValue(fieldName).IsList() && (expectation.Kind() == reflect.Slice || expectation.Kind() == reflect.Array) {
				plcList := readResponse.GetValue(fieldName).GetList()
				for j := 0; j < expectation.Len(); j++ {
					assertEquals(t, expectation.Index(j).Interface(), plcList[j], fmt.Sprintf("%s[%d]", fieldName, j))
				}
			} else {
				assertEquals(t, fmt.Sprint(testCase.ExpectedReadValue), readResponse.GetValue(fieldName).GetString(), fieldName)
			}
		}
	}
}

func assertEquals(t *testing.T, expected interface{}, actual interface{}, message ...string) {
	switch actual.(type) {
	case values.PlcBOOL:
		actual = actual.(values.PlcBOOL).GetBool()
	case values.PlcWORD:
		actual = actual.(values.PlcWORD).GetInt8()
	}
	if expected != actual {
		t.Errorf("actual %v doesn't match expected %v\nmessage: %s", actual, expected, message)
		t.FailNow()
	}
}

func assertNotNil(t *testing.T, actual interface{}, message ...string) {
	if actual == nil {
		t.Errorf("actual %v is nil\nmessage: %v", actual, message)
		t.FailNow()
	}
}
