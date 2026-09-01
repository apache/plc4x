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

package iec608705104

import (
	"strings"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/iec608705104/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// The keys of the struct every subscription event carries.
//
// Why a struct rather than a bare value: IEC 60870-5-104 hangs quality descriptors (SIQ, DIQ, QDS,
// QDP) off nearly every measurement, and for grid telemetry those flags are not decoration - a
// reading flagged invalid, blocked, substituted or not-topical must not reach an operator looking
// like a good reading. plc4j attaches them as PlcValue metadata; plc4go has no metadata mechanism at
// all, on neither PlcValue nor PlcSubscriptionEvent.
//
// Of the two ways out - inventing a metadata mechanism in the Go SPI, or modelling value and quality
// together as a PlcStruct - this driver takes the struct, deliberately:
//
//   - It is additive by construction. A metadata channel would have to be threaded through
//     apiValues.PlcValue (or the event), which every driver in the tree implements.
//   - Metadata is ignorable, and that is the whole failure mode being avoided. With a struct, a
//     consumer which reaches for the reading has to walk past the quality to get at it; treating the
//     event value as a plain boolean or number no longer type-checks at runtime (IsBool and friends
//     all report false), so a flagged-invalid reading cannot be read as a good one by accident.
//   - The wire has more to say than value plus quality: which point of which ASDU fired (a tag may
//     be a wildcard), what type identification the station used, and whether the ASDU was a test
//     frame. All of that is load-bearing and none of it fits into a scalar.
const (
	// fieldAsduAddress is the common address of ASDU the point belongs to.
	fieldAsduAddress = "asduAddress"
	// fieldObjectAddress is the three octet information object address of the point.
	fieldObjectAddress = "informationObjectAddress"
	// fieldTypeIdentification is the numeric type identification of the ASDU the point arrived in.
	fieldTypeIdentification = "typeIdentification"
	// fieldTypeName is the same type identification spelled out.
	fieldTypeName = "typeIdentificationName"
	// fieldCauseOfTransmission is why the station sent this ASDU (spontaneous, interrogated, ...).
	fieldCauseOfTransmission = "causeOfTransmission"
	// fieldTest carries the ASDU's test bit. A frame marked as a test must never drive a decision.
	fieldTest = "test"
	// fieldNegative carries the ASDU's negative-confirmation bit.
	fieldNegative = "negative"
	// fieldValue is the reading itself.
	fieldValue = "value"
	// fieldRawValue is the undecoded wire value, where the decoded one is a conversion.
	fieldRawValue = "rawValue"
	// fieldQuality holds the quality descriptor bits the wire carried with the value.
	fieldQuality = "quality"
	// fieldTimestamp holds the point's own time tag, where it has one.
	fieldTimestamp = "timestamp"
	// fieldElapsedTime is the CP16Time2a elapsed time the protection-equipment ASDUs carry.
	fieldElapsedTime = "elapsedTimeMilliseconds"
	// fieldState spells out a two bit state code (double point, double command, regulating step).
	fieldState = "state"
	// fieldCommandQualifier holds the qualifier of a command or set-point command.
	fieldCommandQualifier = "commandQualifier"
	// fieldParameterQualifier holds the qualifier of a parameter of measured values.
	fieldParameterQualifier = "parameterQualifier"
)

// pointFields collects the named values one information object decodes into.
//
// The names are kept in insertion order. A PlcStruct is a Go map and iterating one is not
// deterministic, so the change detection a change-of-state subscription needs cannot be built from
// the finished struct - it is built here instead, while the order is still known.
type pointFields struct {
	names        []string
	values       []apiValues.PlcValue
	fingerprints []string
}

// add records a scalar field.
func (p *pointFields) add(name string, value apiValues.PlcValue) {
	p.names = append(p.names, name)
	p.values = append(p.values, value)
	p.fingerprints = append(p.fingerprints, fingerprintOf(value))
}

// fingerprintOf renders one scalar for the change-detection fingerprint. String() is used rather
// than GetString(): every scalar plc value renders deterministically through it, while GetString
// panics outright on the ones which have no string form at all (PlcNull, which is what an ASDU whose
// payload cannot be decoded reports).
func fingerprintOf(value apiValues.PlcValue) string {
	if value == nil {
		return "<nil>"
	}
	return value.String()
}

// addStruct records a nested struct field, taking its fingerprint from the nested field order
// rather than from the finished (map-backed) struct.
func (p *pointFields) addStruct(name string, nested pointFields) {
	p.names = append(p.names, name)
	p.values = append(p.values, nested.plcValue())
	p.fingerprints = append(p.fingerprints, nested.fingerprint())
}

// plcValue is the struct handed to the consumer.
func (p pointFields) plcValue() apiValues.PlcValue {
	values := make(map[string]apiValues.PlcValue, len(p.names))
	for i, name := range p.names {
		values[name] = p.values[i]
	}
	return spiValues.NewPlcStruct(values)
}

// fingerprint is a deterministic rendering of everything about the point which does not change on
// its own. It is what a change-of-state subscription compares consecutive reports of a point by.
//
// The time tag, the elapsed time and the cause of transmission are left out: they differ between two
// reports of an unchanged point (a station re-reports the same value after a general interrogation),
// and publishing on those would turn a change-of-state subscription into a firehose. The quality
// flags are in, because invalid becoming valid is a change an operator has to see.
func (p pointFields) fingerprint() string {
	var sb strings.Builder
	for i, name := range p.names {
		if isVolatileField(name) {
			continue
		}
		sb.WriteString(name)
		sb.WriteByte('=')
		sb.WriteString(p.fingerprints[i])
		sb.WriteByte(';')
	}
	return sb.String()
}

// isVolatileField reports whether a field changes between two reports of a point whose state has
// not changed.
func isVolatileField(name string) bool {
	switch name {
	case fieldTimestamp, fieldElapsedTime, fieldCauseOfTransmission:
		return true
	default:
		return false
	}
}

// decodePoint turns one information object of an ASDU into the value a subscription event carries.
// It hands back the value, the fingerprint change detection compares, and the response code the
// event reports: OK for a point which decoded, UNSUPPORTED for an ASDU type whose payload the
// generated model carries nothing decodable for.
func decodePoint(asdu readWriteModel.ASDU, informationObject readWriteModel.InformationObject) (apiValues.PlcValue, string, apiModel.PlcResponseCode) {
	typeIdentification := asdu.GetTypeIdentification()
	fields := &pointFields{}
	fields.add(fieldAsduAddress, spiValues.NewPlcUINT(asdu.GetAsduAddressField()))
	fields.add(fieldObjectAddress, spiValues.NewPlcUDINT(informationObject.GetAddress()))
	fields.add(fieldTypeIdentification, spiValues.NewPlcUSINT(uint8(typeIdentification)))
	fields.add(fieldTypeName, spiValues.NewPlcSTRING(typeIdentification.String()))
	fields.add(fieldCauseOfTransmission, spiValues.NewPlcSTRING(asdu.GetCauseOfTransmission().String()))
	fields.add(fieldTest, spiValues.NewPlcBOOL(asdu.GetTest()))
	fields.add(fieldNegative, spiValues.NewPlcBOOL(asdu.GetNegative()))

	code := decodePayload(fields, informationObject)
	return fields.plcValue(), fields.fingerprint(), code
}

// decodePayload adds the type-specific fields of one information object. The switch is on the
// concrete Go type rather than on the type identification, because the model parses the payload from
// the type identification already - so the concrete type is the type identification, checked by the
// compiler instead of by a second hand-written mapping which could disagree with the first.
func decodePayload(fields *pointFields, informationObject readWriteModel.InformationObject) apiModel.PlcResponseCode {
	switch typed := informationObject.(type) {

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Process information without a time tag
	////////////////////////////////////////////////////////////////////////////////////////////////

	case readWriteModel.InformationObjectWithoutTime_SINGLE_POINT_INFORMATION:
		addSinglePointInformation(fields, typed.GetSiq())
	case readWriteModel.InformationObjectWithoutTime_DOUBLE_POINT_INFORMATION:
		addDoublePointInformation(fields, typed.GetDiq())
	case readWriteModel.InformationObjectWithoutTime_STEP_POSITION_INFORMATION:
		addStepPosition(fields, typed.GetVti(), typed.GetQds())
	case readWriteModel.InformationObjectWithoutTime_BITSTRING_OF_32_BIT:
		addBitString(fields, typed.GetBsi())
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
	case readWriteModel.InformationObjectWithoutTime_MEASURED_VALUE_NORMALISED_VALUE:
		addNormalizedValue(fields, typed.GetNva())
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
	case readWriteModel.InformationObjectWithoutTime_MEASURED_VALUE_NORMALIZED_VALUE_WITHOUT_QUALITY_DESCRIPTOR:
		// The only measurement the protocol sends without any quality at all.
		addNormalizedValue(fields, typed.GetNva())
	case readWriteModel.InformationObjectWithoutTime_MEASURED_VALUE_SCALED_VALUE:
		fields.add(fieldValue, spiValues.NewPlcINT(typed.GetSva().GetValue()))
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
	case readWriteModel.InformationObjectWithoutTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER:
		fields.add(fieldValue, spiValues.NewPlcREAL(typed.GetValue()))
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
	case readWriteModel.InformationObjectWithoutTime_INTEGRATED_TOTALS:
		addBinaryCounterReading(fields, typed.GetBcr())
	case readWriteModel.InformationObjectWithoutTime_PACKED_SINGLE_POINT_INFORMATION_WITH_STATUS_CHANGE_DETECTION:
		addStatusChangeDetection(fields, typed.GetScd())
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Commands, which a station mirrors back as an activation confirmation
	////////////////////////////////////////////////////////////////////////////////////////////////

	case readWriteModel.InformationObjectWithoutTime_SINGLE_COMMAND:
		fields.add(fieldValue, spiValues.NewPlcBOOL(typed.GetSco().GetCommandOn()))
		fields.addStruct(fieldCommandQualifier, qualifierOfCommand(typed.GetSco().GetQoc()))
	case readWriteModel.InformationObjectWithoutTime_DOUBLE_COMMAND:
		fields.add(fieldValue, spiValues.NewPlcUSINT(typed.GetDco().GetDcs()))
		fields.add(fieldState, spiValues.NewPlcSTRING(doublePointStateName(typed.GetDco().GetDcs())))
		fields.addStruct(fieldCommandQualifier, qualifierOfCommand(typed.GetDco().GetQoc()))
	case readWriteModel.InformationObjectWithoutTime_REGULATING_STEP_COMMAND:
		fields.add(fieldValue, spiValues.NewPlcUSINT(typed.GetRco().GetRcs()))
		fields.add(fieldState, spiValues.NewPlcSTRING(regulatingStepStateName(typed.GetRco().GetRcs())))
		fields.addStruct(fieldCommandQualifier, qualifierOfCommand(typed.GetRco().GetQoc()))
	case readWriteModel.InformationObjectWithoutTime_SET_POINT_COMMAND_NORMALISED_VALUE:
		addNormalizedValue(fields, typed.GetNva())
		fields.addStruct(fieldCommandQualifier, qualifierOfSetPointCommand(typed.GetQos()))
	case readWriteModel.InformationObjectWithoutTime_SET_POINT_COMMAND_SCALED_VALUE:
		fields.add(fieldValue, spiValues.NewPlcINT(typed.GetSva().GetValue()))
		fields.addStruct(fieldCommandQualifier, qualifierOfSetPointCommand(typed.GetQos()))
	case readWriteModel.InformationObjectWithoutTime_SET_POINT_COMMAND_SHORT_FLOATING_POINT_NUMBER:
		fields.add(fieldValue, spiValues.NewPlcREAL(typed.GetValue()))
		fields.addStruct(fieldCommandQualifier, qualifierOfSetPointCommand(typed.GetQos()))
	case readWriteModel.InformationObjectWithoutTime_BITSTRING_32_BIT_COMMAND:
		addBitString(fields, typed.GetBsi())

	////////////////////////////////////////////////////////////////////////////////////////////////
	// System information
	////////////////////////////////////////////////////////////////////////////////////////////////

	case readWriteModel.InformationObjectWithoutTime_END_OF_INITIALISATION:
		fields.add(fieldValue, spiValues.NewPlcUSINT(typed.GetCoi().GetQualifier()))
		// The high bit of the COI is "initialisation after a change of local parameters", which the
		// mspec (and so the generated model) calls 'select' after the qualifier-of-command layout.
		fields.add("afterLocalParameterChange", spiValues.NewPlcBOOL(typed.GetCoi().GetSelect()))
	case readWriteModel.InformationObjectWithoutTime_INTERROGATION_COMMAND:
		fields.add(fieldValue, spiValues.NewPlcUSINT(typed.GetQoi().GetQualifierOfCommand()))
	case readWriteModel.InformationObjectWithoutTime_COUNTER_INTERROGATION_COMMAND:
		fields.add(fieldValue, spiValues.NewPlcUSINT(typed.GetQcc().GetRequest()))
		fields.add("freeze", spiValues.NewPlcUSINT(typed.GetQcc().GetFreeze()))
	case readWriteModel.InformationObjectWithoutTime_READ_COMMAND:
		// A read command carries no payload at all - the address is the whole message.
		fields.add(fieldValue, spiValues.NewPlcNULL())
	case readWriteModel.InformationObjectWithoutTime_CLOCK_SYNCHRONISATION_COMMAND:
		fields.add(fieldValue, sevenOctetDateTime(typed.GetCp56Time2a()))
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
	case readWriteModel.InformationObjectWithoutTime_TEST_COMMAND:
		fields.add(fieldValue, spiValues.NewPlcWORD(typed.GetFbp().GetPattern()))
	case readWriteModel.InformationObjectWithoutTime_RESET_PROCESS_COMMAND:
		fields.add(fieldValue, spiValues.NewPlcUSINT(typed.GetQrp().GetQualifier()))
	case readWriteModel.InformationObjectWithoutTime_DELAY_ACQUISITION_COMMAND:
		fields.add(fieldValue, spiValues.NewPlcUINT(typed.GetCp16Time2a().GetMilliseconds()))

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Parameters
	////////////////////////////////////////////////////////////////////////////////////////////////

	case readWriteModel.InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_NORMALIZED_VALUE:
		addNormalizedValue(fields, typed.GetNva())
		fields.addStruct(fieldParameterQualifier, qualifierOfParameter(typed.GetQpm()))
	case readWriteModel.InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_SCALED_VALUE:
		fields.add(fieldValue, spiValues.NewPlcINT(typed.GetSva().GetValue()))
		fields.addStruct(fieldParameterQualifier, qualifierOfParameter(typed.GetQpm()))
	case readWriteModel.InformationObjectWithoutTime_PARAMETER_OF_MEASURED_VALUES_SHORT_FLOATING_POINT_NUMBER:
		fields.add(fieldValue, spiValues.NewPlcREAL(typed.GetValue()))
		fields.addStruct(fieldParameterQualifier, qualifierOfParameter(typed.GetQpm()))
	case readWriteModel.InformationObjectWithoutTime_PARAMETER_ACTIVATION:
		fields.add(fieldValue, spiValues.NewPlcUSINT(typed.GetQpa().GetQualifier()))

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Process information with a CP24Time2a time tag
	////////////////////////////////////////////////////////////////////////////////////////////////

	case readWriteModel.InformationObjectWithTreeByteTime_SINGLE_POINT_INFORMATION:
		addSinglePointInformation(fields, typed.GetSiq())
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))
	case readWriteModel.InformationObjectWithTreeByteTime_DOUBLE_POINT_INFORMATION:
		addDoublePointInformation(fields, typed.GetDiq())
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))
	case readWriteModel.InformationObjectWithTreeByteTime_STEP_POSITION_INFORMATION:
		addStepPosition(fields, typed.GetVti(), typed.GetQds())
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))
	case readWriteModel.InformationObjectWithTreeByteTime_BITSTRING_OF_32_BIT:
		addBitString(fields, typed.GetBsi())
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))
	case readWriteModel.InformationObjectWithTreeByteTime_MEASURED_VALUE_NORMALIZED_VALUE:
		addNormalizedValue(fields, typed.GetNva())
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))
	case readWriteModel.InformationObjectWithTreeByteTime_MEASURED_VALUE_SCALED_VALUE:
		fields.add(fieldValue, spiValues.NewPlcINT(typed.GetSva().GetValue()))
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))
	case readWriteModel.InformationObjectWithTreeByteTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER:
		fields.add(fieldValue, spiValues.NewPlcREAL(typed.GetValue()))
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))
	case readWriteModel.InformationObjectWithTreeByteTime_INTEGRATED_TOTALS:
		addBinaryCounterReading(fields, typed.GetBcr())
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))
	case readWriteModel.InformationObjectWithTreeByteTime_EVENT_OF_PROTECTION_EQUIPMENT:
		// See below: the model has no SEP octet for this type, so the event state is unavailable.
		fields.add(fieldElapsedTime, spiValues.NewPlcUINT(typed.GetCp16Time2a().GetMilliseconds()))
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))
		return apiModel.PlcResponseCode_UNSUPPORTED
	case readWriteModel.InformationObjectWithTreeByteTime_PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT:
		addProtectionEvent(fields, typed.GetSep(), typed.GetQdp(), typed.GetCp16Time2a())
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))
	case readWriteModel.InformationObjectWithTreeByteTime_PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT:
		addOutputCircuitInformation(fields, typed.GetOci(), typed.GetQdp(), typed.GetCp16Time2a())
		fields.addStruct(fieldTimestamp, threeOctetTimestamp(typed.GetCp24Time2a()))

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Process information with a CP56Time2a time tag
	////////////////////////////////////////////////////////////////////////////////////////////////

	case readWriteModel.InformationObjectWithSevenByteTime_SINGLE_POINT_INFORMATION:
		addSinglePointInformation(fields, typed.GetSiq())
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
	case readWriteModel.InformationObjectWithSevenByteTime_DOUBLE_POINT_INFORMATION:
		addDoublePointInformation(fields, typed.GetDiq())
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
	case readWriteModel.InformationObjectWithSevenByteTime_STEP_POSITION_INFORMATION:
		addStepPosition(fields, typed.GetVti(), typed.GetQds())
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
	case readWriteModel.InformationObjectWithSevenByteTime_BITSTRING_OF_32_BIT:
		addBitString(fields, typed.GetBsi())
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
	case readWriteModel.InformationObjectWithSevenByteTime_MEASURED_VALUE_NORMALISED_VALUE:
		addNormalizedValue(fields, typed.GetNva())
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
	case readWriteModel.InformationObjectWithSevenByteTime_MEASURED_VALUE_SCALED_VALUE:
		fields.add(fieldValue, spiValues.NewPlcINT(typed.GetSva().GetValue()))
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
	case readWriteModel.InformationObjectWithSevenByteTime_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER:
		fields.add(fieldValue, spiValues.NewPlcREAL(typed.GetValue()))
		fields.addStruct(fieldQuality, qualityOfQds(typed.GetQds()))
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
	case readWriteModel.InformationObjectWithSevenByteTime_INTEGRATED_TOTALS:
		addBinaryCounterReading(fields, typed.GetBcr())
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
	case readWriteModel.InformationObjectWithSevenByteTime_EVENT_OF_PROTECTION_EQUIPMENT:
		// M_EP_TD_1 carries an SEP octet with the event state, but the mspec for this type lists
		// only the two time fields, so the generated model has nowhere to keep it. The elapsed time
		// and the time tag are all this driver can honestly report, and reporting a point whose
		// state is missing as OK would be exactly the "looks like a reading" failure this driver
		// avoids everywhere else.
		fields.add(fieldElapsedTime, spiValues.NewPlcUINT(typed.GetCp16Time2a().GetMilliseconds()))
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
		return apiModel.PlcResponseCode_UNSUPPORTED
	case readWriteModel.InformationObjectWithSevenByteTime_PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT:
		addProtectionEvent(fields, typed.GetSep(), typed.GetQdp(), typed.GetCp16Time2a())
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))
	case readWriteModel.InformationObjectWithSevenByteTime_PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT:
		addOutputCircuitInformation(fields, typed.GetOci(), typed.GetQdp(), typed.GetCp16Time2a())
		fields.addStruct(fieldTimestamp, sevenOctetTimestamp(typed.GetCp56Time2a()))

	default:
		// The file transfer types (FILE_READY, SECTION_READY, SEGMENT, DIRECTORY and friends) parse,
		// but every one of their payload types - NameOfFile, LengthOfFile, NameOfSection, Checksum,
		// the qualifiers - is an empty type in the mspec, so there is not a single field to read out
		// of them. Anything else here is a type identification the model has no case for.
		fields.add(fieldValue, spiValues.NewPlcNULL())
		return apiModel.PlcResponseCode_UNSUPPORTED
	}
	return apiModel.PlcResponseCode_OK
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Payload decoders
////////////////////////////////////////////////////////////////////////////////////////////////////

