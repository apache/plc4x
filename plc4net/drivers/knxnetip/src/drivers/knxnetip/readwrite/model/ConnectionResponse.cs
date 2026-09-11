//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

// Code generated from the mspec by plc4net-code-gen. DO NOT EDIT.

using System;
using System.Linq;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{
    public partial class ConnectionResponse : KnxNetIpMessage
    {
        public override ushort MsgType => (ushort) (0x0206);

        public byte CommunicationChannelId { get; }
        public Status Status { get; }
        public HPAIDataEndpoint? HpaiDataEndpoint { get; }
        public ConnectionResponseDataBlock? ConnectionResponseDataBlock { get; }

        public ConnectionResponse(byte communicationChannelId, Status status, HPAIDataEndpoint? hpaiDataEndpoint, ConnectionResponseDataBlock? connectionResponseDataBlock)
        {
            CommunicationChannelId = communicationChannelId;
            Status = status;
            HpaiDataEndpoint = hpaiDataEndpoint;
            ConnectionResponseDataBlock = connectionResponseDataBlock;
        }

        public static new ConnectionResponse StaticParse(ReadBuffer readBuffer)
        {
            var communicationChannelId = readBuffer.ReadByte("communicationChannelId", 8);
            var status = (Status) readBuffer.ReadByte("status", 8);
            HPAIDataEndpoint? hpaiDataEndpoint = null;
            if ((status == Status.NO_ERROR))
            {
                hpaiDataEndpoint = HPAIDataEndpoint.StaticParse(readBuffer);
            }
            ConnectionResponseDataBlock? connectionResponseDataBlock = null;
            if ((status == Status.NO_ERROR))
            {
                connectionResponseDataBlock = ConnectionResponseDataBlock.StaticParse(readBuffer);
            }
            return new ConnectionResponse(communicationChannelId, status, hpaiDataEndpoint, connectionResponseDataBlock);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("communicationChannelId", 8, CommunicationChannelId);
            writeBuffer.WriteByte("status", 8, (byte) Status);
            if (HpaiDataEndpoint != null)
            {
                HpaiDataEndpoint.Serialize(writeBuffer);
            }
            if (ConnectionResponseDataBlock != null)
            {
                ConnectionResponseDataBlock.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += (HpaiDataEndpoint?.GetLengthInBits() ?? 0);
            lengthInBits += (ConnectionResponseDataBlock?.GetLengthInBits() ?? 0);
            return lengthInBits;
        }

    }
}
