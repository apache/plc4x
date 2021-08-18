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

package model

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
)

func S7EventHelperRightShift3(readBuffer utils.ReadBuffer) (uint16, error) {
	return 0, nil
}

func S7EventHelperLeftShift3(writeBuffer utils.WriteBuffer, valueLength uint16) error {
	return nil
}

func S7EventHelperEventItemLength(readBuffer utils.ReadBuffer, valueLength uint16) uint16 {
	return 0
}

func S7EventHelperBcdToInt(readBuffer utils.ReadBuffer) (uint8, error) {
	return 0, nil
}

func S7EventHelperByteToBcd(writeBuffer utils.WriteBuffer, value uint8) error {
	return nil
}

func S7EventHelperS7msecToInt(readBuffer utils.ReadBuffer) (uint16, error) {
	return 0, nil
}

func S7EventHelperIntToS7msec(writeBuffer utils.WriteBuffer, value uint16) error {
	return nil
}
