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

package profinet

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/profinet/gsdml"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// Configuration is what a profinet connection string can say. gsd-directory, dap-id and ip-address
// are plc4j's; the device options are what discovery puts in its items, and the timing options
// pick the cycle (send clock x reduction ratio x 31.25 µs; 1 ms x 16 = 16 ms by default).
type Configuration struct {
	gsdDirectory   string
	deviceName     string
	deviceMAC      net.HardwareAddr
	deviceIP       net.IP
	localIP        net.IP
	settings       arSettings
	identifyWindow time.Duration
	// modules are GSDML module ids by slot, for slots the device doesn't report on its own: a
	// device that plugs what the controller expects (p-net's sample device does) or a module to
	// replace what is reported.
	modules map[uint16]string
}

func ParseFromOptions(log zerolog.Logger, connectionOptions map[string][]string) (Configuration, error) {
	reader := options.NewOptionReader(log, connectionOptions)
	defer reader.ReportUnknown(DriverCode)
	reader.Ignore("dap-id") // the DAP is what the device reports in slot 0
	c := Configuration{
		gsdDirectory: "~/.gsd",
		settings: arSettings{
			stationName: "plc4x", sendClockFactor: 32, reductionRatio: 16, watchdogFactor: 3, dataHoldFactor: 3,
			requestTimeout: 5 * time.Second, readyTimeout: 30 * time.Second,
		},
		identifyWindow: 3 * time.Second,
	}
	if v := reader.Get("gsd-directory"); v != "" {
		c.gsdDirectory = v
	}
	c.deviceName = reader.Get("device-name")
	if v := reader.Get("device-mac"); v != "" {
		mac, err := net.ParseMAC(v)
		if err != nil || len(mac) != 6 {
			return Configuration{}, errors.Errorf("bad device-mac %q", v)
		}
		c.deviceMAC = mac
	}
	for _, key := range []string{"device-ip", "ip-address"} {
		if v := reader.Get(key); v != "" {
			if c.deviceIP = net.ParseIP(v).To4(); c.deviceIP == nil {
				return Configuration{}, errors.Errorf("bad %s %q", key, v)
			}
		}
	}
	if v := reader.Get("local-ip"); v != "" {
		if c.localIP = net.ParseIP(v).To4(); c.localIP == nil {
			return Configuration{}, errors.Errorf("bad local-ip %q", v)
		}
	}
	if v := reader.Get("station-name"); v != "" {
		c.settings.stationName = v
	}
	for key, target := range map[string]*uint16{
		"send-clock-factor": &c.settings.sendClockFactor, "reduction-ratio": &c.settings.reductionRatio,
		"watchdog-factor": &c.settings.watchdogFactor, "data-hold-factor": &c.settings.dataHoldFactor,
	} {
		if v := reader.Get(key); v != "" {
			n, err := strconv.ParseUint(v, 10, 16)
			if err != nil || n == 0 {
				return Configuration{}, errors.Errorf("bad %s %q", key, v)
			}
			*target = uint16(n)
		}
	}
	for key, target := range map[string]*time.Duration{
		"request-timeout-ms": &c.settings.requestTimeout, "ready-timeout-ms": &c.settings.readyTimeout, "identify-timeout-ms": &c.identifyWindow,
	} {
		if v := reader.Get(key); v != "" {
			n, err := strconv.ParseUint(v, 10, 32)
			if err != nil || n == 0 {
				return Configuration{}, errors.Errorf("bad %s %q", key, v)
			}
			*target = time.Duration(n) * time.Millisecond
		}
	}
	// rpc-fragment-size is the most RPC body (bytes) per datagram; longer requests go out as
	// fragments. The default, 1360, fills the 1440 byte UDP payload devices take; a link with a
	// smaller MTU needs less.
	if v := reader.Get("rpc-fragment-size"); v != "" {
		n, err := strconv.ParseUint(v, 10, 16)
		if err != nil || n < 64 || n > defaultRPCFragmentBody {
			return Configuration{}, errors.Errorf("bad rpc-fragment-size %q: 64 to %d bytes", v, defaultRPCFragmentBody)
		}
		c.settings.rpcFragmentSize = int(n)
	}
	for key, values := range connectionOptions {
		slotText, ok := strings.CutPrefix(key, "module.")
		if !ok || len(values) == 0 {
			continue
		}
		reader.Ignore(key)
		slot, err := strconv.ParseUint(slotText, 10, 16)
		if err != nil || slot == 0 {
			return Configuration{}, errors.Errorf("bad %s: modules go in slots 1 and up", key)
		}
		if c.modules == nil {
			c.modules = map[uint16]string{}
		}
		c.modules[uint16(slot)] = values[0]
	}
	if c.deviceName == "" && c.deviceMAC == nil {
		return Configuration{}, errors.New("name the device with device-name or device-mac (discovery puts both in its items)")
	}
	return c, nil
}

