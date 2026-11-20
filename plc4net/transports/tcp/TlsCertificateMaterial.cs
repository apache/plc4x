// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
// https://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied. See the License for the
// specific language governing permissions and limitations
// under the License.

using System;
using System.Net.Security;
using System.Security.Cryptography.X509Certificates;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.transports.tcp
{
    internal sealed class TlsCertificateMaterial : IDisposable
    {
        private readonly X509Certificate2Collection _owned = new X509Certificate2Collection();
        public X509ChainPolicy? ChainPolicy { get; private set; }
        public X509CertificateCollection ClientCertificates { get; } = new X509CertificateCollection();
        public SslStreamCertificateContext? ClientCertificateContext { get; private set; }

        public TlsCertificateMaterial(TcpTransportConfiguration configuration)
        {
            try
            {
                ChainPolicy = configuration.TlsChainPolicy;
                if (configuration.TlsTrustStore != null)
                {
                    if (ChainPolicy != null)
                        throw new TransportException("Specify either a TLS trust store or a chain policy.");
                    var roots = new X509Certificate2Collection();
                    roots.ImportFromPemFile(configuration.TlsTrustStore);
                    _owned.AddRange(roots);
                    if (roots.Count == 0) throw new TransportException("The TLS trust store contains no certificates.");
                    ChainPolicy = new X509ChainPolicy
                    {
                        TrustMode = X509ChainTrustMode.CustomRootTrust,
                        RevocationMode = X509RevocationMode.Online
                    };
                    ChainPolicy.CustomTrustStore.AddRange(roots);
                }
                if (configuration.TlsClientCertificates != null)
                    ClientCertificates.AddRange(configuration.TlsClientCertificates);
                if (configuration.TlsKeyStore != null)
                {
                    if (ClientCertificates.Count != 0)
                        throw new TransportException("Specify either a TLS client keystore or client certificates.");
                    var imported = new X509Certificate2Collection();
                    var password = configuration.TlsKeyStorePasswordProvider?.Invoke();
                    try
                    {
                        imported.Import(configuration.TlsKeyStore,
                            password,
                            // macOS does not support ephemeral key imports; use its temporary keychain import.
                            OperatingSystem.IsLinux() ? X509KeyStorageFlags.EphemeralKeySet : X509KeyStorageFlags.DefaultKeySet);
                    }
                    finally { _owned.AddRange(imported); }
                    var identities = 0;
                    foreach (var certificate in imported)
                        if (certificate.HasPrivateKey) identities++;
                    if (identities != 1)
                        throw new TransportException("Specify exactly one TLS client identity with a private key.");
                    if (OperatingSystem.IsWindows())
                    {
                        foreach (var certificate in imported)
                            if (!certificate.HasPrivateKey) ClientCertificates.Add(certificate);
                        // Preserve the single-certificate loader's Schannel-compatible key handle.
                        var selected = new X509Certificate2(configuration.TlsKeyStore, password);
                        _owned.Add(selected);
                        ClientCertificates.Add(selected);
                    }
                    else ClientCertificates.AddRange(imported);
                }
                X509Certificate2? identity = null;
                var additional = new X509Certificate2Collection();
                foreach (X509Certificate certificate in ClientCertificates)
                {
                    if (certificate is not X509Certificate2 client)
                        throw new TransportException("TLS client certificates must be X509Certificate2 instances.");
                    if (client.HasPrivateKey)
                    {
                        if (identity != null)
                            throw new TransportException("Specify exactly one TLS client identity with a private key.");
                        identity = client;
                    }
                    else additional.Add(client);
                }
                if (ClientCertificates.Count != 0)
                {
                    if (identity == null)
                        throw new TransportException("TLS client certificates must include a private key.");
                    // Supply the chain explicitly; importing only the identity loses PFX intermediates.
                    ClientCertificateContext = SslStreamCertificateContext.Create(identity, additional, offline: true);
                }
            }
            catch
            {
                Dispose();
                throw;
            }
        }

        public void Dispose()
        {
            foreach (var certificate in _owned) certificate.Dispose();
            _owned.Clear();
        }
    }
}
