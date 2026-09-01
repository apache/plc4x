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

package knxnetip

import (
	"context"
	"runtime/debug"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Writer implements group-address writes. Everything else (device properties and device
// memory) is read-only in this driver, so only GroupAddressTag is supported here.
type Writer struct {
	connection *Connection

	log zerolog.Logger
}

func NewWriter(connection *Connection, _options ...options.WithOption) Writer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return Writer{
		connection: connection,
		log:        customLogger,
	}
}

// Write sends a GroupValueWrite for every tag of the request. The returned channel is
// always completed exactly once, no matter how many tags the request contains and no
// matter if the individual writes succeed, fail or time out.
func (m Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	m.connection.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				// Without this the caller would wait for a response which never comes.
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil,
					errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()

		responseCodes := map[string]apiModel.PlcResponseCode{}
		for _, tagName := range writeRequest.GetTagNames() {
			responseCodes[tagName] = m.writeTag(ctx, writeRequest, tagName)
		}

		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(
			writeRequest,
			spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes),
			nil,
		))
	})
	return result
}

// writeTag writes a single tag of the request and maps whatever happened to a response code.
func (m Writer) writeTag(ctx context.Context, writeRequest apiModel.PlcWriteRequest, tagName string) apiModel.PlcResponseCode {
	groupAddressTag, err := CastToGroupAddressTagFromPlcTag(writeRequest.GetTag(tagName))
	if err != nil {
		m.log.Debug().Err(err).Str("tagName", tagName).Msg("only group-addresses can be written")
		return apiModel.PlcResponseCode_INVALID_ADDRESS
	}

	// Writing to a pattern would mean writing to an unknown number of devices, so unlike a
	// read a write only accepts tags which resolve to exactly one group address.
	rawAddresses, err := m.resolveGroupAddresses(groupAddressTag)
	if err != nil {
		m.log.Debug().Err(err).Str("tagName", tagName).Msg("error resolving addresses")
		return apiModel.PlcResponseCode_INVALID_ADDRESS
	}
	if len(rawAddresses) != 1 {
		m.log.Debug().Str("tagName", tagName).Int("numAddresses", len(rawAddresses)).
			Msg("a write tag has to address exactly one group address")
		return apiModel.PlcResponseCode_INVALID_ADDRESS
	}
	numericAddress := rawAddresses[0]
	groupAddress := []byte{byte(numericAddress >> 8), byte(numericAddress & 0xFF)}

	value := writeRequest.GetValue(tagName)
	if value == nil {
		m.log.Debug().Str("tagName", tagName).Msg("no value to write")
		return apiModel.PlcResponseCode_INVALID_DATA
	}

	writeResults := m.connection.WriteGroupAddress(ctx, groupAddress, groupAddressTag.GetTagType(), value)
	select {
	case writeResult := <-writeResults:
		if writeResult.err != nil {
			m.log.Debug().Err(writeResult.err).Str("tagName", tagName).Msg("error writing group address")
			return errorToResponseCode(writeResult.err)
		}
		return apiModel.PlcResponseCode_OK
	case <-ctx.Done():
		// Without this a write would block forever if the gateway never answers.
		m.log.Debug().Err(ctx.Err()).Str("tagName", tagName).Msg("context done while writing group address")
		return apiModel.PlcResponseCode_REQUEST_TIMEOUT
	}
}

// resolveGroupAddresses maps a (possibly pattern-based) group-address tag to the numeric
// group addresses it refers to, reusing the resolution the Reader implements.
func (m Writer) resolveGroupAddresses(tag GroupAddressTag) ([]uint16, error) {
	return NewReader(m.connection, options.WithCustomLogger(m.log)).resolveAddresses(tag)
}

// errorToResponseCode tells a timeout apart from any other failure, as the api has a
// dedicated response code for it.
func errorToResponseCode(err error) apiModel.PlcResponseCode {
	var timeoutError utils.TimeoutError
	if errors.As(err, &timeoutError) ||
		errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) {
		return apiModel.PlcResponseCode_REQUEST_TIMEOUT
	}
	return apiModel.PlcResponseCode_INTERNAL_ERROR
}
