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
package drivers

import (
	"encoding/hex"
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip"
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports/udp"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func TestParser(t *testing.T) {
	request, err := hex.DecodeString("06100202004e080100000000000036010200ff00000000c501010e37e000170c00246d00ec7830302d32342d36442d30302d45432d3738000000000000000000000000000a020201030104010501")
	if err != nil {
		// Output an error ...
	}
	rb := utils.NewReadBuffer(request)
	knxMessage, err := driverModel.KnxNetIpMessageParse(rb)
	if err != nil {
		fmt.Printf("Got error parsing message: %s\n", err.Error())
		// TODO: Possibly clean up ...
		return
	}
	print(knxMessage)
}

func TestKnxNetIpPlc4goBrowse(t *testing.T) {
	Init()

	startTime := time.Now()

	log.Debug("Initializing PLC4X")
	driverManager := plc4go.NewPlcDriverManager()
	log.Debug("Registering KNXnet/IP driver")
	driverManager.RegisterDriver(knxnetip.NewKnxNetIpDriver())
	log.Debug("Registering UDP transport")
	driverManager.RegisterTransport(udp.NewUdpTransport())

	// Create a connection string from the discovery result.
	connectionString := "knxnet-ip:udp://192.168.42.11:3671"
	crc := driverManager.GetConnection(connectionString)
	connectionResult := <-crc
	if connectionResult.Err != nil {
		log.Errorf("Got an error getting a connection: %s", connectionResult.Err.Error())
		t.Fail()
		return
	}
	log.Info("Got a connection")
	connection := connectionResult.Connection
	defer connection.Close()

	// Build a browse request, to scan the KNX network for KNX devices
	// (Limiting the range to only the actually used range of addresses)
	browseRequestBuilder := connection.BrowseRequestBuilder()
	browseRequestBuilder.AddItem("findAllKnxDevices", "[1-3].[1-6].[1-60]")
	browseRequest, err := browseRequestBuilder.Build()
	if err != nil {
		log.Errorf("Got an error preparing browse-request: %s", connectionResult.Err.Error())
		t.Fail()
		return
	}

	// Execute the browse-request
	log.Info("Scanning for KNX devices")
	brr := browseRequest.Execute()
	browseRequestResults := <-brr
	if browseRequestResults.Err != nil {
		log.Errorf("Got an error scanning for KNX devices: %s", connectionResult.Err.Error())
		t.Fail()
		return
	}

	// Output the addresses found
	for _, queryName := range browseRequestResults.Response.GetQueryNames() {
		results := browseRequestResults.Response.GetQueryResults(queryName)
		for _, result := range results {
			log.Infof("Found KNX device at address: %v querying device information: \n", result.Address)

			// Create a read-request to read the manufacturer and hardware ids
			readRequestBuilder := connection.ReadRequestBuilder()
			readRequestBuilder.AddItem("manufacturerId", result.Address+"/0/12")
			readRequestBuilder.AddItem("applicationProgramVersion", result.Address+"/3/13")
			readRequestBuilder.AddItem("interfaceProgramVersion", result.Address+"/4/13")
			readRequest, err := readRequestBuilder.Build()
			if err != nil {
				log.Errorf("Got an error creating read-request: %s", err.Error())
				t.Fail()
				return
			}

			// Execute the read-requests
			rrr := readRequest.Execute()
			readResult := <-rrr
			if readResult.Err != nil {
				log.Errorf("got an error executing read-request: %s", readResult.Err.Error())
				t.Fail()
				return
			}

			// Check the response
			readResponse := readResult.Response
			if readResponse.GetResponseCode("manufacturerId") != apiModel.PlcResponseCode_OK {
				log.Errorf("Got an error response code %d for field 'manufacturerId'", readResponse.GetResponseCode("manufacturerId"))
				t.Fail()
				continue
			}
			if readResponse.GetResponseCode("applicationProgramVersion") != apiModel.PlcResponseCode_OK && readResponse.GetResponseCode("interfaceProgramVersion") != apiModel.PlcResponseCode_OK {
				log.Errorf("Got response code %d for address %s ('programVersion')",
					readResponse.GetResponseCode("applicationProgramVersion"), result.Address+"/3/13")
				log.Errorf("Got response code %d for address %s ('programVersion')",
					readResponse.GetResponseCode("interfaceProgramVersion"), result.Address+"/4/13")
				t.Fail()
			}

			manufacturerId := readResponse.GetValue("manufacturerId").GetUint16()
			if readResponse.GetResponseCode("applicationProgramVersion") == apiModel.PlcResponseCode_OK {
				programVersion := readResponse.GetValue("applicationProgramVersion")
				programVersionBytes := PlcValueUint8ListToByteArray(programVersion)
				log.Infof(" - Manufacturer Id: %d, Application Program Version: %s\n", manufacturerId, hex.EncodeToString(programVersionBytes))
			} else if readResponse.GetResponseCode("interfaceProgramVersion") == apiModel.PlcResponseCode_OK {
				programVersion := readResponse.GetValue("interfaceProgramVersion")
				programVersionBytes := PlcValueUint8ListToByteArray(programVersion)
				log.Infof(" - Manufacturer Id: %d, Interface Program Version: %s\n", manufacturerId, hex.EncodeToString(programVersionBytes))
			}
		}
	}

	log.Infof("Operation finished in %s", time.Since(startTime))
}

