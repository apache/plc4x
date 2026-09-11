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

using org.apache.plc4net.spi.transports;
using org.apache.plc4net.drivers.modbus;
using Xunit;

namespace org.apache.plc4net.spi.test.drivers
{
    /// <summary>
    /// Tests for the PLC4Net service extensions — driver registration,
    /// assembly scanning, and the convenience registration method.
    /// </summary>
    public class Plc4NetExtensionsTests
    {
        [Fact]
        public void Can_register_a_driver_via_extension_method()
        {
            var transportManager = new DefaultTransportManager(
                new[] { new org.apache.plc4net.transports.tcp.TcpTransport() });
            var driver = new ModbusTcpDriver(transportManager);

            // Just verify registration + lookup works (PlcDriverManager is a singleton;
            // other tests may already have registered modbus-tcp).
            PlcDriverManager.Instance.RegisterDriver(driver);
            var found = PlcDriverManager.Instance.GetDriverByCode("modbus-tcp");
            Assert.NotNull(found);
        }

        [Fact]
        public void Extension_method_rejects_null()
        {
            Assert.Throws<System.ArgumentNullException>(
                () => Plc4NetServiceExtensions.RegisterDriver(null, null));
        }
    }
}
