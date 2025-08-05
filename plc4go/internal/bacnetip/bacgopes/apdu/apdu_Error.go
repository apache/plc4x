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

package apdu

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/basetypes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructeddata"
)

type Error struct {
	*ErrorSequence
}

func NewError(args Args, kwArgs KWArgs) (*Error, error) {
	e := &Error{}
	var err error
	e.ErrorSequence, err = NewErrorSequence(args, kwArgs, WithErrorSequenceExtension(e))
	if err != nil {
		return e, errors.Wrap(err, "Error creating new ErrorSequence")
	}
	return e, nil
}

func (e *Error) SetErrorSequence(es *ErrorSequence) {
	e.ErrorSequence = es
}

func (e *Error) GetSequenceElements() []Element {
	errorType, _ := basetypes.NewErrorType(nil) // TODO: check if is meant to be like that
	return errorType.GetSequenceElements()
}

func (e *Error) String() string {
	errorClass, _ := e.GetAttr("errorClass")
	errorCode, _ := e.GetAttr("errorCode")
	return fmt.Sprintf("%s:%s", errorClass, errorCode)
}
