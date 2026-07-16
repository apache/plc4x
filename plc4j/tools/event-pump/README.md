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

# Event Pump

An event-driven system for fetching PLC tag data and delivering it to consumers.

## Overview

The Event Pump is a flexible component that allows you to fetch sets of tags from PLCs based on various triggers. One instance can manage multiple batches of tags from the same PLC or from multiple PLCs simultaneously.

### Key Features

- **Multiple Batch Support**: Fetch different sets of tags with different triggers
- **Multi-PLC Capability**: Connect to and fetch from multiple PLCs simultaneously
- **Flexible Triggers**:
  - **Time-based**: Fetch at regular intervals
  - **Subscription-based**: Fetch when a tag value changes
  - **Custom**: Implement your own trigger logic
- **Configuration Support**: Define batches in XML, JSON, or YAML
- **Thread-Safe**: Concurrent batch management
- **AutoCloseable**: Proper resource cleanup

## Quick Start

### Add Dependency

```xml
<dependency>
    <groupId>org.apache.plc4x</groupId>
    <artifactId>plc4j-tools-event-pump</artifactId>
    <version>1.0.0-SNAPSHOT</version>
</dependency>
```

### Basic Usage (Programmatic)

```java
import org.apache.plc4x.java.tools.eventpump.*;
import org.apache.plc4x.java.tools.eventpump.triggers.TimerTrigger;
import org.apache.plc4x.java.tools.eventpump.triggers.Trigger;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.spi.PlcDriverManager;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

// Create connection
PlcConnectionManager connectionManager = PlcDriverManager.getDefault();
PlcConnection connection = connectionManager.getConnection("ads:tcp://192.168.1.1");

// Create a batch using the Builder pattern
TagBatch batch = TagBatch.builder()
    .withBatchId("batch1")
    .withConnection(connection)
    .addTag("temperature", "MAIN.temperature")
    .addTag("pressure", "MAIN.pressure")
    .addTag("status", "MAIN.status")
    .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
    .withListener((batchRef, response) -> {
        System.out.println("Temperature: " + response.getInteger("temperature"));
        System.out.println("Pressure: " + response.getInteger("pressure"));
        System.out.println("Status: " + response.getBoolean("status"));
    })
    .build();

// Add to pump and start
EventPump pump = new EventPump();
pump.addBatch(batch);
pump.startAll();

// ... application runs ...

// Cleanup
pump.close();
connection.close();
```

### Configuration-Based Usage

Create a configuration file (YAML example):

```yaml
# config.yaml
connections:
  - id: plc1
    url: "ads:tcp://192.168.1.1"
    username: admin
    password: secret

batches:
  - id: sensor-batch
    connectionId: plc1
    tags:
      temperature: "MAIN.temperature"
      pressure: "MAIN.pressure"
      humidity: "MAIN.humidity"
    trigger:
      type: timer
      intervalSeconds: 5
      initialDelaySeconds: 0

  - id: alarm-batch
    connectionId: plc1
    tags:
      alarmStatus: "MAIN.alarmStatus"
      alarmCode: "MAIN.alarmCode"
    trigger:
      type: subscription
      tagName: alarmTrigger
      tagAddress: "MAIN.alarmTrigger"
```

Load and use the configuration:

```java
import org.apache.plc4x.java.tools.eventpump.config.*;
import java.io.File;

// Load configuration
PlcConnectionManager connectionManager = PlcDriverManager.getDefault();

EventPump pump = EventPumpFactory.fromYaml(
    new File("config.yaml"),
    connectionManager,
    (batch, response) -> {
        System.out.println("Batch '" + batch.getBatchId() + "' fetched:");
        response.getTagNames().forEach(tagName -> {
            System.out.println("  " + tagName + ": " + response.getObject(tagName));
        });
    }
);

// Start all batches
pump.startAll();

// ... application runs ...

// Cleanup
pump.close();
```

## Core Concepts

### EventPump

The main component that manages multiple tag batches. Think of it as a container for batches.

