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
package org.apache.plc4x.kafka.source;

import org.apache.kafka.common.config.ConfigException;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.errors.ConnectException;
import org.apache.kafka.connect.source.SourceRecord;
import org.apache.kafka.connect.source.SourceTask;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.atomic.AtomicBoolean;

public class Plc4xSourceTask extends SourceTask {

    static final Logger log = LoggerFactory.getLogger(Plc4xSourceTask.class);

    private Plc4xSourceConfig config;
    private PlcConnection plcConnection;
    private PlcReader reader;
    private PlcReadRequest readRequest;
    private AtomicBoolean running = new AtomicBoolean(false);
    private String topic;
    private Schema keySchema = Schema.STRING_SCHEMA;
    private Schema valueSchema;
    private long offset = 0;
    private final Map<Class<?>, Schema> typeSchemas = initTypeSchemas();

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public void start(Map<String, String> properties) {
        try {
            config = new Plc4xSourceConfig(properties);
        } catch (ConfigException e) {
            throw new ConnectException("Couldn't start Plc4xSourceTask due to configuration error", e);
        }
        final String url = config.getString(Plc4xSourceConfig.PLC_CONNECTION_STRING_CONFIG);

        try {
            plcConnection = new PlcDriverManager().getConnection(url);
            Optional<PlcReader> readerOptional = plcConnection.getReader();
            if(!readerOptional.isPresent()) {
                throw new ConnectException("PlcReader not available for this type of connection");
            }
            reader = readerOptional.get();
            Class<?> dataType = config.getClass(Plc4xSourceConfig.PLC_DATATYPE_CONFIG);
            String addressString = config.getString(Plc4xSourceConfig.PLC_ADDRESS);
            readRequest = reader.readRequestBuilder().addItem("value", addressString).build();
            topic = config.getString(Plc4xSourceConfig.PLC_TOPIC);
            valueSchema = typeSchemas.get(dataType);
            running.set(true);
        } catch (PlcConnectionException e) {
            throw new ConnectException("Caught exception while connecting to PLC", e);
        }
    }

    private Map<Class<?>, Schema> initTypeSchemas() {
        Map<Class<?>, Schema> map = new HashMap<>();
        map.put(Boolean.class, Schema.BOOLEAN_SCHEMA);
        map.put(Integer.class, Schema.INT32_SCHEMA);
        // TODO add other
        return map;
    }

    @Override
    public void stop() {
        if(plcConnection != null) {
            running.set(false);
            try {
                plcConnection.close();
            } catch (Exception e) {
                throw new RuntimeException("Caught exception while closing connection to PLC", e);
            }
        }
    }

    @Override
    public List<SourceRecord> poll() throws InterruptedException {
        if((plcConnection != null) && plcConnection.isConnected() && (reader != null)) {
            final List<SourceRecord> results = new LinkedList<>();

            try {
                PlcReadResponse<?> plcReadResponse = reader.read(readRequest).get();
                for (String fieldName : plcReadResponse.getFieldNames()) {
                    for (int i = 0; i < plcReadResponse.getNumberOfValues(fieldName); i++) {
                        Object value = plcReadResponse.getObject(fieldName, i);
                        Map<String, String> sourcePartition = Collections.singletonMap("field-name", fieldName);
                        Map<String, Long> sourceOffset = Collections.singletonMap("offset", offset);
                        SourceRecord record = new SourceRecord(sourcePartition, sourceOffset, topic, keySchema, fieldName, valueSchema, value);
                        results.add(record);
                        offset++; // TODO: figure out how to track offsets
                    }
                }
            } catch (ExecutionException e) {
                log.error("Error reading values from PLC", e);
            }

            return results;
        }
        return null;
    }

}