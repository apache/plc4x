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

namespace org.apache.plc4net.transports.tcp
{
    /// <summary>
    /// Settings for the TCP transport. Defaults match the Java SPI3 transport so the
    /// same connection string behaves the same from either language.
    /// </summary>
    public class TcpTransportConfiguration : ITransportConfiguration
    {
        /// <summary>Port used when the connection string does not name one.</summary>
        public int DefaultPort { get; set; } = -1;

        /// <summary>Connect timeout in milliseconds.</summary>
        public int ConnectTimeout { get; set; } = 5000;

        /// <summary>Disable Nagle's algorithm. On by default: PLC traffic is small and latency sensitive.</summary>
        public bool TcpNoDelay { get; set; } = true;

        public bool KeepAlive { get; set; } = true;

        /// <summary>Socket send buffer size; 0 leaves the OS default in place.</summary>
        public int SendBufferSize { get; set; }

        /// <summary>
        /// Size of the ring buffer between the read loop and the codec, and of the
        /// socket receive buffer.
        /// </summary>
        public int ReceiveBufferSize { get; set; } = 81920;

        /// <summary>Optional local address to bind to before connecting.</summary>
        public string LocalAddress { get; set; }

        public int LocalPort { get; set; }
    }
}
