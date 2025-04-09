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
package org.apache.plc4x.java.opm.issues.i1935;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.connection.MockDevice;
import org.apache.plc4x.java.opm.PlcEntityManager;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;

public class Plc4jMockConnectionTest {

    @Test
    @Disabled("This is only used for manual execution")
    void shouldRead() throws PlcConnectionException {
        MockDevice mockDevice = Mockito.mock(MockDevice.class);
        DefaultPlcDriverManager driverManager = new DefaultPlcDriverManager();
        MockConnection connection = (MockConnection) driverManager.getConnection("mock:test");
        when(mockDevice.read(any())).thenAnswer(invocation -> {
            System.out.println("Reading: " + invocation.getArgument(0));
            return new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcSTRING("1")) {};
        });
        connection.setDevice(mockDevice);

        PlcEntityManager entityManager = new PlcEntityManager(driverManager);

        try {
            FooRxEntity entity = entityManager.connect(FooRxEntity.class, "mock:test");
            while (true) {
                try {
                    entity.updateAllTheTags();
                    System.out.println(entity.getBf102Setpoint());
                    System.out.println(entity.getBf112());
                    System.out.println(entity.getBp112());
                    Thread.sleep(10);
                } catch (InterruptedException e) {
                    throw new RuntimeException(e);
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
