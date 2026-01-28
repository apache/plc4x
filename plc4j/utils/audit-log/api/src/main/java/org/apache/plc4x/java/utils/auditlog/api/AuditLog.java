/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.utils.auditlog.api;

import org.apache.plc4x.java.utils.auditlog.api.config.AuditLogConfiguration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Main API for the audit log system.
 * This class provides a simple interface for logging events to an audit log.
 * The actual implementation is only active if the audit-log-impl module is on the classpath.
 * <p>
 * Example using the Builder pattern:
 * <pre>
 * AuditLog auditLog = AuditLog.builder()
 *     .withSource("connection-1")
 *     .withDriverTestsuiteFile("/path/to/testsuite.xml")
 *     .build();
 *
 * auditLog.write(AuditLogEventType.CONNECT, "Connected to PLC");
 * </pre>
 */
public abstract class AuditLog {

    private static final Logger LOGGER = LoggerFactory.getLogger(AuditLog.class);
    private static final String IMPL_CLASS_NAME = "org.apache.plc4x.java.utils.auditlog.impl.AuditLogImpl";

    protected final AuditLogConfiguration config;
    protected final String source;

    /**
     * Protected constructor to be called by implementations
     * @param config the audit log configuration
     * @param source the source identifier for log entries
     */
    protected AuditLog(AuditLogConfiguration config, String source) {
        this.config = config;
        this.source = source != null ? source : "unknown";
    }

    /**
     * Create a new Builder for constructing AuditLog instances.
     *
     * @return A new Builder
     */
    public static Builder builder() {
        return new Builder();
    }

    /**
     * Checks if the audit log is enabled.
     *
     * @return true if the audit log is enabled and the event should be logged, false otherwise
     */
    public boolean isEnabled() {
        return config.isEnabled();
    }

    /**
     * Writes an event to the audit log.
     *
     * @param eventType the type of event
     * @param message the event message
     */
    public abstract void write(AuditLogEventType eventType, String message);

    /**
     * Writes an event to the audit log with an object that will be serialized to JSON.
     * The default implementation converts the object to a string using toString().
     * Implementations can override this to provide JSON serialization.
     *
     * @param eventType the type of event
     * @param message the event message
     * @param data the data object to serialize (will be converted to JSON by implementations)
     */
    public void write(AuditLogEventType eventType, String message, Object data) {
        // Default implementation: just append toString() representation
        String dataString = (data != null) ? data.toString() : "null";
        write(eventType, message + ": " + dataString);
    }

    /**
     * Closes the audit log and releases any resources.
     */
    public abstract void close();

    /**
     * Gets the configuration for this audit log instance
     * @return the audit log configuration
     */
    public AuditLogConfiguration getConfig() {
        return config;
    }

    /**
     * Gets the source identifier for this audit log instance
     * @return the source identifier
     */
    public String getSource() {
        return source;
    }

    /**
     * Builder for creating AuditLog instances using a fluent API.
     * <p>
     * Example:
     * <pre>
     * AuditLog auditLog = AuditLog.builder()
     *     .withSource("my-connection")
     *     .withDriverTestsuiteFile("/path/to/testsuite.xml")
     *     .build();
     * </pre>
     */
    public static class Builder {
        private String auditLogFile;
        private String source;

        /**
         * Set the source identifier for all log entries.
         * This identifies where the log entries come from (e.g., connection ID, driver instance).
         *
         * @param source the source identifier
         * @return this builder
         */
        public Builder withSource(String source) {
            this.source = source;
            return this;
        }

        /**
         * Set the path to the audit-log file.
         * If not set or set to null/empty, audit logging will be disabled.
         *
         * @param auditLogFile the path to the audit-log file
         * @return this builder
         */
        public Builder withAuditLogFile(String auditLogFile) {
            this.auditLogFile = auditLogFile;
            return this;
        }

        /**
         * Set the configuration directly.
         * This allows using an existing AuditLogConfiguration object.
         *
         * @param config the audit log configuration
         * @return this builder
         */
        public Builder withConfiguration(AuditLogConfiguration config) {
            if (config != null) {
                this.auditLogFile = config.auditLogFile;
            }
            return this;
        }

        /**
         * Build the AuditLog instance.
         * If the impl module is not available on the classpath, returns a no-op implementation.
         *
         * @return a new AuditLog instance
         */
        public AuditLog build() {
            AuditLogConfiguration config = new AuditLogConfiguration();
            config.auditLogFile = this.auditLogFile;

            if (!config.isEnabled()) {
                return new NoOpAuditLog(config, this.source);
            }

            try {
                // Try to load the implementation class
                Class<?> implClass = AuditLog.class.getClassLoader().loadClass(IMPL_CLASS_NAME);
                LOGGER.info("Audit log implementation found on classpath. Audit logging enabled.");
                return (AuditLog) implClass.getDeclaredConstructor(AuditLogConfiguration.class, String.class)
                    .newInstance(config, this.source);
            } catch (ClassNotFoundException e) {
                LOGGER.info("Audit log implementation not found on classpath. Audit logging will be disabled.");
                return new NoOpAuditLog(config, this.source);
            } catch (Exception e) {
                LOGGER.warn("Failed to instantiate audit log implementation. Audit logging will be disabled.", e);
                return new NoOpAuditLog(config, this.source);
            }
        }
    }

    /**
     * No-op implementation used when audit logging is disabled or the impl is not available
     */
    private static class NoOpAuditLog extends AuditLog {

        protected NoOpAuditLog(AuditLogConfiguration config, String source) {
            super(config, source);
        }

        @Override
        public void write(AuditLogEventType eventType, String message) {
            // No-op
        }

        @Override
        public void close() {
            // No-op
        }
    }
}
