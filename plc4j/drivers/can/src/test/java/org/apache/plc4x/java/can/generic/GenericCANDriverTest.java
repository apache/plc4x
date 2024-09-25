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
package org.apache.plc4x.java.can.generic;

import io.netty.buffer.ByteBuf;
import io.netty.channel.Channel;
import io.netty.channel.embedded.Plc4xEmbeddedChannel;
import java.lang.reflect.Array;
import java.util.Map;
import java.util.Map.Entry;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicReference;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest.Builder;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.spi.connection.ChannelExposingConnection;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import static java.util.Map.entry;
import static org.junit.jupiter.api.Assertions.*;

/**
 * Test of generic can driver with virtual can transport.
 *
 * This test have additional role of confirming end to end behavior of driver and transport layer.
 * The virtual can transport rely on netty channel/pipeline infrastructure, hence it does not have
 * any backend such as memory queue. It simply converts message to a stream and expect it to be read.
 */
public class GenericCANDriverTest {

    @Test
    void testConnection() throws PlcConnectionException {
        PlcConnection connection = new DefaultPlcDriverManager().getConnection("genericcan:virtualcan://");

        assertNotNull(connection);
        assertTrue(connection.isConnected());
        assertFalse(connection.getMetadata().isReadSupported());
        assertTrue(connection.getMetadata().isWriteSupported());
        assertTrue(connection.getMetadata().isSubscribeSupported());
    }

    @Test
    void testSubscribeAndWrite() throws Exception {
        Map<String, Entry<String, Object>> tagMap = Map.ofEntries(
            entry("tag1", entry("200:BYTE", (short) 0x0A)),
            entry("tag2", entry("200:INTEGER8", (byte) 20)),
            entry("tag3", entry("200:INTEGER8", (byte) 30))
        );

        subscribeAndWrite(tagMap);
    }

    @Test
    void testSubscribeAndWriteRawArray() throws Exception {
        Map<String, Entry<String, Object>> tagMap = Map.ofEntries(
            entry("arr1", entry("201:RAW", new byte[] {
                (short) 0, (short) 1, (short) 2, (short) 3, (short) 4, (short) 5, (short) 6, (short) 7}
            ))
        );

        subscribeAndWrite(tagMap);
    }

    @Test
    @Disabled("Writing arrays requires use of RAW type")
    void testSubscribeAndWriteByteArray() throws Exception {
        Map<String, Entry<String, Object>> tagMap = Map.ofEntries(
            entry("arr1", entry("201:BYTE[8]", new byte[] {
                (short) 0, (short) 1, (short) 2, (short) 3, (short) 4, (short) 5, (short) 6, (short) 7}
            ))
        );

        subscribeAndWrite(tagMap);
    }


    void subscribeAndWrite(Map<String, Entry<String, Object>> entries) throws Exception {
        PlcConnection connection = new DefaultPlcDriverManager().getConnection("genericcan:virtualcan://");

        Plc4xEmbeddedChannel subscribeChannel = null;
        Plc4xEmbeddedChannel writeChannel = null;
        if (connection instanceof ChannelExposingConnection) {
            Channel channel = ((ChannelExposingConnection) connection).getChannel();
            if (channel instanceof Plc4xEmbeddedChannel) {
                subscribeChannel = (Plc4xEmbeddedChannel) channel;
                writeChannel = (Plc4xEmbeddedChannel) channel;
            }
        }
        if (subscribeChannel == null) {
            throw new IllegalArgumentException("Invalid configuration");
        }

        CountDownLatch latch = new CountDownLatch(1);

        final AtomicReference<PlcSubscriptionEvent> plcEvent = new AtomicReference<>();
        Builder subscriptionRequestBuilder = connection.subscriptionRequestBuilder();
        entries.forEach((k, v) -> {
            subscriptionRequestBuilder.addEventTagAddress(k, v.getKey());
        });

        PlcSubscriptionRequest subscriptionRequest = subscriptionRequestBuilder.build();
        subscriptionRequest.execute().whenComplete((reply, error) -> {
            if (error != null) {
                fail(error);
                return;
            }

            reply.getSubscriptionHandle(entries.keySet().iterator().next()).register(event -> {
                plcEvent.set(event);
                latch.countDown();
            });
        });

        PlcWriteRequest.Builder writeRequestBuilder = connection.writeRequestBuilder();
        entries.forEach((k, v) -> {
            writeRequestBuilder.addTagAddress(k, v.getKey(), v.getValue());
        });
        PlcWriteRequest writeRequest = writeRequestBuilder.build();
        writeRequest.execute().whenComplete((reply, error) -> {
            if (error != null) {
                fail(error);
            }
        }).get();

        // copy outbound message to inbound queue to confirm that transport API works and subscription
        // is properly matched against incoming message
        ByteBuf outgoing = writeChannel.flushOutbound().readOutbound();
        subscribeChannel.writeInbound(outgoing);
        subscribeChannel.flushInbound();

        latch.await();

        PlcSubscriptionEvent event = plcEvent.get();
        entries.forEach((k, v) -> {
            Object object = event.getObject(k);
            if (!v.getValue().getClass().isArray()) {
                assertEquals(v.getValue(), object);
                return;
            }

            // comparing arrays is a bit of nightmare due to primitives
            int length = Array.getLength(v.getValue());
            int readLength = Array.getLength(object);
            if (readLength != length) {
                throw new IllegalArgumentException("Return value length do not match reference value");
            }
            for (int index = 0; index < length; index++) {
                assertEquals(Array.get(v.getValue(), index), Array.get(object, index));
            }
        });

    }
}
