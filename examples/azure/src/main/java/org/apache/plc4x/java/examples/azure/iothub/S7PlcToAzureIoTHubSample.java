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

import com.microsoft.azure.sdk.iot.device.DeviceClient;
import com.microsoft.azure.sdk.iot.device.IotHubClientProtocol;
import com.microsoft.azure.sdk.iot.device.Message;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.TimeUnit;

public class S7PlcToAzureIoTHubSample {

    private static final Logger LOGGER = LoggerFactory.getLogger(S7PlcToAzureIoTHubSample.class);

    private static final String FIELD_NAME = "value";

    /**
     * Example code do demonstrate sending events from an S7 device to Microsoft Azure IoT Hub
     *
     * @param args Expected: [plc4x connection string, plc4x field address, IoT-Hub connection string].
     */
    public static void main(String[] args) throws Exception {
        if (args.length != 3) {
            System.out.println("Usage: S7PlcToAzureIoTHubSample " +
                "{plc4x-connection-string} {plc4x-field-address} {iot-hub-connection-string}");
            return;
        }

        String plc4xConnectionString = args[0];
        String addressString = args[1];
        String iotConnectionString = args[2];
        LOGGER.info("Connecting {}, {}, {}", plc4xConnectionString, addressString, iotConnectionString);

        // Open a connection to the remote PLC.
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(plc4xConnectionString)) {
            LOGGER.info("Connected");

            // Open a connection to the cloud service.
            DeviceClient client = new DeviceClient(iotConnectionString, IotHubClientProtocol.MQTT);
            client.open();


            // Prepare a read request.
            PlcReadRequest request = plcConnection.readRequestBuilder().get().addItem(FIELD_NAME, addressString).build();

            while (!Thread.currentThread().isInterrupted()) {
                // Simulate telemetry.
                PlcReadResponse response = request.execute().get();
                response.getAllLongs(FIELD_NAME)
                    .forEach(longValue -> {
                            String result = Long.toBinaryString(longValue);
                            LOGGER.info("Outputs {}", result);
                            Message msg = new Message("{ \"bits\" : \"" + result + "\"}");

                            // Send the message.
                            client.sendEventAsync(msg, (responseStatus, callbackContext) -> LOGGER.info("Received status: ", responseStatus), new Object());
                        }
                    );

                // Wait a second.
                TimeUnit.SECONDS.sleep(1);
            }

        }
    }
}