// addSinglePointInformation decodes an SIQ: one status bit plus its four quality flags.
func addSinglePointInformation(fields *pointFields, siq readWriteModel.SinglePointInformation) {
	fields.add(fieldValue, spiValues.NewPlcBOOL(siq.GetStausOn()))
	quality := pointFields{}
	quality.add("invalid", spiValues.NewPlcBOOL(siq.GetInvalid()))
	quality.add("notTopical", spiValues.NewPlcBOOL(siq.GetNotTopical()))
	quality.add("substituted", spiValues.NewPlcBOOL(siq.GetSubstituted()))
	quality.add("blocked", spiValues.NewPlcBOOL(siq.GetBlocked()))
	fields.addStruct(fieldQuality, quality)
}

// addDoublePointInformation decodes a DIQ. The two bit DPI code is a state, not a pair of booleans:
// 0 and 3 both mean the station cannot tell, which is why plc4j's rendering as a list of the two
// bits is worse than useless - it makes "indeterminate" look like a reading.
func addDoublePointInformation(fields *pointFields, diq readWriteModel.DoublePointInformation) {
	fields.add(fieldValue, spiValues.NewPlcUSINT(diq.GetDpiCode()))
	fields.add(fieldState, spiValues.NewPlcSTRING(doublePointStateName(diq.GetDpiCode())))
	quality := pointFields{}
	quality.add("invalid", spiValues.NewPlcBOOL(diq.GetInvalid()))
	quality.add("notTopical", spiValues.NewPlcBOOL(diq.GetNotTopical()))
	quality.add("substituted", spiValues.NewPlcBOOL(diq.GetSubstituted()))
	quality.add("blocked", spiValues.NewPlcBOOL(diq.GetBlocked()))
	// A DPI code of 0 or 3 is the station saying it does not know the state.
	quality.add("indeterminate", spiValues.NewPlcBOOL(diq.GetDpiCode() == 0 || diq.GetDpiCode() == 3))
	fields.addStruct(fieldQuality, quality)
}