// Connection is an IO controller connection to one Profinet device.
type Connection struct {
	_default.DefaultConnection

	configuration Configuration
	codec         *stationCodec
	ar            *applicationRelation
	subscriber    _default.DefaultPollingSubscriber
	alarms        *alarmSubscriber
	closed        chan struct{}
	closeOnce     sync.Once

	wg  sync.WaitGroup
	log zerolog.Logger
	opt []options.WithOption
}

var _ _default.DefaultPollingSubscriberRequirements = (*Connection)(nil)

func newConnection(st station, configuration Configuration, _options ...options.WithOption) *Connection {
	log := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	_options = slices.Clip(_options)
	c := &Connection{configuration: configuration, codec: &stationCodec{station: st}, log: log, opt: _options}
	c.DefaultConnection = _default.NewDefaultConnection(c, append(_options,
		_default.WithPlcTagHandler(NewTagHandler()),
		_default.WithPlcValueHandler(spiValues.NewDefaultValueHandler(_options...)))...)
	c.subscriber = _default.NewDefaultPollingSubscriber(c, _options...)
	c.alarms = newAlarmSubscriber(_options...)
	return c
}

func (c *Connection) GetConnection() plc4go.PlcConnection { return c }
func (c *Connection) GetMessageCodec() spi.MessageCodec   { return c.codec }

// Connect identifies the device, loads its GSDML and establishes the application relation.
func (c *Connection) Connect(ctx context.Context) error {
	if err := c.DefaultConnection.Connect(ctx); err != nil {
		return errors.Wrap(err, "error opening the network interface")
	}
	if err := c.establish(ctx); err != nil {
		_ = c.Close()
		return err
	}
	c.closed = make(chan struct{})
	c.wg.Go(func() {
		select {
		case err := <-c.ar.lost:
			c.log.Warn().Err(err).Msg("Profinet connection lost")
			// Invalidate closes the connection, which waits for this goroutine: not from here.
			go c.Invalidate()
		case <-c.closed:
		}
	})
	return nil
}

func (c *Connection) establish(ctx context.Context) error {
	st := c.codec.station
	cfg := c.configuration
	var found identifiedDevice
	if err := identifyOn(ctx, st, cfg.identifyWindow, c.log, func(d identifiedDevice) bool {
		if (cfg.deviceName != "" && strings.EqualFold(d.name, cfg.deviceName)) || (cfg.deviceMAC != nil && d.mac == cfg.deviceMAC.String()) {
			found = d
			return true
		}
		return false
	}); err != nil {
		return errors.Wrap(err, "error identifying the device")
	}
	if found.mac == "" {
		return errors.Errorf("no Profinet device %s answered on %s", firstOf(cfg.deviceName, cfg.deviceMAC.String()), st)
	}
	if found.vendor == 0 && found.product == 0 {
		return errors.New("the device did not report its vendor and device id")
	}
	mac, _ := net.ParseMAC(found.mac)
	ip := cfg.deviceIP
	if ip == nil {
		ip = net.ParseIP(found.ip).To4()
	}
	if ip == nil || ip.IsUnspecified() {
		return errors.Errorf("device %s has no IP address; set one with device-ip or give it one with an engineering tool first", found.name)
	}
	profile, path, err := gsdml.FindInDirectory(cfg.gsdDirectory, found.vendor, found.product)
	if err != nil {
		return err
	}
	c.log.Debug().Str("gsdml", path).Str("device", found.name).Msg("using GSDML")
	localIP := cfg.localIP
	if localIP == nil {
		localIP = interfaceIPv4(st)
	}
	if localIP == nil {
		return errors.Errorf("%s has no IPv4 address to talk to the device from; set local-ip", st)
	}
	c.ar = newApplicationRelation(cfg.settings, device{name: found.name, mac: mac, ip: ip, vendorID: found.vendor, deviceID: found.product},
		gsdml.NewIndex(profile), st, localIP, c.log)
	c.ar.modules = cfg.modules
	c.ar.onAlarm = c.alarms.publish
	return c.ar.establish(ctx)
}

func firstOf(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// interfaceIPv4 is the station's interface's IPv4 address, for a real raw socket.
func interfaceIPv4(st station) net.IP {
	type localAddresser interface{ LocalAddress() net.Addr }
	if l, ok := st.(localAddresser); ok {
		if ipNet, ok := l.LocalAddress().(*net.IPNet); ok {
			return ipNet.IP.To4()
		}
	}
	return nil
}

func (c *Connection) Close() error {
	c.closeOnce.Do(func() {
		if c.closed != nil {
			close(c.closed)
		}
	})
	if c.ar != nil {
		c.ar.close()
	}
	c.subscriber.Close()
	c.alarms.Close()
	err := c.DefaultConnection.Close()
	c.wg.Wait()
	return err
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{ProvidesReading: true, ProvidesWriting: true, ProvidesSubscribing: true}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(c.GetPlcTagHandler(), readerFunc(c.read))
}

func (c *Connection) PollingReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return c.ReadRequestBuilder()
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(c.GetPlcTagHandler(), c.GetPlcValueHandler(), writerFunc(c.write))
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(c.GetPlcTagHandler(), c.GetPlcValueHandler(), subscriberRouter{polling: c.subscriber, alarms: c.alarms})
}

