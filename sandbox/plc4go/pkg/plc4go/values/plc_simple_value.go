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

type plcSimpleValueAdapter struct {
	plcValueAdapter
}

func (m plcSimpleValueAdapter) IsSimple() bool {
	return true
}

func (m plcSimpleValueAdapter) GetLength() int {
	return 1
}

type plcSimpleNumericValueAdapter struct {
	plcSimpleValueAdapter
}

func (m plcSimpleNumericValueAdapter) IsBoolean() bool {
	return true
}

func (m plcSimpleNumericValueAdapter) IsUint8() bool {
	return true
}

func (m plcSimpleNumericValueAdapter) IsUint16() bool {
	return true
}

func (m plcSimpleNumericValueAdapter) IsUint32() bool {
	return true
}

func (m plcSimpleNumericValueAdapter) IsUint64() bool {
	return true
}

func (m plcSimpleNumericValueAdapter) IsUnt8() bool {
	return true
}

func (m plcSimpleNumericValueAdapter) IsUnt16() bool {
	return true
}

func (m plcSimpleNumericValueAdapter) IsUnt32() bool {
	return true
}

func (m plcSimpleNumericValueAdapter) IsUnt64() bool {
	return true
}

func (m plcSimpleNumericValueAdapter) IsFloat32() bool {
	return true
}

func (m plcSimpleNumericValueAdapter) IsFloat64() bool {
	return true
}
