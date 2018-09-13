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

package org.apache.plc4x.java.ads.connection;

import io.netty.channel.Channel;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import static org.junit.Assert.assertEquals;
import static org.mockito.Mockito.RETURNS_DEEP_STUBS;
import static org.mockito.Mockito.mock;

public class AdsTcpPlcConnectionTests {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsTcpPlcConnectionTests.class);

    private AdsTcpPlcConnection SUT;

    private Channel channelMock;

    private ExecutorService executorService;

    @Before
    public void setUp() throws Exception {
        SUT = AdsTcpPlcConnection.of(InetAddress.getByName("localhost"), AmsNetId.of("0.0.0.0.0.0"), AmsPort.of(13));
        // TODO: Refactor this to use the TestChannelFactory instead.
        channelMock = mock(Channel.class, RETURNS_DEEP_STUBS);
        FieldUtils.writeField(SUT, "channel", channelMock, true);
        executorService = Executors.newFixedThreadPool(10);
    }

    @After
    public void tearDown() {
        executorService.shutdownNow();
        SUT = null;
    }

    @Test
    public void initialState() {
        assertEquals(SUT.getTargetAmsNetId().toString(), "0.0.0.0.0.0");
        assertEquals(SUT.getTargetAmsPort().toString(), "13");
    }

    @Test
    public void implementMeTestNewAndMissingMethods() {
        // TODO: implement me
    }
}