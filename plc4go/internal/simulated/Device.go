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
	"github.com/rs/zerolog"
	"math/rand"

	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
)

type Device struct {
	Name  string
	State map[simulatedTag]*values.PlcValue

	log zerolog.Logger
}

func NewDevice(name string, _options ...options.WithOption) *Device {
	return &Device{
		Name:  name,
		State: make(map[simulatedTag]*values.PlcValue),
		log:   options.ExtractCustomLogger(_options...),
	}
}

func (d *Device) Get(tag simulatedTag) *values.PlcValue {
	switch tag.TagType {
	case TagState:
		return d.State[tag]
	case TagRandom:
		return d.getRandomValue(tag)
	}
	return nil
}

func (d *Device) Set(tag simulatedTag, value *values.PlcValue) {
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

func (d *Device) getRandomValue(tag simulatedTag) *values.PlcValue {
	size := tag.GetDataTypeSize().DataTypeSize()
	data := make([]byte, uint16(size)*tag.Quantity)
	rand.Read(data)
	plcValue, err := model.DataItemParse(context.Background(), data, tag.DataTypeSize.String(), tag.Quantity)
	if err != nil {
		panic("Unable to parse random bytes")
	}
	return &plcValue
}
