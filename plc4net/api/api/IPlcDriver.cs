/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

using System.Threading.Tasks;
using org.apache.plc4net.api.authentication;

namespace org.apache.plc4net.api
{
    /// <summary>
    /// Interface for PLC drivers to be implemented
    /// </summary>
    public interface IPlcDriver
    {
        /// <summary>
        /// Get the code of the implemented protocol
        /// </summary>
        string ProtocolCode { get; }

        /// <summary>
        /// Full Name of the implemented protocol
        /// </summary>
        string ProtocolName { get; }

        /// <summary>
        /// Connects to the PLC identified by the connection string
        /// </summary>
        /// <param name="connectionString">Connection string identifying the PLC to connect to</param>
        /// <returns>Awaitable task returning the <see cref="IPlcConnection"/> to which the connection was established</returns>
        /// <exception cref="org.apache.plc4net.exceptions.PlcConnectionException">Thrown on connection failure</exception>
        Task<IPlcConnection> ConnectAsync(string connectionString);

        /// <summary>
        /// Connects to the PLC identified by the connection string and using the 
        /// </summary>
        /// <param name="connectionString"></param>
        /// <param name="authentication"></param>
        /// <returns>Awaitable task returning the <see cref="IPlcConnection"/> to which the connection was established</returns>
        /// <exception cref="org.apache.plc4net.exceptions.PlcConnectionException">Thrown on connection failure</exception>
        Task<IPlcConnection> ConnectAsync(string connectionString, IPlcAuthentication authentication);
    }
}