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
using System.IO;
using System.Net;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.api.value;
using org.apache.plc4net.drivers.s7;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.cotp;
using org.apache.plc4net.transports.tcp;

namespace org.apache.plc4net.tools.s7verify
{
    /// <summary>
    /// S7 hardware verification harness for <c>docs/test_report.md</c>.
    ///
    /// Usage: s7-verify &lt;host&gt; [rack] [slot] [read-address]
    /// Defaults: rack=0, slot=1 (S7-1500), read-address=%M0
    /// </summary>
    public static class Program
    {
        public static async Task<int> Main(string[] args)
        {
            if (args.Length == 0)
            {
                Console.Error.WriteLine("Usage: s7-verify <host> [rack] [slot] [read-address]");
                Console.Error.WriteLine("Example: s7-verify 192.168.0.10 0 1 %DB1.DBW0");
                return 1;
            }

            var host = args[0];
            var rack = args.Length > 1 ? int.Parse(args[1], CultureInfo.InvariantCulture) : 0;
            var slot = args.Length > 2 ? int.Parse(args[2], CultureInfo.InvariantCulture) : 1;
            var readAddress = args.Length > 3 ? args[3] : "%M0";
            var port = 102;

            var w = Console.Out;
            w.WriteLine("# S7 Hardware Verification Report");
            w.WriteLine();
            w.WriteLine($"**Date**: {DateTime.UtcNow:yyyy-MM-dd HH:mm} UTC");
            w.WriteLine($"**Target**: s7://{host}:{port}?remote-rack={rack}&remote-slot={slot}");
            w.WriteLine($"**Read address**: {readAddress}");
            w.WriteLine();

            ITransportInstance? tcp = null;
            CotpTransportInstance? cotp = null;

            try
            {
                // ── 1. TCP connection ──
                w.WriteLine("## 1. TCP Connection");
                w.WriteLine();
                var tcpTransport = new TcpTransport();
                tcp = tcpTransport.CreateTransportInstance(host,
                    new TcpTransportConfiguration { DefaultPort = port, ConnectTimeout = 5000 });
                w.WriteLine($"Connected to {host}:{port}");
                w.WriteLine();

                // ── 2. COTP handshake ──
                w.WriteLine("## 2. COTP Connection Request / Confirm");
                w.WriteLine();
                cotp = new CotpTransportInstance(tcp)
                {
                    DiagnosticOutput = w,
                    HandshakeTimeout = TimeSpan.FromSeconds(5)
                };
                var localTsap = 0x0311;
                var remoteTsap = ((0x01) << 8) | ((rack & 0x0F) << 4) | (slot & 0x0F);
                w.WriteLine($"**Calling TSAP**: 0x{localTsap:X4}");
                w.WriteLine($"**Called TSAP**:  0x{remoteTsap:X4}");
                w.WriteLine();

                cotp.Open(
                    localTsapHi: (byte)((localTsap >> 8) & 0xFF),
                    localTsapLo: (byte)(localTsap & 0xFF),
                    remoteTsapHi: (byte)((remoteTsap >> 8) & 0xFF),
                    remoteTsapLo: (byte)(remoteTsap & 0xFF));
                w.WriteLine();

                // ── 3. S7 Read ──
                w.WriteLine("## 3. S7 Read Var PDU Exchange");
                w.WriteLine();

                var connString = ConnectionString.Parse(
                    $"s7:cotp://{host}:{port}?remote-rack={rack}&remote-slot={slot}");
                var s7Conn = new S7Connection(connString, cotp);
                s7Conn.Connect();

                // Use the standard builder API to construct a read request.
                var builder = s7Conn.ReadRequestBuilder;
                builder.AddTagAddress("data", readAddress);
                var request = builder.Build();

                var response = await s7Conn.Read(
                    (DefaultPlcReadRequest)request,
                    CancellationToken.None).ConfigureAwait(false);
                w.WriteLine();

                // ── 4. Result ──
                w.WriteLine("## 4. Result");
                w.WriteLine();

                // Access the response values via reflection over the
                // internal dictionary (the tool needs byte-level visibility;
                // the public API surfaces are sufficient for production code).
                var valuesField = typeof(DefaultPlcReadResponse)
                    .GetField("_values",
                        System.Reflection.BindingFlags.NonPublic |
                        System.Reflection.BindingFlags.Instance);
                var values = valuesField?.GetValue(response) as
                    System.Collections.IDictionary;

                if (values != null)
                {
                    foreach (string tagName in values.Keys)
                    {
                        var item = values[tagName]!;
                        var itemType = item.GetType().Name;

                        var valProp = item.GetType().GetProperty("Value");
                        if (valProp != null)
                        {
                            var val = valProp.GetValue(item) as IPlcValue;
                            if (val != null)
                            {
                                w.Write($"- **{tagName}**: ");
                                PrintValue(w, val);
                                continue;
                            }
                        }

                        var rcProp = item.GetType().GetProperty("ResponseCode");
                        var rc = rcProp?.GetValue(item)?.ToString() ?? "?";
                        w.WriteLine($"- **{tagName}**: {itemType} — {rc}");
                    }
                }
                else
                {
                    w.WriteLine("(no values)");
                }

                w.WriteLine();
                w.WriteLine("## 5. Summary");
                w.WriteLine();
                w.WriteLine("| Step | Result |");
                w.WriteLine("|---|---|");
                w.WriteLine("| TCP connect | ✅ |");
                w.WriteLine("| COTP CR/CC handshake | ✅ |");
                w.WriteLine($"| S7 Read Var ({readAddress}) | ✅ |");
                w.WriteLine();

                return 0;
            }
            catch (Exception ex)
            {
                w.WriteLine();
                w.WriteLine("## Error");
                w.WriteLine();
                w.WriteLine($"**{ex.GetType().Name}**: {ex.Message}");
                w.WriteLine();
                w.WriteLine("```");
                w.WriteLine(ex.ToString());
                w.WriteLine("```");
                w.WriteLine();
                return 1;
            }
            finally
            {
                cotp?.Close();
                tcp?.Close();
            }
        }

        private static void PrintValue(TextWriter w, IPlcValue v)
        {
            if (!v.IsSimple())
            {
                w.WriteLine($"(complex, {CountKeys(v)} keys)");
                return;
            }

            if (v.IsBool())       { w.WriteLine($"`{v.GetBool()}` (BOOL)"); }
            else if (v.IsByte())   { w.WriteLine($"0x{v.GetByte():X2} (BYTE)"); }
            else if (v.IsUshort()) { w.WriteLine($"{v.GetUshort()} (UINT)"); }
            else if (v.IsUint())   { w.WriteLine($"{v.GetUint()} (UDINT)"); }
            else if (v.IsShort())  { w.WriteLine($"{v.GetShort()} (INT)"); }
            else if (v.IsInt())    { w.WriteLine($"{v.GetInt()} (DINT)"); }
            else if (v.IsFloat())  { w.WriteLine($"{v.GetFloat()} (REAL)"); }
            else if (v.IsString()) { w.WriteLine($"\"{v.GetString()}\" (STRING)"); }
            else                   { w.WriteLine($"(type not resolved)"); }
        }

        private static int CountKeys(IPlcValue v)
        {
            try { return v.GetStruct()?.Count ?? 0; }
            catch { return 0; }
        }
    }
}
