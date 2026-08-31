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
    /// Default implementation of a read request.
    /// Mirrors the Java SPI3 <c>DefaultPlcReadRequest</c>.
    /// </summary>
    public class DefaultPlcReadRequest : IPlcReadRequest, IPlcTagRequest
    {
        private readonly PlcReader _reader;
        private readonly Dictionary<string, PlcTagItem<IPlcTag>> _tags;

        public DefaultPlcReadRequest(PlcReader reader, Dictionary<string, PlcTagItem<IPlcTag>> tags)
        {
            _reader = reader;
            _tags = tags ?? new Dictionary<string, PlcTagItem<IPlcTag>>();
        }

        public int TagCount => _tags.Count;

        public IEnumerable<string> TagNames => _tags.Keys;

        public IPlcTag GetTagByName(string name)
        {
            return _tags.TryGetValue(name, out var item) ? item.Tag : null;
        }

        public IEnumerable<IPlcTag> Tags => _tags.Values.Select(v => v.Tag);

        async Task<IPlcResponse> IPlcRequest.ExecuteAsync()
        {
            return (IPlcResponse)await _reader.Read(this).ConfigureAwait(false);
        }
    }
}
