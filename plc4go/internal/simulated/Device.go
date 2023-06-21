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
	"github.com/rs/zerolog"
	"math/rand"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
)

type Device struct {
	Name  string
	State map[simulatedTag]*apiValues.PlcValue

	passLogToModel bool
	log            zerolog.Logger
}

func NewDevice(name string, _options ...options.WithOption) *Device {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger, _ := options.ExtractCustomLogger(_options...)
	return &Device{
		Name:           name,
		State:          make(map[simulatedTag]*apiValues.PlcValue),
		passLogToModel: passLoggerToModel,
		log:            customLogger,
	}
}

func (d *Device) Get(tag simulatedTag) *apiValues.PlcValue {
	switch tag.TagType {
	case TagState:
		return d.State[tag]
	case TagRandom:
		return d.getRandomValue(tag)
	}
	return nil
}

func (d *Device) Set(tag simulatedTag, value *apiValues.PlcValue) {
	switch tag.TagType {
	case TagState:
		d.State[tag] = value
		break
	case TagRandom:
		// TODO: Doesn'd really make any sense to write a random
		break
	case TagStdOut:
		d.log.Debug().Msgf("TEST PLC STDOUT [%s]: %s", tag.Name, (*value).GetString())
		break
	}
}

func (d *Device) getRandomValue(tag simulatedTag) *apiValues.PlcValue {
	size := tag.GetDataTypeSize().DataTypeSize()
	data := make([]byte, uint16(size)*tag.Quantity)
	rand.Read(data)
	ctxForModel := options.GetLoggerContextForModel(context.TODO(), d.log, options.WithPassLoggerToModel(d.passLogToModel))
	plcValue, err := readWriteModel.DataItemParse(ctxForModel, data, tag.DataTypeSize.String(), tag.Quantity)
	if err != nil {
		d.log.Err(err).Msg("Unable to parse random bytes")
		return nil
	}
	return &plcValue
}
