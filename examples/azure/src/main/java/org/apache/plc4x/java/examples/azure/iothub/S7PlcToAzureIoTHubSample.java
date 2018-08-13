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
package org.apache.plc4x.java.examples.azure.iothub;

import com.microsoft.azure.sdk.iot.device.*;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadResponse;
import org.apache.plc4x.java.api.model.Address;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Optional;

public class S7PlcToAzureIoTHubSample {

    private static final Logger logger = LoggerFactory.getLogger(S7PlcToAzureIoTHubSample.class);

    private static IotHubClientProtocol protocol = IotHubClientProtocol.MQTT;
    private static DeviceClient client;

    public static class Callback implements IotHubEventCallback {

        @Override
        public void execute(IotHubStatusCode iotHubStatusCode, Object ctx) {
            System.out.println("Received status: " + iotHubStatusCode.name());
        }
    }

    /**
     * Example code do demonstrate sending events from an S7 device to Microsoft Azure IoT Hub
     *
     * @param args Expected: [plc4x connection string, plc4x address, IoT-Hub connection string].
     */
    public static void main(String[] args) {
        logger.info("Connecting");
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(args[0])) {
            logger.info("Connected");

            client = new DeviceClient(args[2], protocol);
            client.open();

            Optional<PlcReader> reader = plcConnection.getReader();

            if (reader.isPresent()) {
                PlcReader plcReader = reader.get();
                Address outputs = plcConnection.parseAddress(args[1]);

                while (true) {
                    // Simulate telemetry.
                    TypeSafePlcReadResponse<Byte> plcReadResponse = plcReader.read(
                        new TypeSafePlcReadRequest<>(Byte.class, outputs)).get();

                    System.out.println("Outputs: " + Long.toBinaryString(plcReadResponse.getResponseItem()
                        .orElseThrow(() -> new IllegalStateException("No response available"))
                        .getValues().get(0).longValue()));

                    plcReadResponse.getResponseItem().map(byt -> {
                            String result = Long.toBinaryString(byt.getValues().get(0).longValue());
                            Message msg = new Message("{ \"bits\" : \"" + result + "\"}");
                            // Send the message.
                            Callback callback = new Callback();

                            client.sendEventAsync(msg, callback, new Object());
                            return byt;
                        }
                    );

                    Thread.sleep(1000);
                }
            }

        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
