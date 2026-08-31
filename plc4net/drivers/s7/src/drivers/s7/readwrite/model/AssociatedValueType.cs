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
    public partial class AssociatedValueType : IMessage
    {
        public DataTransportErrorCode ReturnCode { get; }
        public DataTransportSize TransportSize { get; }
        public ushort ValueLength { get; }
        public System.Collections.Generic.List<byte> Data { get; }

        public AssociatedValueType(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort valueLength, System.Collections.Generic.List<byte> data)
        {
            ReturnCode = returnCode;
            TransportSize = transportSize;
            ValueLength = valueLength;
            Data = data;
        }

        public static AssociatedValueType StaticParse(ReadBuffer readBuffer)
        {
            var returnCode = (DataTransportErrorCode) readBuffer.ReadByte("returnCode", 8);
            var transportSize = (DataTransportSize) readBuffer.ReadByte("transportSize", 8);
            var valueLength = (ushort) (S7StaticHelper.RightShift3(readBuffer, transportSize));
            var data = new System.Collections.Generic.List<byte>();
            var _dataCnt = (int) (S7StaticHelper.EventItemLength(readBuffer, valueLength));
            for (var _dataI = 0; _dataI < _dataCnt; _dataI++)
            {
                data.Add(readBuffer.ReadByte("data", 8));
            }
            return new AssociatedValueType(returnCode, transportSize, valueLength, data);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("returnCode", 8, (byte) ReturnCode);
            writeBuffer.WriteByte("transportSize", 8, (byte) TransportSize);
            S7StaticHelper.LeftShift3(writeBuffer, ValueLength);
            foreach (var _e in Data)
            {
                writeBuffer.WriteByte("data", 8, _e);
            }
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += ((2) * 8);
            lengthInBits += (Data.Count * 8);
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
