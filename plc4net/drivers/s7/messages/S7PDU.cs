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

using System.Collections.Generic;
using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.spi.generation;

namespace org.apache.plc4net.drivers.s7.messages
{
    /// <summary>
    /// S7 protocol constants and PDU construction. Every request is built with the
    /// generated S7 wire model (<see cref="S7Message"/>) - it computes header and
    /// item lengths, the S7-ANY transport-size codes, the memory-area codes and the
    /// bit/byte length quirk of a data item, and the same model round-trips the
    /// shared <c>ParserSerializerTestsuite.xml</c> vectors byte-identical.
    /// </summary>
    public static class S7Constants
    {
        public const byte ProtocolId = 0x32;

        // Message types (ROSCTR).
        public const byte JobRequest = 0x01;
        public const byte Ack = 0x02;
        public const byte AckData = 0x03;
        public const byte UserData = 0x07;

        // Parameter function codes.
        public const byte SetupCommunication = 0xF0;
        public const byte ReadVar = 0x04;
        public const byte WriteVar = 0x05;

        // Setup Communication request values - matching the shared S7 test vector.
        // The CPU negotiates the PDU length down to its own capability.
        public const ushort RequestedPduLength = 1008;
        public const ushort RequestedMaxAmq = 8;

        /// <summary>
        /// Setup Communication request (S7 job, function 0xF0). Must be exchanged once,
        /// right after the COTP CR/CC handshake - a real CPU ignores Read / Write Var
        /// until it has run.
        /// </summary>
        public static byte[] BuildSetupCommunication(ushort tpduRef) =>
            Serialize(new S7MessageRequest(tpduRef,
                new S7ParameterSetupCommunication(RequestedMaxAmq, RequestedMaxAmq, RequestedPduLength),
                null));

        /// <summary>Read Var request PDU (S7 header + parameter part) for one or more tags.</summary>
        public static byte[] BuildReadRequest(ushort tpduRef, IReadOnlyList<S7Tag> tags)
        {
            var items = new List<S7VarRequestParameterItem>(tags.Count);
            foreach (var tag in tags)
            {
                items.Add(new S7VarRequestParameterItemAddress(ToAddress(tag)));
            }
            return Serialize(new S7MessageRequest(tpduRef, new S7ParameterReadVarRequest(items), null));
        }

        /// <summary>Write Var request PDU (S7 header + parameter part + data part).</summary>
        public static byte[] BuildWriteRequest(ushort tpduRef, IReadOnlyList<(S7Tag Tag, byte[] Data)> items)
        {
            var paramItems = new List<S7VarRequestParameterItem>(items.Count);
            var payloadItems = new List<S7VarPayloadDataItem>(items.Count);
            foreach (var (tag, data) in items)
            {
                paramItems.Add(new S7VarRequestParameterItemAddress(ToAddress(tag)));
                payloadItems.Add(new S7VarPayloadDataItem(
                    DataTransportErrorCode.OK, WritePayloadSize(tag), data));
            }
            return Serialize(new S7MessageRequest(tpduRef,
                new S7ParameterWriteVarRequest(paramItems),
                new S7PayloadWriteVarRequest(payloadItems)));
        }

        private static byte[] Serialize(S7Message message)
        {
            var writeBuffer = new WriteBuffer();
            message.Serialize(writeBuffer);
            return writeBuffer.GetBytes();
        }

        internal static S7AddressAny ToAddress(S7Tag tag) =>
            new S7AddressAny(
                AddressTransportSize(tag),
                numberOfElements: 1,
                dbNumber: (ushort) tag.DbNumber,
                area: MemoryAreaFor(tag.Area),
                byteAddress: (ushort) tag.ByteOffset,
                bitAddress: (byte) (tag.BitOffset >= 0 ? tag.BitOffset : 0));

        private static TransportSize AddressTransportSize(S7Tag tag)
        {
            // A counter / timer is addressed with its own transport size, not
            // as a WORD — the CPU reads the S7-ANY area (COUNTERS / TIMERS)
            // together with TransportSize.COUNTER.
            if (tag.Area is S7Tag.AreaType.Counter or S7Tag.AreaType.Timer)
            {
                return TransportSize.COUNTER;
            }
            if (tag.BitOffset >= 0)
            {
                return TransportSize.BOOL;
            }
            return tag.DataTypeSize switch
            {
                1 => TransportSize.BYTE,
                2 => TransportSize.WORD,
                4 => TransportSize.DWORD,
                _ => TransportSize.BYTE,
            };
        }

        private static DataTransportSize WritePayloadSize(S7Tag tag) =>
            tag.BitOffset >= 0 ? DataTransportSize.BIT : DataTransportSize.BYTE_WORD_DWORD;

        internal static MemoryArea MemoryAreaFor(S7Tag.AreaType area) => area switch
        {
            S7Tag.AreaType.DataBlock => MemoryArea.DATA_BLOCKS,
            S7Tag.AreaType.Merker => MemoryArea.FLAGS_MARKERS,
            S7Tag.AreaType.Input => MemoryArea.INPUTS,
            S7Tag.AreaType.Output => MemoryArea.OUTPUTS,
            S7Tag.AreaType.Counter => MemoryArea.COUNTERS,
            S7Tag.AreaType.Timer => MemoryArea.TIMERS,
            _ => MemoryArea.FLAGS_MARKERS,
        };
    }
}
