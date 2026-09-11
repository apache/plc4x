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
using org.apache.plc4net.model;

namespace org.apache.plc4net.drivers.modbus
{
    /// <summary>
    /// A Modbus tag — an address on a Modbus device.
    /// </summary>
    /// <remarks>
    /// Address syntax (case-insensitive):
    ///   coil:{address}            — read/write single bit (0x01/0x05)
    ///   holding:{address}         — read/write 16-bit register (0x03/0x06/0x10)
    ///   discrete:{address}        — read-only single bit (0x02)
    ///   input:{address}           — read-only 16-bit register (0x04)
    /// </remarks>
    public class ModbusTag : IPlcTag
    {
        public enum TagType
        {
            Coil,
            HoldingRegister,
            DiscreteInput,
            InputRegister
        }

        public ModbusTag(TagType type, ushort address)
        {
            Type = type;
            Address = address;
        }

        public TagType Type { get; }
        public ushort Address { get; }

        /// <summary>Parses a Modbus address string like "coil:0" or "holding:10".</summary>
        public static ModbusTag Parse(string tagAddress)
        {
            if (string.IsNullOrWhiteSpace(tagAddress))
                throw new ModbusDriverException("Tag address must not be empty.");

            var colon = tagAddress.IndexOf(':');
            if (colon < 0)
                throw new ModbusDriverException(
                    $"Modbus tag address '{tagAddress}' is missing a type prefix " +
                    "(coil:, holding:, discrete:, input:).");

            var typePart = tagAddress.Substring(0, colon).Trim();
            var addrPart = tagAddress.Substring(colon + 1).Trim();

            if (!ushort.TryParse(addrPart, NumberStyles.Integer, CultureInfo.InvariantCulture, out var address))
                throw new ModbusDriverException(
                    $"Modbus tag address '{addrPart}' is not a valid unsigned 16-bit address.");

            var type = typePart.ToLowerInvariant() switch
            {
                "coil" => TagType.Coil,
                "holding" => TagType.HoldingRegister,
                "discrete" => TagType.DiscreteInput,
                "input" => TagType.InputRegister,
                _ => throw new ModbusDriverException(
                    $"Unknown Modbus tag type '{typePart}'. " +
                    "Expected: coil, holding, discrete, or input.")
            };

            return new ModbusTag(type, address);
        }

        public override string ToString()
        {
            var prefix = Type switch
            {
                TagType.Coil => "coil",
                TagType.HoldingRegister => "holding",
                TagType.DiscreteInput => "discrete",
                TagType.InputRegister => "input",
                _ => "?"
            };
            return $"{prefix}:{Address}";
        }
    }
}
