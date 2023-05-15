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
	"strconv"
	"strings"
	"time"

	"github.com/IBM/netaddr"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/libp2p/go-reuseport"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
)

type Discoverer struct {
	messageCodec spi.MessageCodec
}

func NewDiscoverer() *Discoverer {
	return &Discoverer{}
}

func (d *Discoverer) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	// TODO: handle ctx
	interfaces, err := extractInterfaces(discoveryOptions)
	if err != nil {
		return errors.Wrap(err, "error extracting interfaces")
	}

	specificOptions, err := extractProtocolSpecificOptions(discoveryOptions)
	if err != nil {
		return errors.Wrap(err, "error extracting protocol specific options")
	}

	communicationChannels, err := buildupCommunicationChannels(interfaces, specificOptions.bacNetPort)
	if err != nil {
		return errors.Wrap(err, "error building communication channels")
	}

	// TODO: make adjustable
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*60)
	defer func() {
		cancelFunc()
	}()
	incomingBVLCChannel, err := broadcastAndDiscover(ctx, communicationChannels, specificOptions)
	if err != nil {
		return errors.Wrap(err, "error broadcasting and discovering")
	}
	handleIncomingBVLCs(ctx, callback, incomingBVLCChannel)
	// TODO: make adjustable
	time.Sleep(time.Second * 60)
	for _, channel := range communicationChannels {
		_ = channel.Close()
	}
	return nil
}

