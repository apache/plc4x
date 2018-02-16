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
package org.apache.plc4x.java.s7.connection;

import io.netty.channel.Channel;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.junit.After;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.Timeout;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import static org.hamcrest.core.Is.is;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;

public class S7PlcConnectionIT {

    private static final Logger logger = LoggerFactory.getLogger(S7PlcConnectionIT.class);

    @Rule
    public Timeout globalTimeout = Timeout.seconds(2); // 2 seconds max per method tested

    private S7PlcTestConnection  s7PlcConnection;
    private Channel channel;

    @Before
    public void setUp() {
        try {
            s7PlcConnection = new S7PlcTestConnection(1, 2, "");
            s7PlcConnection.connect();
            channel = s7PlcConnection.getChannel();
        } catch (PlcConnectionException e) {
            logger.error("Error initializing connection", e);
            fail("Error initializing connection");
        }
    }

    @After
    public void tearDown() {
        if(s7PlcConnection.isConnected()) {
            s7PlcConnection.close();
        }
        s7PlcConnection = null;
        channel = null;
    }

    @Test
    public void connect() {
        assertThat(s7PlcConnection, notNullValue());
        assertThat("The connection should be 'connected'", s7PlcConnection.isConnected(), is( true) );
    }

    // TODO more tests for connect, close, read and write

}
