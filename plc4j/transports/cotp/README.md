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
# COTP Transport

Connection Oriented Transport Protocol (COTP) transport layer with TPKT framing over TCP.

## Overview

This transport implements **ISO 8073 (COTP)** with **TPKT framing (RFC 1006)** on top of TCP. It provides connection-oriented communication commonly used in industrial protocols like S7.

## Features

- **TPKT Framing**: Automatic 4-byte header handling (version, reserved, length)
- **COTP Connection Management**: Handles CR (Connection Request) and CC (Connection Confirm)
- **COTP Data Transfer**: Wraps/unwraps data in COTP Data (DT) packets  
- **COTP Disconnection**: Proper shutdown with DR (Disconnect Request)
- **Payload Extraction**: Transparently passes only the payload to drivers

## Usage

### Connection String

```
cotp://hostname:port
```

Example:
```
cotp://192.168.1.100:102
```

### Configuration

```java
CotpTransportConfiguration config = new CotpTransportConfiguration();
config.localTsap = 0x0100;      // Local TSAP
config.remoteTsap = 0x0102;     // Remote TSAP
config.cotpTpduSize = 1024;     // PDU size in bytes
config.cotpConnectionTimeout = 5000; // Connection timeout in ms
config.protocolClass = 0;       // COTP protocol class

// Configure underlying TCP settings
config.setDefaultPort(102);
config.tcpNoDelay = true;
config.connectTimeout = 5000;
```

### Driver Usage

Drivers using COTP transport receive and send pure payload data - all TPKT/COTP framing is handled automatically:

```java
// Read data - TPKT/COTP headers are automatically removed
byte[] payload = transportInstance.read(numBytes);

// Write data - automatically wrapped in COTP Data packet + TPKT
transportInstance.write(payloadBytes);
```

## Protocol Details

### TPKT Format

```
+--------+--------+--------+--------+
| 0x03   | 0x00   | Length (MSB,LSB)|  
+--------+--------+--------+--------+
| COTP Packet ...                   |
+-----------------------------------+
```

### COTP Connection Setup

**Connection Request (CR):**
- TPDU Code: 0xE0
- Contains: Destination/Source references, Protocol class, Parameters (TSAPs, TPDU size)

**Connection Confirm (CC):**
- TPDU Code: 0xD0  
- Response from server confirming connection

### COTP Data Transfer

**Data Packet (DT):**
- TPDU Code: 0xF0
- Contains: EOT flag, TPDU number, Payload

The transport automatically:
- Wraps outgoing payload in DT packets with TPKT framing
- Unwraps incoming TPKT/DT packets and returns only the payload

### COTP Disconnection

**Disconnect Request (DR):**
- TPDU Code: 0x80
- Sent when closing the connection

## TSAP Configuration

TSAP (Transport Service Access Point) identifies communication endpoints:

- **Local TSAP**: Your application's identifier (default: 0x0100)
- **Remote TSAP**: Target PLC/device identifier (default: 0x0102)

Common TSAP values:
- **S7-300**: 0x0102 (for PG communication)
- **S7-400**: 0x0102 (for PG communication)  
- **S7-1200/1500**: 0x0100 (for PG), 0x0102 (for OP)

## TPDU Size

Valid sizes: 128, 256, 512, 1024, 2048, 4096, 8192 bytes

Default: 1024 bytes

## Protocol Class

- **Class 0**: Simple class, no flow control (most common)
- **Class 1-4**: Advanced features

Default: Class 0

## Implementation Notes

1. **Uses Generated Protocol Classes**: All COTP packet handling uses code generated from the mspec protocol definition
2. **ReadBuffer/WriteBuffer Pattern**: Follows the same pattern as drivers for consistent serialization
3. **TCP Transport**: Wraps TCP transport and adds COTP/TPKT framing on top
4. **RingBuffer for Payload Management**: Uses RingBuffer to efficiently store extracted payloads, avoiding repeated packet parsing
5. **Thread-Safe**: Uses underlying TCP transport with proper locking
6. **Connection Timeout**: Configurable timeout for COTP handshake
7. **Error Handling**: Detects COTP errors and connection failures
8. **Automatic Framing**: Drivers work with pure payload, never see TPKT/COTP headers

## Architecture

```
Driver API
    ↓
CotpTransportInstance
    ├─ RingBuffer (payload storage)
    ├─ Uses generated classes: TPKTPacket, COTPPacketData, etc.
    ├─ ReadBufferByteBased / WriteBufferByteBased for parsing/serialization
    └─ TcpTransportInstance (underlying connection)
```

The transport leverages the same code generation infrastructure that drivers use,
ensuring consistent and correct protocol handling.

### Efficient Payload Handling

The transport uses a RingBuffer to efficiently manage payload data:

1. **fillBuffer()**: Reads TPKT packets from TCP, parses them once using generated classes, extracts payload, and stores in RingBuffer
2. **getNumBytesAvailable()**: Calls fillBuffer() and returns buffered payload size
3. **read()/peek()**: Operate directly on buffered payload data without re-parsing

This approach is much more efficient than parsing the same packet multiple times.

## Examples

### S7 PLC Connection

```java
CotpTransportConfiguration config = new CotpTransportConfiguration();
config.setDefaultPort(102);
config.remoteTsap = 0x0102;  // S7 PG communication
config.localTsap = 0x0100;

Transport<CotpTransportConfiguration> transport = new CotpTransport();
TransportInstance instance = transport.createTransportInstance(
    "192.168.1.100:102", config);

// Use the transport - COTP connection is automatically established
```

### Custom TSAP Configuration

```java
// For specific rack/slot addressing in S7
int rack = 0;
int slot = 2;
config.remoteTsap = 0x0100 | (rack << 8) | slot;
```

## Debugging

Enable debug logging to see COTP protocol details:

```xml
<logger name="org.apache.plc4x.java.transport.cotp" level="DEBUG"/>
```

Log output includes:
- Connection handshake details
- TPKT packet sizes
- COTP packet types
- Payload sizes
- Connection lifecycle events

## Related Protocols

- **TCP Transport**: Underlying transport layer
- **S7**: Uses COTP for S7comm protocol  
- **ISO-on-TCP**: Generic term for TPKT/COTP over TCP

## References

- RFC 1006: ISO Transport Service on top of the TCP
- ISO 8073: Connection Oriented Transport Protocol
- [COTP Protocol Specification](/protocols/cotp/src/main/resources/protocols/cotp/cotp.mspec)
