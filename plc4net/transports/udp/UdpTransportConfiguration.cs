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

namespace org.apache.plc4net.transports.udp
{
    /// <summary>
    /// Settings for the UDP transport: default port, local bind address / port, and
    /// receive buffer capacity.
    /// </summary>
    public class UdpTransportConfiguration : ITransportConfiguration
    {
        /// <summary>Port used when the connection string does not name one.</summary>
        public int DefaultPort { get; set; } = -1;

        /// <summary>
        /// Local address to bind the socket to. Empty binds to the wildcard address for
        /// the remote address family and lets the OS choose the outbound interface.
        /// </summary>
        public string LocalAddress { get; set; } = string.Empty;

        /// <summary>Local port to bind to; 0 lets the OS choose an ephemeral port.</summary>
        public int LocalPort { get; set; }

        /// <summary>
        /// Capacity of the user-space ring buffer between the receive loop and the
        /// codec, not the kernel socket buffer. It is clamped to at least 65535 (one
        /// maximum-size datagram) by <see cref="UdpTransport.CreateConfiguration"/>.
        /// </summary>
        public int ReceiveBufferSize { get; set; } = 81920;

    }
}
