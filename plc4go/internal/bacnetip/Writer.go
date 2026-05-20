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

package bacnetip

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
)

// Writer issues BACnet WriteProperty / WritePropertyMultiple confirmed-service
// requests and waits for SimpleAck / Error / Reject / Abort responses.
type Writer struct {
	invokeIdGenerator *InvokeIdGenerator
	messageCodec      spi.MessageCodec
	tm                transactions.RequestTransactionManager
	driverContext     DriverContext

	wg sync.WaitGroup

	log zerolog.Logger
}

func NewWriter(
	invokeIdGenerator *InvokeIdGenerator,
	messageCodec spi.MessageCodec,
	tm transactions.RequestTransactionManager,
	driverContext DriverContext,
	_options ...options.WithOption,
) *Writer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Writer{
		invokeIdGenerator: invokeIdGenerator,
		messageCodec:      messageCodec,
		tm:                tm,
		driverContext:     driverContext,
		log:               customLogger,
	}
}

func (m *Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	m.log.Trace().Msg("Writing")
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	m.wg.Go(func() {
		defer func() {
			if r := recover(); r != nil {
				result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", r, debug.Stack()))
			}
		}()
		tagNames := writeRequest.GetTagNames()
		if len(tagNames) == 0 {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.New("at least one tag required"))
			return
		}

		serviceRequest, err := m.buildServiceRequest(writeRequest)
		if err != nil {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrap(err, "Error building WriteProperty request"))
			return
		}

		invokeId := m.invokeIdGenerator.getAndIncrement()
		apdu := readWriteModel.NewAPDUConfirmedRequest(
			false,
			false,
			true,
			m.driverContext.maxSegmentsAccepted,
			m.driverContext.maxApduLengthAccepted,
			invokeId,
			nil,
			nil,
			serviceRequest,
			nil,
			nil,
		)

		transaction := m.tm.StartTransaction("write")
		transaction.Submit("writeOperation", func(transactionContext context.Context, transaction transactions.RequestTransaction) {
			ctx, cancel := context.WithCancel(ctx)
			context.AfterFunc(transactionContext, cancel)

			err := m.messageCodec.SendRequest(ctx, "write", apdu, func(message spi.Message) bool {
				return m.acceptsResponse(message, invokeId)
			}, func(message spi.Message) error {
				bvlc := message.(readWriteModel.BVLC)
				responseApdu := bvlc.(interface{ GetNpdu() readWriteModel.NPDU }).GetNpdu().GetApdu()
				writeResponse := m.toPlcWriteResponse(responseApdu, writeRequest)
				result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, writeResponse, nil)
				return transaction.EndRequest()
			}, func(err error) error {
				result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrap(err, "got timeout while waiting for write response"))
				return transaction.EndRequest()
			})
			if err != nil {
				result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrap(err, "error sending message"))
				if failErr := transaction.FailRequest(err); failErr != nil {
					m.log.Debug().Err(failErr).Msg("Error failing request")
				}
			}
		})
	})
	return result
}

// buildServiceRequest assembles a single WriteProperty (one tag, no multi-prop)
// or a WritePropertyMultiple (>1 tag or any tag with >1 property) request from
// the PlcWriteRequest.
func (m *Writer) buildServiceRequest(writeRequest apiModel.PlcWriteRequest) (readWriteModel.BACnetConfirmedServiceRequest, error) {
	tagNames := writeRequest.GetTagNames()

	if len(tagNames) == 1 {
		tag, ok := writeRequest.GetTag(tagNames[0]).(BacNetPlcTag)
		if !ok {
			return nil, errors.Errorf("tag %s is not a BACnet tag", tagNames[0])
		}
		if len(tag.GetProperties()) == 1 {
			return m.buildSingleWriteProperty(tag, writeRequest.GetValue(tagNames[0]))
		}
	}

	// Multi: collapse all (tag, properties...) into a WritePropertyMultiple
	// with one BACnetWriteAccessSpecification per tag and one
	// BACnetPropertyWriteDefinition per (tag, property) pair. WritePriority is
	// taken from the same nil/!nil dance as WriteProperty.
	var specs []readWriteModel.BACnetWriteAccessSpecification
	for _, tagName := range tagNames {
		tag, ok := writeRequest.GetTag(tagName).(BacNetPlcTag)
		if !ok {
			return nil, errors.Errorf("tag %s is not a BACnet tag", tagName)
		}
		objectIdTag := readWriteModel.CreateBACnetContextTagObjectIdentifier(0, tag.GetObjectId().getId(), tag.GetObjectId().ObjectIdInstance)
		var defs []readWriteModel.BACnetPropertyWriteDefinition
		for _, prop := range tag.GetProperties() {
			plcValue := writeRequest.GetValue(tagName)
			appTag, err := plcValueToApplicationTag(plcValue.(apiValues.PlcValue), hintForProperty(prop.getId()))
			if err != nil {
				return nil, errors.Wrapf(err, "tag %s property %s", tagName, prop.String())
			}
			cd := constructedDataFromAppTag(appTag)
			propId := readWriteModel.CreateBACnetPropertyIdentifierTagged(2, prop.getId())
			var arrayIndex readWriteModel.BACnetContextTagUnsignedInteger
			if prop.ArrayIndex != nil {
				arrayIndex = readWriteModel.CreateBACnetContextTagUnsignedInteger(3, *prop.ArrayIndex)
			}
			defs = append(defs, readWriteModel.NewBACnetPropertyWriteDefinition(propId, arrayIndex, cd, nil))
		}
		specs = append(specs, readWriteModel.NewBACnetWriteAccessSpecification(
			objectIdTag,
			readWriteModel.CreateBACnetOpeningTag(1),
			defs,
			readWriteModel.CreateBACnetClosingTag(1),
		))
	}
	return readWriteModel.NewBACnetConfirmedServiceRequestWritePropertyMultiple(0, specs), nil
}

