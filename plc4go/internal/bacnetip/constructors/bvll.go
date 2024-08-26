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

package constructors

import (
	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func Result(i uint16) *bacnetip.Result {
	result, err := bacnetip.NewResult(bacnetip.WithResultBvlciResultCode(readWriteModel.BVLCResultCode(i)))
	if err != nil {
		panic(err)
	}
	return result
}

func WriteBroadcastDistributionTable(bdt ...*bacnetip.Address) *bacnetip.WriteBroadcastDistributionTable {
	writeBroadcastDistributionTable, err := bacnetip.NewWriteBroadcastDistributionTable(bacnetip.WithWriteBroadcastDistributionTableBDT(bdt...))
	if err != nil {
		panic(err)
	}
	return writeBroadcastDistributionTable
}

func ReadBroadcastDistributionTable() *bacnetip.ReadBroadcastDistributionTable {
	readBroadcastDistributionTable, err := bacnetip.NewReadBroadcastDistributionTable()
	if err != nil {
		panic(err)
	}
	return readBroadcastDistributionTable
}

func ReadBroadcastDistributionTableAck(bdt ...*bacnetip.Address) *bacnetip.ReadBroadcastDistributionTableAck {
	readBroadcastDistributionTable, err := bacnetip.NewReadBroadcastDistributionTableAck(bacnetip.WithReadBroadcastDistributionTableAckBDT(bdt...))
	if err != nil {
		panic(err)
	}
	return readBroadcastDistributionTable
}

func ForwardedNPDU(addr *bacnetip.Address, pduBytes []byte) *bacnetip.ForwardedNPDU {
	npdu, err := bacnetip.NewForwardedNPDU(bacnetip.NewPDU(bacnetip.NewMessageBridge(pduBytes...)), bacnetip.WithForwardedNPDUAddress(addr))
	if err != nil {
		panic(err)
	}
	return npdu
}

func RegisterForeignDevice(ttl uint16) *bacnetip.RegisterForeignDevice {
	registerForeignDevice, err := bacnetip.NewRegisterForeignDevice(bacnetip.WithRegisterForeignDeviceBvlciTimeToLive(ttl))
	if err != nil {
		panic(err)
	}
	return registerForeignDevice
}

func ReadForeignDeviceTable() *bacnetip.ReadForeignDeviceTable {
	readForeignDeviceTable, err := bacnetip.NewReadForeignDeviceTable()
	if err != nil {
		panic(err)
	}
	return readForeignDeviceTable
}

func FDTEntry() (entry bacnetip.FDTEntry) {
	return
}

func ReadForeignDeviceTableAck(fdts ...bacnetip.FDTEntry) *bacnetip.ReadForeignDeviceTableAck {
	readForeignDeviceTableAck, err := bacnetip.NewReadForeignDeviceTableAck(bacnetip.WithReadForeignDeviceTableAckFDT(fdts...))
	if err != nil {
		panic(err)
	}
	return readForeignDeviceTableAck
}

func DeleteForeignDeviceTableEntry(address *bacnetip.Address) *bacnetip.DeleteForeignDeviceTableEntry {
	deleteForeignDeviceTableEntry, err := bacnetip.NewDeleteForeignDeviceTableEntry(bacnetip.WithDeleteForeignDeviceTableEntryAddress(address))
	if err != nil {
		panic(err)
	}
	return deleteForeignDeviceTableEntry
}

func DistributeBroadcastToNetwork(pduBytes []byte) *bacnetip.DistributeBroadcastToNetwork {
	distributeBroadcastToNetwork, err := bacnetip.NewDistributeBroadcastToNetwork(bacnetip.NewPDU(bacnetip.NewMessageBridge(pduBytes...)))
	if err != nil {
		panic(err)
	}
	return distributeBroadcastToNetwork
}

func OriginalUnicastNPDU(pduBytes []byte) *bacnetip.OriginalUnicastNPDU {
	npdu, err := bacnetip.NewOriginalUnicastNPDU(bacnetip.NewPDU(bacnetip.NewMessageBridge(pduBytes...)))
	if err != nil {
		panic(err)
	}
	return npdu
}

func OriginalBroadcastNPDU(pduBytes []byte) *bacnetip.OriginalBroadcastNPDU {
	npdu, err := bacnetip.NewOriginalBroadcastNPDU(bacnetip.NewPDU(bacnetip.NewMessageBridge(pduBytes...)))
	if err != nil {
		panic(err)
	}
	return npdu
}
