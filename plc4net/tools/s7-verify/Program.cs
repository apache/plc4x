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
using System.Threading.Tasks;
using org.apache.plc4net.api;
using org.apache.plc4net.api.value;
using org.apache.plc4net.drivers.s7;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.types;

namespace org.apache.plc4net.tools.s7verify
{
    /// <summary>
    /// End-to-end S7 hardware check. Drives the public driver API exactly as a
    /// NuGet consumer would - <c>new S7Driver(...).Connect("s7://host?...")</c> -
    /// then reads every scalar type from a data block, round-trips a write, and
    /// exercises an error path. Prints a Markdown report to stdout.
    ///
    ///   s7-verify &lt;host&gt; [--rack N] [--slot N] [--db N]
    ///             [--device-group PG_OR_PC|OS|OTHERS] [--remote-tsap 0xNNNN]
    ///
    /// Defaults: rack 0, slot 1, db 100. Build the data block per
    /// docs/s7-hardware-verification.md before running.
    /// </summary>
    public static class Program
    {
        public static async Task<int> Main(string[] args)
        {
            if (args.Length == 0 || args[0] is "-h" or "--help")
            {
                Console.Error.WriteLine(
                    "usage: s7-verify <host> [--rack N] [--slot N] [--db N] " +
                    "[--device-group PG_OR_PC|OS|OTHERS] [--remote-tsap 0xNNNN] [--read <address>]");
                return 2;
            }

            var host = args[0];
            var opts = ParseOptions(args);
            var db = opts.TryGetValue("db", out var dbRaw)
                ? int.Parse(dbRaw, CultureInfo.InvariantCulture) : 100;

            var query = new List<string>
            {
                $"remote-rack={Get(opts, "rack", "0")}",
                $"remote-slot={Get(opts, "slot", "1")}",
                "request-timeout=5000",
            };
            if (opts.TryGetValue("device-group", out var dg)) query.Add($"remote-device-group={dg}");
            if (opts.TryGetValue("remote-tsap", out var tsap)) query.Add($"remote-tsap={tsap}");
            var connectionString = $"s7://{host}?{string.Join("&", query)}";

            // --read <address>: connect, read one tag, print the outcome, exit. A
            // focused probe for troubleshooting a single address without running the
            // whole suite or needing DB100 to exist.
            if (opts.TryGetValue("read", out var probeAddress))
            {
                return await ProbeSingleRead(connectionString, probeAddress);
            }

            var report = new Report();
            report.Line("# S7 Hardware Verification Report");
            report.Line();
            report.Line($"- **Date**: {DateTime.Now:yyyy-MM-dd HH:mm}");
            report.Line($"- **Connection string**: `{connectionString}`");
            report.Line($"- **Data block**: DB{db}");
            report.Line();

            IPlcConnection? connection = null;
            try
            {
                // ── connect + Setup Communication ──
                var driver = new S7Driver(new DefaultTransportManager());
                connection = driver.Connect(connectionString);
                var s7 = (S7Connection)connection;
                report.Pass("Connect", $"COTP + Setup Communication ok, negotiated PDU length {s7.NegotiatedPduLength} bytes");

                var reader = (PlcReader)connection;
                var writer = (PlcWriter)connection;

                // ── reads ──
                await ReadScalar(reader, connection, report, "BOOL", $"%DB{db}.DBX0.0",
                    v => v.GetBool() ? "true" : "false", "true");
                await ReadScalar(reader, connection, report, "BYTE", $"%DB{db}.DBB1",
                    v => $"0x{v.GetByte():X2}", "0xA5");
                await ReadScalar(reader, connection, report, "INT", $"%DB{db}.DBW2",
                    v => unchecked((short)v.GetUshort()).ToString(CultureInfo.InvariantCulture), "-12345");
                await ReadScalar(reader, connection, report, "DINT", $"%DB{db}.DBD4",
                    v => unchecked((int)v.GetUint()).ToString(CultureInfo.InvariantCulture), "-1000000");
                await ReadScalar(reader, connection, report, "REAL", $"%DB{db}.DBD8",
                    v => BitConverter.Int32BitsToSingle(unchecked((int)v.GetUint()))
                        .ToString("0.####", CultureInfo.InvariantCulture), "3.1416");
                await ReadScalar(reader, connection, report, "WORD", $"%DB{db}.DBW12",
                    v => $"0x{v.GetUshort():X4}", "0xBEEF");
                await ReadScalar(reader, connection, report, "DWORD", $"%DB{db}.DBD14",
                    v => $"0x{v.GetUint():X8}", "0xDEADBEEF");

                // ── multi-tag single request ──
                await MultiRead(reader, connection, report, db);

                // ── write round-trip ──
                await WriteRoundTrip(reader, writer, connection, report, $"%DB{db}.DBW18",
                    b => b.AddTag("w", $"%DB{db}.DBW18", (short)6789),
                    v => unchecked((short)v.GetUshort()).ToString(CultureInfo.InvariantCulture), "6789");
                // 12345.5 is exact in IEEE-754 single, so a byte-perfect round-trip
                // renders back to the same string. (A value like 12345.678f is really
                // 12345.6787..., which would fail a literal string compare.)
                await WriteRoundTrip(reader, writer, connection, report, $"%DB{db}.DBD20",
                    b => b.AddTag("w", $"%DB{db}.DBD20", 12345.5f),
                    v => BitConverter.Int32BitsToSingle(unchecked((int)v.GetUint()))
                        .ToString("0.###", CultureInfo.InvariantCulture), "12345.5");

                // ── error path ──
                await ErrorPath(reader, connection, report);

                return report.Failures == 0 ? 0 : 1;
            }
            catch (Exception e)
            {
                report.Line();
                report.Line("## Result: FAIL - connection error");
                report.Line();
                report.Line($"**{e.GetType().Name}**: {e.Message}");
                report.Line();
                report.Line("```");
                report.Line(e.ToString());
                report.Line("```");
                report.Line();
                report.Line("Common causes: PUT/GET not enabled on the CPU; wrong rack/slot; " +
                            "S7-1200/1500 needs `--device-group OTHERS` or `--remote-tsap 0x0301`; " +
                            "TCP 102 blocked by a firewall.");
                return 1;
            }
            finally
            {
                report.Flush();
                connection?.Close();
            }
        }

