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
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/profinet/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

const (
	// etherTypeProfinet carries DCP and real-time frames.
	etherTypeProfinet = 0x8892
	// identifyResponseDelay spreads the answers over up to 2.56 s (factor x 10 ms), so a segment full
	// of devices doesn't answer in one burst. plc4j sends the same value.
	identifyResponseDelay = 256
	// defaultDiscoveryWindow is how long a scan listens: past the response delay spread.
	defaultDiscoveryWindow = 4 * time.Second
	// DiscoveryWindowOption is the protocol specific discovery option that changes the window, as a
	// time.Duration or as milliseconds.
	DiscoveryWindowOption = "discovery-window"
)

// dcpMulticast is the Profinet DCP identify multicast address.
var dcpMulticast = net.HardwareAddr{0x01, 0x0E, 0xCF, 0x00, 0x00, 0x00}

// station is a transport instance that knows its own MAC address, which a DCP request carries as its
// source.
type station interface {
	transports.TransportInstance
	HardwareAddress() net.HardwareAddr
}

// Discoverer finds Profinet devices with a DCP Identify All request on every selected interface,
// the port of plc4j's ProfinetDiscoverer.
type Discoverer struct {
	// openStation opens a station on an interface, filtered to Profinet frames.
	openStation func(interfaceName string) (station, error)
	// interfaces lists the interfaces to scan when no device name option selects some.
	interfaces func() ([]string, error)
	window     time.Duration

	wg sync.WaitGroup

	log zerolog.Logger
}

func NewDiscoverer(_options ...options.WithOption) *Discoverer {
	log := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Discoverer{
		openStation: func(interfaceName string) (station, error) {
			return rawsocket.NewTransportInstance(interfaceName, rawsocket.Filter{EtherTypes: []uint16{etherTypeProfinet}}, rawsocket.NewTransport(), _options...), nil
		},
		interfaces: ethernetInterfaces,
		window:     defaultDiscoveryWindow,
		log:        log,
	}
}

// ethernetInterfaces are the interfaces that are up, aren't loopback and have a MAC address.
func ethernetInterfaces() ([]string, error) {
	all, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, networkInterface := range all {
		if networkInterface.Flags&net.FlagUp == 0 || networkInterface.Flags&net.FlagLoopback != 0 || len(networkInterface.HardwareAddr) != 6 {
			continue
		}
		names = append(names, networkInterface.Name)
	}
	return names, nil
}

// Discover starts a scan of each selected interface in the background and returns; items reach the
// callback as devices answer, until the window closes or ctx ends. Interfaces are selected with
// device name options, otherwise every Ethernet interface is scanned.
func (d *Discoverer) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	interfaceNames, window, err := d.plan(discoveryOptions)
	if err != nil {
		return err
	}
	for _, interfaceName := range interfaceNames {
		d.wg.Go(func() {
			defer func() {
				if err := recover(); err != nil {
					d.log.Error().Str("stack", string(debug.Stack())).Interface("err", err).Msg("panic-ed")
				}
			}()
			if err := d.scan(ctx, interfaceName, window, callback); err != nil {
				d.log.Debug().Err(err).Str("interface", interfaceName).Msg("DCP scan failed")
			}
		})
	}
	return nil
}

func (d *Discoverer) plan(discoveryOptions []options.WithDiscoveryOption) ([]string, time.Duration, error) {
	window := d.window
	for _, option := range options.FilterDiscoveryOptionProtocolSpecific(discoveryOptions) {
		if option.GetKey() != DiscoveryWindowOption {
			continue
		}
		switch value := option.GetValue().(type) {
		case time.Duration:
			window = value
		case int:
			window = time.Duration(value) * time.Millisecond
		default:
			return nil, 0, errors.Errorf("%s is a time.Duration or milliseconds, got %T", DiscoveryWindowOption, value)
		}
	}
	var interfaceNames []string
	for _, option := range options.FilterDiscoveryOptionsDeviceName(discoveryOptions) {
		interfaceNames = append(interfaceNames, option.GetDeviceName())
	}
	if len(interfaceNames) == 0 {
		var err error
		if interfaceNames, err = d.interfaces(); err != nil {
			return nil, 0, errors.Wrap(err, "error listing the network interfaces")
		}
	}
	return interfaceNames, window, nil
}

