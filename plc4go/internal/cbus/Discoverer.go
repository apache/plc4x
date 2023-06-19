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
	"github.com/apache/plc4x/plc4go/spi/pool"
	"github.com/apache/plc4x/plc4go/spi/transports/tcp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net"
	"net/url"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Discoverer struct {
	transportInstanceCreationWorkItemId atomic.Int32
	transportInstanceCreationQueue      pool.Executor
	deviceScanningWorkItemId            atomic.Int32
	deviceScanningQueue                 pool.Executor

	log zerolog.Logger
}

func NewDiscoverer(_options ...options.WithOption) *Discoverer {
	return &Discoverer{
		// TODO: maybe a dynamic executor would be better to not waste cycles when not in use
		transportInstanceCreationQueue: pool.NewFixedSizeExecutor(50, 100, _options...),
		deviceScanningQueue:            pool.NewFixedSizeExecutor(50, 100, _options...),

		log: options.ExtractCustomLogger(_options...),
	}
}

func (d *Discoverer) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	d.transportInstanceCreationQueue.Start()
	d.deviceScanningQueue.Start()

	deviceNames := d.extractDeviceNames(discoveryOptions...)
	interfaces, err := addressProviderRetriever(d.log, deviceNames)
	if err != nil {
		return errors.Wrap(err, "error getting addresses")
	}
	if d.log.Debug().Enabled() {
		for _, provider := range interfaces {
			d.log.Debug().Msgf("Discover on %s", provider)
			d.log.Trace().Msgf("Discover on %#v", provider.containedInterface())
		}
	}

	transportInstances := make(chan transports.TransportInstance)
	wg := &sync.WaitGroup{}
	tcpTransport := tcp.NewTransport()
	// Iterate over all network devices of this system.
	for _, netInterface := range interfaces {
		interfaceLog := d.log.With().Stringer("interface", netInterface).Logger()
		interfaceLog.Debug().Msg("Scanning")
		addrs, err := netInterface.Addrs()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(netInterface addressProvider, interfaceLog zerolog.Logger) {
			defer func() {
				if err := recover(); err != nil {
					interfaceLog.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
				}
			}()
			defer func() { wg.Done() }()
			// Iterate over all addresses the current interface has configured
			for _, addr := range addrs {
				addressLogger := interfaceLog.With().Stringer("address", addr).Logger()
				addressLogger.Debug().Msg("looking into")
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
				addresses, err := utils.GetIPAddresses(d.log, ctx, netInterface.containedInterface(), false)
				if err != nil {
					addressLogger.Warn().Err(err).Msgf("Can't get addresses for %v", netInterface)
					continue
				}
				wg.Add(1)
				go func(addressLogger zerolog.Logger) {
					defer func() {
						if err := recover(); err != nil {
							addressLogger.Error().Msgf("panic-ed %v. Stack: %s; ", err, debug.Stack())
						}
					}()
					defer func() { wg.Done() }()
					for ip := range addresses {
						addressLogger.Trace().Msgf("Handling found ip %v", ip)
						d.transportInstanceCreationQueue.Submit(
							ctx,
							d.transportInstanceCreationWorkItemId.Add(1),
							d.createTransportInstanceDispatcher(
								ctx,
								wg,
								ip,
								tcpTransport,
								transportInstances,
								readWriteModel.CBusConstants_CBUSTCPDEFAULTPORT,
								addressLogger,
							),
						)
					}
				}(addressLogger)
			}
		}(netInterface, interfaceLog)
	}
	go func() {
		wg.Wait()
		d.log.Trace().Msg("Closing transport instance channel")
		close(transportInstances)
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				d.log.Error().Msgf("panic-ed %v. Stack: %s; ", err, debug.Stack())
			}
		}()
		deviceScanWg := sync.WaitGroup{}
		for transportInstance := range transportInstances {
			d.log.Debug().Stringer("transportInstance", transportInstance).Msg("submitting device scan")
			completionFuture := d.deviceScanningQueue.Submit(ctx, d.deviceScanningWorkItemId.Add(1), d.createDeviceScanDispatcher(transportInstance.(*tcp.TransportInstance), callback))
			deviceScanWg.Add(1)
			go func() {
				defer deviceScanWg.Done()
				if err := completionFuture.AwaitCompletion(context.TODO()); err != nil {
					d.log.Debug().Err(err).Msg("error waiting for completion")
				}
			}()
			deviceScanWg.Wait()
			d.log.Info().Msg("Discovery done")
			d.transportInstanceCreationQueue.Stop()
			d.deviceScanningQueue.Stop()
			// TODO: do we maybe want a callback for that? As option for example
		}
	}()
	return nil
}

func (d *Discoverer) createTransportInstanceDispatcher(ctx context.Context, wg *sync.WaitGroup, ip net.IP, tcpTransport *tcp.Transport, transportInstances chan transports.TransportInstance, cBusPort uint16, addressLogger zerolog.Logger) pool.Runnable {
	wg.Add(1)
	return func() {
		defer wg.Done()
		// Create a new "connection" (Actually open a local udp socket and target outgoing packets to that address)
		var connectionUrl url.URL
		{
			connectionUrlParsed, err := url.Parse(fmt.Sprintf("tcp://%s:%d", ip, cBusPort))
			if err != nil {
				addressLogger.Error().Err(err).Msgf("Error parsing url for lookup")
				return
			}
			connectionUrl = *connectionUrlParsed
		}

		transportInstance, err := tcpTransport.CreateTransportInstance(connectionUrl, nil)
		if err != nil {
			addressLogger.Error().Err(err).Msgf("Error creating transport instance")
			return
		}
		addressLogger.Trace().Msgf("trying %v", connectionUrl)
		err = transportInstance.ConnectWithContext(ctx)
		if err != nil {
			secondErr := transportInstance.ConnectWithContext(ctx)
			if secondErr != nil {
				addressLogger.Trace().Err(err).Msgf("Error connecting transport instance")
				return
			}
		}
		addressLogger.Debug().Msgf("Adding transport instance to scan %v", transportInstance)
		transportInstances <- transportInstance
	}
}

