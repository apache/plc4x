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

using System.Collections.Generic;
using System.Globalization;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.transports.test
{
    /// <summary>
    /// In-memory transport for testing drivers without hardware or a network.
    /// Mirrors the Java SPI3 <c>TestTransport</c>.
    /// </summary>
    /// <remarks>
    /// A test wires this transport into a driver and uses
    /// <see cref="TestTransportInstance.InjectTestData"/> to simulate the
    /// bytes that would arrive from the wire, then reads the driver's
    /// response with <see cref="TestTransportInstance.GetAllWrittenData"/>.
    /// </remarks>
    public class TestTransport : ITransport
    {
        public string TransportCode => "test";

        public string TransportName => "Test Transport (in-memory)";

        public ITransportConfiguration CreateConfiguration(IReadOnlyDictionary<string, string> parameters)
        {
            var configuration = new TestTransportConfiguration();
            if (parameters == null)
            {
                return configuration;
            }

            if (parameters.TryGetValue("receive-buffer-size", out var raw) &&
                int.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture, out var size) && size > 0)
            {
                configuration.ReceiveBufferSize = size;
            }

            return configuration;
        }

        public ITransportInstance CreateTransportInstance(
            string transportConfig, ITransportConfiguration configuration)
        {
            if (!(configuration is TestTransportConfiguration testConfig))
            {
                throw new TransportException(
                    $"Expected a {nameof(TestTransportConfiguration)} but got " +
                    $"{configuration?.GetType().Name ?? "null"}.");
            }

            return new TestTransportInstance(testConfig);
        }
    }
}
