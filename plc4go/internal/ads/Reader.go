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

package ads

import (
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Reader struct {
	transactionIdentifier uint32
	targetAmsNetId        readWriteModel.AmsNetId
	targetAmsPort         uint16
	sourceAmsNetId        readWriteModel.AmsNetId
	sourceAmsPort         uint16
	messageCodec          spi.MessageCodec
	fieldMapping          map[SymbolicPlcField]DirectPlcField
	mappingLock           sync.Mutex
}

func NewReader(messageCodec spi.MessageCodec, targetAmsNetId readWriteModel.AmsNetId, targetAmsPort uint16, sourceAmsNetId readWriteModel.AmsNetId, sourceAmsPort uint16) *Reader {
	return &Reader{
		transactionIdentifier: 0,
		targetAmsNetId:        targetAmsNetId,
		targetAmsPort:         targetAmsPort,
		sourceAmsNetId:        sourceAmsNetId,
		sourceAmsPort:         sourceAmsPort,
		messageCodec:          messageCodec,
		fieldMapping:          make(map[SymbolicPlcField]DirectPlcField),
	}
}

func (m *Reader) Read(readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	log.Trace().Msg("Reading")
	result := make(chan model.PlcReadRequestResult)
	go func() {
		if len(readRequest.GetFieldNames()) <= 1 {
			m.singleRead(readRequest, result)
		} else {
			m.multiRead(readRequest, result)
		}
	}()
	return result
}

func (m *Reader) singleRead(readRequest model.PlcReadRequest, result chan model.PlcReadRequestResult) {
	if len(readRequest.GetFieldNames()) != 1 {
		result <- &plc4goModel.DefaultPlcReadRequestResult{
			Request:  readRequest,
			Response: nil,
			Err:      errors.New("ads only supports single-item requests"),
		}
		log.Debug().Msgf("ads only supports single-item requests. Got %d fields", len(readRequest.GetFieldNames()))
		return
	}
	// If we are requesting only one field, use a
	fieldName := readRequest.GetFieldNames()[0]
	field := readRequest.GetField(fieldName)
	if needsResolving(field) {
		adsField, err := castToSymbolicPlcFieldFromPlcField(field)
		if err != nil {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid field item type"),
			}
			log.Debug().Msgf("Invalid field item type %T", field)
			return
		}
		field, err = m.resolveField(adsField)
		if err != nil {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid field item type"),
			}
			log.Debug().Msgf("Invalid field item type %T", field)
			return
		}
	}
	adsField, err := castToDirectAdsFieldFromPlcField(field)
	if err != nil {
		result <- &plc4goModel.DefaultPlcReadRequestResult{
			Request:  readRequest,
			Response: nil,
			Err:      errors.Wrap(err, "invalid field item type"),
		}
		log.Debug().Msgf("Invalid field item type %T", field)
		return
	}
	userdata := readWriteModel.NewAmsPacket(
		m.targetAmsNetId,
		m.targetAmsPort,
		m.sourceAmsNetId,
		m.sourceAmsPort,
		readWriteModel.CommandId_ADS_READ,
		readWriteModel.NewState(false, false, false, false, false, true, false, false, false),
		0,
		0,
		nil,
	)

	readLength := uint32(adsField.Datatype.NumBytes())
	switch {
	case adsField.GetDatatype() == readWriteModel.AdsDataType_STRING:
		// If an explicit size is given with the string, use this, if not use 256
		if adsField.GetStringLength() != 0 {
			readLength = uint32(adsField.GetStringLength())
		} else {
			readLength = 256
		}
	case adsField.GetDatatype() == readWriteModel.AdsDataType_WSTRING:
		// If an explicit size is given with the string, use this, if not use 512
		if adsField.GetStringLength() != 0 {
			readLength = uint32(adsField.GetStringLength() * 2)
		} else {
			readLength = 512
		}
	default:
		readLength = uint32(adsField.Datatype.NumBytes())
	}
	userdata.Data = readWriteModel.NewAdsReadRequest(adsField.IndexGroup, adsField.IndexOffset, readLength)

	m.sendOverTheWire(userdata, readRequest, result)
}

