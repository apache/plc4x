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
	"fmt"
	"strings"
	"time"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
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

// blockTypeDataBlock is the "list blocks of type" wire code for DBs: the two ASCII chars "0A".
const blockTypeDataBlock = uint16(0x3041)

func buildListBlocksOfTypeRequest(tpduId uint16, blockType uint16) readWriteModel.TPKTPacket {
	return newUserDataMessage(
		tpduId,
		newCpuFunctionsRequestParameter(0x03, 0x02), // block functions / list blocks of type
		readWriteModel.NewS7PayloadUserDataItemCpuFunctionListBlocksOfTypeRequest(
			readWriteModel.DataTransportErrorCode_OK,
			readWriteModel.DataTransportSize_OCTET_STRING,
			2, // dataLength: the packed block-type code
			blockType,
		),
	)
}

// parseListBlocksOfTypeResponse extracts the block numbers from a list-blocks-of-type
// response. Each entry is 4 bytes: blockNumber (uint16 BE), flags, language.
func parseListBlocksOfTypeResponse(message readWriteModel.S7Message) ([]uint16, error) {
	messageUserData, ok := message.(readWriteModel.S7MessageUserData)
	if !ok {
		return nil, errors.Errorf("expected S7MessageUserData, got %T", message)
	}
	if err := checkUserDataErrorCode(messageUserData); err != nil {
		return nil, err
	}
	payloadUserData, ok := messageUserData.GetPayload().(readWriteModel.S7PayloadUserData)
	if !ok {
		return nil, errors.Errorf("expected S7PayloadUserData, got %T", messageUserData.GetPayload())
	}
	for _, item := range payloadUserData.GetItems() {
		response, ok := item.(readWriteModel.S7PayloadUserDataItemCpuFunctionListBlocksOfTypeResponse)
		if !ok {
			continue
		}
		data := response.GetItems()
		blockNumbers := make([]uint16, 0, len(data)/4)
		for i := 0; i+4 <= len(data); i += 4 {
			blockNumbers = append(blockNumbers, uint16(data[i])<<8|uint16(data[i+1]))
		}
		return blockNumbers, nil
	}
	return nil, errors.New("no list-blocks-of-type response item in payload")
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

// msgSubscriptionMagicKey is the magic key the PLC echoes back in subscription confirmations.
// The legacy driver used "HmiRtm  " (HMI runtime, padded to 8 ASCII bytes); keep it for
// compatibility with PLCs matching the key pattern in some firmware paths.
const msgSubscriptionMagicKey = "HmiRtm  "

// alarmStateForController picks the alarm stream to arm: S7-400 uses the block-of-8 ALARM_8
// family, everything else the ALARM_S/SQ family (SFC17/18).
func alarmStateForController(controllerType readWriteModel.ControllerType, abort bool) readWriteModel.AlarmStateType {
	if controllerType == readWriteModel.ControllerType_S7_400 {
		if abort {
			return readWriteModel.AlarmStateType_ALARM_ABORT
		}
		return readWriteModel.AlarmStateType_ALARM_INITIATE
	}
	if abort {
		return readWriteModel.AlarmStateType_ALARM_S_ABORT
	}
	return readWriteModel.AlarmStateType_ALARM_S_INITIATE
}

// buildMsgSubscriptionRequest arms (or with an abort alarm state: cancels) the ALARM_S/SQ
// push stream on this connection (UserData group 0x04, subfunction 0x02, 12-byte variant).
func buildMsgSubscriptionRequest(tpduId uint16, alarmState readWriteModel.AlarmStateType) readWriteModel.TPKTPacket {
	reserve := uint8(0)
	return newUserDataMessage(
		tpduId,
		newCpuFunctionsRequestParameter(0x04, 0x02), // CPU functions / MsgSubscription
		readWriteModel.NewS7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest(
			readWriteModel.DataTransportErrorCode_OK,
			readWriteModel.DataTransportSize_OCTET_STRING,
			12,   // dataLength: the ALARM_S/SQ variant with appended alarm state
			0x80, // subscription bitmask: ALM
			msgSubscriptionMagicKey,
			&alarmState,
			&reserve,
		),
	)
}

func parseMsgSubscriptionResponse(message readWriteModel.S7Message) (bool, error) {
	messageUserData, ok := message.(readWriteModel.S7MessageUserData)
	if !ok {
		return false, errors.Errorf("expected S7MessageUserData, got %T", message)
	}
	if err := checkUserDataErrorCode(messageUserData); err != nil {
		return false, err
	}
	payloadUserData, ok := messageUserData.GetPayload().(readWriteModel.S7PayloadUserData)
	if !ok {
		return false, errors.Errorf("expected S7PayloadUserData, got %T", messageUserData.GetPayload())
	}
	for _, item := range payloadUserData.GetItems() {
		switch response := item.(type) {
		case readWriteModel.S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse:
			result := response.GetResult()
			return result == 0x00 || result == 0x02, nil
		case readWriteModel.S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse:
			// Plain "ok" response with no body.
			return true, nil
		}
	}
	return false, nil
}

// cyclicInterval is the TimeBase + factor pair approximating a requested cycle time.
type cyclicInterval struct {
	base   readWriteModel.TimeBase
	factor uint8
}

func (c cyclicInterval) toDuration() time.Duration {
	base := 100 * time.Millisecond
	switch c.base {
	case readWriteModel.TimeBase_B1SEC:
		base = time.Second
	case readWriteModel.TimeBase_B10SEC:
		base = 10 * time.Second
	}
	return base * time.Duration(c.factor)
}

// pickCyclicInterval prefers the coarsest base whose factor still fits into a byte -
// empirically required for S7-300 firmware (>= V3.2), which rejects the 100ms base with
// errCode 0xD804 and only honors B1SEC/B10SEC. Range: 100ms ... 2550s; non-positive
// durations default to 1s.
func pickCyclicInterval(requested time.Duration) cyclicInterval {
	ms := int64(1000)
	if requested > 0 {
		ms = requested.Milliseconds()
		if ms < 100 {
			ms = 100
		}
	}
	if ms >= 10_000 && ms <= 2_550_000 && ms%10_000 == 0 {
		return cyclicInterval{readWriteModel.TimeBase_B10SEC, uint8(ms / 10_000)}
	}
	if ms >= 1_000 && ms <= 255_000 && ms%1_000 == 0 {
		return cyclicInterval{readWriteModel.TimeBase_B1SEC, uint8(ms / 1_000)}
	}
	// Round near-second cadences to full seconds rather than using the often-rejected
	// 100ms base when the cadence changes by less than ~10%.
	if ms >= 900 && ms <= 255_000 {
		factor := (ms + 500) / 1_000
		if factor < 1 {
			factor = 1
		}
		return cyclicInterval{readWriteModel.TimeBase_B1SEC, uint8(factor)}
	}
	if ms >= 9_000 && ms <= 2_550_000 {
		factor := (ms + 5_000) / 10_000
		if factor < 1 {
			factor = 1
		}
		return cyclicInterval{readWriteModel.TimeBase_B10SEC, uint8(factor)}
	}
	// Genuine sub-second request - only the 100ms base can express it; some firmwares
	// reject this with 0xD804, callers should treat that as "not supported".
	factor := (ms + 50) / 100
	if factor < 1 {
		factor = 1
	}
	if factor > 255 {
		factor = 255
	}
	return cyclicInterval{readWriteModel.TimeBase_B01SEC, uint8(factor)}
}

func buildCyclicSubscribeRequest(tpduId uint16, tags []PlcTag, interval cyclicInterval) readWriteModel.TPKTPacket {
	items := make([]readWriteModel.CycServiceItemType, len(tags))
	for i, tag := range tags {
		// The 24-bit address encodes bit-addressed memory: byteOffset shifted left by 3
		// plus the bit number - same as in a regular S7AddressAny.
		address := uint32(tag.GetByteOffset())<<3 | uint32(tag.GetBitOffset()&0x07)
		items[i] = readWriteModel.NewCycServiceItemAnyType(
			0x0a, // byteLength: 10 bytes after the byteLength field
			0x10, // syntaxId: ANY
			tag.GetDataType(),
			tag.GetNumElements(),
			tag.GetBlockNumber(),
			tag.GetMemoryArea(),
			address,
		)
	}
	// dataLength: itemsCount(2) + timeBase(1) + timeFactor(1) + 12 bytes per ANY item.
	payloadLength := uint16(4 + len(items)*12)
	return newUserDataMessage(
		tpduId,
		newCpuFunctionsRequestParameter(0x02, 0x01), // cyclic services / subscribe
		readWriteModel.NewS7PayloadUserDataItemCyclicServicesSubscribeRequest(
			readWriteModel.DataTransportErrorCode_OK,
			readWriteModel.DataTransportSize_OCTET_STRING,
			payloadLength,
			uint16(len(items)),
			interval.base,
			interval.factor,
			items,
		),
	)
}

func buildCyclicUnsubscribeRequest(tpduId uint16, jobId uint8) readWriteModel.TPKTPacket {
	return newUserDataMessage(
		tpduId,
		newCpuFunctionsRequestParameter(0x02, 0x04), // cyclic services / unsubscribe
		readWriteModel.NewS7PayloadUserDataItemCyclicServicesUnsubscribeRequest(
			readWriteModel.DataTransportErrorCode_OK,
			readWriteModel.DataTransportSize_OCTET_STRING,
			2,
			0x05, // function: unsubscribe single job
			jobId,
		),
	)
}

// parseCyclicSubscribeResponse returns the jobId the PLC assigned to the subscription
// (carried as the parameter's sequence number).
func parseCyclicSubscribeResponse(message readWriteModel.S7Message) (uint8, error) {
	messageUserData, ok := message.(readWriteModel.S7MessageUserData)
	if !ok {
		return 0, errors.Errorf("expected S7MessageUserData, got %T", message)
	}
	if err := checkUserDataErrorCode(messageUserData); err != nil {
		return 0, err
	}
	payloadUserData, ok := messageUserData.GetPayload().(readWriteModel.S7PayloadUserData)
	if !ok {
		return 0, errors.Errorf("expected S7PayloadUserData, got %T", messageUserData.GetPayload())
	}
	for _, item := range payloadUserData.GetItems() {
		switch item.(type) {
		// Two success shapes: a response carrying the initial values, or an empty response
		// when the PLC doesn't bundle them - pushes start arriving on the cycle anyway.
		case readWriteModel.S7PayloadUserDataItemCyclicServicesSubscribeResponse,
			readWriteModel.S7PayloadUserDataItemCyclicServicesSubscribeEmptyResponse:
			jobId, ok := userDataSequenceNumber(messageUserData)
			if !ok {
				return 0, errors.New("cyclic subscribe response carries no sequence number")
			}
			return jobId, nil
		case readWriteModel.S7PayloadUserDataItemCyclicServicesErrorResponse:
			return 0, errors.New("cyclic subscribe rejected by PLC")
		}
	}
	return 0, errors.New("no cyclic subscribe response item in payload")
}

func userDataSequenceNumber(message readWriteModel.S7MessageUserData) (uint8, bool) {
	parameterUserData, ok := message.GetParameter().(readWriteModel.S7ParameterUserData)
	if !ok {
		return 0, false
	}
	for _, item := range parameterUserData.GetItems() {
		if cpuFunctions, ok := item.(readWriteModel.S7ParameterUserDataItemCPUFunctions); ok {
			return cpuFunctions.GetSequenceNumber(), true
		}
	}
	return 0, false
}

// userDataPushKey extracts the (functionGroup, functionType, subfunction) routing triple.
func userDataPushKey(message readWriteModel.S7MessageUserData) (group uint8, functionType uint8, subfunction uint8, ok bool) {
	parameterUserData, isUserData := message.GetParameter().(readWriteModel.S7ParameterUserData)
	if !isUserData {
		return 0, 0, 0, false
	}
	for _, item := range parameterUserData.GetItems() {
		if cpuFunctions, isCpuFunctions := item.(readWriteModel.S7ParameterUserDataItemCPUFunctions); isCpuFunctions {
			return cpuFunctions.GetCpuFunctionGroup(), cpuFunctions.GetCpuFunctionType(), cpuFunctions.GetCpuSubfunction(), true
		}
	}
	return 0, 0, 0, false
}

// extractCyclicPushItems returns the per-item raw payloads of a cyclic services push.
func extractCyclicPushItems(message readWriteModel.S7MessageUserData) [][]byte {
	payloadUserData, ok := message.GetPayload().(readWriteModel.S7PayloadUserData)
	if !ok {
		return nil
	}
	toByteSlices := func(values []readWriteModel.AssociatedValueType) [][]byte {
		result := make([][]byte, len(values))
		for i, value := range values {
			result[i] = value.GetData()
		}
		return result
	}
	toByteSlicesQuery := func(values []readWriteModel.AssociatedQueryValueType) [][]byte {
		result := make([][]byte, len(values))
		for i, value := range values {
			result[i] = value.GetData()
		}
		return result
	}
	for _, item := range payloadUserData.GetItems() {
		switch push := item.(type) {
		case readWriteModel.S7PayloadUserDataItemCyclicServicesPush:
			return toByteSlices(push.GetItems())
		case readWriteModel.S7PayloadUserDataItemCyclicServicesChangeDrivenPush:
			return toByteSlicesQuery(push.GetItems())
		}
	}
	return nil
}

func buildAlarmQueryRequest(tpduId uint16, queryType readWriteModel.QueryType) readWriteModel.TPKTPacket {
	alarmType := readWriteModel.AlarmType_ALARM_S
	if queryType == readWriteModel.QueryType_ALARM_8 {
		alarmType = readWriteModel.AlarmType_ALARM_8
	}
	return newUserDataMessage(
		tpduId,
		newCpuFunctionsRequestParameter(0x04, 0x13), // CPU functions / AlarmQuery
		readWriteModel.NewS7PayloadUserDataItemCpuFunctionAlarmQueryRequest(
			readWriteModel.DataTransportErrorCode_OK,
			readWriteModel.DataTransportSize_OCTET_STRING,
			// 4 const bytes + 5 variable bytes = 9; the legacy driver used 12 and real
			// S7-300s absorb the extra framing, so keep the proven value.
			12,
			readWriteModel.SyntaxIdType_ALARM_QUERYREQSET,
			queryType,
			alarmType,
		),
	)
}

// parseAlarmQueryResponse extracts the raw byte payload of an AlarmQuery response.
func parseAlarmQueryResponse(message readWriteModel.S7Message) ([]byte, error) {
	messageUserData, ok := message.(readWriteModel.S7MessageUserData)
	if !ok {
		return nil, errors.Errorf("expected S7MessageUserData, got %T", message)
	}
	if err := checkUserDataErrorCode(messageUserData); err != nil {
		return nil, err
	}
	payloadUserData, ok := messageUserData.GetPayload().(readWriteModel.S7PayloadUserData)
	if !ok {
		return nil, errors.Errorf("expected S7PayloadUserData, got %T", messageUserData.GetPayload())
	}
	for _, item := range payloadUserData.GetItems() {
		if response, ok := item.(readWriteModel.S7PayloadUserDataItemCpuFunctionAlarmQueryResponse); ok {
			return response.GetItems(), nil
		}
	}
	return nil, errors.New("no alarm query response item in payload")
}

// alarmIndicationSubfunctions are the push subfunctions carrying alarm indication payloads.
var alarmIndicationSubfunctions = map[uint8]struct{}{
	0x05: {}, // ALARM8
	0x06: {}, // NOTIFY
	0x0c: {}, // ALARM_ACK_IND
	0x11: {}, // ALARM_SQ
	0x12: {}, // ALARM_S
	0x13: {}, // ALARM_SC
	0x16: {}, // NOTIFY8
}

// parseAlarmIndication turns a pushed alarm message into a PlcStruct. The struct keys match
// the plc4j driver: plcTimestamp, functionId, numberOfObjects, eventId, eventState,
// localState, ackStateGoing, ackStateComing, receivedAt.
func parseAlarmIndication(message readWriteModel.S7MessageUserData) apiValues.PlcValue {
	payloadUserData, ok := message.GetPayload().(readWriteModel.S7PayloadUserData)
	if !ok {
		return nil
	}
	for _, item := range payloadUserData.GetItems() {
		switch payload := item.(type) {
		case readWriteModel.S7PayloadAlarmS:
			return alarmStructFromPush(payload.GetAlarmMessage())
		case readWriteModel.S7PayloadAlarmSQ:
			return alarmStructFromPush(payload.GetAlarmMessage())
		case readWriteModel.S7PayloadAlarmSC:
			return alarmStructFromPush(payload.GetAlarmMessage())
		case readWriteModel.S7PayloadAlarm8:
			return alarmStructFromPush(payload.GetAlarmMessage())
		case readWriteModel.S7PayloadNotify:
			return alarmStructFromPush(payload.GetAlarmMessage())
		case readWriteModel.S7PayloadNotify8:
			return alarmStructFromPush(payload.GetAlarmMessage())
		case readWriteModel.S7PayloadAlarmAckInd:
			return alarmStructFromAckPush(payload.GetAlarmMessage())
		}
	}
	return nil
}

func alarmStructFromPush(push readWriteModel.AlarmMessagePushType) apiValues.PlcValue {
	if push == nil {
		return nil
	}
	children := map[string]apiValues.PlcValue{}
	addAlarmHeader(children, push.GetTimeStamp(), push.GetFunctionId(), push.GetNumberOfObjects())
	// Surface the first object directly - the common case is exactly one object per
	// indication; multi-object indications can be inspected via numberOfObjects.
	if objects := push.GetMessageObjects(); len(objects) > 0 {
		first := objects[0]
		children["eventId"] = spiValues.NewPlcUDINT(first.GetEventId())
		children["eventState"] = spiValues.NewPlcUDINT(stateToMask(first.GetEventState()))
		children["localState"] = spiValues.NewPlcUDINT(stateToMask(first.GetLocalState()))
		children["ackStateGoing"] = spiValues.NewPlcUDINT(stateToMask(first.GetAckStateGoing()))
		children["ackStateComing"] = spiValues.NewPlcUDINT(stateToMask(first.GetAckStateComing()))
	}
	children["receivedAt"] = spiValues.NewPlcLINT(time.Now().UnixMilli())
	return spiValues.NewPlcStruct(children)
}

func alarmStructFromAckPush(push readWriteModel.AlarmMessageAckPushType) apiValues.PlcValue {
	if push == nil {
		return nil
	}
	children := map[string]apiValues.PlcValue{}
	addAlarmHeader(children, push.GetTimeStamp(), push.GetFunctionId(), push.GetNumberOfObjects())
	if objects := push.GetMessageObjects(); len(objects) > 0 {
		first := objects[0]
		children["eventId"] = spiValues.NewPlcUDINT(first.GetEventId())
		children["ackStateGoing"] = spiValues.NewPlcUDINT(stateToMask(first.GetAckStateGoing()))
		children["ackStateComing"] = spiValues.NewPlcUDINT(stateToMask(first.GetAckStateComing()))
	}
	children["receivedAt"] = spiValues.NewPlcLINT(time.Now().UnixMilli())
	return spiValues.NewPlcStruct(children)
}

func addAlarmHeader(children map[string]apiValues.PlcValue, timestamp readWriteModel.DateAndTime, functionId uint8, numberOfObjects uint8) {
	children["plcTimestamp"] = spiValues.NewPlcSTRING(formatAlarmDateAndTime(timestamp))
	children["functionId"] = spiValues.NewPlcUDINT(uint32(functionId))
	children["numberOfObjects"] = spiValues.NewPlcUDINT(uint32(numberOfObjects))
}

func stateToMask(state readWriteModel.State) uint32 {
	if state == nil {
		return 0
	}
	var mask uint32
	if state.GetSIG_1() {
		mask |= 0x01
	}
	if state.GetSIG_2() {
		mask |= 0x02
	}
	if state.GetSIG_3() {
		mask |= 0x04
	}
	if state.GetSIG_4() {
		mask |= 0x08
	}
	if state.GetSIG_5() {
		mask |= 0x10
	}
	if state.GetSIG_6() {
		mask |= 0x20
	}
	if state.GetSIG_7() {
		mask |= 0x40
	}
	if state.GetSIG_8() {
		mask |= 0x80
	}
	return mask
}

func formatAlarmDateAndTime(dateAndTime readWriteModel.DateAndTime) string {
	if dateAndTime == nil {
		return ""
	}
	// The BCD year is 2-digit; S7 convention: < 90 -> 2000s, >= 90 -> 1900s.
	year := uint16(dateAndTime.GetYear())
	if year < 90 {
		year += 2000
	} else {
		year += 1900
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%03d",
		year, dateAndTime.GetMonth(), dateAndTime.GetDay(),
		dateAndTime.GetHour(), dateAndTime.GetMinutes(), dateAndTime.GetSeconds(),
		dateAndTime.GetMsec())
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