**Key Methods:**
- `addBatch(TagBatch)`: Add a batch to the pump
- `removeBatch(String)`: Remove a batch by ID
- `getBatch(String)`: Get a specific batch by ID
- `startBatch(String)`: Start a specific batch
- `stopBatch(String)`: Stop a specific batch
- `getAllBatches()`: Get all batches as a map
- `getBatchIds()`: Get all batch IDs as a set
- `getBatchCount()`: Get the number of batches
- `getStartedBatchCount()`: Get the number of started batches
- `startAll()`: Start all batches
- `stopAll()`: Stop all batches
- `close()`: Stop and cleanup all batches

### TagBatch

Represents a batch of tags to be fetched from a PLC. Each batch has:
- A unique ID
- A PLC connection
- A map of tag names to tag addresses
- A trigger that determines when to fetch
- A listener that receives results

**Creating a TagBatch:**

Use the Builder pattern to create TagBatch instances:

```java
TagBatch batch = TagBatch.builder()
    .withBatchId("my-batch")
    .withConnection(connection)
    .addTag("tag1", "MAIN.tag1")
    .addTag("tag2", "MAIN.tag2")
    .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
    .withListener((b, response) -> {
        // Handle results
    })
    .build();
```

**Key Methods:**
- `start()`: Start fetching (starts the trigger)
- `stop()`: Stop fetching (stops the trigger)
- `fetchTags()`: Manually trigger a fetch (regardless of trigger state)
- `setListener(TagBatchListener)`: Set/update the result handler (can also be set via builder)

### Triggers

Triggers determine when a batch should be fetched.

#### TimerTrigger

Fires at regular time intervals.

```java
// Fetch every 10 seconds
Trigger trigger = new TimerTrigger(10, TimeUnit.SECONDS);

// Fetch every 5 seconds, with a 2 second initial delay
Trigger trigger = new TimerTrigger(5, 2, TimeUnit.SECONDS);
```

#### SubscriptionTrigger

Fires when a subscribed tag value changes.

```java
// Trigger on any change
Trigger trigger = new SubscriptionTrigger(
    connection,
    "triggerTag",
    "MAIN.trigger"
);

// Trigger only when a value changes to true
Trigger trigger = new SubscriptionTrigger(
    connection,
    "triggerTag",
    "MAIN.trigger",
    event -> event.getBoolean("triggerTag") == Boolean.TRUE
);
```

#### Custom Triggers

Implement the `Trigger` interface to create custom trigger logic:

```java
public class MyCustomTrigger implements Trigger {
    // Implement trigger methods
    @Override
    public void start(TriggerListener listener) {
        // Start your trigger logic
        // Call listener.onTrigger(this) when it should fire
    }

    @Override
    public void stop() {
        // Stop your trigger logic
    }

    // ... other methods
}
```

### Listeners

The `TagBatchListener` interface handles fetch results. You can set a listener via the builder or using `setListener()`:

```java
// Option 1: Set listener via builder (recommended)
TagBatch batch = TagBatch.builder()
    .withBatchId("batch1")
    .withConnection(connection)
    .addTag("temp", "MAIN.temperature")
    .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
    .withListener(new TagBatch.TagBatchListener() {
        @Override
        public void onTagsFetched(TagBatch batch, PlcReadResponse response) {
            // Process successful fetch
            for (String tagName : response.getTagNames()) {
                Object value = response.getObject(tagName);
                System.out.println(tagName + " = " + value);
            }
        }

        @Override
        public void onError(TagBatch batch, Throwable error) {
            // Handle errors
            System.err.println("Error in batch " + batch.getBatchId() + ": " + error.getMessage());
        }
    })
    .build();

// Option 2: Use a lambda (for simple cases)
TagBatch batch2 = TagBatch.builder()
    .withBatchId("batch2")
    .withConnection(connection)
    .addTag("pressure", "MAIN.pressure")
    .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
    .withListener((b, response) -> {
        // Process results
        System.out.println("Pressure: " + response.getObject("pressure"));
    })
    .build();

// Option 3: Set/update listener after creation
batch.setListener((b, response) -> {
    // Process results
});
```

## Configuration Formats

The Event Pump supports configuration in three formats: YAML, JSON, and XML.

### YAML Configuration

```yaml
connections:
  - id: plc1
    url: "ads:tcp://192.168.1.1"
  - id: plc2
    url: "modbus:tcp://192.168.1.2"

batches:
  - id: batch1
    connectionId: plc1
    tags:
      temp1: "MAIN.temperature1"
      temp2: "MAIN.temperature2"
    trigger:
      type: timer
      intervalSeconds: 10

  - id: batch2
    connectionId: plc2
    tags:
      pressure: "40001"
    trigger:
      type: subscription
      tagName: trigger
      tagAddress: "40100"
```

