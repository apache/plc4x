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
using org.apache.plc4net.drivers.knxnetip.readwrite.model;
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
        private const string Dpt = @"(?::(?<dpt>DPT\d+(?:\.\d+)?))?";

        private static readonly Regex ThreeLevel = new Regex(
            @"^(?<main>\d{1,2}|\*)/(?<middle>\d|\*)/(?<sub>\d{1,3}|\*)" + Dpt + "$",
            RegexOptions.IgnoreCase | RegexOptions.Compiled);

        private static readonly Regex TwoLevel = new Regex(
            @"^(?<main>\d{1,2}|\*)/(?<sub>\d{1,4}|\*)" + Dpt + "$",
            RegexOptions.IgnoreCase | RegexOptions.Compiled);

        private static readonly Regex FreeLevel = new Regex(
            @"^(?<sub>\d{1,5}|\*)" + Dpt + "$",
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
        /// 3-level is 5/3/8 bits, 2-level 5/11, free-level 16.
        /// </summary>
        public byte[] ToWireAddress(int connectionLevels)
        {
            if (HasWildcard)
            {
                throw new InvalidOperationException(
                    $"'{AddressString}' has a wildcard and cannot be encoded to a single address.");
            }

            var writeBuffer = new org.apache.plc4net.spi.generation.WriteBuffer();
            switch (connectionLevels)
            {
                case 3:
                    writeBuffer.WriteByte("main", 5, Byte(MainGroup));
                    writeBuffer.WriteByte("middle", 3, Byte(MiddleGroup));
                    writeBuffer.WriteByte("sub", 8, Byte(SubGroup));
                    break;
                case 2:
                    writeBuffer.WriteByte("main", 5, Byte(MainGroup));
                    writeBuffer.WriteUshort("sub", 11, Ushort(SubGroup));
                    break;
                case 1:
                    writeBuffer.WriteUshort("sub", 16, Ushort(SubGroup));
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
            var other = FromWire(wire, connectionLevels);
            if (other.Levels != Levels)
            {
                return false;
            }
            if (Levels == 3 && MiddleGroup != "*" && MiddleGroup != other.MiddleGroup)
            {
                return false;
            }
            if (Levels >= 2 && SubGroup != "*" && SubGroup != other.SubGroup)
            {
                return false;
            }
            return MainGroup == "*" || MainGroup == other.MainGroup;
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

        private static byte Byte(string s) => byte.Parse(s, CultureInfo.InvariantCulture);

        private static ushort Ushort(string s) => ushort.Parse(s, CultureInfo.InvariantCulture);

        /// <summary>
        /// Resolves the <c>:DPT…</c> hint to a generated <see cref="KnxDatapointType"/>.
        /// An exact <c>DPT&lt;main&gt;.&lt;sub&gt;</c> is matched against the enum's
        /// <c>id</c> (<c>DPST-main-sub</c>); a bare <c>DPT&lt;main&gt;</c> resolves to
        /// the first datapoint of that main type. Returns null when unresolved.
        /// </summary>
        public KnxDatapointType? ResolveDatapointType()
        {
            if (DptId == null)
            {
                return null;
            }

            var body = DptId.StartsWith("DPT", StringComparison.OrdinalIgnoreCase)
                ? DptId.Substring(3)
                : DptId;
            var dot = body.IndexOf('.');
            string wantId;
            int mainNum;
            if (dot >= 0)
            {
                if (!int.TryParse(body.Substring(0, dot), out mainNum)
                    || !int.TryParse(body.Substring(dot + 1), out var subNum))
                {
                    return null;
                }
                wantId = $"DPST-{mainNum}-{subNum}";
            }
            else
            {
                if (!int.TryParse(body, out mainNum))
                {
                    return null;
                }
                wantId = null;
            }

            KnxDatapointType? firstOfMain = null;
            foreach (KnxDatapointType candidate in Enum.GetValues(typeof(KnxDatapointType)))
            {
                var id = candidate.GetId();
                if (wantId != null)
                {
                    if (id == wantId)
                    {
                        return candidate;
                    }
                }
                else if (firstOfMain == null && id != null
                         && id.StartsWith($"DPST-{mainNum}-", StringComparison.Ordinal))
                {
                    firstOfMain = candidate;
                }
            }
            return firstOfMain;
        }
    }
}
