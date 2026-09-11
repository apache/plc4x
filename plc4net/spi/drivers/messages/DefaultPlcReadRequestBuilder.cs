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

using System;
using System.Collections.Generic;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages.items;

namespace org.apache.plc4net.spi.drivers.messages
{
    /// <summary>
    /// Default builder for read requests.
    /// </summary>
    public class DefaultPlcReadRequestBuilder : IPlcReadRequestBuilder
    {
        private readonly PlcReader _reader;
        private readonly Func<string, IPlcTag> _tagParser;
        private readonly Dictionary<string, PlcTagItem<IPlcTag>> _tags
            = new Dictionary<string, PlcTagItem<IPlcTag>>();

        public DefaultPlcReadRequestBuilder(PlcReader reader, Func<string, IPlcTag> tagParser)
        {
            _reader = reader;
            _tagParser = tagParser ?? (s => new GenericTag(s));
        }

        public IPlcReadRequestBuilder AddTagAddress(string name, string tagAddress)
        {
            var tag = _tagParser(tagAddress);
            _tags[name] = new DefaultPlcTagItem<IPlcTag>(name, tag);
            return this;
        }

        public IPlcReadRequest Build()
        {
            return new DefaultPlcReadRequest(_reader, _tags);
        }
    }

    /// <summary>
    /// Fallback tag carrying an address string — used when no protocol-specific
    /// parser is available.
    /// </summary>
    internal class GenericTag : IPlcTag
    {
        public GenericTag(string address) { Address = address; }
        public string Address { get; }
        public override string ToString() => Address;
    }
}
