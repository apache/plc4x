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
package org.apache.plc4x.java.s7.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.PooledByteBufAllocator;
import io.netty.channel.embedded.EmbeddedChannel;
import io.netty.util.ResourceLeakDetector;
import org.apache.plc4x.java.s7.readwrite.configuration.S7Configuration;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

/**
 * Regression test for issue #2248.
 *
 * <p>S7HMuxImpl is installed in the EmbeddedChannel pipeline of S7HPlcConnection. Its encode()
 * forwards each outbound message to the active TCP channel and propagates a sentinel downstream.
 * Before the fix, encode() pushed a freshly-copied direct ByteBuf into the EmbeddedChannel's
 * outbound queue: nothing in plc4j ever drained that queue, so every sent message leaked one
 * direct buffer until DirectMemory was exhausted (OutOfMemoryError reported on long-running
 * applications).
 */
public class S7HMuxLeakTest {

    private static final int ITERATIONS = 2000;
    private static final int PAYLOAD_BYTES = 256;

    @Test
    void embeddedOutboundQueueShouldNotAccumulateRealByteBufs() throws Exception {
        ResourceLeakDetector.Level previous = ResourceLeakDetector.getLevel();
        ResourceLeakDetector.setLevel(ResourceLeakDetector.Level.PARANOID);
        try {
            S7HMuxImpl mux = new S7HMuxImpl();
            EmbeddedChannel embedded = new EmbeddedChannel(mux);
            EmbeddedChannel tcp = new EmbeddedChannel();

            mux.setEmbededhannel(embedded, new S7Configuration());
            mux.setPrimaryChannel(tcp);

            int forwardedToTcp = 0;
            for (int i = 0; i < ITERATIONS; i++) {
                ByteBuf out = PooledByteBufAllocator.DEFAULT.directBuffer(PAYLOAD_BYTES);
                out.writeZero(PAYLOAD_BYTES);
                embedded.writeOutbound(out);

                Object msg;
                while ((msg = tcp.readOutbound()) != null) {
                    if (msg instanceof ByteBuf) {
                        ((ByteBuf) msg).release();
                        forwardedToTcp++;
                    }
                }
            }

            assertEquals(ITERATIONS, forwardedToTcp,
                "every message written to the embedded channel must be forwarded once to the TCP channel");

            int leakedBuffers = 0;
            long leakedBytes = 0;
            Object pending;
            while ((pending = embedded.readOutbound()) != null) {
                if (pending instanceof ByteBuf) {
                    ByteBuf bb = (ByteBuf) pending;
                    if (bb.capacity() > 0) {
                        leakedBuffers++;
                        leakedBytes += bb.capacity();
                    }
                    bb.release();
                }
            }
            assertEquals(0, leakedBuffers,
                "EmbeddedChannel outbound queue must not accumulate non-empty ByteBufs (leaked "
                    + leakedBytes + " bytes across " + leakedBuffers + " buffers)");

            embedded.finishAndReleaseAll();
            tcp.finishAndReleaseAll();
        } finally {
            ResourceLeakDetector.setLevel(previous);
        }
    }
}
