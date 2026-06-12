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
package org.apache.plc4x.java.s7.manual;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;

import java.util.concurrent.TimeUnit;

/**
 * Manual smoke-test for the S7Comm UserData "alarm push" subscription path against an
 * S7-300 (or any classic-family CPU). Connects, subscribes to the {@code ALM} tag, and
 * prints every alarm indication the PLC pushes for the next 60 seconds.
 *
 * <p>Trigger an alarm in the PLC's diagnostic buffer (e.g. force a STOP/RUN transition,
 * trigger an SFB OB, or use Step7 to inject a test message via {@code Test &gt; Modify
 * Variable}) to see entries appear here.
 *
 * <p>Run with: {@code mvn -pl plc4j/drivers/s7 -DskipTests=true compile test-compile
 * exec:java -Dexec.mainClass=org.apache.plc4x.java.s7.manual.ManualWallS7300AlarmSubscribe}
 */
public class ManualWallS7300AlarmSubscribe {

    private static final String CONNECTION_URL = "s7://192.168.24.60?local-device-group=PG_OR_PC";
    private static final long LISTEN_DURATION_MS = 600_000L;

    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionManager().getConnection(CONNECTION_URL)) {
            System.out.printf("Connected. metadata.subscribeSupported=%s%n",
                connection.getMetadata().isSubscribeSupported());

            PlcSubscriptionResponse subResp = connection.subscriptionRequestBuilder()
                .addEventTagAddress("alarms", "ALM")
                .build()
                .execute()
                .get(5, TimeUnit.SECONDS);

            PlcResponseCode code = subResp.getResponseCode("alarms");
            System.out.printf("Subscribe response: %s%n", code);
            if (code != PlcResponseCode.OK) {
                System.err.println("Subscription failed — aborting");
                return;
            }

            PlcSubscriptionHandle handle = subResp.getSubscriptionHandle("alarms");
            handle.register(event -> {
                PlcValue value = event.getPlcValue("alarms");
                System.out.printf("[%s] alarm: %s%n", event.getTimestamp(), formatStruct(value));
            });

            System.out.printf("Listening for alarms for %ds — trigger an alarm in the PLC...%n",
                LISTEN_DURATION_MS / 1000);
            Thread.sleep(LISTEN_DURATION_MS);

            connection.unsubscriptionRequestBuilder()
                .addHandles(handle)
                .build()
                .execute()
                .get(5, TimeUnit.SECONDS);
            System.out.println("Unsubscribed.");
        }
    }

    private static String formatStruct(PlcValue value) {
        if (value == null || !value.isStruct()) {
            return String.valueOf(value);
        }
        StringBuilder sb = new StringBuilder("{");
        boolean first = true;
        for (String key : value.getKeys()) {
            if (!first) sb.append(", ");
            first = false;
            PlcValue v = value.getValue(key);
            // Render state masks as hex for readability, eventId likewise.
            if ("eventId".equals(key) || key.endsWith("State")) {
                sb.append(key).append("=0x").append(Long.toHexString(v.getLong()).toUpperCase());
            } else {
                sb.append(key).append("=").append(v);
            }
        }
        return sb.append("}").toString();
    }
}
