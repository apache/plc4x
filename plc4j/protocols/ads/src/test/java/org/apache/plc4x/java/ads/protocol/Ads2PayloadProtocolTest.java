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
package org.apache.plc4x.java.ads.protocol;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.generic.AmsPacket;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Collection;
import java.util.stream.Collectors;

import static org.hamcrest.Matchers.*;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertThat;

@RunWith(Parameterized.class)
public class Ads2PayloadProtocolTest extends AbstractProtocolTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(Ads2PayloadProtocolTest.class);

    private Ads2PayloadProtocol SUT;

    @Parameterized.Parameter
    public AmsPacket amsPacket;

    @Parameterized.Parameter(1)
    public String clazzName;

    @Parameterized.Parameters(name = "{index} {1}")
    public static Collection<Object[]> data() {
        return amsPacketStream()
            .map(amsPacket -> new Object[]{
                amsPacket,
                amsPacket.getClass().getSimpleName()
            })
            .collect(Collectors.toList());
    }

    @Before
    public void setUp() throws Exception {
        SUT = new Ads2PayloadProtocol();
        byte[] bytes = amsPacket.getBytes();
        LOGGER.info("amsPacket:\n{} has \n{}bytes\nHexDump:\n{}", amsPacket, bytes.length, amsPacket.dump());
    }

    @Test
    public void encode() throws Exception {
        ArrayList<Object> out = new ArrayList<>();
        SUT.encode(null, amsPacket, out);
        assertEquals(1, out.size());
        assertThat(out, hasSize(1));
    }

    @Test
    public void decode() throws Exception {
        ArrayList<Object> out = new ArrayList<>();
        SUT.decode(null, amsPacket.getByteBuf(), out);
        assertThat(out, hasSize(1));
    }

    @Test
    public void roundTrip() throws Exception {
        ArrayList<Object> outbound = new ArrayList<>();
        SUT.encode(null, amsPacket, outbound);
        assertEquals(1, outbound.size());
        assertThat(outbound, hasSize(1));
        assertThat(outbound.get(0), instanceOf(ByteBuf.class));
        ByteBuf byteBuf = (ByteBuf) outbound.get(0);
        ArrayList<Object> inbound = new ArrayList<>();
        SUT.decode(null, byteBuf, inbound);
        assertEquals(1, inbound.size());
        assertThat(inbound, hasSize(1));
        assertThat(inbound.get(0), instanceOf(AmsPacket.class));
        AmsPacket inboundAmsPacket = (AmsPacket) inbound.get(0);
        assertThat("inbound divers from outbound", this.amsPacket, equalTo(inboundAmsPacket));
    }
}