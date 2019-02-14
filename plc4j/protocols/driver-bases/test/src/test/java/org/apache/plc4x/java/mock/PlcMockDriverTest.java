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

package org.apache.plc4x.java.mock;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.items.DefaultIntegerFieldItem;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.concurrent.ExecutionException;

import static org.junit.Assert.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;


public class PlcMockDriverTest implements WithAssertions {

    private final PlcMockDriver driver = new PlcMockDriver();

    @Test
    public void fetchTwoConnections_areEqual() throws PlcConnectionException {
        PlcConnection conn1 = driver.connect("mock:123");
        PlcConnection conn2 = driver.connect("mock:123");

        assertEquals(conn1, conn2);
    }

    /**
     * Example of Mock Usage.
     * The idea is to fetch a mock connection with a specific name and prepare a {@link MockDevice} which is set there.
     * <p>
     * Some application code which uses the same Connection String will then automatically get the same connection
     * and operate against the same {@link MockDevice} without the necessity to also mock field queries or other things.
     * <p>
     * In this example the {@link #someCodeWhichShouldDoPlcManipulation(String)} function represents the Business Logic
     * which should be tested and where only the connection string is manipulated for the test.
     */
    @Test
    public void testScenarioExample() throws PlcConnectionException, ExecutionException, InterruptedException {
        PlcMockConnection preparingConnection = ((PlcMockConnection) driver.connect("test:123"));
        MockDevice mock = Mockito.mock(MockDevice.class);
        when(mock.read(any())).thenReturn(Pair.of(PlcResponseCode.OK, new DefaultIntegerFieldItem(1)));
        preparingConnection.setDevice(mock);

        // Now we can simply inject this URL into our code and automatically have our mock
        someCodeWhichShouldDoPlcManipulation("test:123");

        // Verify that the code did indeed what we wanted it to do
        verify(mock, times(1)).read("DB2.DBD17:INT");
    }

    @Test
    public void wrongDevice() {
        assertThatThrownBy(() -> driver.connect("mock:"))
            .isInstanceOf(PlcConnectionException.class);
    }

    /**
     * Example function that does some reading from a siemens plc using Siemens Syntax
     *
     * @param connection Connection String, e.g., from config
     */
    private void someCodeWhichShouldDoPlcManipulation(String connection) throws PlcConnectionException, ExecutionException, InterruptedException {
        // Normally this would be from the driver manager
        PlcConnection connect = driver.connect(connection);
        PlcReadResponse response = connect.readRequestBuilder()
            .addItem("value", "DB2.DBD17:INT")
            .build()
            .execute()
            .get();
        // Normally do something with the response
    }
}