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
	adsIO "github.com/apache/plc4x/plc4go/internal/plc4go/ads/readwrite"
	adsModel "github.com/apache/plc4x/plc4go/internal/plc4go/ads/readwrite/model"
	knxIO "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite"
	knxModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	modbusIO "github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite"
	modbusModel "github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite/model"
	s7IO "github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite"
	s7Model "github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports/test"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	api "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/subchen/go-xmldom"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"
)

type DriverTestsuite struct {
	name             string
	protocolName     string
	outputFlavor     string
	driverName       string
	driverParameters map[string]string
	setupSteps       []TestStep
	teardownSteps    []TestStep
	testcases        []Testcase
}

func (m DriverTestsuite) Run(driverManager plc4go.PlcDriverManager, testcase Testcase) error {
	var options []string
	for key, value := range m.driverParameters {
		options = append(options, fmt.Sprintf("%s=%s", key, value))
	}
	optionsString := ""
	if len(options) > 0 {
		optionsString = "?" + strings.Join(options, "&")
	}
	// Get a connection
	connectionChan := driverManager.GetConnection(m.driverName + ":test://hurz" + optionsString)
	connectionResult := <-connectionChan

	if connectionResult.Err != nil {
		return errors.Wrap(connectionResult.Err, "error getting a connection")
	}

	log.Info().Msgf("\n-------------------------------------------------------\nExecuting testcase: %s \n-------------------------------------------------------\n", testcase.name)

	// Run the setup steps
	log.Info().Msgf("\n-------------------------------------------------------\nPerforming setup for: %s \n-------------------------------------------------------\n", testcase.name)
	for _, testStep := range m.setupSteps {
		err := m.ExecuteStep(connectionResult.Connection, &testcase, testStep)
		if err != nil {
			return errors.Wrap(err, "error in setup step "+testStep.name)
		}
	}

	// Run the actual scenario steps
	log.Info().Msgf("\n-------------------------------------------------------\nRunning testcases for: %s \n-------------------------------------------------------\n", testcase.name)
	for _, testStep := range testcase.steps {
		err := m.ExecuteStep(connectionResult.Connection, &testcase, testStep)
		if err != nil {
			return errors.Wrap(err, "error in step "+testStep.name)
		}
	}

	// Run the teardown steps
	log.Info().Msgf("\n-------------------------------------------------------\nPerforming teardown for: %s \n-------------------------------------------------------\n", testcase.name)
	for _, testStep := range m.teardownSteps {
		err := m.ExecuteStep(connectionResult.Connection, &testcase, testStep)
		if err != nil {
			return errors.Wrap(err, "error in teardown step "+testStep.name)
		}
	}

	log.Info().Msgf("\n-------------------------------------------------------\nDone\n-------------------------------------------------------\n")
	return nil
}