func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcUnsubscriptionRequestBuilder()
}

func (c *Connection) String() string {
	return fmt.Sprintf("profinet.Connection{%s}", c.codec.station)
}

// read answers from the latest cyclic data: inputs as the device last sent them, outputs as they
// go out. An input whose provider status isn't good reads NOT_READY.
func (c *Connection) read(_ context.Context, request apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	result := make(chan apiModel.PlcReadRequestResult, 1)
	c.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(c.log, result, spiModel.NewDefaultPlcReadRequestResult(request, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		codes := map[string]apiModel.PlcResponseCode{}
		values := map[string]apiValues.PlcValue{}
		input, _, last := c.ar.inputImage()
		output := c.ar.outputImage()
		for _, name := range request.GetTagNames() {
			tag, ok := castToTag(request.GetTag(name))
			if !ok {
				codes[name] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			location, err := c.ar.plan.locate(tag)
			if err != nil {
				c.log.Debug().Err(err).Str("tag", name).Msg("no such data item")
				codes[name] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			image := input
			if tag.Output {
				image = output
			} else {
				s, _ := c.ar.plan.Find(tag.Slot, tag.Subslot)
				if last.IsZero() || input[s.InputIOPS] != iopsGood {
					codes[name] = apiModel.PlcResponseCode_NOT_READY
					continue
				}
			}
			value, err := decodeItem(tag, image[location.offset:location.offset+location.size])
			if err != nil {
				c.log.Debug().Err(err).Str("tag", name).Msg("value doesn't decode as the tag's type")
				codes[name] = apiModel.PlcResponseCode_INVALID_DATATYPE
				continue
			}
			codes[name] = apiModel.PlcResponseCode_OK
			values[name] = value
		}
		utils.DeliverResult(c.log, result, spiModel.NewDefaultPlcReadRequestResult(request, spiModel.NewDefaultPlcReadResponse(request, codes, values), nil))
	})
	return result
}

// write changes output data, which goes out from the next cycle on.
func (c *Connection) write(_ context.Context, request apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	c.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(c.log, result, spiModel.NewDefaultPlcWriteRequestResult(request, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		codes := map[string]apiModel.PlcResponseCode{}
		for _, name := range request.GetTagNames() {
			tag, ok := castToTag(request.GetTag(name))
			if !ok {
				codes[name] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			if !tag.Output {
				codes[name] = apiModel.PlcResponseCode_ACCESS_DENIED
				continue
			}
			location, err := c.ar.plan.locate(tag)
			if err != nil {
				codes[name] = apiModel.PlcResponseCode_INVALID_ADDRESS
				continue
			}
			data, mask, err := encodeItem(tag, location.size, request.GetValue(name))
			if err != nil {
				c.log.Debug().Err(err).Str("tag", name).Msg("value doesn't encode as the tag's type")
				codes[name] = apiModel.PlcResponseCode_INVALID_DATA
				continue
			}
			if err := c.ar.updateOutput(location.offset, data, mask); err != nil {
				codes[name] = apiModel.PlcResponseCode_INTERNAL_ERROR
				continue
			}
			codes[name] = apiModel.PlcResponseCode_OK
		}
		utils.DeliverResult(c.log, result, spiModel.NewDefaultPlcWriteRequestResult(request, spiModel.NewDefaultPlcWriteResponse(request, codes), nil))
	})
	return result
}

type readerFunc func(context.Context, apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult

func (f readerFunc) Read(ctx context.Context, r apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	return f(ctx, r)
}

type writerFunc func(context.Context, apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult

func (f writerFunc) Write(ctx context.Context, r apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	return f(ctx, r)
}

// stationCodec lets DefaultConnection open and close the station; the relation does the protocol.
type stationCodec struct {
	station station
}

var (
	_ spi.MessageCodec             = (*stationCodec)(nil)
	_ spi.TransportInstanceExposer = (*stationCodec)(nil)
)

func (s *stationCodec) GetTransportInstance() transports.TransportInstance { return s.station }

func (s *stationCodec) Connect(ctx context.Context) error {
	if s.station.IsConnected() {
		return nil
	}
	return s.station.Connect(ctx)
}

func (s *stationCodec) Disconnect() error { return nil }
func (s *stationCodec) IsRunning() bool   { return s.station.IsConnected() }
func (s *stationCodec) Send(context.Context, string, spi.Message) error {
	return errors.New("the profinet connection sends through its application relation")
}
func (s *stationCodec) Expect(context.Context, string, spi.AcceptsMessage, spi.HandleMessage, spi.HandleError) {
}
func (s *stationCodec) SendRequest(context.Context, string, spi.Message, spi.AcceptsMessage, spi.HandleMessage, spi.HandleError) error {
	return errors.New("the profinet connection sends through its application relation")
}
func (s *stationCodec) GetDefaultIncomingMessageChannel() chan spi.Message { return nil }