// scan sends one Identify All on an interface and reports each device that answers, once, until
// the window closes.
func (d *Discoverer) scan(ctx context.Context, interfaceName string, window time.Duration, callback func(event apiModel.PlcDiscoveryItem)) error {
	instance, err := d.openStation(interfaceName)
	if err != nil {
		return err
	}
	if err := instance.Connect(ctx); err != nil {
		return err
	}
	defer func() { _ = instance.Close() }()
	return identifyOn(ctx, instance, window, d.log, func(device identifiedDevice) bool {
		callback(device.item(interfaceName, instance.HardwareAddress().String()))
		return false
	})
}

// identifyOn sends one Identify All on an open station and hands each device that answers, once,
// to found, until the window closes or found returns true.
func identifyOn(ctx context.Context, instance station, window time.Duration, log zerolog.Logger, found func(identifiedDevice) bool) error {
	local := instance.HardwareAddress()
	xid := randomXid()
	request, err := identifyRequest(local, xid)
	if err != nil {
		return err
	}
	if err := instance.Write(ctx, request); err != nil {
		return errors.Wrap(err, "error sending the DCP identify request")
	}
	scanCtx, cancel := context.WithTimeout(ctx, window)
	defer cancel()
	seen := map[string]bool{}
	for {
		frame, err := readFrame(scanCtx, instance)
		if err != nil {
			// The socket's read deadline is the window's end too, and it can report that a moment
			// before the context does.
			if scanCtx.Err() != nil || errors.Is(err, os.ErrDeadlineExceeded) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}
		device, ok := identifiedFrom(scanCtx, local, xid, frame, log)
		if !ok || seen[device.mac] {
			continue
		}
		seen[device.mac] = true
		if found(device) {
			return nil
		}
	}
}

func randomXid() uint32 {
	var b [4]byte
	_, _ = rand.Read(b[:])
	return binary.BigEndian.Uint32(b[:])
}

// identifyRequest is a DCP Identify All, VLAN tagged with priority 0 the way plc4j sends it, as a
// raw-socket record.
func identifyRequest(source net.HardwareAddr, xid uint32) ([]byte, error) {
	frame := readWriteModel.NewEthernet_Frame(
		readWriteModel.NewMacAddress(dcpMulticast),
		readWriteModel.NewMacAddress(source),
		readWriteModel.NewEthernet_FramePayload_VirtualLan(readWriteModel.VirtualLanPriority_BEST_EFFORT, false, 0,
			readWriteModel.NewEthernet_FramePayload_PnDcp(
				readWriteModel.NewPnDcp_Pdu_IdentifyReq(uint16(readWriteModel.PnDcp_FrameId_DCP_Identify_ReqPDU), xid, identifyResponseDelay,
					[]readWriteModel.PnDcp_Block{readWriteModel.NewPnDcp_Block_ALLSelector()}))))
	data, err := frame.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "error serializing the DCP identify request")
	}
	if len(data) < 60 {
		// Pad to the Ethernet minimum ourselves: a virtual or software interface won't.
		data = append(data, make([]byte, 60-len(data))...)
	}
	return rawsocket.EncodeRecord(data)
}

// readFrame reads the next frame record.
func readFrame(ctx context.Context, instance transports.TransportInstance) ([]byte, error) {
	header, err := instance.PeekReadableBytes(ctx, rawsocket.RecordHeaderSize)
	if err != nil {
		return nil, err
	}
	length, err := rawsocket.RecordFrameLength(header)
	if err != nil {
		return nil, err
	}
	record, err := instance.Read(ctx, uint32(rawsocket.RecordHeaderSize+length))
	if err != nil {
		return nil, err
	}
	return record[rawsocket.RecordHeaderSize:], nil
}

