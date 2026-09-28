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
using System.Net;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.tcp;
using Xunit;

namespace org.apache.plc4net.spi.test.transports
{
    /// <summary>
    /// Opening a connection through <see cref="TcpTransport"/>: host-name resolution,
    /// then the connect, both bounded by the connect timeout.
    /// </summary>
    /// <remarks>
    /// The lookup is injected, so a test can answer instantly, fail, or never answer at
    /// all without depending on what the machine's resolver does.
    /// </remarks>
    public class TcpTransportConnectTests : IDisposable
    {
        private const int SettleMillis = 2000;

        private readonly TcpListener _listener;
        private readonly int _port;

        public TcpTransportConnectTests()
        {
            _listener = new TcpListener(IPAddress.Loopback, 0);
            _listener.Start();
            _port = ((IPEndPoint) _listener.LocalEndpoint).Port;
        }

        public void Dispose()
        {
            _listener.Stop();
        }

        private static TcpTransportConfiguration Configuration(int connectTimeout = 5000)
        {
            return new TcpTransportConfiguration { ConnectTimeout = connectTimeout };
        }

        private static Func<string, CancellationToken, Task<IPAddress[]>> Answering(params IPAddress[] addresses)
        {
            return (host, cancellationToken) => Task.FromResult(addresses);
        }

        /// <summary>Opens a connection to the test listener and returns both ends.</summary>
        private TcpTransportInstance Open(
            TcpTransport transport, string host, TcpTransportConfiguration configuration, out Socket serverSide)
        {
            var accepting = _listener.AcceptSocketAsync();
            var instance = Assert.IsType<TcpTransportInstance>(
                transport.CreateTransportInstance($"{host}:{_port}", configuration));
            Assert.True(accepting.Wait(SettleMillis), "the listener never accepted the connection");
            serverSide = accepting.Result;
            return instance;
        }

        [Fact]
        public void A_literal_address_is_connected_without_a_lookup()
        {
            var transport = new TcpTransport(-1, (host, cancellationToken) =>
                throw new InvalidOperationException("a literal address must not be looked up"));

            using (var instance = Open(transport, "127.0.0.1", Configuration(), out var serverSide))
            {
                Assert.True(instance.IsOpen);
                serverSide.Dispose();
            }
        }

        [Fact]
        public void A_host_name_is_resolved_through_the_lookup_and_connected()
        {
            string? asked = null;
            var transport = new TcpTransport(-1, (host, cancellationToken) =>
            {
                asked = host;
                return Task.FromResult(new[] { IPAddress.Loopback });
            });

            using (var instance = Open(transport, "plc.example.test", Configuration(), out var serverSide))
            {
                Assert.Equal("plc.example.test", asked);
                Assert.Equal(new IPEndPoint(IPAddress.Loopback, _port), instance.RemoteAddress);
                serverSide.Dispose();
            }
        }

        [Fact]
        public void An_IPv4_address_is_preferred_when_the_lookup_returns_several()
        {
            // The listener is IPv4 only, so connecting at all proves the v4 address won.
            var transport = new TcpTransport(-1, Answering(IPAddress.IPv6Loopback, IPAddress.Loopback));

            using (var instance = Open(transport, "plc.example.test", Configuration(), out var serverSide))
            {
                Assert.Equal(IPAddress.Loopback, instance.RemoteAddress.Address);
                serverSide.Dispose();
            }
        }

        [Fact]
        public void A_failed_lookup_is_reported_as_a_transport_error()
        {
            var transport = new TcpTransport(-1, (host, cancellationToken) =>
                throw new SocketException((int) SocketError.HostNotFound));

            var e = Assert.Throws<TransportException>(
                () => transport.CreateTransportInstance($"nowhere.example.test:{_port}", Configuration()));

            Assert.Contains("Unable to resolve host 'nowhere.example.test'", e.Message);
            Assert.IsType<SocketException>(e.InnerException);
        }

        [Fact]
        public void A_lookup_that_returns_nothing_is_reported()
        {
            var transport = new TcpTransport(-1, Answering());

            var e = Assert.Throws<TransportException>(
                () => transport.CreateTransportInstance($"empty.example.test:{_port}", Configuration()));

            Assert.Contains("did not resolve to any address", e.Message);
        }

        [Fact]
        public async Task A_lookup_that_never_answers_is_abandoned_when_the_timeout_expires()
        {
            var transport = new TcpTransport(-1, async (host, cancellationToken) =>
            {
                // Answers only by being cancelled, like a resolver that has gone silent.
                await Task.Delay(Timeout.Infinite, cancellationToken);
                return Array.Empty<IPAddress>();
            });

            var call = Task.Run(
                () => transport.CreateTransportInstance($"silent.example.test:{_port}", Configuration(connectTimeout: 200)));

            var finished = await Task.WhenAny(call, Task.Delay(10 * SettleMillis));
            Assert.Same(call, finished);

            var e = await Assert.ThrowsAsync<TransportException>(() => call);
            Assert.Contains("timed out after 200 ms", e.Message);
        }

        [Theory]
        [InlineData(0)]
        [InlineData(-1)]
        public void A_connect_timeout_of_zero_or_less_means_no_timeout(int connectTimeout)
        {
            // Zero used to mean "give up at once", which no caller can want.
            var transport = new TcpTransport();

            using (var instance = Open(transport, "127.0.0.1", Configuration(connectTimeout), out var serverSide))
            {
                Assert.True(instance.IsOpen);
                serverSide.Dispose();
            }
        }
    }
}