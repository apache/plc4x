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

package comp

import (
	"fmt"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
)

type KnownKey string

type KWArgs map[KnownKey]any

var NoKWArgs = NewKWArgs

func NewKWArgs(kw ...any) KWArgs {
	if len(kw)%2 != 0 {
		panic("KWArgs must have an even number of arguments")
	}
	r := make(KWArgs)
	for i := 0; i < len(kw)-1; i += 2 {
		key, ok := kw[i].(KnownKey)
		if !ok {
			panic("keys must be of type KnownKey")
		}
		r[key] = kw[i+1]
	}
	return r
}

// NKW is a shortcut for NewKWArgs
var NKW = NewKWArgs

func (k KWArgs) Format(s fmt.State, verb rune) {
	switch verb {
	case 'r':
		_, _ = fmt.Fprint(s, "{")
		currentE := 0
		nElements := len(k)
		for kk, element := range k {
			currentE++
			_, _ = fmt.Fprintf(s, "'%s': ", kk)
			switch element := element.(type) {
			case interface{ StructHeader() []byte }:
				_, _ = s.Write(element.StructHeader())
			case fmt.Formatter:
				element.Format(s, verb)
			case fmt.Stringer:
				_, _ = fmt.Fprint(s, element.String())
			default:
				_, _ = fmt.Fprintf(s, "%v", element)
			}
			if currentE < nElements {
				_, _ = fmt.Fprint(s, ", ")
			}
		}
		_, _ = fmt.Fprint(s, "}")
	case 's', 'v':
		_, _ = fmt.Fprint(s, k.String())
	}
}

func (k KWArgs) String() string {
	r := ""
	for kk, ea := range k {
		switch tea := ea.(type) {
		case []byte:
			ea = Btox(tea, ".")
		}
		if IsNil(ea) {
			ea = fmt.Sprintf("<nil>(%T)", ea)
		}
		r += fmt.Sprintf("'%s': %v, ", kk, ea)
	}
	if r != "" {
		r = r[:len(r)-2]
	}
	return "{" + r + "}"
}

// KW gets a value from KWArgs and if not present panics
func KW[T any](kwArgs KWArgs, key KnownKey) T {
	r, ok := kwArgs[key]
	if !ok {
		panic(fmt.Sprintf("key %v not found in kwArgs", key))
	}
	delete(kwArgs, key) // usually that means this argument was consumed so we get rid of it
	return r.(T)
}

// KWO gets a value from KWArgs and if not present returns the supplied default value
func KWO[T any](kwArgs KWArgs, key KnownKey, defaultValue T) (T, bool) {
	r, ok := kwArgs[key]
	if !ok {
		return defaultValue, false
	}
	v, ok := r.(T)
	if !ok {
		return defaultValue, false
	}
	delete(kwArgs, key) // usually that means this argument was consumed so we get rid of it
	return v, true
}
