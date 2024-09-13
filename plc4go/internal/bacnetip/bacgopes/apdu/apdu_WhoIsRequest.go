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
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructeddata"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type WhoIsRequest struct {
	*UnconfirmedRequestSequence

	serviceChoice    readWriteModel.BACnetUnconfirmedServiceChoice
	sequenceElements []Element
}

func NewWhoIsRequest() (*WhoIsRequest, error) {
	w := &WhoIsRequest{
		serviceChoice: readWriteModel.BACnetUnconfirmedServiceChoice_WHO_IS,
		sequenceElements: []Element{
			NewElement("deviceInstanceRangeLowLimit", V2E(NewUnsigned), WithElementContext(0), WithElementOptional(true)),
			NewElement("deviceInstanceRangeHighLimit", V2E(NewUnsigned), WithElementContext(1), WithElementOptional(true)),
		},
	}
	var err error
	w.UnconfirmedRequestSequence, err = NewUnconfirmedRequestSequence(
		readWriteModel.NewBACnetUnconfirmedServiceRequestWhoIs(
			readWriteModel.CreateBACnetContextTagUnsignedInteger(0, 0),
			readWriteModel.CreateBACnetContextTagUnsignedInteger(1, 0),
			0,
		),
		NoKWArgs(),
		WithUnconfirmedRequestSequenceExtension(w),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating UnconfirmedRequestSequence")
	}

	return w, nil
}

func (w *WhoIsRequest) GetServiceChoice() *readWriteModel.BACnetUnconfirmedServiceChoice {
	return &w.serviceChoice
}

func (w *WhoIsRequest) GetSequenceElements() []Element {
	return w.sequenceElements
}

func (w *WhoIsRequest) SetUnconfirmedRequestSequence(u *UnconfirmedRequestSequence) {
	w.UnconfirmedRequestSequence = u
}