func broadcastAndDiscover(ctx context.Context, communicationChannels []communicationChannel, specificOptions *protocolSpecificOptions) (chan receivedBvlcMessage, error) {
	incomingBVLCChannel := make(chan receivedBvlcMessage)
	for _, communicationChannelInstance := range communicationChannels {
		// Prepare the discovery packet data
		{
			var lowLimit driverModel.BACnetContextTagUnsignedInteger
			var highLimit driverModel.BACnetContextTagUnsignedInteger
			if whoIsOptions := specificOptions.whoIsOptions; whoIsOptions != nil && whoIsOptions.limits != nil {
				lowLimit = driverModel.CreateBACnetContextTagUnsignedInteger(0, whoIsOptions.limits.low)
				highLimit = driverModel.CreateBACnetContextTagUnsignedInteger(1, whoIsOptions.limits.high)
			}
			requestWhoIs := driverModel.NewBACnetUnconfirmedServiceRequestWhoIs(lowLimit, highLimit, 0)
			apdu := driverModel.NewAPDUUnconfirmedRequest(requestWhoIs, 0)

			control := driverModel.NewNPDUControl(false, false, false, false, driverModel.NPDUNetworkPriority_NORMAL_MESSAGE)
			npdu := driverModel.NewNPDU(1, control, nil, nil, nil, nil, nil, nil, nil, nil, apdu, 0)
			bvlc := driverModel.NewBVLCOriginalUnicastNPDU(npdu, 0)

			// Send the search request.
			theBytes, err := bvlc.Serialize()
			if err != nil {
				return nil, err
			}
			if _, err := communicationChannelInstance.broadcastConnection.WriteTo(theBytes, communicationChannelInstance.broadcastConnection.LocalAddr()); err != nil {
				log.Debug().Err(err).Msg("Error sending broadcast")
			}
		}
		if whoHasOptions := specificOptions.whoHasOptions; whoHasOptions != nil {
			var lowLimit driverModel.BACnetContextTagUnsignedInteger
			var highLimit driverModel.BACnetContextTagUnsignedInteger
			if limits := whoHasOptions.limits; limits != nil {
				lowLimit = driverModel.CreateBACnetContextTagUnsignedInteger(0, limits.deviceInstanceRangeLow)
				highLimit = driverModel.CreateBACnetContextTagUnsignedInteger(1, limits.deviceInstanceRangeHigh)
			}
			var object driverModel.BACnetUnconfirmedServiceRequestWhoHasObject
			if identifier := whoHasOptions.object.identifier; identifier != nil {
				var objectType uint16
				objectTypeByName, ok := driverModel.BACnetObjectTypeByName(identifier.type_)
				if ok {
					parseUint, err := strconv.ParseUint(identifier.type_, 10, 16)
					if err != nil {
						return nil, err
					}
					objectType = uint16(parseUint)
				} else {
					objectType = uint16(objectTypeByName)
				}
				objectIdentifier := driverModel.CreateBACnetContextTagObjectIdentifier(2, objectType, uint32(identifier.instance))
				object = driverModel.NewBACnetUnconfirmedServiceRequestWhoHasObjectIdentifier(objectIdentifier, objectIdentifier.GetHeader())
			} else if name := whoHasOptions.object.name; name != nil {
				characterString := driverModel.CreateBACnetContextTagCharacterString(3, driverModel.BACnetCharacterEncoding_ISO_10646, *name)
				object = driverModel.NewBACnetUnconfirmedServiceRequestWhoHasObjectName(characterString, characterString.GetHeader())
			} else {
				panic("Invalid state")
			}
			requestWhoHas := driverModel.NewBACnetUnconfirmedServiceRequestWhoHas(lowLimit, highLimit, object, 0)
			apdu := driverModel.NewAPDUUnconfirmedRequest(requestWhoHas, 0)

			control := driverModel.NewNPDUControl(false, false, false, false, driverModel.NPDUNetworkPriority_NORMAL_MESSAGE)
			npdu := driverModel.NewNPDU(1, control, nil, nil, nil, nil, nil, nil, nil, nil, apdu, 0)
			bvlc := driverModel.NewBVLCOriginalUnicastNPDU(npdu, 0)

			// Send the search request.
			theBytes, err := bvlc.Serialize()
			if err != nil {
				return nil, err
			}
			if _, err := communicationChannelInstance.broadcastConnection.WriteTo(theBytes, communicationChannelInstance.broadcastConnection.LocalAddr()); err != nil {
				log.Debug().Err(err).Msg("Error sending broadcast")
			}
		}

		go func(communicationChannelInstance communicationChannel) {
			for {
				blockingReadChan := make(chan bool)
				go func() {
					buf := make([]byte, 4096)
					n, addr, err := communicationChannelInstance.unicastConnection.ReadFrom(buf)
					if err != nil {
						log.Debug().Err(err).Msg("Ending unicast receive")
						blockingReadChan <- false
						return
					}
					log.Debug().Stringer("addr", addr).Msg("Received broadcast bvlc")
					incomingBvlc, err := driverModel.BVLCParse(buf[:n])
					if err != nil {
						log.Warn().Err(err).Msg("Could not parse bvlc")
						blockingReadChan <- true
						return
					}
					incomingBVLCChannel <- receivedBvlcMessage{incomingBvlc, addr}
					blockingReadChan <- true
				}()
				select {
				case ok := <-blockingReadChan:
					if !ok {
						log.Debug().Msg("Ending unicast reading")
						return
					}
					log.Trace().Msg("Received something unicast")
				case <-ctx.Done():
					log.Debug().Err(ctx.Err()).Msg("Ending unicast receive")
					return
				}
			}
		}(communicationChannelInstance)

		go func(communicationChannelInstance communicationChannel) {
			for {
				blockingReadChan := make(chan bool)
				go func() {
					buf := make([]byte, 4096)
					n, addr, err := communicationChannelInstance.broadcastConnection.ReadFrom(buf)
					if err != nil {
						log.Debug().Err(err).Msg("Ending broadcast receive")
						blockingReadChan <- false
						return
					}
					log.Debug().Stringer("addr", addr).Msg("Received broadcast bvlc")
					incomingBvlc, err := driverModel.BVLCParse(buf[:n])
					if err != nil {
						log.Warn().Err(err).Msg("Could not parse bvlc")
						blockingReadChan <- true
					}
					incomingBVLCChannel <- receivedBvlcMessage{incomingBvlc, addr}
					blockingReadChan <- true
				}()
				select {
				case ok := <-blockingReadChan:
					if !ok {
						log.Debug().Msg("Ending broadcast reading")
						return
					}
					log.Trace().Msg("Received something broadcast")
				case <-ctx.Done():
					log.Debug().Err(ctx.Err()).Msg("Ending broadcast receive")
					return
				}
			}
		}(communicationChannelInstance)
	}
	return incomingBVLCChannel, nil
}

