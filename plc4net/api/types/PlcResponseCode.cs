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

namespace org.apache.plc4net.types
{
    /// <summary>
    /// Possible responses
    /// </summary>
    public enum PlcResponseCode
    {
        Ok,
        NotFound,
        AccessDenied,
        InvalidAddress,
        InvalidDatatype,
        InternalError,
        ResponsePending,

        /// <summary>The device did not answer within the request timeout.</summary>
        RequestTimeout,

        /// <summary>The device rejected the request as unsupported (e.g. an
        /// illegal function code).</summary>
        Unsupported,

        // Added after the others so that the values above keep their numbers.

        /// <summary>A value carried by the request is malformed.</summary>
        InvalidData,

        /// <summary>The device is currently unable to process the request due to overload.</summary>
        RemoteBusy,

        /// <summary>Something went wrong on the remote side.</summary>
        RemoteError,

        /// <summary>The request and the connection are healthy, but the data the tag points at
        /// is not available yet - typically an asynchronous source (cyclic IO, a subscription,
        /// a cache) has not delivered its first value. Retry shortly; this is not permanent.
        /// Unlike <see cref="ResponsePending"/>, there is no request in flight.</summary>
        NotReady,

        /// <summary>The request is well-formed, but a value it carries is outside the range the
        /// target accepts (e.g. a number too large for the datatype, or an address beyond the
        /// device's extent). Unlike <see cref="InvalidData"/>, the value is not malformed.</summary>
        OutOfRange
    }
}