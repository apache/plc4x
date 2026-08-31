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

using System.Globalization;
using System.Text.RegularExpressions;
using org.apache.plc4net.model;

namespace org.apache.plc4net.drivers.s7
{
    /// <summary>
    /// An S7 tag — an address in a Siemens S7 PLC.
    /// </summary>
    /// <remarks>
    /// Supported address syntax:
    ///   %DB{num}.DB{type}{byteOffset}  — Data block
    ///   %M{byteOffset}.{bitOffset}     — Merker / flag
    ///   %I{byteOffset}.{bitOffset}     — Input
    ///   %Q{byteOffset}.{bitOffset}     — Output
    ///   %C{num}                         — Counter
    ///   %T{num}                         — Timer
    ///
    /// Examples: %DB1.DBW10, %DB1.DBX0.0, %M0.0, %I0.0, %Q4.5
    /// </remarks>
    public class S7Tag : IPlcTag
    {
        private static readonly Regex DbPattern = new Regex(
            @"^%?DB(?<db>\d+)\.DB(?<type>[XBWDC])(?<offset>\d+)(\.(?<bit>\d))?$",
            RegexOptions.IgnoreCase | RegexOptions.Compiled);

        // C/T checked before M/I/Q to avoid the [MIQC] class capturing C and T.
        private static readonly Regex CTpattern = new Regex(
            @"^%?(?<area>[TC])(?<num>\d+)$",
            RegexOptions.IgnoreCase | RegexOptions.Compiled);

        private static readonly Regex MioPattern = new Regex(
            @"^%?(?<area>[MIQ])(?<offset>\d+)(\.(?<bit>\d))?$",
            RegexOptions.IgnoreCase | RegexOptions.Compiled);

        public enum AreaType { DataBlock, Merker, Input, Output, Counter, Timer }

        public AreaType Area { get; }
        public int DbNumber { get; }
        public int ByteOffset { get; }
        public int BitOffset { get; } // -1 = whole byte/word/dword
        public int DataTypeSize { get; } // 1=byte, 2=word, 4=dword

        private S7Tag(AreaType area, int dbNumber, int byteOffset, int bitOffset, int size)
        {
            Area = area;
            DbNumber = dbNumber;
            ByteOffset = byteOffset;
            BitOffset = bitOffset;
            DataTypeSize = size;
        }

        /// <summary>Parses an S7 address string.</summary>
        public static S7Tag Parse(string tagAddress)
        {
            if (string.IsNullOrWhiteSpace(tagAddress))
                throw new S7DriverException("Tag address must not be empty.");

            // DB: %DB1.DBW10, %DB1.DBX0.0, %DB1.DBD20
            var dbMatch = DbPattern.Match(tagAddress);
            if (dbMatch.Success)
            {
                var db = int.Parse(dbMatch.Groups["db"].Value, CultureInfo.InvariantCulture);
                var offset = int.Parse(dbMatch.Groups["offset"].Value, CultureInfo.InvariantCulture);
                var type = dbMatch.Groups["type"].Value.ToUpperInvariant();
                var size = type switch { "X" => 1, "B" => 1, "W" => 2, "D" => 4, "C" => 1, _ => 1 };
                var bit = dbMatch.Groups["bit"].Success
                    ? int.Parse(dbMatch.Groups["bit"].Value, CultureInfo.InvariantCulture) : -1;
                return new S7Tag(AreaType.DataBlock, db, offset, bit, size);
            }

            // T/C: %T0, %C5 — checked before MioPattern to avoid [MIQ] capture.
            var ctMatch = CTpattern.Match(tagAddress);
            if (ctMatch.Success)
            {
                var area = ctMatch.Groups["area"].Value.ToUpperInvariant() == "T"
                    ? AreaType.Timer : AreaType.Counter;
                var num = int.Parse(ctMatch.Groups["num"].Value, CultureInfo.InvariantCulture);
                return new S7Tag(area, 0, num, -1, 2);
            }

            // M/I/Q: %M0.0, %I0.0, %Q4.5
            var mioMatch = MioPattern.Match(tagAddress);
            if (mioMatch.Success)
            {
                var area = mioMatch.Groups["area"].Value.ToUpperInvariant() switch
                {
                    "M" => AreaType.Merker,
                    "I" => AreaType.Input,
                    "Q" => AreaType.Output,
                    _ => AreaType.Merker
                };
                var offset = int.Parse(mioMatch.Groups["offset"].Value, CultureInfo.InvariantCulture);
                var bit = mioMatch.Groups["bit"].Success
                    ? int.Parse(mioMatch.Groups["bit"].Value, CultureInfo.InvariantCulture) : -1;
                return new S7Tag(area, 0, offset, bit, 1);
            }

            throw new S7DriverException(
                $"Cannot parse S7 tag address '{tagAddress}'. " +
                "Expected format: %DBn.DBTm, %MBn, %IBn, %QBn, %Cn, or %Tn.");
        }

        public override string ToString()
        {
            var prefix = Area switch
            {
                AreaType.DataBlock => $"%DB{DbNumber}.DB",
                AreaType.Merker => "%M",
                AreaType.Input => "%I",
                AreaType.Output => "%Q",
                AreaType.Counter => "%C",
                AreaType.Timer => "%T",
                _ => "?"
            };

            if (Area == AreaType.DataBlock)
            {
                var typeSuffix = DataTypeSize switch { 1 => BitOffset >= 0 ? "X" : "B", 2 => "W", 4 => "D", _ => "B" };
                var bitSuffix = BitOffset >= 0 ? $".{BitOffset}" : "";
                return $"{prefix}{typeSuffix}{ByteOffset}{bitSuffix}";
            }

            return BitOffset >= 0 ? $"{prefix}{ByteOffset}.{BitOffset}" : $"{prefix}{ByteOffset}";
        }
    }
}
