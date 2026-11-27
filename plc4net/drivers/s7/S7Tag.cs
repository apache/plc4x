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
using System.Globalization;
using System.Linq;
using System.Text.RegularExpressions;
using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.model;

namespace org.apache.plc4net.drivers.s7
{
    /// <summary>
    /// An S7 tag: an address in a Siemens S7 PLC.
    /// </summary>
    /// <remarks>
    /// <para>
    /// The syntax is the one plc4j uses, so the same address names the same variable from either
    /// language. The data type is part of the address:
    /// </para>
    /// <list type="bullet">
    /// <item><c>%DB{block}.DB{X|B|W|D}{byte}[.{bit}]:{TYPE}</c> - a data block, for example
    /// <c>%DB1.DBW10:INT</c> or <c>%DB1.DBX0.0:BOOL</c>.</item>
    /// <item><c>%DB{block}:{byte}[.{bit}]:{TYPE}</c> - the short form of the same.</item>
    /// <item><c>%DB{block}:{byte}:STRING({length})</c> - a fixed-length string in a data block.</item>
    /// <item><c>%{I|Q|M}[X|B|W|D]{byte}[.{bit}]:{TYPE}</c> - inputs, outputs and flags, for example
    /// <c>%M0.0:BOOL</c> or <c>%QW4:WORD</c>.</item>
    /// </list>
    /// <para>
    /// A <c>BOOL</c> carries its bit; no other type does. The size letter is optional, and when
    /// it is there it has to be the one the type goes by: <c>W</c> for an <c>INT</c>, <c>D</c> for
    /// a <c>REAL</c>.
    /// </para>
    /// <para>
    /// Array selections, the 64-bit types, the date and time types, counters, timers and the
    /// other memory areas are not supported here yet; they are rejected rather than read as something
    /// else. Fixed-length STRING and WSTRING are supported only in data blocks.
    /// </para>
    /// </remarks>
    public sealed class S7Tag : IPlcTag
    {
        /// <summary>The data types a tag can have: the scalars that one S7-ANY item reads in one piece.</summary>
        public static readonly IReadOnlyList<TransportSize> SupportedTypes = new[]
        {
            TransportSize.BOOL,
            TransportSize.BYTE,
            TransportSize.SINT,
            TransportSize.USINT,
            TransportSize.WORD,
            TransportSize.INT,
            TransportSize.UINT,
            TransportSize.DWORD,
            TransportSize.DINT,
            TransportSize.UDINT,
            TransportSize.REAL,
        };

        /// <summary>The highest data block number: it travels in 16 bits.</summary>
        public const int MaxDataBlock = ushort.MaxValue;

        /// <summary>The highest byte offset: the byte address of an S7-ANY item is 16 bits wide.</summary>
        public const int MaxByteOffset = ushort.MaxValue;

        // [0-9], not \d: in .NET \d matches every Unicode decimal digit, so a non-ASCII number would
        // pass the pattern and then blow up int.Parse instead of failing here as an invalid tag.
        private const string Tail =
            @"(?<byte>[0-9]{1,7})(\.(?<bit>[0-7]))?(?<array>\[[^\]]*\])?:(?<type>(S5)?[A-Za-z_]+)(\((?<stringLength>[0-9]{1,3})\))?$";

        private static readonly Regex DataBlockPattern = new Regex(
            @"^%DB(?<db>[0-9]{1,5})\.DB(?<size>[XBWD]?)" + Tail,
            RegexOptions.CultureInvariant | RegexOptions.IgnoreCase);

        private static readonly Regex DataBlockShortPattern = new Regex(
            @"^%DB(?<db>[0-9]{1,5}):" + Tail,
            RegexOptions.CultureInvariant | RegexOptions.IgnoreCase);

        private static readonly Regex AreaPattern = new Regex(
            @"^%(?<area>[A-Za-z])(?<size>[XBWD]?)" + Tail,
            RegexOptions.CultureInvariant | RegexOptions.IgnoreCase);

        public S7Tag(MemoryArea area, int dbNumber, int byteOffset, byte bitOffset, TransportSize dataType)
            : this(area, dbNumber, byteOffset, bitOffset, dataType, null)
        {
        }

