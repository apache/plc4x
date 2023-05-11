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
	"context"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (m *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(m.GetPlcTagHandler(), m.GetPlcValueHandler(), m)
}

func (m *Connection) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	log.Trace().Msg("Writing")
	result := make(chan apiModel.PlcWriteRequestResult)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- internalModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Errorf("panic-ed %v", err))
			}
		}()
		if len(writeRequest.GetTagNames()) <= 1 {
			m.singleWrite(ctx, writeRequest, result)
		} else {
			m.multiWrite(ctx, writeRequest, result)
		}
	}()
	return result
}

func (m *Connection) singleWrite(ctx context.Context, writeRequest apiModel.PlcWriteRequest, result chan apiModel.PlcWriteRequestResult) {
	if len(writeRequest.GetTagNames()) != 1 {
		result <- internalModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.New("this part of the ads driver only supports single-item requests"))
		log.Debug().Msgf("this part of the ads driver only supports single-item requests. Got %d tags", len(writeRequest.GetTagNames()))
		return
	}

	// Here we can be sure that we're only handling a single request.
	tagName := writeRequest.GetTagNames()[0]
	tag := writeRequest.GetTag(tagName)
	if model.NeedsResolving(tag) {
		adsField, err := model.CastToSymbolicPlcTagFromPlcTag(tag)
		if err != nil {
			result <- &internalModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid tag item type"),
			}
			log.Debug().Msgf("Invalid tag item type %T", tag)
			return
		}
		// Replace the symbolic tag with a direct one
		tag, err = m.resolveSymbolicTag(ctx, adsField)
		if err != nil {
			result <- &internalModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid tag item type"),
			}
			log.Debug().Msgf("Invalid tag item type %T", tag)
			return
		}
	}
	directAdsTag, ok := tag.(*model.DirectPlcTag)
	if !ok {
		result <- &internalModel.DefaultPlcWriteRequestResult{
			Request:  writeRequest,
			Response: nil,
			Err:      errors.New("invalid tag item type"),
		}
		log.Debug().Msgf("Invalid tag item type %T", tag)
		return
	}

	// Get the value from the request and serialize it to a byte array
	value := writeRequest.GetValue(tagName)
	io := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	err := m.serializePlcValue(directAdsTag.DataType, directAdsTag.GetArrayInfo(), value, io)
	if err != nil {
		result <- &internalModel.DefaultPlcWriteRequestResult{
			Request:  writeRequest,
			Response: nil,
			Err:      errors.Wrap(err, "error serializing plc value"),
		}
		return
	}
	data := io.GetBytes()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- internalModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Errorf("panic-ed %v", err))
			}
		}()
		response, err := m.ExecuteAdsWriteRequest(ctx, directAdsTag.IndexGroup, directAdsTag.IndexOffset, data)
		if err != nil {
			result <- internalModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrap(err, "got error executing the write request"))
			return
		}

		if response.GetErrorCode() != 0x00000000 {
			// TODO: Handle this ...
		}

		responseCodes := map[string]apiModel.PlcResponseCode{}
		if response.GetErrorCode() != 0x00000000 {
			// TODO: Correctly handle this.
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
		} else {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
		}
		// Return the response to the caller.
		result <- &internalModel.DefaultPlcWriteRequestResult{
			Request:  writeRequest,
			Response: internalModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes),
			Err:      nil,
		}
	}()
}

func (m *Connection) multiWrite(ctx context.Context, writeRequest apiModel.PlcWriteRequest, result chan apiModel.PlcWriteRequestResult) {
	// Calculate the size of all tags together.
	// Calculate the expected size of the response data.
	expectedResponseDataSize := uint32(0)
	directAdsTags := map[string]*model.DirectPlcTag{}
	requestItems := make([]driverModel.AdsMultiRequestItem, 0)
	io := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	for _, tagName := range writeRequest.GetTagNames() {
		tag := writeRequest.GetTag(tagName)
		if model.NeedsResolving(tag) {
			adsField, err := model.CastToSymbolicPlcTagFromPlcTag(tag)
			if err != nil {
				result <- &internalModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrap(err, "invalid tag item type"),
				}
				log.Debug().Msgf("Invalid tag item type %T", tag)
				return
			}
			// Replace the symbolic tag with a direct one
			tag, err = m.resolveSymbolicTag(ctx, adsField)
			if err != nil {
				result <- &internalModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrap(err, "invalid tag item type"),
				}
				log.Debug().Msgf("Invalid tag item type %T", tag)
				return
			}
		}
		directAdsTag, ok := tag.(*model.DirectPlcTag)
		if !ok {
			result <- &internalModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("invalid tag item type"),
			}
			log.Debug().Msgf("Invalid tag item type %T", tag)
			return
		}

		directAdsTags[tagName] = directAdsTag

		// Serialize the plc value
		err := m.serializePlcValue(directAdsTag.DataType, directAdsTag.GetArrayInfo(), writeRequest.GetValue(tagName), io)
		if err != nil {
			result <- &internalModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.Wrap(err, "error serializing plc value"),
			}
			return
		}

		// Size of one element.
		size := directAdsTag.DataType.GetSize()

		// Calculate how many elements in total we'll be reading.
		arraySize := uint32(1)
		if len(tag.GetArrayInfo()) > 0 {
			for _, arrayInfo := range tag.GetArrayInfo() {
				arraySize = arraySize * arrayInfo.GetSize()
			}
		}

		// Status code + payload size
		expectedResponseDataSize += 4

		requestItems = append(requestItems, driverModel.NewAdsMultiRequestItemWrite(
			directAdsTag.IndexGroup, directAdsTag.IndexOffset, size*arraySize))
	}

	response, err := m.ExecuteAdsReadWriteRequest(ctx,
		uint32(driverModel.ReservedIndexGroups_ADSIGRP_MULTIPLE_WRITE), uint32(len(directAdsTags)),
		expectedResponseDataSize, requestItems, io.GetBytes())
	if err != nil {
		result <- &internalModel.DefaultPlcWriteRequestResult{
			Request:  writeRequest,
			Response: nil,
			Err:      errors.Wrap(err, "error executing multi-item write request"),
		}
		return
	}

	if response.GetResult() != driverModel.ReturnCode_OK {
		result <- &internalModel.DefaultPlcWriteRequestResult{
			Request:  writeRequest,
			Response: nil,
			Err:      fmt.Errorf("got return result %s from remote", response.GetResult().String()),
		}
		return
	}

	rb := utils.NewReadBufferByteBased(response.GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))

	// Read in the response codes first.
	responseCodes := map[string]apiModel.PlcResponseCode{}
	for _, tagName := range writeRequest.GetTagNames() {
		returnCodeValue, err := rb.ReadUint32("returnCode", 32)
		if err != nil {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
		} else if returnCodeValue != 0x00000000 {
			// TODO: Correctly handle this.
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
		} else {
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
		}
	}

	// Return the response to the caller.
	result <- &internalModel.DefaultPlcWriteRequestResult{
		Request:  writeRequest,
		Response: internalModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes),
		Err:      nil,
	}
}

