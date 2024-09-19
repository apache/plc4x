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

package object

import . "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"

// TODO: big WIP
type Property interface {
	ReadProperty()
	WriteProperty()
}

type PropertyKlass interface {
	Encode(Arg) error
}

func NewProperty(name string, klass func(Args, KWArgs) (PropertyKlass, error), opts ...func(*PropertyKlass)) Property {
	return &_Property{}
}

type _Property struct {
	Name     string
	Klass    func(Args, KWArgs) (PropertyKlass, error)
	Optional bool
}

var _ Property = (*_Property)(nil)

func WithPropertyOptional(optional bool) func(*_Property) {
	return func(e *_Property) {
		e.Optional = optional
	}
}

func (p *_Property) ReadProperty() {
	//TODO implement me
	panic("implement me")
}

func (p *_Property) WriteProperty() {
	//TODO implement me
	panic("implement me")
}
