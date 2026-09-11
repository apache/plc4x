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

namespace org.apache.plc4net.spi.transports
{
    /// <summary>
    /// Factory for a kind of transport (tcp, udp, serial, cotp, ...).
    /// </summary>
    public interface ITransport
    {
        /// <summary>The code used in a connection string to select this transport, e.g. "tcp".</summary>
        string TransportCode { get; }

        /// <summary>Human-readable name.</summary>
        string TransportName { get; }

        /// <summary>
        /// Builds a configuration for this transport from the connection string's
        /// parameters, applying this transport's defaults.
        /// </summary>
        ITransportConfiguration CreateConfiguration(IReadOnlyDictionary<string, string> parameters);

        /// <summary>
        /// Opens a connection.
        /// </summary>
        /// <param name="transportConfig">the address portion of the URL, e.g. "192.168.0.1:102"</param>
        /// <param name="configuration">configuration produced by <see cref="CreateConfiguration"/></param>
        ITransportInstance CreateTransportInstance(string transportConfig, ITransportConfiguration configuration);
    }

    /// <summary>Registry of the transports available to drivers.</summary>
    public interface ITransportManager
    {
        /// <summary>Returns the transport for a code, or null when none is registered.</summary>
        ITransport GetTransport(string transportCode);

        IReadOnlyCollection<string> GetTransportCodes();
    }

    /// <summary>
    /// In-memory <see cref="ITransportManager"/>. Transports register themselves here
    /// rather than being discovered, which keeps the .NET port free of the assembly
    /// scanning the Java side does via ServiceLoader.
    /// </summary>
    public class DefaultTransportManager : ITransportManager
    {
        private readonly Dictionary<string, ITransport> _transports =
            new Dictionary<string, ITransport>(StringComparer.OrdinalIgnoreCase);

        public DefaultTransportManager()
        {
        }

        public DefaultTransportManager(IEnumerable<ITransport> transports)
        {
            if (transports == null)
            {
                return;
            }
            foreach (var transport in transports)
            {
                Register(transport);
            }
        }

        public void Register(ITransport transport)
        {
            if (transport == null)
            {
                throw new ArgumentNullException(nameof(transport));
            }
            _transports[transport.TransportCode] = transport;
        }

        public ITransport GetTransport(string transportCode)
        {
            if (transportCode == null)
            {
                return null;
            }
            return _transports.TryGetValue(transportCode, out var transport) ? transport : null;
        }

        public IReadOnlyCollection<string> GetTransportCodes()
        {
            return new List<string>(_transports.Keys);
        }
    }
}
