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
using System.Globalization;
using System.Text.RegularExpressions;
using org.apache.plc4net.drivers.modbus.readwrite.model;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.model;

namespace org.apache.plc4net.drivers.modbus
{
    /// <summary>
    /// A Modbus tag: an address on a Modbus device.
    /// </summary>
    /// <remarks>
    /// <para>
    /// The syntax is the one plc4j uses, so the same address names the same register from
    /// either language: <c>{area}:{address}[:{TYPE}]</c>.
    /// </para>
    /// <list type="bullet">
    /// <item><c>coil:{address}</c> - a single bit, read and written (functions 0x01, 0x05).</item>
    /// <item><c>discrete-input:{address}</c> - a single bit, read only (0x02).</item>
    /// <item><c>holding-register:{address}[:{TYPE}]</c> - registers read and written with
    /// functions 0x03, 0x06 and 0x10; <c>INT</c> unless stated.</item>
    /// <item><c>input-register:{address}[:{TYPE}]</c> - registers read only with function 0x04.</item>
    /// </list>
    /// <para>
    /// The address is the register number as a Modbus user counts it, from 1: the request
    /// carries <c>address - 1</c>, so <c>holding-register:1</c> is the register at wire
    /// address 0.
    /// </para>
    /// <para>
    /// The optional plc4j-compatible <c>{byte-order:'...'}</c> configuration controls the byte
    /// layout of a register value. Array selections, fixed strings, extended registers and short
    /// forms such as <c>40001</c> and <c>4x00001</c> are supported.
    /// </para>
    /// </remarks>
    public sealed class ModbusTag : IPlcTag
    {
        public enum TagType
        {
            Coil,
            HoldingRegister,
            DiscreteInput,
            InputRegister,
            ExtendedRegister
        }

        /// <summary>The highest register number: wire address 65535.</summary>
        public const int MaxAddress = 65536;

        // [0-9], not \d: in .NET \d matches every Unicode decimal digit, so a
        // non-ASCII address would pass the pattern and blow up int.Parse below
        // instead of failing here as a PlcInvalidFieldException.
        private static readonly Regex AddressPattern = new Regex(
            @"^(?<area>coil|discrete-input|holding-register|input-register|extended-register):(?<address>[0-9]{1,9})" +
            @"(?<selection>\[(?<start>[0-9]+)(\.\.(?<end>[0-9]+))?\])?" +
            @"(:(?<datatype>[A-Za-z_]+)(\((?<stringlength>[0-9]+)\))?)?" +
            @"(\{(?<config>[^{}]+)\})?$",
            RegexOptions.CultureInvariant | RegexOptions.IgnoreCase);

        private static readonly Regex ShortAddressPattern = new Regex(
            @"^(?<prefix>[01346])[xX]?(?<address>[0-9]{4,5})(?<tail>.*)$",
            RegexOptions.CultureInvariant);

        public ModbusTag(TagType type, int address, ModbusDataType dataType,
            ModbusByteOrder byteOrder = ModbusByteOrder.BigEndian, ushort quantity = 1,
            ushort stringLength = 0, bool explicitRange = false, byte? unitIdentifier = null)
        {
            // File numbers are 1..65535 and each file contains 10000 logical records.
            const int maximumExtendedAddress = 655349999;
            var maximumAddress = type == TagType.ExtendedRegister ? maximumExtendedAddress : MaxAddress;
            if (address < 1 || address > maximumAddress)
            {
                throw new PlcInvalidFieldException(
                    $"The Modbus address must be between 1 and {maximumAddress}, but was {address}.");
            }
            if (!IsSupported(type, dataType))
            {
                throw new PlcInvalidFieldException(
                    $"The data type {dataType} is not supported for {AreaName(type)} tags yet; " +
                    $"supported: {SupportedTypes(type)}.");
            }
            if (quantity == 0)
            {
                throw new PlcInvalidFieldException("The Modbus quantity must be greater than zero.");
            }
            if (IsString(dataType) != (stringLength > 0))
            {
                throw new PlcInvalidFieldException(IsString(dataType)
                    ? $"{dataType} requires a positive string length, for example {dataType}(20)."
                    : $"A string length is not valid for {dataType}.");
            }

            Type = type;
            Address = address;
            DataType = dataType;
            ByteOrder = byteOrder;
            Quantity = quantity;
            StringLength = stringLength;
            ExplicitRange = explicitRange;
            UnitIdentifier = unitIdentifier;
            if (RegisterCount > (IsBit ? 2000 : 125))
            {
                throw new PlcInvalidFieldException(
                    $"The Modbus selection occupies {RegisterCount} {(IsBit ? "bits" : "registers")}; " +
                    $"at most {(IsBit ? 2000 : 125)} fit in one request.");
            }
            if (type != TagType.ExtendedRegister && (long)WireAddress + RegisterCount > 65536)
            {
                throw new PlcInvalidFieldException("The Modbus selection runs past the end of the address space.");
            }
            if (type == TagType.ExtendedRegister && (long)Address + RegisterCount - 1 > maximumExtendedAddress)
                throw new PlcInvalidFieldException("The Modbus selection runs past the final file record.");
        }

