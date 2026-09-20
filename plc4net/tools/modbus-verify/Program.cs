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
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.api.value;
using org.apache.plc4net.drivers.modbus;
using org.apache.plc4net.drivers.modbus.messages;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.serial;
using org.apache.plc4net.transports.tcp;
using org.apache.plc4net.types;

namespace org.apache.plc4net.tools.modbusverify
{
    /// <summary>
    /// Modbus hardware verification harness. Writes a Markdown run log to stdout.
    ///
    /// Modbus TCP:
    ///   modbus-verify &lt;host&gt; [port] [unit-id] [read-address]
    ///   Defaults: port=502, unit-id=1, read-address=holding:0
    ///
    /// Modbus RTU (serial) — selected automatically when the first argument
    /// names a serial port (COMx, /dev/ttyUSB0, /dev/serial/...):
    ///   modbus-verify &lt;port&gt; [unit-id] [read-address]
    ///       [--baud 19200] [--parity Even|None|Odd] [--stop-bits One|Two] [--data-bits 8]
    ///   Defaults: unit-id=1, read-address=holding:0, 19200-8-E-1
    ///
    /// A full connection string is also accepted verbatim as the first argument
    /// (modbus-tcp:tcp://host:port?... or modbus-rtu://COMx?...).
    /// </summary>
    public static class Program
    {
        public static async Task<int> Main(string[] args)
        {
            if (args.Length == 0)
            {
                PrintUsage(Console.Error);
                return 1;
            }

            var positional = new List<string>();
            var flags = new Dictionary<string, string>(StringComparer.OrdinalIgnoreCase);
            for (var i = 0; i < args.Length; i++)
            {
                if (args[i].StartsWith("--", StringComparison.Ordinal))
                {
                    var key = args[i].Substring(2);
                    var value = i + 1 < args.Length && !args[i + 1].StartsWith("--", StringComparison.Ordinal)
                        ? args[++i]
                        : "true";
                    flags[key] = value;
                }
                else
                {
                    positional.Add(args[i]);
                }
            }

            var target = positional[0];
            return LooksLikeSerialTarget(target)
                ? await VerifyRtu(positional, flags).ConfigureAwait(false)
                : await VerifyTcp(positional, flags).ConfigureAwait(false);
        }

        private static void PrintUsage(TextWriter w)
        {
            w.WriteLine("Modbus TCP:  modbus-verify <host> [port] [unit-id] [read-address]");
            w.WriteLine("             modbus-verify 192.168.0.9 502 1 holding:0");
            w.WriteLine();
            w.WriteLine("Modbus RTU:  modbus-verify <COMx|/dev/ttyUSB0> [unit-id] [read-address] \\");
            w.WriteLine("                 [--baud 19200] [--parity Even] [--stop-bits One] [--data-bits 8]");
            w.WriteLine("                 [--quantity 1]");
            w.WriteLine("             modbus-verify COM3 1 holding:0 --baud 19200 --parity Even");
        }

        // Serial targets: COM1, COM12, /dev/ttyUSB0, /dev/ttyS1, /dev/serial/by-id/...
        private static bool LooksLikeSerialTarget(string s) =>
            !s.Contains("://")
            && (s.StartsWith("COM", StringComparison.OrdinalIgnoreCase)
                || s.StartsWith("/dev/tty", StringComparison.Ordinal)
                || s.StartsWith("/dev/serial", StringComparison.Ordinal));

        // ────────────────────────────────────────────────────────────
        //  Modbus RTU (serial)
        // ────────────────────────────────────────────────────────────

