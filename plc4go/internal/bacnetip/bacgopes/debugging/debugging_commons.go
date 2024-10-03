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

package debugging

import (
	"encoding/hex"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func Btox(data []byte, sep string) string {
	hexString := hex.EncodeToString(data)
	if sep != "" {
		pairs := make([]string, len(hexString)/2)
		for i := 0; i < len(hexString)-1; i += 2 {
			pairs[i/2] = hexString[i : i+2]
		}
		hexString = strings.Join(pairs, ".")
	}
	return hexString
}

func Xtob(hexString string) ([]byte, error) {
	compile, err := regexp.Compile("[^0-9a-fA-F]")
	if err != nil {
		return nil, err
	}
	replaceAll := compile.ReplaceAll([]byte(hexString), nil)
	decodeString, err := hex.DecodeString(string(replaceAll))
	if err != nil {
		return nil, err
	}
	return decodeString, nil
}

func VerbForType(value any, printVerb rune) rune {
	if isNil(value) {
		return 'v'
	}
	switch value.(type) {
	case string:
		printVerb = 's'
	case bool:
		printVerb = 't'
	case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, uintptr:
		printVerb = 'd'
	case *int8, *uint8, *int16, *uint16, *int32, *uint32, *int64, *uint64, *int, *uint, *uintptr:
		printVerb = 'v'
	case float32, float64:
		printVerb = 'f'
	case complex64, complex128:
		printVerb = 'v' // TODO: what is it meant for?
	case time.Time, time.Duration:
		printVerb = 's'
	case []int8, []uint8, []int16, []uint16, []int32, []uint32, []int64, []uint64, []int, []uint, []uintptr:
		printVerb = 'v'
	case interface{ PLC4XEnumName() string }: // shortcut to handle model enums
		printVerb = 'd'
	}
	return printVerb
}

// clone from comp to avoid circular dependencies // TODO: maybe move Btox somewhere else or come up with something smarter there
func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	valueOf := reflect.ValueOf(v)
	switch valueOf.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return valueOf.IsNil()
	default:
		return false
	}
}

func StructName() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	dir := path.Dir(file)
	rootIndex := strings.Index(dir, projectName)
	dir = dir[rootIndex:]
	dirPrefix := path.Base(dir) + "_"
	base := path.Base(file)
	prefix := strings.TrimSuffix(base, ".go")
	return strings.TrimPrefix(prefix, dirPrefix)
}