// addStepPosition decodes a VTI plus its QDS. The seven bit value is a two's complement step
// position in the range -64..63 (IEC 60870-5-101, 7.2.6.5), which plc4j hands out as a raw unsigned
// byte - a transformer tap on the negative side then reads as a large positive number.
func addStepPosition(fields *pointFields, vti readWriteModel.ValueWithTransientStateIndication, qds readWriteModel.QualityDescriptor) {
	fields.add(fieldValue, spiValues.NewPlcSINT(signedSevenBit(vti.GetValue())))
	fields.add(fieldRawValue, spiValues.NewPlcUSINT(vti.GetValue()))
	quality := qualityOfQds(qds)
	// The transient bit says the equipment is still moving, so the position is on its way somewhere.
	quality.add("transientState", spiValues.NewPlcBOOL(vti.GetTransientState()))
	fields.addStruct(fieldQuality, quality)
}

// addBitString decodes a BSI, which is 32 free bits the station gives no meaning to. plc4j drops it
// on the floor (its processBinaryStateInformation returns null).
func addBitString(fields *pointFields, bsi readWriteModel.BinaryStateInformation) {
	fields.add(fieldValue, spiValues.NewPlcDWORD(bsi.GetBits()))
}

// addNormalizedValue decodes an NVA. A normalized value is a two's complement fraction of full
// scale in [-1, 1) (IEC 60870-5-101, 7.2.6.6); the generated model hands out the raw 16 bits, and
// plc4j passes those straight on as an unsigned integer, which turns every negative measurement into
// a number near 65535.
func addNormalizedValue(fields *pointFields, nva readWriteModel.NormalizedValue) {
	raw := int16(nva.GetValue())
	fields.add(fieldValue, spiValues.NewPlcREAL(float32(raw)/32768.0))
	fields.add(fieldRawValue, spiValues.NewPlcINT(raw))
}