        private static async Task<int> VerifyRtu(
            List<string> positional, IReadOnlyDictionary<string, string> flags)
        {
            var w = Console.Out;

            string url;
            string port;
            int unitId;
            string readAddress;
            string baud, parity, stopBits, dataBits;
            ushort quantity;
            try
            {
                quantity = ParseQuantity(flags);
            }
            catch (Exception ex)
            {
                w.WriteLine();
                w.WriteLine("## Error");
                w.WriteLine();
                w.WriteLine($"**{ex.GetType().Name}**: {ex.Message}");
                w.WriteLine();
                return 1;
            }

            if (positional[0].Contains("://"))
            {
                url = positional[0];
                var pre = ConnectionString.Parse(url);
                port = pre.TransportConfig;
                unitId = pre.GetIntParameter("unit-identifier", 1);
                readAddress = positional.Count > 1 ? positional[1] : "holding:0";
                baud = pre.GetParameter("baud-rate") ?? "19200";
                parity = pre.GetParameter("parity") ?? "Even";
                stopBits = pre.GetParameter("stop-bits") ?? "One";
                dataBits = pre.GetParameter("data-bits") ?? "8";
            }
            else
            {
                port = positional[0];
                unitId = positional.Count > 1 ? int.Parse(positional[1], CultureInfo.InvariantCulture) : 1;
                readAddress = positional.Count > 2 ? positional[2] : "holding:0";
                baud = Flag(flags, "19200", "baud", "baud-rate");
                parity = Flag(flags, "Even", "parity");
                stopBits = Flag(flags, "One", "stop-bits", "stopbits");
                dataBits = Flag(flags, "8", "data-bits", "databits");
                url = $"modbus-rtu://{port}?unit-identifier={unitId}&baud-rate={baud}"
                      + $"&parity={parity}&stop-bits={stopBits}&data-bits={dataBits}";
            }

            w.WriteLine("# Modbus RTU hardware verification report");
            w.WriteLine();
            w.WriteLine($"**Date**: {DateTime.UtcNow:yyyy-MM-dd HH:mm} UTC");
            w.WriteLine($"**Target**: `{url}`");
            w.WriteLine($"**Line settings**: {baud} {dataBits}{ParityLetter(parity)}{(StopBitCount(stopBits))}");
            w.WriteLine($"**Read address**: `{readAddress}`");
            w.WriteLine($"**Quantity**: {quantity}");
            w.WriteLine();

            var parsed = ConnectionString.Parse(url);
            var serialTransport = new SerialTransport();
            var config = serialTransport.CreateConfiguration(parsed.Parameters);

            ITransportInstance instance = null;
            var rawOk = false;
            var crcOk = false;
            PlcResponseCode driverCode = PlcResponseCode.InternalError;
            string disconnectNote = null;

            try
            {
                // ── 1. Open the port ──
                w.WriteLine("## 1. Serial port");
                w.WriteLine();
                try
                {
                    instance = serialTransport.CreateTransportInstance(parsed.TransportConfig, config);
                    if (instance is IAsyncTransportInstance asyncTransport)
                    {
                        // Recorded, not written directly: this callback runs on the
                        // background receive-loop thread, and interleaving it with the
                        // main thread's own sequential WriteLine calls would garble the
                        // report. Printed once, from the main thread, in the finally below.
                        asyncTransport.RegisterDisconnectListener(ex =>
                            Interlocked.Exchange(ref disconnectNote,
                                $"Serial receive loop stopped: {ex.GetType().Name} - {ex.Message}"));
                    }
                }
                catch (Exception ex)
                {
                    w.WriteLine($"Could not open `{port}`: **{ex.GetType().Name}** — {ex.Message}");
                    w.WriteLine();
                    w.WriteLine("Check the port name (Device Manager on Windows / `ls /dev/tty*`), that no other");
                    w.WriteLine("program holds it open, and that the USB↔RS-485 adapter is plugged in.");
                    return 1;
                }

                w.WriteLine($"Opened `{port}`.");
                w.WriteLine();

                // Allow USB serial adapters and the background receive loop to
                // settle before the first frame is transmitted.
                await Task.Delay(100).ConfigureAwait(false);

                var tag = ModbusTag.Parse(readAddress);
                var requestPdu = BuildReadPdu(tag, quantity);
                var requestFrame = BuildRtuFrame((byte)unitId, requestPdu);

                // ── 2. Raw frame exchange (independent of the driver's framing) ──
                w.WriteLine("## 2. Raw frame exchange");
                w.WriteLine();
                w.WriteLine($"→ `{ToHex(requestFrame)}`  — unit {unitId}, {DescribeRead(tag)}");

                Drain(instance);
                instance.Write(requestFrame);
                var raw = await CollectResponse(instance, 1000).ConfigureAwait(false);

                // Some half-duplex adapters echo the transmitter onto the receiver.
                if (raw.Length >= requestFrame.Length
                    && raw.Take(requestFrame.Length).SequenceEqual(requestFrame))
                {
                    w.WriteLine($"← `{ToHex(raw)}`  ({raw.Length} bytes)");
                    w.WriteLine();
                    w.WriteLine("⚠️ The response begins with our own request bytes — the adapter is echoing its");
                    w.WriteLine("transmitter (half-duplex self-receive). `ModbusRtuConnection` does not strip that;");
                    w.WriteLine("use an auto-direction adapter that does not echo, or the driver needs a fix.");
                    raw = raw.Skip(requestFrame.Length).ToArray();
                    w.WriteLine();
                    w.WriteLine($"Remainder after the echo: `{(raw.Length == 0 ? "(nothing)" : ToHex(raw))}`");
                }
                else
                {
                    w.WriteLine($"← `{(raw.Length == 0 ? "(nothing)" : ToHex(raw))}`  ({raw.Length} bytes)");
                }

                w.WriteLine();
                rawOk = raw.Length >= 5;
                crcOk = raw.Length >= 4 && ModbusCRC.Validate(raw, raw.Length);
                var rawShapeOk = HasExpectedReadShape(raw, (byte)unitId, tag, quantity);
                AnalyzeRaw(w, raw, (byte)unitId);
                if (rawOk)
                {
                    w.WriteLine($"- Expected response shape {(rawShapeOk ? "valid" : "**invalid**")}");
                    PrintRawValues(w, raw, tag, quantity, rawShapeOk);
                }
                w.WriteLine();

                if (quantity != 1)
                {
                    w.WriteLine("## 3. Driver read — `ModbusRtuConnection`");
                    w.WriteLine();
                    w.WriteLine("Skipped: the current plc4net ModbusTag and driver read path issue a single-value");
                    w.WriteLine("request. The raw exchange above validates this fixed multi-register slave request.");
                    w.WriteLine();
                    w.WriteLine("## 4. Summary");
                    w.WriteLine();
                    w.WriteLine("| Step | Result |");
                    w.WriteLine("|---|---|");
                    w.WriteLine($"| Open `{port}` | {Mark(true)} |");
                    w.WriteLine($"| Raw response received | {Mark(rawOk)} |");
                    w.WriteLine($"| Raw response shape valid | {Mark(rawShapeOk)} |");
                    w.WriteLine($"| Raw response CRC valid | {Mark(crcOk)} |");
                    w.WriteLine($"| Driver read (`{readAddress}`) | skipped (single-value API) |");
                    w.WriteLine();
                    return rawOk && rawShapeOk && crcOk ? 0 : 1;
                }

                // ── 3. Driver read path ──
                w.WriteLine("## 3. Driver read — `ModbusRtuConnection`");
                w.WriteLine();

                var connection = new ModbusRtuConnection(parsed, instance);
                connection.Connect();

                var builder = connection.ReadRequestBuilder;
                builder.AddTagAddress("data", readAddress);
                var request = (DefaultPlcReadRequest)builder.Build();

                var response = (DefaultPlcReadResponse)await connection
                    .Read(request, CancellationToken.None).ConfigureAwait(false);

                driverCode = response.GetResponseCode("data");
                if (driverCode == PlcResponseCode.Ok)
                {
                    PrintValue(w, "data", response.GetValue("data"), tag.Type);
                }
                else
                {
                    w.WriteLine($"- **data**: `{driverCode}`");
                    if (!crcOk && rawOk)
                    {
                        w.WriteLine();
                        w.WriteLine("The raw exchange above got a reply but the driver did not decode it. If the raw");
                        w.WriteLine("CRC is invalid only under the driver, suspect `ModbusRtuConnection.SendAndReceive` —");
                        w.WriteLine("it reads whatever is available once ≥ 4 bytes have arrived, with no expected-length");
                        w.WriteLine("or t3.5 inter-frame-gap check, so a byte-at-a-time UART can hand it a partial frame.");
                    }
                }
                w.WriteLine();

                // ── 4. Summary ──
                w.WriteLine("## 4. Summary");
                w.WriteLine();
                w.WriteLine("| Step | Result |");
                w.WriteLine("|---|---|");
                w.WriteLine($"| Open `{port}` | {Mark(true)} |");
                w.WriteLine($"| Raw response received | {Mark(rawOk)} |");
                w.WriteLine($"| Raw response shape valid | {Mark(rawShapeOk)} |");
                w.WriteLine($"| Raw response CRC valid | {Mark(crcOk)} |");
                w.WriteLine($"| Driver read (`{readAddress}`) | {(driverCode == PlcResponseCode.Ok ? "✅" : "❌ " + driverCode)} |");
                w.WriteLine();

                return driverCode == PlcResponseCode.Ok && crcOk ? 0 : 1;
            }
            catch (Exception ex)
            {
                w.WriteLine();
                w.WriteLine("## Error");
                w.WriteLine();
                w.WriteLine($"**{ex.GetType().Name}**: {ex.Message}");
                w.WriteLine();
                w.WriteLine("```");
                w.WriteLine(ex);
                w.WriteLine("```");
                return 1;
            }
            finally
            {
                if (disconnectNote != null)
                {
                    w.WriteLine();
                    w.WriteLine(disconnectNote);
                }
                instance?.Close();
            }
        }

