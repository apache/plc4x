//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{

    public class CEMIFrameDataExt : CEMIFrame
    {
        // Accessors for discriminator values.
        public override bool GetNotAckFrame() {
            return true;
        }
        public override bool GetStandardFrame() {
            return false;
        }
        public override bool GetPolling() {
            return false;
        }

        // Properties.
        public bool GroupAddress { get; }
        public byte HopCount { get; }
        public byte ExtendedFrameFormat { get; }
        public KnxAddress SourceAddress { get; }
        public sbyte[] DestinationAddress { get; }
        public byte DataLength { get; }
        public TPCI Tcpi { get; }
        public byte Counter { get; }
        public APCI Apci { get; }
        public sbyte DataFirstByte { get; }
        public sbyte[] Data { get; }
        public byte Crc { get; }

        public CEMIFrameDataExt(bool repeated, CEMIPriority priority, bool acknowledgeRequested, bool errorFlag, bool groupAddress, byte hopCount, byte extendedFrameFormat, KnxAddress sourceAddress, sbyte[] destinationAddress, byte dataLength, TPCI tcpi, byte counter, APCI apci, sbyte dataFirstByte, sbyte[] data, byte crc)
            : base(repeated, priority, acknowledgeRequested, errorFlag)
        {
            GroupAddress = groupAddress;
            HopCount = hopCount;
            ExtendedFrameFormat = extendedFrameFormat;
            SourceAddress = sourceAddress;
            DestinationAddress = destinationAddress;
            DataLength = dataLength;
            Tcpi = tcpi;
            Counter = counter;
            Apci = apci;
            DataFirstByte = dataFirstByte;
            Data = data;
            Crc = crc;
        }

    }

}
