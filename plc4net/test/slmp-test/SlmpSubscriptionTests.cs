// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
// https://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied. See the License for the
// specific language governing permissions and limitations
// under the License.

using System;
using System.Threading.Tasks;
using org.apache.plc4net.drivers.slmp;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.messages;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.drivers.subscriptions;
using org.apache.plc4net.transports.test;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.test.slmp
{
    public class SlmpSubscriptionTests
    {
        [Theory]
        [InlineData("D100:WORD", "D00000FFFF0300040000003412", "500000FFFF03000C00100001040000640000A80100")]
        [InlineData("M100:BOOL", "D00000FFFF03000300000010", "500000FFFF03000C00100001040100640000900100")]
        public async Task Cyclic_subscription_uses_existing_batch_read_and_connection_remains_readable(
            string address, string reply, string expected)
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            using var connection = new SlmpConnection(ConnectionString.Parse("slmp://localhost"), transport);
            Assert.True(connection.PlcConnectionMetadata.CanSubscribe);
            var subscribed = (IPlcSubscriptionResponse)await connection.SubscriptionRequestBuilder!
                .AddCyclicField("x", address, TimeSpan.FromMilliseconds(100)).Build().ExecuteAsync();
            var result = new TaskCompletionSource<SubscriptionValueEventArgs>(TaskCreationOptions.RunContinuationsAsynchronously);
            var handle = subscribed.GetSubscriptionHandle("x");
            handle.PlcSubscriptionEvent += (_, e) =>
            {
                connection.UnsubscriptionRequestBuilder!.AddHandles(handle).Build().ExecuteAsync().GetAwaiter().GetResult();
                result.TrySetResult((SubscriptionValueEventArgs)e);
            };
            transport.InjectTestData(Convert.FromHexString(reply));
            var value = await result.Task.WaitAsync(TimeSpan.FromSeconds(5));
            Assert.Equal(PlcResponseCode.Ok, value.ResponseCode);
            Assert.Equal(Convert.FromHexString(expected), transport.GetAllWrittenData());
            if (address.StartsWith("M", StringComparison.Ordinal)) Assert.True(value.Value!.GetBool());
            else Assert.Equal(0x1234, value.Value!.GetInt());

            transport.InjectTestData(Convert.FromHexString(reply));
            var read = (DefaultPlcReadResponse)await connection.ReadRequestBuilder!.AddTagAddress("again", address).Build().ExecuteAsync();
            Assert.Equal(PlcResponseCode.Ok, read.GetResponseCode("again"));
            connection.Close();
            await Assert.ThrowsAsync<ObjectDisposedException>(() => connection.SubscriptionRequestBuilder!
                .AddCyclicField("closed", address, TimeSpan.FromSeconds(1)).Build().ExecuteAsync());
        }

        [Theory]
        [InlineData("0")]
        [InlineData("-1")]
        [InlineData("invalid")]
        public void Invalid_change_interval_is_rejected(string interval)
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            Assert.Throws<PlcConnectionException>(() => new SlmpConnection(
                ConnectionString.Parse("slmp://localhost?subscription-polling-interval-ms=" + interval), transport));
        }
    }
}
