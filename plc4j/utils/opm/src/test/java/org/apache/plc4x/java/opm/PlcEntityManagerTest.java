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
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.items.DefaultStringFieldItem;
import org.apache.plc4x.java.mock.MockDevice;
import org.apache.plc4x.java.mock.PlcMockConnection;
import org.junit.Test;
import org.mockito.Mockito;
import org.mockito.invocation.InvocationOnMock;
import org.mockito.stubbing.Answer;

import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedQueue;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.*;

public class PlcEntityManagerTest {

    @Test
    public void read_throwsInvalidFieldException_rethrows() throws PlcConnectionException {
        // Prepare the Mock
        PlcDriverManager driverManager = Mockito.mock(PlcDriverManager.class);
        PlcConnection connection = Mockito.mock(PlcConnection.class);
        PlcConnectionMetadata metadata = Mockito.mock(PlcConnectionMetadata.class);
        PlcReadRequest.Builder builder = Mockito.mock(PlcReadRequest.Builder.class);
        when(metadata.canRead()).thenReturn(true);
        when(connection.readRequestBuilder()).thenReturn(builder);
        when(connection.getMetadata()).thenReturn(metadata);
        when(builder.build()).thenThrow(new PlcInvalidFieldException("field1"));
        when(driverManager.getConnection(any())).thenReturn(connection);

        // Create Entity Manager
        PlcEntityManager entityManager = new PlcEntityManager(driverManager);

        // Issue Call to trigger interception
        String message = null;
        try {
            BadEntity entity = entityManager.read(BadEntity.class, "mock:test");
        } catch (Exception e) {
            message = e.getMessage();
        }

        assertEquals("Unable to parse field 'field1'", message);
    }

    @Test
    public void read_unableToConnect_rethrows() throws PlcConnectionException {
        // Prepare the Mock
        PlcDriverManager driverManager = Mockito.mock(PlcDriverManager.class);
        when(driverManager.getConnection(any())).thenThrow(new PlcConnectionException(""));

        // Create Entity Manager
        PlcEntityManager entityManager = new PlcEntityManager(driverManager);

        // Issue Call to trigger interception
        String message = null;
        try {
            BadEntity entity = entityManager.read(BadEntity.class, "mock:test");
        } catch (Exception e) {
            message = e.getMessage();
        }

        assertEquals("Unable to get connection with url 'mock:test'", message);
    }

    @Test(expected = OPMException.class)
    public void read_timeoutOnGet_throwsException() throws PlcConnectionException, OPMException {
        // Prepare the Mock
        MockDevice mockDevice = Mockito.mock(MockDevice.class);
        PlcDriverManager driverManager = new PlcDriverManager();
        PlcMockConnection connection = (PlcMockConnection) driverManager.getConnection("mock:test");
        when(mockDevice.read(any())).thenAnswer(new Answer<Object>() {
            @Override
            public Object answer(InvocationOnMock invocation) throws Throwable {
                // Sleep for 3s
                Thread.sleep(3_000);
                return Pair.of(PlcResponseCode.OK, new DefaultStringFieldItem("Hallo"));
            }
        });
        connection.setDevice(mockDevice);

        // Create Entity Manager
        PlcEntityManager entityManager = new PlcEntityManager(driverManager);

        // Issue Call which SHOULD timeout
        BadEntity entity = entityManager.read(BadEntity.class, "mock:test");
    }

    @Test(expected = OPMException.class)
    public void read_uninstantiableEntity_throws() throws OPMException {
        PlcEntityManager entityManager = new PlcEntityManager();

        UninstantiableEntity entity = entityManager.read(UninstantiableEntity.class, "mock:test");
    }

    /**
     * Class is private, so EntityManager has no access to it
     * @throws OPMException
     */
    @Test(expected = OPMException.class)
    public void connect_uninstantiableEntity_throws() throws OPMException {
        PlcEntityManager entityManager = new PlcEntityManager();

        UninstantiableEntity entity = entityManager.connect(UninstantiableEntity.class, "mock:test");
    }

    @Test
    public void read_resolveAlias_works() throws OPMException, PlcConnectionException {
        SimpleAliasRegistry registry = new SimpleAliasRegistry();
        registry.register("alias", "real_field");

        // Mock
        PlcDriverManager driverManager = new PlcDriverManager();
        PlcMockConnection connection = (PlcMockConnection) driverManager.getConnection("mock:test");
        MockDevice mockDevice = Mockito.mock(MockDevice.class);
        when(mockDevice.read(any())).thenReturn(Pair.of(PlcResponseCode.OK, new DefaultStringFieldItem("value")));
        connection.setDevice(mockDevice);

        PlcEntityManager entityManager = new PlcEntityManager(driverManager, registry);
        entityManager.read(AliasEntity.class, "mock:test");

        // Assert that "field" was queried
        verify(mockDevice).read(eq("real_field"));
    }

    @Test
    public void connect_resolveAlias_works() throws PlcConnectionException, OPMException {
        SimpleAliasRegistry registry = new SimpleAliasRegistry();
        registry.register("alias", "real_field");

        // Mock
        PlcDriverManager driverManager = new PlcDriverManager();
        PlcMockConnection connection = (PlcMockConnection) driverManager.getConnection("mock:test");
        MockDevice mockDevice = Mockito.mock(MockDevice.class);
        when(mockDevice.read(any())).thenReturn(Pair.of(PlcResponseCode.OK, new DefaultStringFieldItem("value")));
        connection.setDevice(mockDevice);

        PlcEntityManager entityManager = new PlcEntityManager(driverManager, registry);
        entityManager.connect(AliasEntity.class, "mock:test");

        // Assert that "field" was queried
        verify(mockDevice, times(1)).read(eq("real_field"));
    }

    @Test(expected = OPMException.class)
    public void read_unknownAlias_throws() throws OPMException {
        PlcEntityManager entityManager = new PlcEntityManager();

        entityManager.read(AliasEntity.class, "mock:test");
    }

    @Test
    public void read_badAlias_throws() {
        PlcEntityManager entityManager = new PlcEntityManager();

        String message = null;
        try {
            entityManager.read(BadAliasEntity.class, "mock:test");
        } catch (OPMException e) {
            message = e.getMessage();
        }

        assertNotNull(message);
        assertTrue(message.contains("Invalid Syntax, either use field address (no starting $) or an alias with Syntax ${xxx}. But given was"));
    }

    @PlcEntity
    private static class UninstantiableEntity {

        public UninstantiableEntity() throws InstantiationException {
            throw new InstantiationException("Do not instantiate");
        }

    }

    @PlcEntity
    public static class BadEntity {

        @PlcField("field1")
        private String field1;

        public BadEntity() {
            // for OPM
        }

        public String getField1() {
            return field1;
        }
    }

    @PlcEntity
    public static class AliasEntity {

        @PlcField("${alias}")
        private String aliasedField;

        public AliasEntity() {
            // for OPM
        }

        public String getAliasedField() {
            return aliasedField;
        }
    }

    @PlcEntity
    public static class BadAliasEntity {

        @PlcField("${alias")
        private String aliasedField;

        public BadAliasEntity() {
            // for OPM
        }

        public String getAliasedField() {
            return aliasedField;
        }
    }

}