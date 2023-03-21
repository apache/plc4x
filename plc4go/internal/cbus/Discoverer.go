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

package cbus

import (
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/transports/tcp"
	"net"
	"net/url"
	"sync"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/rs/zerolog/log"
)

type Discoverer struct {
	transportInstanceCreationQueue utils.Executor
	deviceScanningQueue            utils.Executor
}

func NewDiscoverer() *Discoverer {
	return &Discoverer{
		// TODO: maybe a dynamic executor would be better to not waste cycles when not in use
		transportInstanceCreationQueue: utils.NewFixedSizeExecutor(50, 100),
		deviceScanningQueue:            utils.NewFixedSizeExecutor(50, 100),
	}
}

func (d *Discoverer) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	d.transportInstanceCreationQueue.Start()
	d.deviceScanningQueue.Start()

	tcpTransport := tcp.NewTransport()

	allInterfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	// If no device is explicitly selected via option, simply use all of them
	// However if a discovery option is present to select a device by name, only
	// add those devices matching any of the given names.
	var interfaces []net.Interface
	deviceNames := options.FilterDiscoveryOptionsDeviceName(discoveryOptions)
	if len(deviceNames) > 0 {
		for _, curInterface := range allInterfaces {
			for _, deviceNameOption := range deviceNames {
				if curInterface.Name == deviceNameOption.GetDeviceName() {
					interfaces = append(interfaces, curInterface)
					break
				}
			}
		}
	} else {
		interfaces = allInterfaces
	}

	transportInstances := make(chan transports.TransportInstance)
	wg := &sync.WaitGroup{}
	// Iterate over all network devices of this system.
	for _, netInterface := range interfaces {
		addrs, err := netInterface.Addrs()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(netInterface net.Interface) {
			defer func() { wg.Done() }()
			// Iterate over all addresses the current interface has configured
			for _, addr := range addrs {
				var ipv4Addr net.IP
				switch addr.(type) {
				// If the device is configured to communicate with a subnet
				case *net.IPNet:
					ipv4Addr = addr.(*net.IPNet).IP.To4()

				// If the device is configured for a point-to-point connection
				case *net.IPAddr:
					ipv4Addr = addr.(*net.IPAddr).IP.To4()
				}

				// If we found an IPv4 address and this is not a loopback address,
				// add it to the list of devices we will open ports and send discovery
				// messages from.
				if ipv4Addr == nil || ipv4Addr.IsLoopback() {
					continue
				}
				addresses, err := utils.GetIPAddresses(ctx, netInterface, false)
				if err != nil {
					log.Warn().Err(err).Msgf("Can't get addresses for %v", netInterface)
					continue
				}
				wg.Add(1)
				go func() {
					defer func() { wg.Done() }()
					for ip := range addresses {
						log.Trace().Msgf("Handling found ip %v", ip)
						d.transportInstanceCreationQueue.Submit(ctx, 0, d.createTransportInstanceDispatcher(ctx, wg, ip, tcpTransport, transportInstances))
					}
				}()
			}
		}(netInterface)
	}
	go func() {
		wg.Wait()
		log.Trace().Msg("Closing transport instance channel")
		close(transportInstances)
	}()

	for transportInstance := range transportInstances {
		d.deviceScanningQueue.Submit(ctx, 0, d.createDeviceScanDispatcher(transportInstance.(*tcp.TransportInstance), callback))
	}
	return nil
}

func (d *Discoverer) createTransportInstanceDispatcher(ctx context.Context, wg *sync.WaitGroup, ip net.IP, tcpTransport *tcp.Transport, transportInstances chan transports.TransportInstance) utils.Runnable {
	wg.Add(1)
	return func() {
		defer wg.Done()
		// Create a new "connection" (Actually open a local udp socket and target outgoing packets to that address)
		var connectionUrl url.URL
		{
			connectionUrlParsed, err := url.Parse(fmt.Sprintf("tcp://%s:%d", ip, readWriteModel.CBusConstants_CBUSTCPDEFAULTPORT))
			if err != nil {
				log.Error().Err(err).Msgf("Error parsing url for lookup")
				return
			}
			connectionUrl = *connectionUrlParsed
		}

		transportInstance, err := tcpTransport.CreateTransportInstance(connectionUrl, nil)
		if err != nil {
			log.Error().Err(err).Msgf("Error creating transport instance")
			return
		}
		log.Trace().Msgf("trying %v", connectionUrl)
		err = transportInstance.ConnectWithContext(ctx)
		if err != nil {
			secondErr := transportInstance.ConnectWithContext(ctx)
			if secondErr != nil {
				log.Trace().Err(err).Msgf("Error connecting transport instance")
				return
			}
		}
		transportInstances <- transportInstance
	}
}

