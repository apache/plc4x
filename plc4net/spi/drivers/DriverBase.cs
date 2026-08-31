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
using System.Linq;
using System.Text.RegularExpressions;
using org.apache.plc4net.api;
using org.apache.plc4net.api.authentication;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.spi.drivers
{
    /// <summary>
    /// Shared logic for every driver: connection-string dispatch, transport selection,
    /// configuration parsing, supported-transport enforcement.
    /// </summary>
    /// <remarks>
    /// Mirrors the Java SPI3 <c>DriverBase</c>. A concrete driver implements
    /// <see cref="GetProtocolCode"/>, <see cref="GetProtocolName"/>, and
    /// <see cref="CreateConnection"/>; everything else is handled here.
    /// </remarks>
    public abstract class DriverBase : IPlcDriver
    {
        private readonly ITransportManager _transportManager;

        protected DriverBase(ITransportManager transportManager)
        {
            _transportManager = transportManager ?? throw new ArgumentNullException(nameof(transportManager));
        }

        protected void RegisterTransport(ITransport transport)
        {
            if (_transportManager is DefaultTransportManager defaultManager)
            {
                defaultManager.Register(transport);
            }
        }

        // ----- abstract members a driver must implement -----

        /// <summary>e.g. "s7", "modbus-tcp".</summary>
        public abstract string ProtocolCode { get; }

        /// <summary>Human-readable name, e.g. "Siemens S7".</summary>
        public abstract string ProtocolName { get; }

        /// <summary>
        /// Optional default transport when the connection string omits one
        /// (e.g. "s7" defaults to "cotp").
        /// </summary>
        public abstract string DefaultTransportCode { get; }

        /// <summary>
        /// The transports this driver supports. A connection string that requests a
        /// different transport is rejected unless <c>allow-unsupported-transport=true</c>
        /// is set.
        /// </summary>
        protected abstract string[] SupportedTransportCodes { get; }

        /// <summary>
        /// Constructs a protocol-specific connection. Called after the transport has
        /// been opened and validated. Implementations record configuration only; the
        /// protocol handshake belongs in <see cref="ConnectionBase.OnConnect"/>, which
        /// <see cref="Connect(string, IPlcAuthentication)"/> invokes on the result.
        /// </summary>
        protected abstract ConnectionBase CreateConnection(
            ConnectionString connectionString,
            ITransportInstance transportInstance,
            IPlcAuthentication authentication);

        // ----- public driver API -----

        public IPlcConnection Connect(string connectionString)
        {
            return Connect(connectionString, null);
        }

        public IPlcConnection Connect(string connectionString, IPlcAuthentication authentication)
        {
            if (string.IsNullOrWhiteSpace(connectionString))
            {
                throw new PlcConnectionException("Connection string must not be empty.");
            }

            var parsed = ConnectionString.Parse(connectionString);

            if (!string.Equals(parsed.ProtocolCode, ProtocolCode, StringComparison.OrdinalIgnoreCase))
            {
                throw new PlcConnectionException(
                    $"This driver ({ProtocolCode}) is not suited for protocol '{parsed.ProtocolCode}'.");
            }

            var transportCode = ResolveTransportCode(parsed);

            var transport = _transportManager.GetTransport(transportCode);
            if (transport == null)
            {
                throw new PlcConnectionException(
                    $"No transport registered for code '{transportCode}'. " +
                    $"Registered transports: {string.Join(", ", _transportManager.GetTransportCodes())}");
            }

            var rawTransportConfig = parsed.TransportConfig;
            // Per SPI3, the driver-config portion of the URL is the part after the
            // transport's parsed host:port. The TCP transport returns it as /milo in
            // opcua:tcp://localhost:4840/milo; for drivers that don't consume one it
            // stays empty.
            var transportConfig = transport.CreateConfiguration(parsed.Parameters);
            var transportInstance = transport.CreateTransportInstance(rawTransportConfig, transportConfig);

            // The transport is already open by this point, so a driver constructor that
            // throws (e.g. parameter validation) must not leak it: nothing else holds a
            // reference. Closing the instance is chain-safe — a wrapper transport (COTP
            // over TCP) closes its inner one.
            ConnectionBase connection;
            try
            {
                connection = CreateConnection(parsed, transportInstance, authentication);
            }
            catch
            {
                try
                {
                    transportInstance.Close();
                }
                catch
                {
                    // Preserve the original failure; cleanup errors must not mask it.
                }
                throw;
            }

            // Establish the protocol-level session. Construction above only records
            // configuration; anything that talks to the device happens here, so that a
            // driver's constructor never blocks on I/O. If setup fails the connection is
            // closed before the exception escapes — the caller never receives a reference,
            // so nobody else can release the socket.
            //
            // Closing goes through the connection rather than the transport instance
            // created above: a driver may wrap it (S7 puts COTP on top of TCP) and it is
            // the wrapper the connection holds. Closing the inner one directly would leave
            // the wrapper's own state untouched.
            try
            {
                connection.Connect();
            }
            catch (PlcConnectionException)
            {
                try
                {
                    connection.Close();
                }
                catch
                {
                    // Preserve the original failure; cleanup errors must not mask it.
                }
                throw;
            }
            catch (Exception ex)
            {
                try
                {
                    connection.Close();
                }
                catch
                {
                    // Preserve the original failure; cleanup errors must not mask it.
                }

                // Every other Connect() failure path reports PlcConnectionException;
                // a handshake failure that surfaces as TransportException must not
                // escape that contract. Java's SPI3 wraps both getConnection() and
                // connect() failures the same way.
                throw new PlcConnectionException(
                    $"Failed to establish the connection: {ex.Message}", ex);
            }

            return connection;
        }

        /// <summary>
        /// Async equivalent of <see cref="Connect(string, IPlcAuthentication)"/>.
        /// Runs the transport creation and protocol handshake on a background thread
        /// so that UI threads and ASP.NET request threads are not blocked.
        /// </summary>
        public async System.Threading.Tasks.Task<IPlcConnection> ConnectAsync(
            string connectionString,
            IPlcAuthentication authentication = null)
        {
            var parsed = ConnectionString.Parse(connectionString);

            if (!string.Equals(parsed.ProtocolCode, ProtocolCode, StringComparison.OrdinalIgnoreCase))
            {
                throw new PlcConnectionException(
                    $"This driver ({ProtocolCode}) is not suited for protocol '{parsed.ProtocolCode}'.");
            }

            var transportCode = ResolveTransportCode(parsed);

            var transport = _transportManager.GetTransport(transportCode);
            if (transport == null)
            {
                throw new PlcConnectionException(
                    $"No transport registered for code '{transportCode}'. " +
                    $"Registered transports: {string.Join(", ", _transportManager.GetTransportCodes())}");
            }

            var transportConfig = transport.CreateConfiguration(parsed.Parameters);
            var transportInstance = transport.CreateTransportInstance(
                parsed.TransportConfig, transportConfig);

            ConnectionBase connection;
            try
            {
                connection = CreateConnection(parsed, transportInstance, authentication);
            }
            catch
            {
                try { transportInstance.Close(); }
                catch { }
                throw;
            }

            try
            {
                await connection.ConnectAsync().ConfigureAwait(false);
            }
            catch (PlcConnectionException)
            {
                try { connection.Close(); }
                catch { }
                throw;
            }
            catch (Exception ex)
            {
                try { connection.Close(); }
                catch { }
                throw new PlcConnectionException(
                    $"Failed to establish the connection: {ex.Message}", ex);
            }

            return connection;
        }

        private string ResolveTransportCode(ConnectionString parsed)
        {
            if (!string.IsNullOrWhiteSpace(parsed.TransportCode))
            {
                var requested = parsed.TransportCode;
                if (!SupportedTransportCodes.Contains(requested, StringComparer.OrdinalIgnoreCase))
                {
                    var allowUnsupported = parsed.GetBoolParameter("allow-unsupported-transport", false);
                    if (!allowUnsupported)
                    {
                        throw new PlcConnectionException(
                            $"Transport '{requested}' is not supported by driver '{ProtocolCode}'. " +
                            $"Supported transports: {string.Join(", ", SupportedTransportCodes)}. " +
                            "Set 'allow-unsupported-transport=true' to bypass this check.",
                            new Exception("Unsupported transport"));
                    }
                }
                return requested;
            }

            if (!string.IsNullOrWhiteSpace(DefaultTransportCode))
            {
                return DefaultTransportCode;
            }

            throw new PlcConnectionException(
                $"Driver '{ProtocolCode}' has no default transport and the connection string did not specify one.");
        }
    }
}
