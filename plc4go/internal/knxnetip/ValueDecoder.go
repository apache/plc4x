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

package knxnetip

import (
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	api "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
)

type ValueDecoder struct {
	rb utils.ReadBuffer
}

func NewValueDecoder(rb utils.ReadBuffer) ValueDecoder {
	return ValueDecoder{
		rb: rb,
	}
}

func (m ValueDecoder) Decode(typeName string) api.PlcValue {
	datatype := driverModel.KnxDatapointTypeByName(typeName)
	plcValue, err := driverModel.KnxDatapointParse(m.rb, datatype)
	if err != nil {
		return nil
	}
	return plcValue
}
