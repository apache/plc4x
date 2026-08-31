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
    public partial class S7VarPayloadDataItem : IMessage
    {
        public DataTransportErrorCode ReturnCode { get; }
        public DataTransportSize TransportSize { get; }
        public byte[] Data { get; }

        public S7VarPayloadDataItem(DataTransportErrorCode returnCode, DataTransportSize transportSize, byte[] data)
        {
            ReturnCode = returnCode;
            TransportSize = transportSize;
            Data = data;
        }

        public static S7VarPayloadDataItem StaticParse(ReadBuffer readBuffer, bool _lastItem = false)
        {
            var returnCode = (DataTransportErrorCode) readBuffer.ReadByte("returnCode", 8);
            var transportSize = (DataTransportSize) readBuffer.ReadByte("transportSize", 8);
            var dataLength = readBuffer.ReadUshort("dataLength", 16);
            var data = readBuffer.ReadByteArray("data", (int) ((transportSize.GetSizeInBits() ? (int) System.Math.Ceiling((double) ((dataLength / 8.0))) : dataLength)) * 8);
            {
                var _timesPadding = (int) ((!_lastItem ? (data.Length % 2) : 0));
                while (_timesPadding-- > 0)
                {
                    readBuffer.ReadByte("padding", 8);
                }
            }
            return new S7VarPayloadDataItem(returnCode, transportSize, data);
        }

        public void Serialize(WriteBuffer writeBuffer) => Serialize(writeBuffer, false);

        public void Serialize(WriteBuffer writeBuffer, bool _lastItem)
        {
            writeBuffer.WriteByte("returnCode", 8, (byte) ReturnCode);
            writeBuffer.WriteByte("transportSize", 8, (byte) TransportSize);
            writeBuffer.WriteUshort("dataLength", 16, (ushort) ((Data.Length * ((TransportSize == DataTransportSize.BIT) ? 1 : (TransportSize.GetSizeInBits() ? 8 : 1)))));
            writeBuffer.WriteByteArray("data", Data);
            {
                var _timesPadding = (int) ((!_lastItem ? (Data.Length % 2) : 0));
                while (_timesPadding-- > 0)
                {
                    writeBuffer.WriteByte("padding", 8, (byte) (0x00));
                }
            }
        }

        public int GetLengthInBits() => GetLengthInBits(false);

        public int GetLengthInBits(bool _lastItem)
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += (Data.Length * 8);
            lengthInBits += (((!_lastItem ? (Data.Length % 2) : 0)) * 8);
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits(false) / 8;
    }
}
