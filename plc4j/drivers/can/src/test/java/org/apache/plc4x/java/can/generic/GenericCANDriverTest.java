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
package org.apache.plc4x.java.can.generic;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicReference;
import java.util.function.Consumer;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

public class GenericCANDriverTest {

    @Test
    void testConnection() throws PlcConnectionException {
        PlcConnection connection = new PlcDriverManager().getConnection("genericcan:virtualcan://");

        assertNotNull(connection);
        assertTrue(connection.isConnected());
        assertFalse(connection.getMetadata().canRead());
        assertTrue(connection.getMetadata().canWrite());
        assertTrue(connection.getMetadata().canSubscribe());
    }

    @Test
    @Disabled("This test requires working virtual CAN transport to be truly platform independent")
    void testSubscribeAndWrite() throws Exception {
//        PlcConnection connection1 = new PlcDriverManager().getConnection("genericcan:socketcan://vcan0");
//        PlcConnection connection2 = new PlcDriverManager().getConnection("genericcan:socketcan://vcan0");
        PlcConnection connection1 = new PlcDriverManager().getConnection("genericcan:virtualcan://");
        PlcConnection connection2 = connection1;

        CountDownLatch latch = new CountDownLatch(1);
        Byte field1 = 0x55;
        short field2 = 10;
        short field3 = 50;

        final AtomicReference<PlcSubscriptionEvent> plcEvent = new AtomicReference<>();
        connection1.subscriptionRequestBuilder()
            .addEventField("field1", "200:BYTE")
            .addEventField("field2", "200:UNSIGNED8")
            .addEventField("field3", "200:UNSIGNED8")
            .build().execute().whenComplete((reply, error) -> {
                if (error != null) {
                    fail(error);
                    return;
                }

                reply.getSubscriptionHandle("field1").register(new Consumer<PlcSubscriptionEvent>() {
                    @Override
                    public void accept(PlcSubscriptionEvent event) {
                        plcEvent.set(event);
                        latch.countDown();
                    }
                });
            });

        connection2.writeRequestBuilder()
            .addItem("f1", "200:BYTE", field1)
            .addItem("f2", "200:UNSIGNED8", field2)
            .addItem("f3", "200:UNSIGNED8", field3)
            .build().execute().whenComplete((reply, error) -> {
                if (error != null) {
                    fail(error);
                }
            }).get();

        latch.await();

        PlcSubscriptionEvent event = plcEvent.get();
        assertEquals(field1, event.getByte("field1"));
        assertEquals(field2, event.getShort("field2"));
        assertEquals(field3, event.getShort("field3"));

    }
}
