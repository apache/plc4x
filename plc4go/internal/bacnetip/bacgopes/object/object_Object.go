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

import (
	"io"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/basetypes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

// TODO: big WIP
type Object interface {
	DebugContentPrinter
	GetAttr(name string) (any, bool)
	SetAttr(name string, value any)
	AddProperty(prop Property)
	DeleteProperty(prop string)
	ReadProperty() error
	WriteProperty() error
}

type _Object struct {
	// TODO: debug contents
	_objectSupportsCov bool

	properties []Property

	_leafName string
}

func NewObject(kwArgs KWArgs, options ...Option) (Object, error) {
	o := &_Object{
		properties: []Property{
			NewObjectIdentifierProperty("objectIdentifier", Vs2P(NewObjectIdentifier), WithPropertyOptional(false)),
			NewReadableProperty("objectName", V2P(NewCharacterString), WithPropertyOptional(false)),
			NewOptionalProperty("description", V2P(NewCharacterString)),
			NewOptionalProperty("profileName", V2P(NewCharacterString)),
			NewReadableProperty("propertyList", ArrayOfP(NewPropertyIdentifier, 0, 0)),
			NewOptionalProperty("auditLevel", V2P(NewAuditLevel)),
			NewOptionalProperty("auditableOperations", V2P(NewAuditOperationFlags)),
			NewOptionalProperty("tags", ArrayOfsP(NewNameValue, 0, 0)),
			NewOptionalProperty("profileLocation", V2P(NewCharacterString)),
			NewOptionalProperty("profileName", V2P(NewCharacterString)),
		},
		_leafName: ExtractLeafName(options, StructName()),
	}
	if _debug != nil {
		_debug("__init__(%s) %r", o._leafName, kwArgs)
	}
	panic("implement me")
	return o, nil
}

func (o *_Object) PrintDebugContents(indent int, file io.Writer, _ids []uintptr) {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) GetAttr(name string) (any, bool) {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) SetAttr(name string, value any) {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) AddProperty(prop Property) {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) DeleteProperty(prop string) {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) ReadProperty() error {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) WriteProperty() error {
	//TODO implement me
	panic("implement me")
}