        private static async Task<int> ProbeSingleRead(string connectionString, string address)
        {
            Console.WriteLine($"Connecting: {connectionString}");
            IPlcConnection? connection = null;
            try
            {
                var driver = new S7Driver(new DefaultTransportManager());
                connection = driver.Connect(connectionString);
                var s7 = (S7Connection)connection;
                Console.WriteLine($"Connected. Negotiated PDU length: {s7.NegotiatedPduLength} bytes");

                var reader = (PlcReader)connection;
                var rb = (DefaultPlcReadRequestBuilder)connection.ReadRequestBuilder;
                rb.AddTagAddress("tag", address);
                var resp = (DefaultPlcReadResponse)await reader.Read((DefaultPlcReadRequest)rb.Build());

                var code = resp.GetResponseCode("tag");
                Console.WriteLine($"Read {address}: {code}");
                if (code == PlcResponseCode.Ok)
                {
                    Console.WriteLine($"  value: {Render(resp.GetValue("tag"))}");
                    return 0;
                }
                if (code == PlcResponseCode.AccessDenied)
                {
                    Console.WriteLine(
                        "  The CPU refused access (S7 error 0x8104). Enable \"Permit access with " +
                        "PUT/GET communication from remote partner\" in the CPU's Protection & " +
                        "Security settings, then recompile and download.");
                }
                return 1;
            }
            catch (Exception e)
            {
                Console.WriteLine($"FAIL - {e.GetType().Name}: {e.Message}");
                return 1;
            }
            finally
            {
                connection?.Close();
            }
        }

        /// <summary>Renders a scalar <see cref="IPlcValue"/> without assuming its concrete type.</summary>
        private static string Render(IPlcValue v)
        {
            if (v.IsBool()) return v.GetBool().ToString();
            if (v.IsByte()) return $"0x{v.GetByte():X2}";
            if (v.IsUshort()) return $"0x{v.GetUshort():X4} ({unchecked((short)v.GetUshort())})";
            if (v.IsUint()) return $"0x{v.GetUint():X8} ({unchecked((int)v.GetUint())})";
            var raw = v.GetRaw();
            return raw != null ? BitConverter.ToString(raw) : v.GetString();
        }

        private static async Task ReadScalar(PlcReader reader, IPlcConnection connection, Report report,
            string type, string address, Func<org.apache.plc4net.api.value.IPlcValue, string> render, string expected)
        {
            try
            {
                var rb = (DefaultPlcReadRequestBuilder)connection.ReadRequestBuilder;
                rb.AddTagAddress("v", address);
                var resp = (DefaultPlcReadResponse)await reader.Read((DefaultPlcReadRequest)rb.Build());
                var code = resp.GetResponseCode("v");
                if (code != PlcResponseCode.Ok)
                {
                    report.Fail($"Read {type} {address}", $"response code {code}");
                    return;
                }
                var actual = render(resp.GetValue("v"));
                if (actual == expected)
                {
                    report.Pass($"Read {type} {address}", $"= {actual}");
                }
                else
                {
                    report.Fail($"Read {type} {address}", $"got {actual}, expected {expected}");
                }
            }
            catch (Exception e)
            {
                report.Fail($"Read {type} {address}", e.Message);
            }
        }

