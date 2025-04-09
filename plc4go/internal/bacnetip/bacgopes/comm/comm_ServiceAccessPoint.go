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

package comm

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type ServiceAccessPoint interface {
	ServiceAccessPointContract
	ServiceAccessPointRequirements
}

// ServiceAccessPointContract provides a set of functions which can be overwritten by a sub struct
type ServiceAccessPointContract interface {
	fmt.Stringer
	utils.Serializable
	GetServiceID() *int
	SapRequest(Args, KWArgs) error
	SapResponse(Args, KWArgs) error
	_setServiceElement(serviceElement ServiceElement)
	GetServiceElement() ServiceElement
}

// ServiceAccessPointRequirements provides a set of functions which must be overwritten by a sub struct
type ServiceAccessPointRequirements interface {
	SapIndication(Args, KWArgs) error
	SapConfirmation(Args, KWArgs) error
}

type ServiceElement interface {
	Indication(args Args, kwArgs KWArgs) error
	Confirmation(args Args, kwArgs KWArgs) error
}

//go:generate go tool plc4xGenerator -type=serviceAccessPoint -prefix=comm_
type serviceAccessPoint struct {
	serviceID      *int
	serviceElement ServiceElement `asPtr:"true"`

	// arguments
	argSAPExtension ServiceAccessPoint `ignore:"true"`

	log zerolog.Logger
}

var _ ServiceAccessPointContract = (*serviceAccessPoint)(nil)

func NewServiceAccessPoint(localLog zerolog.Logger, options ...Option) (ServiceAccessPointContract, error) {
	s := &serviceAccessPoint{
		log: localLog,
	}
	ApplyAppliers(options, s)
	if _debug != nil {
		_debug("__init__(%v)", s.serviceID)
	}
	if s.serviceID != nil {
		sapID := *s.serviceID
		if _, ok := serviceMap[sapID]; ok {
			return nil, errors.Errorf("already a server %d", sapID)
		}
		serviceMap[sapID] = s

		// automatically bind
		if element, ok := elementMap[sapID]; ok {
			if element.elementService != nil {
				return nil, errors.Errorf("application service element %d already bound", sapID)
			}

			// Note: we need to pass the requirements (which should contain s as delegate) here
			if err := Bind(s.log, element, s.argSAPExtension); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return s, nil
}

func WithServiceAccessPointSapID(sapID int, sap ServiceAccessPoint) GenericApplier[*serviceAccessPoint] {
	if sap == nil {
		panic("saq required (completely build sap)") // TODO: might be hard because initialization not yet done
	}
	return WrapGenericApplier(func(s *serviceAccessPoint) {
		s.serviceID = &sapID
		s.argSAPExtension = sap
	})
}

func (s *serviceAccessPoint) GetServiceID() *int {
	return s.serviceID
}

func (s *serviceAccessPoint) SapRequest(args Args, kwArgs KWArgs) error {
	if _debug != nil {
		_debug("sap_request(%v) %r %r", s.serviceID, args, kwArgs)
	}
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Interface("serviceID", s.serviceID).Msg("SapRequest")

	if s.serviceElement == nil {
		return errors.New("unbound service access point")
	}
	return s.serviceElement.Indication(args, kwArgs)
}

func (s *serviceAccessPoint) SapResponse(args Args, kwArgs KWArgs) error {
	if _debug != nil {
		_debug("sap_response(%v) %r %r", s.serviceID, args, kwArgs)
	}
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Interface("serviceID", s.serviceID).Msg("SapResponse")

	if s.serviceElement == nil {
		return errors.New("unbound service access point")
	}
	return s.serviceElement.Confirmation(args, kwArgs)
}

func (s *serviceAccessPoint) _setServiceElement(serviceElement ServiceElement) {
	s.log.Trace().Interface("serviceElement", serviceMap).Msg("setting service element")
	s.serviceElement = serviceElement
}

func (s *serviceAccessPoint) GetServiceElement() ServiceElement {
	return s.serviceElement
}
