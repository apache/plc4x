/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	model2 "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"strconv"
	"time"
)

type Reader struct {
	device  *Device
	options map[string][]string
}

func NewReader(device *Device, options map[string][]string) Reader {
	return Reader{
		device:  device,
		options: options,
	}
}

func (r Reader) Read(readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	ch := make(chan model.PlcReadRequestResult)
	go func() {
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
		for _, fieldName := range readRequest.GetFieldNames() {
			field := readRequest.GetField(fieldName)
			simulatedField, ok := field.(SimulatedField)
			if !ok {
				responseCodes[fieldName] = model.PlcResponseCode_INVALID_ADDRESS
				responseValues[fieldName] = nil
			} else {
				value := r.device.Get(simulatedField)
				if value == nil {
					responseCodes[fieldName] = model.PlcResponseCode_NOT_FOUND
					responseValues[fieldName] = nil
				} else {
					responseCodes[fieldName] = model.PlcResponseCode_OK
					responseValues[fieldName] = *value
				}
			}
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