        public S7Tag(
            MemoryArea area,
            int dbNumber,
            int byteOffset,
            byte bitOffset,
            TransportSize dataType,
            int? stringLength)
        {
            if (area != MemoryArea.DATA_BLOCKS && area != MemoryArea.INPUTS
                && area != MemoryArea.OUTPUTS && area != MemoryArea.FLAGS_MARKERS)
            {
                throw new PlcInvalidFieldException(
                    $"The memory area {area} is not supported yet; supported: data blocks, inputs, outputs and flags.");
            }
            if (!SupportedTypes.Contains(dataType) && dataType != TransportSize.STRING
                && dataType != TransportSize.WSTRING)
            {
                throw new PlcInvalidFieldException(UnsupportedType(dataType.ToString()));
            }
            if (dataType == TransportSize.STRING || dataType == TransportSize.WSTRING)
            {
                if (area != MemoryArea.DATA_BLOCKS)
                {
                    throw new PlcInvalidFieldException("S7 STRING and WSTRING tags are supported only in data blocks.");
                }
                if (!stringLength.HasValue || stringLength.Value < 1 || stringLength.Value > 254)
                {
                    throw new PlcInvalidFieldException("An S7 STRING or WSTRING length must be between 1 and 254 characters.");
                }
                if (bitOffset != 0)
                {
                    throw new PlcInvalidFieldException("An S7 STRING or WSTRING cannot have a bit offset.");
                }
                var sizeInBytes = (stringLength.Value + 2) * (dataType == TransportSize.WSTRING ? 2 : 1);
                if (byteOffset > MaxByteOffset - (sizeInBytes - 1))
                {
                    throw new PlcInvalidFieldException(
                        "The S7 STRING or WSTRING extends past the highest byte offset the address can carry.");
                }
            }
            else if (stringLength.HasValue)
            {
                throw new PlcInvalidFieldException("A string length can be specified only for an S7 STRING or WSTRING tag.");
            }
            if (dbNumber < 0 || dbNumber > MaxDataBlock)
            {
                throw new PlcInvalidFieldException(
                    $"The S7 data block number must be between 0 and {MaxDataBlock}, but was {dbNumber}.");
            }
            if (area != MemoryArea.DATA_BLOCKS && dbNumber != 0)
            {
                throw new PlcInvalidFieldException(
                    $"The S7 memory area {area} must use data block number 0, but was {dbNumber}.");
            }
            if (byteOffset < 0 || byteOffset > MaxByteOffset)
            {
                throw new PlcInvalidFieldException(
                    $"The S7 byte offset must be between 0 and {MaxByteOffset}, but was {byteOffset}.");
            }
            if (bitOffset > 7 || (dataType != TransportSize.BOOL && bitOffset != 0))
            {
                throw new PlcInvalidFieldException(
                    "The S7 bit offset must be 0 to 7, and only a BOOL has one that is not 0.");
            }

            Area = area;
            DbNumber = dbNumber;
            ByteOffset = byteOffset;
            BitOffset = bitOffset;
            DataType = dataType;
            StringLength = stringLength;
        }

        public MemoryArea Area { get; }

        /// <summary>The data block number; 0 for inputs, outputs and flags.</summary>
        public int DbNumber { get; }

        public int ByteOffset { get; }

        /// <summary>The bit within the byte; 0 unless the tag is a <c>BOOL</c>.</summary>
        public byte BitOffset { get; }

        public TransportSize DataType { get; }

        /// <summary>The declared maximum character count for a fixed-length STRING; null otherwise.</summary>
        public int? StringLength { get; }

        /// <summary>The bytes the value takes on the wire, including a string's length header.</summary>
        public int SizeInBytes => StringLength.HasValue
            ? (StringLength.Value + 2) * (DataType == TransportSize.WSTRING ? 2 : 1)
            : DataType.GetSizeInBytes();

        /// <summary>
        /// Parses a tag address such as <c>%DB1.DBW10:INT</c> or <c>%M0.0:BOOL</c>.
        /// </summary>
        /// <exception cref="PlcInvalidFieldException">The address is not a supported S7 tag address.</exception>
        public static S7Tag Parse(string tagAddress)
        {
            if (string.IsNullOrWhiteSpace(tagAddress))
            {
                throw new PlcInvalidFieldException("The S7 tag address must not be empty.");
            }

            var address = tagAddress.Trim();

            var match = DataBlockPattern.Match(address);
            if (match.Success)
            {
                return FromMatch(match, MemoryArea.DATA_BLOCKS, int.Parse(
                    match.Groups["db"].Value, CultureInfo.InvariantCulture));
            }

            match = DataBlockShortPattern.Match(address);
            if (match.Success)
            {
                return FromMatch(match, MemoryArea.DATA_BLOCKS, int.Parse(
                    match.Groups["db"].Value, CultureInfo.InvariantCulture));
            }

            match = AreaPattern.Match(address);
            if (match.Success)
            {
                return FromMatch(match, AreaOf(match.Groups["area"].Value), 0);
            }

            throw new PlcInvalidFieldException(
                $"'{tagAddress}' is not a supported S7 tag address. Expected " +
                "%DB{block}.DB{X|B|W|D}{byte}[.{bit}]:{TYPE}, %DB{block}:{byte}[.{bit}]:{TYPE} or " +
                "%{I|Q|M}[X|B|W|D]{byte}[.{bit}]:{TYPE}, for example %DB1.DBW10:INT or %M0.0:BOOL; " +
                "a data-block string uses %DB{block}:{byte}:STRING({length}) or :WSTRING({length}). The other memory areas " +
                "are not supported yet.");
        }

