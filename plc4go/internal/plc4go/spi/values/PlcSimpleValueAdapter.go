//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package values

type PlcSimpleValueAdapter struct {
	PlcValueAdapter
}

func (m PlcSimpleValueAdapter) IsSimple() bool {
	return true
}

func (m PlcSimpleValueAdapter) GetLength() uint32 {
	return 1
}

type PlcSimpleNumericValueAdapter struct {
	PlcSimpleValueAdapter
}

func (m PlcSimpleNumericValueAdapter) IsBool() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUint8() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUint16() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUint32() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUint64() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUnt8() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUnt16() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUnt32() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUnt64() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsFloat32() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsFloat64() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsString() bool {
	return true
}
