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
	"context"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/app"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/appservice"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/iocb"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
)

// ApplicationLayerMessageCodec is a wrapper for MessageCodec which takes care of segmentation, retries etc.
//
//go:generate go tool plc4xGenerator -type=ApplicationLayerMessageCodec
type ApplicationLayerMessageCodec struct {
	bipSimpleApplication *app.BIPSimpleApplication
	messageCode          *MessageCodec
	deviceInfoCache      appservice.DeviceInfoCache `directSerialize:"true"`

	localAddress  *net.UDPAddr `stringer:"true"`
	remoteAddress *net.UDPAddr `stringer:"true"`

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewApplicationLayerMessageCodec(localLog zerolog.Logger, udpTransport *udp.Transport, transportUrl url.URL, options map[string][]string, localAddress *net.UDPAddr, remoteAddress *net.UDPAddr) (*ApplicationLayerMessageCodec, error) {
	// TODO: currently this is done by the BIP down below
	// Have the transport create a new transport-instance.
	//transportInstance, err := udpTransport.CreateTransportInstanceForLocalAddress(transportUrl, options, localAddress)
	//if err != nil {
	//	return nil, errors.Wrap(err, "error creating transport instance")
	//}
	a := &ApplicationLayerMessageCodec{
		localAddress:  localAddress,
		remoteAddress: remoteAddress,

		log: localLog,
	}
	address, err := pdu.NewAddress(comp.NewArgs(localAddress))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}
	// TODO: workaround for strange address parsing
	address.AddrTuple = pdu.NewAddressTuple(fmt.Sprintf("%d.%d.%d.%d", address.AddrAddress[0], address.AddrAddress[1], address.AddrAddress[2], address.AddrAddress[3]), *address.AddrPort)
	localDeviceObject, err := device.NewLocalDeviceObject(comp.NoArgs,
		comp.NewKWArgs(comp.KWNumberOfAPDURetries, func() *uint { retries := uint(10); return &retries }()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating local device object")
	}
	application, err := app.NewBIPSimpleApplication(localLog, localDeviceObject, *address, app.WithApplicationDeviceInfoCache(&a.deviceInfoCache))
	if err != nil {
		return nil, errors.Wrap(err, "error creating application")
	}
	a.bipSimpleApplication = application
	// TODO: this is currently done by the BIP
	//a.messageCode = NewMessageCodec(transportInstance)
	return a, nil
}

func (m *ApplicationLayerMessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *ApplicationLayerMessageCodec) Connect(ctx context.Context) error {
	// TODO: this is currently done by the BIP
	//	return m.messageCode.Connect(ctx)
	return nil
}

func (m *ApplicationLayerMessageCodec) Disconnect() error {
	if err := m.bipSimpleApplication.Close(); err != nil {
		m.log.Error().Err(err).Msg("error closing application")
	}
	return m.messageCode.Disconnect()
}

func (m *ApplicationLayerMessageCodec) IsRunning() bool {
	return m.messageCode.IsRunning()
}

func (m *ApplicationLayerMessageCodec) Send(ctx context.Context, message spi.Message) error {
	address, err := pdu.NewAddress(comp.NewArgs(m.remoteAddress))
	if err != nil {
		return err
	}
	_iocb, err := iocb.NewIOCB(m.log, pdu.NewPDU(comp.NewArgs(message), comp.NewKWArgs(comp.KWCPCIDestination, address)), address)
	if err != nil {
		return errors.Wrap(err, "error creating IOCB")
	}
	m.wg.Go(func() {
		m.wg.Go(func() {
			if err := m.bipSimpleApplication.RequestIO(_iocb); err != nil {
				m.log.Debug().Err(err).Msg("errored")
			}
		})
		_iocb.Wait()
		if err := _iocb.GetIOError(); err != nil {
			// TODO: handle error
			fmt.Printf("Err: %v\n", err)
		} else if _iocb.GetIOResponse() != nil {
			// TODO: response?
			fmt.Printf("Response: %v\n", _iocb.GetIOResponse())
		} else {
			// TODO: what now?
		}
	})
	return nil
}

func (m *ApplicationLayerMessageCodec) Expect(ctx context.Context, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) {
	// TODO: implement me
	panic("not yet implemented")
}

func (m *ApplicationLayerMessageCodec) SendRequest(ctx context.Context, message spi.Message, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	address, err := pdu.NewAddress(comp.NewArgs(m.remoteAddress))
	if err != nil {
		return err
	}
	_iocb, err := iocb.NewIOCB(m.log, pdu.NewPDU(comp.NewArgs(message), comp.NewKWArgs(comp.KWCPCIDestination, address)), address)
	if err != nil {
		return errors.Wrap(err, "error creating IOCB")
	}
	m.wg.Go(func() {
		m.wg.Go(func() {
			if err := m.bipSimpleApplication.RequestIO(_iocb); err != nil {
				m.log.Error().Err(err).Msg("errored")
			}
		})
		_iocb.Wait()
		if err := _iocb.GetIOError(); err != nil {
			if err := handleError(err); err != nil {
				m.log.Debug().Err(err).Msg("error handling error")
				return
			}
		} else if response := _iocb.GetIOResponse(); response != nil {
			// TODO: we wrap it into a BVLC for now. Once we change the Readers etc. to accept apdus we can remove that
			tempBVLC := model.NewBVLCOriginalUnicastNPDU(
				model.NewNPDU(
					0,
					model.NewNPDUControl(
						false,
						false,
						false,
						false,
						model.NPDUNetworkPriority_NORMAL_MESSAGE,
					),
					nil,
					nil,
					nil,
					nil,
					nil,
					nil,
					nil,
					nil,
					response.GetRootMessage().(model.APDU),
				),
			)
			if acceptsMessage(tempBVLC) {
				if err := handleMessage(
					tempBVLC,
				); err != nil {
					m.log.Debug().Err(err).Msg("error handling message")
					return
				}
			}
		} else {
			// TODO: what now?
		}
	})
	return nil
}

func (m *ApplicationLayerMessageCodec) GetDefaultIncomingMessageChannel() chan spi.Message {
	// TODO: this is currently done by the BIP
	//return m.messageCode.GetDefaultIncomingMessageChannel()
	return make(chan spi.Message)
}
