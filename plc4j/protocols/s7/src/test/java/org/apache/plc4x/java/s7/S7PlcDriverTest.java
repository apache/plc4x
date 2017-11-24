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
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.Address;
import org.apache.plc4x.java.api.messages.PlcSimpleReadRequest;
import org.apache.plc4x.java.api.messages.PlcSimpleReadResponse;
import org.apache.plc4x.java.s7.connection.S7PlcConnection;
import org.apache.plc4x.java.api.types.ByteValue;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.Calendar;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;

public class S7PlcDriverTest {

    @Test
    @Tag("fast")
    void getConnectionTest() throws PlcException {
        S7PlcConnection s7Connection = (S7PlcConnection)
            new PlcDriverManager().getConnection("s7://localhost/1/2");
        Assertions.assertEquals(s7Connection.getHostName(), "localhost");
        Assertions.assertEquals(s7Connection.getRack(), 1);
        Assertions.assertEquals(s7Connection.getSlot(), 2);
    }

    /**
     * In this test case the 's7' driver should report an invalid url format.
     *
     * @throws PlcException something went wrong
     */
    @Test
    @Tag("fast")
    void getConnectionInvalidUrlTest() throws PlcException {
        Assertions.assertThrows(PlcConnectionException.class,
            () -> new PlcDriverManager().getConnection("s7://localhost/hurz/2"));
    }

    /**
     * In this test case the 's7' driver should report an error as this protocol
     * doesn't support authentication.
     *
     * @throws PlcException something went wrong
     */
    @Test
    @Tag("fast")
    void getConnectionWithAuthenticationTest() throws PlcException {
        Assertions.assertThrows(PlcConnectionException.class,
            () -> new PlcDriverManager().getConnection("s7://localhost/1/2",
                new PlcUsernamePasswordAuthentication("user", "pass")));
    }

    /**
     * Example code do demonstrate using the S7 Plc Driver.
     *
     * @param args ignored.
     * @throws Exception something went wrong.
     */
    public static void main(String[] args) throws Exception {
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection("s7://192.168.0.1/0/0")){
            // Create a connection to the S7 PLC (s7://{hostname/ip}/{racknumber}/{slotnumber})
            plcConnection.connect();

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
                PlcSimpleReadResponse<ByteValue> plcReadResponse = plcReader.read(
                    new PlcSimpleReadRequest<>(ByteValue.class, inputs)).get();
                ByteValue data = plcReadResponse.getValue();
                System.out.println("Inputs: " + data.getValue());

                //////////////////////////////////////////////////////////
                // Read asynchronously ...
                Calendar start = Calendar.getInstance();
                CompletableFuture<PlcSimpleReadResponse<ByteValue>> asyncResponse = plcReader.read(
                    new PlcSimpleReadRequest<>(ByteValue.class, outputs));

                // Simulate doing something else ...
                System.out.println("Processing: ");
                while (true) {
                    // I had to make sleep this small or it would have printed only one "."
                    // On my system the average response time with a siemens s7-1200 was 5ms.
                    Thread.sleep(1);
                    System.out.print(".");
                    if (asyncResponse.isDone()) {
                        break;
                    }
                }
                System.out.println();

                Calendar end = Calendar.getInstance();
                plcReadResponse = asyncResponse.get();
                data = plcReadResponse.getValue();
                System.out.println("Outputs: " + data.getValue() + " (in " + (end.getTimeInMillis() - start.getTimeInMillis()) + "ms)");
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
