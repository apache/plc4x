# Audit Log Utility

The audit-log utility provides a lightweight, configurable logging system for recording driver communication events. It is primarily designed to capture communication traces that can be used for debugging and creating driver test suites.

## Module Structure

The audit-log is split into two modules following an API/implementation separation pattern:

- **api** (`plc4j-utils-audit-log-api`) - Contains the public API, configuration classes, and a no-op fallback implementation
- **impl** (`plc4j-utils-audit-log-impl`) - Contains the actual file-based implementation using Logback

This separation allows code to depend only on the API module at compile time while the implementation is optionally included at runtime.

## Features

- **Conditional Activation**: Audit logging is only active when the implementation module is on the classpath AND a log file path is configured
- **No-Op Fallback**: When disabled or unconfigured, uses a zero-overhead no-op implementation
- **Rolling File Appender**: Automatically rolls and compresses log files based on size (10MB) and time (daily)
- **JSON Serialization**: Supports logging objects as JSON using Jackson
- **Source Identification**: Each log instance can be associated with a source identifier (e.g., connection ID)

## Usage

### Basic Usage with Builder

```java
AuditLog auditLog = AuditLog.builder()
    .withSource("connection-1")
    .withAuditLogFile("/path/to/audit-log.txt")
    .build();

// Log events
auditLog.write(AuditLogEventType.CONNECT, "Connected to PLC at 192.168.1.100");
auditLog.write(AuditLogEventType.OUTGOING_BYTES, "Sent request", requestBytes);
auditLog.write(AuditLogEventType.INCOMING_MESSAGE, "Received response", responseObject);

// Close when done
auditLog.close();
```

### Configuration via Connection String

The audit log can be configured via the `log.audit-log-file` connection string parameter:

```
ads://192.168.1.100?log.audit-log-file=/tmp/ads-debug.log
```

### Using AuditLogProvider Interface

Classes that need audit logging can implement the `AuditLogProvider` interface:

```java
public class MyConnection implements AuditLogProvider {
    private final AuditLog auditLog;

    public MyConnection(AuditLogConfiguration config) {
        this.auditLog = AuditLog.builder()
            .withConfiguration(config)
            .withSource("my-connection")
            .build();
    }

    @Override
    public AuditLog getAuditLog() {
        return auditLog;
    }
}
```

## Event Types

The `AuditLogEventType` enum defines the following event types:

| Event Type | Description |
|------------|-------------|
| `CONFIG` | General configuration information |
| `SYSTEM` | General system information |
| `CONNECT` | Connection establishment events |
| `OUTGOING_BYTES` | Raw bytes being sent |
| `OUTGOING_MESSAGE` | Parsed messages being sent |
| `INCOMING_BYTES` | Raw bytes received |
| `INCOMING_MESSAGE` | Parsed messages received |
| `API_REQUEST` | Incoming API requests |
| `API_RESPONSE` | Outgoing API responses |
| `CLOSE` | Connection close events |
| `ERROR` | Error events |

## Log Format

Log entries follow this format:

```
[timestamp] [eventType] [source] message
```

Example:
```
[2025-01-28 14:30:45.123] [CONNECT] [connection-1] Connected to PLC at 192.168.1.100
[2025-01-28 14:30:45.234] [OUTGOING_MESSAGE] [connection-1] Sending read request: {"address":"DB1.DBD0","count":4}
```

## File Rolling Policy

When enabled, the implementation uses a size-and-time-based rolling policy:

- **Maximum file size**: 10MB per file
- **Rolling pattern**: `{filename}.{date}.{index}.gz` (e.g., `audit.txt.2025-01-28.0.gz`)
- **Maximum history**: 30 days
- **Total size cap**: 1GB for all archived files
- **Compression**: GZIP compression for archived files

## Dependencies

### API Module
- `plc4j-spi-config` - Configuration annotations
- `slf4j-api` - Logging facade

### Implementation Module
- `plc4j-utils-audit-log-api` - The API module
- `logback-classic` / `logback-core` - File appender implementation
- `jackson-core` / `jackson-databind` - JSON serialization

## Enabling Audit Logging

To enable audit logging in your application:

1. Add the API dependency (compile scope):
```xml
<dependency>
    <groupId>org.apache.plc4x</groupId>
    <artifactId>plc4j-utils-audit-log-api</artifactId>
    <version>${project.version}</version>
</dependency>
```

2. Add the implementation dependency (runtime scope):
```xml
<dependency>
    <groupId>org.apache.plc4x</groupId>
    <artifactId>plc4j-utils-audit-log-impl</artifactId>
    <version>${project.version}</version>
    <scope>runtime</scope>
</dependency>
```

3. Configure the `audit-log-file` parameter with a valid file path

If the implementation module is not on the classpath or no file path is configured, audit logging silently falls back to a no-op implementation with zero overhead.
