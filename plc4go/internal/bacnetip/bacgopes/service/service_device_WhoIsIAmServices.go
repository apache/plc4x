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

package service

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/apdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/capability"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type WhoIsIAmServicesRequirements interface {
	Request(args Args, kwArgs KWArgs) error
	RegisterHelperFn(name string, fn func(apdu APDU) error) error
}

//go:generate plc4xGenerator -type=WhoIsIAmServices -prefix=service_device_
type WhoIsIAmServices struct {
	_requirements WhoIsIAmServicesRequirements `ignore:"true"`
	Capability

	localDevice *LocalDeviceObject

	log zerolog.Logger
}

func NewWhoIsIAmServices(localLog zerolog.Logger, whoIsIAmServicesRequirements WhoIsIAmServicesRequirements, opts ...func(*WhoIsIAmServices)) (*WhoIsIAmServices, error) {
	w := &WhoIsIAmServices{
		_requirements: whoIsIAmServicesRequirements,
		log:           localLog,
	}
	for _, opt := range opts {
		opt(w)
	}
	w.Capability = NewCapability()
	if err := w._requirements.RegisterHelperFn(fmt.Sprintf("Do_%T", &WhoIsRequest{}), w.DoWhoIsRequest); err != nil {
		return nil, errors.Wrap(err, "registering function failed")
	}
	return w, nil
}

func WithWhoIsIAmServicesLocalDevice(localDevice *LocalDeviceObject) func(*WhoIsIAmServices) {
	return func(w *WhoIsIAmServices) {
		w.localDevice = localDevice
	}
}

func (w *WhoIsIAmServices) Startup() error {
	w.log.Debug().Msg("Startup")

	// send a global broadcast I-Am
	return w.IAm(nil)
}

func (w *WhoIsIAmServices) WhoIs(lowLimit, highLimit *uint, address *Address) error {
	w.log.Debug().Msg("WhoIs")

	var deviceInstanceRangeLowLimit, deviceInstanceRangeHighLimit uint
	if lowLimit != nil {
		if highLimit == nil {
			return errors.New("highLimit required")
		}
		if *lowLimit < 0 || *lowLimit > 4194303 {
			return errors.New("lowLimit out of range")
		}

		// low limit is fine
		deviceInstanceRangeLowLimit = *lowLimit
	}
	if highLimit != nil {
		if lowLimit == nil {
			return errors.New("lowLimit required")
		}
		if *highLimit < 0 || *highLimit > 4194303 {
			return errors.New("highLimit out of range")
		}

		// low limit is fine
		deviceInstanceRangeHighLimit = *highLimit
	}

	// Build a request
	whoIs := readWriteModel.NewBACnetUnconfirmedServiceRequestWhoIs(readWriteModel.CreateBACnetContextTagUnsignedInteger(0, deviceInstanceRangeLowLimit), readWriteModel.CreateBACnetContextTagUnsignedInteger(1, deviceInstanceRangeHighLimit), 0)

	w.log.Debug().Stringer("whoIs", whoIs).Msg("WhoIs")

	return w._requirements.Request(NA(NewPDU(NoArgs, NKW(KWCompRootMessage, whoIs, KWCPCIDestination, address))), NoKWArgs)
}

// DoWhoIsRequest respond to a Who-Is request.
func (w *WhoIsIAmServices) DoWhoIsRequest(apdu APDU) error {
	w.log.Debug().Stringer("apdu", apdu).Msg("DoWhoIsRequest")

	// ignore this if there's no local device
	if w.localDevice == nil {
		w.log.Debug().Msg("No local device")
	}

	// TODO: ugly hacky hacky, better feat from the orginal api
	whois := apdu.GetRootMessage().(readWriteModel.APDUUnconfirmedRequest).GetServiceRequest().(readWriteModel.BACnetUnconfirmedServiceRequestWhoIs)

	// extract the parameters
	var lowLimit, highLimit *uint
	if deviceInstanceRangeLowLimit := whois.GetDeviceInstanceRangeLowLimit(); deviceInstanceRangeLowLimit != nil {
		_lowLimit := uint(deviceInstanceRangeLowLimit.GetActualValue())
		lowLimit = &_lowLimit
	}
	if deviceInstanceRangeHighLimit := whois.GetDeviceInstanceRangeHighLimit(); deviceInstanceRangeHighLimit != nil {
		_highLimit := uint(deviceInstanceRangeHighLimit.GetActualValue())
		highLimit = &_highLimit
	}
	// check for consistent parameters
	if lowLimit != nil {
		if highLimit == nil {
			return errors.New("deviceInstanceRangeHighLimit required")
		}
		if *lowLimit < 0 || *lowLimit > 4194303 {
			return errors.New("deviceInstanceRangeHighLimit out of range")
		}
	}
	if highLimit != nil {
		if lowLimit == nil {
			return errors.New("deviceInstanceRangeLowLimit required")
		}
		if *highLimit < 0 || *highLimit > 4194303 {
			return errors.New("deviceInstanceRangeHighLimit out of range")
		}
	}

	// see we should respond
	if lowLimit != nil {
		if uint(w.localDevice.ObjectIdentifier[1]) < *lowLimit {
			return nil
		}
	}
	if highLimit != nil {
		if uint(w.localDevice.ObjectIdentifier[1]) > *highLimit {
			return nil
		}
	}

	// generate an I-Am
	return w.IAm(apdu.GetPDUSource())
}

func (w *WhoIsIAmServices) IAm(address *Address) error {
	w.log.Debug().Msg("IAm")

	// this requires a local device
	if w.localDevice == nil {
		w.log.Debug().Msg("no local device")
		return nil
	}

	iAm, err := NewIAmRequest(
		NoArgs,
		NKW(
			KnownKey("iAmDeviceIdentifier"), w.localDevice.ObjectIdentifier,
			KnownKey("maxAPDULengthAccepted"), w.localDevice.MaximumApduLengthAccepted,
			KnownKey("segmentationSupported"), w.localDevice.SegmentationSupported,
			KnownKey("vendorID"), w.localDevice.VendorIdentifier,
		),
	)
	if err != nil {
		return errors.Wrap(err, "IAm creation failed")
	}

	// defaults to a global broadcast
	if address == nil {
		address = NewGlobalBroadcast(nil)
	}
	iAm.SetPDUDestination(address)
	w.log.Debug().Stringer("iAm", iAm).Msg("IAm")

	return w._requirements.Request(NA(NewPDU(NoArgs, NKW(KWCompRootMessage, iAm, KWCPCIDestination, address))), NoKWArgs)
}

// DoIAmRequest responds to an I-Am request.
func (w *WhoIsIAmServices) DoIAmRequest(apdu PDU) error {
	// TODO: implement me... upstream impl empty
	return nil
}
