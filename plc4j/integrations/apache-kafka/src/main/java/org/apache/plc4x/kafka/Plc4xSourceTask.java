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
package org.apache.plc4x.kafka;

import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.connect.data.*;
import org.apache.kafka.connect.data.Date;
import org.apache.kafka.connect.errors.ConnectException;
import org.apache.kafka.connect.source.SourceRecord;
import org.apache.kafka.connect.source.SourceTask;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.scraper.config.triggeredscraper.JobConfigurationTriggeredImplBuilder;
import org.apache.plc4x.java.scraper.config.triggeredscraper.ScraperConfigurationTriggeredImpl;
import org.apache.plc4x.java.scraper.config.triggeredscraper.ScraperConfigurationTriggeredImplBuilder;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollector;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollectorImpl;
import org.apache.plc4x.java.utils.connectionpool2.CachedDriverManager;
import org.apache.plc4x.java.utils.connectionpool2.PooledDriverManager;
import org.apache.plc4x.kafka.config.Constants;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.security.SecureRandom;
import java.util.stream.Collectors;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.TimeUnit;

/**
 * Source Connector Task polling the data source at a given rate.
 * A timer thread is scheduled which sets the fetch flag to true every rate milliseconds.
 * When poll() is invoked, the calling thread waits until the fetch flag is set for WAIT_LIMIT_MILLIS.
 * If the flag does not become true, the method returns null, otherwise a fetch is performed.
 */
public class Plc4xSourceTask extends SourceTask {

    private static final Logger log = LoggerFactory.getLogger(Plc4xSourceTask.class);

    private static final ConfigDef CONFIG_DEF = new ConfigDef()
        .define(Constants.CONNECTION_NAME_CONFIG,
            ConfigDef.Type.STRING,
            ConfigDef.Importance.HIGH,
            Constants.CONNECTION_NAME_STRING_DOC)
        .define(Constants.CONNECTION_STRING_CONFIG,
            ConfigDef.Type.STRING,
            ConfigDef.Importance.HIGH,
            Constants.CONNECTION_STRING_DOC)
        .define(Constants.KAFKA_POLL_RETURN_CONFIG,
            ConfigDef.Type.INT,
            Constants.KAFKA_POLL_RETURN_DEFAULT,
            ConfigDef.Importance.HIGH,
            Constants.KAFKA_POLL_RETURN_DOC)
        .define(Constants.BUFFER_SIZE_CONFIG,
            ConfigDef.Type.INT,
            Constants.BUFFER_SIZE_DEFAULT,
            ConfigDef.Importance.HIGH,
            Constants.BUFFER_SIZE_DOC)
        .define(Constants.QUERIES_CONFIG,
            ConfigDef.Type.LIST,
            ConfigDef.Importance.HIGH,
            Constants.QUERIES_DOC);


    private static final Schema KEY_SCHEMA =
        new SchemaBuilder(Schema.Type.STRUCT)
            .field(Constants.SOURCE_NAME_FIELD, Schema.STRING_SCHEMA)
            .field(Constants.JOB_NAME_FIELD, Schema.STRING_SCHEMA)
            .build();

