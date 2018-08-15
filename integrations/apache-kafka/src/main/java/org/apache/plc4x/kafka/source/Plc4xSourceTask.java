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
import org.apache.kafka.common.utils.SystemTime;
import org.apache.kafka.common.utils.Time;
import org.apache.kafka.connect.errors.ConnectException;
import org.apache.kafka.connect.source.SourceRecord;
import org.apache.kafka.connect.source.SourceTask;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.atomic.AtomicBoolean;

public class Plc4xSourceTask extends SourceTask {

    static final Logger log = LoggerFactory.getLogger(Plc4xSourceTask.class);

    private Plc4xSourceConfig config;
    private PlcConnection plcConnection;
    private PlcReader reader;
    private PlcReadRequest readRequest;
    private AtomicBoolean running = new AtomicBoolean(false);

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public void start(Map<String, String> properties) {
        try {
            config = new Plc4xSourceConfig(properties);
        } catch (ConfigException e) {
            throw new ConnectException("Couldn't start JdbcSourceTask due to configuration error", e);
        }
        final String url = config.getString(Plc4xSourceConfig.PLC_CONNECTION_STRING_CONFIG);

        try {
            plcConnection = new PlcDriverManager().getConnection(url);
            Optional<PlcReader> readerOptional = plcConnection.getReader();
            if(!readerOptional.isPresent()) {
                throw new ConnectException("PlcReader not available for this type of connection");
            }
            reader = readerOptional.get();
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
    public List<SourceRecord> poll() throws InterruptedException {
        if((plcConnection != null) && plcConnection.isConnected() && (reader != null)) {
            final List<SourceRecord> results = new LinkedList<>();

            try {
                PlcReadResponse plcReadResponse = reader.read(readRequest).get();

                for (ReadResponseItem<?> responseItem : plcReadResponse.getResponseItems()) {
                    Address address = responseItem.getRequestItem().getAddress();
                    List<?> values = responseItem.getValues();

                    // TODO: Implement Sending this information to Kafka ...
                    //results.add(new SourceRecord())
                }
            } catch (ExecutionException e) {
                log.error("Error reading values from PLC", e);
            }

            return results;
        }
        return null;
    }

}