func (m *Reader) multiRead(readRequest model.PlcReadRequest, result chan model.PlcReadRequestResult) {
	// Calculate the size of all fields together.
	// Calculate the expected size of the response data.
	expectedResponseDataSize := uint32(0)
	for _, fieldName := range readRequest.GetFieldNames() {
		field, err := castToAdsFieldFromPlcField(readRequest.GetField(fieldName))
		if err != nil {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Wrap(err, "error casting field"),
			}
			return
		}
		size := uint32(0)
		switch field.GetDatatype() {
		case readWriteModel.AdsDataType_STRING:
			// If an explicit size is given with the string, use this, if not use 256
			if field.GetStringLength() != 0 {
				size = uint32(field.GetStringLength())
			} else {
				size = 256
			}
		case readWriteModel.AdsDataType_WSTRING:
			// If an explicit size is given with the string, use this, if not use 512
			if field.GetStringLength() != 0 {
				size = uint32(field.GetStringLength() * 2)
			} else {
				size = 512
			}
		default:
			size = uint32(field.GetDatatype().NumBytes())
		}
		// Status code + payload size
		expectedResponseDataSize += 4 + (size * field.GetNumberOfElements())
	}

	userdata := readWriteModel.NewAmsPacket(
		m.targetAmsNetId,
		m.targetAmsPort,
		m.sourceAmsNetId,
		m.sourceAmsPort,
		readWriteModel.CommandId_ADS_READ_WRITE,
		readWriteModel.NewState(false, false, false, false, false, true, false, false, false),
		0,
		0,
		nil,
	)

	items := make([]readWriteModel.AdsMultiRequestItem, len(readRequest.GetFieldNames()))
	for i, fieldName := range readRequest.GetFieldNames() {
		field := readRequest.GetField(fieldName)
		if needsResolving(field) {
			adsField, err := castToSymbolicPlcFieldFromPlcField(field)
			if err != nil {
				result <- &plc4goModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: nil,
					Err:      errors.Wrap(err, "invalid field item type"),
				}
				log.Debug().Msgf("Invalid field item type %T", field)
				return
			}
			field, err = m.resolveField(adsField)
			if err != nil {
				result <- &plc4goModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: nil,
					Err:      errors.Wrap(err, "invalid field item type"),
				}
				log.Debug().Msgf("Invalid field item type %T", field)
				return
			}
		}
		adsField, err := castToDirectAdsFieldFromPlcField(field)
		if err != nil {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid field item type"),
			}
			log.Debug().Msgf("Invalid field item type %T", field)
			return
		}
		// With multi-requests, the index-group is fixed and the index offset indicates the number of elements.
		items[i] = readWriteModel.NewAdsMultiRequestItemRead(adsField.IndexGroup, adsField.IndexOffset, uint32(adsField.GetDatatype().NumBytes())*adsField.NumberOfElements)
	}
	userdata.Data = readWriteModel.NewAdsReadWriteRequest(uint32(readWriteModel.ReservedIndexGroups_ADSIGRP_MULTIPLE_READ), uint32(len(readRequest.GetFieldNames())), expectedResponseDataSize, items, nil)

	m.sendOverTheWire(userdata, readRequest, result)
}

func (m *Reader) sendOverTheWire(userdata readWriteModel.AmsPacket, readRequest model.PlcReadRequest, result chan model.PlcReadRequestResult) {
	// Calculate a new transaction identifier
	transactionIdentifier := atomic.AddUint32(&m.transactionIdentifier, 1)
	if transactionIdentifier > math.MaxUint8 {
		transactionIdentifier = 1
		atomic.StoreUint32(&m.transactionIdentifier, 1)
	}
	log.Debug().Msgf("Calculated transaction identifier %x", transactionIdentifier)
	userdata = readWriteModel.NewAmsPacket(
		userdata.GetTargetAmsNetId(),
		userdata.GetTargetAmsPort(),
		userdata.GetSourceAmsNetId(),
		userdata.GetSourceAmsPort(),
		userdata.GetCommandId(),
		userdata.GetState(),
		userdata.GetErrorCode(),
		transactionIdentifier,
		userdata.GetData(),
	)

	// Assemble the finished tcp paket
	log.Trace().Msg("Assemble tcp paket")
	amsTcpPaket := readWriteModel.NewAmsTCPPacket(userdata)

	// Send the TCP Paket over the wire
	log.Trace().Msg("Send TCP Paket")
	if err := m.messageCodec.SendRequest(
		amsTcpPaket,
		func(message spi.Message) bool {
			paket := message.(readWriteModel.AmsTCPPacket)
			return paket.GetUserdata().GetInvokeId() == transactionIdentifier
		},
		func(message spi.Message) error {
			// Convert the response into an amsTcpPaket
			log.Trace().Msg("convert response to amsTcpPaket")
			receivedAmsTcpPaket := message.(readWriteModel.AmsTCPPacket)
			// Convert the ads response into a PLC4X response
			log.Trace().Msg("convert response to PLC4X response")
			readResponse, err := m.ToPlc4xReadResponse(receivedAmsTcpPaket, readRequest)

			if err != nil {
				result <- &plc4goModel.DefaultPlcReadRequestResult{
					Request: readRequest,
					Err:     errors.Wrap(err, "Error decoding response"),
				}
				// TODO: should we return the error here?
				return nil
			}
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: readResponse,
			}
			return nil
		},
		func(err error) error {
			result <- &plc4goModel.DefaultPlcReadRequestResult{
				Request: readRequest,
				Err:     errors.Wrap(err, "got timeout while waiting for response"),
			}
			return nil
		},
		time.Second*1); err != nil {
		result <- &plc4goModel.DefaultPlcReadRequestResult{
			Request:  readRequest,
			Response: nil,
			Err:      errors.Wrap(err, "error sending message"),
		}
	}
}