func (m *Connection) serializePlcValue(dataType driverModel.AdsDataTypeTableEntry, arrayInfo []apiModel.ArrayInfo, plcValue values.PlcValue, wb utils.WriteBufferByteBased) error {
	// Decode the data according to the information from the request
	// Based on the AdsDataTypeTableEntry in tag.DataType() parse the data
	if len(arrayInfo) > 0 {
		// This is an Array/List type.
		curArrayInfo := arrayInfo[0]
		// Do some initial checks
		if !plcValue.IsList() {
			return fmt.Errorf("expecting a plc value of type list")
		}
		plcValues := plcValue.GetList()
		if uint32(len(plcValues)) != curArrayInfo.GetSize() {
			return fmt.Errorf("expecting exactly %d items in the list", len(plcValues))
		}

		arrayItemTypeName := dataType.GetDataTypeName()[strings.Index(dataType.GetDataTypeName(), " OF ")+4:]
		arrayItemType, ok := m.driverContext.dataTypeTable[arrayItemTypeName]
		if !ok {
			return fmt.Errorf("couldn't resolve array item type %s", arrayItemTypeName)
		}

		for _, plcValue := range plcValues {
			restArrayInfo := arrayInfo[1:]
			err := m.serializePlcValue(arrayItemType, restArrayInfo, plcValue, wb)
			if err != nil {
				return errors.Wrap(err, "error encoding list item")
			}
		}
		return nil
	} else if len(dataType.GetChildren()) > 0 {
		// Do some initial checks
		if !plcValue.IsStruct() {
			return fmt.Errorf("expecting a plc value of type struct")
		}
		plcValues := plcValue.GetStruct()
		if len(plcValues) != len(dataType.GetChildren()) {
			return fmt.Errorf("expecting exactly %d children in struct, but got %d",
				len(plcValues), len(dataType.GetChildren()))
		}

		// This is a Struct type.
		startPos := uint32(wb.GetPos())
		curPos := uint32(0)
		for _, child := range dataType.GetChildren() {
			childName := child.GetPropertyName()
			childDataType, ok := m.driverContext.dataTypeTable[child.GetDataTypeName()]
			if !ok {
				return fmt.Errorf("couldn't find data type named %s for property %s of type %s", child.GetDataTypeName(), childName, dataType.GetDataTypeName())
			}
			if child.GetOffset() > curPos {
				skipBytes := child.GetOffset() - curPos
				for i := uint32(0); i < skipBytes; i++ {
					_ = wb.WriteByte("", 0x00)
				}
			}
			var childArrayInfo []apiModel.ArrayInfo
			for _, adsArrayInfo := range childDataType.GetArrayInfo() {
				childArrayInfo = append(childArrayInfo, &internalModel.DefaultArrayInfo{
					LowerBound: adsArrayInfo.GetLowerBound(),
					UpperBound: adsArrayInfo.GetUpperBound(),
				})
			}
			err := m.serializePlcValue(childDataType, childArrayInfo, plcValues[childName], wb)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error parsing propery %s of type %s", childName, dataType.GetDataTypeName()))
			}
			curPos = uint32(wb.GetPos()) - startPos
		}
		return nil
	} else {
		// This is a primitive type.
		valueType, stringLength := m.getPlcValueForAdsDataTypeTableEntry(dataType)
		if valueType == values.NULL {
			return errors.New(fmt.Sprintf("error converting %s into plc4x plc-value type", dataType.GetDataTypeName()))
		}
		adsValueType, ok := driverModel.PlcValueTypeByName(valueType.String())
		if !ok {
			return errors.New(fmt.Sprintf("error converting plc4x plc-value type %s into ads plc-value type", valueType.String()))
		}
		return driverModel.DataItemSerializeWithWriteBuffer(context.Background(), wb, plcValue, adsValueType, stringLength)
	}
}
