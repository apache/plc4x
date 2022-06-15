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
	"github.com/apache/plc4x/plc4go/internal/spi"
	model2 "github.com/apache/plc4x/plc4go/internal/spi/model"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"strconv"
	"time"
)

type Writer struct {
	device  *Device
	options map[string][]string
	tracer  *spi.Tracer
}

func NewWriter(device *Device, options map[string][]string, tracer *spi.Tracer) Writer {
	return Writer{
		device:  device,
		options: options,
		tracer:  tracer,
	}
}

func (w Writer) Write(writeRequest model.PlcWriteRequest) <-chan model.PlcWriteRequestResult {
	ch := make(chan model.PlcWriteRequestResult)
	go func() {
		var txId string
		if w.tracer != nil {
			txId = w.tracer.AddTransactionalStartTrace("write", "started")
		}
		// Possibly add a delay.
		if delayString, ok := w.options["writeDelay"]; ok {
			if len(delayString) == 1 {
				delay, err := strconv.Atoi(delayString[0])
				if err == nil {
					time.Sleep(time.Duration(delay) * time.Millisecond)
				}
			}
		}

		// Process the request
		responseCodes := map[string]model.PlcResponseCode{}
		for _, fieldName := range writeRequest.GetFieldNames() {
			field := writeRequest.GetField(fieldName)
			simulatedField, ok := field.(SimulatedField)
			if !ok {
				responseCodes[fieldName] = model.PlcResponseCode_INVALID_ADDRESS
			} else {
				plcValue := writeRequest.GetValue(fieldName)
				w.device.Set(simulatedField, &plcValue)
				responseCodes[fieldName] = model.PlcResponseCode_OK
			}
		}

		if w.tracer != nil {
			w.tracer.AddTransactionalTrace(txId, "write", "success")
		}
		// Emit the response
		ch <- &model2.DefaultPlcWriteRequestResult{
			Request:  writeRequest,
			Response: model2.NewDefaultPlcWriteResponse(writeRequest, responseCodes),
			Err:      nil,
		}
	}()
	return ch
}
