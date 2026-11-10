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
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

using System;
using System.Collections.Generic;
using System.Linq;
using org.apache.plc4net.api.value;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.drivers.messages.items;

namespace org.apache.plc4net.drivers.slmp
{
    /// <summary>Builds SLMP write requests, preserving multiple values for range tags.</summary>
    public sealed class SlmpWriteRequestBuilder : IPlcWriteRequestBuilder
    {
        private readonly PlcWriter _writer;
        private readonly Func<string, IPlcTag> _tagParser;
        private readonly Dictionary<string, PlcTagValueItem<IPlcTag>> _tags = new();

        public SlmpWriteRequestBuilder(PlcWriter writer, Func<string, IPlcTag> tagParser)
        {
            _writer = writer;
            _tagParser = tagParser;
        }

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params bool[] values)
            => AddValues(name, tagAddress, values);
        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params byte[] values)
            => AddValues(name, tagAddress, values);
        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params short[] values)
            => AddValues(name, tagAddress, values);
        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params int[] values)
            => AddValues(name, tagAddress, values);
        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params long[] values)
            => AddValues(name, tagAddress, values);
        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params float[] values)
            => AddValues(name, tagAddress, values);
        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params double[] values)
            => AddValues(name, tagAddress, values);
        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params Decimal[] values)
            => AddValues(name, tagAddress, values);
        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params string[] values)
            => AddValues(name, tagAddress, values);
        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params DateTime[] values)
            => AddValues(name, tagAddress, values);

        public IPlcWriteRequestBuilder AddTag<T>(string name, string tagAddress, params T[] values)
            => AddValues(name, tagAddress, values);

        public IPlcWriteRequest Build() => new DefaultPlcWriteRequest(_writer, _tags);

        private IPlcWriteRequestBuilder AddValues<T>(string name, string tagAddress, T[] values)
        {
            if (_tags.ContainsKey(name))
            {
                throw new ArgumentException($"Duplicate tag definition '{name}'.", nameof(name));
            }
            if (values == null || values.Length == 0)
            {
                throw new ArgumentException($"Tag '{name}' needs at least one value to write.", nameof(values));
            }

            var tag = _tagParser(tagAddress);
            object? value;
            if (values.Length == 1 && values[0] is IPlcValue plcValue && plcValue.IsList())
            {
                value = plcValue.GetList().Cast<object>().ToArray();
            }
            else if ((tag is SlmpTag slmpTag && slmpTag.ExplicitRange) || values.Length > 1)
            {
                value = values.Cast<object>().ToArray();
            }
            else
            {
                value = values[0] as object;
            }
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(name, tag, value);
            return this;
        }
    }
}
