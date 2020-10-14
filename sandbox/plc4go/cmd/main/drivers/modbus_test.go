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
	"net"
	"os"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/modbus/readwrite/model"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"strings"
	"testing"
)

func TestModbus(t *testing.T) {
	test(t, "000000000006ff0408d20002", false)
	test(t, "7cfe000000c9ff04c600000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000004000000000000000000000000000001db000001d600004a380000000000000000000000000000000000000000000000000000000000006461696d006e0000000000000000000000000000303100300000000000000000000000000000000000000000000000000000000000000000000000000000", true)
	test(t, "000a0000001101140e060003270e000206000400000008", false)
	test(t, "000a0000001b011418050600000000110600000000000000000000000000000000", true)
	test(t, "000a0000000c011509060002000000010008", false)
	test(t, "000a00000015011512060001270F00010000060002000000010000", false)
}

func test(t *testing.T, rawMessage string, response bool) {
	// Create the input data
	// "000a 0000 0006 01 03 00 00 00 04"
	request, err := hex.DecodeString(rawMessage)
	if err != nil {
		t.Errorf("Error decoding test input")
	}
	rb := spi.ReadBufferNew(request)
	adu, err := model.ModbusTcpADUParse(rb, response)
	if err != nil {
		t.Errorf("Error parsing: %s", err)
	}
	if adu != nil {
		wb := spi.WriteBufferNew()
		val := model.CastIModbusTcpADU(adu)
		val.Serialize(*wb)
		serializedMessage := hex.EncodeToString(wb.GetBytes())
		if strings.ToUpper(serializedMessage) != strings.ToUpper(rawMessage) {
			t.Errorf("The serilized result doesn't match the input")
		}
	}
}

func TestConnection(t *testing.T) {

	pdu := model.ModbusPDUReadInputRegistersRequest{
		StartingAddress: 1,
		Quantity:        1,
	}
	adu := model.ModbusTcpADU{
		TransactionIdentifier: 0,
		UnitIdentifier:        255,
		Pdu:                   &pdu,
	}

	wb := spi.WriteBufferNew()
	adu.Serialize(*wb)

	servAddr := "192.168.23.30:502"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write(wb.GetBytes())
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, 1024)

	numBytes, err := conn.Read(buffer)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	rb := spi.ReadBufferNew(buffer[0:numBytes])
	response, err := model.ModbusTcpADUParse(rb, true)
	if err != nil {
		println("Parsing response failed:", err.Error())
		os.Exit(1)
	}

	fmt.Println(response)

	conn.Close()
}