// addBinaryCounterReading decodes a BCR.
//
// Careful with the naming here: the mspec reads the qualifier byte most significant bit first and
// calls the first bit 'counterValid', but bit 7 of that byte is IV - which is *set* when the reading
// is invalid (IEC 60870-5-101, 7.2.6.9). So the model's GetCounterValid is the invalid flag, and
// plc4j publishes it under the name "counterValid", inverting its meaning: a consumer reading
// counterValid=true is looking at a counter the station has just declared broken.
func addBinaryCounterReading(fields *pointFields, bcr readWriteModel.BinaryCounterReading) {
	fields.add(fieldValue, spiValues.NewPlcUDINT(bcr.GetCounterValue()))
	quality := pointFields{}
	quality.add("invalid", spiValues.NewPlcBOOL(bcr.GetCounterValid()))
	quality.add("counterAdjusted", spiValues.NewPlcBOOL(bcr.GetCounterAdjusted()))
	quality.add("carry", spiValues.NewPlcBOOL(bcr.GetCarry()))
	quality.add("sequenceNumber", spiValues.NewPlcUSINT(bcr.GetSequenceNumber()))
	fields.addStruct(fieldQuality, quality)
}

// addStatusChangeDetection decodes an SCD, which packs 16 status bits and the 16 bits saying which
// of them have changed since the last report into one 32 bit word. plc4j drops it (its
// processStatusChangeDetection returns null).
func addStatusChangeDetection(fields *pointFields, scd readWriteModel.StatusChangeDetection) {
	bits := scd.GetBits()
	fields.add(fieldValue, spiValues.NewPlcDWORD(bits))
	fields.add("status", spiValues.NewPlcWORD(uint16(bits&0xFFFF)))
	fields.add("changeDetection", spiValues.NewPlcWORD(uint16(bits>>16)))
}

