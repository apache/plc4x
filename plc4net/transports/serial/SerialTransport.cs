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
using System.IO.Ports;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.transports.serial
{
    /// <summary>
    /// The "serial" transport for RS-232 / RS-485 links.
    ///
    /// Connection string examples:
    /// <code>
    ///   modbus:serial://COM1?baud-rate=19200&amp;data-bits=8&amp;parity=Even&amp;stop-bits=One
    ///   modbus:serial:///dev/ttyUSB0?baud-rate=115200
    /// </code>
    ///
    /// Parameters (all optional; defaults are the Modbus RTU convention):
    ///   baud-rate, data-bits, stop-bits, parity, handshake,
    ///   read-timeout, write-timeout, receive-buffer-size, send-buffer-size.
    /// </summary>
    public class SerialTransport : ITransport
    {
        public string TransportCode => "serial";

        public string TransportName => "Serial Port Transport (RS-232 / RS-485)";

        public ITransportConfiguration CreateConfiguration(
            IReadOnlyDictionary<string, string> parameters)
        {
            var config = new SerialTransportConfiguration();
            if (parameters == null) return config;

            config.BaudRate = GetInt(parameters, "baud-rate", config.BaudRate);
            config.DataBits = GetInt(parameters, "data-bits", config.DataBits);
            config.ReadTimeout = GetInt(parameters, "read-timeout", config.ReadTimeout);
            config.WriteTimeout = GetInt(parameters, "write-timeout", config.WriteTimeout);
            config.ReceiveBufferSize = GetInt(parameters, "receive-buffer-size", config.ReceiveBufferSize);
            config.SendBufferSize = GetInt(parameters, "send-buffer-size", config.SendBufferSize);

            config.Parity = GetEnum(parameters, "parity", config.Parity);
            config.StopBits = GetEnum(parameters, "stop-bits", config.StopBits);
            config.Handshake = GetEnum(parameters, "handshake", config.Handshake);

            return config;
        }

        public ITransportInstance CreateTransportInstance(
            string transportConfig,
            ITransportConfiguration configuration)
        {
            if (!(configuration is SerialTransportConfiguration serialConfig))
            {
                throw new ArgumentException(
                    $"Serial transport requires a {nameof(SerialTransportConfiguration)} " +
                    $"but got {configuration?.GetType().Name ?? "null"}.");
            }

            // The transport config is the port name (COM1, /dev/ttyUSB0, …).
            return new SerialTransportInstance(transportConfig, serialConfig);
        }

        private static int GetInt(IReadOnlyDictionary<string, string> parameters,
            string key, int fallback)
        {
            if (!parameters.TryGetValue(key, out var raw)) return fallback;
            return int.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture,
                       out var value) ? value : fallback;
        }

        private static T GetEnum<T>(IReadOnlyDictionary<string, string> parameters,
            string key, T fallback) where T : struct
        {
            if (!parameters.TryGetValue(key, out var raw)) return fallback;
            return Enum.TryParse<T>(raw, ignoreCase: true, out var value) ? value : fallback;
        }
    }
}
