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
using org.apache.plc4net.api;
using org.apache.plc4net.api.authentication;
using org.apache.plc4net.api.metadata;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.transports;
using Xunit;

namespace org.apache.plc4net.spi.test.drivers
{
    /// <summary>
    /// A minimal transport for testing driver logic without a real network.
    /// </summary>
    internal class FakeTransport : ITransport
    {
        public FakeTransport(string transportCode = "fake")
        {
            TransportCode = transportCode;
        }

        public string TransportCode { get; }
        public string TransportName => "Fake Transport";

        public ITransportConfiguration CreateConfiguration(IReadOnlyDictionary<string, string> parameters)
        {
            return new FakeTransportConfiguration();
        }

        public ITransportInstance CreateTransportInstance(
            string transportConfig, ITransportConfiguration configuration)
        {
            return new FakeTransportInstance(configuration);
        }
    }

    internal class FakeTransportConfiguration : ITransportConfiguration
    {
    }

    internal class FakeTransportInstance : BaseTransportInstance
    {
        public FakeTransportInstance(ITransportConfiguration config) : base(config) { }
        public override bool IsOpen => true;
        public override int GetNumBytesAvailable() => 0;
        public override byte[] PeekReadableBytes(int numBytes) => Array.Empty<byte>();
        public override byte[] Read(int numBytes) => Array.Empty<byte>();
        public override void Write(byte[] bytes) { }
        public override void Close() { }
    }

    /// <summary>
    /// A stub connection returned by the test driver.
    /// </summary>
    internal class FakeConnection : ConnectionBase
    {
        public FakeConnection(ConnectionString cs, ITransportInstance transport)
            : base(cs, transport) { }

        public override IPlcConnectionMetadata PlcConnectionMetadata => null;
        public override IPlcTag Parse(string tagQuery) => null;
        public override IPlcReadRequestBuilder ReadRequestBuilder => null;
        public override IPlcWriteRequestBuilder WriteRequestBuilder => null;
        public override IPlcSubscriptionRequestBuilder SubscriptionRequestBuilder => null;
        public override IPlcUnsubscriptionRequestBuilder UnsubscriptionRequestBuilder => null;
    }

    /// <summary>
    /// A stub driver wired to the "fake" transport.
    /// </summary>
    internal class StubDriver : DriverBase
    {
        private readonly string _protocolCode;

        public StubDriver(ITransportManager transportManager, string protocolCode = "stub")
            : base(transportManager)
        {
            _protocolCode = protocolCode;
        }

        public override string ProtocolCode => _protocolCode;
        public override string ProtocolName => "Stub Driver";
        public override string DefaultTransportCode => "fake";
        protected override string[] SupportedTransportCodes => new[] { "fake" };

        protected override ConnectionBase CreateConnection(
            ConnectionString connectionString,
            ITransportInstance transportInstance,
            IPlcAuthentication authentication)
        {
            return new FakeConnection(connectionString, transportInstance);
        }
    }

    public class DriverBaseTests
    {
        [Fact]
        public void Connect_resolves_the_default_transport()
        {
            var transportManager = new DefaultTransportManager(new[] { new FakeTransport() });
            var driver = new StubDriver(transportManager);

            var connection = driver.Connect("stub://192.168.0.1");
            Assert.True(connection.IsConnected);
        }

        [Fact]
        public void Connect_rejects_an_unsupported_transport()
        {
            var transportManager = new DefaultTransportManager(new[] { new FakeTransport() });
            var driver = new StubDriver(transportManager);

            var ex = Assert.Throws<PlcConnectionException>(
                () => driver.Connect("stub:unregistered://host"));
            Assert.Contains("unregistered", ex.Message);
        }

        [Fact]
        public void Connect_rejects_a_registered_but_unsupported_transport()
        {
            // "other" is registered, so it exists — but StubDriver does not list it in
            // SupportedTransportCodes, so the driver's own check has to reject it.
            var transportManager = new DefaultTransportManager(
                new[] { new FakeTransport(), new FakeTransport("other") });
            var driver = new StubDriver(transportManager);

            var ex = Assert.Throws<PlcConnectionException>(
                () => driver.Connect("stub:other://host"));
            Assert.Contains("allow-unsupported-transport", ex.Message);
        }

        [Fact]
        public void Connect_allows_unsupported_transport_when_opted_in()
        {
            var transportManager = new DefaultTransportManager(
                new[] { new FakeTransport(), new FakeTransport("other") });
            var driver = new StubDriver(transportManager);

            // Same connection string as above, plus the opt-in flag.
            var connection = driver.Connect("stub:other://host?allow-unsupported-transport=true");
            Assert.True(connection.IsConnected);
        }

        [Fact]
        public void Connect_rejects_a_wrong_protocol_code()
        {
            var transportManager = new DefaultTransportManager(new[] { new FakeTransport() });
            var driver = new StubDriver(transportManager);

            Assert.Throws<PlcConnectionException>(
                () => driver.Connect("wrong-protocol://host"));
        }