        public TagType Type { get; }

        /// <summary>The register number, counted from 1.</summary>
        public int Address { get; }

        /// <summary>The address the request carries: <see cref="Address"/> - 1.</summary>
        public ushort WireAddress => checked((ushort)(Type == TagType.ExtendedRegister ? Address : Address - 1));

        public ModbusDataType DataType { get; }

        public ModbusByteOrder ByteOrder { get; }

        public ushort Quantity { get; }

        public ushort StringLength { get; }

        public bool ExplicitRange { get; }

        public byte? UnitIdentifier { get; }

        /// <summary>The number of 16-bit registers occupied by one value.</summary>
        public ushort RegisterCount
        {
            get
            {
                if (IsBit) return Quantity;
                var bytesPerValue = IsString(DataType)
                    ? (long)StringLength * (DataType == ModbusDataType.WSTRING ? 2 : 1)
                    : NaturalSize(DataType);
                var registers = Math.Max(1, (bytesPerValue * Quantity + 1) / 2);
                if (registers > ushort.MaxValue)
                {
                    throw new PlcInvalidFieldException(
                        $"The Modbus selection needs {registers} registers, which does not fit a request.");
                }
                return (ushort)registers;
            }
        }

        /// <summary>Whether the tag is a single bit (a coil or a discrete input) rather than a register.</summary>
        public bool IsBit => Type == TagType.Coil || Type == TagType.DiscreteInput;

        /// <summary>
        /// Parses a Modbus tag address such as <c>coil:1</c> or <c>holding-register:10:UINT</c>.
        /// </summary>
        /// <exception cref="PlcInvalidFieldException">The address is not a supported Modbus tag address.</exception>
        public static ModbusTag Parse(string tagAddress)
            => Parse(tagAddress, ModbusByteOrder.BigEndian);

