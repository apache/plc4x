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
	"github.com/pkg/errors"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/values"

	"github.com/stretchr/testify/assert"
)

type ManualTestCase struct {
	Address           string
	ExpectedReadValue any
	WriteValue        any
	UnwrappedValue    bool
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

func (m *ManualTestSuite) AddTestCase(address string, expectedReadValue any, testCaseOptions ...WithTestCaseOption) {
	testCase := ManualTestCase{Address: address, ExpectedReadValue: expectedReadValue, UnwrappedValue: true}
	for _, testCaseOption := range testCaseOptions {
		testCaseOption(testCase)
	}
	m.TestCases = append(m.TestCases, testCase)
}

type WithTestCaseOption func(testCase ManualTestCase)

// WithUnwrappedValue is a WithTestCaseOption which can be used to control if plc4go.PlcValue
func WithUnwrappedValue(unwrap bool) WithTestCaseOption {
	return func(testCase ManualTestCase) {
		testCase.UnwrappedValue = unwrap
	}
}

func (m *ManualTestSuite) Run(t *testing.T) plc4go.PlcConnection {
	connectionResult := <-m.DriverManager.GetConnection(m.ConnectionString)
	if err := connectionResult.GetErr(); err != nil {
		tracer, ok := errors.Cause(err).(interface{ StackTrace() errors.StackTrace })
		if ok {
			stackTrace := tracer.StackTrace()
			for _, f := range stackTrace {
				fmt.Printf("%+s:%d\n", f, f)
			}
			fmt.Println()
		}
		m.t.Error(err)
		m.t.FailNow()
	}
	connection := connectionResult.GetConnection()
	m.t.Cleanup(func() {
		connection.Close()
	})
	t.Log("Reading all types in separate requests")
	// Run all entries separately:
	for _, testCase := range m.TestCases {
		tagName := testCase.Address
		m.t.Run(tagName, func(t *testing.T) {
			m.runSingleTest(t, connection, tagName, testCase)
		})
	}
	m.t.Run("combinedTest", func(t *testing.T) {
		m.runBurstTest(t, connection)
	})
	return connection
}

func (m *ManualTestSuite) runSingleTest(t *testing.T, connection plc4go.PlcConnection, tagName string, testCase ManualTestCase) {
	// Prepare the read-request
	readRequestBuilder := connection.ReadRequestBuilder()
	readRequestBuilder.AddTagAddress(tagName, testCase.Address)
	readRequest, err := readRequestBuilder.Build()
	if err != nil {
		t.Fatal(err)
		return
	}

	// Execute the read request
	readResponseResult := <-readRequest.Execute()
	if readResponseResult.GetErr() != nil {
		t.Fatalf("Error getting response %v", readResponseResult.GetErr())
		return
	}
	readResponse := readResponseResult.GetResponse()

	// Check the result
	assert.Equalf(t, 1, len(readResponse.GetTagNames()), "response should have a tag for %s", tagName)
	assert.Equalf(t, tagName, readResponse.GetTagNames()[0], "first tag should be equal to %s", tagName)
	assert.Equalf(t, model.PlcResponseCode_OK, readResponse.GetResponseCode(tagName), "response code should be ok for %s", tagName)
	assert.NotNil(t, readResponse.GetValue(tagName), tagName)
	expectation := reflect.ValueOf(testCase.ExpectedReadValue)
	if readResponse.GetValue(tagName).IsList() && (expectation.Kind() == reflect.Slice || expectation.Kind() == reflect.Array) {
		plcList := readResponse.GetValue(tagName).GetList()
		for j := 0; j < expectation.Len(); j++ {
			var actual any
			actual = plcList[j]
			if testCase.UnwrappedValue {
				switch actualCasted := actual.(type) {
				case values.PlcBOOL:
					actual = actualCasted.GetBool()
				case values.PlcWORD:
					actual = actualCasted.GetInt8()
				default:
					t.Fatalf("%T not yet mapped", actualCasted)
				}
			}
			assert.Equal(t, expectation.Index(j).Interface(), actual, fmt.Sprintf("%s[%d]", tagName, j))
		}
	} else {
		assert.Equal(t, fmt.Sprint(testCase.ExpectedReadValue), readResponse.GetValue(tagName).GetString(), tagName)
	}
}

func (m *ManualTestSuite) runBurstTest(t *testing.T, connection plc4go.PlcConnection) {
	// Read all items in one big request.
	// Shuffle the list of test cases and run the test 10 times.
	t.Log("Reading all items together in random order")
	for i := 0; i < 100; i++ {
		t.Logf(" - run number %d of %d", i, 100)
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
		t.Logf("       using order: %s", sb.String())

		builder := connection.ReadRequestBuilder()
		for _, testCase := range shuffledTestcases {
			tagName := testCase.Address
			builder.AddTagAddress(tagName, testCase.Address)
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
		assert.Equal(t, len(shuffledTestcases), len(readResponse.GetTagNames()))
		for _, testCase := range shuffledTestcases {
			tagName := testCase.Address
			assert.Equalf(t, model.PlcResponseCode_OK, readResponse.GetResponseCode(tagName), "response code should be ok for %s", tagName)
			assert.NotNil(t, readResponse.GetValue(tagName))
			expectation := reflect.ValueOf(testCase.ExpectedReadValue)
			if readResponse.GetValue(tagName).IsList() && (expectation.Kind() == reflect.Slice || expectation.Kind() == reflect.Array) {
				plcList := readResponse.GetValue(tagName).GetList()
				for j := 0; j < expectation.Len(); j++ {
					var actual any
					actual = plcList[j]
					if testCase.UnwrappedValue {
						switch actualCasted := actual.(type) {
						case values.PlcBOOL:
							actual = actualCasted.GetBool()
						case values.PlcWORD:
							actual = actualCasted.GetInt8()
						default:
							t.Fatalf("%T not yet mapped", actualCasted)
						}
					}
					assert.Equal(t, expectation.Index(j).Interface(), actual, fmt.Sprintf("%s[%d]", tagName, j))
				}
			} else {
				assert.Equal(t, fmt.Sprint(testCase.ExpectedReadValue), readResponse.GetValue(tagName).GetString(), tagName)
			}
		}
	}
}
