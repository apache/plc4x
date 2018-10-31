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

package org.apache.plc4x.java.mock;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.items.DefaultLongFieldItem;
import org.junit.Test;
import org.mockito.Mockito;

import static junit.framework.TestCase.assertFalse;
import static junit.framework.TestCase.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

public class PlcMockDriverTest {

    private final PlcDriverManager driverManager = new PlcDriverManager();

    @Test
    public void useMockDriver_noDevice_isNotConnected() throws Exception {
        PlcConnection connection = driverManager.getConnection("mock:dummy");

        assertFalse(connection.isConnected());
    }

    @Test
    public void useMockDriver_deviceSet_isConnected() throws Exception {
        PlcMockConnection connection = (PlcMockConnection) driverManager.getConnection("mock:dummy");
        MockDevice mock = Mockito.mock(MockDevice.class);
        connection.setDevice(mock);

        assertTrue(connection.isConnected());
    }

    @Test
    public void mockDriver_assertSimpleRequest() throws PlcConnectionException {
        PlcMockConnection connection = (PlcMockConnection) driverManager.getConnection("mock:dummy");
        MockDevice mock = Mockito.mock(MockDevice.class);
        when(mock.read(any())).thenReturn(Pair.of(PlcResponseCode.OK, new DefaultLongFieldItem(1L)));
        connection.setDevice(mock);

        connection.readRequestBuilder()
            .addItem("item1", "myPlcField")
            .build()
            .execute();

        // Verify the call
        verify(mock, times(1)).read(eq("myPlcField"));
    }
}