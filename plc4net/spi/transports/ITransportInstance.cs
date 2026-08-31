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

namespace org.apache.plc4net.spi.transports
{
    /// <summary>
    /// Marker for a transport's configuration object. Concrete transports define their
    /// own (host, port, timeouts, buffer sizes, ...).
    /// </summary>
    public interface ITransportConfiguration
    {
    }

    /// <summary>
    /// A live connection to a device, exposed as a byte stream.
    /// </summary>
    /// <remarks>
    /// Mirrors the Java SPI3 <c>TransportInstance</c> contract. Deliberately byte-oriented
    /// and peekable: most industrial protocols put the message length inside a fixed-size
    /// header, so a codec peeks the header, works out the full length, and only then
    /// consumes the whole frame.
    /// </remarks>
    public interface ITransportInstance : IDisposable
    {
        /// <summary>The transport's configuration.</summary>
        ITransportConfiguration Configuration { get; }

        /// <summary>
        /// The part of the connection URL the transport did not consume, which is
        /// meaningful only to the driver. For <c>opcua:tcp://localhost:4840/milo</c> the
        /// transport takes <c>localhost:4840</c> and this returns <c>/milo</c>. Empty
        /// when the address carries no trailing part.
        /// </summary>
        string DriverConfig { get; }

        /// <summary>Whether the underlying transport is still usable.</summary>
        bool IsOpen { get; }

        /// <summary>Number of bytes buffered and ready to read.</summary>
        int GetNumBytesAvailable();

        /// <summary>Returns the next bytes without advancing the read position.</summary>
        byte[] PeekReadableBytes(int numBytes);

        /// <summary>Returns the next bytes and marks them consumed.</summary>
        byte[] Read(int numBytes);

        /// <summary>Writes the given bytes to the transport.</summary>
        void Write(byte[] bytes);

        /// <summary>Closes the transport. Safe to call more than once.</summary>
        void Close();
    }

    /// <summary>
    /// A transport that can push a notification when bytes arrive, so the driver does not
    /// have to poll <see cref="ITransportInstance.GetNumBytesAvailable"/>.
    /// </summary>
    public interface IAsyncTransportInstance : ITransportInstance
    {
        /// <summary>
        /// Called after new bytes land in the buffer. Invoked on the transport's read
        /// thread, so implementations should hand off rather than block.
        /// </summary>
        void RegisterDataListener(Action listener);

        void RemoveDataListener();

        /// <summary>
        /// Called when the connection drops. The argument is the causing exception, or
        /// null for an orderly close by the remote end.
        /// </summary>
        void RegisterDisconnectListener(Action<Exception> listener);

        void RemoveDisconnectListener();
    }
}