        [Fact]
        public void Connect_matches_the_protocol_code_case_insensitively()
        {
            // Connection strings are lowercase by grammar, but a driver is free to
            // declare its code in any casing.
            var transportManager = new DefaultTransportManager(new[] { new FakeTransport() });
            var driver = new StubDriver(transportManager, "Stub");

            var connection = driver.Connect("stub://192.168.0.1");
            Assert.True(connection.IsConnected);
        }

        [Fact]
        public void Connect_runs_the_protocol_level_setup_hook()
        {
            var transportManager = new DefaultTransportManager(new[] { new FakeTransport() });
            var driver = new HookedDriver(transportManager);

            var connection = (HookedConnection)driver.Connect("hooked://192.168.0.1");

            // The hook must have run by the time Connect() returns, so callers never
            // receive a connection whose protocol session is still unestablished.
            Assert.Equal(1, connection.OnConnectCalls);
        }

        [Fact]
        public void Connect_closes_the_transport_when_the_setup_hook_fails()
        {
            var transportManager = new DefaultTransportManager(new[] { new FakeTransport() });
            var driver = new HookedDriver(transportManager, failOnConnect: true);

            Assert.Throws<PlcConnectionException>(() => driver.Connect("hooked://192.168.0.1"));

            // The caller never gets a reference to the connection, so if Connect() did not
            // close the transport itself the socket would stay bound with no way to reach it.
            Assert.True(driver.LastTransportInstance.Closed);
        }

        [Fact]
        public void Connect_wraps_a_non_connection_exception_and_still_closes()
        {
            var transportManager = new DefaultTransportManager(new[] { new FakeTransport() });
            // A handshake failure surfaces as TransportException, not PlcConnectionException
            // (this is exactly what a real S7 COTP timeout does). The driver API promises
            // PlcConnectionException on every failure path, so the wrap must kick in and
            // the transport must still be closed.
            var driver = new HookedDriver(transportManager, failOnConnect: true, failWithTransportException: true);

            var ex = Assert.Throws<PlcConnectionException>(() => driver.Connect("hooked://192.168.0.1"));
            Assert.IsType<TransportException>(ex.InnerException);
            Assert.True(driver.LastTransportInstance.Closed);
        }
    }

    /// <summary>
    /// A connection that records whether the protocol-level setup hook ran, and can be
    /// told to fail inside it.
    /// </summary>
    internal class HookedConnection : ConnectionBase
    {
        private readonly bool _failOnConnect;
        private readonly bool _failWithTransportException;

        public HookedConnection(ConnectionString cs, ITransportInstance transport,
            bool failOnConnect, bool failWithTransportException)
            : base(cs, transport)
        {
            _failOnConnect = failOnConnect;
            _failWithTransportException = failWithTransportException;
        }

        public int OnConnectCalls { get; private set; }

        protected override void OnConnect()
        {
            OnConnectCalls++;
            if (_failOnConnect)
            {
                throw _failWithTransportException
                    ? new TransportException("Simulated handshake failure.")
                    : new PlcConnectionException("Simulated handshake failure.");
            }
        }

        public override IPlcConnectionMetadata PlcConnectionMetadata => null;
        public override IPlcTag Parse(string tagQuery) => null;
        public override IPlcReadRequestBuilder ReadRequestBuilder => null;
        public override IPlcWriteRequestBuilder WriteRequestBuilder => null;
        public override IPlcSubscriptionRequestBuilder SubscriptionRequestBuilder => null;
        public override IPlcUnsubscriptionRequestBuilder UnsubscriptionRequestBuilder => null;
    }

    /// <summary>
    /// A transport instance that reports whether it was closed.
    /// </summary>
    internal class ClosableTransportInstance : BaseTransportInstance
    {
        public ClosableTransportInstance(ITransportConfiguration config) : base(config) { }

        public bool Closed { get; private set; }

        public override bool IsOpen => !Closed;
        public override int GetNumBytesAvailable() => 0;
        public override byte[] PeekReadableBytes(int numBytes) => Array.Empty<byte>();
        public override byte[] Read(int numBytes) => Array.Empty<byte>();
        public override void Write(byte[] bytes) { }
        public override void Close() => Closed = true;
    }

    internal class HookedDriver : DriverBase
    {
        private readonly bool _failOnConnect;
        private readonly bool _failWithTransportException;

        public HookedDriver(ITransportManager transportManager,
            bool failOnConnect = false, bool failWithTransportException = false)
            : base(transportManager)
        {
            _failOnConnect = failOnConnect;
            _failWithTransportException = failWithTransportException;
        }

        public ClosableTransportInstance LastTransportInstance { get; private set; }

        public override string ProtocolCode => "hooked";
        public override string ProtocolName => "Hooked Driver";
        public override string DefaultTransportCode => "fake";
        protected override string[] SupportedTransportCodes => new[] { "fake" };

        protected override ConnectionBase CreateConnection(
            ConnectionString connectionString,
            ITransportInstance transportInstance,
            IPlcAuthentication authentication)
        {
            // Swap in an instance that reports closure, so the failure path can be asserted.
            LastTransportInstance = new ClosableTransportInstance(transportInstance.Configuration);
            return new HookedConnection(connectionString, LastTransportInstance,
                _failOnConnect, _failWithTransportException);
        }
    }
}
