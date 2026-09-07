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

package firmata

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Writer writes digital pins. Ported from plc4j's FirmataConnection.buildWriteMessages: a pin is
// switched to output mode the first time it is written to, and the value is then shipped as a
// set-digital-pin-value command.
//
// Firmata acknowledges nothing, so a write is done once its bytes are on the wire - there is no
// response to wait for and no way for the board to report a problem.
type Writer struct {
	connection *Connection

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewWriter(connection *Connection, _options ...options.WithOption) *Writer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Writer{
		connection: connection,
		log:        customLogger,
	}
}

func (w *Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	w.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(w.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		if err := ctx.Err(); err != nil {
			utils.DeliverResult(w.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, err))
			return
		}

		// Every tag is written on its own and answered on its own. plc4j fails the whole request
		// as soon as one tag doesn't work out, which isn't any more atomic - firmata has no
		// transaction - and loses the tags that would have gone through.
		responseCodes := map[string]apiModel.PlcResponseCode{}
		for _, tagName := range writeRequest.GetTagNames() {
			responseCodes[tagName] = w.writeTag(ctx, writeRequest, tagName)
		}
		utils.DeliverResult(w.log, result, spiModel.NewDefaultPlcWriteRequestResult(
			writeRequest,
			spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes),
			nil,
		))
	})
	return result
}

// writeTag ships one tag and reports how it went.
func (w *Writer) writeTag(ctx context.Context, writeRequest apiModel.PlcWriteRequest, tagName string) apiModel.PlcResponseCode {
	digitalTagVar, ok := writeRequest.GetTag(tagName).(digitalTag)
	if !ok {
		// Analog pins are read-only in this driver, the same as in plc4j: writing one needs the
		// extended-analog sysex command, which neither driver speaks.
		w.log.Debug().
			Str("tagName", tagName).
			Interface("tag", writeRequest.GetTag(tagName)).
			Msg("Writing is only supported for digital pins")
		return apiModel.PlcResponseCode_UNSUPPORTED
	}

	values, err := boolValues(writeRequest.GetValue(tagName), digitalTagVar.quantity)
	if err != nil {
		w.log.Debug().Err(err).Str("tagName", tagName).Msg("Invalid value")
		return apiModel.PlcResponseCode_INVALID_DATA
	}

	messages, claim, err := w.connection.claimOutputPins(digitalTagVar.address, digitalTagVar.quantity)
	if err != nil {
		w.log.Debug().Err(err).Str("tagName", tagName).Msg("Unable to use the pins of this tag as outputs")
		return apiModel.PlcResponseCode_INVALID_ADDRESS
	}
	for i, value := range values {
		pin := digitalTagVar.address + uint8(i)
		messages = append(messages, readWriteModel.NewFirmataMessageCommand(
			readWriteModel.NewFirmataCommandSetDigitalPinValue(pin, value)))
	}

	if err := w.connection.sendAll(ctx, "write", messages); err != nil {
		// The set-pin-mode messages are what configures the pins, so a claim whose messages didn't
		// make it out has to be given back - otherwise a retry would consider the pins configured
		// and never send them.
		claim.rollback()
		w.log.Debug().Err(err).Str("tagName", tagName).Msg("Error sending the write messages")
		return apiModel.PlcResponseCode_INTERNAL_ERROR
	}
	return apiModel.PlcResponseCode_OK
}

// boolValues unpacks the value of a write request into one boolean per pin the tag covers.
func boolValues(plcValue apiValues.PlcValue, quantity uint8) ([]bool, error) {
	if plcValue == nil {
		return nil, errors.New("no value to write")
	}
	if plcValue.IsList() {
		elements := plcValue.GetList()
		if len(elements) != int(quantity) {
			return nil, errors.Errorf("expected %d values but got %d", quantity, len(elements))
		}
		values := make([]bool, 0, len(elements))
		for _, element := range elements {
			if element == nil || !element.IsBool() {
				return nil, errors.New("expecting only BOOL values when writing digital pins")
			}
			values = append(values, element.GetBool())
		}
		return values, nil
	}
	if quantity != 1 {
		return nil, errors.Errorf("expected %d values but got a single one", quantity)
	}
	if !plcValue.IsBool() {
		return nil, errors.New("expecting a BOOL value when writing a digital pin")
	}
	return []bool{plcValue.GetBool()}, nil
}
