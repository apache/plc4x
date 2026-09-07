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

# Subscription Emulation

This module provides polling-based subscription emulation for drivers that don't natively support subscriptions.

## PollingSubscriptionConnectionBase

The `PollingSubscriptionConnectionBase` class provides polling-based subscription emulation for drivers that don't natively support subscriptions (such as Modbus, BACnet, etc.).

### Features

- **CYCLIC subscriptions** - Polls tags at specified intervals
- **CHANGE_OF_STATE subscriptions** - Polls tags and only fires events when values change
- **EVENT subscriptions** - Not supported (throws `UnsupportedOperationException`)
- **Automatic grouping** - Tags with the same subscription type and polling interval are grouped into a single batch for efficiency
- **Consumer registration** - Supports request-level, tag-specific, and registered consumers
- **Thread-safe** - Uses concurrent data structures for thread-safe operation

### Basic Usage

To add polling-based subscription support to your driver:

1. Add this dependency to your driver's `pom.xml`:

```xml
<dependency>
    <groupId>org.apache.plc4x</groupId>
    <artifactId>plc4j-utils-subscription-emulation</artifactId>
    <version>1.0.0-SNAPSHOT</version>
</dependency>
```

2. Extend `PollingSubscriptionConnectionBase` instead of `ConnectionBase` and implement `PlcReader`:

```java
import org.apache.plc4x.java.utils.subscriptionemulation.PollingSubscriptionConnectionBase;

public class ModbusTcpConnection
        extends PollingSubscriptionConnectionBase<ModbusTcpConfiguration>
        implements PlcReader, PlcWriter {

    public ModbusTcpConnection(ModbusTcpConfiguration configuration,
                               TransportInstance<?> transportInstance) {
        super(configuration, transportInstance);
    }

    // Implement PlcReader interface
    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        // Your read implementation
    }

    // Rest of your driver implementation...
}
```

That's it! Your driver now supports subscriptions via polling.

### How It Works

When a subscription is created:

1. The base class groups tags by subscription type and polling interval
2. For each group, it creates an EventPump batch with a TimerTrigger
3. The batch polls the tags using your driver's `read()` method
4. For CYCLIC subscriptions, events are fired on every poll
5. For CHANGE_OF_STATE subscriptions, events are only fired when values change

### Customization

You can override these methods to customize behavior:

#### Custom Default Polling Interval

```java
@Override
protected long getDefaultPollingInterval() {
    return 500; // 500ms instead of default 1000ms
}
```

#### Custom Value Comparison

```java
@Override
protected boolean valuesEqual(PlcValue v1, PlcValue v2) {
    // Custom comparison logic
    // For example, treat small floating point differences as equal
    if (v1 == v2) return true;
    if (v1 == null || v2 == null) return false;

    if (v1.isDouble() && v2.isDouble()) {
        double diff = Math.abs(v1.getDouble() - v2.getDouble());
        return diff < 0.001; // Treat differences < 0.001 as equal
    }

    return Objects.equals(v1.getObject(), v2.getObject());
}
```

### Usage Examples

#### CYCLIC Subscription

```java
PlcSubscriptionRequest request = DefaultPlcSubscriptionRequest.builder()
    .addCyclicTagAddress("temperature", "MAIN.temperature", Duration.ofMillis(100))
    .addCyclicTagAddress("pressure", "MAIN.pressure", Duration.ofMillis(100))
    .setConsumer(event -> {
        System.out.println("Temperature: " + event.getInteger("temperature"));
        System.out.println("Pressure: " + event.getInteger("pressure"));
    })
    .build();

CompletableFuture<PlcSubscriptionResponse> future = connection.subscribe(request);
PlcSubscriptionResponse response = future.get();

// Later, to unsubscribe:
PlcUnsubscriptionRequest unsubRequest = DefaultPlcUnsubscriptionRequest.builder()
    .addHandles(response.getSubscriptionHandle("temperature"))
    .addHandles(response.getSubscriptionHandle("pressure"))
    .build();

connection.unsubscribe(unsubRequest).get();
```

#### CHANGE_OF_STATE Subscription

```java
PlcSubscriptionRequest request = DefaultPlcSubscriptionRequest.builder()
    .addChangeOfStateTagAddress("alarm", "MAIN.alarm", Duration.ofMillis(50))
    .setConsumer(event -> {
        System.out.println("Alarm state changed: " + event.getBoolean("alarm"));
    })
    .build();

connection.subscribe(request);
```

#### Tag-Specific Consumers

```java
PlcSubscriptionRequest request = DefaultPlcSubscriptionRequest.builder()
    .addCyclicTagAddress("tag1", "MAIN.tag1", Duration.ofMillis(100))
    .setTagConsumer("tag1", event -> {
        System.out.println("Tag1 event: " + event.getInteger("tag1"));
    })
    .addCyclicTagAddress("tag2", "MAIN.tag2", Duration.ofMillis(100))
    .setTagConsumer("tag2", event -> {
        System.out.println("Tag2 event: " + event.getInteger("tag2"));
    })
    .build();

connection.subscribe(request);
```

#### Consumer Registration

```java
// Subscribe
PlcSubscriptionResponse response = connection.subscribe(request).get();
PlcSubscriptionHandle handle1 = response.getSubscriptionHandle("tag1");
PlcSubscriptionHandle handle2 = response.getSubscriptionHandle("tag2");

// Register a consumer for specific handles
PlcConsumerRegistration registration = connection.registerConsumer(
    event -> {
        System.out.println("Registered consumer received event");
    },
    Arrays.asList(handle1, handle2)
);

// Later, unregister
connection.unregisterConsumer(registration);
```

### Performance Considerations

- **Grouping**: Tags with the same subscription type and polling interval are automatically grouped into a single batch, reducing overhead
- **Concurrency**: All subscription management uses thread-safe concurrent data structures
- **Polling Intervals**: Choose appropriate polling intervals to balance responsiveness and network/PLC load
  - Too frequent polling can overload the PLC or network
  - Too infrequent polling can miss value changes in CHANGE_OF_STATE subscriptions

### Requirements

Your driver must:
- Extend `PollingSubscriptionConnectionBase` instead of `ConnectionBase`
- Implement the `PlcReader` interface
- Add the `plc4j-utils-subscription-emulation` dependency

### Limitations

- **EVENT subscriptions are not supported** - Only CYCLIC and CHANGE_OF_STATE are supported
- **Polling overhead** - All subscriptions are emulated via polling, which may not be as efficient as native protocol subscriptions
- **Value change detection** - CHANGE_OF_STATE subscriptions compare values using `Objects.equals()` by default, which may not be suitable for all data types (e.g., floating point numbers with small differences)

## Implementation Details

For implementation details, see:
- `PollingSubscriptionConnectionBase.java:73` - Main implementation
- `EventPump` - The underlying polling engine (see `java/tools/event-pump/README.md`)