func (d *Discoverer) createDeviceScanDispatcher(tcpTransportInstance *tcp.TransportInstance, callback func(event apiModel.PlcDiscoveryItem)) pool.Runnable {
	return func() {
		transportInstanceLogger := d.log.With().Stringer("transportInstance", tcpTransportInstance).Logger()
		transportInstanceLogger.Debug().Msgf("Scanning %v", tcpTransportInstance)
		// Create a codec for sending and receiving messages.
		codec := NewMessageCodec(tcpTransportInstance, options.WithCustomLogger(d.log))
		// Explicitly start the worker
		if err := codec.ConnectWithContext(context.TODO()); err != nil {
			transportInstanceLogger.Debug().Err(err).Msg("Error connecting")
			return
		}
		defer func() {
			// Disconnect codec when done
			d.log.Debug().Msg("Shutting down codec")
			if err := codec.Disconnect(); err != nil {
				d.log.Warn().Err(err).Msg("Error disconnecting codec")
			}
			d.log.Trace().Msg("done")
		}()

		// Prepare the discovery packet data
		cBusOptions := readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, true)
		requestContext := readWriteModel.NewRequestContext(false)
		calData := readWriteModel.NewCALDataIdentify(readWriteModel.Attribute_Manufacturer, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, requestContext)
		alpha := readWriteModel.NewAlpha('x')
		request := readWriteModel.NewRequestDirectCommandAccess(calData, alpha, 0x0, nil, nil, readWriteModel.RequestType_DIRECT_COMMAND, readWriteModel.NewRequestTermination(), cBusOptions)
		cBusMessageToServer := readWriteModel.NewCBusMessageToServer(request, requestContext, cBusOptions)
		// Send the search request.
		if err := codec.Send(cBusMessageToServer); err != nil {
			transportInstanceLogger.Debug().Err(err).Msgf("Error sending message:\n%s", cBusMessageToServer)
			return
		}
		// Keep on reading responses till the timeout is done.
		// TODO: Make this configurable
		timeout := time.NewTimer(1 * time.Second)
		defer utils.CleanupTimer(timeout)
		for start := time.Now(); time.Since(start) < 5*time.Second; {
			timeout.Reset(1 * time.Second)
			select {
			case receivedMessage := <-codec.GetDefaultIncomingMessageChannel():
				// Cleanup, going to be resetted again
				utils.CleanupTimer(timeout)
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
						transportInstanceLogger.Error().Err(err).Msg("Error creating url")
						continue
					}
					remoteUrl = *remoteUrlParse
				}
				// TODO: manufacturer + type would be good but this means two requests then
				deviceName := identifyReplyCommand.GetManufacturerName()
				discoveryEvent := spiModel.NewDefaultPlcDiscoveryItem(
					"c-bus",
					"tcp",
					remoteUrl,
					nil,
					deviceName,
					nil,
				)
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

func (d *Discoverer) extractDeviceNames(discoveryOptions ...options.WithDiscoveryOption) []string {
	deviceNamesOptions := options.FilterDiscoveryOptionsDeviceName(discoveryOptions)
	deviceNames := make([]string, len(deviceNamesOptions))
	for i, option := range deviceNamesOptions {
		deviceNames[i] = option.GetDeviceName()
	}
	return deviceNames
}

func (d *Discoverer) Close() error {
	d.transportInstanceCreationQueue.Stop()
	d.deviceScanningQueue.Stop()
	return nil
}

// addressProvider is used to make discover testable
type addressProvider interface {
	fmt.Stringer
	// Addrs is implemented by net.Interface#Addrs
	Addrs() ([]net.Addr, error)
	name() string
	containedInterface() net.Interface
}

// wrappedInterface extends net.Interface with name() and containedInterface()
type wrappedInterface struct {
	*net.Interface
}

func (w *wrappedInterface) name() string {
	return w.Interface.Name
}

func (w *wrappedInterface) containedInterface() net.Interface {
	return *w.Interface
}

func (w *wrappedInterface) String() string {
	return w.name()
}

// allInterfaceRetriever can be exchanged in tests
var allInterfaceRetriever = func(localLog zerolog.Logger) ([]addressProvider, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve all interfaces")
	}
	localLog.Debug().Msgf("Mapping %d interfaces", len(interfaces))
	addressProviders := make([]addressProvider, len(interfaces))
	for i, networkInterface := range interfaces {
		var copyInterface = networkInterface
		addressProviders[i] = &wrappedInterface{&copyInterface}
	}
	return addressProviders, nil
}

// addressProviderRetriever can be exchanged in tests
var addressProviderRetriever = func(localLog zerolog.Logger, deviceNames []string) ([]addressProvider, error) {
	allInterfaces, err := allInterfaceRetriever(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error getting all interfaces")
	}

	// If no device is explicitly selected via option, simply use all of them
	// However if a discovery option is present to select a device by name, only
	// add those devices matching any of the given names.
	if len(deviceNames) <= 0 {
		localLog.Info().Msgf("no devices selected, use all devices (%d)", len(allInterfaces))
		return allInterfaces, nil
	}

	var interfaces []addressProvider
	for _, curInterface := range allInterfaces {
		for _, deviceName := range deviceNames {
			if curInterface.name() == deviceName {
				interfaces = append(interfaces, curInterface)
				break
			}
		}
	}
	return interfaces, nil
}