func (m DriverTestsuite) ExecuteStep(connection plc4go.PlcConnection, testcase *Testcase, step TestStep) error {
	mc, ok := connection.(spi.TransportInstanceExposer)
	if !ok {
		return errors.New("couldn't access connections transport instance")
	}
	testTransportInstance, ok := mc.GetTransportInstance().(transports.TestTransportInstance)
	if !ok {
		return errors.New("transport must be of type TestTransport")
	}

	start := time.Now()
	log.Info().Msgf("\n-------------------------------------------------------\n - Executing step: %s \n-------------------------------------------------------\n", step.name)

	log.Debug().Stringer("stepType", step.stepType).Msg("Handling step")
	switch step.stepType {
	case StepTypeApiRequest:
		switch step.payload.Name {
		case "TestReadRequest":
			// Assemble a read-request according to the information in the test xml
			log.Trace().Msg("Assemble read request")
			rrb := connection.ReadRequestBuilder()
			for _, fieldNode := range step.payload.GetChild("fields").GetChildren("field") {
				fieldName := fieldNode.GetChild("name").Text
				fieldAddress := fieldNode.GetChild("address").Text
				rrb.AddQuery(fieldName, fieldAddress)
			}
			readRequest, err := rrb.Build()
			if err != nil {
				return errors.Wrap(err, "Error creating read-request")
			}

			// Execute the read-request and store the response-channel in the testcase.
			log.Trace().Msg("Execute read request")
			if testcase.readRequestResultChannel != nil {
				return errors.New("testcase read-request result channel already occupied")
			}
			testcase.readRequestResultChannel = readRequest.Execute()
		case "TestWriteRequest":
			log.Trace().Msg("Assemble write request")
			wrb := connection.WriteRequestBuilder()
			for _, fieldNode := range step.payload.GetChild("fields").GetChildren("field") {
				fieldName := fieldNode.GetChild("name").Text
				fieldAddress := fieldNode.GetChild("address").Text

				he, ok := connection.(spi.HandlerExposer)
				if !ok {
					return errors.New("connection is not a HandlerExposer")
				}
				field, err := he.GetPlcFieldHandler().ParseQuery(fieldAddress)
				if err != nil {
					return errors.Wrapf(err, "error parsing address: %s", fieldAddress)
				}
				if field.GetQuantity() > 1 {
					var fieldValue []string
					for _, valueChild := range fieldNode.GetChildren("value") {
						fieldValue = append(fieldValue, valueChild.Text)
					}
					wrb.AddQuery(fieldName, fieldAddress, fieldValue)
				} else {
					fieldValue := fieldNode.GetChild("value").Text
					wrb.AddQuery(fieldName, fieldAddress, fieldValue)
				}
			}
			writeRequest, err := wrb.Build()
			if err != nil {
				return errors.Wrap(err, "Error creating write-request")
			}
			log.Trace().Msg("Execute write request")
			if testcase.writeRequestResultChannel != nil {
				return errors.New("testcase write-request result channel already occupied")
			}
			testcase.writeRequestResultChannel = writeRequest.Execute()
		}
	case StepTypeApiResponse:
		switch step.payload.Name {
		case "PlcReadResponse":
			if testcase.readRequestResultChannel == nil {
				return errors.New("no read response expected")
			}
			log.Trace().Msg("Waiting for read request result")
			readRequestResult := <-testcase.readRequestResultChannel
			if readRequestResult.Err != nil {
				return errors.Wrap(readRequestResult.Err, "error sending response")
			}
			// Serialize the response to XML
			actualResponse, err := xml.Marshal(readRequestResult.Response)
			if err != nil {
				return errors.Wrap(err, "error serializing response")
			}
			// Get the reference XML
			referenceSerialized := step.payload.XML()
			// Compare the results
			err = CompareResults(actualResponse, []byte(referenceSerialized))
			if err != nil {
				return errors.Wrap(err, "Error comparing the results")
			}
			// Reset read channel
			testcase.readRequestResultChannel = nil
		case "PlcWriteResponse":
			if testcase.writeRequestResultChannel == nil {
				return errors.New("no write response expected")
			}
			log.Trace().Msg("Waiting for write request result")
			writeResponseResult := <-testcase.writeRequestResultChannel
			if writeResponseResult.Err != nil {
				return errors.Wrap(writeResponseResult.Err, "error sending response")
			}
			// Serialize the response to XML
			actualResponse, err := xml.Marshal(writeResponseResult.Response)
			if err != nil {
				return errors.Wrap(err, "error serializing response")
			}
			// Get the reference XML
			referenceSerialized := step.payload.XML()
			// Compare the results
			err = CompareResults(actualResponse, []byte(referenceSerialized))
			if err != nil {
				return errors.Wrap(err, "Error comparing the results")
			}
			// Reset write channel
			testcase.writeRequestResultChannel = nil
		}
	case StepTypeOutgoingPlcMessage:
		typeName := step.payload.Name
		payloadString := step.payload.XML()

		// Parse the xml into a real model
		log.Trace().Msg("parsing xml")
		expectedMessage, err := parseMessage(m.protocolName, typeName, payloadString, step)

		// Serialize the model into bytes
		log.Trace().Msg("Write to bytes")
		expectedSerializable, ok := expectedMessage.(utils.Serializable)
		if !ok {
			return errors.New("error converting type into Serializable type")
		}
		var expectedWriteBuffer utils.WriteBufferByteBased
		switch m.driverName {
		case "ads":
			expectedWriteBuffer = utils.NewLittleEndianWriteBuffer()
		default:
			expectedWriteBuffer = utils.NewWriteBuffer()
		}
		err = expectedSerializable.Serialize(expectedWriteBuffer)
		if err != nil {
			return errors.Wrap(err, "error serializing expectedMessage")
		}
		expectedRawOutput := expectedWriteBuffer.GetBytes()
		expectedRawOutputLength := uint32(len(expectedRawOutput))

		now := time.Now()
		// Read exactly this amount of bytes from the transport
		log.Trace().Uint32("expectedRawOutputLength", expectedRawOutputLength).Msg("Reading bytes")
		for testTransportInstance.GetNumDrainableBytes() < expectedRawOutputLength {
			if time.Now().Sub(now) > 2*time.Second {
				return errors.Errorf("error getting bytes from transport. Not enough data available: actual(%d)<expected(%d)", testTransportInstance.GetNumDrainableBytes(), expectedRawOutputLength)
			}
			time.Sleep(10 * time.Millisecond)
		}
		actualRawOutput, err := testTransportInstance.DrainWriteBuffer(expectedRawOutputLength)
		if testTransportInstance.GetNumDrainableBytes() != 0 {
			//panic(fmt.Sprintf("leftover drainable bytes (%d)", testTransportInstance.GetNumDrainableBytes()))
		}
		if err != nil {
			return errors.Wrap(err, "error getting bytes from transport")
		}

		// Compare the bytes read with the ones we expect
		log.Trace().Msg("Comparing outputs")
		for i := range expectedRawOutput {
			if expectedRawOutput[i] != actualRawOutput[i] {
				switch m.driverName {
				case "modbus":
					expectation := expectedSerializable.(*modbusModel.ModbusTcpADU)
					actual, err := modbusModel.ModbusTcpADUParse(utils.NewReadBuffer(actualRawOutput), false)
					log.Error().Err(err).Msgf("A readabled render of expectation:\n%v\nvs actual paket\n%v\n", expectation, actual)
				case "ads":
					expectation := expectedSerializable.(*adsModel.AmsTCPPacket)
					actual, err := adsModel.AmsTCPPacketParse(utils.NewLittleEndianReadBuffer(actualRawOutput))
					log.Error().Err(err).Msgf("A readabled render of expectation:\n%v\nvs actual paket\n%v\n", expectation, actual)
				case "s7":
					expectation := expectedSerializable.(*s7Model.TPKTPacket)
					actual, err := s7Model.TPKTPacketParse(utils.NewReadBuffer(actualRawOutput))
					log.Error().Err(err).Msgf("A readabled render of expectation:\n%v\nvs actual paket\n%v\n", expectation, actual)
				case "knx":
					expectation := expectedSerializable.(*knxModel.KnxNetIpMessage)
					actual, err := knxModel.KnxNetIpMessageParse(utils.NewReadBuffer(actualRawOutput))
					log.Error().Err(err).Msgf("A readabled render of expectation:\n%v\nvs actual paket\n%v\n", expectation, actual)
				}
				return errors.Errorf("actual output doesn't match expected output:\nactual:   0x%X\nexpected: 0x%X", actualRawOutput, expectedRawOutput)
			}
		}
		// If there's a difference, parse the input and display it to simplify debugging
	case StepTypeOutgoingPlcBytes:
		// Read exactly this amount of bytes from the transport
		log.Trace().Msg("Reading bytes")
		expectedRawInput, err := hex.DecodeString(step.payload.Text)
		if err != nil {
			return errors.Wrap(err, "error decoding hex-encoded byte data")
		}
		rawInput, err := testTransportInstance.DrainWriteBuffer(uint32(len(expectedRawInput)))
		if err != nil {
			return errors.Wrap(err, "error getting bytes from transport")
		}

		// Compare the bytes read with the ones we expect
		log.Trace().Msg("Comparing bytes")
		for i := range expectedRawInput {
			if expectedRawInput[i] != rawInput[i] {
				return errors.Errorf("actual output doesn't match expected output:\nactual:   0x%X\nexpected: 0x%X", rawInput, expectedRawInput)
			}
		}
		// If there's a difference, parse the input and display it to simplify debugging
	case StepTypeIncomingPlcMessage:
		typeName := step.payload.Name
		payloadString := step.payload.XML()

		// Parse the xml into a real model
		log.Trace().Msg("Parsing model")
		expectedMessage, err := parseMessage(m.protocolName, typeName, payloadString, step)
		if err != nil {
			return errors.Wrap(err, "error parsing message")
		}

		// Serialize the model into bytes
		log.Trace().Msg("Serializing bytes")
		expectedSerializable, ok := expectedMessage.(utils.Serializable)
		if !ok {
			return errors.New("error converting type into Serializable type")
		}
		var wb utils.WriteBufferByteBased
		switch m.driverName {
		case "ads":
			wb = utils.NewLittleEndianWriteBuffer()
		default:
			wb = utils.NewWriteBuffer()
		}
		err = expectedSerializable.Serialize(wb)
		if err != nil {
			return errors.Wrap(err, "error serializing expectedMessage")
		}

		// Send these bytes to the transport
		log.Trace().Msg("Writing to transport")
		err = testTransportInstance.FillReadBuffer(wb.GetBytes())
		if err != nil {
			return errors.Wrap(err, "error writing data to transport")
		}
	case StepTypeIncomingPlcBytes:
		// Get the raw hex-data.
		log.Trace().Msg("Get hex data")
		rawInput, err := hex.DecodeString(step.payload.Text)
		if err != nil {
			return errors.Wrap(err, "error decoding hex-encoded byte data: ")
		}

		// Send these bytes to the transport
		log.Trace().Msg("Writing bytes to transport")
		err = testTransportInstance.FillReadBuffer(rawInput)
		if err != nil {
			return errors.Wrap(err, "error writing data to transport")
		}
	case StepTypeDelay:
		// Get the number of milliseconds
		log.Trace().Msg("Getting millis")
		delay, err := strconv.Atoi(step.payload.Text)
		if err != nil {
			return errors.Wrap(err, "invalid delay format")
		}
		// Sleep for that long
		log.Debug().Int("delay", delay).Msg("Sleeping")
		time.Sleep(time.Millisecond * time.Duration(delay))
	case StepTypeTerminate:
		// Simply close the transport connection
		log.Trace().Msg("closing transport")
		err := testTransportInstance.Close()
		if err != nil {
			return errors.Wrap(err, "error closing transport")
		}
	}
	log.Info().Msgf("\n\n-------------------------------------------------------\n - Finished step: %s after %vms \n-------------------------------------------------------", step.name, time.Now().Sub(start).Milliseconds())
	return nil
}

