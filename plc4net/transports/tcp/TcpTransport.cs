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
using System.Net;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.transports.tcp
{
    /// <summary>
    /// The "tcp" transport.
    /// </summary>
    public class TcpTransport : ITransport
    {
        public string TransportCode => "tcp";

        public string TransportName => "TCP/IP Socket Transport";

        public ITransportConfiguration CreateConfiguration(IReadOnlyDictionary<string, string> parameters)
        {
            var configuration = new TcpTransportConfiguration();
            if (parameters == null)
            {
                return configuration;
            }

            // Transport options are prefixed with the transport code in the connection
            // string, matching the Java side: ConfigurationFactory builds the prefix from
            // the transport code and strips it before matching the option name, so the
            // option Java declares as "tcp-no-delay" is written "tcp.tcp-no-delay" and
            // "connect-timeout" is written "tcp.connect-timeout". The unprefixed form is
            // accepted too, since drivers with a single transport commonly use it.
            configuration.ConnectTimeout = GetInt(parameters, "connect-timeout", configuration.ConnectTimeout);
            configuration.TcpNoDelay = GetBool(parameters, "tcp-no-delay", configuration.TcpNoDelay);
            configuration.KeepAlive = GetBool(parameters, "keep-alive", configuration.KeepAlive);
            configuration.SendBufferSize = GetInt(parameters, "send-buffer-size", configuration.SendBufferSize);
            configuration.ReceiveBufferSize = GetInt(parameters, "receive-buffer-size", configuration.ReceiveBufferSize);
            configuration.LocalAddress = GetString(parameters, "local-address", configuration.LocalAddress);
            configuration.LocalPort = GetInt(parameters, "local-port", configuration.LocalPort);
            configuration.DefaultPort = GetInt(parameters, "default-port", configuration.DefaultPort);

            // A zero or negative receive buffer size would crash the RingBuffer
            // constructor, which requires a positive capacity. Fall back to the
            // default in that case.
            if (configuration.ReceiveBufferSize <= 0)
            {
                configuration.ReceiveBufferSize = new TcpTransportConfiguration().ReceiveBufferSize;
            }

            return configuration;
        }

        public ITransportInstance CreateTransportInstance(
            string transportConfig, ITransportConfiguration configuration)
        {
            if (!(configuration is TcpTransportConfiguration tcpConfiguration))
            {
                throw new TransportException(
                    $"Expected a {nameof(TcpTransportConfiguration)} but got {configuration?.GetType().Name ?? "null"}.");
            }

            var (host, port, driverConfig) = ParseAddress(transportConfig, tcpConfiguration.DefaultPort);
            var endPoint = new IPEndPoint(ResolveHost(host), port);

            var instance = new TcpTransportInstance(endPoint, tcpConfiguration);
            instance.SetDriverConfig(driverConfig);
            return instance;
        }

        /// <summary>
        /// Splits "host:port/driver-config" into its parts. The port may be omitted when
        /// the driver supplies a default. Bracketed IPv6 literals are supported.
        /// </summary>
        internal static (string Host, int Port, string DriverConfig) ParseAddress(string address, int defaultPort)
        {
            if (string.IsNullOrWhiteSpace(address))
            {
                throw new TransportException("Transport address must not be empty.");
            }

            var driverConfig = string.Empty;
            var slash = address.IndexOf('/');
            if (slash >= 0)
            {
                driverConfig = address.Substring(slash);
                address = address.Substring(0, slash);
            }

            string host;
            int port;

            if (address.StartsWith("[", StringComparison.Ordinal))
            {
                // IPv6 literal, e.g. [::1]:502
                var closing = address.IndexOf(']');
                if (closing < 0)
                {
                    throw new TransportException($"Malformed IPv6 address '{address}'.");
                }
                host = address.Substring(1, closing - 1);
                var rest = address.Substring(closing + 1);
                port = rest.StartsWith(":", StringComparison.Ordinal)
                    ? ParsePort(rest.Substring(1))
                    : RequireDefaultPort(defaultPort, address);
            }
            else
            {
                var parts = address.Split(':');
                if (parts.Length == 1)
                {
                    host = parts[0];
                    port = RequireDefaultPort(defaultPort, address);
                }
                else if (parts.Length == 2)
                {
                    host = parts[0];
                    port = ParsePort(parts[1]);
                }
                else
                {
                    throw new TransportException(
                        $"Malformed address '{address}'. Use host, host:port, or [ipv6]:port.");
                }
            }

            if (string.IsNullOrWhiteSpace(host))
            {
                throw new TransportException($"Address '{address}' has no host part.");
            }

            return (host, port, driverConfig);
        }

        private static int RequireDefaultPort(int defaultPort, string address)
        {
            if (defaultPort <= 0)
            {
                throw new TransportException(
                    $"Address '{address}' has no port and the driver declares no default port.");
            }
            return defaultPort;
        }

        private static int ParsePort(string raw)
        {
            if (!int.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture, out var port)
                || port <= 0 || port > 65535)
            {
                throw new TransportException($"'{raw}' is not a valid port number.");
            }
            return port;
        }

        private static IPAddress ResolveHost(string host)
        {
            if (IPAddress.TryParse(host, out var address))
            {
                return address;
            }
            try
            {
                var resolved = Dns.GetHostAddresses(host).FirstOrDefault();
                if (resolved == null)
                {
                    throw new TransportException($"Host '{host}' did not resolve to any address.");
                }
                return resolved;
            }
            catch (Exception e) when (!(e is TransportException))
            {
                throw new TransportException($"Unable to resolve host '{host}'.", e);
            }
        }

        private static string GetString(IReadOnlyDictionary<string, string> p, string key, string fallback)
        {
            if (p.TryGetValue("tcp." + key, out var prefixed))
            {
                return prefixed;
            }
            return p.TryGetValue(key, out var plain) ? plain : fallback;
        }

        private static int GetInt(IReadOnlyDictionary<string, string> p, string key, int fallback)
        {
            var raw = GetString(p, key, null);
            return int.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture, out var value)
                ? value
                : fallback;
        }

        private static bool GetBool(IReadOnlyDictionary<string, string> p, string key, bool fallback)
        {
            var raw = GetString(p, key, null);
            return bool.TryParse(raw, out var value) ? value : fallback;
        }
    }
}
