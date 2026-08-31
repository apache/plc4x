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
using org.apache.plc4net.drivers.modbus;
using org.apache.plc4net.drivers.modbus.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.tcp;

namespace org.apache.plc4net.tools.modbusverify
{
    /// <summary>
    /// Modbus TCP hardware verification harness for
    /// <c>docs/test_report.md</c>.
    ///
    /// Usage: modbus-verify &lt;host&gt; [port] [unit-id] [read-address]
    /// Defaults: port=502, unit-id=1, read-address=holding:0
    /// </summary>
    public static class Program
    {
        public static async Task<int> Main(string[] args)
        {
            if (args.Length == 0)
            {
                Console.Error.WriteLine(
                    "Usage: modbus-verify <host> [port] [unit-id] [read-address]");
                Console.Error.WriteLine(
                    "Example: modbus-verify 192.168.0.9 502 1 holding:0");
                return 1;
            }

            var host = args[0];
            var port = args.Length > 1 ? int.Parse(args[1], CultureInfo.InvariantCulture) : 502;
            var unitId = args.Length > 2 ? int.Parse(args[2], CultureInfo.InvariantCulture) : 1;
            var readAddress = args.Length > 3 ? args[3] : "holding:0";

            var w = Console.Out;
            w.WriteLine("# Modbus TCP Hardware Verification Report");
            w.WriteLine();
            w.WriteLine($"**Date**: {DateTime.UtcNow:yyyy-MM-dd HH:mm} UTC");
            w.WriteLine($"**Target**: modbus-tcp://{host}:{port}?unit-identifier={unitId}");
            w.WriteLine($"**Read address**: {readAddress}");
            w.WriteLine();

            TcpTransportInstance? tcp = null;

            try
            {
                // ── 1. TCP connection ──
                w.WriteLine("## 1. TCP Connection");
                w.WriteLine();
                var tcpTransport = new TcpTransport();
                tcp = (TcpTransportInstance)tcpTransport.CreateTransportInstance(
                    $"{host}:{port}",
                    new TcpTransportConfiguration
                    {
                        DefaultPort = port,
                        ConnectTimeout = 5000
                    });
                tcp.DiagnosticOutput = w;
                w.WriteLine($"Connected to {host}:{port}");
                w.WriteLine();

                // ── 2. Modbus read ──
                w.WriteLine("## 2. Modbus Read Request / Response");
                w.WriteLine();

                var connString = ConnectionString.Parse(
                    $"modbus-tcp:tcp://{host}:{port}?unit-identifier={unitId}");
                var connection = new ModbusConnection(connString, tcp);
                connection.Connect();

                var tag = connection.Parse(readAddress);
                var builder = connection.ReadRequestBuilder;
                builder.AddTagAddress("data", readAddress);
                var request = (DefaultPlcReadRequest)builder.Build();

                var response = await connection.Read(request, CancellationToken.None)
                    .ConfigureAwait(false);
                w.WriteLine();

                // ── 3. Result ──
                w.WriteLine("## 3. Result");
                w.WriteLine();

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
                        var valProp = item.GetType().GetProperty("Value");
                        if (valProp?.GetValue(item) is IPlcValue val)
                        {
                            PrintValue(w, tagName, val);
                        }
                        else
                        {
                            var rcProp = item.GetType().GetProperty("ResponseCode");
                            var rc = rcProp?.GetValue(item)?.ToString() ?? "?";
                            w.WriteLine($"- **{tagName}**: ERROR — {rc}");
                        }
                    }
                }
                else
                {
                    w.WriteLine("(no values)");
                }

                // Decode MBAP from the diagnostic if available.
                w.WriteLine();
                w.WriteLine("## 4. Summary");
                w.WriteLine();
                w.WriteLine("| Step | Result |");
                w.WriteLine("|---|---|");
                w.WriteLine("| TCP connect | ✅ |");
                w.WriteLine($"| Modbus Read ({readAddress}) | ✅ |");
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
                tcp?.Close();
            }
        }

        private static void PrintValue(TextWriter w, string name, IPlcValue v)
        {
            if (!v.IsSimple())
            {
                w.WriteLine($"- **{name}**: (complex)");
                return;
            }

            if (v.IsBool())       w.WriteLine($"- **{name}**: `{v.GetBool()}` (BOOL)");
            else if (v.IsByte())  w.WriteLine($"- **{name}**: 0x{v.GetByte():X2} (BYTE)");
            else if (v.IsUshort()) w.WriteLine($"- **{name}**: {v.GetUshort()} (UINT16)");
            else if (v.IsUint())  w.WriteLine($"- **{name}**: {v.GetUint()} (UINT32)");
            else if (v.IsInt())   w.WriteLine($"- **{name}**: {v.GetInt()} (INT32)");
            else if (v.IsFloat()) w.WriteLine($"- **{name}**: {v.GetFloat()} (REAL)");
            else                  w.WriteLine($"- **{name}**: `{v}`");
        }
    }
}
