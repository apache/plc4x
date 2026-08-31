//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using org.apache.plc4net.api;
using org.apache.plc4net.api.authentication;
using org.apache.plc4net.exceptions;

namespace org.apache.plc4net
{
    /// <summary>
    /// Registry of PLC protocol drivers. A driver registers itself here;
    /// connection requests are dispatched to the driver whose protocol code
    /// matches the connection string scheme.
    /// </summary>
    /// <remarks>
    /// All read and write paths through the driver dictionary are safe for
    /// concurrent use so that a background initialisation task can register
    /// a driver without racing against an incoming connection request.
    /// </remarks>
    public class PlcDriverManager
    {
        // Static initializer is thread-safe; the CLR guarantees it runs exactly once.
        private static readonly PlcDriverManager _instance = new PlcDriverManager();

        // Protocol codes are matched case-insensitively, the same way
        // DefaultTransportManager matches transport codes.
        private readonly ConcurrentDictionary<string, IPlcDriver> _drivers
            = new ConcurrentDictionary<string, IPlcDriver>(StringComparer.OrdinalIgnoreCase);

        private PlcDriverManager()
        {
        }

        public static PlcDriverManager Instance => _instance;

        public void RegisterDriver(IPlcDriver driver)
        {
            if (driver == null)
            {
                throw new ArgumentNullException(nameof(driver));
            }
            _drivers[driver.ProtocolCode] = driver;
        }

        public IPlcDriver GetDriver(string url)
        {
            // Extract the protocol code from the connection string without relying
            // on System.Uri, which may reject or misinterpret non-standard schemes
            // like "s7:cotp". The pattern is always {protocol}(:{transport})?://...
            if (string.IsNullOrWhiteSpace(url))
                throw new PlcConnectionException("Connection string must not be null or empty.");

            var schemeEnd = url.IndexOf("://", StringComparison.Ordinal);
            if (schemeEnd < 0)
                throw new PlcConnectionException(
                    $"Invalid connection string '{url}': no :// separator found.");

            var scheme = schemeEnd > 0 ? url.Substring(0, schemeEnd) : string.Empty;
            if (scheme.Length == 0)
                throw new PlcConnectionException(
                    $"Invalid connection string '{url}': scheme is empty.");

            // For two-scheme forms like s7:cotp://host, take the protocol
            // portion before the colon (s7), matching what Java's
            // URI.getScheme() returns for the same input.
            var colon = scheme.IndexOf(':');
            if (colon >= 0)
                scheme = scheme.Substring(0, colon);

            return GetDriverByCode(scheme);
        }

        public IPlcDriver GetDriverByCode(string protocolCode)
        {
            if (protocolCode != null
                && _drivers.TryGetValue(protocolCode, out var driver)
                && driver != null)
            {
                return driver;
            }
            throw new PlcConnectionException(
                $"No driver registered for protocol '{protocolCode}'.");
        }

        public IPlcConnection GetConnection(string connectionString)
        {
            return GetConnection(connectionString, null);
        }

        public IPlcConnection GetConnection(string connectionString, IPlcAuthentication authentication)
        {
            var driver = GetDriver(connectionString);
            return driver.Connect(connectionString, authentication);
        }
    }
}
