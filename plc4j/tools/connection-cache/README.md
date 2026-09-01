<!--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
  -->

# Connection Cache

In contrast to typical database-connections, plc-connections have numerous disadvantages:

1. The number of connections a PLC is able to accept might be very limited.
2. Connecting can require a number of round-trips and a lot of data being transmitted.

If multiple places in the application want to access a PLC, the option of `every part opens its own connection` will quickly drain the resources of the PLC. Opening and closing the connections in order to reduce this, will waste a lot of resources.

Therefore, it's the goal of the connection-cache to allow gaining access to a PLC the usual way, however as soon as a client is finished with its work, it doesn't close the connection, but gives it back to the cache for the next client to be able to use the same connection without having to re-connect.

The `PlcConnectionCache` implements the `PlcConnectionManager` interface, which extends the `PlcConnectionFactory` interface the `PlcDriverManager` implements, so it can generally be used instead of the un-cached version. It does not create connections itself: it is built with a `PlcConnectionFactory` (usually the driver manager) that it delegates the actual connecting to.

For how to configure and use it, see the [Connection Cache page on the website](https://plc4x.apache.org/plc4x/latest/users/tools/connection-cache.html). This document describes how it works inside.

## Architecture

The `PlcConnectionCache` contains a map of `ConnectionContainer` objects.
Each of these generally have a reference to a real `PlcConnection` as well as all properties for managing it's state.

In general there are just three properties:

- A reference to the connection this container handles
- A reference to the current connection-lease (`null`, if the `ConnectionContainer` is idle)
- A queue where all further lease-requests are lined up

Whenever a `PlcConnection` is required, instead of returning a real `PlcConnection`, the `PlcConnectionCache` returns a `LeasedPlcConnection`. 

This object is a volatile container for a `PlcConnection`, allowing the container to invalidate the `PlcConnection`. 

Whenever a `PlcConnectionCache`'s `getConnection` method is used, it returns a new instance of such a container. Whenever the client calls `close` on this connection, it is however not really closed, but the reference to the real connection is cleared, hereby rendering the connection-lease useless and the connection is returned to the `ConnectionContainer`. Also, if the client holds on to the connection-lease for longer than the `maxLeaseTime` the container invalidates the connection-lease. 

If a `PlcConnectionCache` is used to get a connection that is currently being used, instead of returning a reference to it, a Future is generated and added to a queue. As soon as the connection is released, the container checks if there are any requests waiting.
If there are it takes the oldest request and completes that future with a new lease.

Before handing out a connection that has been sitting idle for longer than the `idlePingThreshold`, the container validates it with a `ping()`. If that fails, the connection is closed and a new one is established, transparently to the client asking for the lease.

## Timeouts

Every timeout is driven by a `ScheduledExecutorService` the cache owns (or one supplied through the builder), rather than by threads of its own:

- **idle timeout** - a container whose connection has not been leased for `maxIdleTime` closes it and drops out of the cache. The next request for that connection string builds a new container.
- **lease timeout** - a lease not returned within `maxLeaseTime` is invalidated, which is what keeps a client that forgot to `close()` from holding a connection forever.
- **wait timeout** - each queued lease-request carries its own scheduled task, so a waiter that is never served gives up after `maxWaitTime` and its caller gets a `PlcConnectionException`.

Two further bounds protect the container's lock from a wedged socket, because connect and close are executed while holding it:

- `connectBounded()` runs the connect on a separate thread and gives up after `maxWaitTime`; a connection that arrives after the caller gave up is closed rather than leaked.
- `closeBounded()` does the same for closing, abandoning a close that exceeds `closeTimeout` to a background daemon thread so the cache can carry on. Both bounds can be disabled by configuring the respective timeout as `0`.

## Surviving a reconnect

Subscriptions and event listeners belong to a connection, so recreating one underneath a client would silently lose them. Each `ConnectionContainer` therefore owns a `ConnectionStateTracker`, and the `LeasedPlcConnection` wraps the subscription and unsubscription request builders to report what the client registers.

When the container establishes a replacement connection - after a failed validation ping, or after the previous connection was invalidated - it asks the tracker to restore that state on the new connection. Each subscription is held as a `SubscriptionRecord`, which keeps both the handle the client originally received and the handle currently valid for each tag, so a client still holding a pre-reconnect handle can be mapped onto the new one (for example to unsubscribe).

## Shutting down

The cache holds the connections it hands out, so it has to be closed. `close()` shuts the scheduler down first - so no timeout task can fire into a half-torn-down cache - then closes every container and empties the map. A closed cache rejects further `getConnection()` calls with a `PlcConnectionCacheClosedException`, and closing it again does nothing.

Single connections can be evicted without closing the cache using `removeCachedConnection(connectionString)`, which force-closes the connection even if it is currently leased.