func handleIncomingBVLCs(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), incomingBVLCChannel chan receivedBvlcMessage) {
	for {
		select {
		case receivedBvlc := <-incomingBVLCChannel:
			var npdu driverModel.NPDU
			if bvlc, ok := receivedBvlc.bvlc.(interface{ GetNpdu() driverModel.NPDU }); ok {
				npdu = bvlc.GetNpdu()
			}
			_ = npdu
			if apdu := npdu.GetApdu(); apdu == nil {
				nlm := npdu.GetNlm()
				log.Debug().Msgf("Got nlm\n%v", nlm)
				continue
			}
			apdu := npdu.GetApdu()
			if _, ok := apdu.(driverModel.APDUConfirmedRequestExactly); ok {
				log.Debug().Msgf("Got apdu \n%v", apdu)
				continue
			}
			apduUnconfirmedRequest := apdu.(driverModel.APDUUnconfirmedRequestExactly)
			serviceRequest := apduUnconfirmedRequest.GetServiceRequest()
			switch serviceRequest := serviceRequest.(type) {
			case driverModel.BACnetUnconfirmedServiceRequestIAmExactly:
				iAm := serviceRequest
				remoteUrl, err := url.Parse("udp://" + receivedBvlc.addr.String())
				if err != nil {
					log.Debug().Err(err).Msg("Error parsing url")
				}
				discoveryEvent := internalModel.NewDefaultPlcDiscoveryItem(
					"bacnet-ip",
					"udp",
					*remoteUrl,
					nil,
					fmt.Sprintf("device %v:%v", iAm.GetDeviceIdentifier().GetObjectType(), iAm.GetDeviceIdentifier().GetInstanceNumber()),
					nil,
				)

				// Pass the event back to the callback
				callback(discoveryEvent)
			case driverModel.BACnetUnconfirmedServiceRequestIHaveExactly:
				iHave := serviceRequest
				remoteUrl, err := url.Parse("udp://" + receivedBvlc.addr.String())
				if err != nil {
					log.Debug().Err(err).Msg("Error parsing url")
				}
				discoveryEvent := internalModel.NewDefaultPlcDiscoveryItem(
					"bacnet-ip",
					"udp",
					*remoteUrl,
					nil,
					fmt.Sprintf("device %v:%v with %v:%v and %v", iHave.GetDeviceIdentifier().GetObjectType(), iHave.GetDeviceIdentifier().GetInstanceNumber(), iHave.GetObjectIdentifier().GetObjectType(), iHave.GetObjectIdentifier().GetInstanceNumber(), iHave.GetObjectName().GetValue()),
					nil,
				)

				// Pass the event back to the callback
				callback(discoveryEvent)
			}
		case <-ctx.Done():
			log.Debug().Err(ctx.Err()).Msg("Ending unicast receive")
			return
		}
	}
}