    // Internal buffer into which all incoming scraper responses are written to.
    private ArrayBlockingQueue<SourceRecord> buffer;
    private Integer pollReturnInterval;
    private TriggeredScraperImpl scraper;
    private final SecureRandom random = new SecureRandom();

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public void start(Map<String, String> props) {
        AbstractConfig config = new AbstractConfig(CONFIG_DEF, props);
        String connectionName = config.getString(Constants.CONNECTION_NAME_CONFIG);
        String plc4xConnectionString = config.getString(Constants.CONNECTION_STRING_CONFIG);
        pollReturnInterval = config.getInt(Constants.KAFKA_POLL_RETURN_CONFIG);
        Integer bufferSize = config.getInt(Constants.BUFFER_SIZE_CONFIG);

        Map<String, String> topics = new HashMap<>();
        // Create a buffer with a capacity of BUFFER_SIZE_CONFIG elements which schedules access in a fair way.
        buffer = new ArrayBlockingQueue<>(bufferSize, true);

        ScraperConfigurationTriggeredImplBuilder builder = new ScraperConfigurationTriggeredImplBuilder();
        builder.addSource(connectionName, plc4xConnectionString);

        List<String> jobConfigs = config.getList(Constants.QUERIES_CONFIG);
        for (String jobConfig : jobConfigs) {
            String[] jobConfigSegments = jobConfig.split("\\|");
            if (jobConfigSegments.length < 4) {
                log.warn("Error in job configuration '{}'. " +
                    "The configuration expects at least 4 segments: " +
                    "{job-name}|{topic}|{rate}(|{field-alias}#{field-address})+", jobConfig);
                continue;
            }

            String jobName = jobConfigSegments[0];
            String topic = jobConfigSegments[1];
            Integer rate = Integer.valueOf(jobConfigSegments[2]);
            JobConfigurationTriggeredImplBuilder jobBuilder = builder.job(
                jobName, String.format("(SCHEDULED,%s)", rate)).source(connectionName);
            for (int i = 3; i < jobConfigSegments.length; i++) {
                String[] fieldSegments = jobConfigSegments[i].split("#");
                if (fieldSegments.length != 2) {
                    log.warn("Error in job configuration '{}'. " +
                            "The field segment expects a format {field-alias}#{field-address}, but got '%s'",
                        jobName, jobConfigSegments[i]);
                    continue;
                }
                String fieldAlias = fieldSegments[0];
                String fieldAddress = fieldSegments[1];
                jobBuilder.field(fieldAlias, fieldAddress);
                topics.put(jobName, topic);
            }
            jobBuilder.build();
        }

        ScraperConfigurationTriggeredImpl scraperConfig = builder.build();

        try {
            PlcDriverManager manager = new PlcDriverManager();
            PlcDriverManager plcDriverManager = new CachedDriverManager(plc4xConnectionString,() -> manager.getConnection(plc4xConnectionString));
            TriggerCollector triggerCollector = new TriggerCollectorImpl(plcDriverManager);
            scraper = new TriggeredScraperImpl(scraperConfig, (jobName, sourceName, results) -> {
                try {
                    Long timestamp = System.currentTimeMillis();

                    Map<String, String> sourcePartition = new HashMap<>();
                    sourcePartition.put("sourceName", sourceName);
                    sourcePartition.put("jobName", jobName);

                    Map<String, Long> sourceOffset = Collections.singletonMap("offset", timestamp);

                    String topic = topics.get(jobName);

                    // Prepare the key structure.
                    Struct key = new Struct(KEY_SCHEMA)
                        .put(Constants.SOURCE_NAME_FIELD, sourceName)
                        .put(Constants.JOB_NAME_FIELD, jobName);

                    // Build the Schema for the result struct.
                    SchemaBuilder fieldSchemaBuilder = SchemaBuilder.struct()
                        .name("org.apache.plc4x.kafka.schema.Field");


                    for (Map.Entry<String, Object> result : results.entrySet()) {
                        // Get field-name and -value from the results.
                        String fieldName = result.getKey();
                        Object fieldValue = result.getValue();

                        // Get the schema for the given value type.
                        Schema valueSchema = getSchema(fieldValue);

                        // Add the schema description for the current field.
                        fieldSchemaBuilder.field(fieldName, valueSchema);
                    }
                    Schema fieldSchema = fieldSchemaBuilder.build();

                    Schema recordSchema = SchemaBuilder.struct()
                        .name("org.apache.plc4x.kafka.schema.JobResult")
                        .doc("PLC Job result. This contains all of the received PLCValues as well as a recieved timestamp")
                        .field(Constants.FIELDS_CONFIG, fieldSchema)
                        .field(Constants.TIMESTAMP_CONFIG, Schema.INT64_SCHEMA)
                        .field(Constants.EXPIRES_CONFIG, Schema.OPTIONAL_INT64_SCHEMA)
                        .build();

                    // Build the struct itself.
                    Struct fieldStruct = new Struct(fieldSchema);
                    for (Map.Entry<String, Object> result : results.entrySet()) {
                        // Get field-name and -value from the results.
                        String fieldName = result.getKey();
                        Object fieldValue = result.getValue();

                        if (fieldSchema.field(fieldName).schema().type() == Schema.Type.ARRAY) {
                            fieldStruct.put(fieldName, ((List) fieldValue).stream().map(p -> ((PlcValue) p).getObject()).collect(Collectors.toList()));
                        } else {
                            fieldStruct.put(fieldName, fieldValue);
                        }
                    }

                    Struct recordStruct = new Struct(recordSchema)
                        .put(Constants.FIELDS_CONFIG, fieldStruct)
                        .put(Constants.TIMESTAMP_CONFIG, timestamp);

                    // Prepare the source-record element.
                    SourceRecord sourceRecord = new SourceRecord(
                        sourcePartition, sourceOffset,
                        topic,
                        KEY_SCHEMA, key,
                        recordSchema, recordStruct
                    );

                    // Add the new source-record to the buffer.
                    buffer.add(sourceRecord);
                } catch (Exception e) {
                    log.error("Error while parsing returned values", e);
                }
            }, triggerCollector);
            scraper.start();
            triggerCollector.start();
        } catch (ScraperException e) {
            log.error("Error starting the scraper", e);
        }
    }

