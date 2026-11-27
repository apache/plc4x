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
	"net"
	"net/netip"
	"runtime/debug"
	"strconv"
	"time"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

// defaultForeignDeviceTTL is the registration TTL in seconds used when ForeignDeviceTTL is 0.
const defaultForeignDeviceTTL uint16 = 300

// defaultForeignDeviceResultTimeout bounds the wait for the BBMD's BVLC-Result when ApduTimeoutMs is 0.
const defaultForeignDeviceResultTimeout = 3 * time.Second

// foreignDeviceTTL is the TTL this connection registers with.
func (c *Connection) foreignDeviceTTL() uint16 {
	if c.configuration.ForeignDeviceTTL == 0 {
		return defaultForeignDeviceTTL
	}
	return c.configuration.ForeignDeviceTTL
}

// foreignDeviceRenewalInterval is how often a registration with the given TTL is renewed: at
// half the TTL, so one lost renewal still leaves the entry alive.
func foreignDeviceRenewalInterval(ttl uint16) time.Duration {
	return time.Duration(ttl) * time.Second / 2
}

// splitForeignDeviceBBMD splits ForeignDeviceBBMD ("host" or "host:port") into host and port,
// defaulting the port to the BACnet/IP port 47808.
func splitForeignDeviceBBMD(bbmd string) (host string, port uint16, err error) {
	host, portString, splitErr := net.SplitHostPort(bbmd)
	if splitErr != nil {
		// No port present: the whole string is the host.
		host, portString = bbmd, strconv.Itoa(int(readWriteModel.Constant_BACNETUDPDEFAULTPORT))
	}
	if host == "" {
		return "", 0, errors.Errorf("no host in ForeignDeviceBBMD %q", bbmd)
	}
	parsedPort, err := strconv.ParseUint(portString, 10, 16)
	if err != nil || parsedPort == 0 {
		return "", 0, errors.Errorf("invalid port in ForeignDeviceBBMD %q", bbmd)
	}
	return host, uint16(parsedPort), nil
}

// resolveForeignDeviceBBMD resolves the BBMD's address. It runs on every registration, so a BBMD
// given by name follows a DNS change and a failed lookup is retried with the next renewal.
func resolveForeignDeviceBBMD(ctx context.Context, host string, port uint16) (*net.UDPAddr, error) {
	if ip, err := netip.ParseAddr(host); err == nil {
		return net.UDPAddrFromAddrPort(netip.AddrPortFrom(ip, port)), nil
	}
	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip4", host)
	if err != nil {
		return nil, errors.Wrapf(err, "error resolving BBMD %q", host)
	}
	if len(ips) == 0 {
		return nil, errors.Errorf("BBMD %q resolves to no IPv4 address", host)
	}
	return net.UDPAddrFromAddrPort(netip.AddrPortFrom(ips[0], port)), nil
}

// startForeignDeviceRegistration registers this connection as a foreign device with the BBMD in
// ForeignDeviceBBMD and keeps the registration alive until ctx ends. It does nothing without a
// BBMD. The loop runs on the connection's wait group, so Close waits for it.
func (c *Connection) startForeignDeviceRegistration(ctx context.Context) {
	bbmd := c.configuration.ForeignDeviceBBMD
	if bbmd == "" {
		return
	}
	log := c.log.With().Str("bbmd", bbmd).Logger()
	host, port, err := splitForeignDeviceBBMD(bbmd)
	if err != nil {
		log.Error().Err(err).Msg("not registering as a foreign device")
		return
	}
	writer, ok := c.foreignDeviceWriter()
	if !ok {
		log.Error().Msg("not registering as a foreign device: the transport cannot send to the BBMD from the connection's socket")
		return
	}
	done := make(chan struct{})
	c.foreignDeviceDone = done
	c.wg.Go(func() {
		defer close(done)
		c.maintainForeignDeviceRegistration(ctx, writer, host, port)
	})
}

// foreignDeviceWriter is the connection's transport instance, if it can send to the BBMD.
func (c *Connection) foreignDeviceWriter() (transports.AddressedDatagramWriter, bool) {
	exposer, ok := c.messageCodec.(spi.TransportInstanceExposer)
	if !ok {
		return nil, false
	}
	writer, ok := exposer.GetTransportInstance().(transports.AddressedDatagramWriter)
	return writer, ok
}

// maintainForeignDeviceRegistration registers right away and renews at TTL/2 until ctx ends. A
// failed registration is logged and retried on the next tick.
func (c *Connection) maintainForeignDeviceRegistration(ctx context.Context, writer transports.AddressedDatagramWriter, host string, port uint16) {
	defer func() {
		if err := recover(); err != nil {
			c.log.Error().
				Str("stack", string(debug.Stack())).
				Interface("err", err).
				Msg("panic-ed")
		}
	}()
	ttl := c.foreignDeviceTTL()
	interval := foreignDeviceRenewalInterval(ttl)
	log := c.log.With().Str("bbmd", c.configuration.ForeignDeviceBBMD).Uint16("ttl", ttl).Dur("renewalInterval", interval).Logger()
	log.Info().Msg("Registering as a foreign device")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := c.registerForeignDevice(ctx, writer, host, port, ttl); err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Warn().Err(err).Msg("foreign-device registration failed; retrying with the next renewal")
		} else {
			log.Debug().Msg("registered as a foreign device")
		}
		select {
		case <-ctx.Done():
			log.Debug().Msg("Ending foreign-device registration")
			return
		case <-ticker.C:
		}
	}
}

// registerForeignDevice sends one BVLC Register-Foreign-Device to the BBMD and waits for its
// BVLC-Result. It fails on a NAK, on no answer within the APDU timeout and when sending fails.
func (c *Connection) registerForeignDevice(ctx context.Context, writer transports.AddressedDatagramWriter, host string, port uint16, ttl uint16) error {
	timeout := time.Duration(c.configuration.ApduTimeoutMs) * time.Millisecond
	if timeout == 0 {
		timeout = defaultForeignDeviceResultTimeout
	}
	attemptCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	address, err := resolveForeignDeviceBBMD(attemptCtx, host, port)
	if err != nil {
		return err
	}
	request, err := readWriteModel.NewBVLCRegisterForeignDevice(ttl).Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing Register-Foreign-Device")
	}

	// Buffered and never blocking: the codec calls these from its own workers.
	result := make(chan error, 1)
	report := func(err error) {
		select {
		case result <- err:
		default:
		}
	}
	// A BVLC-Result carries no invoke id, so the first one to arrive answers the registration.
	// Nothing else this connection sends asks for one.
	c.messageCodec.Expect(attemptCtx, "Register-Foreign-Device",
		func(message spi.Message) bool {
			_, ok := message.(readWriteModel.BVLCResult)
			return ok
		},
		func(message spi.Message) error {
			if code := message.(readWriteModel.BVLCResult).GetCode(); code != readWriteModel.BVLCResultCode_SUCCESSFUL_COMPLETION {
				report(errors.Errorf("BBMD %s answered %s", address, code))
			} else {
				report(nil)
			}
			return nil
		},
		func(err error) error {
			report(err)
			return nil
		},
	)
	if err := writer.WriteToAddress(attemptCtx, request, address); err != nil {
		return errors.Wrapf(err, "error sending Register-Foreign-Device to %s", address)
	}
	select {
	case err := <-result:
		return err
	case <-attemptCtx.Done():
		return errors.Wrapf(context.Cause(attemptCtx), "no BVLC-Result from BBMD %s", address)
	}
}
