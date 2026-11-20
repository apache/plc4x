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
using System.Threading.Tasks;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.types;
using Xunit;
using static org.apache.plc4net.drivers.modbus.test.ConnectionTestSupport;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// The requests and answers of plc4j's own driver test data,
    /// <c>protocols/modbus/src/test/resources/protocols/modbus/tcp/ManualFactoryModbusTCPDriverTest-testsuite.xml</c>,
    /// for the tags this driver supports. The same tag has to put the same PDU on the wire
    /// and decode the same answer; the bytes below are copied from that file (the transaction
    /// identifier and the unit identifier are the connection's own).
    /// </summary>
    public class Plc4jParityTests
    {
        private static byte[] Pdu(string hex) => Convert.FromHexString(hex);

        [Theory]
        [InlineData("max-registers-per-request=2", "holding-register:1:UINT", "holding-register:3:UINT", "03020001", "03020003", "0300000001", "0300020001")]
        [InlineData("max-coils-per-request=8", "coil:1", "coil:9", "010101", "010101", "0100000001", "0100080001")]
        public async Task Configured_window_keeps_distant_tags_in_separate_requests(
            string options, string first, string last, string firstReply, string lastReply,
            string firstRequest, string lastRequest)
        {
            using var transport = NewTransport();
            using var connection = NewTcp(transport, "?" + options);
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu(firstReply)));
            var server = Task.Run(() =>
            {
                var requests = transport.WaitForWrittenData(24, 5000);
                transport.InjectTestData(ModbusFrames.Tcp(2, 1, Pdu(lastReply)));
                return requests;
            });
            var response = (DefaultPlcReadResponse)await connection.ReadRequestBuilder!
                .AddTagAddress("first", first).AddTagAddress("last", last).Build().ExecuteAsync();
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("first"));
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("last"));
            var expected = new List<byte>(ModbusFrames.Tcp(1, 1, Pdu(firstRequest)));
            expected.AddRange(ModbusFrames.Tcp(2, 1, Pdu(lastRequest)));
            Assert.Equal(expected.ToArray(), await server);
        }

        [Theory]
        [InlineData("holding-register:1:UINT", "holding-register:3:UINT", "0306000100000003", "0300000003")]
        [InlineData("input-register:1:UINT", "input-register:3:UINT", "0406000100000003", "0400000003")]
        [InlineData("extended-register:65536:UINT", "extended-register:65538:UINT", "14080706000100000003", "140706000715A00003")]
        public async Task Read_window_includes_gaps_and_splits_values_by_address(
            string first, string last, string reply, string request)
        {
            using var transport = NewTransport();
            using var connection = NewTcp(transport);
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu(reply)));
            var response = (DefaultPlcReadResponse)await connection.ReadRequestBuilder!
                .AddTagAddress("last", last).AddTagAddress("first", first).Build().ExecuteAsync();
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("first"));
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("last"));
            Assert.Equal((ushort)1, response.GetValue("first")!.GetUshort());
            Assert.Equal((ushort)3, response.GetValue("last")!.GetUshort());
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu(request)), transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("coil:1", "coil:9", "01020101", "0100000009")]
        [InlineData("discrete-input:1", "discrete-input:9", "02020101", "0200000009")]
        public async Task Bit_read_window_includes_gaps(string first, string last, string reply, string request)
        {
            using var transport = NewTransport();
            using var connection = NewTcp(transport);
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu(reply)));
            var response = (DefaultPlcReadResponse)await connection.ReadRequestBuilder!
                .AddTagAddress("first", first).AddTagAddress("last", last).Build().ExecuteAsync();
            Assert.True(response.GetValue("first")!.GetBool());
            Assert.True(response.GetValue("last")!.GetBool());
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu(request)), transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("-1")]
        [InlineData("65536")]
        [InlineData("invalid")]
        public void Invalid_optimizer_limit_is_rejected(string value)
        {
            using var transport = NewTransport();
            Assert.Throws<org.apache.plc4net.exceptions.PlcConnectionException>(() =>
                NewTcp(transport, "?max-registers-per-request=" + value));
        }

        [Theory]
        [InlineData("0")]
        [InlineData("65535")]
        public async Task Optimizer_limit_uses_protocol_ceiling_for_zero_or_excess(string value)
        {
            using var transport = NewTransport();
            using var connection = NewTcp(transport, "?max-registers-per-request=" + value);
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("0306000100000003")));
            var response = (DefaultPlcReadResponse)await connection.ReadRequestBuilder!
                .AddTagAddress("a", "holding-register:1:UINT")
                .AddTagAddress("b", "holding-register:3:UINT").Build().ExecuteAsync();
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("b"));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000003")), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Discovery_cancels_a_stalled_tls_handshake()
        {
            using var listener = new System.Net.Sockets.TcpListener(System.Net.IPAddress.Loopback, 0);
            listener.Start();
            using var cancellation = new System.Threading.CancellationTokenSource();
            var driver = new ModbusTcpDriver(new org.apache.plc4net.spi.transports.DefaultTransportManager());
            var pending = driver.DiscoverAsync($"modbus-tcp:tls://127.0.0.1:{((System.Net.IPEndPoint)listener.LocalEndpoint).Port}?connect-timeout=0", cancellation.Token);
            using var peer = await listener.AcceptTcpClientAsync().WaitAsync(TimeSpan.FromSeconds(5));
            cancellation.Cancel();
            await Assert.ThrowsAnyAsync<OperationCanceledException>(() => pending.WaitAsync(TimeSpan.FromSeconds(5)));
            var bytes = new byte[4096];
            try
            {
                while (await peer.GetStream().ReadAsync(bytes.AsMemory()).AsTask().WaitAsync(TimeSpan.FromSeconds(5)) != 0) { }
            }
            catch (System.IO.IOException ex) when (ex.InnerException is System.Net.Sockets.SocketException
            { SocketErrorCode: System.Net.Sockets.SocketError.ConnectionReset })
            {
                // Linux may reset a socket closed with unread handshake bytes instead of returning EOF.
            }
        }

        [Fact]
        public async Task Factory_does_not_mutate_shared_tls_hostname()
        {
            using var listener = new System.Net.Sockets.TcpListener(System.Net.IPAddress.Loopback, 0);
            listener.Start();
            var config = new org.apache.plc4net.transports.tcp.TcpTransportConfiguration
            { Tls = true, TlsTargetHost = "original", ConnectTimeout = 0 };
            var transport = new org.apache.plc4net.transports.tcp.TlsTransport();
            using var cancellation = new System.Threading.CancellationTokenSource();
            var pending = transport.CreateTransportInstanceAsync($"localhost:{((System.Net.IPEndPoint)listener.LocalEndpoint).Port}", config, cancellation.Token);
            using var peer = await listener.AcceptTcpClientAsync().WaitAsync(TimeSpan.FromSeconds(5));
            Assert.Equal("original", config.TlsTargetHost);
            cancellation.Cancel();
            await Assert.ThrowsAnyAsync<OperationCanceledException>(() => pending.WaitAsync(TimeSpan.FromSeconds(5)));
        }

        [Fact]
        public void Tls_transport_uses_its_own_configuration_prefix()
        {
            var transport = new org.apache.plc4net.transports.tcp.TlsTransport();
            var configuration = (org.apache.plc4net.transports.tcp.TcpTransportConfiguration)
                transport.CreateConfiguration(new Dictionary<string, string>
                {
                    ["tls.connect-timeout"] = "321",
                    ["tcp.connect-timeout"] = "999"
                });
            Assert.True(configuration.Tls);
            Assert.Equal(802, configuration.DefaultPort);
            Assert.Equal(321, configuration.ConnectTimeout);
        }

        [Fact]
        public async Task Subscription_can_unsubscribe_in_callback_and_close_prevents_new_subscriptions()
        {
            using var transport = NewTransport();
            using var connection = NewTcp(transport);
            var subscribed = (org.apache.plc4net.messages.IPlcSubscriptionResponse)await connection.SubscriptionRequestBuilder!
                .AddCyclicField("v", "holding-register:1:UINT", TimeSpan.FromMilliseconds(100)).Build().ExecuteAsync();
            var result = new TaskCompletionSource<org.apache.plc4net.spi.drivers.subscriptions.SubscriptionValueEventArgs>(
                TaskCreationOptions.RunContinuationsAsynchronously);
            var handle = subscribed.GetSubscriptionHandle("v");
            handle.PlcSubscriptionEvent += (_, e) =>
            {
                connection.UnsubscriptionRequestBuilder!.AddHandles(handle).Build().ExecuteAsync().GetAwaiter().GetResult();
                result.TrySetResult((org.apache.plc4net.spi.drivers.subscriptions.SubscriptionValueEventArgs)e);
            };
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("03021234")));
            var value = await result.Task.WaitAsync(TimeSpan.FromSeconds(5));
            Assert.Equal(PlcResponseCode.Ok, value.ResponseCode);
            Assert.Equal((ushort)0x1234, value.Value!.GetUshort());
            transport.InjectTestData(ModbusFrames.Tcp(2, 1, Pdu("03025678")));
            Assert.Equal((ushort)0x5678, ValueOf(await ReadAsync(connection, "holding-register:1:UINT")).GetUshort());
            connection.Close();
            await Assert.ThrowsAsync<ObjectDisposedException>(() => connection.SubscriptionRequestBuilder!
                .AddCyclicField("closed", "holding-register:1:UINT", TimeSpan.FromSeconds(1)).Build().ExecuteAsync());
        }

        [Theory]
        [InlineData(true, "localhost", true)]
        [InlineData(true, "wrong-host", false)]
        [InlineData(false, "localhost", false)]
        public async Task Tls_validates_trust_and_hostname(bool trusted, string hostname, bool succeeds)
        {
            using var key = System.Security.Cryptography.RSA.Create(2048);
            var request = new System.Security.Cryptography.X509Certificates.CertificateRequest("CN=localhost", key,
                System.Security.Cryptography.HashAlgorithmName.SHA256, System.Security.Cryptography.RSASignaturePadding.Pkcs1);
            request.CertificateExtensions.Add(new System.Security.Cryptography.X509Certificates.X509BasicConstraintsExtension(true, false, 0, true));
            var names = new System.Security.Cryptography.X509Certificates.SubjectAlternativeNameBuilder();
            names.AddDnsName("localhost");
            request.CertificateExtensions.Add(names.Build());
            using var original = request.CreateSelfSigned(DateTimeOffset.UtcNow.AddDays(-1), DateTimeOffset.UtcNow.AddDays(1));
            using var certificate = new System.Security.Cryptography.X509Certificates.X509Certificate2(
                original.Export(System.Security.Cryptography.X509Certificates.X509ContentType.Pfx));
            var policy = new System.Security.Cryptography.X509Certificates.X509ChainPolicy
            {
                TrustMode = System.Security.Cryptography.X509Certificates.X509ChainTrustMode.CustomRootTrust,
                RevocationMode = System.Security.Cryptography.X509Certificates.X509RevocationMode.NoCheck
            };
            if (trusted) policy.CustomTrustStore.Add(certificate);
            using var listener = new System.Net.Sockets.TcpListener(System.Net.IPAddress.Loopback, 0);
            listener.Start();
            using var deadline = new System.Threading.CancellationTokenSource(TimeSpan.FromSeconds(10));
            var server = Task.Run(async () =>
            {
                using var client = await listener.AcceptTcpClientAsync(deadline.Token);
                using var stream = new System.Net.Security.SslStream(client.GetStream());
                try
                {
                    await stream.AuthenticateAsServerAsync(new System.Net.Security.SslServerAuthenticationOptions { ServerCertificate = certificate }, deadline.Token);
                    if (succeeds)
                    {
                        var bytes = new byte[12];
                        await stream.ReadExactlyAsync(bytes, deadline.Token);
                        Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000001")), bytes);
                        await stream.WriteAsync(ModbusFrames.Tcp(1, 1, Pdu("03021234")), deadline.Token);
                    }
                }
                catch (System.Security.Authentication.AuthenticationException) when (!succeeds) { }
                catch (System.IO.IOException) when (!succeeds) { }
            });
            var config = new org.apache.plc4net.transports.tcp.TcpTransportConfiguration
            { Tls = true, TlsTargetHost = hostname, TlsChainPolicy = policy };
            if (succeeds)
            {
                using var transport = new org.apache.plc4net.transports.tcp.TcpTransportInstance((System.Net.IPEndPoint)listener.LocalEndpoint, config);
                using var connection = new ModbusTcpConnection(org.apache.plc4net.spi.drivers.ConnectionString.Parse("modbus-tcp://localhost"), transport);
                Assert.Equal((ushort)0x1234, ValueOf(await ReadAsync(connection, "holding-register:1:UINT", deadline.Token)).GetUshort());
            }
            else
            {
                var error = Assert.Throws<org.apache.plc4net.spi.transports.TransportException>(() =>
                    new org.apache.plc4net.transports.tcp.TcpTransportInstance((System.Net.IPEndPoint)listener.LocalEndpoint, config));
                Assert.IsType<System.Security.Authentication.AuthenticationException>(error.InnerException);
            }
            await server;
        }

        [Theory]
        [InlineData(false)]
        [InlineData(true)]
        public async Task File_records_cross_a_file_boundary_over_tcp(bool write)
        {
            using var listener = new System.Net.Sockets.TcpListener(System.Net.IPAddress.Loopback, 0);
            listener.Start();
            using var deadline = new System.Threading.CancellationTokenSource(TimeSpan.FromSeconds(5));
            var server = Task.Run(async () =>
            {
                using var client = await listener.AcceptTcpClientAsync(deadline.Token);
                var stream = client.GetStream();
                for (ushort transaction = 1; transaction <= 2; transaction++)
                {
                    var item = transaction == 1 ? "060001270F0001" : "06000200000001";
                    var value = transaction == 1 ? "0001" : "0002";
                    var request = ModbusFrames.Tcp(transaction, 1, Pdu(write ? "1509" + item + value : "1407" + item));
                    var actual = new byte[request.Length];
                    await stream.ReadExactlyAsync(actual, deadline.Token);
                    Assert.Equal(request, actual);
                    await stream.WriteAsync(write ? request : ModbusFrames.Tcp(transaction, 1, Pdu("14040306" + value)), deadline.Token);
                }
            });
            var driver = new ModbusTcpDriver(new org.apache.plc4net.spi.transports.DefaultTransportManager());
            using var connection = driver.Connect($"modbus-tcp://127.0.0.1:{((System.Net.IPEndPoint)listener.LocalEndpoint).Port}");
            if (write)
                Assert.Equal(PlcResponseCode.Ok, (await WriteAsync(connection, "extended-register:9999[0..1]:INT", new List<short> { 1, 2 }, deadline.Token)).GetResponseCode("v"));
            else
                Assert.Equal(new short[] { 1, 2 }, ValueOf(await ReadAsync(connection, "extended-register:9999[0..1]:INT", deadline.Token)).GetList().ConvertAll(v => v.GetShort()));
            await server;
        }

        [Fact]
        public async Task Discovery_keeps_success_and_failure_for_multiple_endpoints()
        {
            using var listener = new System.Net.Sockets.TcpListener(System.Net.IPAddress.Loopback, 0);
            listener.Start();
            using var deadline = new System.Threading.CancellationTokenSource(TimeSpan.FromSeconds(5));
            var server = Task.Run(async () =>
            {
                using var client = await listener.AcceptTcpClientAsync(deadline.Token);
                var bytes = new byte[11];
                await client.GetStream().ReadExactlyAsync(bytes, deadline.Token);
                Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("2B0E0100")), bytes);
                await client.GetStream().WriteAsync(ModbusFrames.Tcp(1, 1, Pdu("2B0E01010000010003414243")), deadline.Token);
            });
            var driver = new ModbusTcpDriver(new org.apache.plc4net.spi.transports.DefaultTransportManager());
            var results = await driver.DiscoverAsync(new[] { "invalid", $"modbus-tcp://127.0.0.1:{((System.Net.IPEndPoint)listener.LocalEndpoint).Port}" }, deadline.Token);
            Assert.NotNull(results[0].Error);
            Assert.Null(results[1].Error);
            Assert.Equal(Pdu("414243"), results[1].Objects![0]);
            await server;
        }

        [Fact]
        public async Task Tls_rejects_an_untrusted_server_certificate()
        {
            using var key = System.Security.Cryptography.RSA.Create(2048);
            var request = new System.Security.Cryptography.X509Certificates.CertificateRequest("CN=localhost", key,
                System.Security.Cryptography.HashAlgorithmName.SHA256, System.Security.Cryptography.RSASignaturePadding.Pkcs1);
            using var certificate = request.CreateSelfSigned(DateTimeOffset.UtcNow.AddDays(-1), DateTimeOffset.UtcNow.AddDays(1));
            using var listener = new System.Net.Sockets.TcpListener(System.Net.IPAddress.Loopback, 0);
            listener.Start();
            using var deadline = new System.Threading.CancellationTokenSource(TimeSpan.FromSeconds(5));
            var server = Task.Run(async () =>
            {
                using var client = await listener.AcceptTcpClientAsync(deadline.Token);
                using var stream = new System.Net.Security.SslStream(client.GetStream());
                try { await stream.AuthenticateAsServerAsync(new System.Net.Security.SslServerAuthenticationOptions { ServerCertificate = certificate }, deadline.Token); }
                catch (System.Security.Authentication.AuthenticationException) { }
                catch (System.IO.IOException) { }
            });
            var driver = new ModbusTcpDriver(new org.apache.plc4net.spi.transports.DefaultTransportManager());
            Assert.Throws<org.apache.plc4net.spi.transports.TransportException>(() =>
                driver.Connect($"modbus-tcp:tls://127.0.0.1:{((System.Net.IPEndPoint)listener.LocalEndpoint).Port}"));
            await server;
        }

        [Theory]
        [InlineData("modbus-tcp", "0300000001", "03021234")]
        [InlineData("modbus-rtu", "0300000001", "03021234")]
        [InlineData("modbus-ascii", "", "")]
        public async Task Drivers_read_over_real_tcp(string protocol, string request, string response)
        {
            using var listener = new System.Net.Sockets.TcpListener(System.Net.IPAddress.Loopback, 0);
            listener.Start();
            using var deadline = new System.Threading.CancellationTokenSource(TimeSpan.FromSeconds(5));
            var expected = protocol == "modbus-tcp" ? ModbusFrames.Tcp(1, 1, Pdu(request))
                : protocol == "modbus-rtu" ? ModbusFrames.Rtu(1, Pdu(request))
                : System.Text.Encoding.ASCII.GetBytes(":010300000001FB\r\n");
            var reply = protocol == "modbus-tcp" ? ModbusFrames.Tcp(1, 1, Pdu(response))
                : protocol == "modbus-rtu" ? ModbusFrames.Rtu(1, Pdu(response))
                : System.Text.Encoding.ASCII.GetBytes(":0103021234B4\r\n");
            var server = Task.Run(async () =>
            {
                using var client = await listener.AcceptTcpClientAsync(deadline.Token);
                var stream = client.GetStream();
                var bytes = new byte[expected.Length];
                await stream.ReadExactlyAsync(bytes, deadline.Token);
                Assert.Equal(expected, bytes);
                await stream.WriteAsync(reply, deadline.Token);
            });
            var manager = new org.apache.plc4net.spi.transports.DefaultTransportManager();
            org.apache.plc4net.spi.drivers.DriverBase driver = protocol == "modbus-tcp"
                ? new ModbusTcpDriver(manager) : protocol == "modbus-rtu"
                    ? new ModbusRtuDriver(manager) : new ModbusAsciiDriver(manager);
            using var connection = driver.Connect($"{protocol}:tcp://127.0.0.1:{((System.Net.IPEndPoint)listener.LocalEndpoint).Port}");
            Assert.Equal((ushort)0x1234, ValueOf(await ReadAsync(connection, "holding-register:1:UINT", deadline.Token)).GetUshort());
            await server;
        }

        [Theory]
        [InlineData("modbus-tcp")]
        [InlineData("modbus-rtu")]
        [InlineData("modbus-ascii")]
        public async Task Drivers_read_over_real_udp(string protocol)
        {
            using var server = new System.Net.Sockets.UdpClient(new System.Net.IPEndPoint(System.Net.IPAddress.Loopback, 0));
            using var deadline = new System.Threading.CancellationTokenSource(TimeSpan.FromSeconds(5));
            var expected = protocol == "modbus-tcp" ? ModbusFrames.Tcp(1, 1, Pdu("0300000001"))
                : protocol == "modbus-rtu" ? ModbusFrames.Rtu(1, Pdu("0300000001"))
                : System.Text.Encoding.ASCII.GetBytes(":010300000001FB\r\n");
            var response = protocol == "modbus-tcp" ? ModbusFrames.Tcp(1, 1, Pdu("03021234"))
                : protocol == "modbus-rtu" ? ModbusFrames.Rtu(1, Pdu("03021234"))
                : System.Text.Encoding.ASCII.GetBytes(":0103021234B4\r\n");
            var reply = Task.Run(async () =>
            {
                var request = await server.ReceiveAsync(deadline.Token);
                Assert.Equal(expected, request.Buffer);
                await server.SendAsync(response, request.RemoteEndPoint, deadline.Token);
            });
            var manager = new org.apache.plc4net.spi.transports.DefaultTransportManager();
            org.apache.plc4net.spi.drivers.DriverBase driver = protocol == "modbus-tcp"
                ? new ModbusTcpDriver(manager) : protocol == "modbus-rtu"
                    ? new ModbusRtuDriver(manager) : new ModbusAsciiDriver(manager);
            using var connection = driver.Connect($"{protocol}:udp://127.0.0.1:{((System.Net.IPEndPoint)server.Client.LocalEndPoint!).Port}");
            Assert.Equal((ushort)0x1234, ValueOf(await ReadAsync(connection, "holding-register:1:UINT", deadline.Token)).GetUshort());
            await reply;
        }

        [Fact]
        public async Task Bit_blocks_split_at_unaligned_offsets()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("01025501")));
            using var connection = NewTcp(transport);
            var builder = connection.ReadRequestBuilder!;
            builder.AddTagAddress("last", "coil:8[0..1]");
            builder.AddTagAddress("first", "coil:1[0..7]");
            var response = (DefaultPlcReadResponse)await ((PlcReader)connection).Read((DefaultPlcReadRequest)builder.Build());
            Assert.Equal(new[] { false, true }, response.GetValue("last")!.GetList().ConvertAll(v => v.GetBool()));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0100000009")), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Identification_reads_basic_objects()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("2B0E01010000010003414243")));
            using var connection = NewTcp(transport);
            var objects = await connection.IdentifyAsync();
            Assert.Equal(Pdu("414243"), objects[0]);
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("2B0E0100")), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Identification_rejects_non_advancing_pagination()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("2B0E0101FF0000")));
            using var connection = NewTcp(transport);
            await Assert.ThrowsAsync<ModbusDriverException>(() => connection.IdentifyAsync());
        }

        [Fact]
        public async Task Overlapping_tags_are_split_by_address_after_sorting()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("030400010002")));
            using var connection = NewTcp(transport);
            var builder = connection.ReadRequestBuilder!;
            builder.AddTagAddress("last", "holding-register:2:INT");
            builder.AddTagAddress("both", "holding-register:1[0..1]:INT");
            var response = (DefaultPlcReadResponse)await ((PlcReader)connection)
                .Read((DefaultPlcReadRequest)builder.Build());
            Assert.Equal((short)2, response.GetValue("last")!.GetShort());
            Assert.Equal(new short[] { 1, 2 }, response.GetValue("both")!.GetList().ConvertAll(v => v.GetShort()));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000002")), transport.GetAllWrittenData());
        }

        [Fact]
        public void Extended_selection_cannot_cross_the_final_file()
        {
            Assert.Equal(655349999, ModbusTag.Parse("extended-register:655349999:INT").Address);
            Assert.Throws<org.apache.plc4net.exceptions.PlcInvalidFieldException>(
                () => ModbusTag.Parse("extended-register:655349999:DINT"));
        }

        [Fact]
        public async Task Ping_reads_the_configured_address_and_unit()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 3, Pdu("03020001")));
            using var connection = NewTcp(transport,
                "?ping-address=holding-register:2:INT{unit-id:3}");
            await connection.PingAsync();
            Assert.Equal(ModbusFrames.Tcp(1, 3, Pdu("0300010001")), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Ping_times_out_when_the_device_does_not_answer()
        {
            using var connection = NewTcp(NewTransport(), "?request-timeout-ms=50");
            await Assert.ThrowsAsync<TimeoutException>(() => connection.PingAsync());
        }

        [Fact]
        public async Task Ping_propagates_device_errors()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("8302")));
            using var connection = NewTcp(transport);
            await Assert.ThrowsAsync<ModbusDriverException>(() => connection.PingAsync());
        }

        [Theory]
        // Read holding-register:3:WORD   (request 000500000006 01 0300020001, answer 000500000005 01 0302A5B8)
        [InlineData("holding-register:3:WORD", "0300020001", "0302A5B8", 42424)]
        // Read holding-register:12:INT   (-2424)
        [InlineData("holding-register:12:INT", "03000B0001", "0302F688", -2424)]
        // Read holding-register:13:UINT
        [InlineData("holding-register:13:UINT", "03000C0001", "0302A5B8", 42424)]
        public async Task A_read_puts_the_same_pdu_on_the_wire_and_decodes_the_same_answer(
            string tag, string requestPdu, string answerPdu, int expected)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu(answerPdu)));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, tag);

            var value = ValueOf(response);
            var actual = expected < 0 ? value.GetShort() : (long)value.GetUshort();
            Assert.Equal(expected, actual);
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu(requestPdu)), transport.GetAllWrittenData());
        }

        [Theory]
        // Write holding-register:3:WORD  (request 00060000 0006 01 060002A5B8, the answer echoes it)
        [InlineData("holding-register:3:WORD", 42424, "060002A5B8")]
        // Write holding-register:12:INT
        [InlineData("holding-register:12:INT", -2424, "06000BF688")]
        // Write holding-register:13:UINT
        [InlineData("holding-register:13:UINT", 42424, "06000CA5B8")]
        public async Task A_write_puts_the_same_pdu_on_the_wire(string tag, int value, string requestPdu)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu(requestPdu)));
            var connection = NewTcp(transport);

            var response = await WriteAsync(connection, tag, value);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu(requestPdu)), transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("BIG_ENDIAN", "11223344")]
        [InlineData("LITTLE_ENDIAN", "44332211")]
        [InlineData("BIG_ENDIAN_BYTE_SWAP", "22114433")]
        [InlineData("LITTLE_ENDIAN_BYTE_SWAP", "33441122")]
        public async Task A_wide_value_is_read_with_the_same_byte_orders_as_plc4j(
            string byteOrder, string wireValue)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("0304" + wireValue)));
            var connection = NewTcp(transport);

            var response = await ReadAsync(
                connection, $"holding-register:1:DWORD{{byte-order:'{byteOrder}'}}");

            Assert.Equal(0x11223344U, ValueOf(response).GetUint());
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000002")), transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("BIG_ENDIAN", "11223344")]
        [InlineData("LITTLE_ENDIAN", "44332211")]
        [InlineData("BIG_ENDIAN_BYTE_SWAP", "22114433")]
        [InlineData("LITTLE_ENDIAN_BYTE_SWAP", "33441122")]
        public async Task A_wide_value_is_written_with_the_same_byte_orders_as_plc4j(
            string byteOrder, string wireValue)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("1000000002")));
            var connection = NewTcp(transport);

            var response = await WriteAsync(
                connection, $"holding-register:1:DWORD{{byte-order:'{byteOrder}'}}", 0x11223344U);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(
                ModbusFrames.Tcp(1, 1, Pdu("100000000204" + wireValue)), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task A_register_array_is_read_in_one_request()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("0306000100020003")));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, "holding-register:1[0..2]:INT");

            Assert.Equal(new short[] { 1, 2, 3 }, ValueOf(response).GetList().ConvertAll(v => v.GetShort()));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000003")), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task A_register_array_is_written_in_one_request()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("1000000003")));
            var connection = NewTcp(transport);

            var response = await WriteAsync(
                connection, "holding-register:1[0..2]:INT", new List<short> { 1, 2, 3 });

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("100000000306000100020003")), transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("LITTLE_ENDIAN", "4433221188776655")]
        [InlineData("LITTLE_ENDIAN_BYTE_SWAP", "3344112277885566")]
        public async Task Byte_order_is_applied_to_each_array_element(
            string byteOrder, string wireValue)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("0308" + wireValue)));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection,
                $"holding-register:1[0..1]:DWORD{{byte-order:'{byteOrder}'}}");

            Assert.Equal(new uint[] { 0x11223344, 0x55667788 },
                ValueOf(response).GetList().ConvertAll(v => v.GetUint()));
        }

        [Fact]
        public async Task Public_builder_params_write_a_register_array()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("1000000003")));
            var connection = NewTcp(transport);
            var builder = connection.WriteRequestBuilder!;
            builder.AddTag("v", "holding-register:1[0..2]:INT", (short)1, (short)2, (short)3);

            var response = (DefaultPlcWriteResponse)await ((PlcWriter)connection)
                .Write((DefaultPlcWriteRequest)builder.Build());

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("100000000306000100020003")),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Connection_default_byte_order_and_tag_unit_id_are_used()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 3, Pdu("030444332211")));
            var connection = NewTcp(transport, "?default-payload-byte-order=LITTLE_ENDIAN");

            var response = await ReadAsync(connection, "holding-register:1:DWORD{unit-id:3}");

            Assert.Equal(0x11223344U, ValueOf(response).GetUint());
            Assert.Equal(ModbusFrames.Tcp(1, 3, Pdu("0300000002")), transport.GetAllWrittenData());
        }

        [Fact]
        public void Extended_register_keeps_a_wide_flat_address()
        {
            var tag = ModbusTag.Parse("extended-register:65536:INT");

            Assert.Equal(65536, tag.Address);
        }

        [Fact]
        public async Task A_125_register_write_is_split_at_the_123_register_pdu_limit()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("100000007B")));
            var secondResponse = Task.Run(async () =>
            {
                while (transport.GetNumBytesWritten() < 259) await Task.Delay(1);
                transport.InjectTestData(ModbusFrames.Tcp(2, 1, Pdu("10007B0002")));
            });
            var connection = NewTcp(transport);
            var values = new short[125];
            var builder = connection.WriteRequestBuilder!;
            builder.AddTag("v", "holding-register:1[0..124]:INT", values);

            var response = (DefaultPlcWriteResponse)await ((PlcWriter)connection)
                .Write((DefaultPlcWriteRequest)builder.Build());
            await secondResponse;

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(276, transport.GetAllWrittenData().Length);
        }

        [Fact]
        public async Task A_2000_coil_write_is_split_at_the_1968_coil_pdu_limit()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("0F000007B0")));
            var secondResponse = Task.Run(async () =>
            {
                while (transport.GetNumBytesWritten() < 259) await Task.Delay(1);
                transport.InjectTestData(ModbusFrames.Tcp(2, 1, Pdu("0F07B00020")));
            });
            var connection = NewTcp(transport);
            var values = new bool[2000];
            var builder = connection.WriteRequestBuilder!;
            builder.AddTag("v", "coil:1[0..1999]:BOOL", values);

            var response = (DefaultPlcWriteResponse)await ((PlcWriter)connection)
                .Write((DefaultPlcWriteRequest)builder.Build());
            await secondResponse;

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(276, transport.GetAllWrittenData().Length);
        }

        [Fact]
        public async Task A_fixed_length_string_is_read_from_its_registers()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("030441424344")));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, "holding-register:1:STRING(4)");

            Assert.Equal("ABCD", ValueOf(response).GetString());
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000002")), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task An_extended_register_uses_read_file_record()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("140403061234")));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, "extended-register:1:INT");

            Assert.Equal(0x1234, ValueOf(response).GetShort());
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("140706000100010001")), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Adjacent_register_tags_are_read_as_one_block()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("030400010002")));
            var connection = NewTcp(transport);
            var builder = connection.ReadRequestBuilder!;
            builder.AddTagAddress("a", "holding-register:1:INT");
            builder.AddTagAddress("b", "holding-register:2:INT");
            var request = (DefaultPlcReadRequest)builder.Build();

            var response = (DefaultPlcReadResponse)await ((PlcReader)connection).Read(request);

            Assert.Equal((short)1, response.GetValue("a")!.GetShort());
            Assert.Equal((short)2, response.GetValue("b")!.GetShort());
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000002")), transport.GetAllWrittenData());
        }
    }
}
