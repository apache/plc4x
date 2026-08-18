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

package eip

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/pool"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// EipUdpDiscoveryDefaultPort is the UDP port EIP/CIP devices listen on for
// ListIdentity broadcast discovery requests (mirrors plc4j's
// Constants.EIPUDPDISCOVERYDEFAULTPORT).
const EipUdpDiscoveryDefaultPort = 44818

type Discoverer struct {
	transportInstanceCreationWorkItemId atomic.Int32
	transportInstanceCreationQueue      pool.Executor
	deviceScanningWorkItemId            atomic.Int32
	deviceScanningQueue                 pool.Executor

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewDiscoverer(_options ...options.WithOption) *Discoverer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Discoverer{
		// TODO: maybe a dynamic executor would be better to not waste cycles when not in use
		transportInstanceCreationQueue: pool.NewFixedSizeExecutor(50, 100, _options...),
		deviceScanningQueue:            pool.NewFixedSizeExecutor(50, 100, _options...),
		log:                            customLogger,
		_options:                       _options,
	}
}

func buildListIdentityRequest() readWriteModel.EipListIdentityRequest {
	return readWriteModel.NewEipListIdentityRequest(0, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0)
}

func (d *Discoverer) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	d.transportInstanceCreationQueue.Start()
	d.deviceScanningQueue.Start()

	udpTransport := udp.NewTransport()

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
			if err := ctx.Err(); err != nil {
				return err
			}
			for _, deviceNameOption := range deviceNames {
				if err := ctx.Err(); err != nil {
					return err
				}
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
		if err := ctx.Err(); err != nil {
			return err
		}
		addrs, err := netInterface.Addrs()
		if err != nil {
			return err
		}
		wg.Go(func() {
			defer func() {
				if err := recover(); err != nil {
					d.log.Error().
						Str("stack", string(debug.Stack())).
						Interface("err", err).
						Msg("panic-ed")
				}
			}()
			// Iterate over all addresses the current interface has configured.
			// We're only interested in IPv4 subnets, as we need to compute a
			// broadcast address to send the ListIdentity request to.
			for _, addr := range addrs {
				if err := ctx.Err(); err != nil {
					d.log.Debug().Err(err).Msg("done")
					return
				}
				ipNet, ok := addr.(*net.IPNet)
				if !ok {
					continue
				}
				ipv4Addr := ipNet.IP.To4()
				if ipv4Addr == nil || ipv4Addr.IsLoopback() {
					continue
				}
				mask := ipNet.Mask
				if len(mask) != len(ipv4Addr) {
					// Normalize a possibly 16-byte mask down to 4 bytes for an IPv4 address.
					if ones, bits := ipNet.Mask.Size(); bits == 32 || bits == 128 {
						mask = net.CIDRMask(ones, 32)
					} else {
						continue
					}
				}
				broadcastAddr := make(net.IP, len(ipv4Addr))
				for i := range ipv4Addr {
					broadcastAddr[i] = ipv4Addr[i] | ^mask[i]
				}

				connectionUrl, err := url.Parse(fmt.Sprintf("udp://%s:%d", broadcastAddr.String(), EipUdpDiscoveryDefaultPort))
				if err != nil {
					d.log.Error().Err(err).Msg("error building broadcast connection url")
					continue
				}

				d.transportInstanceCreationQueue.Submit(ctx, d.transportInstanceCreationWorkItemId.Add(1), d.createTransportInstanceDispatcher(ctx, wg, connectionUrl, ipv4Addr, udpTransport, transportInstances))
			}
		})
	}
	d.wg.Go(func() {
		wg.Wait()
		d.log.Trace().Msg("Closing transport instance channel")
		close(transportInstances)
	})

	d.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				d.log.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
			}
		}()
		for transportInstance := range transportInstances {
			if transportInstance == nil {
				d.log.Trace().Msg("channel closed")
				break
			}
			d.deviceScanningQueue.Submit(ctx, d.deviceScanningWorkItemId.Add(1), d.createDeviceScanDispatcher(ctx, transportInstance.(*udp.TransportInstance), callback))
		}
	})
	return nil
}

