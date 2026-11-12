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
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Channels;
using System.Threading.Tasks;
using org.apache.plc4net.api.value;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.drivers.messages.items;
using org.apache.plc4net.spi.drivers.subscriptions;
using org.apache.plc4net.spi.model.values;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.test.spi.drivers
{
    public class PollingSubscriptionTests
    {
        private sealed class Tag : IPlcTag { }
        private sealed record ReadCall(DefaultPlcReadRequest Request, TaskCompletionSource<IPlcReadResponse> Reply);
        private sealed class Reader : PlcReader
        {
            public readonly Channel<ReadCall> Calls = Channel.CreateUnbounded<ReadCall>();
            public Task<IPlcReadResponse> Read(DefaultPlcReadRequest request, CancellationToken cancellationToken = default)
            {
                var reply = new TaskCompletionSource<IPlcReadResponse>(TaskCreationOptions.RunContinuationsAsynchronously);
                Calls.Writer.TryWrite(new ReadCall(request, reply));
                return reply.Task;
            }
        }

        private static readonly TimeSpan Interval = TimeSpan.FromMilliseconds(20);
        private static readonly TimeSpan Timeout = TimeSpan.FromSeconds(5);
        private static PollingSubscriptionManager Manager(Reader reader, Action<Exception>? errors = null)
            => new PollingSubscriptionManager(reader, _ => new Tag(), Interval, errors ?? (_ => { }));
        private static async Task<ReadCall> Next(Reader reader)
            => await reader.Calls.Reader.ReadAsync().AsTask().WaitAsync(Timeout);
        private static void Reply(ReadCall call, IPlcValue? value, PlcResponseCode code = PlcResponseCode.Ok)
            => call.Reply.SetResult(new DefaultPlcReadResponse(call.Request,
                new Dictionary<string, PlcResponseItem<IPlcValue>>
                {
                    [call.Request.TagNames.Single()] = new DefaultPlcResponseItem<IPlcValue>(code, value)
                }));
        private static async Task<IPlcSubscriptionResponse> Subscribe(PollingSubscriptionManager manager, bool changes = false)
            => Assert.IsAssignableFrom<IPlcSubscriptionResponse>(await (changes
                ? manager.SubscriptionRequestBuilder.AddChangeOfStateField("x", "D0")
                : manager.SubscriptionRequestBuilder.AddCyclicField("x", "D0", Interval)).Build().ExecuteAsync());

        [Fact]
        public async Task Cyclic_subscription_delivers_values_and_unsubscribe_drains_a_pending_read()
        {
            var reader = new Reader();
            await using var manager = Manager(reader);
            var response = await Subscribe(manager);
            var handle = response.GetSubscriptionHandle("x");
            var events = Channel.CreateUnbounded<SubscriptionValueEventArgs>();
            handle.PlcSubscriptionEvent += (_, e) => events.Writer.TryWrite(Assert.IsType<SubscriptionValueEventArgs>(e));
            Reply(await Next(reader), new PlcINT(123));
            var received = await events.Reader.ReadAsync().AsTask().WaitAsync(Timeout);
            Assert.Equal(123, received.Value!.GetInt());
            Assert.Equal(DateTimeKind.Utc, received.Timestamp.Kind);
            var pending = await Next(reader);
            var stop = manager.UnsubscriptionRequestBuilder.AddHandles(handle).Build().ExecuteAsync();
            Assert.False(stop.IsCompleted);
            Reply(pending, new PlcINT(124));
            await stop.WaitAsync(Timeout);
            Assert.False(events.Reader.TryRead(out _));
            await manager.DisposeAsync();
            Assert.False(reader.Calls.Reader.TryRead(out _));
        }

        [Theory]
        [InlineData(false)]
        [InlineData(true)]
        public async Task Changes_compare_numeric_values_and_array_contents_without_notifying_the_baseline(bool list)
        {
            var reader = new Reader();
            await using var manager = Manager(reader);
            var response = await Subscribe(manager, true);
            var events = Channel.CreateUnbounded<SubscriptionValueEventArgs>();
            response.GetSubscriptionHandle("x").PlcSubscriptionEvent += (_, e)
                => events.Writer.TryWrite((SubscriptionValueEventArgs)e);
            IPlcValue Value(float x) => list ? new PlcList(new List<IPlcValue> { new PlcREAL(x) }) : new PlcREAL(x);
            Reply(await Next(reader), Value(1.25f));
            Reply(await Next(reader), Value(1.25f));
            var third = await Next(reader);
            Assert.False(events.Reader.TryRead(out _));
            Reply(third, Value(1.75f));
            var changed = await events.Reader.ReadAsync().AsTask().WaitAsync(Timeout);
            Assert.Equal(1.75f, list ? changed.Value!.GetIndex(0).GetFloat() : changed.Value!.GetFloat());
            await manager.UnsubscriptionRequestBuilder.AddHandles(response.SubscriptionHandles).Build().ExecuteAsync();
        }

        [Fact]
        public async Task Error_event_preserves_last_successful_baseline_and_recovery_can_notify()
        {
            var reader = new Reader();
            await using var manager = Manager(reader);
            var response = await Subscribe(manager, true);
            var events = Channel.CreateUnbounded<SubscriptionValueEventArgs>();
            response.GetSubscriptionHandle("x").PlcSubscriptionEvent += (_, e) => events.Writer.TryWrite((SubscriptionValueEventArgs)e);
            Reply(await Next(reader), new PlcINT(1));
            Reply(await Next(reader), null, PlcResponseCode.RequestTimeout);
            var error = await events.Reader.ReadAsync().AsTask().WaitAsync(Timeout);
            Assert.Equal(PlcResponseCode.RequestTimeout, error.ResponseCode);
            Assert.Null(error.Value);
            Reply(await Next(reader), new PlcINT(2));
            Assert.Equal(2, (await events.Reader.ReadAsync().AsTask().WaitAsync(Timeout)).Value!.GetInt());
            await manager.UnsubscriptionRequestBuilder.AddHandles(response.SubscriptionHandles).Build().ExecuteAsync();
        }

        [Fact]
        public async Task Listener_exception_is_isolated_and_callback_can_unsubscribe_synchronously()
        {
            var reader = new Reader();
            var errors = new ConcurrentQueue<Exception>();
            await using var manager = Manager(reader, errors.Enqueue);
            var response = await Subscribe(manager);
            var handle = response.GetSubscriptionHandle("x");
            var done = new TaskCompletionSource(TaskCreationOptions.RunContinuationsAsynchronously);
            handle.PlcSubscriptionEvent += (_, _) => throw new InvalidOperationException("observer failed");
            handle.PlcSubscriptionEvent += (_, _) =>
            {
                manager.UnsubscriptionRequestBuilder.AddHandles(handle).Build().ExecuteAsync().GetAwaiter().GetResult();
                done.SetResult();
            };
            Reply(await Next(reader), new PlcBOOL(true));
            await done.Task.WaitAsync(Timeout);
            Assert.Single(errors);
        }

        [Fact]
        public async Task Closing_from_callback_does_not_deadlock_and_prevents_new_subscriptions()
        {
            var reader = new Reader();
            await using var manager = Manager(reader);
            var response = await Subscribe(manager);
            var done = new TaskCompletionSource(TaskCreationOptions.RunContinuationsAsynchronously);
            response.GetSubscriptionHandle("x").PlcSubscriptionEvent += (_, _) =>
            {
                manager.DisposeAsync().AsTask().GetAwaiter().GetResult();
                done.SetResult();
            };
            Reply(await Next(reader), new PlcWORD(7));
            await done.Task.WaitAsync(Timeout);
            await Assert.ThrowsAsync<ObjectDisposedException>(() => Subscribe(manager));
        }

        [Fact]
        public async Task Error_observer_can_close_the_manager_without_waiting_for_itself()
        {
            var reader = new Reader();
            PollingSubscriptionManager? owner = null;
            var done = new TaskCompletionSource(TaskCreationOptions.RunContinuationsAsynchronously);
            await using var manager = Manager(reader, _ =>
            {
                owner!.DisposeAsync().AsTask().GetAwaiter().GetResult();
                done.SetResult();
            });
            owner = manager;
            await Subscribe(manager);
            (await Next(reader)).Reply.SetException(new InvalidOperationException("read failed"));
            await done.Task.WaitAsync(Timeout);
        }

        [Fact]
        public async Task Multiple_subscriptions_share_a_serial_reader_and_close_drains_it()
        {
            var reader = new Reader();
            await using var manager = Manager(reader);
            await Subscribe(manager);
            await Subscribe(manager);
            var pending = await Next(reader);
            await Task.Delay(Interval * 3);
            Assert.False(reader.Calls.Reader.TryRead(out _));
            var stop = manager.DisposeAsync().AsTask();
            Assert.False(stop.IsCompleted);
            Reply(pending, new PlcINT(9));
            await stop.WaitAsync(Timeout);
            Assert.False(reader.Calls.Reader.TryRead(out _));
        }

        [Fact]
        public async Task Builders_snapshot_fields_reject_invalid_intervals_events_and_foreign_handles()
        {
            var reader = new Reader();
            await using var manager = Manager(reader);
            await using var other = Manager(reader);
            var builder = manager.SubscriptionRequestBuilder;
            Assert.Throws<InvalidOperationException>(() => builder.Build());
            Assert.Throws<ArgumentOutOfRangeException>(() => builder.AddCyclicField("x", "D0", TimeSpan.Zero));
            Assert.Throws<NotSupportedException>(() => builder.AddEventField("x", "D0"));
            var request = builder.AddCyclicField("x", "D0", TimeSpan.FromHours(1)).Build();
            builder.AddCyclicField("y", "D1", Interval);
            Assert.Single(request.TagNames);
            var response = (IPlcSubscriptionResponse)await request.ExecuteAsync();
            Assert.Throws<ArgumentException>(() => other.UnsubscriptionRequestBuilder.AddHandles(response.SubscriptionHandles));
        }
    }
}
