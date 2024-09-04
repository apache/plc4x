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

package bacgopes

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type ServiceAccessPoint interface {
	ServiceAccessPointContract
	ServiceAccessPointRequirements
}

// ServiceAccessPointContract provides a set of functions which can be overwritten by a sub struct
type ServiceAccessPointContract interface {
	fmt.Stringer
	SapRequest(Args, KWArgs) error
	SapResponse(Args, KWArgs) error
	_setServiceElement(serviceElement ServiceElement)
	_getServiceElement() ServiceElement
}

// ServiceAccessPointRequirements provides a set of functions which must be overwritten by a sub struct
type ServiceAccessPointRequirements interface {
	SapIndication(Args, KWArgs) error
	SapConfirmation(Args, KWArgs) error
}

type ServiceElement interface {
	Indication(args Args, kwargs KWArgs) error
	Confirmation(args Args, kwargs KWArgs) error
}

type serviceAccessPoint struct {
	serviceID      *int
	serviceElement ServiceElement

	// arguments
	argSAPExtension ServiceAccessPoint

	log zerolog.Logger
}

var _ ServiceAccessPointContract = (*serviceAccessPoint)(nil)

func NewServiceAccessPoint(localLog zerolog.Logger, opts ...func(point *serviceAccessPoint)) (ServiceAccessPointContract, error) {
	s := &serviceAccessPoint{
		log: localLog,
	}
	for _, opt := range opts {
		opt(s)
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
			if err := Bind(localLog, element, s.argSAPExtension); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return s, nil
}

func WithServiceAccessPointSapID(sapID int, sap ServiceAccessPoint) func(*serviceAccessPoint) {
	if sap == nil {
		panic("saq required (completely build sap)") // TODO: might be hard because initialization not yet done
	}
	return func(s *serviceAccessPoint) {
		s.serviceID = &sapID
		s.argSAPExtension = sap
	}
}

func (s *serviceAccessPoint) String() string {
	serviceID := "-"
	if s.serviceID != nil {
		serviceID = strconv.Itoa(*s.serviceID)
	}
	return fmt.Sprintf("ServiceAccessPoint(serviceID:%v, serviceElement: %s)", serviceID, s.serviceElement)
}

func (s *serviceAccessPoint) SapRequest(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Interface("serviceID", s.serviceID).Msg("SapRequest")

	if s.serviceElement == nil {
		return errors.New("unbound service access point")
	}
	return s.serviceElement.Indication(args, kwargs)
}

func (s *serviceAccessPoint) SapResponse(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Interface("serviceID", s.serviceID).Msg("SapResponse")

	if s.serviceElement == nil {
		return errors.New("unbound service access point")
	}
	return s.serviceElement.Confirmation(args, kwargs)
}

func (s *serviceAccessPoint) _setServiceElement(serviceElement ServiceElement) {
	s.log.Trace().Interface("serviceElement", serviceMap).Msg("setting service element")
	s.serviceElement = serviceElement
}

func (s *serviceAccessPoint) _getServiceElement() ServiceElement {
	return s.serviceElement
}