        private static S7Tag FromMatch(Match match, MemoryArea area, int dbNumber)
        {
            if (match.Groups["array"].Success)
            {
                throw new PlcInvalidFieldException(
                    $"'{match.Value}': array selections are not supported yet.");
            }

            var dataType = TypeOf(match.Groups["type"].Value);
            var byteOffset = int.Parse(match.Groups["byte"].Value, CultureInfo.InvariantCulture);
            int? stringLength = null;
            if (match.Groups["stringLength"].Success)
            {
                if (dataType != TransportSize.STRING && dataType != TransportSize.WSTRING)
                {
                    throw new PlcInvalidFieldException($"'{match.Value}': a length is supported only for STRING or WSTRING.");
                }
                stringLength = int.Parse(match.Groups["stringLength"].Value, CultureInfo.InvariantCulture);
            }
            else if (dataType == TransportSize.STRING || dataType == TransportSize.WSTRING)
            {
                throw new PlcInvalidFieldException($"'{match.Value}': STRING and WSTRING need a declared length, as in :STRING(10).");
            }

            byte bitOffset = 0;
            if (match.Groups["bit"].Success)
            {
                bitOffset = byte.Parse(match.Groups["bit"].Value, CultureInfo.InvariantCulture);
            }
            else if (dataType == TransportSize.BOOL)
            {
                throw new PlcInvalidFieldException($"'{match.Value}': a BOOL needs a bit offset, as in %M0.0:BOOL.");
            }

            var size = match.Groups["size"];
            if (size.Success && size.Length > 0
                && char.ToUpperInvariant(size.Value[0]) != (char)dataType.GetShortName())
            {
                throw new PlcInvalidFieldException(
                    $"'{match.Value}': the size letter '{size.Value}' does not go with the data type {dataType}, " +
                    $"which is '{(char)dataType.GetShortName()}'.");
            }

            return new S7Tag(area, dbNumber, byteOffset, bitOffset, dataType, stringLength);
        }

        private static MemoryArea AreaOf(string letter)
        {
            switch (letter.ToUpperInvariant())
            {
                case "I":
                    return MemoryArea.INPUTS;
                case "Q":
                    return MemoryArea.OUTPUTS;
                case "M":
                    return MemoryArea.FLAGS_MARKERS;
                default:
                    throw new PlcInvalidFieldException(
                        $"The memory area '{letter}' is not supported yet; supported: DB, I, Q and M.");
            }
        }

        private static TransportSize TypeOf(string name)
        {
            // Enum.TryParse would also accept a number; only the names are data types.
            foreach (var candidate in Enum.GetNames(typeof(TransportSize)))
            {
                if (string.Equals(candidate, name, StringComparison.OrdinalIgnoreCase))
                {
                    var type = (TransportSize)Enum.Parse(typeof(TransportSize), candidate);
                    if (!SupportedTypes.Contains(type) && type != TransportSize.STRING
                        && type != TransportSize.WSTRING)
                    {
                        throw new PlcInvalidFieldException(UnsupportedType(candidate));
                    }
                    return type;
                }
            }
            throw new PlcInvalidFieldException($"Unknown S7 data type '{name}'.");
        }

        private static string UnsupportedType(string name)
            => $"The data type {name} is not supported yet; supported: " +
               string.Join(", ", SupportedTypes) + ".";

        /// <summary>The address in the form plc4j gives it back, which parses to the same tag.</summary>
        public override string ToString()
        {
            var bit = DataType == TransportSize.BOOL ? $".{BitOffset}" : string.Empty;
            var type = StringLength.HasValue ? $"{DataType}({StringLength})" : DataType.ToString();
            switch (Area)
            {
                case MemoryArea.DATA_BLOCKS:
                    return $"%DB{DbNumber}.DB{ByteOffset}{bit}:{type}";
                case MemoryArea.INPUTS:
                    return $"%I{ByteOffset}{bit}:{type}";
                case MemoryArea.OUTPUTS:
                    return $"%Q{ByteOffset}{bit}:{type}";
                default:
                    return $"%M{ByteOffset}{bit}:{type}";
            }
        }
    }
}
