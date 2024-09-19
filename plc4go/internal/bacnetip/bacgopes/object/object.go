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

// TODO: big WIP
package object

import (
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
)

var _debug = CreateDebugPrinter()

var registeredObjectTypes map[any]struct{}

func init() {
	registeredObjectTypes = make(map[any]struct{})
}

// V2P accepts a function which takes an Arg and maps it to a PropertyKlass
func V2P[T any](b func(arg Arg) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	return func(args Args, kwArgs KWArgs) (PropertyKlass, error) {
		var arg any
		if len(args) == 1 {
			arg = args[0]
		}
		r, err := b(arg)
		return any(r).(PropertyKlass), err
	}
}

// Vs2P accepts a function which takes an Args and maps it to a PropertyKlass
func Vs2P[T any](b func(args Args) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	return func(args Args, kwArgs KWArgs) (PropertyKlass, error) {
		r, err := b(args)
		return any(r).(PropertyKlass), err
	}
}

// TODO: finish
func SequenceOfP[T any](b func(arg Arg) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	panic("finish me")
}

func SequenceOfsP[T any](b func(args Args) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	panic("finish me")
}

// TODO: finish // convert to kwArgs and check wtf we are doing here...
func ArrayOfP[T any](b func(arg Arg) (*T, error), fixedLength int, prototype any) func(Args, KWArgs) (PropertyKlass, error) {
	panic("finish me")
}

// TODO: finish // convert to kwArgs and check wtf we are doing here...
func ArrayOfsP[T any](b func(args Args) (*T, error), fixedLength int, prototype any) func(Args, KWArgs) (PropertyKlass, error) {
	panic("finish me")
}

// TODO: finish // convert to kwArgs and check wtf we are doing here...
func ListOfP[T any](b func(arg Arg) (*T, error)) func(Args, KWArgs) (PropertyKlass, error) {
	panic("finish me")
}