        private static byte[] BuildReadPdu(ModbusTag tag, ushort quantity) => tag.Type switch
        {
            ModbusTag.TagType.Coil => ModbusPDU.BuildReadBitsRequest(
                ModbusFunctionCodes.ReadCoils, tag.Address, quantity),
            ModbusTag.TagType.DiscreteInput => ModbusPDU.BuildReadBitsRequest(
                ModbusFunctionCodes.ReadDiscreteInputs, tag.Address, quantity),
            ModbusTag.TagType.HoldingRegister => ModbusPDU.BuildReadRegistersRequest(
                ModbusFunctionCodes.ReadHoldingRegisters, tag.Address, quantity),
            ModbusTag.TagType.InputRegister => ModbusPDU.BuildReadRegistersRequest(
                ModbusFunctionCodes.ReadInputRegisters, tag.Address, quantity),
            _ => throw new ArgumentException($"Cannot read a {tag.Type} tag.")
        };

        private static byte[] BuildRtuFrame(byte slaveAddress, byte[] pdu)
        {
            var frame = new byte[1 + pdu.Length + 2];
            frame[0] = slaveAddress;
            Array.Copy(pdu, 0, frame, 1, pdu.Length);
            var crc = ModbusCRC.Compute(frame, 0, frame.Length - 2);
            frame[frame.Length - 2] = crc[0];
            frame[frame.Length - 1] = crc[1];
            return frame;
        }

