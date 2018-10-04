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
import org.apache.kafka.connect.errors.ConnectException;
import org.apache.kafka.connect.sink.SinkRecord;
import org.apache.kafka.connect.sink.SinkTask;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.kafka.util.VersionUtil;

import java.util.Collection;
import java.util.Map;
import java.util.concurrent.ExecutionException;

public class Plc4xSinkTask extends SinkTask {
    private String url;

    private PlcConnection plcConnection;

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public void start(Map<String, String> props) {
        AbstractConfig config = new AbstractConfig(Plc4xSinkConnector.CONFIG_DEF, props);
        url = config.getString(Plc4xSinkConnector.URL_CONFIG);

        openConnection();

        if (!plcConnection.writeRequestBuilder().isPresent()) {
            throw new ConnectException("Writing not supported on this connection");
        }
    }

    @Override
    public void stop() {
        closeConnection();
    }

    @Override
    public void put(Collection<SinkRecord> records) {
        for (SinkRecord record: records) {
            String query = record.key().toString();
            Object value = record.value();
            PlcWriteRequest.Builder builder = plcConnection.writeRequestBuilder().get();
            PlcWriteRequest plcRequest = addToBuilder(builder, query, value).build();
            doWrite(plcRequest);
        }
    }

    // TODO: fix this
    private PlcWriteRequest.Builder addToBuilder(PlcWriteRequest.Builder builder, String query, Object obj) {
        Class<?> type = obj.getClass();

        if (type.equals(Integer.class)) {
            int value = (int) obj;
            builder.addItem(query, query, value);
        } else if (type.equals(String.class)) {
            String value = (String) obj;
            builder.addItem(query, query, value);
        }

        return builder;
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

    private void doWrite(PlcWriteRequest request) {
        try {
            request.execute().get();
        } catch (ExecutionException | InterruptedException e) {
            throw new ConnectException("Caught exception during write", e);
        }
    }

}