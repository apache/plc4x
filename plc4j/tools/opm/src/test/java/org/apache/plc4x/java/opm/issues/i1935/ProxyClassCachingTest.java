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
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.connection.MockDevice;
import org.apache.plc4x.java.opm.OPMException;
import org.apache.plc4x.java.opm.PlcEntity;
import org.apache.plc4x.java.opm.PlcEntityManager;
import org.apache.plc4x.java.opm.PlcTag;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.HashSet;
import java.util.Set;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;

/**
 * The generated proxy type depends only on the entity class, so it has to be generated once and
 * reused - see GH-1935. Generating one per call leaks a class per read, which is what filled up
 * the reporter's heap.
 */
@ExtendWith(MockitoExtension.class)
public class ProxyClassCachingTest {

    private PlcEntityManager entityManager;

    @Mock
    MockDevice mockDevice;

    @BeforeEach
    void setUp() throws Exception {
        DefaultPlcDriverManager driverManager = new DefaultPlcDriverManager();
        MockConnection connection = (MockConnection) driverManager.getConnection("mock:test");
        when(mockDevice.read(any()))
            .thenReturn(new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcSTRING("hallo")));
        connection.setDevice(mockDevice);
        entityManager = new PlcEntityManager(driverManager);
    }

    @Test
    void repeatedReadsReuseTheSameProxyClass() throws OPMException {
        Set<Class<?>> proxyClasses = new HashSet<>();
        for (int i = 0; i < 25; i++) {
            proxyClasses.add(entityManager.read(SimpleEntity.class, "mock:test").getClass());
        }

        assertEquals(1, proxyClasses.size(),
            "every read generated a new proxy class - the type cache is not being used");
    }

    @Test
    void connectAndReadShareTheSameProxyClass() throws OPMException {
        Object connected = entityManager.connect(SimpleEntity.class, "mock:test");
        Object read = entityManager.read(SimpleEntity.class, "mock:test");

        assertSame(connected.getClass(), read.getClass());
    }

    @Test
    void separateEntityManagersShareTheSameProxyClass() throws Exception {
        DefaultPlcDriverManager otherDriverManager = new DefaultPlcDriverManager();
        MockConnection otherConnection = (MockConnection) otherDriverManager.getConnection("mock:other");
        otherConnection.setDevice(mockDevice);
        PlcEntityManager otherEntityManager = new PlcEntityManager(otherDriverManager);

        Object first = entityManager.read(SimpleEntity.class, "mock:test");
        Object second = otherEntityManager.read(SimpleEntity.class, "mock:other");

        assertSame(first.getClass(), second.getClass());
    }

    @Test
    void differentEntitiesGetDifferentProxyClasses() throws OPMException {
        Object simple = entityManager.read(SimpleEntity.class, "mock:test");
        Object other = entityManager.read(OtherEntity.class, "mock:test");

        assertEquals(2, new HashSet<>(Set.of(simple.getClass(), other.getClass())).size());
    }

    @PlcEntity
    public static class SimpleEntity {

        @PlcTag("address")
        private String field;

        public SimpleEntity() {
            // For OPM
        }

        public String getField() {
            return field;
        }
    }

    @PlcEntity
    public static class OtherEntity {

        @PlcTag("address")
        private String field;

        public OtherEntity() {
            // For OPM
        }

        public String getField() {
            return field;
        }
    }
}