        private static async Task MultiRead(PlcReader reader, IPlcConnection connection, Report report, int db)
        {
            try
            {
                var rb = (DefaultPlcReadRequestBuilder)connection.ReadRequestBuilder;
                rb.AddTagAddress("bit", $"%DB{db}.DBX0.0");
                rb.AddTagAddress("byte", $"%DB{db}.DBB1");
                rb.AddTagAddress("word", $"%DB{db}.DBW12");
                var resp = (DefaultPlcReadResponse)await reader.Read((DefaultPlcReadRequest)rb.Build());
                var ok = resp.GetResponseCode("bit") == PlcResponseCode.Ok
                         && resp.GetResponseCode("byte") == PlcResponseCode.Ok
                         && resp.GetResponseCode("word") == PlcResponseCode.Ok
                         && resp.GetValue("bit").GetBool()
                         && resp.GetValue("byte").GetByte() == 0xA5
                         && resp.GetValue("word").GetUshort() == 0xBEEF;
                if (ok) report.Pass("Read 3 tags in one request", "all three correct");
                else report.Fail("Read 3 tags in one request", "a tag was wrong or errored");
            }
            catch (Exception e)
            {
                report.Fail("Read 3 tags in one request", e.Message);
            }
        }

        private static async Task WriteRoundTrip(PlcReader reader, PlcWriter writer, IPlcConnection connection,
            Report report, string address, Action<DefaultPlcWriteRequestBuilder> addWrite,
            Func<org.apache.plc4net.api.value.IPlcValue, string> render, string expected)
        {
            try
            {
                var wb = (DefaultPlcWriteRequestBuilder)connection.WriteRequestBuilder;
                addWrite(wb);
                var wResp = (DefaultPlcWriteResponse)await writer.Write((DefaultPlcWriteRequest)wb.Build());
                if (wResp.GetResponseCode("w") != PlcResponseCode.Ok)
                {
                    report.Fail($"Write {address}", $"write code {wResp.GetResponseCode("w")}");
                    return;
                }

                var rb = (DefaultPlcReadRequestBuilder)connection.ReadRequestBuilder;
                rb.AddTagAddress("v", address);
                var rResp = (DefaultPlcReadResponse)await reader.Read((DefaultPlcReadRequest)rb.Build());
                var actual = render(rResp.GetValue("v"));
                if (actual == expected) report.Pass($"Write + read-back {address}", $"= {actual}");
                else report.Fail($"Write + read-back {address}", $"read back {actual}, expected {expected}");
            }
            catch (Exception e)
            {
                report.Fail($"Write + read-back {address}", e.Message);
            }
        }

        private static async Task ErrorPath(PlcReader reader, IPlcConnection connection, Report report)
        {
            try
            {
                var rb = (DefaultPlcReadRequestBuilder)connection.ReadRequestBuilder;
                rb.AddTagAddress("v", "%DB31999.DBW0");
                var resp = (DefaultPlcReadResponse)await reader.Read((DefaultPlcReadRequest)rb.Build());
                var code = resp.GetResponseCode("v");
                if (code != PlcResponseCode.Ok)
                {
                    report.Pass("Read a non-existent DB", $"rejected with {code} (connection survived)");
                }
                else
                {
                    report.Fail("Read a non-existent DB", "unexpectedly returned Ok");
                }
            }
            catch (Exception e)
            {
                report.Fail("Read a non-existent DB", e.Message);
            }
        }

        private static Dictionary<string, string> ParseOptions(string[] args)
        {
            var opts = new Dictionary<string, string>(StringComparer.OrdinalIgnoreCase);
            for (var i = 1; i < args.Length - 1; i++)
            {
                if (args[i].StartsWith("--", StringComparison.Ordinal))
                {
                    opts[args[i].Substring(2)] = args[i + 1];
                    i++;
                }
            }
            return opts;
        }

        private static string Get(Dictionary<string, string> opts, string key, string fallback) =>
            opts.TryGetValue(key, out var v) ? v : fallback;

        private sealed class Report
        {
            private readonly List<string> _lines = new();
            private readonly List<string> _rows = new();
            public int Passes { get; private set; }
            public int Failures { get; private set; }

            public void Line(string text = "") => _lines.Add(text);

            public void Pass(string step, string detail)
            {
                Passes++;
                _rows.Add($"| ✅ | {step} | {detail} |");
            }

            public void Fail(string step, string detail)
            {
                Failures++;
                _rows.Add($"| ❌ | {step} | {detail} |");
            }

            public void Flush()
            {
                var w = Console.Out;
                foreach (var l in _lines) w.WriteLine(l);
                if (_rows.Count > 0)
                {
                    w.WriteLine();
                    w.WriteLine("| | Step | Detail |");
                    w.WriteLine("|---|---|---|");
                    foreach (var r in _rows) w.WriteLine(r);
                    w.WriteLine();
                    w.WriteLine(Failures == 0
                        ? $"## Result: PASS ({Passes}/{Passes})"
                        : $"## Result: FAIL ({Passes} passed, {Failures} failed)");
                }
                w.Flush();
                _lines.Clear();
                _rows.Clear();
            }
        }
    }
}
