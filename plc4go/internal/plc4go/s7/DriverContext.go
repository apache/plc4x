/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package s7

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
)

type DriverContext struct {
	PassiveMode             bool
	CallingTsapId           uint16
	CalledTsapId            uint16
	CotpTpduSize            model.COTPTpduSize
	PduSize                 uint16
	MaxAmqCaller            uint16
	MaxAmqCallee            uint16
	ControllerType          ControllerType
	awaitSetupComplete      bool
	awaitDisconnectComplete bool
}

func NewDriverContext(configuration Configuration) (DriverContext, error) {
	callingTsapId := encodeS7TsapId(model.DeviceGroup_OTHERS, configuration.localRack, configuration.localSlot)
	calledTsapId := encodeS7TsapId(model.DeviceGroup_PG_OR_PC, configuration.remoteRack, configuration.remoteSlot)

	controllerType := configuration.controllerType
	if controllerType == ControllerType_UNKNOWN {
		controllerType = ControllerType_ANY
	}

	pduSize := configuration.pduSize
	// The Siemens LOGO device seems to only work with very limited settings,
	// so we're overriding some of the defaults.
	if controllerType == ControllerType_LOGO && pduSize == 1024 {
		pduSize = 480
	}

	// Initialize the parameters with initial version (Will be updated during the login process)
	cotpTpduSize := getNearestMatchingTpduSize(pduSize)
	// The PDU size is theoretically not bound by the COTP TPDU size, however having a larger
	// PDU size would make the code extremely complex. But even if the protocol would allow this
	// I have never seen this happen in reality. Making is smaller would unnecessarily limit the
	// size, so we're setting it to the maximum that can be included.
	pduSize = cotpTpduSize.SizeInBytes() - 16
	maxAmqCaller := configuration.maxAmqCaller
	maxAmqCallee := configuration.maxAmqCallee
	return DriverContext{
		CallingTsapId:  callingTsapId,
		CalledTsapId:   calledTsapId,
		ControllerType: controllerType,
		CotpTpduSize:   cotpTpduSize,
		PduSize:        pduSize,
		MaxAmqCaller:   maxAmqCaller,
		MaxAmqCallee:   maxAmqCallee,
	}, nil
}

func getNearestMatchingTpduSize(tpduSizeParameter uint16) model.COTPTpduSize {
	switch {
	case model.COTPTpduSize_SIZE_128.SizeInBytes() >= tpduSizeParameter:
		return model.COTPTpduSize_SIZE_128
	case model.COTPTpduSize_SIZE_256.SizeInBytes() >= tpduSizeParameter:
		return model.COTPTpduSize_SIZE_256
	case model.COTPTpduSize_SIZE_512.SizeInBytes() >= tpduSizeParameter:
		return model.COTPTpduSize_SIZE_512
	case model.COTPTpduSize_SIZE_1024.SizeInBytes() >= tpduSizeParameter:
		return model.COTPTpduSize_SIZE_1024
	case model.COTPTpduSize_SIZE_2048.SizeInBytes() >= tpduSizeParameter:
		return model.COTPTpduSize_SIZE_2048
	case model.COTPTpduSize_SIZE_4096.SizeInBytes() >= tpduSizeParameter:
		return model.COTPTpduSize_SIZE_4096
	case model.COTPTpduSize_SIZE_8192.SizeInBytes() >= tpduSizeParameter:
		return model.COTPTpduSize_SIZE_8192
	}
	return model.COTPTpduSize(0)
}