        private static void Drain(ITransportInstance t)
        {
            var n = t.GetNumBytesAvailable();
            if (n > 0)
            {
                t.Read(n);
            }
        }

        /// <summary>
        /// Collects every byte the slave sends, returning once the line has been
        /// quiet for ~60 ms after the first byte (a generous t3.5) or the overall
        /// deadline passes. Deliberately does not assume a frame length — the point
        /// is to see the whole reply exactly as it arrives.
        /// </summary>
        private static async Task<byte[]> CollectResponse(ITransportInstance t, int totalMs)
        {
            var buffer = new List<byte>();
            var deadline = Environment.TickCount64 + totalMs;
            long lastByteAt = 0;

            while (Environment.TickCount64 < deadline)
            {
                var n = t.GetNumBytesAvailable();
                if (n > 0)
                {
                    buffer.AddRange(t.Read(n));
                    lastByteAt = Environment.TickCount64;
                }
                else if (lastByteAt != 0 && Environment.TickCount64 - lastByteAt > 60)
                {
                    break;
                }

                await Task.Delay(5).ConfigureAwait(false);
            }

            return buffer.ToArray();
        }

        private static void AnalyzeRaw(TextWriter w, byte[] r, byte expectedAddr)
        {
            if (r.Length == 0)
            {
                w.WriteLine("No response. Check, in order:");
                w.WriteLine();
                w.WriteLine("- A/B (D+/D−) not swapped; signal grounds tied together");
                w.WriteLine("- the slave address matches (`MB_ADDR` on an S7-1200 `MB_SLAVE`)");
                w.WriteLine("- baud / parity / stop bits match `MB_COMM_LOAD`");
                w.WriteLine("- the CM 1241 is set to *half-duplex RS485 (2-wire)* in the device config");
                w.WriteLine("- `MB_SLAVE` is actually being called every scan, and `MB_COMM_LOAD` ran once at startup");
                w.WriteLine("- bias / termination — enable the CM 1241's internal bias for a short bench link");
                return;
            }

            if (r.Length < 4)
            {
                w.WriteLine($"Only {r.Length} byte(s) — shorter than any valid Modbus RTU frame (min 4).");
                w.WriteLine("Likely a bias/termination problem on the bus, or a truncated first frame.");
                return;
            }

            var crcValid = ModbusCRC.Validate(r, r.Length);
            w.WriteLine($"- Address byte `0x{r[AddressOffset]:X2}` ({r[AddressOffset]})"
                        + (r[AddressOffset] == expectedAddr ? " — matches" : $" — expected {expectedAddr}"));

            var fc = r[FunctionOffset];
            if ((fc & ModbusFunctionCodes.ErrorOffset) != 0)
            {
                var code = r.Length > ByteCountOffset ? r[ByteCountOffset] : (byte)0;
                var name = Enum.IsDefined(typeof(ModbusErrorCode), code)
                    ? ((ModbusErrorCode)code).ToString()
                    : "unknown";
                w.WriteLine($"- Function `0x{fc & 0x7F:X2}` with the error bit set — "
                            + $"exception `0x{code:X2}` ({name})");
                if (code == (byte)ModbusErrorCode.IllegalDataAddress)
                {
                    w.WriteLine("  The address is outside the slave's map. On an S7-1200 `MB_SLAVE`, holding");
                    w.WriteLine("  register 0 is the first word of the `MB_HOLD_REG` DB — make that DB big enough.");
                }
            }
            else
            {
                w.WriteLine($"- Function `0x{fc:X2}`"
                            + (r.Length > ByteCountOffset ? $", byte count {r[ByteCountOffset]}" : string.Empty));
            }

            w.WriteLine($"- CRC {(crcValid ? "valid ✅" : "**invalid** ❌ — corrupt, truncated, or contains echoed request bytes")}");
        }