func TestKnxNetIpPlc4goBlockingBrowseWithCallback(t *testing.T) {
	Init()

	startTime := time.Now()

	log.Debug("Initializing PLC4X")
	driverManager := plc4go.NewPlcDriverManager()
	log.Debug("Registering KNXnet/IP driver")
	driverManager.RegisterDriver(knxnetip.NewKnxNetIpDriver())
	log.Debug("Registering UDP transport")
	driverManager.RegisterTransport(udp.NewUdpTransport())

	// Create a connection string from the discovery result.
	connectionString := "knxnet-ip:udp://192.168.42.11:3671"
	crc := driverManager.GetConnection(connectionString)
	connectionResult := <-crc
	if connectionResult.Err != nil {
		log.Errorf("Got an error getting a connection: %s", connectionResult.Err.Error())
		t.Fail()
		return
	}
	log.Info("Got a connection")
	connection := connectionResult.Connection
	defer connection.Close()

	// Build a browse request, to scan the KNX network for KNX devices
	// (Limiting the range to only the actually used range of addresses)
	browseRequestBuilder := connection.BrowseRequestBuilder()
	browseRequestBuilder.AddItem("findAllKnxDevices", "[1-3].[1-6].[1-60]")
	browseRequest, err := browseRequestBuilder.Build()
	if err != nil {
		log.Errorf("Got an error preparing browse-request: %s", connectionResult.Err.Error())
		t.Fail()
		return
	}

	// Execute the browse-request
	log.Info("Scanning for KNX devices")

	brr := browseRequest.ExecuteWithInterceptor(func(result apiModel.PlcBrowseEvent) bool {
		if result.Err != nil {
			return false
		}

		// Create a read-request to read the manufacturer and hardware ids
		readRequestBuilder := connection.ReadRequestBuilder()
		readRequestBuilder.AddItem("manufacturerId", result.Result.Address+"/0/12")
		readRequestBuilder.AddItem("applicationProgramVersion", result.Result.Address+"/3/13")
		readRequestBuilder.AddItem("interfaceProgramVersion", result.Result.Address+"/4/13")
		readRequest, err := readRequestBuilder.Build()
		if err != nil {
			log.Errorf("Got an error creating read-request: %s", err.Error())
			t.Fail()
			return false
		}

		// Execute the read-requests
		rrr := readRequest.Execute()
		readResult := <-rrr
		if readResult.Err != nil {
			log.Errorf("got an error executing read-request: %s", readResult.Err.Error())
			t.Fail()
			return false
		}

		// Check the response
		readResponse := readResult.Response
		if readResponse.GetResponseCode("manufacturerId") != apiModel.PlcResponseCode_OK {
			log.Errorf("Got an error response code %d for field 'manufacturerId'", readResponse.GetResponseCode("manufacturerId"))
			t.Fail()
			return false
		}
		if readResponse.GetResponseCode("applicationProgramVersion") != apiModel.PlcResponseCode_OK && readResponse.GetResponseCode("interfaceProgramVersion") != apiModel.PlcResponseCode_OK {
			log.Errorf("Got response code %d for address %s ('programVersion')",
				readResponse.GetResponseCode("applicationProgramVersion"), result.Result.Address+"/3/13")
			log.Errorf("Got response code %d for address %s ('programVersion')",
				readResponse.GetResponseCode("interfaceProgramVersion"), result.Result.Address+"/4/13")
			t.Fail()
		}

		manufacturerId := readResponse.GetValue("manufacturerId").GetUint16()
		if readResponse.GetResponseCode("applicationProgramVersion") == apiModel.PlcResponseCode_OK {
			programVersion := readResponse.GetValue("applicationProgramVersion")
			programVersionBytes := PlcValueUint8ListToByteArray(programVersion)
			log.Infof(" - Manufacturer Id: %d, Application Program Version: %s\n", manufacturerId, hex.EncodeToString(programVersionBytes))
		} else if readResponse.GetResponseCode("interfaceProgramVersion") == apiModel.PlcResponseCode_OK {
			programVersion := readResponse.GetValue("interfaceProgramVersion")
			programVersionBytes := PlcValueUint8ListToByteArray(programVersion)
			log.Infof(" - Manufacturer Id: %d, Interface Program Version: %s\n", manufacturerId, hex.EncodeToString(programVersionBytes))
		}
		return true
	})
	browseRequestResults := <-brr
	if browseRequestResults.Err != nil {
		log.Errorf("Got an error scanning for KNX devices: %s", connectionResult.Err.Error())
		t.Fail()
		return
	}

	log.Infof("Operation finished in %s", time.Since(startTime))
}

