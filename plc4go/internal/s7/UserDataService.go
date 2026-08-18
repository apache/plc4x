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

package s7

import (
	"strings"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// Helpers for the S7Comm UserData services (SZL identification, block functions,
// message subscriptions, cyclic services). Ported from the plc4j userdata services.

// newCpuFunctionsRequestParameter builds the CPUFunctions parameter for a UserData request.
func newCpuFunctionsRequestParameter(cpuFunctionGroup uint8, cpuSubfunction uint8) readWriteModel.S7ParameterUserDataItem {
	return readWriteModel.NewS7ParameterUserDataItemCPUFunctions(
		0x11, // method: request
		0x04, // cpuFunctionType: request
		cpuFunctionGroup,
		cpuSubfunction,
		0x00, // sequenceNumber
		nil,
		nil,
		nil,
	)
}

// newUserDataMessage wraps a single parameter/payload item pair into an S7MessageUserData TPKT packet.
func newUserDataMessage(tpduId uint16, parameter readWriteModel.S7ParameterUserDataItem, payload readWriteModel.S7PayloadUserDataItem) readWriteModel.TPKTPacket {
	message := readWriteModel.NewS7MessageUserData(
		tpduId,
		readWriteModel.NewS7ParameterUserData([]readWriteModel.S7ParameterUserDataItem{parameter}),
		readWriteModel.NewS7PayloadUserData([]readWriteModel.S7PayloadUserDataItem{payload}),
	)
	return readWriteModel.NewTPKTPacket(readWriteModel.NewCOTPPacketData(nil, message, true, uint8(tpduId)))
}

// szlIdModuleIdentification works on S7-300/400.
func szlIdModuleIdentification() readWriteModel.SzlId {
	return readWriteModel.NewSzlId(readWriteModel.SzlModuleTypeClass_CPU, 0x00, readWriteModel.SzlSublist_MODULE_IDENTIFICATION)
}

// szlIdComponentIdentification is the standard for S7-1200/1500; with index 0x0001 the
// response carries the CPU order number (MLFB) in ASCII.
func szlIdComponentIdentification() readWriteModel.SzlId {
	return readWriteModel.NewSzlId(readWriteModel.SzlModuleTypeClass_CPU, 0x00, readWriteModel.SzlSublist_COMPONENT_IDENTIFICATION)
}

func buildSzlRequest(tpduId uint16, szlId readWriteModel.SzlId, szlIndex uint16) readWriteModel.TPKTPacket {
	return newUserDataMessage(
		tpduId,
		newCpuFunctionsRequestParameter(0x04, 0x01), // CPU functions / ReadSZL
		readWriteModel.NewS7PayloadUserDataItemCpuFunctionReadSzlRequest(
			readWriteModel.DataTransportErrorCode_OK,
			readWriteModel.DataTransportSize_OCTET_STRING,
			4, // dataLength: SzlId(2) + SzlIndex(2)
			szlId,
			szlIndex,
		),
	)
}

// checkUserDataErrorCode returns an error if any CPUFunctions parameter item reports a non-zero error code.
func checkUserDataErrorCode(message readWriteModel.S7MessageUserData) error {
	parameterUserData, ok := message.GetParameter().(readWriteModel.S7ParameterUserData)
	if !ok {
		return errors.Errorf("expected S7ParameterUserData, got %T", message.GetParameter())
	}
	for _, item := range parameterUserData.GetItems() {
		if cpuFunctions, ok := item.(readWriteModel.S7ParameterUserDataItemCPUFunctions); ok {
			if errorCode := cpuFunctions.GetErrorCode(); errorCode != nil && *errorCode != 0 {
				return errors.Errorf("request rejected by PLC, errorCode=0x%x", *errorCode)
			}
		}
	}
	return nil
}

// parseSzlProbeResponse extracts the device descriptor (article number or CPU model name)
// from an SZL identification response and derives the controller type from it.
func parseSzlProbeResponse(message readWriteModel.S7Message) (string, readWriteModel.ControllerType, error) {
	messageUserData, ok := message.(readWriteModel.S7MessageUserData)
	if !ok {
		return "", 0, errors.Errorf("expected S7MessageUserData, got %T", message)
	}
	if err := checkUserDataErrorCode(messageUserData); err != nil {
		return "", 0, err
	}
	payloadUserData, ok := messageUserData.GetPayload().(readWriteModel.S7PayloadUserData)
	if !ok {
		return "", 0, errors.Errorf("expected S7PayloadUserData, got %T", messageUserData.GetPayload())
	}
	for _, item := range payloadUserData.GetItems() {
		szlResponse, ok := item.(readWriteModel.S7PayloadUserDataItemCpuFunctionReadSzlResponse)
		if !ok {
			continue
		}
		data := szlResponse.GetItems()
		if len(data) < 4 {
			return "", 0, errors.Errorf("SZL response too short: %d bytes", len(data))
		}
		// Best-effort identification: the order-number prefix first, then the model-name
		// pattern (e.g. "CPU 315-2 PN/DP" on S7-300). A structurally valid response without
		// either still proves the device speaks UserData services.
		descriptor := findArticleNumber(data)
		if descriptor == "" {
			descriptor = findCpuModelName(data)
		}
		return descriptor, decodeControllerType(descriptor), nil
	}
	return "", 0, errors.New("no SZL response item in payload")
}

// findArticleNumber scans for a printable ASCII run starting with a Siemens order-number
// prefix ("6ES7"/"6GK") anywhere in the raw SZL item bytes.
func findArticleNumber(data []byte) string {
	for start := 0; start+12 <= len(data); start++ {
		if data[start] != '6' {
			continue
		}
		end := start
		maxEnd := start + 32
		if maxEnd > len(data) {
			maxEnd = len(data)
		}
		for end < maxEnd && data[end] >= 0x20 && data[end] < 0x7F {
			end++
		}
		if end-start < 12 {
			continue
		}
		candidate := strings.TrimSpace(string(data[start:end]))
		if strings.HasPrefix(candidate, "6ES7") || strings.HasPrefix(candidate, "6GK") {
			return candidate
		}
	}
	return ""
}

// findCpuModelName scans for a "CPU ..." model-name run in the raw SZL item bytes.
func findCpuModelName(data []byte) string {
	for start := 0; start+6 <= len(data); start++ {
		if data[start] != 'C' || data[start+1] != 'P' || data[start+2] != 'U' || data[start+3] != ' ' {
			continue
		}
		end := start
		maxEnd := start + 32
		if maxEnd > len(data) {
			maxEnd = len(data)
		}
		for end < maxEnd && data[end] >= 0x20 && data[end] < 0x7F {
			end++
		}
		candidate := strings.TrimSpace(string(data[start:end]))
		if len(candidate) > 4 {
			return candidate
		}
	}
	return ""
}

// decodeControllerType derives the controller family from an order number
// (e.g. "6ES7 212-1AE40-0XB0") or a model name (e.g. "CPU 315-2 PN/DP").
func decodeControllerType(descriptor string) readWriteModel.ControllerType {
	if strings.HasPrefix(descriptor, "6ES7") {
		rest := strings.TrimLeft(descriptor[4:], " ")
		if rest == "" {
			return readWriteModel.ControllerType_ANY
		}
		switch rest[0] {
		case '2':
			return readWriteModel.ControllerType_S7_1200
		case '5':
			return readWriteModel.ControllerType_S7_1500
		case '3':
			return readWriteModel.ControllerType_S7_300
		case '4':
			return readWriteModel.ControllerType_S7_400
		default:
			return readWriteModel.ControllerType_ANY
		}
	}
	if strings.HasPrefix(descriptor, "CPU ") {
		rest := strings.TrimLeft(descriptor[4:], " ")
		digits := 0
		number := 0
		for digits < len(rest) && digits < 4 && rest[digits] >= '0' && rest[digits] <= '9' {
			number = number*10 + int(rest[digits]-'0')
			digits++
		}
		if digits == 3 {
			switch rest[0] {
			case '3':
				return readWriteModel.ControllerType_S7_300
			case '4':
				return readWriteModel.ControllerType_S7_400
			}
		}
		if digits == 4 {
			if number >= 1200 && number <= 1299 {
				return readWriteModel.ControllerType_S7_1200
			}
			if number >= 1500 && number <= 1599 {
				return readWriteModel.ControllerType_S7_1500
			}
		}
	}
	return readWriteModel.ControllerType_ANY
}

// supportsUserDataServices reports whether the controller family supports the S7Comm
// UserData services (browse, alarms, cyclic subscriptions).
func supportsUserDataServices(controllerType readWriteModel.ControllerType) bool {
	switch controllerType {
	case readWriteModel.ControllerType_S7_300,
		readWriteModel.ControllerType_S7_400,
		readWriteModel.ControllerType_S7_1200,
		readWriteModel.ControllerType_S7_1500:
		return true
	default: // S7_200, LOGO, ANY
		return false
	}
}
