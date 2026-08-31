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

using System.Collections.Generic;
using org.apache.plc4net.api.value;
using org.apache.plc4net.messages;
using org.apache.plc4net.spi.drivers.messages.items;
using org.apache.plc4net.types;

namespace org.apache.plc4net.spi.drivers.messages
{
    /// <summary>
    /// Default implementation of a read response.
    /// Mirrors the Java SPI3 <c>DefaultPlcReadResponse</c>.
    /// </summary>
    public class DefaultPlcReadResponse : IPlcReadResponse
    {
        private readonly IPlcReadRequest _request;
        private readonly Dictionary<string, PlcResponseItem<IPlcValue>> _values;

        public DefaultPlcReadResponse(
            IPlcReadRequest request,
            Dictionary<string, PlcResponseItem<IPlcValue>> values)
        {
            _request = request;
            _values = values ?? new Dictionary<string, PlcResponseItem<IPlcValue>>();
        }

        // Explicit interface implementation: IPlcResponse.Request returns IPlcRequest.
        IPlcRequest IPlcResponse.Request => _request;

        /// <summary>The originating read request.</summary>
        public IPlcReadRequest Request => _request;

        /// <summary>
        /// Returns the response code for a named tag.
        /// </summary>
        public PlcResponseCode GetResponseCode(string name)
        {
            return _values.TryGetValue(name, out var item)
                ? item.Code
                : PlcResponseCode.NotFound;
        }

        /// <summary>
        /// Returns the value for a named tag, or null if the tag is absent or in error.
        /// </summary>
        public IPlcValue GetValue(string name)
        {
            return _values.TryGetValue(name, out var item) && item.Code == PlcResponseCode.Ok
                ? item.Value
                : null;
        }

        public IEnumerable<string> TagNames => _values.Keys;
    }
}