func TestKnxNetIpPlc4goGroupAddressRead(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(knxnetip.NewKnxNetIpDriver())
	driverManager.RegisterTransport(udp.NewUdpTransport())

	// Get a connection to a remote PLC
	crc := driverManager.GetConnection("knxnet-ip://192.168.42.11")

	// Wait for the driver to connect (or not)
	connectionResult := <-crc
	if connectionResult.Err != nil {
		t.Errorf("error connecting to PLC: %s", connectionResult.Err.Error())
		t.Fail()
		return
	}
	connection := connectionResult.Connection
	defer connection.Close()

	attributes := connection.GetMetadata().GetConnectionAttributes()
	fmt.Printf("Successfully connected to KNXnet/IP Gateway '%s' with KNX address '%s' got assigned client KNX address '%s'\n",
		attributes["GatewayName"],
		attributes["GatewayKnxAddress"],
		attributes["ClientKnxAddress"])

	// Try to ping the remote device
	pingResultChannel := connection.Ping()
	pingResult := <-pingResultChannel
	if pingResult.Err != nil {
		t.Errorf("couldn't ping device: %s", pingResult.Err.Error())
		t.Fail()
		return
	}

	srb := connection.SubscriptionRequestBuilder()
	srb.AddChangeOfStateItem("heating-actual-temperature", "*/*/10:DPT_Value_Temp")
	srb.AddChangeOfStateItem("heating-target-temperature", "*/*/11:DPT_Value_Temp")
	srb.AddChangeOfStateItem("heating-valve-open", "*/*/12:DPT_OpenClose")
	srb.AddItemHandler(knxEventHandler)
	subscriptionRequest, err := srb.Build()
	if err != nil {
		t.Errorf("error preparing subscription-request: %s", connectionResult.Err.Error())
		t.Fail()
		return
	}

	// Execute a subscription-request
	rrc := subscriptionRequest.Execute()

	// Wait for the response to finish
	rrr := <-rrc
	if rrr.Err != nil {
		t.Errorf("error executing read-request: %s", rrr.Err.Error())
		t.Fail()
		return
	}

	// Wait 2 minutes
	time.Sleep(120 * time.Second)

	// Execute a read request
	rrb := connection.ReadRequestBuilder()
	rrb.AddItem("energy-consumption", "1/1/211:DPT_Value_Power")
	rrb.AddItem("actual-temperatures", "*/*/10:DPT_Value_Temp")
	rrb.AddItem("set-temperatures", "*/*/11:DPT_Value_Temp")
	rrb.AddItem("window-status", "*/*/[60,64]:DPT_Value_Temp")
	rrb.AddItem("power-consumption", "*/*/[111,121,131,141]:DPT_Value_Temp")
	readRequest, err := rrb.Build()
	if err == nil {
		rrr := readRequest.Execute()
		readRequestResult := <-rrr
		if readRequestResult.Err == nil {
			for _, fieldName := range readRequestResult.Response.GetFieldNames() {
				if readRequestResult.Response.GetResponseCode(fieldName) == apiModel.PlcResponseCode_OK {
					fmt.Printf(" - Field %s Value %s\n", fieldName, readRequestResult.Response.GetValue(fieldName).GetString())
				}
			}
		}
	}
}

