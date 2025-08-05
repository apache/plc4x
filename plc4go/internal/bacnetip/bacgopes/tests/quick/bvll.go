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

package quick

import (
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func Result(i uint16) *bvll.Result {
	result, err := bvll.NewResult(ToPtr(readWriteModel.BVLCResultCode(i)), NoArgs, NoKWArgs())
	if err != nil {
		panic(err)
	}
	return result
}

func WriteBroadcastDistributionTable(bdt ...*pdu.Address) *bvll.WriteBroadcastDistributionTable {
	writeBroadcastDistributionTable, err := bvll.NewWriteBroadcastDistributionTable(bdt, NoArgs, NoKWArgs())
	if err != nil {
		panic(err)
	}
	return writeBroadcastDistributionTable
}

func ReadBroadcastDistributionTable() *bvll.ReadBroadcastDistributionTable {
	readBroadcastDistributionTable, err := bvll.NewReadBroadcastDistributionTable(Nothing())
	if err != nil {
		panic(err)
	}
	return readBroadcastDistributionTable
}

func ReadBroadcastDistributionTableAck(bdt ...*pdu.Address) *bvll.ReadBroadcastDistributionTableAck {
	readBroadcastDistributionTable, err := bvll.NewReadBroadcastDistributionTableAck(bdt, NoArgs, NoKWArgs())
	if err != nil {
		panic(err)
	}
	return readBroadcastDistributionTable
}

func ForwardedNPDU(addr *pdu.Address, pduBytes []byte) *bvll.ForwardedNPDU {
	npdu, err := bvll.NewForwardedNPDU(addr, NA(pduBytes), NoKWArgs())
	if err != nil {
		panic(err)
	}
	return npdu
}

func RegisterForeignDevice(ttl uint16) *bvll.RegisterForeignDevice {
	registerForeignDevice, err := bvll.NewRegisterForeignDevice(&ttl, NoArgs, NoKWArgs())
	if err != nil {
		panic(err)
	}
	return registerForeignDevice
}

func ReadForeignDeviceTable(addr *pdu.Address) *bvll.ReadForeignDeviceTable {
	readForeignDeviceTable, err := bvll.NewReadForeignDeviceTable(NoArgs, NKW(KWCPCIDestination, addr))
	if err != nil {
		panic(err)
	}
	return readForeignDeviceTable
}

func FDTEntry() (entry *bvll.FDTEntry) {
	return &bvll.FDTEntry{}
}

func ReadForeignDeviceTableAck(fdts ...*bvll.FDTEntry) *bvll.ReadForeignDeviceTableAck {
	readForeignDeviceTableAck, err := bvll.NewReadForeignDeviceTableAck(fdts, NoArgs, NoKWArgs())
	if err != nil {
		panic(err)
	}
	return readForeignDeviceTableAck
}

func DeleteForeignDeviceTableEntry(address *pdu.Address) *bvll.DeleteForeignDeviceTableEntry {
	deleteForeignDeviceTableEntry, err := bvll.NewDeleteForeignDeviceTableEntry(address, NoArgs, NoKWArgs())
	if err != nil {
		panic(err)
	}
	return deleteForeignDeviceTableEntry
}

func DistributeBroadcastToNetwork(pduBytes []byte) *bvll.DistributeBroadcastToNetwork {
	distributeBroadcastToNetwork, err := bvll.NewDistributeBroadcastToNetwork(NA(pduBytes), NoKWArgs())
	if err != nil {
		panic(err)
	}
	return distributeBroadcastToNetwork
}

func OriginalUnicastNPDU(pduBytes []byte) *bvll.OriginalUnicastNPDU {
	npdu, err := bvll.NewOriginalUnicastNPDU(NA(pduBytes), NoKWArgs())
	if err != nil {
		panic(err)
	}
	return npdu
}

func OriginalBroadcastNPDU(pduBytes []byte) *bvll.OriginalBroadcastNPDU {
	npdu, err := bvll.NewOriginalBroadcastNPDU(NA(pduBytes), NoKWArgs())
	if err != nil {
		panic(err)
	}
	return npdu
}
