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
using System.Linq;

namespace org.apache.plc4net.drivers.s7.test
{
    /// <summary>
    /// S7 frames built by hand from the protocol description, independently of the generated model,
    /// so that a request the driver sends and an answer it reads are checked against the bytes
    /// rather than against the code that produced them.
    /// </summary>
    internal static class S7Frames
    {
        /// <summary>The COTP Connection Confirm that a CPU answers the Connection Request with.</summary>
        public static byte[] ConnectionConfirm() => new byte[]
        {
            0x03, 0x00, 0x00, 0x0B, // TPKT: version, reserved, length
            0x06, 0xD0,             // COTP: length indicator, Connection Confirm
            0x00, 0x01, 0x00, 0x02, // destination and source reference
            0x00,                   // class 0
        };

        /// <summary>Wraps an S7 message in a TPKT header and a COTP Data Transfer header.</summary>
        public static byte[] Dt(params byte[] s7)
        {
            var length = 4 + 3 + s7.Length;
            return new byte[] { 0x03, 0x00, (byte)(length >> 8), (byte)length, 0x02, 0xF0, 0x80 }
                .Concat(s7).ToArray();
        }

        /// <summary>The S7 message of a job: ten header bytes, the parameter, the data.</summary>
        public static byte[] Job(ushort reference, byte[] parameter, byte[] data) => Message(0x01, reference, parameter, data, false);

        /// <summary>An acknowledgement with data (message type 3): twelve header bytes, the parameter, the data.</summary>
        public static byte[] AckData(ushort reference, byte[] parameter, byte[] data, byte errorClass = 0, byte errorCode = 0)
            => Message(0x03, reference, parameter, data, true, errorClass, errorCode);

        /// <summary>A bare acknowledgement (message type 2): how a CPU refuses a job before it looks at the items.</summary>
        public static byte[] Refusal(ushort reference, byte errorClass, byte errorCode)
            => Message(0x02, reference, Array.Empty<byte>(), Array.Empty<byte>(), true, errorClass, errorCode);

        private static byte[] Message(
            byte type, ushort reference, byte[] parameter, byte[] data, bool withError,
            byte errorClass = 0, byte errorCode = 0)
        {
            var message = new List<byte>
            {
                0x32, type, 0x00, 0x00,
                (byte)(reference >> 8), (byte)reference,
                (byte)(parameter.Length >> 8), (byte)parameter.Length,
                (byte)(data.Length >> 8), (byte)data.Length,
            };
            if (withError)
            {
                message.Add(errorClass);
                message.Add(errorCode);
            }
            message.AddRange(parameter);
            message.AddRange(data);
            return message.ToArray();
        }

        /// <summary>The CPU's answer to Setup Communication, wrapped for the transport.</summary>
        public static byte[] SetupAck(ushort reference, ushort pduLength)
            => Dt(AckData(reference,
                new byte[]
                {
                    0xF0, 0x00,             // function, reserved
                    0x00, 0x01, 0x00, 0x01, // maximum jobs the CPU takes and gives
                    (byte)(pduLength >> 8), (byte)pduLength,
                },
                Array.Empty<byte>()));

        /// <summary>
        /// The answer to a Read Var job. Each item is its return code, its transport size and the
        /// data; the length field counts bits for the sizes 3, 4 and 5 and bytes otherwise, and an
        /// item with an odd number of bytes is padded unless it is the last.
        /// </summary>
        public static byte[] ReadAck(ushort reference, params (byte ReturnCode, byte TransportSize, byte[] Data)[] items)
        {
            var data = new List<byte>();
            for (var i = 0; i < items.Length; i++)
            {
                var (returnCode, size, bytes) = items[i];
                var length = size == 0x03 ? 1 : size == 0x04 || size == 0x05 ? bytes.Length * 8 : bytes.Length;
                data.Add(returnCode);
                data.Add(size);
                data.Add((byte)(length >> 8));
                data.Add((byte)length);
                data.AddRange(bytes);
                if (i < items.Length - 1 && bytes.Length % 2 == 1)
                {
                    data.Add(0x00);
                }
            }
            return Dt(AckData(reference, new byte[] { 0x04, (byte)items.Length }, data.ToArray()));
        }

        /// <summary>The answer to a Write Var job: one return code per item.</summary>
        public static byte[] WriteAck(ushort reference, params byte[] returnCodes)
            => Dt(AckData(reference, new byte[] { 0x05, (byte)returnCodes.Length }, returnCodes));

        /// <summary>A Setup Communication request as the driver sends it, for the given PDU length.</summary>
        public static byte[] SetupRequest(ushort reference, ushort pduLength)
            => Dt(Job(reference,
                new byte[]
                {
                    0xF0, 0x00,
                    0x00, 0x08, 0x00, 0x08, // the jobs the driver takes and gives: 8 and 8
                    (byte)(pduLength >> 8), (byte)pduLength,
                },
                Array.Empty<byte>()));

        /// <summary>One S7-ANY item of a request: syntax id 0x10, transport size, elements, DB, area, address.</summary>
        public static byte[] AddressItem(byte transportSize, ushort dbNumber, byte area, int byteOffset, int bitOffset)
        {
            var address = (byteOffset << 3) | bitOffset;
            return new byte[]
            {
                0x12, 0x0A, 0x10, transportSize,
                0x00, 0x01,
                (byte)(dbNumber >> 8), (byte)dbNumber,
                area,
                (byte)(address >> 16), (byte)(address >> 8), (byte)address,
            };
        }

        /// <summary>A Read Var job over the given address items.</summary>
        public static byte[] ReadRequest(ushort reference, params byte[][] items)
            => Dt(Job(reference,
                new byte[] { 0x04, (byte)items.Length }.Concat(items.SelectMany(i => i)).ToArray(),
                Array.Empty<byte>()));

        /// <summary>
        /// A Write Var job over address items and their data items. A data item is a return code
        /// of 0xFF, the transport size, the length and the data, padded to an even length unless
        /// it is the last.
        /// </summary>
        public static byte[] WriteRequest(
            ushort reference, byte[][] items, params (byte TransportSize, byte[] Data)[] values)
        {
            var data = new List<byte>();
            for (var i = 0; i < values.Length; i++)
            {
                var (size, bytes) = values[i];
                var length = size == 0x03 ? 1 : size == 0x04 || size == 0x05 ? bytes.Length * 8 : bytes.Length;
                data.Add(0xFF);
                data.Add(size);
                data.Add((byte)(length >> 8));
                data.Add((byte)length);
                data.AddRange(bytes);
                if (i < values.Length - 1 && bytes.Length % 2 == 1)
                {
                    data.Add(0x00);
                }
            }
            return Dt(Job(reference,
                new byte[] { 0x05, (byte)items.Length }.Concat(items.SelectMany(i => i)).ToArray(),
                data.ToArray()));
        }
    }
}