func TestKnxNetIpPlc4goPropertyRead(t *testing.T) {
	Init()

	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(knxnetip.NewKnxNetIpDriver())
	driverManager.RegisterTransport(udp.NewUdpTransport())

	// Get a connection to a remote PLC
	crc := driverManager.GetConnection("knxnet-ip://192.168.42.11")

	// Wait for the driver to connect (or not)
	connectionResult := <-crc
	if connectionResult.Err != nil {
		t.Errorf("error connecting to PLC: %s", connectionResult.Err.Error())
		t.Fail()
		return
	}
	connection := connectionResult.Connection
	defer connection.Close()

	readRequestBuilder := connection.ReadRequestBuilder()
	readRequestBuilder.AddItem("manufacturerId", "1.1.10/0/12")
	readRequestBuilder.AddItem("programVersion", "1.1.10/3/13")
	readRequest, _ := readRequestBuilder.Build()

	rrr := readRequest.Execute()
	readResult := <-rrr

	fmt.Printf("Got result %v", readResult)
}

func knxEventHandler(event apiModel.PlcSubscriptionEvent) {
	for _, fieldName := range event.GetFieldNames() {
		if event.GetResponseCode(fieldName) == apiModel.PlcResponseCode_OK {
			groupAddress := event.GetAddress(fieldName)
			fmt.Printf("Got update for field %s with address %s. Value changed to: %s\n",
				fieldName, groupAddress, event.GetValue(fieldName).GetString())
		}
	}
}

