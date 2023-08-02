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
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

type WhoIsIAmServicesRequirements interface {
	Request(pdu _PDU) error
}

type WhoIsIAmServices struct {
	WhoIsIAmServicesRequirements
	*Capability
	localDevice *LocalDeviceObject
}

func NewWhoIsIAmServices(whoIsIAmServicesRequirements WhoIsIAmServicesRequirements) (*WhoIsIAmServices, error) {
	w := &WhoIsIAmServices{}
	w.WhoIsIAmServicesRequirements = whoIsIAmServicesRequirements
	w.Capability = NewCapability()
	return w, nil
}

func (w *WhoIsIAmServices) Startup() error {
	log.Debug().Msg("Startup")

	// send a global broadcast I-Am
	return w.IAm(nil)
}

func (w *WhoIsIAmServices) WhoIs(lowLimit, highLimit *uint, address *Address) error {
	log.Debug().Msg("WhoIs")

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
	whoIs := model.NewBACnetUnconfirmedServiceRequestWhoIs(model.CreateBACnetContextTagUnsignedInteger(0, deviceInstanceRangeLowLimit), model.CreateBACnetContextTagUnsignedInteger(1, deviceInstanceRangeHighLimit), 0)

	log.Debug().Stringer("whoIs", whoIs).Msg("WhoIs")

	return w.Request(NewPDU(whoIs, WithPDUDestination(address)))
}

// DoWhoIsRequest respond to a Who-Is request.
func (w *WhoIsIAmServices) DoWhoIsRequest(apdu _PDU) error {
	log.Debug().Stringer("apdu", apdu).Msg("DoWhoIsRequest")

	// ignore this if there's no local device
	if w.localDevice == nil {
		log.Debug().Msg("No local device")
	}

	// extract the parameters
	var lowLimit, highLimit *uint
	if deviceInstanceRangeLowLimit := apdu.GetMessage().(model.BACnetUnconfirmedServiceRequestWhoIs).GetDeviceInstanceRangeLowLimit(); deviceInstanceRangeLowLimit != nil {
		_lowLimit := uint(deviceInstanceRangeLowLimit.GetActualValue())
		lowLimit = &_lowLimit
	}
	if deviceInstanceRangeHighLimit := apdu.GetMessage().(model.BACnetUnconfirmedServiceRequestWhoIs).GetDeviceInstanceRangeHighLimit(); deviceInstanceRangeHighLimit != nil {
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
	log.Debug().Msg("IAm")

	// this requires a local device
	if w.localDevice == nil {
		log.Debug().Msg("no local device")
		return nil
	}

	iAm := model.NewBACnetUnconfirmedServiceRequestIAm(
		model.CreateBACnetApplicationTagObjectIdentifier(ObjectIdentifierStringToTuple(w.localDevice.ObjectIdentifier)),
		model.CreateBACnetApplicationTagUnsignedInteger(uint(*w.localDevice.MaximumApduLengthAccepted)),
		model.CreatBACnetSegmentationTagged(*w.localDevice.SegmentationSupported),
		model.CreateBACnetVendorIdApplicationTagged(w.localDevice.VendorIdentifier),
		0,
	)

	// defaults to a global broadcast
	if address == nil {
		address = NewGlobalBroadcast(nil)
	}
	log.Debug().Stringer("iAm", iAm).Msg("IAm")

	return w.Request(NewPDU(iAm, WithPDUDestination(address)))
}

// DoIAmRequest responds to an I-Am request.
func (w *WhoIsIAmServices) DoIAmRequest(iam *model.BACnetUnconfirmedServiceRequestIAm) error {
	// TODO: implement me... upstream impl empty
	return nil
}

func ObjectIdentifierStringToTuple(objectIdentifier string) (objectType uint16, instance uint32) {
	split := strings.Split(objectIdentifier, ":")
	if len(split) != 2 {
		panic("broken object identifier")
	}
	parsedObjectType, err := strconv.Atoi(split[0])
	if err != nil {
		panic(err)
	}
	objectType = uint16(parsedObjectType)
	parsedInstance, err := strconv.Atoi(split[1])
	if err != nil {
		panic(err)
	}
	instance = uint32(parsedInstance)
	return
}