        internal static ModbusTag Parse(string tagAddress, ModbusByteOrder defaultByteOrder)
        {
            if (string.IsNullOrWhiteSpace(tagAddress))
            {
                throw new PlcInvalidFieldException("The Modbus tag address must not be empty.");
            }

            var addressText = NormalizeShortAddress(tagAddress.Trim());
            var match = AddressPattern.Match(addressText);
            if (!match.Success)
            {
                throw new PlcInvalidFieldException(
                    $"'{tagAddress}' is not a supported Modbus tag address. Expected " +
                    "{area}:{address}[selection]:{TYPE}, for example holding-register:1[0..3]:INT.");
            }

            var type = ParseArea(match.Groups["area"].Value);
            var address = int.Parse(match.Groups["address"].Value, CultureInfo.InvariantCulture);

            ModbusDataType dataType;
            var typeGroup = match.Groups["datatype"];
            if (typeGroup.Success)
            {
                if (!TryParseDataType(typeGroup.Value, out dataType))
                {
                    throw new PlcInvalidFieldException($"Unknown Modbus data type '{typeGroup.Value}'.");
                }
            }
            else
            {
                dataType = type == TagType.Coil || type == TagType.DiscreteInput
                    ? ModbusDataType.BOOL
                    : ModbusDataType.INT;
            }

            var stringLength = (ushort)0;
            if (match.Groups["stringlength"].Success)
            {
                if (!ushort.TryParse(match.Groups["stringlength"].Value, NumberStyles.None,
                    CultureInfo.InvariantCulture, out stringLength) || stringLength == 0)
                {
                    throw new PlcInvalidFieldException("The Modbus string length must be between 1 and 65535.");
                }
            }

            var quantity = (ushort)1;
            var explicitRange = false;
            if (match.Groups["selection"].Success)
            {
                var start = int.Parse(match.Groups["start"].Value, CultureInfo.InvariantCulture);
                var end = match.Groups["end"].Success
                    ? int.Parse(match.Groups["end"].Value, CultureInfo.InvariantCulture)
                    : start;
                if (end < start || end - start + 1 > ushort.MaxValue)
                {
                    throw new PlcInvalidFieldException("The Modbus selection must be an ascending 16-bit range.");
                }
                explicitRange = match.Groups["end"].Success;
                quantity = checked((ushort)(end - start + 1));
                address = checked(address + ElementAddressOffset(type, dataType, stringLength, start));
            }

            var byteOrder = defaultByteOrder;
            byte? unitIdentifier = null;
            if (match.Groups["config"].Success)
            {
                foreach (var entry in match.Groups["config"].Value.Split(','))
                {
                    var pair = entry.Split(new[] { ':' }, 2);
                    if (pair.Length != 2) throw new PlcInvalidFieldException($"Invalid Modbus tag configuration '{entry.Trim()}'.");
                    var key = pair[0].Trim().ToLowerInvariant();
                    var value = pair[1].Trim().Trim('\'', '"');
                    if (key == "byte-order")
                    {
                        if (!TryParseByteOrder(value, out byteOrder))
                            throw new PlcInvalidFieldException($"Unknown Modbus byte order '{value}'.");
                    }
                    else if (key == "unit-id")
                    {
                        if (!byte.TryParse(value, NumberStyles.None, CultureInfo.InvariantCulture, out var parsed))
                            throw new PlcInvalidFieldException($"The Modbus unit-id must be between 0 and 255, but was '{value}'.");
                        unitIdentifier = parsed;
                    }
                    else throw new PlcInvalidFieldException($"Unknown Modbus tag configuration '{key}'.");
                }
            }

            return new ModbusTag(type, address, dataType, byteOrder, quantity, stringLength, explicitRange, unitIdentifier);
        }

        public override string ToString()
        {
            var selection = ExplicitRange ? $"[0..{Quantity - 1}]" : string.Empty;
            var type = IsString(DataType) ? $"{DataType}({StringLength})" : DataType.ToString();
            var settings = new System.Collections.Generic.List<string>();
            if (UnitIdentifier.HasValue) settings.Add($"unit-id:{UnitIdentifier.Value}");
            if (ByteOrder != ModbusByteOrder.BigEndian)
                settings.Add($"byte-order:'{ByteOrderName(ByteOrder)}'");
            var config = settings.Count == 0 ? string.Empty : "{" + string.Join(", ", settings) + "}";
            return $"{AreaName(Type)}:{Address}{selection}:{type}{config}";
        }

        internal static TagType ParseArea(string area)
        {
            switch (area.ToLowerInvariant())
            {
                case "coil":
                    return TagType.Coil;
                case "discrete-input":
                    return TagType.DiscreteInput;
                case "holding-register":
                    return TagType.HoldingRegister;
                case "input-register":
                    return TagType.InputRegister;
                case "extended-register":
                    return TagType.ExtendedRegister;
                default:
                    // The pattern lets only the four names through. A name that gets past it
                    // later must fail here: reading it as an input register would address
                    // other memory than the tag says.
                    throw new PlcInvalidFieldException($"Unknown Modbus area '{area}'.");
            }
        }

        private static string AreaName(TagType type)
        {
            switch (type)
            {
                case TagType.Coil:
                    return "coil";
                case TagType.DiscreteInput:
                    return "discrete-input";
                case TagType.HoldingRegister:
                    return "holding-register";
                case TagType.InputRegister:
                    return "input-register";
                default:
                    return "extended-register";
            }
        }

