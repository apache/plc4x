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

            Validate(config);

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
            var raw = GetValue(parameters, key);
            if (raw == null) return fallback;
            return int.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture,
                       out var value) ? value : fallback;
        }

        private static T GetEnum<T>(IReadOnlyDictionary<string, string> parameters,
            string key, T fallback) where T : struct
        {
            var raw = GetValue(parameters, key);
            if (raw == null) return fallback;
            return Enum.TryParse<T>(raw, ignoreCase: true, out var value)
                   && Enum.IsDefined(typeof(T), value)
                ? value
                : fallback;
        }

        // ConnectionString keeps all query parameters. Follow the transport convention
        // used by TCP: serial.baud-rate takes precedence, while baud-rate remains
        // convenient for a driver that only supports serial links.
        private static string? GetValue(IReadOnlyDictionary<string, string> parameters, string key)
        {
            if (parameters.TryGetValue("serial." + key, out var prefixed))
            {
                return prefixed;
            }

            return parameters.TryGetValue(key, out var plain) ? plain : null;
        }

        private static void Validate(SerialTransportConfiguration config)
        {
            if (config.BaudRate <= 0)
            {
                throw new TransportException($"baud-rate must be positive, but was {config.BaudRate}.");
            }

            if (config.DataBits < 5 || config.DataBits > 8)
            {
                throw new TransportException($"data-bits must be between 5 and 8, but was {config.DataBits}.");
            }

            if (config.ReadTimeout < -1)
            {
                throw new TransportException("read-timeout must be -1 (infinite) or non-negative.");
            }

            // SerialPort.WriteTimeout is stricter than ReadTimeout: it takes a positive value or
            // -1 and rejects 0. Fail here rather than as an unrelated ArgumentOutOfRangeException
            // when the port is opened.
            if (config.WriteTimeout == 0 || config.WriteTimeout < -1)
            {
                throw new TransportException("write-timeout must be -1 (infinite) or positive.");
            }

            if (config.ReceiveBufferSize <= 0 || config.SendBufferSize < 0)
            {
                throw new TransportException("receive-buffer-size must be positive and send-buffer-size must be non-negative.");
            }

            if (config.StopBits == StopBits.None)
            {
                throw new TransportException("stop-bits must not be None.");
            }
        }
    }
}
