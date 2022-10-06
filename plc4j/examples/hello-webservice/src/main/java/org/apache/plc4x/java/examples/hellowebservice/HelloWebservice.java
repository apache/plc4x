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
package org.apache.plc4x.java.examples.hellowebservice;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.JsonObject;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.util.EntityUtils;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.value.PlcValue;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.Scanner;
import java.util.function.Consumer;

public class HelloWebservice {

    private static final Logger logger = LoggerFactory.getLogger(HelloWebservice.class);

    private final CliOptions options;

    public HelloWebservice(CliOptions options) {
        this.options = options;
    }

    public void run() throws Exception {
        // Establish a connection to the plc.
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(options.getConnectionString())) {

            // Check if this connection support subscriptions.
            if (!plcConnection.getMetadata().canSubscribe()) {
                logger.error("This connection doesn't support subscriptions.");
                return;
            }

            // Create a new read request:
            // - Give the single item requested the alias name "value"
            final PlcSubscriptionRequest.Builder builder = plcConnection.subscriptionRequestBuilder();
            for (int i = 0; i < options.getFieldAddress().length; i++) {
                builder.addChangeOfStateField("value-" + i, options.getFieldAddress()[i]);
            }
            PlcSubscriptionRequest subscriptionRequest = builder.build();

            // Execute the subscription response.
            final PlcSubscriptionResponse subscriptionResponse = subscriptionRequest.execute().get();

            // Attach handlers for the incoming data.
            for (String subscriptionName : subscriptionResponse.getFieldNames()) {
                final PlcSubscriptionHandle subscriptionHandle =
                    subscriptionResponse.getSubscriptionHandle(subscriptionName);
                subscriptionHandle.register(new ValueChangeHandler(options.getWebserviceUrl()));
            }

            // Wait for the user to press "Enter" to abort the program.
            Scanner scanner = new Scanner(System.in);
            try {
                logger.info("Please press Enter to exit program.");
                scanner.nextLine();
                logger.info("Finishing");
            } catch (IllegalStateException e) {
                // System.in has been closed
                logger.error("System.in was closed; exiting");
            }
        }
    }

    /**
     * Example code do demonstrate using PLC4X Subcription API.
     */
    public static void main(String[] args) throws Exception {
        CliOptions options = CliOptions.fromArgs(args);
        if (options == null) {
            CliOptions.printHelp();
            // Could not parse.
            System.exit(1);
        }

        HelloWebservice subscriptionApplication = new HelloWebservice(options);

        subscriptionApplication.run();

        System.exit(0);
    }

    private static class ValueChangeHandler implements Consumer<PlcSubscriptionEvent> {

        private final String webserviceUrl;
        private final Gson gson;
        private final DateTimeFormatter formatter;

        public ValueChangeHandler(String webserviceUrl) {
            this.webserviceUrl = webserviceUrl;
            GsonBuilder gsonBuilder = new GsonBuilder();
            gsonBuilder.setPrettyPrinting();
            gson = gsonBuilder.create();
            formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd'T'hh:mm:ssZ");
        }

        @Override
        public void accept(PlcSubscriptionEvent plcSubscriptionEvent) {
            logger.info("Incoming event:");
            // Iterate over all the fields in this event and then simply output
            // them to the console in a JSON format.
            for (String fieldName : plcSubscriptionEvent.getFieldNames()) {
                final PlcValue plcValue = plcSubscriptionEvent.getPlcValue(fieldName);

                // Create a JSON object that fits the structure of my remote webservice.
                JsonObject output = new JsonObject();
                output.addProperty("time", ZonedDateTime.now().format(formatter));
                output.addProperty("device-id", plcValue.getValue("sourceAddress").getString());
                output.addProperty("target-id", plcValue.getValue("targetAddress").getString());
                output.addProperty("type", plcValue.getValue("description").getString());
                output.addProperty("category", plcValue.getValue("function").getString());
                JsonObject data = new JsonObject();
                data.addProperty("unit-of-measurement", plcValue.getValue("unitOfMeasurement").getString());
                data.addProperty("value", plcValue.getValue("value").getString());
                output.add("data", data);

                // Send the the json payload to the remote webservice.
                HttpPost post = new HttpPost(webserviceUrl);
                try {
                    post.setEntity(new StringEntity(gson.toJson(output)));
                } catch (UnsupportedEncodingException e) {
                    logger.error("Error encoding json string entity", e);
                }
                try (CloseableHttpClient httpClient = HttpClients.createDefault();
                     CloseableHttpResponse response = httpClient.execute(post)) {

                    String result = EntityUtils.toString(response.getEntity());
                    logger.info(String.format("Got '%s' from remote", result));
                } catch (IOException e) {
                    logger.error("Error sending payload to remote webservice.", e);
                }
            }
        }
    }

}
