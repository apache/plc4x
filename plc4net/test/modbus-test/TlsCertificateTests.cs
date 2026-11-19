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
using System.Net.Security;
using System.Net.Sockets;
using System.Security.Authentication;
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.tcp;
using Xunit;
using static org.apache.plc4net.drivers.modbus.test.ConnectionTestSupport;

namespace org.apache.plc4net.drivers.modbus.test
{
    public class TlsCertificateTests
    {
        private static X509Certificate2 Certificate(string name, string usage)
        {
            using var key = RSA.Create(2048);
            var request = new CertificateRequest("CN=" + name, key, HashAlgorithmName.SHA256, RSASignaturePadding.Pkcs1);
            request.CertificateExtensions.Add(new X509BasicConstraintsExtension(true, false, 0, true));
            request.CertificateExtensions.Add(new X509EnhancedKeyUsageExtension(new OidCollection { new Oid(usage) }, true));
            var names = new SubjectAlternativeNameBuilder();
            names.AddDnsName(name);
            request.CertificateExtensions.Add(names.Build());
            using var original = request.CreateSelfSigned(DateTimeOffset.UtcNow.AddDays(-1), DateTimeOffset.UtcNow.AddDays(1));
            return new X509Certificate2(original.Export(X509ContentType.Pfx), (string?)null, X509KeyStorageFlags.Exportable);
        }

        [Theory]
        [InlineData(true, true, "localhost", false)]
        [InlineData(false, true, "localhost", false)]
        [InlineData(true, false, "localhost", false)]
        [InlineData(true, true, "127.0.0.1", false)]
        [InlineData(true, true, "localhost", true)]
        public async Task Private_trust_and_client_keystore_work_through_driver_url(
            bool supplyClient, bool trustServer, string hostname, bool programmatic)
        {
            var succeeds = supplyClient && trustServer && hostname == "localhost";
            using var serverCertificate = Certificate("localhost", "1.3.6.1.5.5.7.3.1");
            using var clientCertificate = Certificate("client", "1.3.6.1.5.5.7.3.2");
            var directory = Path.Combine(Path.GetTempPath(), "plc4net-tls-" + Guid.NewGuid().ToString("N"));
            Directory.CreateDirectory(directory);
            var trustFile = Path.Combine(directory, "trust.pem");
            var keyFile = Path.Combine(directory, "client.pfx");
            var passwordVariable = "PLC4NET_TLS_TEST_" + Guid.NewGuid().ToString("N");
            var password = Convert.ToBase64String(RandomNumberGenerator.GetBytes(24));
            try
            {
                File.WriteAllText(trustFile, (trustServer ? serverCertificate : clientCertificate).ExportCertificatePem());
                File.WriteAllBytes(keyFile, clientCertificate.Export(X509ContentType.Pfx, password));
                Environment.SetEnvironmentVariable(passwordVariable, password);
                using var listener = new TcpListener(IPAddress.Loopback, 0);
                listener.Start();
                using var deadline = new CancellationTokenSource(TimeSpan.FromSeconds(10));
                var server = Task.Run(async () =>
                {
                    using var peer = await listener.AcceptTcpClientAsync(deadline.Token);
                    using var stream = new SslStream(peer.GetStream(), false, (_, certificate, _, _) =>
                        certificate != null && certificate.GetCertHashString() == clientCertificate.Thumbprint);
                    try
                    {
                        await stream.AuthenticateAsServerAsync(new SslServerAuthenticationOptions
                        {
                            ServerCertificate = serverCertificate,
                            ClientCertificateRequired = true,
                            EnabledSslProtocols = SslProtocols.Tls12,
                            CertificateRevocationCheckMode = X509RevocationMode.NoCheck
                        }, deadline.Token);
                        Assert.True(stream.IsMutuallyAuthenticated);
                        if (!succeeds)
                        {
                            // The server can complete its side before receiving the client's alert.
                            Assert.Equal(0, await stream.ReadAsync(new byte[1].AsMemory(), deadline.Token));
                            return;
                        }
                        var request = new byte[12];
                        await stream.ReadExactlyAsync(request, deadline.Token);
                        Assert.Equal(ModbusFrames.Tcp(1, 1, Convert.FromHexString("0300000001")), request);
                        await stream.WriteAsync(ModbusFrames.Tcp(1, 1, Convert.FromHexString("03021234")), deadline.Token);
                    }
                    catch (AuthenticationException) when (!succeeds) { }
                    catch (IOException) when (!succeeds) { }
                });
                var options = "?tls.trust-store=" + Uri.EscapeDataString(trustFile);
                if (supplyClient) options += "&tls.keystore=" + Uri.EscapeDataString(keyFile)
                    + "&tls.keystore-password-env=" + passwordVariable;
                var driver = new ModbusTcpDriver(new DefaultTransportManager());
                var address = $"modbus-tcp:tls://{hostname}:{((IPEndPoint)listener.LocalEndpoint).Port}" + options;
                if (programmatic)
                {
                    var config = new TcpTransportConfiguration
                    {
                        Tls = true,
                        TlsTrustStore = trustFile,
                        TlsClientCertificates = new X509CertificateCollection { clientCertificate }
                    };
                    using (var transport = new TlsTransport().CreateTransportInstance(
                        $"localhost:{((IPEndPoint)listener.LocalEndpoint).Port}", config))
                    using (var connection = new ModbusTcpConnection(
                        org.apache.plc4net.spi.drivers.ConnectionString.Parse("modbus-tcp://localhost"), transport))
                    {
                        Assert.Equal((ushort)0x1234, ValueOf(await ReadAsync(connection, "holding-register:1:UINT", deadline.Token)).GetUshort());
                    }
                    Assert.True(clientCertificate.HasPrivateKey);
                    _ = clientCertificate.GetCertHashString();
                }
                else if (succeeds)
                {
                    using var connection = driver.Connect(address);
                    Assert.Equal((ushort)0x1234, ValueOf(await ReadAsync(connection, "holding-register:1:UINT", deadline.Token)).GetUshort());
                }
                else
                {
                    try
                    {
                        using var connection = driver.Connect(address);
                        var response = await ReadAsync(connection, "holding-register:1:UINT", deadline.Token);
                        Assert.NotEqual(org.apache.plc4net.types.PlcResponseCode.Ok, response.GetResponseCode("v"));
                    }
                    catch (TransportException) { }
                }
                await server;
            }
            finally
            {
                Environment.SetEnvironmentVariable(passwordVariable, null);
                File.Delete(keyFile);
                File.Delete(trustFile);
                Directory.Delete(directory);
            }
        }

