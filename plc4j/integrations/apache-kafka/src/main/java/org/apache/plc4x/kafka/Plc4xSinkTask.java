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

import org.apache.kafka.connect.transforms.Transformation;
import org.apache.kafka.connect.sink.SinkRecord;
import org.apache.kafka.connect.sink.SinkTask;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;
import org.apache.plc4x.kafka.util.VersionUtil;

import java.util.concurrent.ExecutionException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

/**
 * Source Connector Task polling the data source at a given rate.
 * A timer thread is scheduled which sets the fetch flag to true every rate milliseconds.
 * When poll() is invoked, the calling thread waits until the fetch flag is set for WAIT_LIMIT_MILLIS.
 * If the flag does not become true, the method returns null, otherwise a fetch is performed.
 */
public class Plc4xSinkTask extends SinkTask {

    private static final Logger log = LoggerFactory.getLogger(Plc4xSinkTask.class);

    /*
     * Config of the task.
     */
    static final String CONNECTION_NAME_CONFIG = "connection-name";
    private static final String CONNECTION_NAME_STRING_DOC = "Connection Name";

    static final String PLC4X_CONNECTION_STRING_CONFIG = "connectionString";
    private static final String PLC4X_CONNECTION_STRING_DOC = "PLC4X Connection String";

    static final String PLC4X_TOPIC_CONFIG = "topic";
    private static final String PLC4X_TOPIC_DOC = "Task Topic";

    private static final ConfigDef CONFIG_DEF = new ConfigDef()
        .define(CONNECTION_NAME_CONFIG,
                ConfigDef.Type.STRING,
                ConfigDef.Importance.HIGH,
                CONNECTION_NAME_STRING_DOC)
        .define(PLC4X_CONNECTION_STRING_CONFIG,
                ConfigDef.Type.STRING,
                ConfigDef.Importance.HIGH,
                PLC4X_CONNECTION_STRING_DOC)
        .define(PLC4X_TOPIC_CONFIG,
                ConfigDef.Type.STRING,
                ConfigDef.Importance.HIGH,
                PLC4X_TOPIC_DOC);

    /*
     * Configuration of the output.
     */
    private static final String SINK_NAME_FIELD = "sink-name";
    private static final String SINK_TOPIC_FIELD = "topic";


    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    private PlcDriverManager driverManager;
    private Transformation<SinkRecord> transformation;
    private String plc4xConnectionString;
    private String plc4xTopic;

    @Override
    public void start(Map<String, String> props) {
        AbstractConfig config = new AbstractConfig(CONFIG_DEF, props);
        String connectionName = config.getString(CONNECTION_NAME_CONFIG);
        plc4xConnectionString = config.getString(PLC4X_CONNECTION_STRING_CONFIG);
        plc4xTopic = config.getString(PLC4X_TOPIC_CONFIG);
        Map<String, String> topics = new HashMap<>();
        log.info("Creating Pooled PLC4x driver manager");
        driverManager = new PooledPlcDriverManager();
    }

    @Override
    public void stop() {
        synchronized (this) {
            notifyAll(); // wake up thread waiting in awaitFetch
        }
    }

    @Override
    public void put(Collection<SinkRecord> records) {
        if (records.isEmpty()) {
            return;
        }

        for (SinkRecord r: records) {
            Struct record = (Struct) r.value();
            String topic = r.topic();
            String address = record.getString("address");
            String value = record.getString("value");
            Long expires = record.getInt64("expires");

            if (!topic.equals(plc4xTopic) || plc4xTopic.equals("")) {
                log.debug("Ignoring write request recived on wrong topic");
                return;
            }

            if ((System.currentTimeMillis() > expires) & !(expires == 0)) {
                log.warn("Write request has expired {}, discarding {}", System.currentTimeMillis(), address);
                return;
            }

            PlcConnection connection = null;
            try {
                connection = driverManager.getConnection(plc4xConnectionString);
            } catch (PlcConnectionException e) {
                log.warn("Failed to Open Connection {}", plc4xConnectionString);
            }

            final PlcWriteRequest.Builder builder = connection.writeRequestBuilder();
            PlcWriteRequest writeRequest;
            try {
                //If an array value is passed instead of a single value then convert to a String array
                if ((value.charAt(0) == '[') && (value.charAt(value.length() - 1) == ']')) {
                    String[] values = value.substring(1,value.length() - 1).split(",");
                    builder.addItem(address, address, values);
                } else {
                    builder.addItem(address, address, value);
                }

                writeRequest = builder.build();
            } catch (Exception e) {
                //When building a request we want to discard the write if there is an error.
                log.warn("Failed to Write to {}", plc4xConnectionString);
                return;
            }

            try {
                writeRequest.execute().get();
                log.info("Wrote {} to device {}", address, plc4xConnectionString);
            } catch (InterruptedException | ExecutionException e) {
                log.warn("Failed to Write to {}", plc4xConnectionString);
            }

            //try {
            //    connection.close();
            //} catch (Exception e) {
            //    log.warn("Failed to Close {}", plc4xConnectionString);
            //}
        }
        return;
    }
}
