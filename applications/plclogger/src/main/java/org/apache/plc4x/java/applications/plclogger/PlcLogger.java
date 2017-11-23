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
package org.apache.plc4x.java.applications.plclogger;

import org.apache.edgent.providers.direct.DirectProvider;
import org.apache.edgent.topology.TStream;
import org.apache.edgent.topology.Topology;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.Address;
import org.apache.plc4x.java.api.messages.PlcSimpleReadRequest;
import org.apache.plc4x.java.api.messages.PlcSimpleReadResponse;
import org.apache.plc4x.java.api.types.ByteValue;

import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

public class PlcLogger {

    private final PlcReader plcReader;
    private final Address resourceAddress;
    private final int interval;

    private PlcLogger(PlcConnection plcConnection, String addressString, int interval) throws Exception {
        if(!(plcConnection instanceof PlcReader)) {
            throw new Exception("PlcConnection must be a PlcReader");
        }
        plcReader = (PlcReader) plcConnection;
        resourceAddress = plcConnection.parseAddress(addressString);
        this.interval = interval;
    }

    private Byte getPlcValue() throws PlcException, ExecutionException, InterruptedException {
        PlcSimpleReadResponse<ByteValue> plcReadResponse = plcReader.read(
            new PlcSimpleReadRequest<>(ByteValue.class, resourceAddress)).get();
        ByteValue data = plcReadResponse.getValue();
        return data.getValue();
    }

    private void run() throws Exception {
        DirectProvider dp = new DirectProvider();
        Topology top = dp.newTopology();
        TStream<Byte> source = top.poll(() -> {
                try {
                    return getPlcValue();
                } catch (Exception e) {
                    e.printStackTrace();
                }
                return null;
            },
            interval, TimeUnit.MILLISECONDS);
        source.print();
        dp.submit(top);
    }

    public static void main(String[] args) {
        if(args.length != 3) {
            System.out.println("Usage: PlcLogger {connection-string} {resource-address-string} {interval-ms}");
            System.out.println("Example: PlcLogger s7://192.168.0.1/0/0 INPUTS/0 10");
        }

        String connectionString = args[0];
        String addressString = args[1];
        Integer interval = Integer.valueOf(args[2]);
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(connectionString)) {
            // Create a connection to the S7 PLC (s7://{hostname/ip}/{racknumber}/{slotnumber})
            plcConnection.connect();

            // Initialize the logger itself
            PlcLogger logger = new PlcLogger(plcConnection, addressString, interval);

            // Start the logging ...
            logger.run();

            // Yeah ... well prevent the application from exiting ;-)
            while (true) {
                Thread.sleep(1000);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }


}
