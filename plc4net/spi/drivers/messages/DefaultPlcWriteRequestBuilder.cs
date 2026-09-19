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
    /// Default builder for write requests.
    /// </summary>
    public class DefaultPlcWriteRequestBuilder : IPlcWriteRequestBuilder
    {
        private readonly PlcWriter _writer;
        private readonly Func<string, IPlcTag> _tagParser;
        private readonly Dictionary<string, PlcTagValueItem<IPlcTag>> _tags
            = new Dictionary<string, PlcTagValueItem<IPlcTag>>();

        public DefaultPlcWriteRequestBuilder(PlcWriter writer, Func<string, IPlcTag> tagParser)
        {
            _writer = writer;
            _tagParser = tagParser ?? (s => new GenericTag(s));
        }

        // For each typed overload, unwrap a single-element params array so the
        // driver's Write method receives a scalar value rather than an array.

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params bool[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params byte[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params short[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params int[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params long[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params float[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params double[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params Decimal[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params string[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequestBuilder AddTag(string name, string tagAddress, params DateTime[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequestBuilder AddTag<T>(string name, string tagAddress, params T[] values)
        {
            _tags[name] = new DefaultPlcTagValueItem<IPlcTag>(
                name, _tagParser(tagAddress), Unwrap(values));
            return this;
        }

        public IPlcWriteRequest Build()
        {
            return new DefaultPlcWriteRequest(_writer, _tags);
        }

        /// <summary>
        /// If the array has one element, return it directly; if more than one,
        /// throw because multi-value writes are not yet supported.
        /// </summary>
        private static object Unwrap<T>(T[] values)
        {
            if (values == null || values.Length == 0) return null;
            if (values.Length == 1) return values[0];
            throw new NotSupportedException(
                $"Multi-value writes (passed {values.Length} values of type {typeof(T).Name}) " +
                "are not yet supported. Write each tag individually.");
        }
    }
}
