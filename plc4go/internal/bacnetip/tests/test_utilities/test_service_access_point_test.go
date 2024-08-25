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

package test_utilities

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/tests"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type EchoAccessPointRequirements interface {
	SapResponse(args bacnetip.Args, kwArgs bacnetip.KWArgs) error
}

type EchoAccessPoint struct {
	*bacnetip.ServiceAccessPoint

	echoAccessPointRequirements EchoAccessPointRequirements

	log zerolog.Logger
}

func NewEchoAccessPoint(localLog zerolog.Logger, serviceAccessPoint *bacnetip.ServiceAccessPoint) *EchoAccessPoint {
	e := &EchoAccessPoint{
		ServiceAccessPoint: serviceAccessPoint,

		log: localLog,
	}
	return e
}

func (e *EchoAccessPoint) SapIndication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	e.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("SapIndication")
	return e.echoAccessPointRequirements.SapResponse(args, bacnetip.NoKWArgs)
}

func (e *EchoAccessPoint) SapConfirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	e.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("SapConfirmation")
	return nil
}

type TrappedEchoAccessPoint struct {
	*bacnetip.ServiceAccessPoint
	*EchoAccessPoint
	*tests.TrappedServiceAccessPoint

	log zerolog.Logger
}

func NewTrappedEchoAccessPoint(localLog zerolog.Logger) (*TrappedEchoAccessPoint, error) {
	t := &TrappedEchoAccessPoint{}
	var err error
	t.ServiceAccessPoint, err = bacnetip.NewServiceAccessPoint(localLog, t)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating service access point")
	}
	t.EchoAccessPoint = NewEchoAccessPoint(localLog, t.ServiceAccessPoint)
	t.TrappedServiceAccessPoint, err = tests.NewTrappedServiceAccessPoint(localLog, t.EchoAccessPoint)
	if err != nil {
		return nil, errors.Wrap(err, "error creating trapped service access point")
	}
	t.EchoAccessPoint.echoAccessPointRequirements = t // TODO: isn't multi-inheritance easy to follow? At this point it is pretty confusing
	return t, nil
}

func (t *TrappedEchoAccessPoint) SapRequest(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	return t.TrappedServiceAccessPoint.SapRequest(args, kwargs)
}

func (t *TrappedEchoAccessPoint) SapIndication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	return t.TrappedServiceAccessPoint.SapIndication(args, kwargs)
}

func (t *TrappedEchoAccessPoint) SapResponse(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	return t.TrappedServiceAccessPoint.SapResponse(args, kwargs)
}

func (t *TrappedEchoAccessPoint) SapConfirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	return t.TrappedServiceAccessPoint.SapConfirmation(args, kwargs)
}

func (t *TrappedEchoAccessPoint) String() string {
	return "TrappedEchoAccessPoint"
}

type EchoServiceElementRequirements interface {
	Response(args bacnetip.Args, kwArgs bacnetip.KWArgs) error
}

type EchoServiceElement struct {
	*bacnetip.ApplicationServiceElement

	echoServiceElementRequirements EchoServiceElementRequirements

	log zerolog.Logger
}

func NewEchoServiceElement(localLog zerolog.Logger, applicationServiceElement *bacnetip.ApplicationServiceElement) *EchoServiceElement {
	e := &EchoServiceElement{
		ApplicationServiceElement: applicationServiceElement,
		log:                       localLog,
	}
	return e
}

func (e *EchoServiceElement) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	e.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")
	return e.echoServiceElementRequirements.Response(args, bacnetip.NoKWArgs)
}

func (e *EchoServiceElement) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	e.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")
	return nil
}

func (e *EchoServiceElement) String() string {
	return "EchoServiceElement" // TODO: fill with some useful
}

type TrappedEchoServiceElement struct {
	*bacnetip.ApplicationServiceElement
	*EchoServiceElement
	*tests.TrappedApplicationServiceElement
}

func NewTrappedEchoServiceElement(localLog zerolog.Logger) (*TrappedEchoServiceElement, error) {
	t := &TrappedEchoServiceElement{}
	var err error
	t.ApplicationServiceElement, err = bacnetip.NewApplicationServiceElement(localLog, t)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating service access point")
	}
	t.EchoServiceElement = NewEchoServiceElement(localLog, t.ApplicationServiceElement)
	t.TrappedApplicationServiceElement, err = tests.NewTrappedApplicationServiceElement(localLog, t.EchoServiceElement)
	if err != nil {
		return nil, errors.Wrap(err, "error creating trapped application service element")
	}
	t.EchoServiceElement.echoServiceElementRequirements = t // TODO: isn't multi-inheritance easy to follow? At this point it is pretty confusing
	return t, nil
}

func (t *TrappedEchoServiceElement) Request(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	return t.TrappedApplicationServiceElement.Request(args, kwargs)
}

func (t *TrappedEchoServiceElement) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	return t.TrappedApplicationServiceElement.Indication(args, kwargs)
}

func (t *TrappedEchoServiceElement) Response(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	return t.TrappedApplicationServiceElement.Response(args, kwargs)
}

func (t *TrappedEchoServiceElement) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	return t.TrappedApplicationServiceElement.Confirmation(args, kwargs)
}

func (t *TrappedEchoServiceElement) String() string {
	return "TrappedEchoServiceElement"
}

type TestApplicationSuite struct {
	suite.Suite

	sap *TrappedEchoAccessPoint
	ase *TrappedEchoServiceElement

	log zerolog.Logger
}

func (suite *TestApplicationSuite) SetupSuite() {
	suite.log = testutils.ProduceTestingLogger(suite.T())

	var err error
	suite.sap, err = NewTrappedEchoAccessPoint(suite.log)
	suite.Require().NoError(err)

	suite.ase, err = NewTrappedEchoServiceElement(suite.log)
	suite.Require().NoError(err)

	err = bacnetip.Bind(suite.log, suite.ase, suite.sap)
	suite.Require().NoError(err)
}

func (suite *TestApplicationSuite) TearDownSuite() {
}

func (suite *TestApplicationSuite) TestSapRequest() {
	// make a pdu
	pdu := bacnetip.NewPDU(tests.NewDummyMessage())

	// service access point is going to request something
	err := suite.sap.SapRequest(bacnetip.NewArgs(pdu), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)

	// make sure the request was sent and received
	suite.Equal(pdu, suite.sap.GetSapRequestSent())
	suite.Equal(pdu, suite.ase.GetIndicationReceived())

	// make sure the echo response was sent and received
	suite.Equal(pdu, suite.ase.GetResponseSent())
	suite.Equal(pdu, suite.sap.GetSapConfirmationReceived())
}

func (suite *TestApplicationSuite) TestAseRequest() {
	// make a pdu
	pdu := bacnetip.NewPDU(tests.NewDummyMessage())

	// service access point is going to request something
	err := suite.ase.Request(bacnetip.NewArgs(pdu), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)

	// make sure the request was sent and received
	suite.Equal(pdu, suite.ase.GetRequestSent())
	suite.Equal(pdu, suite.sap.GetSapIndicationReceived())

	// make sure the echo response was sent and received
	suite.Equal(pdu, suite.sap.GetSapResponseSent())
	suite.Equal(pdu, suite.ase.GetConfirmationReceived())
}

func TestApplicationService(t *testing.T) {
	suite.Run(t, new(TestApplicationSuite))
}
