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
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// A minimal Modbus TCP slave on a loopback socket, written from the specification
    /// and independent of the generated model. It serves coils and registers
    /// (functions 1-6) and has two deliberately awkward behaviours: one register
    /// address is answered a byte at a time, as a TCP-to-serial gateway forwards it,
    /// and another is answered only after a delay.
    /// </summary>
    internal sealed class LoopbackModbusSlave : IDisposable
    {
        public const int SlowAddress = 200;
        public const int LateAddress = 300;
        public const int LateDelayMs = 500;
        public const int FirstIllegalAddress = 900;

        private readonly TcpListener _listener = new TcpListener(IPAddress.Loopback, 0);
        private readonly CancellationTokenSource _stop = new CancellationTokenSource();
        private readonly List<TcpClient> _clients = new List<TcpClient>();
        private readonly Task _acceptLoop;

        public LoopbackModbusSlave()
        {
            _listener.Start();
            _acceptLoop = Task.Run(AcceptLoop);
        }

        public int Port => ((IPEndPoint)_listener.LocalEndpoint).Port;

        public ushort[] Registers { get; } = new ushort[1000];

        public bool[] Coils { get; } = new bool[1000];

        /// <summary>Closes every client connection, as a device that was switched off would.</summary>
        public void DropClients()
        {
            lock (_clients)
            {
                foreach (var client in _clients)
                {
                    client.Close();
                }
                _clients.Clear();
            }
        }

        public void Dispose()
        {
            _stop.Cancel();
            _listener.Stop();
            DropClients();
            try
            {
                _acceptLoop.Wait(2000);
            }
            catch (AggregateException)
            {
                // The accept loop ends with an exception when the listener stops.
            }
        }

        private async Task AcceptLoop()
        {
            try
            {
                while (!_stop.IsCancellationRequested)
                {
                    var client = await _listener.AcceptTcpClientAsync().ConfigureAwait(false);
                    lock (_clients)
                    {
                        _clients.Add(client);
                    }
                    _ = Task.Run(() => Serve(client));
                }
            }
            catch (Exception) when (_stop.IsCancellationRequested)
            {
                // Shutting down.
            }
            catch (ObjectDisposedException)
            {
                // The listener was stopped.
            }
            catch (SocketException)
            {
                // The listener was stopped.
            }
        }

        private void Serve(TcpClient client)
        {
            try
            {
                var stream = client.GetStream();
                var header = new byte[7];
                while (true)
                {
                    stream.ReadExactly(header, 0, 7); // MBAP header + unit identifier
                    var length = (header[4] << 8) | header[5];
                    var pdu = new byte[length - 1];
                    stream.ReadExactly(pdu, 0, pdu.Length);

                    var response = Respond(pdu, out var delayMs, out var slow);
                    if (delayMs > 0)
                    {
                        Thread.Sleep(delayMs);
                    }

                    var frame = new byte[7 + response.Length];
                    frame[0] = header[0];
                    frame[1] = header[1];
                    frame[4] = (byte)((response.Length + 1) >> 8);
                    frame[5] = (byte)(response.Length + 1);
                    frame[6] = header[6];
                    Array.Copy(response, 0, frame, 7, response.Length);

                    if (slow)
                    {
                        foreach (var b in frame)
                        {
                            stream.WriteByte(b);
                            stream.Flush();
                            Thread.Sleep(5);
                        }
                    }
                    else
                    {
                        stream.Write(frame, 0, frame.Length);
                    }
                }
            }
            catch (Exception ex) when (ex is IOException or ObjectDisposedException or InvalidOperationException)
            {
                // The client went away.
            }
        }

        private byte[] Respond(byte[] pdu, out int delayMs, out bool slow)
        {
            delayMs = 0;
            slow = false;
            var function = pdu[0];
            var address = (pdu[1] << 8) | pdu[2];
            var second = (pdu[3] << 8) | pdu[4];

            if (address >= FirstIllegalAddress)
            {
                return new byte[] { (byte)(function | 0x80), 0x02 };
            }

            switch (function)
            {
                case 0x01:
                case 0x02:
                    {
                        var bytes = new byte[(second + 7) / 8];
                        for (var i = 0; i < second; i++)
                        {
                            if (Coils[address + i])
                            {
                                bytes[i / 8] |= (byte)(1 << (i % 8));
                            }
                        }
                        return Concat(new byte[] { function, (byte)bytes.Length }, bytes);
                    }
                case 0x03:
                case 0x04:
                    {
                        if (address == SlowAddress)
                        {
                            slow = true;
                        }
                        if (address == LateAddress)
                        {
                            delayMs = LateDelayMs;
                        }
                        var bytes = new byte[second * 2];
                        for (var i = 0; i < second; i++)
                        {
                            bytes[2 * i] = (byte)(Registers[address + i] >> 8);
                            bytes[2 * i + 1] = (byte)Registers[address + i];
                        }
                        return Concat(new byte[] { function, (byte)bytes.Length }, bytes);
                    }
                case 0x05:
                    Coils[address] = second == 0xFF00;
                    return (byte[])pdu.Clone();
                case 0x06:
                    Registers[address] = (ushort)second;
                    return (byte[])pdu.Clone();
                default:
                    return new byte[] { (byte)(function | 0x80), 0x01 };
            }
        }

        private static byte[] Concat(byte[] first, byte[] second)
        {
            var result = new byte[first.Length + second.Length];
            Array.Copy(first, result, first.Length);
            Array.Copy(second, 0, result, first.Length, second.Length);
            return result;
        }
    }
}