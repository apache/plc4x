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
package org.apache.plc4x.java.slmp.manual;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.time.Duration;
import java.util.concurrent.TimeUnit;

/**
 * Manual smoke-test for the SLMP / MELSEC Communication 3E driver's <em>emulated</em>
 * subscription path.
 * <p>
 * SLMP has no native subscription mechanism; since {@code SlmpConnection} now extends
 * {@code PollingSubscriptionConnectionBase}, subscriptions are emulated by polling the tag via the
 * regular read path. This test registers two subscriptions against data register {@code D350:WORD}:
 * <ul>
 *   <li>a CYCLIC subscription that fires every second, and</li>
 *   <li>a CHANGE_OF_STATE subscription that only fires when the value actually changes.</li>
 * </ul>
 * It listens for 30 seconds, then unsubscribes cleanly.
 * <p>
 * To produce visible value changes, change data register {@code D350} on the PLC while this is running.
 */
public class ManualSlmpSubscriptionTest {

    private static final String CONNECTION_URL = "slmp://192.168.24.41";
    private static final String TAG_ADDRESS = "D350:WORD";
    private static final long LISTEN_DURATION_MS = 30_000L;

    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionFactory().getConnection(CONNECTION_URL)) {
            System.out.printf("Connected. metadata.subscribeSupported=%s%n",
                connection.getMetadata().isSubscribeSupported());

            PlcSubscriptionResponse subResp = connection.subscriptionRequestBuilder()
                .addCyclicTagAddress("cyclic-tag", TAG_ADDRESS, Duration.ofSeconds(1))
                .addChangeOfStateTagAddress("change-tag", TAG_ADDRESS)
                .build()
                .execute()
                .get(5, TimeUnit.SECONDS);

            PlcResponseCode cyclicCode = subResp.getResponseCode("cyclic-tag");
            PlcResponseCode changeCode = subResp.getResponseCode("change-tag");
            System.out.printf("Subscribe response: cyclic=%s, change-of-state=%s%n", cyclicCode, changeCode);
            if (cyclicCode != PlcResponseCode.OK || changeCode != PlcResponseCode.OK) {
                System.err.println("Subscription failed — aborting");
                return;
            }

            PlcSubscriptionHandle cyclicHandle = subResp.getSubscriptionHandle("cyclic-tag");
            PlcSubscriptionHandle changeHandle = subResp.getSubscriptionHandle("change-tag");
            cyclicHandle.register(event -> System.out.printf("[%s] CYCLIC          cyclic-tag = %s%n",
                event.getTimestamp(), event.getPlcValue("cyclic-tag")));
            changeHandle.register(event -> System.out.printf("[%s] CHANGE_OF_STATE change-tag = %s%n",
                event.getTimestamp(), event.getPlcValue("change-tag")));

            System.out.printf("Listening for emulated pushes for %ds...%n", LISTEN_DURATION_MS / 1000);
            Thread.sleep(LISTEN_DURATION_MS);

            connection.unsubscriptionRequestBuilder()
                .addHandles(cyclicHandle, changeHandle)
                .build()
                .execute()
                .get(5, TimeUnit.SECONDS);
            System.out.println("Unsubscribed.");
        }
    }
}
