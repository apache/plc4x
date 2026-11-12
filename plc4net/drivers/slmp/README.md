<!--
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements. See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License. You may obtain a copy of the License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# SLMP polling subscriptions

The PLC4Net SLMP connection emulates subscriptions using its existing 3E
binary TCP batch reads. Cyclic subscriptions publish every sampled result.
Change-of-state subscriptions use the first successful sample as a baseline
and publish subsequent changes. Lists and structures are compared by content;
failed reads publish their response code without replacing that baseline.
These are sampled values: changes between reads can be missed.

```csharp
using org.apache.plc4net.messages;
using org.apache.plc4net.spi.drivers.subscriptions;

// connection is an established SLMP IPlcConnection.
var subscribed = (IPlcSubscriptionResponse)await connection.SubscriptionRequestBuilder!
    .AddCyclicField("temperature", "D100:REAL", TimeSpan.FromSeconds(1))
    .AddChangeOfStateField("running", "M100:BOOL")
    .Build().ExecuteAsync();

var handle = subscribed.GetSubscriptionHandle("temperature");
handle.PlcSubscriptionEvent += (_, args) =>
{
    var sample = (SubscriptionValueEventArgs)args;
    Console.WriteLine($"{sample.Timestamp:o}: {sample.Name} {sample.ResponseCode}");
    // Read sample.Value only when sample.ResponseCode is PlcResponseCode.Ok.
};

await connection.UnsubscriptionRequestBuilder!
    .AddHandles(subscribed.SubscriptionHandles).Build().ExecuteAsync();
```

Change-of-state sampling defaults to 1000 ms. Configure it with
`subscription-polling-interval-ms`, for example
`slmp:tcp://localhost:5007?subscription-polling-interval-ms=250`.
Each cyclic field supplies its own interval. Intervals must be at least 1 ms
and at most `int.MaxValue` milliseconds. Native event subscriptions
(`AddEventField`) are unsupported and throw `NotSupportedException`.

One scheduler per connection serializes polling reads. Each field waits its
interval after its previous read completes; slow reads do not accumulate
timer ticks. Reads still share the SLMP exchange gate with ordinary reads and
writes. Callbacks execute serially on the scheduler, so keep handlers short.
An exception in one handler is logged and does not prevent other handlers.

Unsubscription disables notifications and waits for an already-started read
to finish. It does not cancel that read: 3E frames cannot correlate a late
response after cancellation. Closing the connection stops its subscriptions
before closing the transport. Calling unsubscribe or close from a callback
is supported. The existing request timeout bounds SLMP reads; a timed-out
SLMP connection remains unusable and must be replaced by the application.

This change is covered by in-memory protocol and lifecycle tests. It has not
been verified against physical Mitsubishi hardware.
