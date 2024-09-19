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
	"strings"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
)

type Arg any

type Args []any

var NoArgs = NewArgs()

func NewArgs(args ...any) Args {
	return args
}

// NA is a shortcut for NewArgs
var NA = NewArgs

// GetFromArgs gets a value fromArgs and if not present panics
func GetFromArgs[T any](args Args, index int) T {
	if index > len(args)-1 {
		panic(fmt.Sprintf("index out of bounds: %d(len %d of %s)", index, len(args), args))
	}
	aAtI := args[index]
	v, ok := aAtI.(T)
	if !ok {
		panic(fmt.Sprintf("argument #%d with type %T is not of type %T", index, aAtI, *new(T)))
	}
	return v
}

// GA is a shortcut for GetFromArgs
func GA[T any](args Args, index int) T {
	return GetFromArgs[T](args, index)
}

// GetFromArgsOptional gets a value from Args or return default if not present
func GetFromArgsOptional[T any](args Args, index int, defaultValue T) (T, bool) {
	if index > len(args)-1 {
		return defaultValue, false
	}
	return args[index].(T), true
}

// GAO is a shortcut for GetFromArgsOptional
func GAO[T any](args Args, index int, defaultValue T) (T, bool) {
	return GetFromArgsOptional(args, index, defaultValue)
}

func (a Args) Format(s fmt.State, verb rune) {
	switch verb {
	case 'r':
		_, _ = fmt.Fprintf(s, "(%s)", a.string(false, false)[1:len(a.string(false, false))-1])
	case 's', 'v':
		_, _ = fmt.Fprintf(s, "(%s)", a.String()[1:len(a.String())-1])
	}
}

func (a Args) String() string {
	return a.string(true, true)
}
func (a Args) string(printIndex bool, printType bool) string {
	r := ""
	for i, ea := range a {
		eat := fmt.Sprintf("%T", ea)
		switch tea := ea.(type) {
		case []byte:
			ea = Btox(tea, ".")
		case string:
			ea = "'" + tea + "'"
		case fmt.Stringer:
			if !IsNil(tea) {
				teaString := tea.String()
				ea = teaString
				if strings.Contains(teaString, "\n") {
					ea = "\n" + teaString + "\n"
				}
			}
		}
		if printIndex {
			r += fmt.Sprintf("%d: ", i)
		}
		r += fmt.Sprintf("%v", ea)
		if printType {
			r += fmt.Sprintf(" (%s)", eat)
		}
		r += ", "
	}
	if r != "" {
		r = r[:len(r)-2]
	}
	return "[" + r + "]"
}