// addProtectionEvent decodes the start-events ASDUs: the event state, the elapsed time, and the two
// quality descriptors the wire carries. The flags of the two descriptors have the same five names
// and the same meaning, so they are merged: a point is invalid when either descriptor says so.
//
// The mspec names the first field SEP (single event of protection equipment) where the standard has
// SPE (start events of protection equipment); the octets line up either way, and the generated model
// is what it is.
func addProtectionEvent(fields *pointFields, sep readWriteModel.SingleEventOfProtectionEquipment, qdp readWriteModel.QualityDescriptorForPointsOfProtectionEquipment, elapsed readWriteModel.TwoOctetBinaryTime) {
	fields.add(fieldValue, spiValues.NewPlcUSINT(sep.GetEventState()))
	fields.add(fieldElapsedTime, spiValues.NewPlcUINT(elapsed.GetMilliseconds()))
	quality := pointFields{}
	quality.add("invalid", spiValues.NewPlcBOOL(sep.GetInvalid() || qdp.GetInvalid()))
	quality.add("notTopical", spiValues.NewPlcBOOL(sep.GetNotTopical() || qdp.GetNotTopical()))
	quality.add("substituted", spiValues.NewPlcBOOL(sep.GetSubstituted() || qdp.GetSubstituted()))
	quality.add("blocked", spiValues.NewPlcBOOL(sep.GetBlocked() || qdp.GetBlocked()))
	quality.add("elapsedTimeInvalid", spiValues.NewPlcBOOL(sep.GetElapsedTimeInvalid() || qdp.GetElapsedTimeInvalid()))
	fields.addStruct(fieldQuality, quality)
}

