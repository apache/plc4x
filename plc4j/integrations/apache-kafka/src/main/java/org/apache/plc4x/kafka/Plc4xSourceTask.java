/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.kafka;

import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.connect.data.*;
import org.apache.kafka.connect.data.Date;
import org.apache.kafka.connect.errors.ConnectException;
import org.apache.kafka.connect.source.SourceRecord;
import org.apache.kafka.connect.source.SourceTask;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.scraper.config.triggeredscraper.JobConfigurationTriggeredImplBuilder;
import org.apache.plc4x.java.scraper.config.triggeredscraper.ScraperConfigurationTriggeredImpl;
import org.apache.plc4x.java.scraper.config.triggeredscraper.ScraperConfigurationTriggeredImplBuilder;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollector;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollectorImpl;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigDecimal;
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

    /*
     * Config of the task.
     */
    static final String CONNECTION_NAME_CONFIG = "connection-name";
    private static final String CONNECTION_NAME_STRING_DOC = "Connection Name";

    static final String PLC4X_CONNECTION_STRING_CONFIG = "plc4x-connection-string";
    private static final String PLC4X_CONNECTION_STRING_DOC = "PLC4X Connection String";

    // Syntax for the queries: {job-name}:{topic}:{rate}:{field-alias}#{field-address}:{field-alias}#{field-address}...,{topic}:{rate}:....
    static final String QUERIES_CONFIG = "queries";
    private static final String QUERIES_DOC = "Field queries to be sent to the PLC";

    private static final ConfigDef CONFIG_DEF = new ConfigDef()
        .define(CONNECTION_NAME_CONFIG, ConfigDef.Type.STRING, ConfigDef.Importance.HIGH, CONNECTION_NAME_STRING_DOC)
        .define(PLC4X_CONNECTION_STRING_CONFIG, ConfigDef.Type.STRING, ConfigDef.Importance.HIGH, PLC4X_CONNECTION_STRING_DOC)
        .define(QUERIES_CONFIG, ConfigDef.Type.LIST, ConfigDef.Importance.HIGH, QUERIES_DOC);

    /*
     * Configuration of the output.
     */
    private static final String SOURCE_NAME_FIELD = "source-name";
    private static final String JOB_NAME_FIELD = "job-name";

    private static final Schema KEY_SCHEMA =
        new SchemaBuilder(Schema.Type.STRUCT)
            .field(SOURCE_NAME_FIELD, Schema.STRING_SCHEMA)
            .field(JOB_NAME_FIELD, Schema.STRING_SCHEMA)
            .build();

    // Internal buffer into which all incoming scraper responses are written to.
    private ArrayBlockingQueue<SourceRecord> buffer;

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public void start(Map<String, String> props) {
        AbstractConfig config = new AbstractConfig(CONFIG_DEF, props);
        String connectionName = config.getString(CONNECTION_NAME_CONFIG);
        String plc4xConnectionString = config.getString(PLC4X_CONNECTION_STRING_CONFIG);
        Map<String, String> topics = new HashMap<>();
        // Create a buffer with a capacity of 1000 elements which schedules access in a fair way.
        buffer = new ArrayBlockingQueue<>(1000, true);

        ScraperConfigurationTriggeredImplBuilder builder = new ScraperConfigurationTriggeredImplBuilder();
        builder.addSource(connectionName, plc4xConnectionString);

        List<String> jobConfigs = config.getList(QUERIES_CONFIG);
        for (String jobConfig : jobConfigs) {
            String[] jobConfigSegments = jobConfig.split("\\|");
            if(jobConfigSegments.length < 4) {
                log.warn(String.format("Error in job configuration '%s'. " +
                    "The configuration expects at least 4 segments: " +
                    "{job-name}|{topic}|{rate}(|{field-alias}#{field-address})+", jobConfig));
                continue;
            }

            String jobName = jobConfigSegments[0];
            String topic = jobConfigSegments[1];
            Integer rate = Integer.valueOf(jobConfigSegments[2]);
            JobConfigurationTriggeredImplBuilder jobBuilder = builder.job(
                jobName, String.format("(SCHEDULED,%s)", rate)).source(connectionName);
            for(int i = 3; i < jobConfigSegments.length; i++) {
                String[] fieldSegments = jobConfigSegments[i].split("#");
                if(fieldSegments.length != 2) {
                    log.warn(String.format("Error in job configuration '%s'. " +
                            "The field segment expects a format {field-alias}#{field-address}, but got '%s'",
                        jobName, jobConfigSegments[i]));
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
            PlcDriverManager plcDriverManager = new PooledPlcDriverManager();
            TriggerCollector triggerCollector = new TriggerCollectorImpl(plcDriverManager);
            TriggeredScraperImpl scraper = new TriggeredScraperImpl(scraperConfig, (jobName, sourceName, results) -> {
                Long timestamp = System.currentTimeMillis();

                Map<String, String> sourcePartition = new HashMap<>();
                sourcePartition.put("sourceName", sourceName);
                sourcePartition.put("jobName", jobName);

                Map<String, Long> sourceOffset = Collections.singletonMap("offset", timestamp);

                String topic = topics.get(jobName);

                // Prepare the key structure.
                Struct key = new Struct(KEY_SCHEMA)
                    .put(SOURCE_NAME_FIELD, sourceName)
                    .put(JOB_NAME_FIELD, jobName);

                // Build the Schema for the result struct.
                SchemaBuilder recordSchemaBuilder = SchemaBuilder.struct().name("org.apache.plc4x.kafka.JobResult");
                for (Map.Entry<String, Object> result : results.entrySet()) {
                    // Get field-name and -value from the results.
                    String fieldName = result.getKey();
                    Object fieldValue = result.getValue();

                    // Get the schema for the given value type.
                    Schema valueSchema = getSchema(fieldValue);

                    // Add the schema description for the current field.
                    recordSchemaBuilder.field(fieldName, valueSchema);
                }
                // Add a timestamp
                Schema valueSchema = Schema.STRING_SCHEMA;
                recordSchemaBuilder.field("timestamp", valueSchema);
                Schema recordSchema = recordSchemaBuilder.build();

                // Build the struct itself.
                Struct recordStruct = new Struct(recordSchema);
                for (Map.Entry<String, Object> result : results.entrySet()) {
                    // Get field-name and -value from the results.
                    String fieldName = result.getKey();
                    Object fieldValue = result.getValue();
                    recordStruct.put(fieldName, fieldValue);
                }
                recordStruct.put("timestamp", LocalDateTime.now().toString());

                // Prepare the source-record element.
                SourceRecord record = new SourceRecord(
                    sourcePartition, sourceOffset,
                    topic,
                    KEY_SCHEMA, key,
                    recordSchema, recordStruct
                    );

                // Add the new source-record to the buffer.
                buffer.add(record);
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
            // TODO: Correctly shutdown the scraper.
            notifyAll(); // wake up thread waiting in awaitFetch
        }
    }

    @Override
    public List<SourceRecord> poll() {
        if(!buffer.isEmpty()) {
            int numElements = buffer.size();
            List<SourceRecord> result = new ArrayList<>(numElements);
            buffer.drainTo(result, numElements);
            return result;
        } else {
            try {
                List<SourceRecord> result = new ArrayList<>(1);
                result.add(buffer.poll(5000, TimeUnit.MILLISECONDS));
                return result;
            } catch (InterruptedException e) {
                return null;
            }
        }
    }

    private Schema getSchema(Object value) {
        Objects.requireNonNull(value);

        if(value instanceof List) {
            List list = (List) value;
            if(list.isEmpty()) {
                throw new ConnectException("Unsupported empty lists.");
            }
            // In PLC4X list elements all contain the same type.
            Object firstElement = list.get(0);
            Schema elementSchema = getSchema(firstElement);
            return SchemaBuilder.array(elementSchema).build();
        }
        if (value instanceof BigDecimal) {

        }
        if (value instanceof Boolean) {
            return Schema.BOOLEAN_SCHEMA;
        }
        if (value instanceof byte[]) {
            return Schema.BYTES_SCHEMA;
        }
        if (value instanceof Byte) {
            return Schema.INT8_SCHEMA;
        }
        if (value instanceof Double) {
            return Schema.FLOAT64_SCHEMA;
        }
        if (value instanceof Float) {
            return Schema.FLOAT32_SCHEMA;
        }
        if (value instanceof Integer) {
            return Schema.INT32_SCHEMA;
        }
        if (value instanceof LocalDate) {
            return Date.SCHEMA;
        }
        if (value instanceof LocalDateTime) {
            return Timestamp.SCHEMA;
        }
        if (value instanceof LocalTime) {
            return Time.SCHEMA;
        }
        if (value instanceof Long) {
            return Schema.INT64_SCHEMA;
        }
        if (value instanceof Short) {
            return Schema.INT16_SCHEMA;
        }
        if (value instanceof String) {
            return Schema.STRING_SCHEMA;
        }
        // TODO: add support for collective and complex types
        throw new ConnectException(String.format("Unsupported data type %s", value.getClass().getName()));
    }

}