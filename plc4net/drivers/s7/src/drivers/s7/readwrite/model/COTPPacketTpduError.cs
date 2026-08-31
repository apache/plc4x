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

namespace org.apache.plc4net.drivers.s7.readwrite.model
{
    public partial class COTPPacketTpduError : COTPPacket
    {
        public override byte TpduCode => (byte) (0x70);

        public ushort DestinationReference { get; }
        public byte RejectCause { get; }

        public COTPPacketTpduError(ushort destinationReference, byte rejectCause, System.Collections.Generic.List<COTPParameter> parameters, S7Message? payload) : base(parameters, payload)
        {
            DestinationReference = destinationReference;
            RejectCause = rejectCause;
        }

        public static COTPPacketTpduError StaticParse(ReadBuffer readBuffer, ushort cotpLen, byte headerLength, int _startPos)
        {
            var destinationReference = readBuffer.ReadUshort("destinationReference", 16);
            var rejectCause = readBuffer.ReadByte("rejectCause", 8);
            var parameters = new System.Collections.Generic.List<COTPParameter>();
            var _parametersEnd = readBuffer.GetPos() + (int) (((headerLength + 1) - ((readBuffer.GetPos() - _startPos) / 8))) * 8;
            while (readBuffer.GetPos() < _parametersEnd)
            {
                parameters.Add(COTPParameter.StaticParse(readBuffer, (byte) (((headerLength + 1) - ((readBuffer.GetPos() - _startPos) / 8)))));
            }
            S7Message? payload = null;
            if ((((readBuffer.GetPos() - _startPos) / 8) < cotpLen))
            {
                payload = S7Message.StaticParse(readBuffer);
            }
            return new COTPPacketTpduError(destinationReference, rejectCause, parameters, payload);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("destinationReference", 16, DestinationReference);
            writeBuffer.WriteByte("rejectCause", 8, RejectCause);
            foreach (var _e in Parameters)
            {
                _e.Serialize(writeBuffer);
            }
            if (Payload != null)
            {
                Payload.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += Parameters.Sum(_e => _e.GetLengthInBits());
            lengthInBits += (Payload?.GetLengthInBits() ?? 0);
            return lengthInBits;
        }

    }
}