func (m *Reader) resolveField(symbolicField SymbolicPlcField) (DirectPlcField, error) {
	if directPlcField, ok := m.fieldMapping[symbolicField]; ok {
		return directPlcField, nil
	}
	m.mappingLock.Lock()
	defer m.mappingLock.Unlock()
	// In case a previous one has already
	if directPlcField, ok := m.fieldMapping[symbolicField]; ok {
		return directPlcField, nil
	}
	userdata := readWriteModel.NewAmsPacket(
		m.targetAmsNetId,
		m.targetAmsPort,
		m.sourceAmsNetId,
		m.sourceAmsPort,
		readWriteModel.CommandId_ADS_READ_WRITE,
		readWriteModel.NewState(false, false, false, false, false, true, false, false, false),
		0,
		0,
		nil,
	)
	userdata.Data = readWriteModel.NewAdsReadWriteRequest(
		uint32(readWriteModel.ReservedIndexGroups_ADSIGRP_SYM_HNDBYNAME),
		0,
		4,
		nil,
		[]byte(symbolicField.SymbolicAddress+"\000"),
	)
	result := make(chan model.PlcReadRequestResult)
	go func() {
		dummyRequest := plc4goModel.NewDefaultPlcReadRequest(map[string]model.PlcField{"dummy": DirectPlcField{PlcField: PlcField{Datatype: readWriteModel.AdsDataType_UINT32}}}, []string{"dummy"}, nil, nil)
		m.sendOverTheWire(userdata, dummyRequest, result)
	}()
	// We wait synchronous for the resolution response before we can continue
	response := <-result
	if response.GetErr() != nil {
		log.Debug().Err(response.GetErr()).Msg("Error during resolve")
		return DirectPlcField{}, response.GetErr()
	}
	if response.GetResponse().GetResponseCode("dummy") != model.PlcResponseCode_OK {
		return DirectPlcField{}, errors.Errorf("Got a response error %#v", response.GetResponse().GetResponseCode("dummy"))
	}
	handle := response.GetResponse().GetValue("dummy").GetUint32()
	log.Debug().Uint32("handle", handle).Str("symbolicAddress", symbolicField.SymbolicAddress).Msg("Resolved symbolic address")
	directPlcField := DirectPlcField{
		IndexGroup:  uint32(readWriteModel.ReservedIndexGroups_ADSIGRP_SYM_VALBYHND),
		IndexOffset: handle,
		PlcField:    symbolicField.PlcField,
	}
	switch directPlcField.FieldType {
	case SymbolicAdsField:
		directPlcField.FieldType = DirectAdsField
	case SymbolicAdsStringField:
		directPlcField.FieldType = DirectAdsStringField
	}
	m.fieldMapping[symbolicField] = directPlcField
	return directPlcField, nil
}

func (m *Reader) ToPlc4xReadResponse(amsTcpPaket readWriteModel.AmsTCPPacket, readRequest model.PlcReadRequest) (model.PlcReadResponse, error) {
	var rb utils.ReadBuffer
	responseCodes := map[string]model.PlcResponseCode{}
	switch data := amsTcpPaket.GetUserdata().GetData().(type) {
	case readWriteModel.AdsReadResponse:
		rb = utils.NewLittleEndianReadBufferByteBased(data.GetData())
		for _, fieldName := range readRequest.GetFieldNames() {
			responseCodes[fieldName] = model.PlcResponseCode_OK
		}
	case readWriteModel.AdsReadWriteResponse:
		rb = utils.NewLittleEndianReadBufferByteBased(data.GetData())
		// When parsing a multi-item response, the error codes of each items come
		// in sequence and then come the values.
		for _, fieldName := range readRequest.GetFieldNames() {
			if len(readRequest.GetFieldNames()) <= 1 {
				// TODO: the comment above seems strange as there is no such spec for response codes per field so maybe this is a speciality
				break
			}
			responseCode, err := rb.ReadUint32("responseCode", 32)
			if err != nil {
				log.Error().Err(err).Str("fieldName", fieldName).Msgf("Error parsing field %s", fieldName)
				responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
				continue
			}
			val, _ := readWriteModel.ReturnCodeByValue(responseCode)
			switch val {
			case readWriteModel.ReturnCode_OK:
				responseCodes[fieldName] = model.PlcResponseCode_OK
			default:
				// TODO: Implement this a little more ...
				log.Error().Stringer("adsReturnCode", val).Msgf("Unmapped return code for %s", fieldName)
				responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
			}
		}
	default:
		return nil, errors.Errorf("unsupported response type %T", data)
	}

	plcValues := map[string]values.PlcValue{}
	// Get the field from the request
	for _, fieldName := range readRequest.GetFieldNames() {
		log.Debug().Msgf("get a field from request with name %s", fieldName)
		field, err := castToAdsFieldFromPlcField(readRequest.GetField(fieldName))
		if err != nil {
			return nil, errors.Wrap(err, "error casting to ads-field")
		}

		// Decode the data according to the information from the request
		log.Trace().Msg("decode data")
		value, err := readWriteModel.DataItemParse(rb, field.GetDatatype().DataFormatName(), field.GetStringLength())
		if err != nil {
			log.Error().Err(err).Msg("Error parsing data item")
			responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
			continue
		}
		plcValues[fieldName] = value
		responseCodes[fieldName] = model.PlcResponseCode_OK
	}

	// Return the response
	log.Trace().Msg("Returning the response")
	return plc4goModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}
