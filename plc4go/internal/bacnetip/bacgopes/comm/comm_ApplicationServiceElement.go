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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type ApplicationServiceElement interface {
	ApplicationServiceElementContract
	ApplicationServiceElementRequirements
}

// ApplicationServiceElementRequirements provides a set of functions which must be overwritten by a sub struct
type ApplicationServiceElementRequirements interface {
	Confirmation(args Args, kwArgs KWArgs) error
	Indication(args Args, kwArgs KWArgs) error
}

// ApplicationServiceElementContract provides a set of functions which can be overwritten by a sub struct
type ApplicationServiceElementContract interface {
	utils.Serializable
	Request(args Args, kwArgs KWArgs) error
	Response(args Args, kwArgs KWArgs) error
	GetElementId() *int
	_setElementService(elementService ElementService)
	GetElementService() ElementService
}

// ElementService is required by ApplicationServiceElementContract to work properly
type ElementService interface {
	SapIndication(args Args, kwArgs KWArgs) error
	SapConfirmation(args Args, kwArgs KWArgs) error
}

//go:generate go tool plc4xGenerator -type=applicationServiceElement -prefix=comm_
type applicationServiceElement struct {
	elementID      *int
	elementService ElementService `asPtr:"true"`

	// arguments
	argASEExtension ApplicationServiceElement `ignore:"true"`

	log zerolog.Logger
}

func NewApplicationServiceElement(localLog zerolog.Logger, options ...Option) (ApplicationServiceElementContract, error) {
	a := &applicationServiceElement{
		log: localLog,
	}
	ApplyAppliers(options, a)
	if _debug != nil {
		_debug("__init__(%v)", a.argASEExtension)
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
			if err := Bind(a.log, a.argASEExtension, service); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return a, nil
}

func WithApplicationServiceElementAseID(aseID int, ase ApplicationServiceElement) GenericApplier[*applicationServiceElement] {
	if ase == nil {
		panic("saq required (completely build sap)") // TODO: might be hard because initialization not yet done
	}
	return WrapGenericApplier(func(s *applicationServiceElement) {
		s.elementID = &aseID
		s.argASEExtension = ase
	})
}

func (a *applicationServiceElement) GetElementId() *int {
	return a.elementID
}

func (a *applicationServiceElement) Request(args Args, kwArgs KWArgs) error {
	if _debug != nil {
		_debug("request(%v) %r %r", a.elementID, args, kwArgs)
	}
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Request")

	if a.elementService == nil {
		return errors.New("unbound application service element")
	}

	return a.elementService.SapIndication(args, kwArgs)
}

func (a *applicationServiceElement) Response(args Args, kwArgs KWArgs) error {
	if _debug != nil {
		_debug("response(%v) %r %r", a.elementID, args, kwArgs)
	}
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Response")

	if a.elementService == nil {
		return errors.New("unbound application service element")
	}

	return a.elementService.SapConfirmation(args, kwArgs)
}

func (a *applicationServiceElement) _setElementService(elementService ElementService) {
	a.elementService = elementService
}

func (a *applicationServiceElement) GetElementService() ElementService {
	return a.elementService
}