func parseMessage(protocolName string, typeName string, payloadString string, step TestStep) (interface{}, error) {
	type Parser interface {
		Parse(typeName string, xmlString string, parserArguments ...string) (interface{}, error)
	}
	parserMap := map[string]Parser{
		"modbus":   modbusIO.ModbusXmlParserHelper{},
		"ads":      adsIO.AdsXmlParserHelper{},
		"knxnetip": knxIO.KnxnetipXmlParserHelper{},
		"s7":       s7IO.S7XmlParserHelper{},
	}
	if parser, ok := parserMap[protocolName]; ok {
		expected, err := parser.Parse(typeName, payloadString, step.parserArguments...)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing xml")
		}
		return expected, nil
	} else {
		return nil, errors.Errorf("Protocol name %s has no mapped parser", protocolName)
	}
}

func (m DriverTestsuite) ParseXml(referenceXml *xmldom.Node, parserArguments []string) {
	normalizeXml(referenceXml)
	//referenceSerialized := referenceXml.FirstChild().XML()
}

type Testcase struct {
	name                      string
	steps                     []TestStep
	readRequestResultChannel  <-chan api.PlcReadRequestResult
	writeRequestResultChannel <-chan api.PlcWriteRequestResult
}

type TestStep struct {
	name            string
	stepType        StepType
	parserArguments []string
	payload         xmldom.Node
}