// identifiedFrom reads an Identify response to our request. Frames that are anything else,
// including answers to somebody else's request, are skipped.
func identifiedFrom(ctx context.Context, local net.HardwareAddr, xid uint32, data []byte, log zerolog.Logger) (identifiedDevice, bool) {
	frame, err := readWriteModel.Ethernet_FrameParse(ctx, data)
	if err != nil {
		log.Trace().Err(err).Msg("not a frame the Profinet model parses")
		return identifiedDevice{}, false
	}
	if net.HardwareAddr(frame.GetDestination().GetAddress()).String() != local.String() {
		// Identify responses go unicast to the requester.
		return identifiedDevice{}, false
	}
	payload := frame.GetPayload()
	if vlan, ok := payload.(readWriteModel.Ethernet_FramePayload_VirtualLan); ok {
		payload = vlan.GetPayload()
	}
	dcp, ok := payload.(readWriteModel.Ethernet_FramePayload_PnDcp)
	if !ok {
		return identifiedDevice{}, false
	}
	response, ok := dcp.GetPdu().(readWriteModel.PnDcp_Pdu_IdentifyRes)
	if !ok || response.GetXid() != xid {
		return identifiedDevice{}, false
	}
	device := deviceFrom(response)
	device.mac = net.HardwareAddr(frame.GetSource().GetAddress()).String()
	return device, true
}

// identifiedDevice is what an Identify response says about a device.
type identifiedDevice struct {
	typeName, name, role string
	vendorID, deviceID   string
	ip, subnetMask       string
	mac                  string
	// vendor and product are the numeric ids, zero when the device didn't send them.
	vendor, product uint16
}

func deviceFrom(response readWriteModel.PnDcp_Pdu_IdentifyRes) identifiedDevice {
	device := identifiedDevice{typeName: "unknown", name: "unknown", role: "unknown", vendorID: "unknown", deviceID: "unknown", ip: "unknown", subnetMask: "unknown"}
	for _, block := range response.GetBlocks() {
		switch b := block.(type) {
		case readWriteModel.PnDcp_Block_DevicePropertiesDeviceVendor:
			device.typeName = string(b.GetDeviceVendorValue())
		case readWriteModel.PnDcp_Block_DevicePropertiesNameOfStation:
			device.name = string(b.GetNameOfStation())
		case readWriteModel.PnDcp_Block_DevicePropertiesDeviceId:
			device.vendorID = fmt.Sprintf("%04X", b.GetVendorId())
			device.deviceID = fmt.Sprintf("%04X", b.GetDeviceId())
			device.vendor, device.product = b.GetVendorId(), b.GetDeviceId()
		case readWriteModel.PnDcp_Block_DevicePropertiesDeviceRole:
			var roles []string
			for _, r := range []struct {
				set  bool
				name string
			}{{b.GetPnioSupervisor(), "SUPERVISOR"}, {b.GetPnioMultidevive(), "MULTIDEVICE"}, {b.GetPnioController(), "CONTROLLER"}, {b.GetPnioDevice(), "DEVICE"}} {
				if r.set {
					roles = append(roles, r.name)
				}
			}
			if len(roles) > 0 {
				device.role = strings.Join(roles, ",")
			}
		case readWriteModel.PnDcp_Block_IpParameter:
			device.ip = net.IP(b.GetIpAddress()).String()
			device.subnetMask = net.IP(b.GetSubnetMask()).String()
		}
	}
	return device
}

// item is the discovery item. The connection URL names the interface the device was found on; the
// device itself goes in options: its name of station, which is what Profinet addresses a device by,
// its MAC address and, once it has one, its IP address.
func (d identifiedDevice) item(interfaceName string, localMAC string) apiModel.PlcDiscoveryItem {
	connectionOptions := map[string][]string{
		"device-name": {d.name},
		"device-mac":  {d.mac},
	}
	if d.ip != "unknown" && d.ip != "0.0.0.0" {
		connectionOptions["device-ip"] = []string{d.ip}
	}
	attributes := map[string]apiValues.PlcValue{}
	for key, value := range map[string]string{
		"ipAddress":       d.ip,
		"subnetMask":      d.subnetMask,
		"macAddress":      d.mac,
		"localMacAddress": localMAC,
		"deviceTypeName":  d.typeName,
		"deviceName":      d.name,
		"vendorId":        d.vendorID,
		"deviceId":        d.deviceID,
		"role":            d.role,
		"packetType":      "dcp",
	} {
		attributes[key] = spiValues.NewPlcSTRING(value)
	}
	return spiModel.NewDefaultPlcDiscoveryItem(
		DriverCode,
		rawsocket.TransportCode,
		url.URL{Scheme: rawsocket.TransportCode, Host: interfaceName},
		connectionOptions,
		d.typeName+" - "+d.name,
		attributes,
	)
}

func (d *Discoverer) Close() error {
	d.wg.Wait()
	return nil
}
