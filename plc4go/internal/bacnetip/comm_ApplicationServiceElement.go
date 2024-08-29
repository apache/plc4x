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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type ApplicationServiceElementRequirements interface {
	Confirmation(args Args, kwargs KWArgs) error
}

// ApplicationServiceElementContract provides a set of functions which can be overwritten by a sub struct
type ApplicationServiceElementContract interface {
	Request(args Args, kwargs KWArgs) error
	Indication(args Args, kwargs KWArgs) error
	Response(args Args, kwargs KWArgs) error
	Confirmation(args Args, kwargs KWArgs) error
	_setElementService(elementService ServiceAccessPointContract)
}

type ApplicationServiceElement struct {
	elementID      *int
	elementService ServiceAccessPointContract

	log zerolog.Logger
}

func NewApplicationServiceElement(localLog zerolog.Logger, requirements ApplicationServiceElementRequirements, opts ...func(*ApplicationServiceElement)) (*ApplicationServiceElement, error) {
	a := &ApplicationServiceElement{
		log: localLog,
	}
	for _, opt := range opts {
		opt(a)
	}

	if a.elementID != nil {
		aseID := *a.elementID
		if _, ok := elementMap[aseID]; ok {
			return nil, errors.Errorf("already an application service element %d", aseID)
		}
		elementMap[aseID] = a

		// automatically bind
		if service, ok := serviceMap[aseID]; ok {
			if service.serviceElement != nil {
				return nil, errors.Errorf("service access point %d already bound", aseID)
			}

			// Note: we need to pass the requirements (which should contain us as a delegate) here
			if err := Bind(localLog, requirements, service); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return a, nil
}

func WithApplicationServiceElementAseID(aseID int) func(*ApplicationServiceElement) {
	return func(s *ApplicationServiceElement) {
		s.elementID = &aseID
	}
}

func (a *ApplicationServiceElement) Request(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Request")

	if a.elementService == nil {
		return errors.New("unbound application service element")
	}

	return a.elementService.SapIndication(args, kwargs)
}

func (a *ApplicationServiceElement) Indication(Args, KWArgs) error {
	// TODO: we should remove this asap to check where we have actual caps because we can compile here
	panic("this should be implemented by outer struct")
}

func (a *ApplicationServiceElement) Response(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Response")

	if a.elementService == nil {
		return errors.New("unbound application service element")
	}

	return a.elementService.SapConfirmation(args, kwargs)
}

func (a *ApplicationServiceElement) Confirmation(Args, KWArgs) error {
	// TODO: we should remove this asap to check where we have actual caps because we can compile here
	panic("this should be implemented by outer struct")
}

func (a *ApplicationServiceElement) _setElementService(elementService ServiceAccessPointContract) {
	a.elementService = elementService
}
