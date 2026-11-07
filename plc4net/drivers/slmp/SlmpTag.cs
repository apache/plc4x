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
using org.apache.plc4net.drivers.slmp.readwrite.model;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.model;

namespace org.apache.plc4net.drivers.slmp
{
    public enum SlmpDataType { WORD, INT, UINT, DINT, UDINT, REAL }

    public static class SlmpDataTypeExtensions
    {
        public static int GetWordsPerElement(this SlmpDataType type) => type switch
        {
            SlmpDataType.WORD or SlmpDataType.INT or SlmpDataType.UINT => 1,
            SlmpDataType.DINT or SlmpDataType.UDINT or SlmpDataType.REAL => 2,
            _ => throw new PlcInvalidFieldException($"Unsupported SLMP data type {type}.")
        };
    }

    /// <summary>A word-device address with an optional zero-based index or inclusive range.</summary>
    public sealed class SlmpTag : IPlcTag
    {
        // Same conservative single-frame policy as the current Java SLMP tag.
        public const int MaxPoints = 960;
        public const int MaxDeviceNumber = 0xFFFFFF;
        private static readonly Regex Pattern = new Regex(
            @"\A(?<device>[DWR])(?<address>(?:0[xX])?[0-9A-Fa-f]+)(?:\[(?<lower>[0-9]+)(?:\.\.(?<upper>[0-9]+))?\])?(?::(?<type>[A-Za-z]+))?\z",
            RegexOptions.CultureInvariant | RegexOptions.IgnoreCase);

        public SlmpTag(SlmpDeviceCode deviceCode, int deviceNumber, SlmpDataType dataType,
            int quantity = 1, bool? explicitRange = null)
        {
            if (deviceCode != SlmpDeviceCode.D && deviceCode != SlmpDeviceCode.W && deviceCode != SlmpDeviceCode.R)
                throw new PlcInvalidFieldException("Only SLMP word devices D, W and R are supported.");
            var words = dataType.GetWordsPerElement();
            if (quantity < 1 || quantity > MaxPoints / words)
                throw new PlcInvalidFieldException($"SLMP selections must contain between 1 and {MaxPoints} words.");
            var points = quantity * words;
            if (deviceNumber < 0 || (long)deviceNumber + points - 1 > MaxDeviceNumber)
                throw new PlcInvalidFieldException("SLMP selection exceeds the 24-bit device-address range.");
            if (quantity > 1 && explicitRange == false)
                throw new PlcInvalidFieldException("A multiple-element selection must be a range.");
            DeviceCode = deviceCode;
            DeviceNumber = deviceNumber;
            DataType = dataType;
            Quantity = quantity;
            ExplicitRange = explicitRange ?? quantity > 1;
        }

        public SlmpDeviceCode DeviceCode { get; }
        public int DeviceNumber { get; }
        public SlmpDataType DataType { get; }
        public int Quantity { get; }
        public bool ExplicitRange { get; }
        public int NumberOfPoints => Quantity * DataType.GetWordsPerElement();

        public static SlmpTag Parse(string tagAddress)
        {
            if (string.IsNullOrWhiteSpace(tagAddress) || tagAddress.Length > 128)
                throw new PlcInvalidFieldException("The SLMP tag address is empty or too long.");
            var match = Pattern.Match(tagAddress.Trim());
            if (!match.Success)
                throw new PlcInvalidFieldException("Expected {D|W|R}{address}[index or lower..upper]:{TYPE}.");
            var device = match.Groups["device"].Value.ToUpperInvariant() switch
            {
                "D" => SlmpDeviceCode.D,
                "W" => SlmpDeviceCode.W,
                _ => SlmpDeviceCode.R
            };
            var address = match.Groups["address"].Value;
            if (address.StartsWith("0x", StringComparison.OrdinalIgnoreCase))
            {
                if (device != SlmpDeviceCode.W)
                    throw new PlcInvalidFieldException("The 0x prefix is valid only for W devices.");
                address = address.Substring(2);
            }
            var style = device == SlmpDeviceCode.W ? NumberStyles.AllowHexSpecifier : NumberStyles.None;
            if (!uint.TryParse(address, style, CultureInfo.InvariantCulture, out var number) || number > MaxDeviceNumber)
                throw new PlcInvalidFieldException("Invalid SLMP device number.");
            var typeName = match.Groups["type"].Success ? match.Groups["type"].Value.ToUpperInvariant() : "WORD";
            if (!Enum.TryParse<SlmpDataType>(typeName, out var type) || !Enum.IsDefined(type))
                throw new PlcInvalidFieldException($"Unsupported SLMP data type '{typeName}'.");
            var lower = 0u;
            var upper = 0u;
            if (match.Groups["lower"].Success &&
                !uint.TryParse(match.Groups["lower"].Value, NumberStyles.None, CultureInfo.InvariantCulture, out lower))
                throw new PlcInvalidFieldException("SLMP selection index is too large.");
            var range = match.Groups["upper"].Success;
            if (range && (!uint.TryParse(match.Groups["upper"].Value, NumberStyles.None, CultureInfo.InvariantCulture, out upper) || upper < lower))
                throw new PlcInvalidFieldException("Invalid SLMP selection range.");
            var quantity = range ? (long)upper - lower + 1 : 1;
            var resolved = (long)number + (long)lower * type.GetWordsPerElement();
            if (quantity > MaxPoints || resolved > MaxDeviceNumber)
                throw new PlcInvalidFieldException("SLMP selection exceeds the supported address or word count.");
            return new SlmpTag(device, (int)resolved, type, (int)quantity, range);
        }

        public override string ToString()
        {
            var address = DeviceCode == SlmpDeviceCode.W
                ? "0x" + DeviceNumber.ToString("X", CultureInfo.InvariantCulture)
                : DeviceNumber.ToString(CultureInfo.InvariantCulture);
            var selection = ExplicitRange ? "[0.." + (Quantity - 1).ToString(CultureInfo.InvariantCulture) + "]" : "";
            return $"{DeviceCode}{address}{selection}:{DataType}";
        }
    }
}