### JSON Configuration

```json
{
  "connections": [
    {
      "id": "plc1",
      "url": "ads:tcp://192.168.1.1"
    }
  ],
  "batches": [
    {
      "id": "batch1",
      "connectionId": "plc1",
      "tags": {
        "temp1": "MAIN.temperature1",
        "temp2": "MAIN.temperature2"
      },
      "trigger": {
        "type": "timer",
        "intervalSeconds": 10
      }
    }
  ]
}
```

### XML Configuration

```xml
<EventPumpConfiguration>
  <connections>
    <item>
      <id>plc1</id>
      <url>ads:tcp://192.168.1.1</url>
    </item>
  </connections>
  <batches>
    <item>
      <id>batch1</id>
      <connectionId>plc1</connectionId>
      <tags>
        <temp1>MAIN.temperature1</temp1>
        <temp2>MAIN.temperature2</temp2>
      </tags>
      <trigger>
        <type>timer</type>
        <intervalSeconds>10</intervalSeconds>
      </trigger>
    </item>
  </batches>
</EventPumpConfiguration>
```

### Trigger Configuration

#### Timer Trigger

```yaml
trigger:
  type: timer
  intervalSeconds: 5        # Required: interval in seconds
  initialDelaySeconds: 2    # Optional: initial delay in seconds
```

#### Subscription Trigger

```yaml
trigger:
  type: subscription
  tagName: triggerTag       # Required: name for the subscription
  tagAddress: "MAIN.trigger"  # Required: tag address to subscribe to
  condition: "value == true"  # Optional: condition expression (not yet implemented)
```

## Value Transformations

The Event Pump supports optional value transformations, allowing you to transform raw PLC values using expressions before they reach your listener. This is useful for unit conversions, calculations, and data normalization.

### Overview

Value transformations are applied after values are read from the PLC, but before they're delivered to your listener. Each tag can have an optional transformation expression that manipulates the raw value.

**Key Features:**
- **Optional**: Transformations are completely optional - use them only where needed
- **Per-Tag**: Each tag can have its own transformation or none at all
- **Pluggable**: Default simple expression evaluator, but you can register custom transformers (e.g., SpEL for Spring applications)
- **Context-Aware**: Transformations can access the current value and all other tag values in the batch

### Basic Usage

#### Extended Tag Format

To use transformations, use the extended tag format in your configuration:

```yaml
batches:
  - id: sensor-batch
    connectionId: plc1
    tags:
      # Simple format (no transformation)
      humidity: "MAIN.humidity"

      # Extended format with transformation
      temperature:
        address: "MAIN.temperatureCelsius"
        transform: "value * 1.8 + 32"  # Convert Celsius to Fahrenheit

      pressure:
        address: "MAIN.pressureBar"
        transform: "value * 14.5038"  # Convert bar to PSI
```

#### Programmatic Usage

```java
TagBatch batch = TagBatch.builder()
    .withBatchId("sensor-batch")
    .withConnection(connection)
    .addTag("humidity", "MAIN.humidity")  // No transformation
    .addTag("temperature", "MAIN.temperatureCelsius")
    .addTransform("temperature", "value * 1.8 + 32")  // Transform to Fahrenheit
    .addTag("pressure", "MAIN.pressureBar")
    .addTransform("pressure", "value * 14.5038")  // Transform to PSI
    .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
    .withListener((b, response) -> {
        // Values are already transformed
        System.out.println("Temperature: " + response.getFloat("temperature") + "°F");
        System.out.println("Pressure: " + response.getFloat("pressure") + " PSI");
    })
    .build();
```

### Expression Syntax

The default `SimpleExpressionEvaluator` supports the following operations:

#### Arithmetic Operators
- `+` : Addition
- `-` : Subtraction
- `*` : Multiplication
- `/` : Division
- `%` : Modulo

#### Comparison Operators
- `==` : Equal to
- `!=` : Not equal to
- `<`  : Less than
- `>`  : Greater than
- `<=` : Less than or equal to
- `>=` : Greater than or equal to

