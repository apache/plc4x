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
	readwriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
)

func CreateRequestContext(cBusMessage readwriteModel.CBusMessage) readwriteModel.RequestContext {
	return CreateRequestContextWithInfoCallback(cBusMessage, func(_ string) {})
}

func CreateRequestContextWithInfoCallback(cBusMessage readwriteModel.CBusMessage, infoCallBack func(string)) readwriteModel.RequestContext {
	switch cBusMessage := cBusMessage.(type) {
	case readwriteModel.CBusMessageToServerExactly:
		switch request := cBusMessage.GetRequest().(type) {
		case readwriteModel.RequestDirectCommandAccessExactly:
			sendIdentifyRequestBefore := false
			infoCallBack("CAL request detected")
			switch request.GetCalData().(type) {
			case readwriteModel.CALDataIdentifyExactly:
				sendIdentifyRequestBefore = true
			}
			return readwriteModel.NewRequestContext(true, false, sendIdentifyRequestBefore)
		case readwriteModel.RequestCommandExactly:
			switch command := request.GetCbusCommand().(type) {
			case readwriteModel.CBusCommandDeviceManagementExactly:
				infoCallBack("CAL request detected")
				return readwriteModel.NewRequestContext(true, false, false)
			case readwriteModel.CBusCommandPointToPointExactly:
				sendIdentifyRequestBefore := false
				infoCallBack("CAL request detected")
				switch command.GetCommand().GetCalData().(type) {
				case readwriteModel.CALDataIdentifyExactly:
					sendIdentifyRequestBefore = true
				}
				return readwriteModel.NewRequestContext(true, false, sendIdentifyRequestBefore)
			case readwriteModel.CBusCommandPointToMultiPointExactly:
				switch command := command.GetCommand().(type) {
				case readwriteModel.CBusPointToMultiPointCommandStatusExactly:
					var sendStatusRequestLevelBefore bool
					switch command.GetStatusRequest().(type) {
					case readwriteModel.StatusRequestLevelExactly:
						sendStatusRequestLevelBefore = true
					}
					infoCallBack("SAL status request detected")
					return readwriteModel.NewRequestContext(false, sendStatusRequestLevelBefore, false)
				}
			case readwriteModel.CBusCommandPointToPointToMultiPointExactly:
				switch command := command.GetCommand().(type) {
				case readwriteModel.CBusPointToPointToMultiPointCommandStatusExactly:
					var sendStatusRequestLevelBefore bool
					switch command.GetStatusRequest().(type) {
					case readwriteModel.StatusRequestLevelExactly:
						sendStatusRequestLevelBefore = true
					}
					infoCallBack("SAL status request detected")
					return readwriteModel.NewRequestContext(false, sendStatusRequestLevelBefore, false)
				}
			}
		case readwriteModel.RequestObsoleteExactly:
			sendIdentifyRequestBefore := false
			infoCallBack("CAL request detected")
			switch request.GetCalData().(type) {
			case readwriteModel.CALDataIdentifyExactly:
				sendIdentifyRequestBefore = true
			}
			return readwriteModel.NewRequestContext(true, false, sendIdentifyRequestBefore)
		}
	case readwriteModel.CBusMessageToClientExactly:
		// We received a request so we need to reset our flags
		return readwriteModel.NewRequestContext(false, false, false)
	}
	return readwriteModel.NewRequestContext(false, false, false)
}
