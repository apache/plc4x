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

package bacgopes

import (
	"encoding/hex"
	"regexp"
	"strings"
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

//go:generate plc4xGenerator -type=DebugContents -prefix=debugging_
type DebugContents struct {
	// TODO: implement me
}