        private static bool TryParseDataType(string name, out ModbusDataType dataType)
        {
            // Enum.TryParse would also accept a number; only the names are data types.
            foreach (var candidate in Enum.GetNames(typeof(ModbusDataType)))
            {
                if (string.Equals(candidate, name, StringComparison.OrdinalIgnoreCase))
                {
                    dataType = (ModbusDataType)Enum.Parse(typeof(ModbusDataType), candidate);
                    return true;
                }
            }
            dataType = default;
            return false;
        }

        internal static bool TryParseByteOrder(string name, out ModbusByteOrder byteOrder)
        {
            switch (name.ToUpperInvariant())
            {
                case "BIG_ENDIAN": byteOrder = ModbusByteOrder.BigEndian; return true;
                case "LITTLE_ENDIAN": byteOrder = ModbusByteOrder.LittleEndian; return true;
                case "BIG_ENDIAN_BYTE_SWAP": byteOrder = ModbusByteOrder.BigEndianByteSwap; return true;
                case "LITTLE_ENDIAN_BYTE_SWAP": byteOrder = ModbusByteOrder.LittleEndianByteSwap; return true;
                default: byteOrder = default; return false;
            }
        }

        private static string ByteOrderName(ModbusByteOrder byteOrder)
            => Regex.Replace(byteOrder.ToString(), "([a-z])([A-Z])", "$1_$2").ToUpperInvariant();

        private static bool IsSupported(TagType type, ModbusDataType dataType)
        {
            if (type == TagType.Coil || type == TagType.DiscreteInput)
            {
                return dataType == ModbusDataType.BOOL;
            }
            return dataType == ModbusDataType.BYTE
                   || dataType == ModbusDataType.WORD
                   || dataType == ModbusDataType.DWORD
                   || dataType == ModbusDataType.LWORD
                   || dataType == ModbusDataType.SINT
                   || dataType == ModbusDataType.INT
                   || dataType == ModbusDataType.DINT
                   || dataType == ModbusDataType.LINT
                   || dataType == ModbusDataType.USINT
                   || dataType == ModbusDataType.UINT
                   || dataType == ModbusDataType.UDINT
                   || dataType == ModbusDataType.ULINT
                   || dataType == ModbusDataType.REAL
                   || dataType == ModbusDataType.LREAL
                   || dataType == ModbusDataType.STRING
                   || dataType == ModbusDataType.WSTRING;
        }

        private static string SupportedTypes(TagType type)
            => type == TagType.Coil || type == TagType.DiscreteInput
                ? "BOOL"
                : "BYTE, WORD, DWORD, LWORD, SINT, INT, DINT, LINT, USINT, UINT, UDINT, ULINT, REAL, LREAL, STRING, WSTRING";

        private static bool IsString(ModbusDataType dataType)
            => dataType == ModbusDataType.STRING || dataType == ModbusDataType.WSTRING;

        private static int NaturalSize(ModbusDataType dataType)
            => dataType == ModbusDataType.BYTE || dataType == ModbusDataType.SINT || dataType == ModbusDataType.USINT
                ? 1 : dataType.GetDataTypeSize();

        private static int ElementAddressOffset(
            TagType type, ModbusDataType dataType, ushort stringLength, int elementOffset)
        {
            if (type == TagType.Coil || type == TagType.DiscreteInput) return elementOffset;
            var bytes = IsString(dataType)
                ? stringLength * (dataType == ModbusDataType.WSTRING ? 2 : 1)
                : NaturalSize(dataType);
            var bitOffset = checked(elementOffset * bytes * 8);
            if (bitOffset % 16 != 0)
            {
                throw new PlcInvalidFieldException(
                    $"The selection starts {bitOffset} bits into the value and not on a register boundary.");
            }
            return bitOffset / 16;
        }

        private static string NormalizeShortAddress(string address)
        {
            var match = ShortAddressPattern.Match(address);
            if (!match.Success) return address;
            var area = match.Groups["prefix"].Value switch
            {
                "0" => "coil",
                "1" => "discrete-input",
                "3" => "input-register",
                "4" => "holding-register",
                _ => "extended-register"
            };
            var number = int.Parse(match.Groups["address"].Value, CultureInfo.InvariantCulture);
            return $"{area}:{number}{match.Groups["tail"].Value}";
        }
    }
}
