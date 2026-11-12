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
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.api.value;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.drivers.messages.items;
using org.apache.plc4net.types;

namespace org.apache.plc4net.spi.drivers.subscriptions
{
    /// <summary>Emulates subscriptions with serialized reads on one connection.</summary>
    public sealed class PollingSubscriptionManager : IAsyncDisposable
    {
        private readonly object _gate = new object();
        private readonly List<Handle> _handles = new List<Handle>();
        private readonly PlcReader _reader;
        private readonly Func<string, IPlcTag> _parse;
        private readonly Action<Exception> _reportError;
        private readonly TimeSpan _changeInterval;
        private readonly CancellationTokenSource _stop = new CancellationTokenSource();
        private readonly SemaphoreSlim _wake = new SemaphoreSlim(0, 1);
        private readonly AsyncLocal<bool> _dispatching = new AsyncLocal<bool>();
        private Task? _worker;
        private bool _closed;

        public PollingSubscriptionManager(PlcReader reader, Func<string, IPlcTag> parse,
            TimeSpan changeInterval, Action<Exception> reportError)
        {
            ValidateInterval(changeInterval);
            _reader = reader ?? throw new ArgumentNullException(nameof(reader));
            _parse = parse ?? throw new ArgumentNullException(nameof(parse));
            _reportError = reportError ?? throw new ArgumentNullException(nameof(reportError));
            _changeInterval = changeInterval;
        }

        public IPlcSubscriptionRequestBuilder SubscriptionRequestBuilder => new Builder(this);
        public IPlcUnsubscriptionRequestBuilder UnsubscriptionRequestBuilder => new UnsubscribeBuilder(this);

        private static void ValidateInterval(TimeSpan interval)
        {
            if (interval < TimeSpan.FromMilliseconds(1) || interval.TotalMilliseconds > int.MaxValue)
                throw new ArgumentOutOfRangeException(nameof(interval));
        }

        private SubscriptionResponse Subscribe(SubscriptionRequest request)
        {
            lock (_gate)
            {
                ObjectDisposedException.ThrowIf(_closed, this);
                var handles = request.Fields.Select(f => new Handle(this, f)).ToArray();
                _handles.AddRange(handles);
                if (_wake.CurrentCount == 0) _wake.Release();
                _worker ??= Task.Run(RunAsync);
                return new SubscriptionResponse(request, handles);
            }
        }

        private async Task RunAsync()
        {
            try
            {
                while (true)
                {
                    Handle[] due;
                    int delay;
                    lock (_gate)
                    {
                        due = _handles.Where(h => h.Next <= Environment.TickCount64).ToArray();
                        delay = _handles.Count == 0 ? Timeout.Infinite :
                            (int)Math.Clamp(_handles.Min(h => h.Next) - Environment.TickCount64, 1, int.MaxValue);
                    }
                    if (due.Length == 0)
                    {
                        await _wake.WaitAsync(delay, _stop.Token).ConfigureAwait(false);
                        continue;
                    }
                    _stop.Token.ThrowIfCancellationRequested();
                    foreach (var handle in due)
                    {
                        TaskCompletionSource drained;
                        lock (_gate)
                        {
                            if (_closed || !handle.Active) continue;
                            drained = new TaskCompletionSource(TaskCreationOptions.RunContinuationsAsynchronously);
                            handle.InFlight = drained.Task;
                        }
                        SubscriptionValueEventArgs result;
                        try
                        {
                            var tags = new Dictionary<string, PlcTagItem<IPlcTag>>
                            {
                                [handle.Field.Name] = new DefaultPlcTagItem<IPlcTag>(handle.Field.Name, handle.Field.Tag)
                            };
                            var request = new DefaultPlcReadRequest(_reader, tags);
                            // Drain an in-flight response even when a subscription is cancelled.
                            // Aborting a 3E exchange would make the connection unusable.
                            var response = await _reader.Read(request).ConfigureAwait(false);
                            if (response is not DefaultPlcReadResponse values)
                                throw new InvalidOperationException("Polling requires a value-bearing read response.");
                            result = new SubscriptionValueEventArgs(handle.Field.Name,
                                values.GetResponseCode(handle.Field.Name), values.GetValue(handle.Field.Name));
                        }
                        catch (Exception ex)
                        {
                            ReportError(ex);
                            result = new SubscriptionValueEventArgs(handle.Field.Name,
                                ex is TimeoutException ? PlcResponseCode.RequestTimeout : PlcResponseCode.InternalError, null);
                        }
                        finally
                        {
                            drained.TrySetResult();
                        }
                        lock (_gate)
                            handle.Next = Environment.TickCount64 + (long)Math.Ceiling(handle.Field.Interval.TotalMilliseconds);
                        handle.Publish(result);
                    }
                }
            }
            catch (OperationCanceledException) when (_stop.IsCancellationRequested) { }
            finally { _stop.Dispose(); _wake.Dispose(); }
        }

