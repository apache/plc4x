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
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
)

func CreateRequestContext(cBusMessage readWriteModel.CBusMessage) readWriteModel.RequestContext {
	return CreateRequestContextWithInfoCallback(cBusMessage, nil)
}

func CreateRequestContextWithInfoCallback(cBusMessage readWriteModel.CBusMessage, infoCallBack func(string)) readWriteModel.RequestContext {
	if infoCallBack == nil {
		infoCallBack = func(_ string) {}
	}
	switch cBusMessage := cBusMessage.(type) {
	case readWriteModel.CBusMessageToServerExactly:
		switch request := cBusMessage.GetRequest().(type) {
		case readWriteModel.RequestDirectCommandAccessExactly:
			sendIdentifyRequestBefore := false
			infoCallBack("CAL request detected")
			switch request.GetCalData().(type) {
			case readWriteModel.CALDataIdentifyExactly:
				sendIdentifyRequestBefore = true
			}
			return readWriteModel.NewRequestContext(sendIdentifyRequestBefore)
		case readWriteModel.RequestCommandExactly:
			switch command := request.GetCbusCommand().(type) {
			case readWriteModel.CBusCommandPointToPointExactly:
				sendIdentifyRequestBefore := false
				infoCallBack("CAL request detected")
				switch command.GetCommand().GetCalData().(type) {
				case readWriteModel.CALDataIdentifyExactly:
					sendIdentifyRequestBefore = true
				}
				return readWriteModel.NewRequestContext(sendIdentifyRequestBefore)
			}
		case readWriteModel.RequestObsoleteExactly:
			sendIdentifyRequestBefore := false
			infoCallBack("CAL request detected")
			switch request.GetCalData().(type) {
			case readWriteModel.CALDataIdentifyExactly:
				sendIdentifyRequestBefore = true
			}
			return readWriteModel.NewRequestContext(sendIdentifyRequestBefore)
		}
	case readWriteModel.CBusMessageToClientExactly:
		// We received a request, so we need to reset our flags
		return readWriteModel.NewRequestContext(false)
	}
	return readWriteModel.NewRequestContext(false)
}
