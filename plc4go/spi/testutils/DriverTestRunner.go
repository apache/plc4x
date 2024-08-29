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
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/subchen/go-xmldom"

	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/options/converter"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

var DriverTestsuiteConnectTimeout = 30 * time.Second

type DriverTestsuite struct {
	name             string
	protocolName     string
	outputFlavor     string
	driverName       string
	driverParameters map[string]string
	byteOrder        binary.ByteOrder
	parser           XmlParser
	rootTypeParser   func(utils.ReadBufferByteBased) (any, error)
	setupSteps       []DriverTestStep
	teardownSteps    []DriverTestStep
	testcases        []DriverTestcase
}

type XmlParser interface {
	Parse(typeName string, xmlString string, parserArguments ...string) (any, error)
}

// WithRootTypeParser can be used to output the root type of protocol for better debugging
func WithRootTypeParser(rootTypeParser func(utils.ReadBufferByteBased) (any, error)) config.WithOption {
	return withRootTypeParser{rootTypeParser: rootTypeParser}
}

type withRootTypeParser struct {
	options.Option
	rootTypeParser func(utils.ReadBufferByteBased) (any, error)
}

// WithSkippedTestCases can be used to skip test cases
func WithSkippedTestCases(skippedTestCases ...string) config.WithOption {
	return withSkippedTestCases{skippedTestCases: skippedTestCases}
}

type withSkippedTestCases struct {
	options.Option
	skippedTestCases []string
}

type TestTransportInstance interface {
	transports.TransportInstance
	FillReadBuffer(data []byte)
	GetNumDrainableBytes() uint32
	DrainWriteBuffer(numBytes uint32) []byte
}

func (m DriverTestsuite) Run(t *testing.T, driverManager plc4go.PlcDriverManager, testcase DriverTestcase) error {
	var driverParameters []string
	for key, value := range m.driverParameters {
		driverParameters = append(driverParameters, fmt.Sprintf("%s=%s", key, value))
	}
	optionsString := ""
	if len(driverParameters) > 0 {
		optionsString = "?" + strings.Join(driverParameters, "&")
	}
	// Get a connection
	t.Log("getting a connection")
	connectionChan := driverManager.GetConnection(m.driverName + ":test://hurz" + optionsString)
	timer := time.NewTimer(DriverTestsuiteConnectTimeout)
	t.Cleanup(func() { utils.CleanupTimer(timer) })
	var connectionResult plc4go.PlcConnectionConnectResult
	select {
	case connectionResult = <-connectionChan:
	case <-timer.C:
		t.Fatalf("timeout")
	}

	if connectionResult.GetErr() != nil {
		return errors.Wrap(connectionResult.GetErr(), "error getting a connection")
	}
	connection := connectionResult.GetConnection()
	t.Cleanup(func() {
		timeout := time.NewTimer(30 * time.Second)
		t.Cleanup(func() {
			utils.CleanupTimer(timeout)
		})
		select {
		case result := <-connection.Close():
			assert.NoError(t, result.GetErr())
		case <-timeout.C:
			t.Error("timeout closing connection")
		}
	})
	utils.NewAsciiBoxWriter()
	m.LogDelimiterSection(t, "=", "Executing testcase: %s", testcase.name)

	if len(m.setupSteps) > 0 {
		// Run the setup steps
		m.LogDelimiterSection(t, "-", "Performing setup for: %s", testcase.name)
		for _, testStep := range m.setupSteps {
			err := m.ExecuteStep(t, connection, &testcase, testStep)
			if err != nil {
				return errors.Wrap(err, "error in setup step "+testStep.name)
			}
			// We sleep a bit to not run too fast into the post setup steps and give connections a bit time to settle built up
			time.Sleep(100 * time.Millisecond) // TODO: this is really bad as on CI sometimes those sleeps are not enough...
		}
		m.LogDelimiterSection(t, "-", "setup done for: %s", testcase.name)
	} else {
		t.Log("no setup steps")
	}

	if len(testcase.steps) > 0 {
		// Run the actual scenario steps
		m.LogDelimiterSection(t, "-", "Running testcases for: %s", testcase.name)
		for _, testStep := range testcase.steps {
			err := m.ExecuteStep(t, connection, &testcase, testStep)
			if err != nil {
				return errors.Wrap(err, "error in step "+testStep.name)
			}
			time.Sleep(100 * time.Millisecond) // TODO: this is really bad as on CI sometimes those sleeps are not enough...
		}
		m.LogDelimiterSection(t, "-", "testcases done for: %s", testcase.name)
	} else {
		t.Log("no testcase steps")
	}

	if len(m.teardownSteps) > 0 {
		// Run the teardown steps
		m.LogDelimiterSection(t, "-", "Performing teardown for: %s", testcase.name)
		for _, testStep := range m.teardownSteps {
			err := m.ExecuteStep(t, connection, &testcase, testStep)
			if err != nil {
				return errors.Wrap(err, "error in teardown step "+testStep.name)
			}
			time.Sleep(100 * time.Millisecond) // TODO: this is really bad as on CI sometimes those sleeps are not enough...
		}
		m.LogDelimiterSection(t, "-", "teardown done for: %s", testcase.name)
	} else {
		t.Log("no teardown steps")
	}

	m.LogDelimiterSection(t, "=", "done testcase: %s", testcase.name)
	return nil
}

