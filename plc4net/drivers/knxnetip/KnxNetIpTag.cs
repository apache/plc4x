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
using System.Text.RegularExpressions;
using org.apache.plc4net.drivers.knxnetip.readwrite.model;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.model;

namespace org.apache.plc4net.drivers.knxnetip
{
    /// <summary>
    /// A KNX group address, optionally with a datapoint-type hint.
    /// </summary>
    /// <remarks>
    /// Three notations, matching the Java <c>KnxNetIpTag</c>:
    /// <list type="bullet">
    ///   <item><c>1/2/3</c> - three level (main/middle/sub)</item>
    ///   <item><c>1/2</c> - two level (main/sub)</item>
    ///   <item><c>1</c> - free level (a single number)</item>
    /// </list>
    /// Each segment may be <c>*</c> to match any value in a subscription. An
    /// optional <c>:DPT&lt;n&gt;[.&lt;sub&gt;]</c> suffix (e.g. <c>1/2/3:DPT9.001</c>)
    /// tells the driver how to decode / encode the value when no project file is
    /// loaded; without it only raw byte reads and writes work.
    /// </remarks>
    public sealed class KnxNetIpTag : IPlcTag
    {
        // [0-9] rather than \d: in .NET \d also matches non-ASCII decimal digits
        // (Arabic-Indic, fullwidth, …), which then fail downstream in int.Parse. \z
        // rather than $ so a trailing newline is not silently accepted.
        private const string Dpt = @"(?::(?<dpt>DPT[0-9]+(?:\.[0-9]+)?))?";

        private static readonly Regex ThreeLevel = new Regex(
            @"^(?<main>[0-9]{1,2}|\*)/(?<middle>[0-9]|\*)/(?<sub>[0-9]{1,3}|\*)" + Dpt + @"\z",
            RegexOptions.IgnoreCase | RegexOptions.Compiled);

        private static readonly Regex TwoLevel = new Regex(
            @"^(?<main>[0-9]{1,2}|\*)/(?<sub>[0-9]{1,4}|\*)" + Dpt + @"\z",
            RegexOptions.IgnoreCase | RegexOptions.Compiled);

        private static readonly Regex FreeLevel = new Regex(
            @"^(?<sub>[0-9]{1,5}|\*)" + Dpt + @"\z",
            RegexOptions.IgnoreCase | RegexOptions.Compiled);

        internal KnxNetIpTag(int levels, string main, string middle, string sub, string dptId)
        {
            Levels = levels;
            MainGroup = main;
            MiddleGroup = middle;
            SubGroup = sub;
            DptId = dptId?.ToUpperInvariant();
        }

        public int Levels { get; }

        public string MainGroup { get; }

        public string MiddleGroup { get; }

        public string SubGroup { get; }

        /// <summary>The <c>DPT…</c> hint from the address, or null.</summary>
        public string DptId { get; }

        /// <summary>Whether any segment is the <c>*</c> wildcard.</summary>
        public bool HasWildcard =>
            MainGroup == "*" || MiddleGroup == "*" || SubGroup == "*";

        public string AddressString
        {
            get
            {
                var s = MainGroup;
                if (MiddleGroup != null)
                {
                    s += "/" + MiddleGroup;
                }
                if (Levels >= 2)
                {
                    s += "/" + SubGroup;
                }
                return s;
            }
        }

        public static bool Matches(string tagString) =>
            ThreeLevel.IsMatch(tagString) || TwoLevel.IsMatch(tagString) || FreeLevel.IsMatch(tagString);

        /// <summary>Parses a tag string, most-specific pattern first.</summary>
        public static KnxNetIpTag Parse(string tagString)
        {
            if (tagString == null)
            {
                throw new ArgumentNullException(nameof(tagString));
            }

            var m = ThreeLevel.Match(tagString);
            if (m.Success)
            {
                return new KnxNetIpTag(3, m.Groups["main"].Value, m.Groups["middle"].Value,
                    m.Groups["sub"].Value, Opt(m, "dpt"));
            }
            m = TwoLevel.Match(tagString);
            if (m.Success)
            {
                return new KnxNetIpTag(2, m.Groups["main"].Value, null,
                    m.Groups["sub"].Value, Opt(m, "dpt"));
            }
            m = FreeLevel.Match(tagString);
            if (m.Success)
            {
                return new KnxNetIpTag(1, m.Groups["sub"].Value, null,
                    m.Groups["sub"].Value, Opt(m, "dpt"));
            }
            throw new ArgumentException($"'{tagString}' is not a KNX group address.", nameof(tagString));
        }

