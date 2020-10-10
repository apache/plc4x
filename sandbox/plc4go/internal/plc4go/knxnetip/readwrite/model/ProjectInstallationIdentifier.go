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

import (
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type ProjectInstallationIdentifier struct {
	projectNumber      uint8
	installationNumber uint8
}

// The corresponding interface
type IProjectInstallationIdentifier interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewProjectInstallationIdentifier(projectNumber uint8, installationNumber uint8) spi.Message {
	return &ProjectInstallationIdentifier{projectNumber: projectNumber, installationNumber: installationNumber}
}

func (m ProjectInstallationIdentifier) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Simple field (projectNumber)
	lengthInBits += 8

	// Simple field (installationNumber)
	lengthInBits += 8

	return lengthInBits
}

func (m ProjectInstallationIdentifier) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ProjectInstallationIdentifierParse(io spi.ReadBuffer) (spi.Message, error) {

	// Simple Field (projectNumber)
	var projectNumber uint8 = io.ReadUint8(8)

	// Simple Field (installationNumber)
	var installationNumber uint8 = io.ReadUint8(8)

	// Create the instance
	return NewProjectInstallationIdentifier(projectNumber, installationNumber), nil
}

func (m ProjectInstallationIdentifier) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IProjectInstallationIdentifier); ok {

			// Simple Field (projectNumber)
			var projectNumber uint8 = m.projectNumber
			io.WriteUint8(8, (projectNumber))

			// Simple Field (installationNumber)
			var installationNumber uint8 = m.installationNumber
			io.WriteUint8(8, (installationNumber))
		}
	}
	serializeFunc(m)
}
