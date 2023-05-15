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
	"github.com/pkg/errors"
	"strconv"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

type Writer struct {
	device  *Device
	options map[string][]string
	tracer  *spi.Tracer
}

func NewWriter(device *Device, options map[string][]string, tracer *spi.Tracer) *Writer {
	return &Writer{
		device:  device,
		options: options,
		tracer:  tracer,
	}
}

func (w *Writer) Write(_ context.Context, writeRequest model.PlcWriteRequest) <-chan model.PlcWriteRequestResult {
	ch := make(chan model.PlcWriteRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Errorf("panic-ed %v", err))
			}
		}()
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
		for _, tagName := range writeRequest.GetTagNames() {
			tag := writeRequest.GetTag(tagName)
			simulatedTagVar, ok := tag.(simulatedTag)
			if !ok {
				responseCodes[tagName] = model.PlcResponseCode_INVALID_ADDRESS
			} else {
				plcValue := writeRequest.GetValue(tagName)
				w.device.Set(simulatedTagVar, &plcValue)
				responseCodes[tagName] = model.PlcResponseCode_OK
			}
		}

		if w.tracer != nil {
			w.tracer.AddTransactionalTrace(txId, "write", "success")
		}
		// Emit the response
		ch <- &spiModel.DefaultPlcWriteRequestResult{
			Request:  writeRequest,
			Response: spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes),
			Err:      nil,
		}
	}()
	return ch
}