func buildupCommunicationChannels(interfaces []net.Interface, bacNetPort int) (communicationChannels []communicationChannel, err error) {
	// Iterate over all network devices of this system.
	for _, networkInterface := range interfaces {
		unicastInterfaceAddress, err := networkInterface.Addrs()
		if err != nil {
			return nil, errors.Wrapf(err, "Error getting Addresses for %v", networkInterface)
		}
		// Iterate over all addresses the current interface has configured
		for _, unicastAddress := range unicastInterfaceAddress {
			var ipAddr net.IP
			switch addr := unicastAddress.(type) {
			// If the device is configured to communicate with a subnet
			case *net.IPNet:
				ipAddr = addr.IP.To4()
				if ipAddr == nil {
					// TODO: for now we only support ipv4 (reuse doesn't like v6 address strings atm)
					continue
					ipAddr = addr.IP.To16()
				}

			// If the device is configured for a point-to-point connection
			case *net.IPAddr:
				ipAddr = addr.IP.To4()
				if ipAddr == nil {
					// TODO: for now we only support ipv4 (reuse doesn't like v6 address strings atm)
					continue
					ipAddr = addr.IP.To16()
				}
			default:
				continue
			}

			if !ipAddr.IsGlobalUnicast() {
				continue
			}

			// Handle undirected
			unicastConnection, err := reuseport.ListenPacket("udp4", fmt.Sprintf("%v:%d", ipAddr, bacNetPort))
			if err != nil {
				log.Debug().Err(err).Msg("Error building unicast Port")
				continue
			}

			_, cidr, _ := net.ParseCIDR(unicastAddress.String())
			broadcastAddr := netaddr.BroadcastAddr(cidr)
			// Handle undirected
			broadcastConnection, err := reuseport.ListenPacket("udp4", fmt.Sprintf("%v:%d", broadcastAddr, bacNetPort))
			if err != nil {
				if err := unicastConnection.Close(); err != nil {
					log.Debug().Err(err).Msg("Error closing transport instance")
				}
				log.Debug().Err(err).Msg("Error building broadcast Port")
				continue
			}
			communicationChannels = append(communicationChannels, communicationChannel{
				networkInterface:    networkInterface,
				unicastConnection:   unicastConnection,
				broadcastConnection: broadcastConnection,
			})
		}
	}
	return
}

type receivedBvlcMessage struct {
	bvlc driverModel.BVLC
	addr net.Addr
}

type communicationChannel struct {
	networkInterface    net.Interface
	unicastConnection   net.PacketConn
	broadcastConnection net.PacketConn
}

func (c communicationChannel) Close() error {
	_ = c.unicastConnection.Close()
	_ = c.broadcastConnection.Close()
	return nil
}

func extractInterfaces(discoveryOptions []options.WithDiscoveryOption) ([]net.Interface, error) {
	allInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
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
	return interfaces, nil
}

type protocolSpecificOptions struct {
	bacNetPort   int
	whoIsOptions *struct {
		limits *struct {
			low  uint
			high uint
		}
	}
	whoHasOptions *struct {
		limits *struct {
			deviceInstanceRangeLow  uint
			deviceInstanceRangeHigh uint
		}
		object struct {
			identifier *struct {
				type_    string
				instance uint
			}
			name *string
		}
	}
}

func bacNetPort(port int) option {
	return func(specificOptions *protocolSpecificOptions) error {
		specificOptions.bacNetPort = port
		return nil
	}
}

func whoIsLimits(whoIsLowLimit, whoIsHighLimit uint) option {
	return func(specificOptions *protocolSpecificOptions) error {
		specificOptions.whoIsOptions = &struct {
			limits *struct {
				low  uint
				high uint
			}
		}{&struct {
			low  uint
			high uint
		}{whoIsLowLimit, whoIsHighLimit}}
		return nil
	}
}

func whoHasOption() option {
	return func(specificOptions *protocolSpecificOptions) error {
		specificOptions.whoHasOptions = &struct {
			limits *struct {
				deviceInstanceRangeLow  uint
				deviceInstanceRangeHigh uint
			}
			object struct {
				identifier *struct {
					type_    string
					instance uint
				}
				name *string
			}
		}{}
		return nil
	}
}

func whoHasLimits(whoHasDeviceInstanceRangeLowLimit, whoHasDeviceInstanceRangeHighLimit uint) option {
	return func(specificOptions *protocolSpecificOptions) error {
		if specificOptions.whoHasOptions == nil {
			panic("we should have set this before")
		}
		specificOptions.whoHasOptions.limits = &struct {
			deviceInstanceRangeLow  uint
			deviceInstanceRangeHigh uint
		}{whoHasDeviceInstanceRangeLowLimit, whoHasDeviceInstanceRangeHighLimit}
		return nil
	}
}

