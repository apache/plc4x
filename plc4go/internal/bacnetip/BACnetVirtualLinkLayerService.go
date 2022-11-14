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
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type UDPMultiplexer struct {
	annexJ interface{}
}

func (m *UDPMultiplexer) Close() error {
	panic("implement me")
}

func NewUDPMultiplexer(address interface{}, noBroadcast bool) (*UDPMultiplexer, error) {
	log.Debug().Msgf("NewUDPMultiplexer %v noBroadcast=%t", address, noBroadcast)
	u := &UDPMultiplexer{}

	// TODO: plumb later
	return u, nil
}

type AnnexJCodec struct {
	*Client
	*Server
}

func NewAnnexJCodec(cid *int, sid *int) (*AnnexJCodec, error) {
	log.Debug().Msgf("NewAnnexJCodec cid=%d sid=%d", cid, sid)
	a := &AnnexJCodec{}
	client, err := NewClient(cid, a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	a.Client = client
	server, err := NewServer(sid, a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	a.Server = server
	return a, nil
}

func (b *AnnexJCodec) Indication(apdu readWriteModel.APDU) error {
	panic("we need to implement this with  generics as we handle npdu not apdu here")
}

func (b *AnnexJCodec) Confirmation(apdu readWriteModel.APDU) error {
	panic("we need to implement this with  generics as we handle npdu not apdu here")
}

type _BIPSAP interface {
	_ServiceAccessPoint
	_Client
}

type BIPSAP struct {
	*ServiceAccessPoint
	rootStruct _BIPSAP
}

func NewBIPSAP(sapID *int, rootStruct _BIPSAP) (*BIPSAP, error) {
	log.Debug().Msgf("NewBIPSAP sapID=%d", sapID)
	b := &BIPSAP{}
	serviceAccessPoint, err := NewServiceAccessPoint(sapID, rootStruct)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating service access point")
	}
	b.ServiceAccessPoint = serviceAccessPoint
	b.rootStruct = rootStruct
	return b, nil
}

func (b *BIPSAP) SapIndication(apdu readWriteModel.APDU, pduDestination []byte) error {
	log.Debug().Msgf("SapIndication\n%s\n%s", apdu, pduDestination)
	// TODO: what to do with the destination?
	// this is a request initiated by the ASE, send this downstream
	return b.rootStruct.Request(apdu)
}

func (b *BIPSAP) SapConfirmation(apdu readWriteModel.APDU, pduDestination []byte) error {
	log.Debug().Msgf("SapConfirmation\n%s\n%s", apdu, pduDestination)
	// TODO: what to do with the destination?
	// this is a response from the ASE, send this downstream
	return b.rootStruct.Request(apdu)
}

type BIPSimple struct {
	*BIPSAP
	*Client
	*Server
}

func NewBIPSimple(sapID *int, cid *int, sid *int) (*BIPSimple, error) {
	log.Debug().Msgf("NewBIPSimple sapID=%d cid=%d sid=%d", sapID, cid, sid)
	b := &BIPSimple{}
	bipsap, err := NewBIPSAP(sapID, b)
	if err != nil {
		return nil, errors.Wrap(err, "error creating bisap")
	}
	b.BIPSAP = bipsap
	client, err := NewClient(cid, b)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	b.Client = client
	server, err := NewServer(sid, b)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	b.Server = server
	return b, nil
}

func (b *BIPSimple) Indication(apdu readWriteModel.APDU) error {
	panic("we need to implement this with  generics as we handle npdu not apdu here")
}

func (b *BIPSimple) Response(apdu readWriteModel.APDU) error {
	panic("we need to implement this with  generics as we handle npdu not apdu here")
}