func TestKnxNetIpPlc4goMemoryRead(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(knxnetip.NewKnxNetIpDriver())
	driverManager.RegisterTransport(udp.NewUdpTransport())

	// Get a connection to a remote PLC
	crc := driverManager.GetConnection("knxnet-ip://192.168.42.11")

	// Wait for the driver to connect (or not)
	connectionResult := <-crc
	if connectionResult.Err != nil {
		t.Errorf("error connecting to PLC: %s", connectionResult.Err.Error())
		t.Fail()
		return
	}
	connection := connectionResult.Connection
	defer connection.Close()

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// Group Address Table reading
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	// First of all, request the starting address of the group address table
	readRequestBuilder := connection.ReadRequestBuilder()
	readRequestBuilder.AddItem("groupAddressTableAddress", "1.1.10/1/7")
	readRequest, err := readRequestBuilder.Build()
	if err != nil {
		t.Errorf("error creating read request: %s", err.Error())
		t.Fail()
		return
	}
	rrr := readRequest.Execute()
	readResult := <-rrr
	groupAddressTableStartAddress := readResult.Response.GetValue("groupAddressTableAddress").GetUint16()

	// Then read one byte at the given location.
	// This will return the number of entries in the group address table (each 2 bytes)
	readRequestBuilder = connection.ReadRequestBuilder()
	readRequestBuilder.AddItem("numberOfAddressTableEntries", fmt.Sprintf("1.1.10#%X:USINT",
		groupAddressTableStartAddress))
	readRequest, _ = readRequestBuilder.Build()
	rrr = readRequest.Execute()
	readResult = <-rrr
	numGroupAddresses := readResult.Response.GetValue("numberOfAddressTableEntries").GetUint8()

	// Read the data in the group address table
	readRequestBuilder = connection.ReadRequestBuilder()
	// TODO: This request needs to be automatically split up into multiple requests.
	// Reasons for splitting up:
	// - Max APDU Size exceeded
	// - Max 63 bytes readable in one request, due to max of count field
	readRequestBuilder.AddItem("groupAddressTable", fmt.Sprintf("1.1.10#%X:UINT[%d]",
		groupAddressTableStartAddress+3, numGroupAddresses-1))
	readRequest, _ = readRequestBuilder.Build()
	rrr = readRequest.Execute()
	readResult = <-rrr

	// Output the group addresses
	var knxGroupAddresses []*driverModel.KnxGroupAddress
	for _, groupAddress := range readResult.Response.GetValue("groupAddressTable").GetList() {
		groupAddress := knxnetip.Uint16ToKnxGroupAddress(groupAddress.GetUint16(), 3)
		knxGroupAddresses = append(knxGroupAddresses, groupAddress)
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// Group Address Association Table reading
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	// Now we read the group address association table address
	readRequestBuilder = connection.ReadRequestBuilder()
	readRequestBuilder.AddItem("groupAddressAssociationTableAddress", "1.1.10/2/7")
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		t.Errorf("error creating read request: %s", err.Error())
		t.Fail()
		return
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	groupAddressAssociationTableAddress := readResult.Response.GetValue("groupAddressAssociationTableAddress").GetUint16()

	// Then read one uint16 at the given location.
	// This will return the number of entries in the group address table (each 2 bytes)
	readRequestBuilder = connection.ReadRequestBuilder()
	readRequestBuilder.AddItem("numberOfGroupAddressAssociationTableEntries", fmt.Sprintf("1.1.10#%X:USINT",
		groupAddressAssociationTableAddress))
	readRequest, _ = readRequestBuilder.Build()
	rrr = readRequest.Execute()
	readResult = <-rrr
	numberOfGroupAddressAssociationTableEntries := readResult.Response.GetValue("numberOfGroupAddressAssociationTableEntries").GetUint8()

	// Read the data in the group address table
	readRequestBuilder = connection.ReadRequestBuilder()
	// TODO: This request needs to be automatically split up into multiple requests.
	// Reasons for splitting up:
	// - Max APDU Size exceeded
	// - Max 63 bytes readable in one request, due to max of count field
	readRequestBuilder.AddItem("groupAddressAssociationTable", fmt.Sprintf("1.1.10#%X:UINT[%d]",
		groupAddressAssociationTableAddress+1, numberOfGroupAddressAssociationTableEntries))
	readRequest, _ = readRequestBuilder.Build()
	rrr = readRequest.Execute()
	readResult = <-rrr

	// Output the group addresses
	for _, groupAddressAssociation := range readResult.Response.GetValue("groupAddressAssociationTable").GetList() {
		addressIndex := uint8(groupAddressAssociation.GetUint16() >> 8)
		comObjectNumber := uint8(groupAddressAssociation.GetUint16() & 0xFF)
		if (addressIndex > 0) && (addressIndex < uint8(len(knxGroupAddresses))) {
			groupAddress := knxGroupAddresses[addressIndex-1]
			fmt.Printf("Com Object %d bound to group address %s\n",
				comObjectNumber, knxnetip.GroupAddressToString(groupAddress))
		}
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// Com Object Table reading
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	// Now we read the group address association table address
	readRequestBuilder = connection.ReadRequestBuilder()
	readRequestBuilder.AddItem("comObjectTableAddress", "1.1.10/3/7")
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		t.Errorf("error creating read request: %s", err.Error())
		t.Fail()
		return
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	comObjectTableAddress := readResult.Response.GetValue("comObjectTableAddress").GetUint16()

	fmt.Printf("%d", comObjectTableAddress)
}

func PlcValueUint8ListToByteArray(value values.PlcValue) []byte {
	var result []byte
	for _, valueItem := range value.GetList() {
		result = append(result, valueItem.GetUint8())
	}
	return result
}
