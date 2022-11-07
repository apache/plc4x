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
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api"
	api "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/subchen/go-xmldom"
)

type DriverTestsuite struct {
	name             string
	protocolName     string
	outputFlavor     string
	driverName       string
	driverParameters map[string]string
	byteOrder        binary.ByteOrder
	parser           XmlParser
	rootTypeParser   func(utils.ReadBufferByteBased) (interface{}, error)
	setupSteps       []TestStep
	teardownSteps    []TestStep
	testcases        []Testcase
}

type XmlParser interface {
	Parse(typeName string, xmlString string, parserArguments ...string) (interface{}, error)
}

type WithOption interface {
	isOption() bool
}

type option struct {
}

func (_ option) isOption() bool {
	return true
}

// WithRootTypeParser Can be used to output the root type of a protocol for better debugging
func WithRootTypeParser(rootTypeParser func(utils.ReadBufferByteBased) (interface{}, error)) WithOption {
	return withRootTypeParser{rootTypeParser: rootTypeParser}
}

type TestTransportInstance interface {
	transports.TransportInstance
	FillReadBuffer(data []uint8) error
	GetNumDrainableBytes() uint32
	DrainWriteBuffer(numBytes uint32) ([]uint8, error)
}

type withRootTypeParser struct {
	option
	rootTypeParser func(utils.ReadBufferByteBased) (interface{}, error)
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

	if connectionResult.GetConnection() != nil {
		return errors.Wrap(connectionResult.GetErr(), "error getting a connection")
	}

	log.Info().Msgf("\n-------------------------------------------------------\nExecuting testcase: %s \n-------------------------------------------------------\n", testcase.name)

	// Run the setup steps
	log.Info().Msgf("\n-------------------------------------------------------\nPerforming setup for: %s \n-------------------------------------------------------\n", testcase.name)
	for _, testStep := range m.setupSteps {
		err := m.ExecuteStep(connectionResult.GetConnection(), &testcase, testStep)
		if err != nil {
			return errors.Wrap(err, "error in setup step "+testStep.name)
		}
		// We sleep a bit to not run too fast into the post setup steps and give connections a bit time to settle built up
		time.Sleep(time.Second)
	}

	// Run the actual scenario steps
	log.Info().Msgf("\n-------------------------------------------------------\nRunning testcases for: %s \n-------------------------------------------------------\n", testcase.name)
	for _, testStep := range testcase.steps {
		err := m.ExecuteStep(connectionResult.GetConnection(), &testcase, testStep)
		if err != nil {
			return errors.Wrap(err, "error in step "+testStep.name)
		}
	}

	// Run the teardown steps
	log.Info().Msgf("\n-------------------------------------------------------\nPerforming teardown for: %s \n-------------------------------------------------------\n", testcase.name)
	for _, testStep := range m.teardownSteps {
		err := m.ExecuteStep(connectionResult.GetConnection(), &testcase, testStep)
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
	testTransportInstance, ok := mc.GetTransportInstance().(TestTransportInstance)
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
				rrb.AddFieldQuery(fieldName, fieldAddress)
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
				field, err := he.GetPlcFieldHandler().ParseField(fieldAddress)
				if err != nil {
					return errors.Wrapf(err, "error parsing address: %s", fieldAddress)
				}
				if len(field.GetArrayInfo()) > 0 {
					var fieldValue []string
					for _, valueChild := range fieldNode.GetChildren("value") {
						fieldValue = append(fieldValue, valueChild.Text)
					}
					wrb.AddFieldQuery(fieldName, fieldAddress, fieldValue)
				} else {
					fieldValue := fieldNode.GetChild("value").Text
					wrb.AddFieldQuery(fieldName, fieldAddress, fieldValue)
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
			if readRequestResult.GetErr() != nil {
				return errors.Wrap(readRequestResult.GetErr(), "error sending response")
			}
			// Serialize the response to XML
			xmlWriteBuffer := utils.NewXmlWriteBuffer()
			err := readRequestResult.GetResponse().(utils.Serializable).SerializeWithWriteBuffer(xmlWriteBuffer)
			if err != nil {
				return errors.Wrap(err, "error serializing response")
			}
			actualResponse := xmlWriteBuffer.GetXmlString()
			// Get the reference XML
			referenceSerialized := step.payload.XMLPretty()
			// Compare the results
			err = CompareResults([]byte(actualResponse), []byte(referenceSerialized))
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
			if writeResponseResult.GetErr() != nil {
				return errors.Wrap(writeResponseResult.GetErr(), "error sending response")
			}
			// Serialize the response to XML
			xmlWriteBuffer := utils.NewXmlWriteBuffer()
			err := writeResponseResult.GetResponse().(utils.Serializable).SerializeWithWriteBuffer(xmlWriteBuffer)
			if err != nil {
				return errors.Wrap(err, "error serializing response")
			}
			actualResponse := xmlWriteBuffer.GetXmlString()
			// Get the reference XML
			referenceSerialized := step.payload.XMLPretty()
			// Compare the results
			err = CompareResults([]byte(actualResponse), []byte(referenceSerialized))
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
		log.Trace().Msg("parsing xml")
		expectedMessage, err := m.parseMessage(typeName, payloadString, step)
		if err != nil {
			return errors.Wrap(err, "Error parsing message")
		}

		// Serialize the model into bytes
		log.Trace().Msg("Write to bytes")
		expectedSerializable, ok := expectedMessage.(utils.Serializable)
		if !ok {
			return errors.Errorf("error converting type %t into Serializable type", expectedMessage)
		}
		var expectedWriteBuffer utils.WriteBufferByteBased
		if m.byteOrder == binary.BigEndian {
			expectedWriteBuffer = utils.NewWriteBufferByteBased()
		} else {
			expectedWriteBuffer = utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
		}
		err = expectedSerializable.SerializeWithWriteBuffer(expectedWriteBuffer)
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

		var bufferFactory func([]byte, ...utils.ReadBufferByteBasedOptions) utils.ReadBufferByteBased
		if m.byteOrder == binary.BigEndian {
			bufferFactory = utils.NewReadBufferByteBased
		} else {
			bufferFactory = func(bytes []byte, options ...utils.ReadBufferByteBasedOptions) utils.ReadBufferByteBased {
				return utils.NewReadBufferByteBased(bytes, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
			}
		}
		// Compare the bytes read with the ones we expect
		log.Trace().Msg("Comparing outputs")
		for i := range expectedRawOutput {
			if expectedRawOutput[i] != actualRawOutput[i] {
				if m.rootTypeParser != nil {
					readBufferByteBased := bufferFactory(actualRawOutput)
					actual, err := m.rootTypeParser(readBufferByteBased)
					log.Error().Err(err).Msgf("A readable render of expectation:\n%v\nvs actual paket\n%v\n", expectedSerializable, actual)
				}
				return errors.Errorf("actual output doesn't match expected output:\nactual:\n%s\nexpected:\n%s", utils.Dump(actualRawOutput), utils.Dump(expectedRawOutput))
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
		payloadString := step.payload.XMLPretty()

		// Parse the xml into a real model
		log.Trace().Msg("Parsing model")
		expectedMessage, err := m.parseMessage(typeName, payloadString, step)
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
		if m.byteOrder == binary.BigEndian {
			wb = utils.NewWriteBufferByteBased()
		} else {
			wb = utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
		}
		err = expectedSerializable.SerializeWithWriteBuffer(wb)
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

func (m DriverTestsuite) parseMessage(typeName string, payloadString string, step TestStep) (interface{}, error) {
	if m.parser == nil {
		return nil, errors.Errorf("Protocol name %s has no mapped parser", m.protocolName)
	}
	parse, err := m.parser.Parse(typeName, payloadString, step.parserArguments...)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing xml")
	}
	return parse, err
}

func (m DriverTestsuite) ParseXml(referenceXml *xmldom.Node, parserArguments []string) {
	normalizeXml(referenceXml)
	//referenceSerialized := referenceXml.FirstChild().XMLPretty()
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

func RunDriverTestsuite(t *testing.T, driver plc4go.PlcDriver, testPath string, parser XmlParser, skippedTestCases ...string) {
	RunDriverTestsuiteWithOptions(t, driver, testPath, parser, nil, skippedTestCases...)
}

func RunDriverTestsuiteWithOptions(t *testing.T, driver plc4go.PlcDriver, testPath string, parser XmlParser, options []WithOption, skippedTestCases ...string) {
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

	var rootTypeParser func(utils.ReadBufferByteBased) (interface{}, error)
	for _, withOption := range options {
		switch withOption.(type) {
		case withRootTypeParser:
			log.Info().Msg("Using root type parser for better output")
			rootTypeParser = withOption.(withRootTypeParser).rootTypeParser
		}
	}
	// Parse the contents of the test-specification
	testsuite, err := ParseDriverTestsuite(*rootNode, parser, rootTypeParser)
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

	// We don't want to await completion of connection initialization
	if connectionConnectAwaiter, ok := driver.(ConnectionConnectAwaiter); ok {
		connectionConnectAwaiter.SetAwaitSetupComplete(false)
		connectionConnectAwaiter.SetAwaitDisconnectComplete(false)
	}

	// Initialize the driver manager
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.(spi.TransportAware).RegisterTransport(test.NewTransport())
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
	pathOffset := "/../../../"
	testSuitePath := path + pathOffset + testPath
	info, err := os.Stat(testSuitePath)
	if os.IsNotExist(err) {
		return nil, errors.Wrap(err, "test-File doesn't exist")
	}
	if info.IsDir() {
		return nil, errors.New("test-file refers to a directory")
	}

	// Open a reader for this file
	dat, err := os.Open(testSuitePath)
	if err != nil {
		return nil, errors.Wrap(err, "error opening file")
	}

	// Read the xml
	node := xmldom.Must(xmldom.Parse(dat)).Root
	return node, nil
}

func ParseDriverTestsuite(node xmldom.Node, parser XmlParser, rootTypeParser func(utils.ReadBufferByteBased) (interface{}, error)) (*DriverTestsuite, error) {
	if node.Name != "driver-testsuite" {
		return nil, errors.New("invalid document structure")
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
		byteOrder:        byteOrder,
		parser:           parser,
		rootTypeParser:   rootTypeParser,
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
