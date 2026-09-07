/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.eip.base.manual;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.time.Duration;
import java.util.concurrent.TimeUnit;

/**
 * Manual smoke-test for the EtherNet/IP driver's <em>emulated</em> subscription path.
 * <p>
 * EtherNet/IP has no native subscription mechanism; since {@code EipTcpConnection} now extends
 * {@code PollingSubscriptionConnectionBase}, subscriptions are emulated by polling the tag via the
 * regular read path. This test asks for the value of {@code %out01} once per second for 30 seconds
 * using a CYCLIC subscription, then unsubscribes cleanly.
 * <p>
 * To produce visible value changes, change the tag on the PLC while this is running.
 */
public class ManualEipSubscriptionTest {

    private static final String CONNECTION_URL = "eip://192.168.24.32";
    private static final String TAG_ADDRESS = "%out01";
    private static final long LISTEN_DURATION_MS = 30_000L;

    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionFactory().getConnection(CONNECTION_URL)) {
            System.out.printf("Connected. metadata.subscribeSupported=%s%n",
                connection.getMetadata().isSubscribeSupported());

            PlcSubscriptionResponse subResp = connection.subscriptionRequestBuilder()
                .addCyclicTagAddress("subscription-tag", TAG_ADDRESS, Duration.ofSeconds(1))
                .build()
                .execute()
                .get(5, TimeUnit.SECONDS);

            PlcResponseCode code = subResp.getResponseCode("subscription-tag");
            System.out.printf("Subscribe response: %s%n", code);
            if (code != PlcResponseCode.OK) {
                System.err.println("Subscription failed — aborting");
                return;
            }

            PlcSubscriptionHandle handle = subResp.getSubscriptionHandle("subscription-tag");
            handle.register(event -> System.out.printf("[%s] subscription-tag = %s%n",
                event.getTimestamp(), event.getPlcValue("subscription-tag")));

            System.out.printf("Listening for emulated cyclic pushes for %ds...%n", LISTEN_DURATION_MS / 1000);
            Thread.sleep(LISTEN_DURATION_MS);

            connection.unsubscriptionRequestBuilder()
                .addHandles(handle)
                .build()
                .execute()
                .get(5, TimeUnit.SECONDS);
            System.out.println("Unsubscribed.");
        }
    }
}
