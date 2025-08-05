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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructeddata"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type ConfirmedPrivateTransferRequest struct {
	*ConfirmedRequestSequence

	serviceChoice    readWriteModel.BACnetConfirmedServiceChoice
	sequenceElements []Element
}

func NewConfirmedPrivateTransferRequest(args Args, kwArgs KWArgs, options ...Option) (*ConfirmedPrivateTransferRequest, error) {
	c := &ConfirmedPrivateTransferRequest{
		serviceChoice: readWriteModel.BACnetConfirmedServiceChoice_CONFIRMED_PRIVATE_TRANSFER,
		sequenceElements: []Element{
			NewElement("vendorID", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("serviceNumber", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("serviceParameters", Vs2E(NewAny), WithElementContext(2), WithElementOptional(true)),
		},
	}
	options = AddRootMessageIfAbundant(options, readWriteModel.NewBACnetConfirmedServiceRequestConfirmedPrivateTransfer(
		readWriteModel.CreateBACnetVendorIdContextTagged(0, 0),     // TODO: get right values
		readWriteModel.CreateBACnetContextTagUnsignedInteger(1, 0), // TODO: get right values
		nil,
		0,
	))
	options = AddLeafTypeIfAbundant(options, c)
	var err error
	c.ConfirmedRequestSequence, err = NewConfirmedRequestSequence(args, kwArgs, Combine(options, WithConfirmedRequestSequenceExtension(c))...)
	if err != nil {
		return nil, errors.Wrap(err, "error building confirmed request")
	}
	return c, nil
}

func (c *ConfirmedPrivateTransferRequest) SetConfirmedRequestSequence(crs *ConfirmedRequestSequence) {
	c.ConfirmedRequestSequence = crs
}

func (c *ConfirmedPrivateTransferRequest) String() string {
	return fmt.Sprintf("ConfirmedPrivateTransferRequest{%s}", c.ConfirmedRequestSequence)
}
