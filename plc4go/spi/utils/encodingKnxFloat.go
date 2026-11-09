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

package utils

import (
	"math"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// knxFloatEncoding is the floatEncoding of the 2 byte float of KNX (DPT 9): sign in bit 15,
// exponent in bits 14-11, and a 12 bit two's complement mantissa made of the sign and bits 10-0;
// the value is 0.01 * mantissa * 2^exponent.
// (Java: org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingKnxFloat)
const knxFloatEncoding = "KNXFloat"

// knxFloatMax is the largest magnitude a KNX float encodes: mantissa 2047, exponent 15.
const knxFloatMax = 2047.0 * (1 << 15) / 100.0

func isKnxFloat(encoding string) bool {
	return encoding == knxFloatEncoding
}

func decodeKnxFloat(raw uint16) float32 {
	exponent := (raw >> 11) & 0x0F
	mantissa := int16(raw & 0x07FF)
	if raw&0x8000 != 0 {
		mantissa = int16(raw&0x07FF | 0xF800)
	}
	return float32(0.01 * float64(mantissa) * math.Pow(2, float64(exponent)))
}

// encodeKnxFloat encodes as the java implementation does, truncating the mantissa to 11 bits.
func encodeKnxFloat(value float32) (uint16, error) {
	if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		return 0, errors.Errorf("a KNX float can't encode %v", value)
	}
	if math.Abs(float64(value)) > knxFloatMax {
		return 0, errors.Errorf("%v is outside the range of a KNX float (±%v)", value, knxFloatMax)
	}
	hundredths := int64(math.Round(float64(value) * 100))
	negative := hundredths < 0
	mantissa := hundredths
	if negative {
		mantissa = -mantissa
	}
	exponent := uint16(0)
	for mantissa > 0x07FF {
		mantissa >>= 1
		exponent++
	}
	if negative {
		mantissa = (^mantissa + 1) & 0x07FF
	}
	raw := exponent<<11 | uint16(mantissa)&0x07FF
	if negative {
		raw |= 0x8000
	}
	return raw, nil
}
