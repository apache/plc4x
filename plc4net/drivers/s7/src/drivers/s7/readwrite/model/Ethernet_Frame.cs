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
    public partial class Ethernet_Frame : IMessage
    {
        public MacAddress Destination { get; }
        public MacAddress Source { get; }
        public Ethernet_FramePayload Payload { get; }

        public Ethernet_Frame(MacAddress destination, MacAddress source, Ethernet_FramePayload payload)
        {
            Destination = destination;
            Source = source;
            Payload = payload;
        }

        public static Ethernet_Frame StaticParse(ReadBuffer readBuffer)
        {
            var destination = MacAddress.StaticParse(readBuffer);
            var source = MacAddress.StaticParse(readBuffer);
            var payload = Ethernet_FramePayload.StaticParse(readBuffer);
            return new Ethernet_Frame(destination, source, payload);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            Destination.Serialize(writeBuffer);
            Source.Serialize(writeBuffer);
            Payload.Serialize(writeBuffer);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += Destination.GetLengthInBits();
            lengthInBits += Source.GetLengthInBits();
            lengthInBits += Payload.GetLengthInBits();
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