        // ────────────────────────────────────────────────────────────
        //  Modbus TCP
        // ────────────────────────────────────────────────────────────

        private static async Task<int> VerifyTcp(
            List<string> positional, IReadOnlyDictionary<string, string> flags)
        {
            var w = Console.Out;

            string host;
            int port;
            int unitId;
            string readAddress;

            if (positional[0].Contains("://"))
            {
                var pre = ConnectionString.Parse(positional[0]);
                var hostPort = pre.TransportConfig.Split(':');
                host = hostPort[0];
                port = hostPort.Length > 1 ? int.Parse(hostPort[1], CultureInfo.InvariantCulture) : 502;
                unitId = pre.GetIntParameter("unit-identifier", 1);
                readAddress = positional.Count > 1 ? positional[1] : "holding:0";
            }
            else
            {
                host = positional[0];
                port = positional.Count > 1 ? int.Parse(positional[1], CultureInfo.InvariantCulture) : 502;
                unitId = positional.Count > 2 ? int.Parse(positional[2], CultureInfo.InvariantCulture) : 1;
                readAddress = positional.Count > 3 ? positional[3] : "holding:0";
            }

            w.WriteLine("# Modbus TCP hardware verification report");
            w.WriteLine();
            w.WriteLine($"**Date**: {DateTime.UtcNow:yyyy-MM-dd HH:mm} UTC");
            w.WriteLine($"**Target**: `modbus-tcp://{host}:{port}?unit-identifier={unitId}`");
            w.WriteLine($"**Read address**: `{readAddress}`");
            w.WriteLine();

            TcpTransportInstance tcp = null;
            try
            {
                w.WriteLine("## 1. TCP connection");
                w.WriteLine();
                var tcpTransport = new TcpTransport();
                tcp = (TcpTransportInstance)tcpTransport.CreateTransportInstance(
                    $"{host}:{port}",
                    new TcpTransportConfiguration { DefaultPort = port, ConnectTimeout = 5000 });
                tcp.DiagnosticOutput = w;
                w.WriteLine($"Connected to {host}:{port}.");
                w.WriteLine();

                w.WriteLine("## 2. Modbus read");
                w.WriteLine();
                var connString = ConnectionString.Parse(
                    $"modbus-tcp:tcp://{host}:{port}?unit-identifier={unitId}");
                var connection = new ModbusConnection(connString, tcp);
                connection.Connect();

                var tag = ModbusTag.Parse(readAddress);
                var builder = connection.ReadRequestBuilder;
                builder.AddTagAddress("data", readAddress);
                var request = (DefaultPlcReadRequest)builder.Build();

                var response = (DefaultPlcReadResponse)await connection
                    .Read(request, CancellationToken.None).ConfigureAwait(false);
                w.WriteLine();

                w.WriteLine("## 3. Result");
                w.WriteLine();
                var code = response.GetResponseCode("data");
                if (code == PlcResponseCode.Ok)
                {
                    PrintValue(w, "data", response.GetValue("data"), tag.Type);
                }
                else
                {
                    w.WriteLine($"- **data**: `{code}`");
                }
                w.WriteLine();

                w.WriteLine("## 4. Summary");
                w.WriteLine();
                w.WriteLine("| Step | Result |");
                w.WriteLine("|---|---|");
                w.WriteLine("| TCP connect | ✅ |");
                w.WriteLine($"| Modbus read (`{readAddress}`) | {(code == PlcResponseCode.Ok ? "✅" : "❌ " + code)} |");
                w.WriteLine();

                return code == PlcResponseCode.Ok ? 0 : 1;
            }
            catch (Exception ex)
            {
                w.WriteLine();
                w.WriteLine("## Error");
                w.WriteLine();
                w.WriteLine($"**{ex.GetType().Name}**: {ex.Message}");
                w.WriteLine();
                w.WriteLine("```");
                w.WriteLine(ex);
                w.WriteLine("```");
                return 1;
            }
            finally
            {
                tcp?.Close();
            }
        }