#### Logical Operators
- `&&` : Logical AND
- `||` : Logical OR
- `!`  : Logical NOT

#### Other Features
- **Parentheses** for grouping: `(value + 10) * 2`
- **Numeric literals**: Integer and floating-point numbers
- **Boolean literals**: `true`, `false`
- **Variable reference**: `value` (the current tag value)

#### Operator Precedence (highest to lowest)
1. Parentheses `( )`
2. Unary operators `!`, `-`
3. Multiplicative `*`, `/`, `%`
4. Additive `+`, `-`
5. Comparison `<`, `>`, `<=`, `>=`
6. Equality `==`, `!=`
7. Logical AND `&&`
8. Logical OR `||`

### Common Examples

#### Temperature Conversion

```yaml
tags:
  # Celsius to Fahrenheit
  tempF:
    address: "MAIN.tempCelsius"
    transform: "value * 1.8 + 32"

  # Fahrenheit to Celsius
  tempC:
    address: "MAIN.tempFahrenheit"
    transform: "(value - 32) / 1.8"
```

#### Pressure Conversion

```yaml
tags:
  # Bar to PSI
  pressurePSI:
    address: "MAIN.pressureBar"
    transform: "value * 14.5038"

  # kPa to Bar
  pressureBar:
    address: "MAIN.pressureKpa"
    transform: "value / 100"
```

#### Scaling and Offset

```yaml
tags:
  # Scale raw sensor value
  scaledValue:
    address: "MAIN.rawSensor"
    transform: "value * 0.001"  # Convert from millivolts to volts

  # Apply offset and scale
  calibratedValue:
    address: "MAIN.rawValue"
    transform: "(value - 100) * 2.5"
```

#### Boolean Transformations

```yaml
tags:
  # Status to boolean
  isRunning:
    address: "MAIN.status"
    transform: "value == 1"

  # Alarm detection
  alarmActive:
    address: "MAIN.alarmCode"
    transform: "value != 0"

  # Within range check
  withinRange:
    address: "MAIN.pressure"
    transform: "value >= 90 && value <= 110"
```

### Pluggable Transformers

The transformation system is pluggable, allowing you to register custom expression evaluators. A Spring Expression Language (SpEL) transformer is included and **automatically registered** when Spring is detected on the classpath.

#### Available Transformers

By default, the following transformers are available:
- `simple`: The default SimpleExpressionEvaluator (always registered)

**Note**: The SpEL transformer is provided as an **optional** dependency. It will only be available if you add Spring Expression to your project dependencies.

#### Spring Expression Language (SpEL)

If Spring Expression is available on your classpath, the SpEL transformer is automatically registered and provides advanced expression capabilities:

**Automatic Registration**: Add this dependency to enable SpEL:
```xml
<dependency>
    <groupId>org.springframework</groupId>
    <artifactId>spring-expression</artifactId>
    <version>6.2.15</version>
</dependency>
```

**SpEL Syntax**: Variables must be accessed using `#variableName` syntax:

```yaml
tags:
  tempF:
    address: "MAIN.tempCelsius"
    transform: "#value * 1.8 + 32"  # Note the # prefix

  average:
    address: "MAIN.temp1"
    transform: "(#temp1 + #temp2) / 2"  # Multiple variables with #
```

**SpEL Features**:
- Ternary operator: `#value > 10 ? #value : 0`
- Method calls: `#value.toUpperCase()`, `T(java.lang.Math).sqrt(#value)`
- String operations: `#value + ' suffix'`
- Elvis operator: `#value ?: 100`
- Complex logic: `#value > 50 && #value <= 100 ? 'normal' : 'abnormal'`

#### Using SpEL as Default

To use SpEL for all transformations by default:

```java
import org.apache.plc4x.java.tools.eventpump.transform.*;

// Create a registry with default transformers (SpEL auto-registered if Spring is on the classpath)
ValueTransformerRegistry registry = ValueTransformerRegistry.createDefault();

// Use this registry for all batches in your EventPump
TagBatch batch = TagBatch.builder()
    .withBatchId("batch1")
    .withConnectionManager(connectionManager)
    .withConnectionString("ads:tcp://192.168.1.1")
    .addTag("temp", "MAIN.temp")
    .addTransform("temp", "#value * 1.8 + 32")
    .withTransformerRegistry(registry)  // Share registry across batches
    .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
    .withListener((b, r) -> handleResults(r))
    .build();
```

