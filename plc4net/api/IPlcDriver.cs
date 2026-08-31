//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

using org.apache.plc4net.api.authentication;
using org.apache.plc4net.exceptions;

namespace org.apache.plc4net.api
{
    /// <summary>
    /// A PLC protocol driver.
    /// </summary>
    /// <remarks>
    /// This is the SPI3-aligned contract: drivers are synchronous factories that
    /// open connections. The transport selection and connection-string parsing
    /// is handled by <see cref="org.apache.plc4net.spi.drivers.DriverBase"/>.
    /// </remarks>
    public interface IPlcDriver
    {
        string ProtocolCode { get; }

        string ProtocolName { get; }

        /// <summary>
        /// Opens a connection to the device described by <paramref name="connectionString"/>.
        /// </summary>
        /// <exception cref="PlcConnectionException">on connection failure</exception>
        IPlcConnection Connect(string connectionString);

        IPlcConnection Connect(string connectionString, IPlcAuthentication authentication);
    }
}
