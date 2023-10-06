/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.examples.hellonats;

import io.nats.client.*;
import io.nats.client.api.StorageType;
import io.nats.client.api.StreamConfiguration;
import io.nats.client.api.StreamInfo;
import io.nats.client.support.JsonUtils;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.List;
import java.util.concurrent.ExecutionException;

public class HelloNats {

    private static final Logger logger = LoggerFactory.getLogger(HelloNats.class);

    /**
     * Example code do demonstrate using PLC4X.
     *
     * @param args ignored.
     */
    public static void main(String[] args) throws Exception {
        CliOptions options = CliOptions.fromArgs(args);
        if (options == null) {
            CliOptions.printHelp();
            // Could not parse.
            System.exit(1);
        }

        try (Connection natsConnection = Nats.connect(options.getNatsServerConnectionString())) {
            final JetStreamManagement jsm = natsConnection.jetStreamManagement();

            final StreamConfiguration streamConfiguration = StreamConfiguration.builder().name(options.getNatsNodeName()).subjects(options.getNatsTopic()).storageType(StorageType.Memory).build();
            final StreamInfo streamInfo = jsm.addStream(streamConfiguration);
            JsonUtils.printFormatted(streamInfo);

            final JetStream jetStream = natsConnection.jetStream();

            final Dispatcher dispatcher = natsConnection.createDispatcher();
            final JetStreamSubscription jetStreamSubscription = jetStream.subscribe(options.getNatsTopic(), dispatcher, msg -> {
                final String connectionUrl = msg.getHeaders().getFirst("connection-url");
                final List<String> tags = msg.getHeaders().get("tags");

                // Establish a connection to the plc using the url provided as first argument
                try (PlcConnection plcConnection = new PlcDriverManager().getConnection(connectionUrl)) {

                    // Check if this connection support reading of data.
                    if (!plcConnection.getMetadata().canRead()) {
                        logger.error("This connection doesn't support reading.");
                        return;
                    }

                    // Create a new read request:
                    // - Give the single item requested the alias name "value"
                    PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
                    for (int i = 0; i < tags.size(); i++) {
                        builder.addItem("value-" + i, tags.get(i));
                    }
                    PlcReadRequest readRequest = builder.build();

                    // Actually execute the read request (synchronously)
                    PlcReadResponse response = readRequest.execute().get();

                    for (String tagName : response.getTagNames()) {
                        if (response.getResponseCode(tagName) == PlcResponseCode.OK) {
                            final PlcValue plcValue = response.getPlcValue(tagName);
                        }
                        // Something went wrong, to output an error message instead.
                        else {
                            logger.error("Error[{}]: {}", tagName, response.getResponseCode(tagName).name());
                        }
                    }
                } catch (ExecutionException e) {
                    logger.error("Error[{}]: {}", tagName, response.getResponseCode(tagName).name());
                }
            }, false);
            natsConnection.flush(Duration.ofSeconds(1));
        }
    }

}
