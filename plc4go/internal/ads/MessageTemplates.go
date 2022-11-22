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

package ads

import adsModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"

func (m *Connection) NewAdsReadDeviceInfoRequest() adsModel.AmsTCPPacket {
	return adsModel.NewAmsTCPPacket(
		adsModel.NewAdsReadDeviceInfoRequest(m.configuration.targetAmsNetId, uint16(adsModel.DefaultAmsPorts_RUNTIME_SYSTEM_01),
			// TODO: Replace 800 with constant.
			m.configuration.sourceAmsNetId, 800, 0, m.driverContext.getInvokeId()))
}

func (m *Connection) NewAdsReadRequest(indexGroup uint32, indexOffset uint32, length uint32) adsModel.AmsTCPPacket {
	return adsModel.NewAmsTCPPacket(
		adsModel.NewAdsReadRequest(indexGroup, indexOffset, length,
			m.configuration.targetAmsNetId, m.configuration.targetAmsPort,
			m.configuration.sourceAmsNetId, m.configuration.sourceAmsPort, 0, m.driverContext.getInvokeId()))
}

func (m *Connection) NewAdsWriteRequest(indexGroup uint32, indexOffset uint32, data []byte) adsModel.AmsTCPPacket {
	return adsModel.NewAmsTCPPacket(
		adsModel.NewAdsWriteRequest(
			indexGroup, indexOffset, data,
			m.configuration.targetAmsNetId, m.configuration.targetAmsPort,
			m.configuration.sourceAmsNetId, m.configuration.sourceAmsPort,
			0, m.driverContext.getInvokeId()))
}

func (m *Connection) NewAdsReadWriteRequest(indexGroup uint32, indexOffset uint32, readLength uint32, items []adsModel.AdsMultiRequestItem, writeData []byte) adsModel.AmsTCPPacket {
	return adsModel.NewAmsTCPPacket(
		adsModel.NewAdsReadWriteRequest(
			indexGroup, indexOffset, readLength, items, writeData,
			m.configuration.targetAmsNetId, m.configuration.targetAmsPort,
			m.configuration.sourceAmsNetId, m.configuration.sourceAmsPort,
			0, m.driverContext.getInvokeId()))
}

func (m *Connection) NewAdsAddDeviceNotificationRequest(indexGroup uint32, indexOffset uint32, length uint32, transmissionMode adsModel.AdsTransMode, maxDelay uint32, cycleTime uint32) adsModel.AmsTCPPacket {
	return adsModel.NewAmsTCPPacket(
		adsModel.NewAdsAddDeviceNotificationRequest(
			indexGroup, indexOffset, length, transmissionMode, maxDelay, cycleTime,
			m.configuration.targetAmsNetId, m.configuration.targetAmsPort,
			m.configuration.sourceAmsNetId, m.configuration.sourceAmsPort,
			0, m.driverContext.getInvokeId()))
}

func (m *Connection) NewAdsDeleteDeviceNotificationRequest(notificationHandle uint32) adsModel.AmsTCPPacket {
	return adsModel.NewAmsTCPPacket(
		adsModel.NewAdsDeleteDeviceNotificationRequest(
			notificationHandle,
			m.configuration.targetAmsNetId, m.configuration.targetAmsPort,
			m.configuration.sourceAmsNetId, m.configuration.sourceAmsPort,
			0, m.driverContext.getInvokeId()))
}
