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
package org.apache.plc4x.kafka.sink;

import org.apache.kafka.common.config.ConfigException;
import org.apache.kafka.connect.errors.ConnectException;
import org.apache.kafka.connect.sink.SinkRecord;
import org.apache.kafka.connect.sink.SinkTask;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.atomic.AtomicBoolean;

public class Plc4xSinkTask extends SinkTask {

    private static Logger log = LoggerFactory.getLogger(Plc4xSinkTask.class);

    private Plc4xSinkConfig config;
    private PlcConnection plcConnection;
    private PlcWriter writer;
    private AtomicBoolean running = new AtomicBoolean(false);

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public void start(Map<String, String> properties) {
        try {
            config = new Plc4xSinkConfig(properties);
        } catch (ConfigException e) {
            throw new ConnectException("Couldn't start JdbcSourceTask due to configuration error", e);
        }
        final String url = config.getString(Plc4xSinkConfig.PLC_CONNECTION_STRING_CONFIG);

        try {
            plcConnection = new PlcDriverManager().getConnection(url);
            Optional<PlcWriter> writerOptional = plcConnection.getWriter();
            if(!writerOptional.isPresent()) {
                throw new ConnectException("PlcWriter not available for this type of connection");
            }
            writer = writerOptional.get();
            running.set(true);
        } catch (PlcConnectionException e) {
            throw new ConnectException("Caught exception while connecting to PLC", e);
        }
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
    public void put(Collection<SinkRecord> records) {
        if((plcConnection != null) && plcConnection.isConnected() && (writer != null)) {
            // Prepare the write request.
            PlcWriteRequest.Builder builder = writer.writeRequestBuilder();
            for (SinkRecord record : records) {
                // TODO: Somehow get the payload from the kafka SinkRecord and create a writeRequestItem from that ...
                // TODO: Replace this dummy with something real ...
                Map<String, Object> value = new HashMap<>(); //(Map<String, String>) record.value()
                String addressString = (String) value.get("address");
                List<Byte> values = (List<Byte>) value.get("values");

                builder.addItem(addressString, addressString, values.toArray(new Byte[0]));
            }
            PlcWriteRequest writeRequest = builder.build();

            // Send the write request to the PLC.
            try {
                PlcWriteResponse<?> plcWriteResponse = writer.write(writeRequest).get();
                for (String fieldName : plcWriteResponse.getFieldNames()) {
                    if(plcWriteResponse.getResponseCode(fieldName) != PlcResponseCode.OK) {
                        // TODO: Do Something if writing this particular item wasn't successful ...
                        log.error("Error writing a value to PLC");
                    }
                }
            } catch (ExecutionException | InterruptedException e) {
                log.error("Error writing values to PLC", e);
            }
        }
    }

}