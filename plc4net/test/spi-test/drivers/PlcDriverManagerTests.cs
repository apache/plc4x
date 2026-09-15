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
using org.apache.plc4net.api;
using org.apache.plc4net.api.authentication;
using org.apache.plc4net.exceptions;
using Xunit;

namespace org.apache.plc4net.spi.test.drivers
{
    /// <summary>
    /// A driver that only carries a protocol code; connecting is out of scope here.
    /// </summary>
    internal class RegistryOnlyDriver : IPlcDriver
    {
        public RegistryOnlyDriver(string protocolCode)
        {
            ProtocolCode = protocolCode;
        }

        public string ProtocolCode { get; }

        public string ProtocolName => "Registry Only Driver";

        public IPlcConnection Connect(string connectionString)
        {
            throw new NotSupportedException();
        }

        public IPlcConnection Connect(string connectionString, IPlcAuthentication authentication)
        {
            throw new NotSupportedException();
        }
    }

    /// <summary>
    /// The manager is a process-wide singleton, so each test registers under its own
    /// protocol code rather than trying to reset shared state.
    /// </summary>
    public class PlcDriverManagerTests
    {
        [Fact]
        public void A_driver_is_found_under_the_code_it_registered_with()
        {
            var driver = new RegistryOnlyDriver("registry-plain");
            PlcDriverManager.Instance.RegisterDriver(driver);

            Assert.Same(driver, PlcDriverManager.Instance.GetDriverByCode("registry-plain"));
        }

        [Fact]
        public void A_mixed_case_driver_code_is_found_by_a_lowercase_lookup()
        {
            // Connection strings are lowercase by grammar, so a driver declaring
            // "Registry-Mixed" still has to be reachable via "registry-mixed".
            var driver = new RegistryOnlyDriver("Registry-Mixed");
            PlcDriverManager.Instance.RegisterDriver(driver);

            Assert.Same(driver, PlcDriverManager.Instance.GetDriverByCode("registry-mixed"));
            Assert.Same(driver, PlcDriverManager.Instance.GetDriverByCode("REGISTRY-MIXED"));
        }

        [Fact]
        public void An_unknown_code_reports_which_protocol_was_asked_for()
        {
            var ex = Assert.Throws<PlcConnectionException>(
                () => PlcDriverManager.Instance.GetDriverByCode("registry-absent"));

            Assert.Contains("registry-absent", ex.Message);
        }

        [Fact]
        public void A_null_code_is_rejected_like_an_unknown_one()
        {
            Assert.Throws<PlcConnectionException>(
                () => PlcDriverManager.Instance.GetDriverByCode(null));
        }

        [Fact]
        public void Registering_null_is_rejected()
        {
            Assert.Throws<ArgumentNullException>(
                () => PlcDriverManager.Instance.RegisterDriver(null));
        }

        [Fact]
        public void GetDriver_extracts_the_protocol_code_from_a_connection_string()
        {
            var driver = new RegistryOnlyDriver("s7");
            PlcDriverManager.Instance.RegisterDriver(driver);

            // The api layer only needs the scheme to dispatch to the right driver.
            // This mirrors Java's DefaultPlcDriverManager.getDriverForUrl(), which
            // uses new URI(url).getScheme().
            Assert.Same(driver, PlcDriverManager.Instance.GetDriver("s7://192.168.0.1"));
            Assert.Same(driver, PlcDriverManager.Instance.GetDriver("s7:cotp://10.0.0.5:102"));
        }

        [Fact]
        public void GetDriver_extracts_the_protocol_code_from_a_transport_shorthand()
        {
            // The code must not collide with a real driver's: the manager is a global
            // singleton, so registering a stub under "modbus-tcp" would make whichever
            // test runs second see the other's driver.
            var driver = new RegistryOnlyDriver("registry-shorthand");
            PlcDriverManager.Instance.RegisterDriver(driver);

            Assert.Same(driver, PlcDriverManager.Instance.GetDriver("registry-shorthand://10.0.0.9:502?unit-identifier=1"));
        }
    }
}
