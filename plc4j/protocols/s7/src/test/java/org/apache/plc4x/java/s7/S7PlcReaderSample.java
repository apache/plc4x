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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.Optional;

public class S7PlcReaderSample {

    private final static Logger logger = LoggerFactory.getLogger(S7PlcReaderSample.class);

    /**
     * Example code do demonstrate using the S7 Plc Driver.
     *
     * @param args ignored.
     * @throws Exception something went wrong.
     */
    public static void main(String[] args) throws Exception {
        // Create a connection to the S7 PLC (s7://{hostname/ip}/{racknumber}/{slotnumber})
        logger.info("Connecting");
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection("s7://192.168.0.1/0/0")){
            logger.info("Connected");

            Optional<PlcReader> reader = plcConnection.getReader();
            // Check if this connection support reading of data.
            if (reader.isPresent()) {
                PlcReader plcReader = reader.get();

                // Prepare some address object for accessing fields in the PLC.
                // ({memory-area}/{byte-offset}[/{bit-offset}]
                // "bit-offset is only specified if the requested type is "bit"
                // NOTICE: This format is probably only valid when using a S7 connection.
                Address inputs = plcConnection.parseAddress("INPUTS/0");
                Address outputs = plcConnection.parseAddress("OUTPUTS/0");

                //////////////////////////////////////////////////////////
                // Read synchronously ...
                // NOTICE: the ".get()" immediately lets this thread pause till
                // the response is processed and available.
                PlcReadResponse plcReadResponse = plcReader.read(new PlcReadRequest(Byte.class, inputs)).get();
                List<ReadResponseItem> responseItems = plcReadResponse.getResponseItems();
                System.out.println("Inputs: " + responseItems.get(0).getValues().get(0));

                //////////////////////////////////////////////////////////
                // Read asynchronously ...
                /*Calendar start = Calendar.getInstance();
                CompletableFuture<PlcReadResponse<Byte>> asyncResponse = plcReader.read(
                    new BytePlcReadRequest(outputs));

                asyncResponse.thenAccept(bytePlcReadResponse -> {
                    Calendar end = Calendar.getInstance();
                    Byte dataAsync = bytePlcReadResponse.getValue();
                    System.out.println("Outputs: " + dataAsync + " (in " + (end.getTimeInMillis() - start.getTimeInMillis()) + "ms)");
                });

                // Simulate doing something else ...
                while (true) {
                    Thread.sleep(1);
                    if (asyncResponse.isDone()) {
                        break;
                    }
                }*/
            }
        }
        // Catch any exception or the application won't be able to finish if something goes wrong.
        catch (Exception e) {
            e.printStackTrace();
        }
        // The application would cleanly terminate after several seconds ... this just speeds things up.
        System.exit(0);
    }

}
