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
    /// <item><c>holding-register:{address}[:INT|UINT|WORD]</c> - a 16-bit register, read and
    /// written (0x03, 0x06); <c>INT</c> unless stated.</item>
    /// <item><c>input-register:{address}[:INT|UINT|WORD]</c> - a 16-bit register, read only (0x04).</item>
    /// </list>
    /// <para>
    /// The address is the register number as a Modbus user counts it, from 1: the request
    /// carries <c>address - 1</c>, so <c>holding-register:1</c> is the register at wire
    /// address 0.
    /// </para>
    /// <para>
    /// Array selections, string lengths, tag configuration, extended registers, the short
    /// forms (<c>40001</c>, <c>4x00001</c>) and the wider data types are accepted by plc4j
    /// but not supported here yet; they are rejected rather than read as something else.
    /// </para>
    /// </remarks>
    public sealed class ModbusTag : IPlcTag
    {
        public enum TagType
        {
            Coil,
            HoldingRegister,
            DiscreteInput,
            InputRegister
        }

        /// <summary>The highest register number: wire address 65535.</summary>
        public const int MaxAddress = 65536;

        private static readonly Regex AddressPattern = new Regex(
            @"^(?<area>coil|discrete-input|holding-register|input-register):(?<address>\d{1,9})(:(?<datatype>[A-Za-z_]+))?$",
            RegexOptions.CultureInvariant | RegexOptions.IgnoreCase);

        public ModbusTag(TagType type, int address, ModbusDataType dataType)
        {
            if (address < 1 || address > MaxAddress)
            {
                throw new PlcInvalidFieldException(
                    $"The Modbus address must be between 1 and {MaxAddress}, but was {address}.");
            }
            if (!IsSupported(type, dataType))
            {
                throw new PlcInvalidFieldException(
                    $"The data type {dataType} is not supported for {AreaName(type)} tags yet; " +
                    $"supported: {SupportedTypes(type)}.");
            }

            Type = type;
            Address = address;
            DataType = dataType;
        }

        public TagType Type { get; }

        /// <summary>The register number, counted from 1.</summary>
        public int Address { get; }

        /// <summary>The address the request carries: <see cref="Address"/> - 1.</summary>
        public ushort WireAddress => (ushort)(Address - 1);

        public ModbusDataType DataType { get; }

        /// <summary>Whether the tag is a single bit (a coil or a discrete input) rather than a register.</summary>
        public bool IsBit => Type == TagType.Coil || Type == TagType.DiscreteInput;

        /// <summary>
        /// Parses a Modbus tag address such as <c>coil:1</c> or <c>holding-register:10:UINT</c>.
        /// </summary>
        /// <exception cref="PlcInvalidFieldException">The address is not a supported Modbus tag address.</exception>
        public static ModbusTag Parse(string tagAddress)
        {
            if (string.IsNullOrWhiteSpace(tagAddress))
            {
                throw new PlcInvalidFieldException("The Modbus tag address must not be empty.");
            }

            var match = AddressPattern.Match(tagAddress.Trim());
            if (!match.Success)
            {
                throw new PlcInvalidFieldException(
                    $"'{tagAddress}' is not a supported Modbus tag address. Expected " +
                    "{area}:{address}[:{TYPE}], for example holding-register:1:INT; the areas are coil, " +
                    "discrete-input, holding-register and input-register, and the address counts from 1. " +
                    "Array selections, string lengths, tag configuration, extended registers and short " +
                    "addresses such as 40001 are not supported yet.");
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

            return new ModbusTag(type, address, dataType);
        }

        public override string ToString() => $"{AreaName(Type)}:{Address}:{DataType}";

        private static TagType ParseArea(string area)
        {
            switch (area.ToLowerInvariant())
            {
                case "coil":
                    return TagType.Coil;
                case "discrete-input":
                    return TagType.DiscreteInput;
                case "holding-register":
                    return TagType.HoldingRegister;
                default:
                    return TagType.InputRegister;
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
                default:
                    return "input-register";
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

        private static bool IsSupported(TagType type, ModbusDataType dataType)
        {
            if (type == TagType.Coil || type == TagType.DiscreteInput)
            {
                return dataType == ModbusDataType.BOOL;
            }
            return dataType == ModbusDataType.INT
                   || dataType == ModbusDataType.UINT
                   || dataType == ModbusDataType.WORD;
        }

        private static string SupportedTypes(TagType type)
            => type == TagType.Coil || type == TagType.DiscreteInput ? "BOOL" : "INT, UINT, WORD";
    }
}