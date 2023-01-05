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
	"strconv"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi"
	model2 "github.com/apache/plc4x/plc4go/spi/model"
)

type Reader struct {
	device  *Device
	options map[string][]string
	tracer  *spi.Tracer
}

func NewReader(device *Device, options map[string][]string, tracer *spi.Tracer) Reader {
	return Reader{
		device:  device,
		options: options,
		tracer:  tracer,
	}
}

func (r Reader) Read(ctx context.Context, readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	// TODO: handle ctx
	ch := make(chan model.PlcReadRequestResult)
	go func() {
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
		responseCodes := make(map[string]model.PlcResponseCode)
		responseValues := make(map[string]values.PlcValue)
		for _, tagName := range readRequest.GetTagNames() {
			tag := readRequest.GetTag(tagName)
			simulatedTagVar, ok := tag.(simulatedTag)
			if !ok {
				responseCodes[tagName] = model.PlcResponseCode_INVALID_ADDRESS
				responseValues[tagName] = nil
			} else {
				value := r.device.Get(simulatedTagVar)
				if value == nil {
					responseCodes[tagName] = model.PlcResponseCode_NOT_FOUND
					responseValues[tagName] = nil
				} else {
					responseCodes[tagName] = model.PlcResponseCode_OK
					responseValues[tagName] = *value
				}
			}
		}

		if r.tracer != nil {
			r.tracer.AddTransactionalTrace(txId, "read", "success")
		}
		// Emit the response
		ch <- &model2.DefaultPlcReadRequestResult{
			Request:  readRequest,
			Response: model2.NewDefaultPlcReadResponse(readRequest, responseCodes, responseValues),
			Err:      nil,
		}
	}()
	return ch
}
