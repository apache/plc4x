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

The `ConnectionCache` implements the `PlcConnectionManager` interface, just the same way the `PlcDriverManager` does, so it can generally be used instead of the un-cached version.

## Architecture

The `CachedPlcConnectionManager` contains a map of `ConnectionContainer` objects.
Each of these generally have a reference to a real `PlcConnection` as well as all properties for managing it's state.

In general there are just three properties:

- A reference to the connection this container handles
- A reference to the current connection-lease (`null`, if the `ConnectionContainer` is idle)
- A queue where all further lease-requests are lined up

Whenever a `PlcConnection` is required, instead of returning a real `PlcConnection`, the `CachedPlcConnectionManager` returns a `LeasedPlcConnection`. 

This object is a volatile container for a `PlcConnection`, allowing the container to invalidate the `PlcConnection`. 

Whenever a `CachedPlcConnectionManager`'s `getConnection` method is used, it returns a new instance of such a container. Whenever the client calls `close` on this connection, it is however not really closed, but the reference to the real connection is cleared, hereby rendering the connection-lease useless and the connection is returned to the `ConnectionContainer`. Also, if the client holds on to the connection-lease for longer than the `maxLeaseTime` the container invalidates the connection-lease. 

If a `CachedPlcConnectionManager` is used to get a connection that is currently being used, instead of returning a reference to it, a Future is generated and added to a queue. As soon as the connection is released, the container checks if there are any requests waiting.
If there are it takes the oldest request and completes that future with a new lease.