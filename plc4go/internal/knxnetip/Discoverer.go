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

package knxnetip

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/rs/zerolog/log"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
)

type Discoverer struct {
	transportInstanceCreationWorkItemId atomic.Int32
	transportInstanceCreationQueue      utils.Executor
	deviceScanningWorkItemId            atomic.Int32
	deviceScanningQueue                 utils.Executor
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

	udpTransport := udp.NewTransport()

	// Create a connection string for the KNX broadcast discovery address.
	connectionUrl, err := url.Parse("udp://224.0.23.12:3671")
	if err != nil {
		return err
	}

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
			defer func() {
				if err := recover(); err != nil {
					log.Error().Msgf("panic-ed %v", err)
				}
			}()
			defer func() { wg.Done() }()
			// Iterate over all addresses the current interface has configured
			// For KNX we're only interested in IPv4 addresses, as it doesn't
			// seem to work with IPv6.
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
				d.transportInstanceCreationQueue.Submit(ctx, d.transportInstanceCreationWorkItemId.Add(1), d.createTransportInstanceDispatcher(ctx, wg, connectionUrl, ipv4Addr, udpTransport, transportInstances))
			}
		}(netInterface)
	}
	go func() {
		wg.Wait()
		log.Trace().Msg("Closing transport instance channel")
		close(transportInstances)
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Msgf("panic-ed %v", err)
			}
		}()
		for transportInstance := range transportInstances {
			d.deviceScanningQueue.Submit(ctx, d.deviceScanningWorkItemId.Add(1), d.createDeviceScanDispatcher(transportInstance.(*udp.TransportInstance), callback))
		}
	}()
	return nil
}

func (d *Discoverer) createTransportInstanceDispatcher(ctx context.Context, wg *sync.WaitGroup, connectionUrl *url.URL, ipv4Addr net.IP, udpTransport *udp.Transport, transportInstances chan transports.TransportInstance) utils.Runnable {
	wg.Add(1)
	return func() {
		defer wg.Done()
		// Create a new "connection" (Actually open a local udp socket and target outgoing packets to that address)
		transportInstance, err :=
			udpTransport.CreateTransportInstanceForLocalAddress(*connectionUrl, nil,
				&net.UDPAddr{IP: ipv4Addr, Port: 0})
		if err != nil {
			log.Error().Err(err).Msg("error creating transport instance")
			return
		}
		err = transportInstance.ConnectWithContext(ctx)
		if err != nil {
			log.Debug().Err(err).Msg("Error Connecting")
			return
		}
		log.Debug().Msgf("Adding transport instance to scan %v", transportInstance)
		transportInstances <- transportInstance
	}
}

func (d *Discoverer) createDeviceScanDispatcher(udpTransportInstance *udp.TransportInstance, callback func(event apiModel.PlcDiscoveryItem)) utils.Runnable {
	return func() {
		log.Debug().Msgf("Scanning %v", udpTransportInstance)
		// Create a codec for sending and receiving messages.
		codec := NewMessageCodec(udpTransportInstance, nil)
		// Explicitly start the worker
		if err := codec.Connect(); err != nil {
			log.Error().Err(err).Msg("Error connecting")
			return
		}

		localAddress := udpTransportInstance.LocalAddress
		localAddr := driverModel.NewIPAddress(localAddress.IP)

		// Prepare the discovery packet data
		discoveryEndpoint := driverModel.NewHPAIDiscoveryEndpoint(
			driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port))
		searchRequestMessage := driverModel.NewSearchRequest(discoveryEndpoint)
		// Send the search request.
		if err := codec.Send(searchRequestMessage); err != nil {
			log.Debug().Err(err).Msgf("Error sending message:\n%s", searchRequestMessage)
			return
		}
		// Keep on reading responses till the timeout is done.
		// TODO: Make this configurable
		timeout := time.NewTimer(time.Second * 1)
		timeout.Stop()
		for start := time.Now(); time.Since(start) < time.Second*5; {
			timeout.Reset(time.Second * 1)
			select {
			case message := <-codec.GetDefaultIncomingMessageChannel():
				{
					if !timeout.Stop() {
						<-timeout.C
					}
					searchResponse := message.(driverModel.SearchResponse)
					if searchResponse != nil {
						addr := searchResponse.GetHpaiControlEndpoint().GetIpAddress().GetAddr()
						remoteUrl, err := url.Parse(fmt.Sprintf("udp://%d.%d.%d.%d:%d",
							addr[0], addr[1], addr[2], addr[3], searchResponse.GetHpaiControlEndpoint().GetIpPort()))
						if err != nil {
							continue
						}
						deviceName := string(bytes.Trim(searchResponse.GetDibDeviceInfo().GetDeviceFriendlyName(), "\x00"))
						discoveryEvent := spiModel.NewDefaultPlcDiscoveryItem(
							"knxnet-ip",
							"udp",
							*remoteUrl,
							nil,
							deviceName,
							nil,
						)
						// Pass the event back to the callback
						callback(discoveryEvent)
					}
					continue
				}
			case <-timeout.C:
				{
					timeout.Stop()
					continue
				}
			}
		}
	}
}