func (d *Discoverer) createDeviceScanDispatcher(tcpTransportInstance *tcp.TransportInstance, callback func(event apiModel.PlcDiscoveryItem)) utils.Runnable {
	return func() {
		// Create a codec for sending and receiving messages.
		codec := NewMessageCodec(tcpTransportInstance)
		// Explicitly start the worker
		if err := codec.Connect(); err != nil {
			log.Debug().Err(err).Msg("Error connecting")
			return
		}

		// Prepare the discovery packet data
		cBusOptions := readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, true)
		requestContext := readWriteModel.NewRequestContext(false)
		calData := readWriteModel.NewCALDataIdentify(readWriteModel.Attribute_Manufacturer, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, requestContext)
		alpha := readWriteModel.NewAlpha('x')
		request := readWriteModel.NewRequestDirectCommandAccess(calData, alpha, 0x0, nil, nil, readWriteModel.RequestType_DIRECT_COMMAND, readWriteModel.NewRequestTermination(), cBusOptions)
		cBusMessageToServer := readWriteModel.NewCBusMessageToServer(request, requestContext, cBusOptions)
		// Send the search request.
		if err := codec.Send(cBusMessageToServer); err != nil {
			log.Debug().Err(err).Msgf("Error sending message:\n%s", cBusMessageToServer)
			return
		}
		// Keep on reading responses till the timeout is done.
		// TODO: Make this configurable
		timeout := time.NewTimer(time.Second * 1)
		timeout.Stop()
		for start := time.Now(); time.Since(start) < time.Second*5; {
			timeout.Reset(time.Second * 1)
			select {
			case receivedMessage := <-codec.GetDefaultIncomingMessageChannel():
				if !timeout.Stop() {
					<-timeout.C
				}
				cbusMessage, ok := receivedMessage.(readWriteModel.CBusMessage)
				if !ok {
					continue
				}
				messageToClient, ok := cbusMessage.(readWriteModel.CBusMessageToClient)
				if !ok {
					continue
				}
				replyOrConfirmationConfirmation, ok := messageToClient.GetReply().(readWriteModel.ReplyOrConfirmationConfirmationExactly)
				if !ok {
					continue
				}
				if receivedAlpha := replyOrConfirmationConfirmation.GetConfirmation().GetAlpha(); receivedAlpha != nil && alpha.GetCharacter() != receivedAlpha.GetCharacter() {
					continue
				}
				embeddedReply, ok := replyOrConfirmationConfirmation.GetEmbeddedReply().(readWriteModel.ReplyOrConfirmationReplyExactly)
				if !ok {
					continue
				}
				encodedReply, ok := embeddedReply.GetReply().(readWriteModel.ReplyEncodedReplyExactly)
				if !ok {
					continue
				}
				encodedReplyCALReply, ok := encodedReply.GetEncodedReply().(readWriteModel.EncodedReplyCALReplyExactly)
				if !ok {
					continue
				}
				calDataIdentifyReply, ok := encodedReplyCALReply.GetCalReply().GetCalData().(readWriteModel.CALDataIdentifyReplyExactly)
				if !ok {
					continue
				}
				identifyReplyCommand, ok := calDataIdentifyReply.GetIdentifyReplyCommand().(readWriteModel.IdentifyReplyCommandManufacturerExactly)
				if !ok {
					continue
				}
				var remoteUrl url.URL
				{
					// TODO: we could check for the exact response
					remoteUrlParse, err := url.Parse(fmt.Sprintf("tcp://%s", tcpTransportInstance.RemoteAddress))
					if err != nil {
						log.Error().Err(err).Msg("Error creating url")
						continue
					}
					remoteUrl = *remoteUrlParse
				}
				// TODO: manufacturer + type would be good but this means two requests then
				deviceName := identifyReplyCommand.GetManufacturerName()
				discoveryEvent := &internalModel.DefaultPlcDiscoveryItem{
					ProtocolCode:  "c-bus",
					TransportCode: "tcp",
					TransportUrl:  remoteUrl,
					Options:       nil,
					Name:          deviceName,
				}
				// Pass the event back to the callback
				callback(discoveryEvent)
				continue
			case <-timeout.C:
				timeout.Stop()
				continue
			}
		}
	}
}
