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

package fields

import (
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// convBool converts an arbitrary value to a bool. Booleans pass through, numerics are true when
// non-zero, strings parse via strconv.ParseBool.
func convBool(v any) (bool, error) {
	if v == nil {
		return false, nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		return rv.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0, nil
	case reflect.Complex64, reflect.Complex128:
		return rv.Complex() != 0, nil
	case reflect.String:
		return strconv.ParseBool(rv.String())
	}
	return false, errors.Errorf("cannot convert %T to bool", v)
}

// convInt64 converts an arbitrary value to int64. Booleans yield 0/1, floats truncate, complex
// values drop the imaginary part, strings parse via strconv.ParseInt.
func convInt64(v any) (int64, error) {
	if v == nil {
		return 0, nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(rv.Float()), nil
	case reflect.Complex64, reflect.Complex128:
		return int64(real(rv.Complex())), nil
	case reflect.String:
		return strconv.ParseInt(rv.String(), 10, 64)
	}
	return 0, errors.Errorf("cannot convert %T to int64", v)
}

// convUint64 converts an arbitrary value to uint64. Negative numbers wrap via two's complement.
func convUint64(v any) (uint64, error) {
	if v == nil {
		return 0, nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(rv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return uint64(rv.Float()), nil
	case reflect.Complex64, reflect.Complex128:
		return uint64(real(rv.Complex())), nil
	case reflect.String:
		return strconv.ParseUint(rv.String(), 10, 64)
	}
	return 0, errors.Errorf("cannot convert %T to uint64", v)
}

// convFloat64 converts an arbitrary value to float64.
func convFloat64(v any) (float64, error) {
	if v == nil {
		return 0, nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return rv.Float(), nil
	case reflect.Complex64, reflect.Complex128:
		return real(rv.Complex()), nil
	case reflect.String:
		return strconv.ParseFloat(rv.String(), 64)
	}
	return 0, errors.Errorf("cannot convert %T to float64", v)
}

// convComplex128 converts an arbitrary value to complex128. Numerics become the real part with
// imaginary 0; strings parse via strconv.ParseComplex.
func convComplex128(v any) (complex128, error) {
	if v == nil {
		return 0, nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return complex(float64(rv.Int()), 0), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return complex(float64(rv.Uint()), 0), nil
	case reflect.Float32, reflect.Float64:
		return complex(rv.Float(), 0), nil
	case reflect.Complex64, reflect.Complex128:
		return rv.Complex(), nil
	case reflect.String:
		return strconv.ParseComplex(rv.String(), 128)
	}
	return 0, errors.Errorf("cannot convert %T to complex128", v)
}

// Typed narrow wrappers so the call site can funnel directly: f.toT(convInt8(v)).

func convInt(v any) (int, error)             { n, err := convInt64(v); return int(n), err }
func convInt8(v any) (int8, error)           { n, err := convInt64(v); return int8(n), err }
func convInt16(v any) (int16, error)         { n, err := convInt64(v); return int16(n), err }
func convInt32(v any) (int32, error)         { n, err := convInt64(v); return int32(n), err }
func convUint(v any) (uint, error)           { n, err := convUint64(v); return uint(n), err }
func convUint8(v any) (uint8, error)         { n, err := convUint64(v); return uint8(n), err }
func convUint16(v any) (uint16, error)       { n, err := convUint64(v); return uint16(n), err }
func convUint32(v any) (uint32, error)       { n, err := convUint64(v); return uint32(n), err }
func convUintptr(v any) (uintptr, error)     { n, err := convUint64(v); return uintptr(n), err }
func convFloat32(v any) (float32, error)     { n, err := convFloat64(v); return float32(n), err }
func convComplex64(v any) (complex64, error) { c, err := convComplex128(v); return complex64(c), err }