type StepType uint8

//go:generate stringer -type StepType
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

func RunDriverTestsuite(t *testing.T, driver plc4go.PlcDriver, testPath string, skippedTestCases ...string) {
	skippedTestCasesMap := map[string]bool{}
	for _, skippedTestCase := range skippedTestCases {
		skippedTestCasesMap[skippedTestCase] = true
	}
	// Read the test-specification as XML file
	rootNode, err := ParseDriverTestsuiteXml(testPath)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}

	// Parse the contents of the test-specification
	testsuite, err := ParseDriverTestsuite(*rootNode)
	if err != nil {
		// TODO: zerolog doesn't render stack human readable :(
		fmt.Printf("%+v\n", err)
		log.Error().
			Stack().
			Err(err).
			Msg("Failed to parse test-specification")
		t.Error(err)
		t.Fail()
		return
	}

	// We don't want to await completion of connection initalization
	if connectionConnectAwaiter, ok := driver.(ConnectionConnectAwaiter); ok {
		connectionConnectAwaiter.SetAwaitSetupComplete(false)
		connectionConnectAwaiter.SetAwaitDisconnectComplete(false)
	}

	// Initialize the driver manager
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterTransport(test.NewTransport())
	driverManager.RegisterDriver(driver)

	for _, testcase := range testsuite.testcases {
		t.Run(testcase.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					log.Error().Msgf("\n-------------------------------------------------------\nPanic Failure\n%+v\n%s\n-------------------------------------------------------\n\n", err, debug.Stack())
					t.FailNow()
				}
			}()
			if skippedTestCasesMap[testcase.name] {
				log.Warn().Msgf("Testcase %s skipped", testcase.name)
				t.Skipf("Testcase %s skipped", testcase.name)
				return
			}
			log.Info().Msgf("Running testcase %s", testcase.name)
			if err := testsuite.Run(driverManager, testcase); err != nil {
				log.Error().Err(err).Msgf("\n-------------------------------------------------------\nFailure\n%+v\n-------------------------------------------------------\n\n", err)
				t.Fail()
			}
		})
	}
	// Execute the tests in the testsuite
	log.Info().Msgf(testsuite.name)
}