#### Sharing Registry Across Multiple Batches

For efficiency and consistency, you can share one registry instance across all batches:

```java
// Create one registry for the entire EventPump
ValueTransformerRegistry registry = ValueTransformerRegistry.createDefault();

// Use it for all batches
TagBatch batch1 = TagBatch.builder()
    .withBatchId("batch1")
    .withTransformerRegistry(registry)  // Share the same registry
    // ... other configuration
    .build();

TagBatch batch2 = TagBatch.builder()
    .withBatchId("batch2")
    .withTransformerRegistry(registry)  // Share the same registry
    // ... other configuration
    .build();
```

**Note**: When using `EventPumpFactory` to create batches from configuration files, a single registry is automatically created and shared across all batches in that EventPump.

#### Checking Registered Transformers

```java
ValueTransformerRegistry registry = ValueTransformerRegistry.createDefault();

// Get all registered transformer names
Set<String> names = registry.getRegisteredNames();
System.out.println("Available transformers: " + names);

// Get the default transformer
ValueTransformer defaultTransformer = registry.getDefault();
System.out.println("Default transformer: " + defaultTransformer.getName());
```

#### Custom Transformers

You can also implement your own transformers by implementing the `ValueTransformer` interface:

```java
public class MyCustomTransformer implements ValueTransformer {
    @Override
    public PlcValue transform(String expression, Map<String, PlcValue> context)
            throws TransformException {
        // Your custom transformation logic
        return transformedValue;
    }

    @Override
    public String getName() {
        return "custom";
    }
}

// Register your custom transformer
ValueTransformerRegistry.getInstance().register(new MyCustomTransformer());
```

### Configuration Examples

#### YAML with Transformations

```yaml
connections:
  - id: plc1
    url: "ads:tcp://192.168.1.1"

batches:
  - id: sensor-batch
    connectionId: plc1
    tags:
      # Mix of simple and extended formats
      humidity: "MAIN.humidity"

      tempF:
        address: "MAIN.tempCelsius"
        transform: "value * 1.8 + 32"

      pressurePSI:
        address: "MAIN.pressureBar"
        transform: "value * 14.5038"

      # Extended format without transformation
      status:
        address: "MAIN.status"

    trigger:
      type: timer
      intervalSeconds: 5
```

#### JSON with Transformations

```json
{
  "batches": [
    {
      "id": "sensor-batch",
      "connectionId": "plc1",
      "tags": {
        "humidity": "MAIN.humidity",
        "tempF": {
          "address": "MAIN.tempCelsius",
          "transform": "value * 1.8 + 32"
        },
        "pressurePSI": {
          "address": "MAIN.pressureBar",
          "transform": "value * 14.5038"
        }
      },
      "trigger": {
        "type": "timer",
        "intervalSeconds": 5
      }
    }
  ]
}
```

#### XML with Transformations

```xml
<batches>
  <item>
    <id>sensor-batch</id>
    <connectionId>plc1</connectionId>
    <tags>
      <humidity>MAIN.humidity</humidity>
      <tempF>
        <address>MAIN.tempCelsius</address>
        <transform>value * 1.8 + 32</transform>
      </tempF>
      <pressurePSI>
        <address>MAIN.pressureBar</address>
        <transform>value * 14.5038</transform>
      </pressurePSI>
    </tags>
    <trigger>
      <type>timer</type>
      <intervalSeconds>5</intervalSeconds>
    </trigger>
  </item>
</batches>
```

### Use Cases

#### Industrial Monitoring

```yaml
batches:
  - id: production-monitoring
    connectionId: plc1
    tags:
      # Convert temperatures to preferred unit
      ovenTempF:
        address: "MAIN.ovenTempC"
        transform: "value * 1.8 + 32"

      # Calculate pressure deviation from setpoint
      pressureDeviation:
        address: "MAIN.actualPressure"
        transform: "value - 100"  # Assuming 100 is setpoint

      # Normalize percentage
      efficiency:
        address: "MAIN.efficiencyRaw"
        transform: "value / 100"
```

#### Energy Management