func (m *Writer) buildSingleWriteProperty(tag BacNetPlcTag, plcValue apiValues.PlcValue) (readWriteModel.BACnetConfirmedServiceRequest, error) {
	if plcValue == nil {
		return nil, errors.New("nil PlcValue")
	}
	prop := tag.GetProperties()[0]
	appTag, err := plcValueToApplicationTag(plcValue, hintForProperty(prop.getId()))
	if err != nil {
		return nil, err
	}
	objectIdTag := readWriteModel.CreateBACnetContextTagObjectIdentifier(0, tag.GetObjectId().getId(), tag.GetObjectId().ObjectIdInstance)
	propId := readWriteModel.CreateBACnetPropertyIdentifierTagged(1, prop.getId())
	var arrayIndex readWriteModel.BACnetContextTagUnsignedInteger
	if prop.ArrayIndex != nil {
		arrayIndex = readWriteModel.CreateBACnetContextTagUnsignedInteger(2, *prop.ArrayIndex)
	}
	cd := constructedDataFromAppTag(appTag)
	return readWriteModel.NewBACnetConfirmedServiceRequestWriteProperty(0, objectIdTag, propId, arrayIndex, cd, nil), nil
}

// constructedDataFromAppTag wraps a single ApplicationTag into a generic
// ConstructedDataUnspecified payload (opening tag 3, one element, closing 3) so
// the wire-format builder can serialize it without per-property typed types.
func constructedDataFromAppTag(tag readWriteModel.BACnetApplicationTag) readWriteModel.BACnetConstructedData {
	header := readWriteModel.CreateBACnetTagHeaderBalanced(true, 3, 0)
	element := readWriteModel.NewBACnetConstructedDataElement(header, tag, nil, nil)
	return readWriteModel.NewBACnetConstructedDataUnspecified(
		readWriteModel.CreateBACnetOpeningTag(3),
		header,
		readWriteModel.CreateBACnetClosingTag(3),
		nil,
		[]readWriteModel.BACnetConstructedDataElement{element},
	)
}

// acceptsResponse matches the incoming message against our outstanding invoke id.
func (m *Writer) acceptsResponse(message spi.Message, invokeId uint8) bool {
	bvlc, ok := message.(readWriteModel.BVLC)
	if !ok {
		return false
	}
	npduRetriever, ok := bvlc.(interface{ GetNpdu() readWriteModel.NPDU })
	if !ok {
		return false
	}
	npdu := npduRetriever.GetNpdu()
	if npdu.GetControl().GetMessageTypeFieldPresent() {
		return false
	}
	got, err := getInvokeIdFromApdu(npdu.GetApdu())
	if err != nil {
		return false
	}
	return got == invokeId
}

// toPlcWriteResponse decodes a WriteProperty response APDU into the per-tag
// response codes. Single-tag writes get the same code for the lone tag; multi
// writes fan out errors to the tag whose write failed when the device tells us.
func (m *Writer) toPlcWriteResponse(apdu readWriteModel.APDU, request apiModel.PlcWriteRequest) apiModel.PlcWriteResponse {
	codes := map[string]apiModel.PlcResponseCode{}
	switch apdu := apdu.(type) {
	case readWriteModel.APDUSimpleAck:
		for _, name := range request.GetTagNames() {
			codes[name] = apiModel.PlcResponseCode_OK
		}
	case readWriteModel.APDUError:
		code := mapErrorAPDU(apdu, m.log)
		for _, name := range request.GetTagNames() {
			codes[name] = code
		}
	case readWriteModel.APDUReject:
		m.log.Warn().Stringer("rejectReason", reasonOfReject(apdu)).Msg("BACnet REJECT on write")
		for _, name := range request.GetTagNames() {
			codes[name] = apiModel.PlcResponseCode_INVALID_DATA
		}
	case readWriteModel.APDUAbort:
		reason := apdu.GetAbortReason().GetValue()
		code := apiModel.PlcResponseCode_INTERNAL_ERROR
		if reason == readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED {
			code = apiModel.PlcResponseCode_UNSUPPORTED
		}
		m.log.Warn().Stringer("abortReason", &reason).Msg("BACnet ABORT on write")
		for _, name := range request.GetTagNames() {
			codes[name] = code
		}
	default:
		m.log.Warn().Type("apdu", apdu).Msg("unsupported APDU type on write response")
		for _, name := range request.GetTagNames() {
			codes[name] = apiModel.PlcResponseCode_REMOTE_ERROR
		}
	}
	return spiModel.NewDefaultPlcWriteResponse(request, codes)
}