type ConnectionConnectAwaiter interface {
	SetAwaitSetupComplete(awaitComplete bool)
	SetAwaitDisconnectComplete(awaitComplete bool)
}

func ParseDriverTestsuiteXml(testPath string) (*xmldom.Node, error) {
	// Get the current working directory
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Check if the test-file is available
	info, err := os.Stat(path + "/../../../../" + testPath)
	if os.IsNotExist(err) {
		return nil, errors.Wrap(err, "test-File doesn't exist")
	}
	if info.IsDir() {
		return nil, errors.New("test-file refers to a directory")
	}

	// Open a reader for this file
	dat, err := os.Open(path + "/../../../../" + testPath)
	if err != nil {
		return nil, errors.Wrap(err, "error opening file")
	}

	// Read the xml
	node := xmldom.Must(xmldom.Parse(dat)).Root
	return node, nil
}

func ParseDriverTestsuite(node xmldom.Node) (*DriverTestsuite, error) {
	if node.Name != "driver-testsuite" {
		return nil, errors.New("invalid document structure")
	}
	var testsuiteName string
	var protocolName string
	var outputFlavor string
	var driverName string
	driverParameters := make(map[string]string)
	var setupSteps []TestStep
	var teardownSteps []TestStep
	var testcases []Testcase
	for _, childPtr := range node.Children {
		child := *childPtr
		if child.Name == "name" {
			testsuiteName = child.Text
		} else if child.Name == "protocolName" {
			protocolName = child.Text
		} else if child.Name == "outputFlavor" {
			outputFlavor = child.Text
		} else if child.Name == "driver-name" {
			driverName = child.Text
		} else if child.Name == "driver-parameters" {
			parameterList := child.FindByName("parameter")
			for _, parameter := range parameterList {
				nameElement := parameter.FindOneByName("name")
				valueElement := parameter.FindOneByName("value")
				if nameElement == nil || valueElement == nil {
					return nil, errors.New("invalid parameter found: no present")
				}
				name := nameElement.Text
				value := valueElement.Text
				if name == "" || value == "" {
					return nil, errors.New("invalid parameter found: empty")
				}
				driverParameters[name] = value
			}
		} else if child.Name == "setup" {
			steps, err := ParseDriverTestsuiteSteps(child)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing setup steps")
			}
			setupSteps = steps
		} else if child.Name == "teardown" {
			steps, err := ParseDriverTestsuiteSteps(child)
			if err != nil {
				return nil, errors.Wrap(err, "error teardown setup steps")
			}
			teardownSteps = steps
		} else if child.Name == "testcase" {
			testcaseName := child.FindOneByName("name").Text
			stepsNode := child.FindOneByName("steps")
			steps, err := ParseDriverTestsuiteSteps(*stepsNode)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing testcase "+testcaseName)
			}
			testcase := Testcase{
				name:  testcaseName,
				steps: steps,
			}
			testcases = append(testcases, testcase)
		} else {
			return nil, errors.New("invalid document structure. Unhandled element " + child.Name)
		}
	}
	log.Info().
		Str("testsuite name", testsuiteName).
		Str("driver name", driverName).
		Msgf("Parsed test suite %s", testsuiteName)

	return &DriverTestsuite{
		name:             testsuiteName,
		protocolName:     protocolName,
		outputFlavor:     outputFlavor,
		driverName:       driverName,
		driverParameters: driverParameters,
		setupSteps:       setupSteps,
		teardownSteps:    teardownSteps,
		testcases:        testcases,
	}, nil
}

func ParseDriverTestsuiteSteps(node xmldom.Node) ([]TestStep, error) {
	log.Debug().Str("rootElement", node.Name).Msg("Parsing driver testsuite steps")
	var testSteps []TestStep
	for _, step := range node.Children {
		name := step.GetAttributeValue("name")
		log.Debug().Str("rootElement", node.Name).Str("name", name).Msg("Parsing step")
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
		log.Debug().Str("rootElement", node.Name).Msg("Looking for payload")
		for _, childNode := range step.Children {
			log.Debug().Str("child node name", childNode.Name).Str("rootElement", node.Name).Msg("Found payload candidate")
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
		testSteps = append(testSteps, TestStep{
			name:            name,
			stepType:        stepType,
			parserArguments: parserArguments,
			payload:         *payload,
		})
	}
	return testSteps, nil
}