    @Override
    public void stop() {
        synchronized (this) {
            scraper.stop();
            notifyAll(); // wake up thread waiting in awaitFetch
        }
    }

    @Override
    public List<SourceRecord> poll() {
        if (!buffer.isEmpty()) {
            int numElements = buffer.size();
            List<SourceRecord> result = new ArrayList<>(numElements);
            buffer.drainTo(result, numElements);
            return result;
        }
        try {
            List<SourceRecord> result = new ArrayList<>(1);
            SourceRecord temp = buffer.poll(pollReturnInterval + (long) random.nextInt((int) Math.round(pollReturnInterval * 0.05)), TimeUnit.MILLISECONDS);
            if (temp == null) {
                return Collections.emptyList();
            }
            result.add(temp);
            return result;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            return Collections.emptyList();
        }
    }

    private Schema getSchema(Object value) {
        Objects.requireNonNull(value);

        if (value instanceof PlcValue) {
            value = ((PlcValue) value).getObject();
        }

        if (value instanceof List) {
            List list = (List) value;
            if (list.isEmpty()) {
                throw new ConnectException("Unsupported empty lists.");
            }
            // In PLC4X list elements all contain the same type.
            Object firstElement = list.get(0);
            Schema elementSchema = getSchema(firstElement);
            return SchemaBuilder.array(elementSchema).build();
        }
        if (value instanceof BigInteger) {
            // no support yet
        }
        if (value instanceof BigDecimal) {
            // no support yet
        }
        if (value instanceof Boolean) {
            return Schema.OPTIONAL_BOOLEAN_SCHEMA;
        }
        if (value instanceof byte[]) {
            return Schema.OPTIONAL_BYTES_SCHEMA;
        }
        if (value instanceof Byte) {
            return Schema.OPTIONAL_INT8_SCHEMA;
        }
        if (value instanceof Double) {
            return Schema.OPTIONAL_FLOAT64_SCHEMA;
        }
        if (value instanceof Float) {
            return Schema.OPTIONAL_FLOAT32_SCHEMA;
        }
        if (value instanceof Integer) {
            return Schema.OPTIONAL_INT32_SCHEMA;
        }
        if (value instanceof LocalDate) {
            return Date.builder().optional().build();
        }
        if (value instanceof LocalDateTime) {
            return Timestamp.builder().optional().build();
        }
        if (value instanceof LocalTime) {
            return Time.builder().optional().build();
        }
        if (value instanceof Long) {
            return Schema.OPTIONAL_INT64_SCHEMA;
        }
        if (value instanceof Short) {
            return Schema.OPTIONAL_INT16_SCHEMA;
        }
        if (value instanceof String) {
            return Schema.OPTIONAL_STRING_SCHEMA;
        }
        // TODO: add support for collective and complex types
        throw new ConnectException(String.format("Unsupported data type %s", value.getClass().getName()));
    }

}
