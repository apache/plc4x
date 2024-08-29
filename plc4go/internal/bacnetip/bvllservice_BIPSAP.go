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

package bacnetip

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type BIPSAPRequirements interface {
	ServiceAccessPointContract
	Client
}

type BIPSAP struct {
	*ServiceAccessPoint
	rootStruct BIPSAPRequirements

	// pass through args
	argSapID *int

	log zerolog.Logger
}

func NewBIPSAP(localLog zerolog.Logger, bipSapRequirements BIPSAPRequirements, opts ...func(*BIPSAP)) (*BIPSAP, error) {
	b := &BIPSAP{
		log: localLog,
	}
	for _, opt := range opts {
		opt(b)
	}
	localLog.Debug().
		Interface("sapID", b.argSapID).
		Interface("bipSapRequirements", bipSapRequirements).
		Msg("NewBIPSAP")
	serviceAccessPoint, err := NewServiceAccessPoint(localLog, bipSapRequirements, func(point *ServiceAccessPoint) {
		point.serviceID = b.argSapID
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error creating service access point")
	}
	b.ServiceAccessPoint = serviceAccessPoint
	b.rootStruct = bipSapRequirements
	return b, nil
}

func WithBIPSAPSapID(sapID int) func(*BIPSAP) {
	return func(point *BIPSAP) {
		point.argSapID = &sapID
	}
}

func (b *BIPSAP) String() string {
	return fmt.Sprintf("BIPSAP(SAP: %s)", b.ServiceAccessPoint)
}

func (b *BIPSAP) SapIndication(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapIndication")
	// this is a request initiated by the ASE, send this downstream
	return b.rootStruct.Request(args, kwargs)
}

func (b *BIPSAP) SapConfirmation(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapConfirmation")
	// this is a response from the ASE, send this downstream
	return b.rootStruct.Request(args, kwargs)
}