```yaml
batches:
  - id: energy-monitoring
    connectionId: plc1
    tags:
      # Convert watts to kilowatts
      powerKW:
        address: "MAIN.powerW"
        transform: "value / 1000"

      # Calculate power factor
      powerFactor:
        address: "MAIN.powerFactorPercent"
        transform: "value / 100"

      # Voltage scaling
      voltageKV:
        address: "MAIN.voltageV"
        transform: "value / 1000"
```

### Error Handling

If a transformation fails (e.g., division by zero, invalid expression), the batch's `onError()` method will be called:

```java
batch.setListener(new TagBatch.TagBatchListener() {
    @Override
    public void onTagsFetched(TagBatch batch, PlcReadResponse response) {
        // Process transformed values
    }

    @Override
    public void onError(TagBatch batch, Throwable error) {
        if (error instanceof TransformException) {
            System.err.println("Transformation error: " + error.getMessage());
            // Handle transformation error
        } else {
            // Handle other errors
        }
    }
});
```

## Advanced Usage

### Listing Batches

You can get a list of all batch IDs to iterate or display them:

```java
EventPump pump = new EventPump();

// Add some batches
pump.addBatch(batch1);
pump.addBatch(batch2);
pump.addBatch(batch3);

// Get all batch IDs
java.util.Set<String> batchIds = pump.getBatchIds();

// Print batch information
System.out.println("Available batches:");
for (String batchId : batchIds) {
    TagBatch batch = pump.getBatch(batchId);
    System.out.println("  - " + batchId +
        " (started: " + batch.isStarted() +
        ", tags: " + batch.getTags().size() + ")");
}

// Or check if a specific batch exists
if (batchIds.contains("sensor-batch")) {
    pump.startBatch("sensor-batch");
}
```

### Multiple Batches, One PLC

```java
EventPump pump = new EventPump();

// Batch 1: Fast polling of critical values
TagBatch criticalBatch = TagBatch.builder()
    .withBatchId("critical")
    .withConnection(connection)
    .addTag("emergency", "MAIN.emergencyStop")
    .addTag("alarm", "MAIN.alarm")
    .withTrigger(new TimerTrigger(1, TimeUnit.SECONDS))  // Every second
    .withListener((b, r) -> handleCritical(r))
    .build();

// Batch 2: Slow polling of non-critical values
TagBatch statusBatch = TagBatch.builder()
    .withBatchId("status")
    .withConnection(connection)
    .addTag("uptime", "MAIN.uptime")
    .addTag("cycleCount", "MAIN.cycleCount")
    .withTrigger(new TimerTrigger(30, TimeUnit.SECONDS))  // Every 30 seconds
    .withListener((b, r) -> handleStatus(r))
    .build();

pump.addBatch(criticalBatch);
pump.addBatch(statusBatch);
pump.startAll();
```

### Multiple PLCs

```java
EventPump pump = new EventPump();

PlcConnection plc1 = connectionManager.getConnection("ads:tcp://192.168.1.1");
PlcConnection plc2 = connectionManager.getConnection("modbus:tcp://192.168.1.2");
PlcConnection plc3 = connectionManager.getConnection("s7://192.168.1.3");

// Batch from PLC 1 (ADS)
TagBatch batch1 = TagBatch.builder()
    .withBatchId("batch1")
    .withConnection(plc1)
    .addTag("temp", "MAIN.temperature")
    .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
    .withListener((b, r) -> handlePlc1(r))
    .build();

// Batch from PLC 2 (Modbus)
TagBatch batch2 = TagBatch.builder()
    .withBatchId("batch2")
    .withConnection(plc2)
    .addTag("pressure", "40001")
    .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
    .withListener((b, r) -> handlePlc2(r))
    .build();

// Batch from PLC 3 (S7)
TagBatch batch3 = TagBatch.builder()
    .withBatchId("batch3")
    .withConnection(plc3)
    .addTag("flow", "%DB1.DBD0:REAL")
    .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
    .withListener((b, r) -> handlePlc3(r))
    .build();

pump.addBatch(batch1);
pump.addBatch(batch2);
pump.addBatch(batch3);
pump.startAll();
```

### Event-Driven Fetching

Use subscription triggers to fetch data only when needed:

```java
// Subscribe to a trigger tag
SubscriptionTrigger trigger = new SubscriptionTrigger(
    connection,
    "productionStarted",
    "MAIN.productionStarted"
);

// When productionStarted changes, fetch the batch
TagBatch batch = TagBatch.builder()
    .withBatchId("production-batch")
    .withConnection(connection)
    .addTag("productCount", "MAIN.productCount")
    .addTag("batchNumber", "MAIN.batchNumber")
    .addTag("startTime", "MAIN.startTime")
    .withTrigger(trigger)
    .withListener((b, response) -> {
        // Process production data only when production starts
        logProductionData(response);
    })
    .build();
```

### Manual Fetching

You can manually trigger a fetch at any time:

```java
TagBatch batch = TagBatch.builder()
    .withBatchId("batch1")
    .withConnection(connection)
    .addTag("value1", "MAIN.value1")
    .addTag("value2", "MAIN.value2")
    .withTrigger(new TimerTrigger(10, TimeUnit.SECONDS))
    .withListener((b, r) -> handleResults(r))
    .build();

pump.addBatch(batch);

// Start the batch (trigger will fire automatically)
pump.startBatch("batch1");

// Manually fetch outside of trigger schedule
batch.fetchTags().thenRun(() -> {
    System.out.println("Manual fetch completed");
});
```

## Error Handling

### Batch-Level Error Handling

```java
batch.setListener(new TagBatch.TagBatchListener() {
    @Override
    public void onTagsFetched(TagBatch batch, PlcReadResponse response) {
        // Check individual tag responses
        for (String tagName : response.getTagNames()) {
            if (response.getResponseCode(tagName) == PlcResponseCode.OK) {
                Object value = response.getObject(tagName);
                // Process value
            } else {
                System.err.println("Error reading tag " + tagName + ": " +
                    response.getResponseCode(tagName));
            }
        }
    }

    @Override
    public void onError(TagBatch batch, Throwable error) {
        System.err.println("Batch error: " + error.getMessage());
        error.printStackTrace();

        // Optional: Implement retry logic
        // Optional: Send alert
        // Optional: Stop batch if errors persist
    }
});
```

### Connection Error Handling

If a connection fails, batches using that connection will receive errors via `onError()`:

```java
batch.setListener(new TagBatch.TagBatchListener() {
    @Override
    public void onTagsFetched(TagBatch batch, PlcReadResponse response) {
        // Process results
    }

    @Override
    public void onError(TagBatch batch, Throwable error) {
        if (error instanceof PlcConnectionException) {
            // Connection lost - attempt reconnection
            handleConnectionLoss(batch);
        } else {
            // Other error - log and continue
            logger.error("Batch error", error);
        }
    }
});
```

## Best Practices

### 1. Use Try-With-Resources

```java
// Good - automatic cleanup
try (EventPump pump = new EventPump()) {
    pump.addBatch(batch1);
    pump.addBatch(batch2);
    pump.startAll();

    // Run for some time
    Thread.sleep(60000);

} // Automatic cleanup on exit
```

### 2. Set Appropriate Intervals

Consider your use case when setting timer intervals:

```java
// For critical monitoring: short intervals
new TimerTrigger(1, TimeUnit.SECONDS);

// For status updates: medium intervals
new TimerTrigger(30, TimeUnit.SECONDS);

// For historical logging: long intervals
new TimerTrigger(5, TimeUnit.MINUTES);
```

### 3. Batch Related Tags Together

Group tags that are logically related and need the same trigger:

```java
// Good - temperature sensors together
Map<String, String> tempTags = new HashMap<>();
tempTags.put("temp1", "MAIN.temperature1");
tempTags.put("temp2", "MAIN.temperature2");
tempTags.put("temp3", "MAIN.temperature3");

// Bad - mixing unrelated tags with different update needs
Map<String, String> mixedTags = new HashMap<>();
mixedTags.put("temp1", "MAIN.temperature1");  // Updates every second
mixedTags.put("uptime", "MAIN.uptime");       // Updates every hour
```

### 4. Reuse Connections

Don't create multiple connections to the same PLC:

```java
// Good - one connection, multiple batches
PlcConnection connection = connectionManager.getConnection("ads:tcp://192.168.1.1");

TagBatch batch1 = TagBatch.builder()
    .withBatchId("batch1")
    .withConnection(connection)
    .addTag("temp", "MAIN.temp")
    .withTrigger(new TimerTrigger(5, TimeUnit.SECONDS))
    .withListener((b, r) -> handleBatch1(r))
    .build();

TagBatch batch2 = TagBatch.builder()
    .withBatchId("batch2")
    .withConnection(connection)  // Reuse same connection
    .addTag("pressure", "MAIN.pressure")
    .withTrigger(new TimerTrigger(10, TimeUnit.SECONDS))
    .withListener((b, r) -> handleBatch2(r))
    .build();

// Bad - multiple connections to same PLC
PlcConnection conn1 = connectionManager.getConnection("ads:tcp://192.168.1.1");
PlcConnection conn2 = connectionManager.getConnection("ads:tcp://192.168.1.1");
```

### 5. Monitor Performance

```java
batch.setListener((batch, response) -> {
    long fetchTime = System.currentTimeMillis();
    // Process response

    // Log slow fetches
    long duration = System.currentTimeMillis() - fetchTime;
    if (duration > 1000) {
        logger.warn("Slow fetch in batch {}: {}ms", batch.getBatchId(), duration);
    }
});
```

## Limitations

1. **Subscription Support**: Not all PLC protocols support subscriptions. Check protocol capabilities before using `SubscriptionTrigger`.

2. **Expression Conditions**: Advanced expression-based conditions for subscription triggers are not yet implemented. Use predicate functions instead.

3. **SimpleExpressionEvaluator Limitations**: The default expression evaluator does not support:
   - Ternary operator (`condition ? true : false`)
   - Function calls (e.g., `Math.sqrt(value)`)
   - String operations
   - Cross-tag references in expressions (access to other tag values)

   For advanced expression needs, register a custom transformer like SpEL (see Value Transformations section).

4. **No Built-in Retry Logic**: If a fetch fails, the batch continues on the next trigger. Implement custom retry logic in `onError()` if needed.

5. **No Data Buffering**: Results are delivered immediately to listeners. If you need buffering, implement it in your listener.

6. **No Automatic Reconnection**: If a connection fails, batches will error until the connection is restored. Implement reconnection logic in your application if needed.

## Testing

Run the tests:

```bash
mvn test
```

The module includes comprehensive tests for:
- TimerTrigger functionality
- EventPump management
- Configuration serialization/deserialization
- Value transformations (SimpleExpressionEvaluator)
- Transformation registry and pluggable transformers
- TransformedPlcReadResponse wrapper
- Extended tag configuration format

## Troubleshooting

### Batch Not Firing

**Symptom**: Batch never fetches data

**Possible Causes**:
1. Batch not started: Call `pump.startBatch(id)` or `pump.startAll()`
2. No listener set: Call `batch.setListener()` before starting
3. Trigger not configured correctly

**Solution**:
```java
// Verify batch is started
System.out.println("Started: " + batch.isStarted());

// Verify trigger is running
System.out.println("Trigger running: " + batch.getTrigger().isRunning());

// Test manual fetch
batch.fetchTags().get(); // Throws exception if connection/tags are wrong
```

### Subscription Trigger Not Working

**Symptom**: SubscriptionTrigger never fires

**Possible Causes**:
1. Protocol doesn't support subscriptions
2. Tag address is incorrect
3. Tag value not changing

**Solution**:
```java
// Test if protocol supports subscriptions
try {
    PlcSubscriptionRequest request = connection.subscriptionRequestBuilder()
        .addChangeOfStateTagAddress("test", "MAIN.test")
        .build();
    // If this throws UnsupportedOperationException, subscriptions are not supported
} catch (UnsupportedOperationException e) {
    System.err.println("Protocol does not support subscriptions");
}
```

### High CPU Usage

**Symptom**: CPU usage is high

**Possible Causes**:
1. Too many batches
2. Timer intervals are too short
3. Too many tags per batch

**Solution**:
```java
// Increase timer intervals
new TimerTrigger(5, TimeUnit.SECONDS);  // Instead of 100ms

// Reduce number of batches
// Combine related batches

// Monitor batch count
System.out.println("Active batches: " + pump.getStartedBatchCount());
```

## See Also

- [Connection Cache](../connection-cache/README.md) - For managing PLC connection pooling
- [Apache PLC4X Documentation](https://plc4x.apache.org/) - For PLC protocol details