        private void ReportError(Exception error)
        {
            var dispatching = _dispatching.Value;
            _dispatching.Value = true;
            try { _reportError(error); }
            catch { /* An observer must not terminate the polling scheduler. */ }
            finally { _dispatching.Value = dispatching; }
        }

        private async Task UnsubscribeAsync(Handle[] handles)
        {
            Task[] pending;
            lock (_gate)
            {
                foreach (var handle in handles)
                    if (handle.Owner != this) throw new ArgumentException("Subscription belongs to another connection.");
                foreach (var handle in handles)
                {
                    handle.Active = false;
                    _handles.Remove(handle);
                }
                pending = handles.Select(h => h.InFlight).ToArray();
            }
            foreach (var handle in handles) handle.ClearListeners();
            await Task.WhenAll(pending).ConfigureAwait(false);
        }

        public ValueTask DisposeAsync()
        {
            Handle[] handles;
            Task worker;
            lock (_gate)
            {
                if (!_closed)
                {
                    _closed = true;
                    _stop.Cancel();
                    if (_worker == null) { _stop.Dispose(); _wake.Dispose(); }
                }
                handles = _handles.ToArray();
                foreach (var handle in handles) handle.Active = false;
                _handles.Clear();
                worker = _worker ?? Task.CompletedTask;
            }
            foreach (var handle in handles) handle.ClearListeners();
            // Closing from a callback must not wait for that same callback.
            return new ValueTask(_dispatching.Value ? Task.CompletedTask : worker);
        }

        private sealed record Field(string Name, IPlcTag Tag, TimeSpan Interval, bool ChangesOnly);

        private sealed class Builder : IPlcSubscriptionRequestBuilder
        {
            private readonly PollingSubscriptionManager _owner;
            private readonly Dictionary<string, Field> _fields = new Dictionary<string, Field>();
            public Builder(PollingSubscriptionManager owner) => _owner = owner;
            public IPlcSubscriptionRequestBuilder AddCyclicField(string name, string fieldQuery, TimeSpan pollingInterval)
                => Add(name, fieldQuery, pollingInterval, false);
            public IPlcSubscriptionRequestBuilder AddChangeOfStateField(string name, string fieldQuery)
                => Add(name, fieldQuery, _owner._changeInterval, true);
            public IPlcSubscriptionRequestBuilder AddEventField(string name, string fieldQuery)
                => throw new NotSupportedException("Polling cannot emulate native event subscriptions.");
            private Builder Add(string name, string query, TimeSpan interval, bool changesOnly)
            {
                ArgumentException.ThrowIfNullOrWhiteSpace(name);
                ValidateInterval(interval);
                _fields.Add(name, new Field(name, _owner._parse(query), interval, changesOnly));
                return this;
            }
            public IPlcSubscriptionRequest Build()
            {
                if (_fields.Count == 0) throw new InvalidOperationException("A subscription requires at least one field.");
                return new SubscriptionRequest(_owner, _fields.Values.ToArray());
            }
        }

        private sealed class SubscriptionRequest : IPlcSubscriptionRequest
        {
            private readonly PollingSubscriptionManager _owner;
            public Field[] Fields { get; }
            public SubscriptionRequest(PollingSubscriptionManager owner, Field[] fields) { _owner = owner; Fields = fields; }
            public int TagCount => Fields.Length;
            public IEnumerable<string> TagNames => Fields.Select(f => f.Name);
            public IEnumerable<IPlcTag> Tags => Fields.Select(f => f.Tag);
            public IPlcTag? GetTagByName(string name) => Fields.FirstOrDefault(f => f.Name == name)?.Tag;
            public Task<IPlcResponse> ExecuteAsync() => Task.FromResult<IPlcResponse>(_owner.Subscribe(this));
        }