func (m DriverTestsuite) ExecuteStep(t *testing.T, connection plc4go.PlcConnection, testcase *DriverTestcase, step DriverTestStep) error {
	mc, ok := connection.(spi.TransportInstanceExposer)
	if !ok {
		return errors.New("couldn't access connections transport instance")
	}
	testTransportInstance, ok := mc.GetTransportInstance().(TestTransportInstance)
	if !ok {
		return errors.New("transport must be of type TestTransport")
	}

	start := time.Now()
	m.LogDelimiterSection(t, "┄", " - Executing step (%s): %s", step.stepType, step.name)
	switch step.stepType {
	case StepTypeApiRequest:
		switch step.payload.Name {
		case "TestReadRequest":
			// Assemble a read-request according to the information in the test xml
			t.Log("Assemble read request")
			rrb := connection.ReadRequestBuilder()
			for _, tagNode := range step.payload.GetChild("tags").GetChildren("tag") {
				tagName := tagNode.GetChild("name").Text
				tagAddress := tagNode.GetChild("address").Text
				rrb.AddTagAddress(tagName, tagAddress)
			}
			readRequest, err := rrb.Build()
			if err != nil {
				return errors.Wrap(err, "Error creating read-request")
			}

			// Execute the read-request and store the response-channel in the testcase.
			t.Logf("Execute read request (%T)\n%[1]s", readRequest)
			if testcase.readRequestResultChannel != nil {
				return errors.New("testcase read-request result channel already occupied")
			}
			testcase.readRequestResultChannel = readRequest.Execute()
			t.Log("request executed")
		case "TestWriteRequest":
			t.Log("Assemble write request")
			wrb := connection.WriteRequestBuilder()
			for _, tagNode := range step.payload.GetChild("tags").GetChildren("tag") {
				tagName := tagNode.GetChild("name").Text
				tagAddress := tagNode.GetChild("address").Text

				he, ok := connection.(spi.HandlerExposer)
				if !ok {
					return errors.New("connection is not a HandlerExposer")
				}
				tag, err := he.GetPlcTagHandler().ParseTag(tagAddress)
				if err != nil {
					return errors.Wrapf(err, "error parsing address: %s", tagAddress)
				}
				if len(tag.GetArrayInfo()) > 0 {
					var tagValue []string
					for _, valueChild := range tagNode.GetChildren("value") {
						tagValue = append(tagValue, valueChild.Text)
					}
					wrb.AddTagAddress(tagName, tagAddress, tagValue)
				} else {
					tagValue := tagNode.GetChild("value").Text
					wrb.AddTagAddress(tagName, tagAddress, tagValue)
				}
			}
			writeRequest, err := wrb.Build()
			if err != nil {
				return errors.Wrap(err, "Error creating write-request")
			}
			t.Logf("Execute write request (%T)\n%[1]s", writeRequest)
			if testcase.writeRequestResultChannel != nil {
				return errors.New("testcase write-request result channel already occupied")
			}
			testcase.writeRequestResultChannel = writeRequest.Execute()
			t.Log("request executed")
		}
	case StepTypeApiResponse:
		switch step.payload.Name {
		case "PlcReadResponse":
			if testcase.readRequestResultChannel == nil {
				return errors.New("no read response expected")
			}
			t.Log("Waiting for read request result")
			readRequestResult := <-testcase.readRequestResultChannel
			if readRequestResult.GetErr() != nil {
				return errors.Wrap(readRequestResult.GetErr(), "error sending response")
			}
			// Serialize the response to XML
			xmlWriteBuffer := utils.NewXmlWriteBuffer()
			response := readRequestResult.GetResponse()
			t.Logf("Got response (%T)\n%[1]s", response)
			err := response.(utils.Serializable).SerializeWithWriteBuffer(context.Background(), xmlWriteBuffer)
			if err != nil {
				return errors.Wrap(err, "error serializing response")
			}
			actualResponse := xmlWriteBuffer.GetXmlString()
			// Get the reference XML
			referenceSerialized := step.payload.XMLPretty()
			// Compare the results
			err = CompareResults(t, []byte(actualResponse), []byte(referenceSerialized))
			if err != nil {
				return errors.Wrap(err, "Error comparing the results")
			}
			// Reset read channel
			testcase.readRequestResultChannel = nil
		case "PlcWriteResponse":
			if testcase.writeRequestResultChannel == nil {
				return errors.New("no write response expected")
			}
			t.Log("Waiting for write request result")
			writeResponseResult := <-testcase.writeRequestResultChannel
			if writeResponseResult.GetErr() != nil {
				return errors.Wrap(writeResponseResult.GetErr(), "error sending response")
			}
			// Serialize the response to XML
			xmlWriteBuffer := utils.NewXmlWriteBuffer()
			response := writeResponseResult.GetResponse()
			t.Logf("Got response (%T)\n%[1]s", response)
			err := response.(utils.Serializable).SerializeWithWriteBuffer(context.Background(), xmlWriteBuffer)
			if err != nil {
				return errors.Wrap(err, "error serializing response")
			}
			actualResponse := xmlWriteBuffer.GetXmlString()
			// Get the reference XML
			referenceSerialized := step.payload.XMLPretty()
			// Compare the results
			err = CompareResults(t, []byte(actualResponse), []byte(referenceSerialized))
			if err != nil {
				return errors.Wrap(err, "Error comparing the results")
			}
			// Reset write channel
			testcase.writeRequestResultChannel = nil
		}
	case StepTypeOutgoingPlcMessage:
		typeName := step.payload.Name
		payloadString := step.payload.XMLPretty()

		// Parse the xml into a real model
		t.Log("parsing xml")
		expectedMessage, err := m.parseMessage(typeName, payloadString, step)
		if err != nil {
			return errors.Wrap(err, "Error parsing message")
		}
		t.Logf("Parsed message (%T) from xml\n%[1]s", expectedMessage)

		// Serialize the model into bytes
		t.Log("Serialize the model into bytes and compare output")
		expectedSerializable, ok := expectedMessage.(utils.Serializable)
		if !ok {
			return errors.Errorf("error converting type %t into Serializable type", expectedMessage)
		}
		var expectedWriteBuffer utils.WriteBufferByteBased
		if m.byteOrder == binary.BigEndian {
			expectedWriteBuffer = utils.NewWriteBufferByteBased()
		} else {
			t.Log("using little endian")
			expectedWriteBuffer = utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
		}
		err = expectedSerializable.SerializeWithWriteBuffer(context.Background(), expectedWriteBuffer)
		if err != nil {
			return errors.Wrap(err, "error serializing expectedMessage")
		}
		expectedRawOutput := expectedWriteBuffer.GetBytes()
		expectedRawOutputLength := uint32(len(expectedRawOutput))

		startTransportPolling := time.Now()
		// Read exactly this amount of bytes from the transport
		t.Logf("Reading bytes from transport instance (expectedRawOutputLength %d)", expectedRawOutputLength)
		for testTransportInstance.GetNumDrainableBytes() < expectedRawOutputLength {
			if time.Since(startTransportPolling) > 2*time.Second {
				drainableBytes := testTransportInstance.GetNumDrainableBytes()
				actualRawOutput := testTransportInstance.DrainWriteBuffer(drainableBytes)
				return errors.Errorf("Not enough data available: actual(%d)<expected(%d), \nactual:   %#X\nexpected: %#X\nHexdumps:\n%s",
					drainableBytes, expectedRawOutputLength, actualRawOutput, expectedRawOutput, utils.DiffHex(expectedRawOutput, actualRawOutput))
			}
			time.Sleep(2 * time.Millisecond)
		}
		actualRawOutput := testTransportInstance.DrainWriteBuffer(expectedRawOutputLength)
		if testTransportInstance.GetNumDrainableBytes() != 0 {
			t.Logf("leftover drainable bytes (%d)", testTransportInstance.GetNumDrainableBytes())
		}

		var bufferFactory func([]byte, ...utils.ReadBufferByteBasedOptions) utils.ReadBufferByteBased
		if m.byteOrder == binary.BigEndian {
			bufferFactory = utils.NewReadBufferByteBased
		} else {
			bufferFactory = func(bytes []byte, options ...utils.ReadBufferByteBasedOptions) utils.ReadBufferByteBased {
				t.Log("using little endian")
				return utils.NewReadBufferByteBased(bytes, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
			}
		}
		// Compare the bytes read with the ones we expect
		t.Log("Comparing outputs")
		for i := range expectedRawOutput {
			if expectedRawOutput[i] != actualRawOutput[i] {
				if m.rootTypeParser != nil {
					readBufferByteBased := bufferFactory(actualRawOutput)
					actual, err := m.rootTypeParser(readBufferByteBased)
					t.Logf("A readable render of expectation:\n%v\nvs actual paket\n%v\n%#v", expectedSerializable, actual, err)
				}
				return errors.Errorf("actual output doesn't match expected output:\nactual:   %#X\nexpected: %#X\nHexdumps:\n%s", actualRawOutput, expectedRawOutput, utils.DiffHex(expectedRawOutput, actualRawOutput))
			}
		}
		t.Log("outputs are matching")
		// If there's a difference, parse the input and display it to simplify debugging
	case StepTypeOutgoingPlcBytes:
		// Read exactly this amount of bytes from the transport
		t.Log("Reading bytes")
		expectedRawInput, err := hex.DecodeString(step.payload.Text)
		if err != nil {
			return errors.Wrap(err, "error decoding hex-encoded byte data")
		}
		t.Logf("\n%s", hex.Dump(expectedRawInput))
		rawInput := testTransportInstance.DrainWriteBuffer(uint32(len(expectedRawInput)))

		// Compare the bytes read with the ones we expect
		t.Log("Comparing bytes")
		for i := range expectedRawInput {
			if expectedRawInput[i] != rawInput[i] {
				return errors.Errorf("actual output doesn't match expected output:\nactual:   %#X\nexpected: %#X\nHexdumps:\n%s", rawInput, expectedRawInput, utils.DiffHex(expectedRawInput, rawInput))
			}
		}
		// If there's a difference, parse the input and display it to simplify debugging
	case StepTypeIncomingPlcMessage:
		typeName := step.payload.Name
		payloadString := step.payload.XMLPretty()

		// Parse the xml into a real model
		t.Log("Parsing model")
		expectedMessage, err := m.parseMessage(typeName, payloadString, step)
		if err != nil {
			return errors.Wrap(err, "error parsing message")
		}
		t.Logf("Parsed message (%T)\n%[1]s", expectedMessage)

		// Serialize the model into bytes
		t.Log("Serializing bytes")
		expectedSerializable, ok := expectedMessage.(utils.Serializable)
		if !ok {
			return errors.New("error converting type into Serializable type")
		}
		var wb utils.WriteBufferByteBased
		if m.byteOrder == binary.BigEndian {
			wb = utils.NewWriteBufferByteBased()
		} else {
			t.Log("using little endian")
			wb = utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
		}
		err = expectedSerializable.SerializeWithWriteBuffer(context.Background(), wb)
		if err != nil {
			return errors.Wrap(err, "error serializing expectedMessage")
		}

		// Send these bytes to the transport
		_bytes := wb.GetBytes()
		t.Logf("Writing to transport\n%s", hex.Dump(_bytes))
		testTransportInstance.FillReadBuffer(_bytes)
	case StepTypeIncomingPlcBytes:
		// Get the raw hex-data.
		t.Log("Get hex data")
		rawInput, err := hex.DecodeString(step.payload.Text)
		if err != nil {
			return errors.Wrap(err, "error decoding hex-encoded byte data: ")
		}

		// Send these bytes to the transport
		t.Logf("Writing bytes to transport\n%[1]s", hex.Dump(rawInput))
		testTransportInstance.FillReadBuffer(rawInput)
	case StepTypeDelay:
		// Get the number of milliseconds
		t.Log("Getting millis")
		delay, err := strconv.Atoi(step.payload.Text)
		if err != nil {
			return errors.Wrap(err, "invalid delay format")
		}
		// Sleep for that long
		t.Logf("Sleeping (delay: %dms)", delay)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	case StepTypeTerminate:
		// Simply close the transport connection
		t.Log("closing transport")
		err := testTransportInstance.Close()
		if err != nil {
			return errors.Wrap(err, "error closing transport")
		}
	}
	m.LogDelimiterSection(t, "┄", " - Finished step: %s after %s ", step.name, time.Since(start))
	return nil
}

func (m DriverTestsuite) parseMessage(typeName string, payloadString string, step DriverTestStep) (any, error) {
	if m.parser == nil {
		return nil, errors.Errorf("Protocol name %s has no mapped parser", m.protocolName)
	}
	parse, err := m.parser.Parse(typeName, payloadString, step.parserArguments...)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing xml")
	}
	return parse, err
}

