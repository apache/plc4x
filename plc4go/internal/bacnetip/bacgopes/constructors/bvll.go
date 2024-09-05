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
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func Result(i uint16) *bacgopes.Result {
	result, err := bacgopes.NewResult(bacgopes.WithResultBvlciResultCode(readWriteModel.BVLCResultCode(i)))
	if err != nil {
		panic(err)
	}
	return result
}

func WriteBroadcastDistributionTable(bdt ...*bacgopes.Address) *bacgopes.WriteBroadcastDistributionTable {
	writeBroadcastDistributionTable, err := bacgopes.NewWriteBroadcastDistributionTable(bacgopes.WithWriteBroadcastDistributionTableBDT(bdt...))
	if err != nil {
		panic(err)
	}
	return writeBroadcastDistributionTable
}

func ReadBroadcastDistributionTable() *bacgopes.ReadBroadcastDistributionTable {
	readBroadcastDistributionTable, err := bacgopes.NewReadBroadcastDistributionTable()
	if err != nil {
		panic(err)
	}
	return readBroadcastDistributionTable
}

func ReadBroadcastDistributionTableAck(bdt ...*bacgopes.Address) *bacgopes.ReadBroadcastDistributionTableAck {
	readBroadcastDistributionTable, err := bacgopes.NewReadBroadcastDistributionTableAck(bacgopes.WithReadBroadcastDistributionTableAckBDT(bdt...))
	if err != nil {
		panic(err)
	}
	return readBroadcastDistributionTable
}

func ForwardedNPDU(addr *bacgopes.Address, pduBytes []byte) *bacgopes.ForwardedNPDU {
	npdu, err := bacgopes.NewForwardedNPDU(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...)), bacgopes.WithForwardedNPDUAddress(addr))
	if err != nil {
		panic(err)
	}
	return npdu
}

func RegisterForeignDevice(ttl uint16) *bacgopes.RegisterForeignDevice {
	registerForeignDevice, err := bacgopes.NewRegisterForeignDevice(bacgopes.WithRegisterForeignDeviceBvlciTimeToLive(ttl))
	if err != nil {
		panic(err)
	}
	return registerForeignDevice
}

func ReadForeignDeviceTable() *bacgopes.ReadForeignDeviceTable {
	readForeignDeviceTable, err := bacgopes.NewReadForeignDeviceTable()
	if err != nil {
		panic(err)
	}
	return readForeignDeviceTable
}

func FDTEntry() (entry *bacgopes.FDTEntry) {
	return &bacgopes.FDTEntry{}
}

func ReadForeignDeviceTableAck(fdts ...*bacgopes.FDTEntry) *bacgopes.ReadForeignDeviceTableAck {
	readForeignDeviceTableAck, err := bacgopes.NewReadForeignDeviceTableAck(bacgopes.WithReadForeignDeviceTableAckFDT(fdts...))
	if err != nil {
		panic(err)
	}
	return readForeignDeviceTableAck
}

func DeleteForeignDeviceTableEntry(address *bacgopes.Address) *bacgopes.DeleteForeignDeviceTableEntry {
	deleteForeignDeviceTableEntry, err := bacgopes.NewDeleteForeignDeviceTableEntry(bacgopes.WithDeleteForeignDeviceTableEntryAddress(address))
	if err != nil {
		panic(err)
	}
	return deleteForeignDeviceTableEntry
}

func DistributeBroadcastToNetwork(pduBytes []byte) *bacgopes.DistributeBroadcastToNetwork {
	distributeBroadcastToNetwork, err := bacgopes.NewDistributeBroadcastToNetwork(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...)))
	if err != nil {
		panic(err)
	}
	return distributeBroadcastToNetwork
}

func OriginalUnicastNPDU(pduBytes []byte) *bacgopes.OriginalUnicastNPDU {
	npdu, err := bacgopes.NewOriginalUnicastNPDU(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...)))
	if err != nil {
		panic(err)
	}
	return npdu
}

func OriginalBroadcastNPDU(pduBytes []byte) *bacgopes.OriginalBroadcastNPDU {
	npdu, err := bacgopes.NewOriginalBroadcastNPDU(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...)))
	if err != nil {
		panic(err)
	}
	return npdu
}
