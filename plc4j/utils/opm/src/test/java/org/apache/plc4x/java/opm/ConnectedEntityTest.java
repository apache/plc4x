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

package org.apache.plc4x.java.opm;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.items.DefaultStringFieldItem;
import org.apache.plc4x.java.mock.MockDevice;
import org.apache.plc4x.java.mock.PlcMockConnection;
import org.junit.Test;
import org.mockito.Mockito;

import java.util.stream.IntStream;

import static org.junit.Assert.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

/**
 * Tests for Connected Entities.
 */
public class ConnectedEntityTest {

    @Test
    public void useCache() throws PlcConnectionException, OPMException {
        // Mock
        PlcDriverManager driverManager = new PlcDriverManager();
        PlcMockConnection connection = (PlcMockConnection) driverManager.getConnection("mock:cached");
        MockDevice mock = Mockito.mock(MockDevice.class);
        when(mock.read(any())).thenReturn(Pair.of(PlcResponseCode.OK, new DefaultStringFieldItem("hallo")));
        connection.setDevice(mock);
        PlcEntityManager entityManager = new PlcEntityManager(driverManager);

        // Trigger a fetch
        CachingEntity entity = entityManager.connect(CachingEntity.class, "mock:cached");
        // Trigger second fetch
        assertEquals("hallo", entity.getField());

        verify(mock, timeout(1_000).times(1)).read(any());
    }

    @Test
    public void useCache_timeout_refetches() throws PlcConnectionException, OPMException, InterruptedException {
        // Mock
        PlcDriverManager driverManager = new PlcDriverManager();
        PlcMockConnection connection = (PlcMockConnection) driverManager.getConnection("mock:cached");
        MockDevice mock = Mockito.mock(MockDevice.class);
        when(mock.read(any())).thenReturn(Pair.of(PlcResponseCode.OK, new DefaultStringFieldItem("hallo")));
        connection.setDevice(mock);
        PlcEntityManager entityManager = new PlcEntityManager(driverManager);

        // Trigger a fetch
        CachingEntity entity = entityManager.connect(CachingEntity.class, "mock:cached");
        Thread.sleep(500);
        // Trigger second fetch
        assertEquals("hallo", entity.getField());

        verify(mock, timeout(1_000).times(2)).read(any());
    }

    @Test
    public void cache_manyRequests_onlyOneToPlc() throws PlcConnectionException, OPMException {
        // Mock
        PlcDriverManager driverManager = new PlcDriverManager();
        PlcMockConnection connection = (PlcMockConnection) driverManager.getConnection("mock:cached");
        MockDevice mock = Mockito.mock(MockDevice.class);
        when(mock.read(any())).thenReturn(Pair.of(PlcResponseCode.OK, new DefaultStringFieldItem("hallo")));
        connection.setDevice(mock);
        PlcEntityManager entityManager = new PlcEntityManager(driverManager);

        // Trigger a fetch
        CachingEntity entity = entityManager.connect(CachingEntity.class, "mock:cached");
        // Trigger Many Fetches via getter
        IntStream.range(1,100).forEach(i -> entity.getField());
        IntStream.range(1,100).forEach(i -> entity.dummyMethod());

        verify(mock, timeout(1_000).times(1)).read(any());
    }

    @PlcEntity
    public static class CachingEntity {

        @PlcField(value = "address", cacheDurationMillis = 100)
        private String field;

        public CachingEntity() {
            // For OPM
        }

        public String getField() {
            return field;
        }

        public void dummyMethod() {
            // do nothing
        }
    }
}