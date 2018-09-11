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
import org.apache.plc4x.kafka.util.VersionUtil;

import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.*;

/**
 * Source Connector Task polling the data source at a given rate.
 * A timer thread is scheduled which sets the fetch flag to true every rate milliseconds.
 * When poll() is invoked, the calling thread waits until the fetch flag is set for WAIT_LIMIT_MILLIS.
 * If the flag does not become true, the method returns null, otherwise a fetch is performed.
 */
public class Plc4xSourceTask extends SourceTask {
    private final static long WAIT_LIMIT_MILLIS = 100;
    private final static long TIMEOUT_LIMIT_MILLIS = 5000;

    private String topic;
    private String url;
    private List<String> queries;

    private PlcConnection plcConnection;
    private PlcReader plcReader;
    private PlcReadRequest plcRequest;

    // TODO: should we use shared (static) thread pool for this?
    private ScheduledExecutorService scheduler;
    private ScheduledFuture<?> timer;
    private boolean fetch = true;

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public void start(Map<String, String> props) {
        AbstractConfig config = new AbstractConfig(Plc4xSourceConnector.CONFIG_DEF, props);
        topic = config.getString(Plc4xSourceConnector.TOPIC_CONFIG);
        url = config.getString(Plc4xSourceConnector.URL_CONFIG);
        queries = config.getList(Plc4xSourceConnector.QUERIES_CONFIG);

        openConnection();

        plcReader = plcConnection.getReader()
            .orElseThrow(() -> new ConnectException("PlcReader not available for this type of connection"));


        PlcReadRequest.Builder builder = plcReader.readRequestBuilder();
        for (String query : queries) {
            builder.addItem(query, query);
        }
        plcRequest = builder.build();

        int rate = Integer.valueOf(props.get(Plc4xSourceConnector.RATE_CONFIG));
        scheduler = Executors.newScheduledThreadPool(1);
        timer = scheduler.scheduleAtFixedRate(Plc4xSourceTask.this::scheduleFetch, rate, rate, TimeUnit.MILLISECONDS);
    }

    @Override
    public void stop() {
        timer.cancel(true);
        scheduler.shutdown();
        closeConnection();
    }

    @Override
    public List<SourceRecord> poll() throws InterruptedException {
        return awaitFetch(WAIT_LIMIT_MILLIS) ? doFetch() : null;
    }

    private void openConnection() {
        try {
            plcConnection = new PlcDriverManager().getConnection(url);
            plcConnection.connect();
        } catch (PlcConnectionException e) {
            throw new ConnectException("Could not establish a PLC connection", e);
        }
    }

    private void closeConnection() {
        if (plcConnection != null) {
            try {
                plcConnection.close();
            } catch (Exception e) {
                throw new RuntimeException("Caught exception while closing connection to PLC", e);
            }
        }
    }

    /**
     * Schedule next fetch operation.
     */
    private synchronized void scheduleFetch() {
        fetch = true;
        notify();
    }

    /**
     * Wait for next scheduled fetch operation.
     * @param milliseconds maximum time to wait
     * @throws InterruptedException if the thread is interrupted
     * @return true if a fetch should be performed, false otherwise
     */
    private synchronized boolean awaitFetch(long milliseconds) throws InterruptedException {
        if (!fetch) {
            wait(milliseconds);
        }
        try {
            return fetch;
        } finally {
            fetch = false;
        }
    }

    private List<SourceRecord> doFetch() throws InterruptedException {
        final CompletableFuture<PlcReadResponse<?>> response = plcReader.read(plcRequest);
        try {
            final PlcReadResponse<?> received = response.get(TIMEOUT_LIMIT_MILLIS, TimeUnit.MILLISECONDS);
            return extractValues(received);
        } catch (ExecutionException e) {
            throw new ConnectException("Could not fetch data from source", e);
        } catch (TimeoutException e) {
            throw new ConnectException("Timed out waiting for data from source", e);
        }
    }

    private List<SourceRecord> extractValues(PlcReadResponse<?> response) {
        final List<SourceRecord> result = new LinkedList<>();
        for (String query : queries) {
            final PlcResponseCode rc = response.getResponseCode(query);
            if (!rc.equals(PlcResponseCode.OK))  {
                continue;
            }

            Object rawValue = response.getObject(query);
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
    }

    private Schema getSchema(Class<?> type) {
        if (type.equals(Byte.class))
            return Schema.INT8_SCHEMA;

        if (type.equals(Short.class))
            return Schema.INT16_SCHEMA;

        if (type.equals(Integer.class))
            return Schema.INT32_SCHEMA;

        if (type.equals(Long.class))
            return Schema.INT64_SCHEMA;

        return Schema.STRING_SCHEMA; // default case; invoke .toString on value
    }

}