        // ────────────────────────────────────────────────────────────
        //  shared helpers
        // ────────────────────────────────────────────────────────────

        private static ushort ParseQuantity(IReadOnlyDictionary<string, string> flags)
        {
            var text = Flag(flags, "1", "quantity", "count");
            if (!ushort.TryParse(text, NumberStyles.Integer, CultureInfo.InvariantCulture, out var quantity)
                || quantity < 1 || quantity > 125)
            {
                throw new ArgumentException("--quantity must be an integer from 1 to 125.");
            }
            return quantity;
        }

        // RTU read-response frame layout: address(1) + function(1) + byteCount(1) + data(N) + crc(2).
        private const int AddressOffset = 0;
        private const int FunctionOffset = 1;
        private const int ByteCountOffset = 2;
        private const int DataOffset = 3;

        private static bool HasExpectedReadShape(
            byte[] frame, byte unitId, ModbusTag tag, ushort quantity)
        {
            var function = tag.Type switch
            {
                ModbusTag.TagType.Coil => ModbusFunctionCodes.ReadCoils,
                ModbusTag.TagType.DiscreteInput => ModbusFunctionCodes.ReadDiscreteInputs,
                ModbusTag.TagType.HoldingRegister => ModbusFunctionCodes.ReadHoldingRegisters,
                ModbusTag.TagType.InputRegister => ModbusFunctionCodes.ReadInputRegisters,
                _ => (byte)0
            };
            var dataBytes = tag.Type is ModbusTag.TagType.Coil or ModbusTag.TagType.DiscreteInput
                ? (quantity + 7) / 8
                : quantity * 2;
            return frame.Length == ModbusFunctionCodes.ReadResponseRtuFrameLength(dataBytes)
                   && frame[AddressOffset] == unitId
                   && frame[FunctionOffset] == function
                   && frame[ByteCountOffset] == dataBytes;
        }