func whoHasObjectIdentifier(objectIdentifierType string, objectIdentifierInstance uint) option {
	return func(specificOptions *protocolSpecificOptions) error {
		if specificOptions.whoHasOptions == nil {
			panic("we should have set this before")
		}
		specificOptions.whoHasOptions.object.identifier = &struct {
			type_    string
			instance uint
		}{objectIdentifierType, objectIdentifierInstance}
		return nil
	}
}

func whoHasObjectName(objectName string) option {
	return func(specificOptions *protocolSpecificOptions) error {
		if specificOptions.whoHasOptions == nil {
			panic("we should have set this before")
		}
		specificOptions.whoHasOptions.object.name = &objectName
		return nil
	}
}

func NewProtocolSpecificOptions(options ...option) (*protocolSpecificOptions, error) {
	var specificOptions protocolSpecificOptions
	for _, _option := range options {
		if parseErr := _option(&specificOptions); parseErr != nil {
			return nil, parseErr
		}
	}
	return &specificOptions, nil
}

type option func(specificOptions *protocolSpecificOptions) error

func extractProtocolSpecificOptions(discoveryOptions []options.WithDiscoveryOption) (*protocolSpecificOptions, error) {
	var collectedOptions []option
	filteredOptionMap := make(map[string][]any)
	for _, protocolSpecificOption := range options.FilterDiscoveryOptionProtocolSpecific(discoveryOptions) {
		key := protocolSpecificOption.GetKey()
		value := protocolSpecificOption.GetValue()
		if _, ok := filteredOptionMap[key]; !ok {
			filteredOptionMap[key] = make([]any, 0)
		}
		filteredOptionMap[key] = append(filteredOptionMap[key], value)
	}
	keyDependencies := map[string][]struct {
		key           string
		mustBePresent bool
	}{
		"who-is-low-limit":                         {{"who-is-high-limit", true}},
		"who-is-high-limit":                        {{"who-is-low-limit", true}},
		"who-has-device-instance-range-low-limit":  {{"who-has-device-instance-range-high-limit", true}, {"who-has-object*", true}},
		"who-has-device-instance-range-high-limit": {{"who-has-device-instance-range-low-limit", true}, {"who-has-object*", true}},
		"who-has-object-identifier-type":           {{"who-has-object-identifier-instance", true}, {"who-has-object-name", false}},
		"who-has-object-identifier-instance":       {{"who-has-object-identifier-type", true}, {"who-has-object-name", false}},
		"who-has-object-name":                      {{"who-has-object-identifier-instance", false}, {"who-has-object-identifier-type", false}},
	}
	for key, value := range keyDependencies {
		if _, ok := filteredOptionMap[key]; ok {
			for _, otherKey := range value {
				if strings.HasSuffix(otherKey.key, "*") {
					prefix := strings.TrimSuffix(otherKey.key, "*")
					mustBePresent := otherKey.mustBePresent
					var found bool
					for key := range filteredOptionMap {
						found = found || strings.HasPrefix(key, prefix)
					}
					if mustBePresent && !found {
						return nil, errors.Errorf("When %s is set one of %s must also be set", key, otherKey.key)
					} else if !mustBePresent && found {
						return nil, errors.Errorf("When %s is set none of %s must be set", key, otherKey.key)
					}
				} else if _, otherOk := filteredOptionMap[otherKey.key]; otherOk && !otherKey.mustBePresent {
					return nil, errors.Errorf("When %s is set %s must not be set", key, otherKey.key)
				} else if !otherOk && otherKey.mustBePresent {
					return nil, errors.Errorf("When %s is set %s must be set too", key, otherKey.key)
				}
			}
		}
	}
	if _, ok := filteredOptionMap["bacnet-port"]; ok {
		parsedInt, err := exactlyOneInt(filteredOptionMap, "bacnet-port")
		if err != nil {
			return nil, err
		}
		collectedOptions = append(collectedOptions, bacNetPort(parsedInt))
	} else {
		collectedOptions = append(collectedOptions, bacNetPort(47808))
	}

	if whoIsLow, whoIsHigh, ok, err := func() (whoIsLowLimit uint, whoIsHighLimit uint, ok bool, err error) {
		if _, limitPresent := filteredOptionMap["who-is-low-limit"]; !limitPresent {
			return
		}
		ok = true
		whoIsLowLimit, err = exactlyOneUint(filteredOptionMap, "who-is-low-limit")
		whoIsHighLimit, err = exactlyOneUint(filteredOptionMap, "who-is-high-limit")
		return
	}(); ok {
		collectedOptions = append(collectedOptions, whoIsLimits(whoIsLow, whoIsHigh))
	} else if err != nil {
		return nil, err
	}
	for key := range filteredOptionMap {
		if strings.HasPrefix(key, "who-has-object") {
			collectedOptions = append(collectedOptions, whoHasOption())
			break
		}
	}
	if whoHasDeviceInstanceRangeLowLimit, whoHasDeviceInstanceRangeHighLimit, ok, err := func() (whoIsLowLimit uint, whoIsHighLimit uint, ok bool, err error) {
		if _, limitPresent := filteredOptionMap["who-has-device-instance-range-low-limit"]; !limitPresent {
			return
		}
		ok = true
		whoIsLowLimit, err = exactlyOneUint(filteredOptionMap, "who-has-device-instance-range-low-limit")
		whoIsHighLimit, err = exactlyOneUint(filteredOptionMap, "who-has-device-instance-range-high-limit")
		return
	}(); ok {
		collectedOptions = append(collectedOptions, whoHasLimits(whoHasDeviceInstanceRangeLowLimit, whoHasDeviceInstanceRangeHighLimit))
	} else if err != nil {
		return nil, err
	}

	if whoHasObjectIdentifierType, objectIdentifierInstance, ok, err := func() (whoHasObjectIdentifierType string, whoHasObjectIdentifierInstance uint, ok bool, err error) {
		if _, limitPresent := filteredOptionMap["who-has-object-identifier-type"]; !limitPresent {
			return
		}
		ok = true
		whoHasObjectIdentifierType, err = exactlyOneString(filteredOptionMap, "who-has-object-identifier-type")
		whoHasObjectIdentifierInstance, err = exactlyOneUint(filteredOptionMap, "who-has-object-identifier-instance")
		return
	}(); ok {
		collectedOptions = append(collectedOptions, whoHasObjectIdentifier(whoHasObjectIdentifierType, objectIdentifierInstance))
	} else if err != nil {
		return nil, err
	}

	if _, ok := filteredOptionMap["who-has-object-name"]; ok {
		if name, err := exactlyOneString(filteredOptionMap, "who-has-object-name"); err != nil {
			return nil, err
		} else {
			collectedOptions = append(collectedOptions, whoHasObjectName(name))
		}
	}
	return NewProtocolSpecificOptions(collectedOptions...)
}

func exactlyOneInt(filteredOptionMap map[string][]any, key string) (int, error) {
	value, err := exactlyOne(filteredOptionMap, key)
	if err != nil {
		return 0, err
	}
	parsedInt, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "Error parsing option bacnet-port")
	}
	return int(parsedInt), nil
}

func exactlyOneUint(filteredOptionMap map[string][]any, key string) (uint, error) {
	value, err := exactlyOne(filteredOptionMap, key)
	if err != nil {
		return 0, err
	}
	parsedInt, err := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "Error parsing option bacnet-port")
	}
	return uint(parsedInt), nil
}
func exactlyOneString(filteredOptionMap map[string][]any, key string) (string, error) {
	value, err := exactlyOne(filteredOptionMap, key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", value), nil
}

func exactlyOne(filteredOptionMap map[string][]any, key string) (any, error) {
	values := filteredOptionMap[key]
	if len(values) != 1 {
		return nil, errors.Errorf("%s expects only one value", key)
	}
	return values[0], nil
}