        private sealed class SubscriptionResponse : IPlcSubscriptionResponse
        {
            private readonly SubscriptionRequest _request;
            private readonly Handle[] _handles;
            public SubscriptionResponse(SubscriptionRequest request, Handle[] handles) { _request = request; _handles = handles; }
            public IPlcRequest Request => _request;
            public IEnumerable<string> TagNames => _request.TagNames;
            public IPlcTag? GetTagByName(string name) => _request.GetTagByName(name);
            public PlcResponseCode ResponseCode => PlcResponseCode.Ok;
            public IEnumerable<IPlcSubscriptionHandle> SubscriptionHandles => _handles.ToArray();
            public IPlcSubscriptionHandle GetSubscriptionHandle(string name)
                => _handles.FirstOrDefault(h => h.Field.Name == name) ?? throw new KeyNotFoundException(name);
        }

        private sealed class UnsubscribeBuilder : IPlcUnsubscriptionRequestBuilder
        {
            private readonly PollingSubscriptionManager _owner;
            private readonly List<Handle> _handles = new List<Handle>();
            public UnsubscribeBuilder(PollingSubscriptionManager owner) => _owner = owner;
            public IPlcUnsubscriptionRequestBuilder AddHandles(params IPlcSubscriptionHandle[] handles)
                => AddHandles((IEnumerable<IPlcSubscriptionHandle>)handles);
            public IPlcUnsubscriptionRequestBuilder AddHandles(IEnumerable<IPlcSubscriptionHandle> handles)
            {
                ArgumentNullException.ThrowIfNull(handles);
                var validated = handles.Select(h => h is Handle handle && handle.Owner == _owner
                    ? handle : throw new ArgumentException("Subscription belongs to another connection.")).ToArray();
                _handles.AddRange(validated);
                return this;
            }
            public IPlcUnsubscriptionRequest Build() => new UnsubscribeRequest(_owner, _handles.Distinct().ToArray());
        }

        private sealed class UnsubscribeRequest : IPlcUnsubscriptionRequest
        {
            private readonly PollingSubscriptionManager _owner;
            private readonly Handle[] _handles;
            public UnsubscribeRequest(PollingSubscriptionManager owner, Handle[] handles) { _owner = owner; _handles = handles; }
            public async Task<IPlcResponse> ExecuteAsync()
            {
                await _owner.UnsubscribeAsync(_handles).ConfigureAwait(false);
                return new UnsubscribeResponse(this);
            }
        }

        private sealed record UnsubscribeResponse(IPlcRequest Request) : IPlcResponse;

        private sealed class Handle : IPlcSubscriptionHandle
        {
            private readonly object _listenersGate = new object();
            private EventHandler<IPlcSubscriptionEventArgs>? _listeners;
            private object? _previous;
            private bool _hasPrevious;
            public PollingSubscriptionManager Owner { get; }
            public Field Field { get; }
            public volatile bool Active = true;
            public long Next;
            public Task InFlight = Task.CompletedTask;
            public Handle(PollingSubscriptionManager owner, Field field)
            {
                Owner = owner;
                Field = field;
                Next = Environment.TickCount64 + (long)Math.Ceiling(field.Interval.TotalMilliseconds);
            }
            public event EventHandler<IPlcSubscriptionEventArgs> PlcSubscriptionEvent
            {
                add { lock (_listenersGate) { if (Active) _listeners += value; } }
                remove { lock (_listenersGate) _listeners -= value; }
            }
            public void ClearListeners() { lock (_listenersGate) _listeners = null; }
            public void Publish(SubscriptionValueEventArgs result)
            {
                lock (_listenersGate)
                {
                    if (!Active) return;
                    bool notify;
                    try
                    {
                        var snapshot = ValueSnapshot.Create(result.Value);
                        notify = !Field.ChangesOnly || result.ResponseCode != PlcResponseCode.Ok ||
                            (_hasPrevious && !StructuralComparisons.StructuralEqualityComparer.Equals(_previous, snapshot));
                        if (result.ResponseCode == PlcResponseCode.Ok)
                        {
                            _previous = snapshot;
                            _hasPrevious = true;
                        }
                    }
                    catch (Exception ex) { Owner.ReportError(ex); return; }
                    if (!notify || _listeners == null) return;
                    Owner._dispatching.Value = true;
                    try
                    {
                        foreach (EventHandler<IPlcSubscriptionEventArgs> listener in _listeners.GetInvocationList())
                        {
                            if (!Active) break;
                            try { listener(this, result); }
                            catch (Exception ex) { Owner.ReportError(ex); }
                        }
                    }
                    finally { Owner._dispatching.Value = false; }
                }
            }
        }
    }
}