        private static void PrintRawValues(
            TextWriter w, byte[] frame, ModbusTag tag, ushort quantity, bool shapeOk)
        {
            if (!shapeOk
                || tag.Type is not (ModbusTag.TagType.HoldingRegister
                                    or ModbusTag.TagType.InputRegister))
            {
                return;
            }

            var values = Enumerable.Range(0, quantity)
                .Select(i => (ushort)((frame[DataOffset + i * 2] << 8) | frame[DataOffset + i * 2 + 1]));
            w.WriteLine($"- Register values: `{string.Join(", ", values)}`");
        }

        private static string Flag(
            IReadOnlyDictionary<string, string> flags, string fallback, params string[] names)
        {
            foreach (var name in names)
            {
                if (flags.TryGetValue(name, out var value))
                {
                    return value;
                }
            }
            return fallback;
        }

        private static string DescribeRead(ModbusTag tag) => tag.Type switch
        {
            ModbusTag.TagType.Coil => $"read coil {tag.Address}",
            ModbusTag.TagType.DiscreteInput => $"read discrete input {tag.Address}",
            ModbusTag.TagType.HoldingRegister => $"read holding register {tag.Address}",
            ModbusTag.TagType.InputRegister => $"read input register {tag.Address}",
            _ => $"read {tag.Type} {tag.Address}"
        };

        private static string ToHex(byte[] bytes) =>
            BitConverter.ToString(bytes).Replace('-', ' ');

        private static string Mark(bool ok) => ok ? "✅" : "❌";

        private static char ParityLetter(string parity) =>
            string.IsNullOrEmpty(parity) ? '?' : char.ToUpperInvariant(parity[0]);

        private static int StopBitCount(string stopBits) =>
            stopBits.Equals("Two", StringComparison.OrdinalIgnoreCase) ? 2 : 1;

        /// <summary>
        /// Renders a value by the Modbus tag type that produced it, not by probing
        /// <see cref="IPlcValue"/> predicates: plc4net's value model coerces freely
        /// between related scalar types (a BOOL also answers <c>IsByte()</c> /
        /// <c>IsUshort()</c> / <c>IsInt()</c> true, and every numeric answers
        /// <c>IsBool()</c> true - see PlcBOOL.cs / PlcSimpleNumericValueAdapter.cs),
        /// so "the" type of a value can't be recovered from it after the fact. The
        /// tag the driver was asked to read already says what it is.
        /// </summary>
        private static void PrintValue(TextWriter w, string name, IPlcValue v, ModbusTag.TagType tagType)
        {
            if (v == null)
            {
                w.WriteLine($"- **{name}**: (null)");
                return;
            }

            switch (tagType)
            {
                case ModbusTag.TagType.Coil:
                case ModbusTag.TagType.DiscreteInput:
                    w.WriteLine($"- **{name}**: `{v.GetBool()}` (BOOL)");
                    break;
                case ModbusTag.TagType.HoldingRegister:
                case ModbusTag.TagType.InputRegister:
                    w.WriteLine($"- **{name}**: {v.GetUshort()} (UINT16)");
                    break;
                default:
                    w.WriteLine($"- **{name}**: `{v}`");
                    break;
            }
        }
    }
}
