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

package simulated

import (
	"context"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"runtime/debug"
	"strconv"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

type Reader struct {
	device  *Device
	options map[string][]string
	tracer  tracer.Tracer

	log zerolog.Logger
}

func NewReader(device *Device, readerOptions map[string][]string, tracer tracer.Tracer, _options ...options.WithOption) *Reader {
	return &Reader{
		device:  device,
		options: readerOptions,
		tracer:  tracer,

		log: options.ExtractCustomLogger(_options...),
	}
}

func (r *Reader) Read(_ context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	ch := make(chan apiModel.PlcReadRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		var txId string
		if r.tracer != nil {
			txId = r.tracer.AddTransactionalStartTrace("read", "started")
		}
		// Possibly add a delay.
		if delayString, ok := r.options["readDelay"]; ok {
			if len(delayString) == 1 {
				delay, err := strconv.Atoi(delayString[0])
				if err == nil {
					time.Sleep(time.Duration(delay) * time.Millisecond)
				}
			}
		}

		// Process the request
		responseCodes := make(map[string]apiModel.PlcResponseCode)
		responseValues := make(map[string]apiValues.PlcValue)
		for _, tagName := range readRequest.GetTagNames() {
			tag := readRequest.GetTag(tagName)
			simulatedTagVar, ok := tag.(simulatedTag)
			if !ok {
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				responseValues[tagName] = nil
			} else {
				value := r.device.Get(simulatedTagVar)
				if value == nil {
					responseCodes[tagName] = apiModel.PlcResponseCode_NOT_FOUND
					responseValues[tagName] = nil
				} else {
					responseCodes[tagName] = apiModel.PlcResponseCode_OK
					responseValues[tagName] = *value
				}
			}
		}

		if r.tracer != nil {
			r.tracer.AddTransactionalTrace(txId, "read", "success")
		}
		// Emit the response
		ch <- spiModel.NewDefaultPlcReadRequestResult(
			readRequest,
			spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, responseValues),
			nil,
		)
	}()
	return ch
}
