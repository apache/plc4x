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
package model

import "plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"

type TPCI uint8

const (
	TPCI_UNNUMBERED_DATA_PACKET TPCI = 0x0
	TPCI_UNNUMBERED             TPCI = 0x1
	TPCI_NUMBERED_DATA_PACKET   TPCI = 0x2
	TPCI_NUMBERED_CONTROL_DATA  TPCI = 0x3
)

func CastTPCI(structType interface{}) TPCI {
	castFunc := func(typ interface{}) TPCI {
		if sTPCI, ok := typ.(TPCI); ok {
			return sTPCI
		}
		return 0
	}
	return castFunc(structType)
}

func TPCIParse(io spi.ReadBuffer) (TPCI, error) {
	// TODO: Implement ...
	return 0, nil
}

func (e TPCI) Serialize(io spi.WriteBuffer) {
	// TODO: Implement ...
}
