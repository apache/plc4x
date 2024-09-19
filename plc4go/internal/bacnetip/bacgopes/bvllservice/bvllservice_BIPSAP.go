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

package bvllservice

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type BIPSAPRequirements interface {
	ServiceAccessPoint
	Client
}

type BIPSAP struct {
	ServiceAccessPointContract
	*DefaultRFormatter `ignore:"true"`

	requirements BIPSAPRequirements

	log zerolog.Logger
}

func NewBIPSAP(localLog zerolog.Logger, requirements BIPSAPRequirements, options ...Option) (*BIPSAP, error) {
	b := &BIPSAP{
		DefaultRFormatter: NewDefaultRFormatter(),
		log:               localLog,
	}
	ApplyAppliers(options, b)
	optionsForParent := AddLeafTypeIfAbundant(options, b)
	var err error
	b.ServiceAccessPointContract, err = NewServiceAccessPoint(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating service access point")
	}
	b.requirements = requirements
	if _debug != nil {
		_debug("__init__ sap=%r", b.GetServiceID())
	}
	localLog.Debug().
		Interface("sapID", b.GetServiceElement()).
		Interface("requirements", requirements).
		Msg("NewBIPSAP")
	return b, nil
}

func (b *BIPSAP) String() string {
	return fmt.Sprintf("BIPSAP(SAP: %s)", b.ServiceAccessPointContract)
}

func (b *BIPSAP) SapIndication(args Args, kwArgs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("SapIndication")
	pdu := GetFromArgs[PDU](args, 0)
	if _debug != nil {
		_debug("sap_indication %r", pdu)
	}
	// this is a request initiated by the ASE, send this downstream
	return b.requirements.Request(args, kwArgs)
}

func (b *BIPSAP) SapConfirmation(args Args, kwArgs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("SapConfirmation")
	pdu := GetFromArgs[PDU](args, 0)
	if _debug != nil {
		_debug("sap_confirmation %r", pdu)
	}
	// this is a response from the ASE, send this downstream
	return b.requirements.Request(args, kwArgs)
}
