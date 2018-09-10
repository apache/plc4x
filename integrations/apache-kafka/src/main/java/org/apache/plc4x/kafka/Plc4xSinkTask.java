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

import org.apache.kafka.connect.errors.ConnectException;
import org.apache.kafka.connect.sink.SinkRecord;
import org.apache.kafka.connect.sink.SinkTask;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.kafka.util.VersionUtil;

import java.util.Collection;
import java.util.Map;
import java.util.concurrent.ExecutionException;

public class Plc4xSinkTask extends SinkTask {
    private final static String FIELD_KEY = "key"; // TODO: is this really necessary?

    private String url;
    private String query;

    private PlcConnection plcConnection;
    private PlcWriter plcWriter;

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public void start(Map<String, String> props) {
        url = props.get(Plc4xSinkConnector.URL_CONFIG);
        query = props.get(Plc4xSinkConnector.QUERY_CONFIG);

        openConnection();

        plcWriter = plcConnection.getWriter()
            .orElseThrow(() -> new ConnectException("PlcReader not available for this type of connection"));
    }

    @Override
    public void stop() {
        closeConnection();
    }

    @Override
    public void put(Collection<SinkRecord> records) {
        for (SinkRecord record: records) {
            String value = record.value().toString(); // TODO: implement other data types
            PlcWriteRequest plcRequest = plcWriter.writeRequestBuilder().addItem(FIELD_KEY, query, value).build();
            doWrite(plcRequest);
        }
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
            plcWriter.write(request).get();
        } catch (ExecutionException | InterruptedException e) {
            throw new ConnectException("Caught exception during write", e);
        }
    }

}