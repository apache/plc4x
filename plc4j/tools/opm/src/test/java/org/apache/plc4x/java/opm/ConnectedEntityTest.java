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

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.connection.MockDevice;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.stream.IntStream;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

/**
 * Tests for Connected Entities.
 */
@ExtendWith(MockitoExtension.class)
public class ConnectedEntityTest {

    PlcDriverManager driverManager;

    MockConnection connection;

    PlcEntityManager entityManager;

    @Mock
    MockDevice mockDevice;

    @BeforeEach
    void setUp() throws Exception {
        driverManager = new PlcDriverManager();
        connection = (MockConnection) driverManager.getConnection("mock:cached");
        when(mockDevice.read(any()))
            .thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcSTRING("hallo")));
        connection.setDevice(mockDevice);
        entityManager = new PlcEntityManager(driverManager);
    }

    @Test
    void useCache() throws OPMException {
        // Trigger a fetch
        CachingEntity entity = entityManager.connect(CachingEntity.class, "mock:cached");
        // Trigger second fetch
        assertEquals("hallo", entity.getField());

        verify(mockDevice, timeout(1_000).times(1)).read(any());
    }

    @Test
    void useCache_timeout_refetches() throws OPMException, InterruptedException {
        // Trigger a fetch
        CachingEntity entity = entityManager.connect(CachingEntity.class, "mock:cached");
        Thread.sleep(500);
        // Trigger second fetch
        assertEquals("hallo", entity.getField());

        verify(mockDevice, timeout(1_000).times(2)).read(any());
    }

    @Test
    void cache_manyRequests_onlyOneToPlc() throws OPMException {
        // Mock
        when(mockDevice.write(any(), any()))
            .thenReturn(PlcResponseCode.OK);

        // Trigger a fetch
        CachingEntity entity = entityManager.connect(CachingEntity.class, "mock:cached");
        // Trigger Many Fetches via getter
        IntStream.range(1, 10).forEach(i -> entity.getField());
        IntStream.range(1, 10).forEach(i -> entity.dummyMethod());

        verify(mockDevice, timeout(1_000).times(1)).read(any());
    }

    @PlcEntity
    public static class CachingEntity {

        @PlcField(value = "address", cacheDurationMillis = 500)
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