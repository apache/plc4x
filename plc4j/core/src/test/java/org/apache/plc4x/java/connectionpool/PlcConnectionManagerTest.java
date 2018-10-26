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
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.exceptions.NoConnectionAvailableException;
import org.apache.plc4x.java.exceptions.NotConnectedException;
import org.apache.plc4x.test.FastTests;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import org.mockito.Mockito;

import java.util.HashMap;
import java.util.Map;


public class PlcConnectionManagerTest {

    private PlcConnection plcConnection = Mockito.mock(PlcConnection.class);
    private PlcConnectionManager plcConnectionManager;
    private static final String TEST_STRING = "s7:127.0.0.100";


    @Before
    public void setUp() {
        Map<String, PlcConnection> plcConnectionProxyMap = new HashMap<>();
        Map<String, Boolean> booleanMap = new HashMap<>();
        plcConnectionProxyMap.put(TEST_STRING, plcConnection);
        booleanMap.put(TEST_STRING, false);
        plcConnectionManager = new PlcConnectionManager(plcConnectionProxyMap, booleanMap);
    }

    @After
    public void shutDown() {
        plcConnectionManager.close();
    }

    @Test
    @Category(FastTests.class)
    public void getConnection() throws Exception {
        plcConnectionManager.getConnection(TEST_STRING);
        Mockito.verify(plcConnection, Mockito.times(1)).isConnected();
    }

    @Test(expected = NoConnectionAvailableException.class)
    @Category(FastTests.class)
    public void returnConnectionThrows() throws Exception {
        Mockito.when(plcConnection.isConnected()).thenReturn(true);
        plcConnectionManager.getConnection(TEST_STRING);
        //second attempt should throw exception because connection was not returned yet
        plcConnectionManager.getConnection(TEST_STRING);
    }


    @Test(expected = NotConnectedException.class)
    @Category(FastTests.class)
    public void returnConnectionCannotConnect() throws Exception {
        Mockito.when(plcConnection.isConnected()).thenReturn(false);
        PlcConnection plcConnectionProxy = plcConnectionManager.getConnection(TEST_STRING); // should
        plcConnectionProxy.close();
        Mockito.doThrow(new PlcConnectionException("")).when(plcConnection).connect();
        plcConnectionManager.getConnection(TEST_STRING);

    }

    @Test
    @Category(FastTests.class)
    public void returnConnection() throws Exception {
        Mockito.when(plcConnection.isConnected()).thenReturn(true);
        plcConnectionManager.getConnection(TEST_STRING);
        plcConnectionManager.returnConnection(TEST_STRING);
        plcConnectionManager.getConnection(TEST_STRING);
    }

    @Test
    @Category(FastTests.class)
    public void closeConnections() {
        plcConnectionManager.close();
    }
}