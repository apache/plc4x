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
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// ServiceAccessPointContract provides a set of functions which can be overwritten by a sub struct
type ServiceAccessPointContract interface {
	fmt.Stringer
	SapConfirmation(Args, KWArgs) error
	SapRequest(Args, KWArgs) error
	SapIndication(Args, KWArgs) error
	SapResponse(Args, KWArgs) error
	_setServiceElement(serviceElement ApplicationServiceElementContract)
}

type ServiceAccessPoint struct {
	serviceID      *int
	serviceElement ApplicationServiceElementContract

	log zerolog.Logger
}

var _ ServiceAccessPointContract = (*ServiceAccessPoint)(nil)

func NewServiceAccessPoint(localLog zerolog.Logger, serviceAccessPointContract ServiceAccessPointContract, opts ...func(point *ServiceAccessPoint)) (*ServiceAccessPoint, error) {
	s := &ServiceAccessPoint{
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
			if err := Bind(localLog, element, serviceAccessPointContract); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return s, nil
}

func WithServiceAccessPointSapID(sapID int) func(*ServiceAccessPoint) {
	return func(s *ServiceAccessPoint) {
		s.serviceID = &sapID
	}
}

func (s *ServiceAccessPoint) String() string {
	serviceID := "-"
	if s.serviceID != nil {
		serviceID = strconv.Itoa(*s.serviceID)
	}
	return fmt.Sprintf("ServiceAccessPoint(serviceID:%v, serviceElement: %s)", serviceID, s.serviceElement)
}

func (s *ServiceAccessPoint) SapRequest(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Interface("serviceID", s.serviceID).Msg("SapRequest")

	if s.serviceElement == nil {
		return errors.New("unbound service access point")
	}
	return s.serviceElement.Indication(args, kwargs)
}

func (s *ServiceAccessPoint) SapIndication(Args, KWArgs) error {
	// TODO: we should remove this asap to check where we have actual caps because we can compile here
	panic("this should be implemented by outer struct")
}

func (s *ServiceAccessPoint) SapResponse(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Interface("serviceID", s.serviceID).Msg("SapResponse")

	if s.serviceElement == nil {
		return errors.New("unbound service access point")
	}
	return s.serviceElement.Confirmation(args, kwargs)
}

func (s *ServiceAccessPoint) SapConfirmation(Args, KWArgs) error {
	// TODO: we should remove this asap to check where we have actual caps because we can compile here
	panic("this should be implemented by outer struct")
}

func (s *ServiceAccessPoint) _setServiceElement(serviceElement ApplicationServiceElementContract) {
	s.log.Trace().Interface("serviceElement", serviceMap).Msg("setting service element")
	s.serviceElement = serviceElement
}
