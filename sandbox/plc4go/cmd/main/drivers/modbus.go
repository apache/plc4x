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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/modbus/readwrite/model"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"strings"
)

func Modbus() {
	test("000000000006ff0408d20002", false)
	test("7cfe000000c9ff04c600000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000004000000000000000000000000000001db000001d600004a380000000000000000000000000000000000000000000000000000000000006461696d006e0000000000000000000000000000303100300000000000000000000000000000000000000000000000000000000000000000000000000000", true)
	test("000a0000001101140e060003270e000206000400000008", false)
	test("000a0000001b011418050600000000110600000000000000000000000000000000", true)
	test("000a0000000c011509060002000000010008", false)
	test("000a00000015011512060001270F00010000060002000000010000", false)
}

func test(rawMessage string, response bool) {
	// Create the input data
	// "000a 0000 0006 01 03 00 00 00 04"
	request, err := hex.DecodeString(rawMessage)
	if err != nil {
		fmt.Errorf("error preparing input buffer")
		return
	}
	rb := spi.ReadBufferNew(request)
	adu, err := model.ModbusTcpADUParse(rb, response)
	if err != nil {
		fmt.Errorf("error parsing input")
		return
	}
	if adu != nil {
		wb := spi.WriteBufferNew()
		val := model.CastIModbusTcpADU(adu)
		val.Serialize(*wb)
		serializedMessage := hex.EncodeToString(wb.GetBytes())
		if strings.ToUpper(serializedMessage) == strings.ToUpper(rawMessage) {
			fmt.Println("Success")
		} else {
			fmt.Println("Failure")
		}
	}
}