type DriverTestcase struct {
	name                      string
	steps                     []DriverTestStep
	readRequestResultChannel  <-chan apiModel.PlcReadRequestResult
	writeRequestResultChannel <-chan apiModel.PlcWriteRequestResult
}

type DriverTestStep struct {
	name            string
	stepType        StepType
	parserArguments []string
	payload         xmldom.Node
}

type StepType uint8

//go:generate stringer -type StepType
//go:generate go run ../../tools/plc4xlicenser/gen.go -type=StepType
const (
	StepTypeOutgoingPlcMessage StepType = 0x01
	StepTypeOutgoingPlcBytes   StepType = 0x02
	StepTypeIncomingPlcMessage StepType = 0x03
	StepTypeIncomingPlcBytes   StepType = 0x04
	StepTypeApiRequest         StepType = 0x05
	StepTypeApiResponse        StepType = 0x06
	StepTypeDelay              StepType = 0x07
	StepTypeTerminate          StepType = 0x08
)

func RunDriverTestsuite(t *testing.T, driver plc4go.PlcDriver, testPath string, parser XmlParser, _options ...options.WithOption) {
	t.Log("Extract testsuite options")
	var rootTypeParser func(utils.ReadBufferByteBased) (any, error)
	skippedTestCasesMap := map[string]bool{}
	for _, withOption := range _options {
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
	t.Log("Read the test-specification as XML file")
	// Read the test-specification as XML file
	rootNode := ParseDriverTestsuiteXml(t, testPath)

	t.Log("Parse the contents of the test-specification")
	// Parse the contents of the test-specification
	testsuite := ParseDriverTestsuite(t, *rootNode, parser, rootTypeParser)

	// We don't want to await completion of connection initialization
	if connectionConnectAwaiter, ok := driver.(ConnectionConnectAwaiter); ok {
		t.Log("We don't wait for setup and disconnect")
		connectionConnectAwaiter.SetAwaitSetupComplete(false)
		connectionConnectAwaiter.SetAwaitDisconnectComplete(false)
	}

	t.Log("Initialize the driver manager")
	// Initialize the driver manager
	driverManager := plc4go.NewPlcDriverManager(converter.WithOptionToExternal(_options...)...)
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	transport := test.NewTransport(_options...)
	driverManager.(spi.TransportAware).RegisterTransport(transport)
	driverManager.RegisterDriver(driver)

	t.Logf("Running %d testcases", len(testsuite.testcases))
	for _, testcase := range testsuite.testcases {
		t.Run(testcase.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					testsuite.LogDelimiterSection(t, "=", "Panic Failure\n%+v\n%s", err, debug.Stack())
					t.FailNow()
				}
			}()
			if skippedTestCasesMap[testcase.name] {
				t.Logf("Testcase %s skipped", testcase.name)
				t.Skipf("Testcase %s skipped", testcase.name)
				return
			}
			t.Logf("Running testcase %s", testcase.name)
			if err := testsuite.Run(t, driverManager, testcase); err != nil {
				testsuite.LogDelimiterSection(t, "=", "Failure:\n%s", err)
				t.FailNow()
			}
		})
	}
}