        private static string Opt(Match m, string name) =>
            m.Groups[name].Success && m.Groups[name].Value.Length > 0 ? m.Groups[name].Value : null;

        /// <summary>
        /// Encodes the group address to its 2-byte wire form. The bit layout is
        /// fixed by the number of levels the connection was told to use:
        /// 3-level is 5/3/8 bits, 2-level 5/11, free-level 16. Every field is
        /// range-checked so an out-of-range segment is rejected here rather than
        /// silently truncated onto the wire (a write to the wrong device).
        /// </summary>
        public byte[] ToWireAddress(int connectionLevels)
        {
            if (HasWildcard)
            {
                throw new PlcInvalidFieldException(
                    $"'{AddressString}' has a wildcard and cannot be encoded to a single address.");
            }
            if (Levels != connectionLevels)
            {
                throw new PlcInvalidFieldException(
                    $"Tag '{AddressString}' uses {Levels} address level(s) but the connection is " +
                    $"configured for {connectionLevels}.");
            }

            var writeBuffer = new org.apache.plc4net.spi.generation.WriteBuffer();
            switch (connectionLevels)
            {
                case 3:
                    writeBuffer.WriteByte("main", 5, (byte) Field(MainGroup, 5, "main"));
                    writeBuffer.WriteByte("middle", 3, (byte) Field(MiddleGroup, 3, "middle"));
                    writeBuffer.WriteByte("sub", 8, (byte) Field(SubGroup, 8, "sub"));
                    break;
                case 2:
                    writeBuffer.WriteByte("main", 5, (byte) Field(MainGroup, 5, "main"));
                    writeBuffer.WriteUshort("sub", 11, Field(SubGroup, 11, "sub"));
                    break;
                case 1:
                    writeBuffer.WriteUshort("sub", 16, Field(SubGroup, 16, "sub"));
                    break;
                default:
                    throw new ArgumentOutOfRangeException(nameof(connectionLevels), connectionLevels,
                        "KNX group addresses have 1, 2 or 3 levels.");
            }
            return writeBuffer.GetBytes();
        }

        /// <summary>Whether an inbound group address matches this tag (wildcards included).</summary>
        public bool MatchesWireAddress(byte[] wire, int connectionLevels)
        {
            if (wire == null || wire.Length < 2)
            {
                return false;
            }
            var other = FromWire(wire, connectionLevels);
            if (other.Levels != Levels)
            {
                return false;
            }
            if (Levels == 3 && !SegmentMatches(MiddleGroup, other.MiddleGroup))
            {
                return false;
            }
            if (Levels >= 2 && !SegmentMatches(SubGroup, other.SubGroup))
            {
                return false;
            }
            return SegmentMatches(MainGroup, other.MainGroup);
        }

        // Compare group segments numerically (not as strings), so a tag written
        // with leading zeros - "1/2/03" - still matches the canonical inbound "1/2/3".
        private static bool SegmentMatches(string pattern, string actual)
        {
            if (pattern == "*")
            {
                return true;
            }
            return int.TryParse(pattern, NumberStyles.None, CultureInfo.InvariantCulture, out var p)
                   && int.TryParse(actual, NumberStyles.None, CultureInfo.InvariantCulture, out var a)
                   && p == a;
        }

        private static KnxNetIpTag FromWire(byte[] wire, int connectionLevels)
        {
            var raw = (wire[0] << 8) | wire[1];
            return connectionLevels switch
            {
                3 => new KnxNetIpTag(3,
                    ((raw >> 11) & 0x1F).ToString(CultureInfo.InvariantCulture),
                    ((raw >> 8) & 0x07).ToString(CultureInfo.InvariantCulture),
                    (raw & 0xFF).ToString(CultureInfo.InvariantCulture), null),
                2 => new KnxNetIpTag(2,
                    ((raw >> 11) & 0x1F).ToString(CultureInfo.InvariantCulture), null,
                    (raw & 0x7FF).ToString(CultureInfo.InvariantCulture), null),
                _ => new KnxNetIpTag(1,
                    raw.ToString(CultureInfo.InvariantCulture), null,
                    raw.ToString(CultureInfo.InvariantCulture), null),
            };
        }

