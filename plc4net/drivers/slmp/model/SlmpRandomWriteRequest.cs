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

namespace org.apache.plc4net.drivers.slmp.readwrite.model
{
    public partial class SlmpRandomWriteRequest : SlmpRequestData
    {
        public override ushort Command => (ushort) (0x1402);

        public byte NumberOfWordAccessPoints { get; }
        public byte NumberOfDoubleWordAccessPoints { get; }
        public System.Collections.Generic.List<SlmpWordWriteItem> WordAccessPoints { get; }
        public System.Collections.Generic.List<SlmpDoubleWordWriteItem> DoubleWordAccessPoints { get; }

        public SlmpRandomWriteRequest(byte numberOfWordAccessPoints, byte numberOfDoubleWordAccessPoints, System.Collections.Generic.List<SlmpWordWriteItem> wordAccessPoints, System.Collections.Generic.List<SlmpDoubleWordWriteItem> doubleWordAccessPoints)
        {
            NumberOfWordAccessPoints = numberOfWordAccessPoints;
            NumberOfDoubleWordAccessPoints = numberOfDoubleWordAccessPoints;
            WordAccessPoints = wordAccessPoints;
            DoubleWordAccessPoints = doubleWordAccessPoints;
        }

        public static new SlmpRandomWriteRequest StaticParse(ReadBuffer readBuffer, ushort command, ushort subCommand)
        {
            var numberOfWordAccessPoints = readBuffer.ReadByte("numberOfWordAccessPoints", 8);
            var numberOfDoubleWordAccessPoints = readBuffer.ReadByte("numberOfDoubleWordAccessPoints", 8);
            var wordAccessPoints = new System.Collections.Generic.List<SlmpWordWriteItem>();
            var _wordAccessPointsCnt = (int) (numberOfWordAccessPoints);
            for (var _wordAccessPointsI = 0; _wordAccessPointsI < _wordAccessPointsCnt; _wordAccessPointsI++)
            {
                wordAccessPoints.Add(SlmpWordWriteItem.StaticParse(readBuffer));
            }
            var doubleWordAccessPoints = new System.Collections.Generic.List<SlmpDoubleWordWriteItem>();
            var _doubleWordAccessPointsCnt = (int) (numberOfDoubleWordAccessPoints);
            for (var _doubleWordAccessPointsI = 0; _doubleWordAccessPointsI < _doubleWordAccessPointsCnt; _doubleWordAccessPointsI++)
            {
                doubleWordAccessPoints.Add(SlmpDoubleWordWriteItem.StaticParse(readBuffer));
            }
            return new SlmpRandomWriteRequest(numberOfWordAccessPoints, numberOfDoubleWordAccessPoints, wordAccessPoints, doubleWordAccessPoints);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("numberOfWordAccessPoints", 8, NumberOfWordAccessPoints);
            writeBuffer.WriteByte("numberOfDoubleWordAccessPoints", 8, NumberOfDoubleWordAccessPoints);
            foreach (var _e in WordAccessPoints)
            {
                _e.Serialize(writeBuffer);
            }
            foreach (var _e in DoubleWordAccessPoints)
            {
                _e.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += WordAccessPoints.Sum(_e => _e.GetLengthInBits());
            lengthInBits += DoubleWordAccessPoints.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

    }
}
