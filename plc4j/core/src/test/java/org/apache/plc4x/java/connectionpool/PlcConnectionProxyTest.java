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
package org.apache.plc4x.java.connectionpool;

import org.apache.plc4x.java.api.connection.PlcConnection;
import org.junit.After;
import org.junit.Test;
import org.mockito.Mockito;


public class PlcConnectionProxyTest {

    private final String TEST_STRING = "s7:127.0.0.100";
    private PlcConnectionManager connectionManager = Mockito.mock(PlcConnectionManager.class);
    private PlcConnection plcConnection = Mockito.mock(PlcConnection.class);
    private PlcConnectionProxy plcConnectionProxy = new PlcConnectionProxy(connectionManager, TEST_STRING, plcConnection);

    @After
    public void shutDown() throws Exception {
        plcConnection.close();
    }

    @Test
    public void connect() throws Exception {
        plcConnectionProxy.connect();
        Mockito.verify(plcConnection, Mockito.times(0)).connect();
    }

    @Test
    public void isConnected() {
        plcConnectionProxy.isConnected();
        Mockito.verify(plcConnection, Mockito.times(1)).isConnected();
    }

    @Test
    public void close() {
        plcConnectionProxy.close();
        Mockito.verify(connectionManager, Mockito.times(1)).returnConnection(TEST_STRING);
    }

    @Test
    public void getReader() {
        plcConnectionProxy.getReader();
        Mockito.verify(plcConnection, Mockito.times(1)).getReader();
    }

    @Test
    public void getWriter() {
        plcConnectionProxy.getWriter();
        Mockito.verify(plcConnection, Mockito.times(1)).getWriter();
    }

    @Test
    public void getSubscriber() {
        plcConnectionProxy.getSubscriber();
        Mockito.verify(plcConnection, Mockito.times(1)).getSubscriber();
    }
}