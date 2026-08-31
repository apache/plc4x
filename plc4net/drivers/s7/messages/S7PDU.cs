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

using System;
using System.Collections.Generic;

namespace org.apache.plc4net.drivers.s7.messages
{
    /// <summary>
    /// S7 protocol constants and PDU construction helpers.
    /// </summary>
    public static class S7Constants
    {
        public const byte ProtocolId = 0x32;

        // Message types
        public const byte JobRequest = 0x01;
        public const byte Ack = 0x02;
        public const byte AckData = 0x03;
        public const byte UserData = 0x07;

        // Function codes
        public const byte ReadVar = 0x04;
        public const byte WriteVar = 0x05;

        // Transport sizes
        public const byte TransportSizeBit = 0x03;
        public const byte TransportSizeByte = 0x04;
        public const byte TransportSizeWord = 0x05;
        public const byte TransportSizeDword = 0x06;

        /// <summary>
        /// Builds a Read Var request PDU for one or more S7 tags.
        /// Returns the complete COTP payload (S7 header + parameter part).
        /// </summary>
        public static byte[] BuildReadRequest(ushort pduRef, List<S7Tag> tags)
        {
            // Parameter part: function code + item count + items
            var paramData = new List<byte>();
            paramData.Add(ReadVar);  // function code
            paramData.Add((byte)tags.Count); // item count

            foreach (var tag in tags)
            {
                paramData.Add(0x12); // item spec: variable specification
                paramData.Add((byte)(10)); // length of following address data
                paramData.Add(0x10); // syntax ID: S7-ANY
                paramData.Add(TransportSizeForTag(tag)); // transport size
                paramData.Add(0x00); // length hi (element count, 2 bytes)
                paramData.Add(0x01); // length lo (1 element)
                paramData.Add(0x00); // DB number high
                paramData.Add(0x00); // DB number low
                paramData.Add(AreaCodeForTag(tag)); // area
                // Address: byte.bit (3 bytes, big-endian)
                var addr = (int)((long)tag.ByteOffset * 8 + (tag.BitOffset >= 0 ? tag.BitOffset : 0));
                paramData.Add((byte)((addr >> 16) & 0xFF));
                paramData.Add((byte)((addr >> 8) & 0xFF));
                paramData.Add((byte)(addr & 0xFF));

                // If DB area, insert DB number
                if (tag.Area == S7Tag.AreaType.DataBlock)
                {
                    // Insert DB number before area code (back-patch over the zeros)
                    var dbIdx = paramData.Count - 6; // after transport/granularity
                    paramData[dbIdx] = (byte)(tag.DbNumber >> 8);
                    paramData[dbIdx + 1] = (byte)(tag.DbNumber & 0xFF);
                }
            }

            return BuildS7Header(pduRef, paramData.ToArray(), null);
        }

        /// <summary>
        /// Builds a Write Var request PDU.
        /// </summary>
        public static byte[] BuildWriteRequest(ushort pduRef, List<(S7Tag tag, byte[] data)> items)
        {
            var paramData = new List<byte>();
            paramData.Add(WriteVar);
            paramData.Add((byte)items.Count);

            var dataPart = new List<byte>();
            dataPart.Add(WriteVar); // function code in data part too
            dataPart.Add((byte)items.Count);

            foreach (var (tag, data) in items)
            {
                paramData.Add(0x12);
                paramData.Add(10);
                paramData.Add(0x10);
                paramData.Add(TransportSizeForTag(tag));
                paramData.Add(0x00); // length hi (element count, 2 bytes)
                paramData.Add(0x01); // length lo (1 element)
                paramData.Add((byte)(tag.DbNumber >> 8));
                paramData.Add((byte)(tag.DbNumber & 0xFF));
                paramData.Add(AreaCodeForTag(tag));
                var addr = (int)((long)tag.ByteOffset * 8 + (tag.BitOffset >= 0 ? tag.BitOffset : 0));
                paramData.Add((byte)((addr >> 16) & 0xFF));
                paramData.Add((byte)((addr >> 8) & 0xFF));
                paramData.Add((byte)(addr & 0xFF));

                // Data part
                dataPart.Add(0x00); // return code (filled by PLC)
                dataPart.Add(TransportSizeForTag(tag));
                dataPart.Add((byte)(data.Length >> 8));
                dataPart.Add((byte)(data.Length & 0xFF));
                dataPart.AddRange(data);
            }

            return BuildS7Header(pduRef, paramData.ToArray(), dataPart.ToArray());
        }

        private static byte[] BuildS7Header(ushort pduRef, byte[] paramData, byte[] dataPart)
        {
            var paramLen = paramData.Length;
            var dataLen = dataPart?.Length ?? 0;
            var totalLen = 10 + paramLen + dataLen; // 10-byte S7 header

            var frame = new byte[totalLen];
            frame[0] = ProtocolId;
            frame[1] = JobRequest;
            frame[2] = 0x00; frame[3] = 0x00; // reserved
            frame[4] = (byte)(pduRef >> 8);
            frame[5] = (byte)(pduRef & 0xFF);
            frame[6] = (byte)(paramLen >> 8);
            frame[7] = (byte)(paramLen & 0xFF);
            frame[8] = (byte)(dataLen >> 8);
            frame[9] = (byte)(dataLen & 0xFF);

            Array.Copy(paramData, 0, frame, 10, paramLen);
            if (dataPart != null)
                Array.Copy(dataPart, 0, frame, 10 + paramLen, dataLen);

            return frame;
        }

        public static byte TransportSizeForTag(S7Tag tag)
        {
            if (tag.BitOffset >= 0 && tag.DataTypeSize == 1)
                return TransportSizeBit;
            return tag.DataTypeSize switch
            {
                1 => TransportSizeByte,
                2 => TransportSizeWord,
                4 => TransportSizeDword,
                _ => TransportSizeByte
            };
        }

        public static byte AreaCodeForTag(S7Tag tag)
        {
            // Values per Siemens S7-comm specification, not positional in an array.
            return tag.Area switch
            {
                S7Tag.AreaType.DataBlock => 0x82,
                S7Tag.AreaType.Merker => 0x83,
                S7Tag.AreaType.Input => 0x81,
                S7Tag.AreaType.Output => 0x80,
                S7Tag.AreaType.Counter => 0x1C,
                S7Tag.AreaType.Timer => 0x1D,
                _ => 0x83
            };
        }
    }
}
