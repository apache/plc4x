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
package main

import (
	"encoding/hex"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/modbus/readwrite/model"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

func main() {
	request, err := hex.DecodeString("000a00000006010300000004")
	if err != nil {
		// Output an error ...
	}
	rb := spi.ReadBufferNew(request)
	adu, err := model.ModbusTcpADUParse(*rb, false)
	if err != nil {
		// Output an error ...
	}
	if adu != nil {
		// Output success ...
	}

}