// addOutputCircuitInformation decodes an OCI plus its QDP: which phases of the protection equipment
// have operated.
func addOutputCircuitInformation(fields *pointFields, oci readWriteModel.OutputCircuitInformation, qdp readWriteModel.QualityDescriptorForPointsOfProtectionEquipment, elapsed readWriteModel.TwoOctetBinaryTime) {
	circuit := pointFields{}
	circuit.add("generalStartOfOperation", spiValues.NewPlcBOOL(oci.GetGeneralStartOfOperation()))
	circuit.add("stateOfOperationPhaseL1", spiValues.NewPlcBOOL(oci.GetStateOfOperationPhaseL1()))
	circuit.add("stateOfOperationPhaseL2", spiValues.NewPlcBOOL(oci.GetStateOfOperationPhaseL2()))
	circuit.add("stateOfOperationPhaseL3", spiValues.NewPlcBOOL(oci.GetStateOfOperationPhaseL3()))
	fields.addStruct(fieldValue, circuit)
	fields.add(fieldElapsedTime, spiValues.NewPlcUINT(elapsed.GetMilliseconds()))
	fields.addStruct(fieldQuality, qualityOfQdp(qdp))
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Quality descriptors and qualifiers
////////////////////////////////////////////////////////////////////////////////////////////////////

// qualityOfQds decodes a QDS, the descriptor most measurements travel with.
func qualityOfQds(qds readWriteModel.QualityDescriptor) pointFields {
	quality := pointFields{}
	quality.add("invalid", spiValues.NewPlcBOOL(qds.GetInvalid()))
	quality.add("notTopical", spiValues.NewPlcBOOL(qds.GetNotTopical()))
	quality.add("substituted", spiValues.NewPlcBOOL(qds.GetSubstituted()))
	quality.add("blocked", spiValues.NewPlcBOOL(qds.GetBlocked()))
	quality.add("overflow", spiValues.NewPlcBOOL(qds.GetOverflow()))
	return quality
}

// qualityOfQdp decodes a QDP, the descriptor the protection-equipment ASDUs travel with.
func qualityOfQdp(qdp readWriteModel.QualityDescriptorForPointsOfProtectionEquipment) pointFields {
	quality := pointFields{}
	quality.add("invalid", spiValues.NewPlcBOOL(qdp.GetInvalid()))
	quality.add("notTopical", spiValues.NewPlcBOOL(qdp.GetNotTopical()))
	quality.add("substituted", spiValues.NewPlcBOOL(qdp.GetSubstituted()))
	quality.add("blocked", spiValues.NewPlcBOOL(qdp.GetBlocked()))
	quality.add("elapsedTimeInvalid", spiValues.NewPlcBOOL(qdp.GetElapsedTimeInvalid()))
	return quality
}

// qualifierOfCommand decodes a QOC: whether the command was a select or an execute, and how it is
// to be carried out.
func qualifierOfCommand(qoc readWriteModel.QualifierOfCommand) pointFields {
	qualifier := pointFields{}
	qualifier.add("select", spiValues.NewPlcBOOL(qoc.GetSelect()))
	qualifier.add("qualifier", spiValues.NewPlcUSINT(qoc.GetQualifier()))
	return qualifier
}

// qualifierOfSetPointCommand decodes a QOS.
func qualifierOfSetPointCommand(qos readWriteModel.QualifierOfSetPointCommand) pointFields {
	qualifier := pointFields{}
	qualifier.add("select", spiValues.NewPlcBOOL(qos.GetSelect()))
	qualifier.add("qualifier", spiValues.NewPlcUSINT(qos.GetQualifier()))
	return qualifier
}

// qualifierOfParameter decodes a QPM.
func qualifierOfParameter(qpm readWriteModel.QualifierOfParameterOfMeasuredValues) pointFields {
	qualifier := pointFields{}
	qualifier.add("parameterInOperation", spiValues.NewPlcBOOL(qpm.GetParameterInOperation()))
	qualifier.add("localParameterChange", spiValues.NewPlcBOOL(qpm.GetLocalParameterChange()))
	qualifier.add("kindOfParameter", spiValues.NewPlcUSINT(qpm.GetKindOfParameter()))
	return qualifier
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Time tags
////////////////////////////////////////////////////////////////////////////////////////////////////

// sevenOctetTimestamp decodes a CP56Time2a into the absolute time it spells plus the flags which say
// how far it can be trusted.
func sevenOctetTimestamp(cp56 readWriteModel.SevenOctetBinaryTime) pointFields {
	timestamp := pointFields{}
	timestamp.add("dateTime", sevenOctetDateTime(cp56))
	timestamp.add("invalid", spiValues.NewPlcBOOL(cp56.GetInvalid()))
	timestamp.add("substituted", spiValues.NewPlcBOOL(cp56.GetSubstituted()))
	timestamp.add("daylightSaving", spiValues.NewPlcBOOL(cp56.GetDaylightSaving()))
	timestamp.add("dayOfWeek", spiValues.NewPlcUSINT(cp56.GetDayOfWeek()))
	return timestamp
}

// sevenOctetDateTime is the reading of the station's own clock, taken at face value.
//
// The protocol carries no timezone: a CP56Time2a is whatever the station's clock says, with a single
// bit telling whether daylight saving was in force. plc4j reinterprets it in the JVM's default
// timezone, which silently shifts every timestamp whenever the station and the client sit in
// different zones (and the DST bit is applied on top of the client's DST rules, not the station's).
// Handing the reading over unshifted, in UTC, keeps the station's own numbers intact and leaves the
// zone - which only the operator knows - to the caller.
func sevenOctetDateTime(cp56 readWriteModel.SevenOctetBinaryTime) apiValues.PlcValue {
	milliseconds := int(cp56.GetMilliseconds())
	return spiValues.NewPlcDATE_AND_TIME(time.Date(
		2000+int(cp56.GetYear()),
		time.Month(cp56.GetMonth()),
		int(cp56.GetDay()),
		int(cp56.GetHour()),
		int(cp56.GetMinutes()),
		milliseconds/1000,
		(milliseconds%1000)*int(time.Millisecond),
		time.UTC))
}

// threeOctetTimestamp decodes a CP24Time2a. It carries only the minute of the hour and the
// millisecond within that minute, so there is no absolute time in it at all: plc4j fills the missing
// year, month, day and hour in from the client's clock, which invents four fields' worth of data and
// jumps an hour whenever a report crosses the hour boundary. The parts the wire really carries are
// reported as they are, and the caller can anchor them against the arrival time if it wants to.
func threeOctetTimestamp(cp24 readWriteModel.ThreeOctetBinaryTime) pointFields {
	milliseconds := cp24.GetMilliseconds()
	timestamp := pointFields{}
	timestamp.add("minutes", spiValues.NewPlcUSINT(cp24.GetMinutes()))
	timestamp.add("seconds", spiValues.NewPlcUSINT(uint8(milliseconds/1000)))
	timestamp.add("milliseconds", spiValues.NewPlcUINT(milliseconds))
	timestamp.add("invalid", spiValues.NewPlcBOOL(cp24.GetInvalid()))
	return timestamp
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Small conversions
////////////////////////////////////////////////////////////////////////////////////////////////////

// signedSevenBit reinterprets a seven bit two's complement number, which the generated model hands
// out as a plain byte.
func signedSevenBit(value uint8) int8 {
	value &= 0x7F
	if value >= 0x40 {
		return int8(int16(value) - 0x80)
	}
	return int8(value)
}

// doublePointStateName spells out a two bit double point / double command state code.
func doublePointStateName(code uint8) string {
	switch code & 0x03 {
	case 1:
		return "OFF"
	case 2:
		return "ON"
	default:
		// 0 and 3 both mean the station cannot tell which of the two contacts is closed.
		return "INDETERMINATE"
	}
}

// regulatingStepStateName spells out a two bit regulating step command code.
func regulatingStepStateName(code uint8) string {
	switch code & 0x03 {
	case 1:
		return "NEXT_STEP_LOWER"
	case 2:
		return "NEXT_STEP_HIGHER"
	default:
		return "NOT_PERMITTED"
	}
}