        [Theory]
        [InlineData("tls.keystore-type", "JKS")]
        [InlineData("tls.trust-store-type", "JKS")]
        [InlineData("tls.keystore-password-env", "UNSET")]
        public void Invalid_certificate_configuration_is_rejected(string key, string value)
        {
            var transport = new TlsTransport();
            Assert.Throws<TransportException>(() => transport.CreateConfiguration(new Dictionary<string, string> { [key] = value }));
        }

        [Theory]
        [InlineData("wrong-password")]
        [InlineData("no-private-key")]
        [InlineData("empty-trust")]
        [InlineData("missing-file")]
        public void Unusable_certificate_material_fails_before_authentication(string scenario)
        {
            var path = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString("N"));
            using var certificate = Certificate("client", "1.3.6.1.5.5.7.3.2");
            var password = Convert.ToBase64String(RandomNumberGenerator.GetBytes(24));
            var config = new TcpTransportConfiguration { Tls = true };
            try
            {
                if (scenario == "wrong-password")
                {
                    File.WriteAllBytes(path, certificate.Export(X509ContentType.Pfx, password));
                    config.TlsKeyStore = path;
                    config.TlsKeyStorePasswordProvider = () => Convert.ToBase64String(RandomNumberGenerator.GetBytes(24));
                }
                else if (scenario == "no-private-key")
                {
                    File.WriteAllBytes(path, certificate.Export(X509ContentType.Cert));
                    config.TlsKeyStore = path;
                }
                else
                {
                    if (scenario == "empty-trust") File.WriteAllText(path, string.Empty);
                    config.TlsTrustStore = path;
                }
                using var listener = new TcpListener(IPAddress.Loopback, 0);
                listener.Start();
                Assert.Throws<TransportException>(() => new TcpTransportInstance((IPEndPoint)listener.LocalEndpoint, config));
                Assert.True(certificate.HasPrivateKey);
            }
            finally { File.Delete(path); }
        }
    }
}
