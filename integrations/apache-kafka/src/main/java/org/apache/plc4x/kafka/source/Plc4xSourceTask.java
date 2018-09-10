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
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.kafka.Plc4xSourceConnector;
import org.apache.plc4x.kafka.util.VersionUtil;

import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.*;

public class Plc4xSourceTask extends SourceTask {
    private final static String FIELD_KEY = "key";

    private String topic;
    private String url;
    private String query;
    private Integer rate;

    private volatile boolean running = false;
    private PlcConnection plcConnection;
    private PlcReader plcReader;
    private PlcReadRequest plcRequest;

    private final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);
    private ScheduledFuture<?> timer;
    private boolean fetch = true;

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public void start(Map<String, String> props) {
        topic = props.get(Plc4xSourceConnector.TOPIC_CONFIG);
        url = props.get(Plc4xSourceConnector.URL_CONFIG);
        query = props.get(Plc4xSourceConnector.QUERY_CONFIG);
        rate = Integer.valueOf(props.get(Plc4xSourceConnector.RATE_CONFIG));

        try {
            plcConnection = new PlcDriverManager().getConnection(url);
            plcConnection.connect();
        } catch (PlcConnectionException e) {
            throw new ConnectException("Could not establish a PLC connection", e);
        }

        plcReader = plcConnection.getReader()
            .orElseThrow(() -> new ConnectException("PlcReader not available for this type of connection"));

        plcRequest = plcReader.readRequestBuilder().addItem(FIELD_KEY, query).build();

        timer = scheduler.scheduleAtFixedRate(() -> {
            synchronized (Plc4xSourceTask.this) {
                Plc4xSourceTask.this.fetch = true;
                notify();
            }
        }, 0, rate, TimeUnit.MILLISECONDS);

        running = true;
    }

    @Override
    public void stop() {
        running = false;
        timer.cancel(true);
        if (plcConnection != null) {
            try {
                plcConnection.close();
            } catch (Exception e) {
                throw new RuntimeException("Caught exception while closing connection to PLC", e);
            }
        }
    }

    @Override
    public List<SourceRecord> poll() {
        if (!running)
            return null;

        synchronized (this) {
            while (!fetch) {
                try {
                    wait();
                } catch (InterruptedException e) {
                    // continue
                }
            }
            List<SourceRecord> result = new LinkedList<>();
            try {
                PlcReadResponse<?> response = plcReader.read(plcRequest).get();
                if (response.getResponseCode(FIELD_KEY).equals(PlcResponseCode.OK)) {
                    Object rawValue = response.getObject(FIELD_KEY);
                    Schema valueSchema = getSchema(rawValue.getClass());
                    Object value = valueSchema.equals(Schema.STRING_SCHEMA) ? rawValue.toString() : rawValue;
                    Long timestamp = System.currentTimeMillis();
                    Map<String, String> sourcePartition = Collections.singletonMap("url", url);
                    Map<String, Long> sourceOffset = Collections.singletonMap("offset", timestamp);

                    SourceRecord record =
                        new SourceRecord(
                            sourcePartition,
                            sourceOffset,
                            topic,
                            Schema.STRING_SCHEMA,
                            query,
                            valueSchema,
                            value
                        );

                    result.add(record);
                }
                return result;
            } catch (InterruptedException | ExecutionException e) {
                return null;
            } finally {
                fetch = false;
            }
        }
    }

    private Schema getSchema(Class<?> type) {
        if (type.equals(Integer.class))
            return Schema.INT32_SCHEMA;

        return Schema.STRING_SCHEMA; // default schema
    }

}