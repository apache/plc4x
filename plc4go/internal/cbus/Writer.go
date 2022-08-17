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

package cbus

import (
	"context"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/pkg/errors"
)

type Writer struct {
	alphaGenerator *AlphaGenerator
	messageCodec   spi.MessageCodec
	tm             *spi.RequestTransactionManager
}

func NewWriter(tpduGenerator *AlphaGenerator, messageCodec spi.MessageCodec, tm *spi.RequestTransactionManager) Writer {
	return Writer{
		alphaGenerator: tpduGenerator,
		messageCodec:   messageCodec,
		tm:             tm,
	}
}

func (m Writer) Write(ctx context.Context, writeRequest model.PlcWriteRequest) <-chan model.PlcWriteRequestResult {
	// TODO: handle context
	result := make(chan model.PlcWriteRequestResult)
	go func() {
		result <- &plc4goModel.DefaultPlcWriteRequestResult{
			Request:  writeRequest,
			Response: nil,
			Err:      errors.New("Not yet implemented"),
		}
	}()
	return result
}