        /// <summary>
        /// Parses one group segment and checks it fits in <paramref name="bits"/> bits,
        /// throwing a clear <see cref="PlcInvalidFieldException"/> rather than letting
        /// an over-range value truncate onto the wire or surface as an OverflowException.
        /// </summary>
        private ushort Field(string s, int bits, string name)
        {
            var max = (1 << bits) - 1;
            if (!int.TryParse(s, NumberStyles.None, CultureInfo.InvariantCulture, out var value)
                || value < 0 || value > max)
            {
                throw new PlcInvalidFieldException(
                    $"KNX {name}-group '{s}' in '{AddressString}' is out of range (0..{max}).");
            }
            return (ushort) value;
        }

        /// <summary>
        /// Resolves this tag's <c>:DPT…</c> hint to a generated <see cref="KnxDatapointType"/>.
        /// </summary>
        public KnxDatapointType? ResolveDatapointType() => ResolveDatapointType(DptId);

        /// <summary>
        /// Resolves a <c>DPT&lt;main&gt;[.&lt;sub&gt;]</c> hint to a generated
        /// <see cref="KnxDatapointType"/>. An exact <c>DPT&lt;main&gt;.&lt;sub&gt;</c> is
        /// matched against the enum's <c>id</c> (<c>DPST-main-sub</c>); a bare
        /// <c>DPT&lt;main&gt;</c> resolves to the first datapoint of that main type.
        /// Returns null when unresolved. Backed by a one-time index, not a per-call scan.
        /// </summary>
        public static KnxDatapointType? ResolveDatapointType(string dptId)
        {
            if (dptId == null)
            {
                return null;
            }

            var body = dptId.StartsWith("DPT", StringComparison.OrdinalIgnoreCase)
                ? dptId.Substring(3)
                : dptId;
            var dot = body.IndexOf('.');
            if (dot >= 0)
            {
                if (!int.TryParse(body.Substring(0, dot), NumberStyles.None, CultureInfo.InvariantCulture, out var mainNum)
                    || !int.TryParse(body.Substring(dot + 1), NumberStyles.None, CultureInfo.InvariantCulture, out var subNum))
                {
                    return null;
                }
                return IdToType.TryGetValue($"DPST-{mainNum}-{subNum}", out var exact)
                    ? exact
                    : (KnxDatapointType?) null;
            }

            if (!int.TryParse(body, NumberStyles.None, CultureInfo.InvariantCulture, out var main))
            {
                return null;
            }
            return MainToFirstType.TryGetValue(main, out var first) ? first : (KnxDatapointType?) null;
        }

        private static readonly Dictionary<string, KnxDatapointType> IdToType = BuildIdToType();
        private static readonly Dictionary<int, KnxDatapointType> MainToFirstType = BuildMainToFirstType();

        private static Dictionary<string, KnxDatapointType> BuildIdToType()
        {
            var map = new Dictionary<string, KnxDatapointType>(StringComparer.Ordinal);
            foreach (KnxDatapointType t in Enum.GetValues(typeof(KnxDatapointType)))
            {
                var id = t.GetId();
                if (!string.IsNullOrEmpty(id))
                {
                    map.TryAdd(id, t);
                }
            }
            return map;
        }

        private static Dictionary<int, KnxDatapointType> BuildMainToFirstType()
        {
            var map = new Dictionary<int, KnxDatapointType>();
            foreach (KnxDatapointType t in Enum.GetValues(typeof(KnxDatapointType)))
            {
                var id = t.GetId();
                if (string.IsNullOrEmpty(id) || !id.StartsWith("DPST-", StringComparison.Ordinal))
                {
                    continue;
                }
                var dash = id.IndexOf('-', 5);
                if (dash > 5
                    && int.TryParse(id.Substring(5, dash - 5), NumberStyles.None, CultureInfo.InvariantCulture, out var main))
                {
                    map.TryAdd(main, t);
                }
            }
            return map;
        }
    }
}