type ConnectionConnectAwaiter interface {
	// SetAwaitSetupComplete sets a flag that the driver should await a connection completion
	SetAwaitSetupComplete(awaitComplete bool)
	// SetAwaitDisconnectComplete sets a flag that the driver should await a dis-connection completion
	SetAwaitDisconnectComplete(awaitComplete bool)
}

func ParseDriverTestsuiteXml(t *testing.T, testPath string) *xmldom.Node {
	// Get the current working directory
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Check if the test-file is available
	pathOffset := "/../../../"
	testSuitePath := path + pathOffset + testPath
	info, err := os.Stat(testSuitePath)
	if os.IsNotExist(err) {
		t.Logf("Test-File %s doesn't exist", testSuitePath)
		t.Fatal(err)
	}
	if info.IsDir() {
		t.Logf("Test-File %s refers to a directory", testSuitePath)
		t.Fatal(err)
	}

	// Open a reader for this file
	dat, err := os.Open(testSuitePath)
	if err != nil {
		t.Log("error opening file")
		t.Fatal(err)
	}

	// Read the xml
	node := xmldom.Must(xmldom.Parse(dat)).Root
	return node
}

func ParseDriverTestsuite(t *testing.T, node xmldom.Node, parser XmlParser, rootTypeParser func(utils.ReadBufferByteBased) (any, error)) *DriverTestsuite {
	if node.Name != "driver-testsuite" {
		t.Fatal("invalid document structure")
	}
	var byteOrder binary.ByteOrder
	if node.GetAttributeValue("byteOrder") != "LITTLE_ENDIAN" {
		byteOrder = binary.BigEndian
	} else {
		t.Log("using little endian")
		byteOrder = binary.LittleEndian
	}
	var testsuiteName string
	var protocolName string
	var outputFlavor string
	var driverName string
	driverParameters := make(map[string]string)
	var setupSteps []DriverTestStep
	var teardownSteps []DriverTestStep
	var testcases []DriverTestcase
	for _, childPtr := range node.Children {
		child := *childPtr
		switch child.Name {
		case "name":
			testsuiteName = child.Text
		case "protocolName":
			protocolName = child.Text
		case "outputFlavor":
			outputFlavor = child.Text
		case "driver-name":
			driverName = child.Text
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
		case "setup":
			steps, err := ParseDriverTestsuiteSteps(t, child)
			if err != nil {
				t.Error("error parsing setup steps")
				t.Fatal(err)
			}
			setupSteps = steps
		case "teardown":
			steps, err := ParseDriverTestsuiteSteps(t, child)
			if err != nil {
				t.Error("error parsing teardown steps")
				t.Fatal(err)
			}
			teardownSteps = steps
		case "testcase":
			testcaseName := child.FindOneByName("name").Text
			stepsNode := child.FindOneByName("steps")
			steps, err := ParseDriverTestsuiteSteps(t, *stepsNode)
			if err != nil {
				t.Errorf("error parsing testcase %s", testcaseName)
				t.Fatal(err)
			}
			testcase := DriverTestcase{
				name:  testcaseName,
				steps: steps,
			}
			testcases = append(testcases, testcase)
		default:
			t.Fatalf("invalid document structure. Unhandled element %s", child.Name)
		}
	}
	t.Logf("Parsed testsuite name: %s, driver name: %s", testsuiteName, driverName)

	return &DriverTestsuite{
		name:             testsuiteName,
		protocolName:     protocolName,
		outputFlavor:     outputFlavor,
		driverName:       driverName,
		driverParameters: driverParameters,
		byteOrder:        byteOrder,
		parser:           parser,
		rootTypeParser:   rootTypeParser,
		setupSteps:       setupSteps,
		teardownSteps:    teardownSteps,
		testcases:        testcases,
	}
}

