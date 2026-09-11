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
    public partial class LDataReq : CEMI
    {
        public override byte MessageCode => (byte) (0x11);

        public byte AdditionalInformationLength { get; }
        public System.Collections.Generic.List<CEMIAdditionalInformation> AdditionalInformation { get; }
        public LDataFrame DataFrame { get; }

        public LDataReq(byte additionalInformationLength, System.Collections.Generic.List<CEMIAdditionalInformation> additionalInformation, LDataFrame dataFrame)
        {
            AdditionalInformationLength = additionalInformationLength;
            AdditionalInformation = additionalInformation;
            DataFrame = dataFrame;
        }

        public static new LDataReq StaticParse(ReadBuffer readBuffer, ushort size)
        {
            var additionalInformationLength = readBuffer.ReadByte("additionalInformationLength", 8);
            var additionalInformation = new System.Collections.Generic.List<CEMIAdditionalInformation>();
            var _additionalInformationEnd = readBuffer.GetPos() + (int) (additionalInformationLength) * 8;
            while (readBuffer.GetPos() < _additionalInformationEnd)
            {
                additionalInformation.Add(CEMIAdditionalInformation.StaticParse(readBuffer));
            }
            var dataFrame = LDataFrame.StaticParse(readBuffer);
            return new LDataReq(additionalInformationLength, additionalInformation, dataFrame);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("additionalInformationLength", 8, AdditionalInformationLength);
            foreach (var _e in AdditionalInformation)
            {
                _e.Serialize(writeBuffer);
            }
            DataFrame.Serialize(writeBuffer);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += AdditionalInformation.Sum(_e => _e.GetLengthInBits());
            lengthInBits += DataFrame.GetLengthInBits();
            return lengthInBits;
        }

    }
}