func (d *Discoverer) createTransportInstanceDispatcher(ctx context.Context, wg *sync.WaitGroup, connectionUrl *url.URL, ipv4Addr net.IP, udpTransport *udp.Transport, transportInstances chan transports.TransportInstance) pool.Runnable {
	wg.Add(1)
	return func(workerCtx context.Context) {
		ctx, cancel := context.WithCancel(ctx)
		context.AfterFunc(workerCtx, cancel)
		defer wg.Done()
		// Create a new "connection" (Actually open a local udp socket and target outgoing packets to that address)
		transportInstance, err :=
			udpTransport.CreateTransportInstanceForLocalAddress(*connectionUrl, nil,
				&net.UDPAddr{IP: ipv4Addr, Port: 0})
		if err != nil {
			d.log.Error().Err(err).Msg("error creating transport instance")
			return
		}
		err = transportInstance.Connect(ctx)
		if err != nil {
			d.log.Debug().Err(err).Msg("Error Connecting")
			return
		}
		d.log.Debug().Interface("transportInstance", transportInstance).Msg("Adding transport instance to scan %v")
		transportInstances <- transportInstance
	}
}

func (d *Discoverer) createDeviceScanDispatcher(ctx context.Context, udpTransportInstance *udp.TransportInstance, callback func(event apiModel.PlcDiscoveryItem)) pool.Runnable {
	return func(workerCtx context.Context) {
		ctx, cancel := context.WithCancel(ctx)
		context.AfterFunc(workerCtx, cancel)
		d.log.Debug().Interface("udpTransportInstance", udpTransportInstance).Msg("Scanning")
		// Create a codec for sending and receiving messages.
		codec := NewMessageCodec(
			udpTransportInstance,
			binary.LittleEndian,
			append(d._options, options.WithCustomLogger(d.log))...,
		)
		// Explicitly start the worker
		if err := codec.Connect(ctx); err != nil {
			d.log.Error().Err(err).Msg("Error connecting")
			return
		}

		// Send the ListIdentity broadcast request.
		listIdentityRequest := buildListIdentityRequest()
		if err := codec.Send(ctx, "device_scan_list_identity_request", listIdentityRequest); err != nil {
			d.log.Debug().Err(err).Interface("listIdentityRequest", listIdentityRequest).Msg("Error sending message")
			return
		}
		// Keep on reading responses till the timeout is done.
		// TODO: Make this configurable
		timeout := time.NewTimer(1 * time.Second)
		timeout.Stop()
		for start := time.Now(); time.Since(start) < time.Second*5; {
			if err := ctx.Err(); err != nil {
				d.log.Debug().Err(err).Msg("done")
				return
			}
			timeout.Reset(1 * time.Second)
			select {
			case message := <-codec.GetDefaultIncomingMessageChannel():
				{
					if !timeout.Stop() {
						<-timeout.C
					}
					listIdentityResponse, ok := message.(readWriteModel.EipListIdentityResponse)
					if !ok {
						continue
					}
					for _, item := range listIdentityResponse.GetItems() {
						identity, ok := item.(readWriteModel.CipIdentity)
						if !ok {
							continue
						}
						addr := identity.GetSocketAddressAddress()
						if len(addr) != 4 {
							continue
						}
						remoteUrl, err := url.Parse(fmt.Sprintf("eip://%d.%d.%d.%d:%d",
							addr[0], addr[1], addr[2], addr[3], EipUdpDiscoveryDefaultPort))
						if err != nil {
							continue
						}
						productName := identity.GetProductName()
						discoveryEvent := spiModel.NewDefaultPlcDiscoveryItem(
							"eip",
							"udp",
							*remoteUrl,
							nil,
							productName,
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

func (d *Discoverer) Close() error {
	defer utils.StopWarn(d.log)()
	d.log.Trace().Msg("Closing discoverer")
	var collectedErrors []error
	d.log.Trace().Msg("Closing transport instance creation queue")
	if err := d.transportInstanceCreationQueue.Close(); err != nil {
		collectedErrors = append(collectedErrors, errors.Wrap(err, "error closing transport instance creation queue"))
	}
	d.log.Trace().Msg("Closing device scanning queue")
	if err := d.deviceScanningQueue.Close(); err != nil {
		collectedErrors = append(collectedErrors, errors.Wrap(err, "error closing device scanning queue"))
	}
	d.log.Trace().Msg("waiting for wait group")
	d.wg.Wait()
	if err := errors.Join(collectedErrors...); err != nil {
		return errors.Wrap(err, "error closing discoverer")
	}
	return nil
}
