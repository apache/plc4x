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
    public partial class DateAndTime : IMessage
    {
        public byte Year { get; }
        public byte Month { get; }
        public byte Day { get; }
        public byte Hour { get; }
        public byte Minutes { get; }
        public byte Seconds { get; }
        public ushort Msec { get; }
        public byte Dow { get; }

        public DateAndTime(byte year, byte month, byte day, byte hour, byte minutes, byte seconds, ushort msec, byte dow)
        {
            Year = year;
            Month = month;
            Day = day;
            Hour = hour;
            Minutes = minutes;
            Seconds = seconds;
            Msec = msec;
            Dow = dow;
        }

        public static DateAndTime StaticParse(ReadBuffer readBuffer)
        {
            var year = readBuffer.ReadByte("year", 8);
            var month = readBuffer.ReadByte("month", 8);
            var day = readBuffer.ReadByte("day", 8);
            var hour = readBuffer.ReadByte("hour", 8);
            var minutes = readBuffer.ReadByte("minutes", 8);
            var seconds = readBuffer.ReadByte("seconds", 8);
            var msec = readBuffer.ReadUshort("msec", 12);
            var dow = readBuffer.ReadByte("dow", 4);
            return new DateAndTime(year, month, day, hour, minutes, seconds, msec, dow);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("year", 8, Year);
            writeBuffer.WriteByte("month", 8, Month);
            writeBuffer.WriteByte("day", 8, Day);
            writeBuffer.WriteByte("hour", 8, Hour);
            writeBuffer.WriteByte("minutes", 8, Minutes);
            writeBuffer.WriteByte("seconds", 8, Seconds);
            writeBuffer.WriteUshort("msec", 12, Msec);
            writeBuffer.WriteByte("dow", 4, Dow);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 12;
            lengthInBits += 4;
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
