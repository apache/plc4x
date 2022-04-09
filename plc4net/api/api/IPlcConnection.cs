/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

using System;
using System.Threading.Tasks;
using org.apache.plc4net.api.metadata;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;

namespace org.apache.plc4net.api
{
    /// <summary>
    /// Interface for generalized PLC connections providing
    /// functionality for basic operations like connect / disconnect etc.
    /// </summary>
    public interface IPlcConnection: IDisposable
    {
        /// <summary>
        /// Connect to the PLC asynchronously
        /// </summary>
        /// <returns>Awaitable task</returns>
        /// <exception cref="org.apache.plc4net.exceptions.PlcConnectionException">Thrown if the connection to the PLC fails</exception>
        Task ConnectAsync();

        /// <summary>
        /// Indicates the connection state
        /// </summary>
        bool IsConnected { get; }

        /// <summary>
        /// Close the PLC connection asynchronously
        /// </summary>
        /// <returns>Awaitable task</returns>
        Task CloseAsync();

        /// <summary>
        /// Get the metadata for the connection
        /// </summary>
        IPlcConnectionMetadata PlcConnectionMetadata { get; }

        /// <summary>
        /// Parse the given field query
        /// </summary>
        /// <param name="fieldQuery">Query for the field</param>
        /// <returns>Field parsed from the query string</returns>
        /// <exception cref="org.apache.plc4net.exceptions.PlcInvalidFieldException">Thrown when the query can not be parsed</exception>
        IPlcField Parse(string fieldQuery);

        /// <summary>
        /// Request builder for constructing read requests
        /// </summary>
        /// <returns>null if the connection does not support reading</returns>
        IPlcReadRequestBuilder ReadRequestBuilder { get; }

        /// <summary>
        /// Request builder for constructing write requests
        /// </summary>
        /// <returns>null if the connection does not support writing</returns>
        IPlcWriteRequestBuilder WriteRequestBuilder { get; }

        /// <summary>
        /// Request builder for constructing subscription requests
        /// </summary>
        /// <returns>null if the connection does not support subscriptions</returns>
        IPlcSubscriptionRequestBuilder SubscriptionRequestBuilder { get; }

        /// <summary>
        /// Request builder for unsubscribing
        /// </summary>
        /// <returns>null if the connection does not support subscriptions</returns>
        IPlcUnsubscriptionRequestBuilder UnsubscriptionRequestBuilder { get; }
    }
}