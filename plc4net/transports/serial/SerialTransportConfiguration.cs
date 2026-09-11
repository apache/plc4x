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

using System.IO.Ports;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.transports.serial
{
    /// <summary>
    /// RS-232 / RS-485 transport settings.  Defaults match the most common
    /// Modbus RTU convention: 19200-8-E-1.
    /// </summary>
    public class SerialTransportConfiguration : ITransportConfiguration
    {
        /// <summary>Baud rate: 9600, 19200, 38400, 57600, 115200 …</summary>
        public int BaudRate { get; set; } = 19200;

        /// <summary>Data bits: 5, 6, 7, or 8.</summary>
        public int DataBits { get; set; } = 8;

        /// <summary>Stop bits: <see cref="System.IO.Ports.StopBits"/>.</summary>
        public StopBits StopBits { get; set; } = StopBits.One;

        /// <summary>Parity: <see cref="System.IO.Ports.Parity"/>.</summary>
        public Parity Parity { get; set; } = Parity.Even;

        /// <summary>Handshake: <see cref="System.IO.Ports.Handshake"/>.</summary>
        public Handshake Handshake { get; set; } = Handshake.None;

        /// <summary>Read timeout in milliseconds.</summary>
        public int ReadTimeout { get; set; } = 2000;

        /// <summary>Write timeout in milliseconds.</summary>
        public int WriteTimeout { get; set; } = 2000;

        /// <summary>Size of the ring buffer between the receive thread and the codec.</summary>
        public int ReceiveBufferSize { get; set; } = 81920;

        /// <summary>Size of the ring buffer for the transmit path (queue → wire).</summary>
        public int SendBufferSize { get; set; }
    }
}
