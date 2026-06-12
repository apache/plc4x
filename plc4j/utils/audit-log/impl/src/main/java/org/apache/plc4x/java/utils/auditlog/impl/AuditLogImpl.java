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

package org.apache.plc4x.java.utils.auditlog.impl;

import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.Logger;
import ch.qos.logback.classic.LoggerContext;
import ch.qos.logback.classic.encoder.PatternLayoutEncoder;
import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.core.rolling.RollingFileAppender;
import ch.qos.logback.core.rolling.SizeAndTimeBasedRollingPolicy;
import ch.qos.logback.core.util.FileSize;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.apache.plc4x.java.utils.auditlog.api.config.AuditLogConfiguration;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.time.Instant;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;

/**
 * Implementation of the audit log using Logback.
 * This class is only instantiated if the impl module is on the classpath.
 */
public class AuditLogImpl extends AuditLog {

    private static final DateTimeFormatter TIMESTAMP_FORMATTER =
        DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss.SSS")
            .withZone(ZoneId.systemDefault());

    private static final org.slf4j.Logger slf4jLogger = LoggerFactory.getLogger(AuditLog.class);

    private final Logger logger;
    private final RollingFileAppender<ILoggingEvent> fileAppender;
    private final LoggerContext loggerContext;
    private final ObjectMapper objectMapper;

    /**
     * Creates a new AuditLogImpl instance with the given configuration.
     *
     * @param config the audit log configuration
     * @param source the source identifier for log entries
     */
    public AuditLogImpl(AuditLogConfiguration config, String source) {
        super(config, source);

        // Initialize JSON object mapper
        objectMapper = new ObjectMapper();
        objectMapper.configure(SerializationFeature.FAIL_ON_EMPTY_BEANS, false);
        objectMapper.configure(SerializationFeature.INDENT_OUTPUT, false); // Compact JSON for logs

        // Create a unique logger for this audit log instance
        String loggerName = "AuditLog-" + System.identityHashCode(this);

        loggerContext = (LoggerContext) LoggerFactory.getILoggerFactory();
        logger = loggerContext.getLogger(loggerName);

        // Don't inherit appenders from root logger
        logger.setAdditive(false);

        // Ensure parent directory exists
        ensureParentDirectoryExists(config.auditLogFile);

        // Create a rolling file appender with compression
        fileAppender = new RollingFileAppender<>();
        fileAppender.setContext(loggerContext);
        fileAppender.setFile(config.auditLogFile);
        fileAppender.setImmediateFlush(true);

        // Create encoder
        PatternLayoutEncoder encoder = new PatternLayoutEncoder();
        encoder.setContext(loggerContext);
        encoder.setPattern("%msg%n");
        encoder.start();

        // Create a rolling policy with compression (.gz extension triggers GZIP compression)
        SizeAndTimeBasedRollingPolicy<ILoggingEvent> rollingPolicy = new SizeAndTimeBasedRollingPolicy<>();
        rollingPolicy.setContext(loggerContext);
        rollingPolicy.setParent(fileAppender);
        // Pattern: audit-2025-11-09.0.log.gz, audit-2025-11-09.1.log.gz, etc.
        rollingPolicy.setFileNamePattern(config.auditLogFile + ".%d{yyyy-MM-dd}.%i.gz");
        rollingPolicy.setMaxFileSize(FileSize.valueOf("10MB")); // Roll when file reaches 10MB
        rollingPolicy.setMaxHistory(30); // Keep 30 days of history
        rollingPolicy.setTotalSizeCap(FileSize.valueOf("1GB")); // Max total size of all archived files
        rollingPolicy.start();

        fileAppender.setRollingPolicy(rollingPolicy);
        fileAppender.setEncoder(encoder);
        fileAppender.start();

        logger.addAppender(fileAppender);
        logger.setLevel(Level.INFO);

        slf4jLogger.info("Audit log initialized for file: {}", config.auditLogFile);
    }

    @Override
    public void write(AuditLogEventType eventType, String message) {
        // Check if this event type is enabled
        if (!config.isEnabled()) {
            return;
        }

        // Format: [timestamp] [eventType] [source] message
        String timestamp = TIMESTAMP_FORMATTER.format(Instant.now());
        String formattedMessage = String.format("[%s] [%s] [%s] %s",
            timestamp, eventType, source, message);

        logger.info(formattedMessage);
    }

    @Override
    public void write(AuditLogEventType eventType, String message, Object data) {
        // Check if this event type is enabled
        if (!config.isEnabled()) {
            return;
        }

        // Serialize the data object to JSON
        String jsonData;
        try {
            jsonData = objectMapper.writeValueAsString(data);
        } catch (JsonProcessingException e) {
            // Fallback to toString() if JSON serialization fails
            jsonData = (data != null) ? data.toString() : "null";
            slf4jLogger.warn("Failed to serialize object to JSON, using toString(): {}", e.getMessage());
        }

        // Format: [timestamp] [eventType] [source] message: {json}
        String timestamp = TIMESTAMP_FORMATTER.format(Instant.now());
        String formattedMessage = String.format("[%s] [%s] [%s] %s: %s",
            timestamp, eventType, source, message, jsonData);

        logger.info(formattedMessage);
    }

    @Override
    public void close() {
        if (fileAppender != null) {
            fileAppender.stop();
        }
        if (logger != null) {
            logger.detachAndStopAllAppenders();
        }

        slf4jLogger.info("Audit log closed");
    }

    /**
     * Ensures that the parent directory of the given file path exists.
     * Creates the directory if it doesn't exist.
     *
     * @param filePath the file path
     */
    private void ensureParentDirectoryExists(String filePath) {
        if (filePath == null || filePath.isEmpty()) {
            return;
        }

        File file = new File(filePath);
        File parentDir = file.getParentFile();

        if (parentDir != null && !parentDir.exists()) {
            slf4jLogger.info("Creating parent directory: {}", parentDir.getAbsolutePath());
            boolean created = parentDir.mkdirs();
            if (!created) {
                slf4jLogger.error("Failed to create parent directory: {}", parentDir.getAbsolutePath());
            } else {
                slf4jLogger.info("Successfully created parent directory: {}", parentDir.getAbsolutePath());
            }
        }
    }
}
