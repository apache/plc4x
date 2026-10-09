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
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

using System;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.api.value;
using org.apache.plc4net.messages;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.spi.test.drivers
{
    /// <summary>
    /// Covers the read/write request builders and their default implementations.
    /// </summary>
    public class MessagesTests
    {
        private class NoOpReader : PlcReader
        {
            public Task<IPlcReadResponse> Read(DefaultPlcReadRequest request,
                CancellationToken cancellationToken = default)
            {
                var values = new System.Collections.Generic.Dictionary<string,
                    spi.drivers.messages.items.PlcResponseItem<IPlcValue>>();
                foreach (var name in request.TagNames)
                {
                    values[name] = new spi.drivers.messages.items.DefaultPlcResponseItem<IPlcValue>(
                        PlcResponseCode.Ok, null);
                }
                return Task.FromResult<IPlcReadResponse>(new DefaultPlcReadResponse(request, values));
            }
        }

        private class NoOpWriter : PlcWriter
        {
            public Task<IPlcWriteResponse> Write(DefaultPlcWriteRequest request,
                CancellationToken cancellationToken = default)
            {
                var codes = new System.Collections.Generic.Dictionary<string, PlcResponseCode>();
                foreach (var name in request.TagNames)
                {
                    codes[name] = PlcResponseCode.Ok;
                }
                return Task.FromResult<IPlcWriteResponse>(new DefaultPlcWriteResponse(request, codes));
            }
        }

        [Fact]
        public void A_read_request_can_be_built_with_tag_addresses()
        {
            var builder = new DefaultPlcReadRequestBuilder(new NoOpReader(), null);

            var request = builder
                .AddTagAddress("a", "coil:1")
                .AddTagAddress("b", "coil:2")
                .Build();

            Assert.NotNull(request);
            var tagRequest = (IPlcTagRequest)request;
            Assert.Equal(2, tagRequest.TagCount);
            Assert.Contains("a", tagRequest.TagNames);
            Assert.Contains("b", tagRequest.TagNames);
        }

        [Fact]
        public async Task A_read_request_executes_and_returns_a_response()
        {
            var request = new DefaultPlcReadRequestBuilder(new NoOpReader(), null)
                .AddTagAddress("x", "holding:10")
                .Build();

            var response = (DefaultPlcReadResponse)await request.ExecuteAsync();

            Assert.NotNull(response);
            Assert.NotNull(response.Request);
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("x"));
        }

        [Fact]
        public void A_write_request_can_be_built_with_tag_values()
        {
            var builder = new DefaultPlcWriteRequestBuilder(new NoOpWriter(), null);

            var request = builder
                .AddTag("tag1", "coil:0", true)
                .AddTag("tag2", "holding:1", (short)42)
                .Build();

            Assert.NotNull(request);
            Assert.Equal(2, request.NumberOfValues);
            var tagRequest = (IPlcTagRequest)request;
            Assert.Equal(2, tagRequest.TagCount);
        }

        [Fact]
        public async Task A_write_request_executes_and_returns_a_response()
        {
            var request = new DefaultPlcWriteRequestBuilder(new NoOpWriter(), null)
                .AddTag("tag1", "coil:0", true)
                .Build();

            var response = (DefaultPlcWriteResponse)await request.ExecuteAsync();

            Assert.NotNull(response);
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("tag1"));
        }

        [Fact]
        public void Reusing_a_read_builder_does_not_change_a_request_already_built()
        {
            var builder = new DefaultPlcReadRequestBuilder(new NoOpReader(), null);
            var first = (IPlcTagRequest)builder.AddTagAddress("a", "coil:1").Build();

            builder.AddTagAddress("b", "coil:2");

            Assert.Equal(1, first.TagCount);
            Assert.DoesNotContain("b", first.TagNames);
        }

        [Fact]
        public void Reusing_a_write_builder_does_not_change_a_request_already_built()
        {
            var builder = new DefaultPlcWriteRequestBuilder(new NoOpWriter(), null);
            var first = builder.AddTag("a", "coil:0", true).Build();

            builder.AddTag("b", "coil:1", false);

            Assert.Equal(1, first.NumberOfValues);
        }

        [Fact]
        public void A_read_builder_rejects_a_duplicate_tag_name()
        {
            var builder = new DefaultPlcReadRequestBuilder(new NoOpReader(), null)
                .AddTagAddress("a", "coil:1");

            Assert.Throws<ArgumentException>(() => builder.AddTagAddress("a", "coil:2"));
        }

        [Fact]
        public void A_write_builder_rejects_a_duplicate_tag_name()
        {
            var builder = new DefaultPlcWriteRequestBuilder(new NoOpWriter(), null)
                .AddTag("t", "coil:0", true);

            Assert.Throws<ArgumentException>(() => builder.AddTag("t", "coil:1", false));
        }

        [Fact]
        public void A_write_needs_at_least_one_value()
        {
            var builder = new DefaultPlcWriteRequestBuilder(new NoOpWriter(), null);

            Assert.Throws<ArgumentException>(() => builder.AddTag("empty", "coil:0", Array.Empty<bool>()));
        }
    }
}