func ParseDriverTestsuiteSteps(t *testing.T, node xmldom.Node) ([]DriverTestStep, error) {
	t.Logf("Parsing driver testsuite steps (rootElement %s)", node.Name)
	var testSteps []DriverTestStep
	for _, step := range node.Children {
		name := step.GetAttributeValue("name")
		t.Logf("Parsing step (name: %s, rootElement %s)", name, node.Name)
		var stepType StepType
		switch step.Name {
		case "api-request":
			stepType = StepTypeApiRequest
		case "api-response":
			stepType = StepTypeApiResponse
		case "outgoing-plc-message":
			stepType = StepTypeOutgoingPlcMessage
		case "incoming-plc-message":
			stepType = StepTypeIncomingPlcMessage
		case "outgoing-plc-bytes":
			stepType = StepTypeOutgoingPlcBytes
		case "incoming-plc-bytes":
			stepType = StepTypeIncomingPlcBytes
		case "delay":
			stepType = StepTypeDelay
		case "terminate":
			stepType = StepTypeTerminate
		default:
			return nil, errors.Errorf("Unknown step with name %s", step.Name)
		}
		var parserArguments []string
		var payload *xmldom.Node
		t.Logf("Looking for payload (name: %s, rootElement %s)", name, node.Name)
		for _, childNode := range step.Children {
			t.Logf("Found payload candidate (child node name: %s, name: %s, rootElement %s)", childNode.Name, name, node.Name)
			if childNode.Name == "parser-arguments" {
				for _, parserArgumentNode := range childNode.Children {
					parserArguments = append(parserArguments, parserArgumentNode.Text)
				}
			} else if payload == nil {
				payload = childNode
			} else {
				return nil, errors.New("test step can only contain a single payload element")
			}
		}
		if stepType == StepTypeDelay {
			payload = step
		}
		if payload == nil {
			return nil, errors.New("missing payload element")
		}
		testSteps = append(testSteps, DriverTestStep{
			name:            name,
			stepType:        stepType,
			parserArguments: parserArguments,
			payload:         *payload,
		})
	}
	return testSteps, nil
}

func (DriverTestsuite) LogDelimiterSection(t *testing.T, delimiter, format string, a ...any) {
	t.Helper()
	delimiter = strings.Repeat(delimiter, 55)
	if !shouldNoColor() {
		delimiter = fmt.Sprintf("\x1b[%dm%v\x1b[0m", 32, delimiter)
	}
	message := "\n"
	message += delimiter + "\n"
	message += fmt.Sprintf(format, a...) + "\n"
	message += delimiter
	t.Log(message)
}
