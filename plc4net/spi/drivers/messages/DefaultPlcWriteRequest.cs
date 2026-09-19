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
using System.Linq;
using System.Threading.Tasks;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages.items;

namespace org.apache.plc4net.spi.drivers.messages
{
    /// <summary>
    /// Default implementation of a write request.
    /// Mirrors the Java SPI3 <c>DefaultPlcWriteRequest</c>.
    /// </summary>
    public class DefaultPlcWriteRequest : IPlcWriteRequest, IPlcTagRequest
    {
        private readonly PlcWriter _writer;
        private readonly Dictionary<string, PlcTagValueItem<IPlcTag>> _tags;

        public DefaultPlcWriteRequest(PlcWriter writer, Dictionary<string, PlcTagValueItem<IPlcTag>> tags)
        {
            _writer = writer;
            _tags = tags ?? new Dictionary<string, PlcTagValueItem<IPlcTag>>();
        }

        public int NumberOfValues => _tags.Count;

        public int TagCount => _tags.Count;

        public IEnumerable<string> TagNames => _tags.Keys;

        public IPlcTag GetTagByName(string name)
        {
            return _tags.TryGetValue(name, out var item) ? item.Tag : null;
        }

        public IEnumerable<IPlcTag> Tags => _tags.Values.Select(v => v.Tag);

        /// <summary>Returns the value to write to the named tag, or null.</summary>
        public object GetValue(string name)
        {
            return _tags.TryGetValue(name, out var item) ? item.Value : null;
        }

        async Task<IPlcResponse> IPlcRequest.ExecuteAsync()
        {
            return (IPlcResponse)await _writer.Write(this).ConfigureAwait(false);
        }
    }
}
