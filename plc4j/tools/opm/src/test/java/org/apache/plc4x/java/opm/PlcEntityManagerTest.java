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
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.connection.MockDevice;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Answers;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
public class PlcEntityManagerTest implements WithAssertions {

    @Nested
    class Read {

        @Mock(answer = Answers.RETURNS_DEEP_STUBS)
        PlcDriverManager driverManager;

        @Test
        public void throwsInvalidFieldException_rethrows() throws PlcConnectionException {
            // Prepare the Mock
            when(driverManager.getConnection(any()).readRequestBuilder().build())
                .thenThrow(new PlcInvalidFieldException("field1"));

            // Create Entity Manager
            PlcEntityManager entityManager = new PlcEntityManager(driverManager);

            // Issue Call to trigger interception
            assertThatThrownBy(() -> entityManager.read(BadEntity.class, "mock:test"))
                .hasCauseInstanceOf(PlcInvalidFieldException.class)
                .hasStackTraceContaining("field1 invalid");
        }

        @Test
        public void unableToConnect_rethrows() throws PlcConnectionException {
            // Prepare the Mock
            when(driverManager.getConnection(any()))
                .thenThrow(new PlcConnectionException(""));

            // Create Entity Manager
            PlcEntityManager entityManager = new PlcEntityManager(driverManager);

            // Issue Call to trigger interception
            assertThatThrownBy(() -> entityManager.read(BadEntity.class, "mock:test"))
                .hasCauseInstanceOf(PlcConnectionException.class)
                .hasStackTraceContaining("Problem during processing");
        }

        @Test
        public void timeoutOnGet_throwsException() throws PlcConnectionException {
            // Prepare the Mock
            MockDevice mockDevice = Mockito.mock(MockDevice.class);
            PlcDriverManager driverManager = new PlcDriverManager();
            MockConnection connection = (MockConnection) driverManager.getConnection("mock:test");
            when(mockDevice.read(any())).thenAnswer(invocation -> {
                // Sleep for 3s
                Thread.sleep(3_000);
                return Pair.of(PlcResponseCode.OK, new PlcSTRING("Hallo"));
            });
            connection.setDevice(mockDevice);

            // Create Entity Manager
            PlcEntityManager entityManager = new PlcEntityManager(driverManager);

            // Issue Call which SHOULD timeout
            assertThatThrownBy(() -> entityManager.read(BadEntity.class, "mock:test"))
                .isInstanceOf(OPMException.class);
        }

        @Test
        public void uninstantiableEntity_throws() {
            PlcEntityManager entityManager = new PlcEntityManager();

            assertThatThrownBy(() -> entityManager.read(UninstantiableEntity.class, "mock:test"))
                .isInstanceOf(OPMException.class);
        }

        @Test
        public void resolveAlias_works() throws OPMException, PlcConnectionException {
            SimpleAliasRegistry registry = new SimpleAliasRegistry();
            registry.register("alias", "real_field");

            // Mock
            PlcDriverManager driverManager = new PlcDriverManager();
            MockConnection connection = (MockConnection) driverManager.getConnection("mock:test");
            MockDevice mockDevice = Mockito.mock(MockDevice.class);
            when(mockDevice.read(any())).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcSTRING("value")));
            connection.setDevice(mockDevice);

            PlcEntityManager entityManager = new PlcEntityManager(driverManager, registry);
            entityManager.read(AliasEntity.class, "mock:test");

            // Assert that "field" was queried
            verify(mockDevice).read(eq("real_field"));
        }


        @Test
        public void unknownAlias_throws() {
            PlcEntityManager entityManager = new PlcEntityManager();

            assertThatThrownBy(() -> entityManager.read(AliasEntity.class, "mock:test"))
                .isInstanceOf(IllegalArgumentException.class);
        }

        @Test
        public void badAlias_throws() {
            PlcEntityManager entityManager = new PlcEntityManager();

            String message = null;
            try {
                entityManager.read(BadAliasEntity.class, "mock:test");
            } catch (IllegalArgumentException e) {
                message = e.getMessage();
            } catch (OPMException e) {
                fail("Unexpected Exception" + e);
            }

            assertNotNull(message);
            assertTrue(message.contains("Invalid Syntax, either use field address (no starting $) or an alias with Syntax ${xxx}. But given was"));
        }
    }

    @Nested
    class Write {

        @Test
        void simpleWrite() throws Exception {
            SimpleAliasRegistry registry = new SimpleAliasRegistry();
            registry.register("alias", "real_field");

            // Mock
            PlcDriverManager driverManager = new PlcDriverManager();
            MockConnection connection = (MockConnection) driverManager.getConnection("mock:test");
            MockDevice mockDevice = Mockito.mock(MockDevice.class);
            when(mockDevice.write(anyString(), any())).thenReturn(PlcResponseCode.OK);
            when(mockDevice.read(any())).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcSTRING("value")));
            connection.setDevice(mockDevice);

            PlcEntityManager entityManager = new PlcEntityManager(driverManager, registry);
            AliasEntity object = new AliasEntity();
            object.setAliasedField("changed");
            AliasEntity connected = entityManager.write(AliasEntity.class, "mock:test", object);
            connected.setAliasedField("changed2");
            connected.getAliasedField();
            verify(mockDevice, times(0)).read(eq("real_field"));
            verify(mockDevice, times(1)).write(eq("real_field"), any());
            AliasEntity merge = entityManager.merge(AliasEntity.class, "mock:test", connected);
            merge.setAliasedField("changed2");
            merge.getAliasedField();

            // Assert that "field" was queried
            verify(mockDevice, times(1)).read(eq("real_field"));
            verify(mockDevice, times(3)).write(eq("real_field"), any());

            entityManager.disconnect(merge);
            assertThat(merge.getAliasedField()).isEqualTo("value");
        }

        @Test
        void simpleWrite_uses_getter() throws Exception {
            // Mock
            PlcDriverManager driverManager = new PlcDriverManager();
            MockConnection connection = (MockConnection) driverManager.getConnection("mock:test");
            MockDevice mockDevice = Mockito.mock(MockDevice.class);
            when(mockDevice.read(any())).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcSTRING("value")));
            connection.setDevice(mockDevice);

            PlcEntityManager entityManager = new PlcEntityManager(driverManager);
            CustomGetterEntity connect = entityManager.connect(CustomGetterEntity.class, "mock:test");
            assertThat(connect.getAsd()).isEqualTo("value!");
        }
    }

    @Nested
    class Lifecycle {
        /**
         * Class is private, so EntityManager has no access to it
         *
         * @throws OPMException
         */
        @Test
        public void connect_uninstantiableEntity_throws() {
            PlcEntityManager entityManager = new PlcEntityManager();

            assertThatThrownBy(() -> entityManager.connect(UninstantiableEntity.class, "mock:test"))
                .isInstanceOf(OPMException.class);
        }

        @Test
        public void connect_resolveAlias_works() throws PlcConnectionException, OPMException {
            SimpleAliasRegistry registry = new SimpleAliasRegistry();
            registry.register("alias", "real_field");

            // Mock
            PlcDriverManager driverManager = new PlcDriverManager();
            MockConnection connection = (MockConnection) driverManager.getConnection("mock:test");
            MockDevice mockDevice = Mockito.mock(MockDevice.class);
            when(mockDevice.read(any())).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcSTRING("value")));
            connection.setDevice(mockDevice);

            PlcEntityManager entityManager = new PlcEntityManager(driverManager, registry);
            entityManager.connect(AliasEntity.class, "mock:test");

            // Assert that "field" was queried
            verify(mockDevice, times(1)).read(eq("real_field"));
        }
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

        public void setAliasedField(String aliasedField) {
            this.aliasedField = aliasedField;
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

    @PlcEntity
    public static class CustomGetterEntity {

        @PlcField("asd")
        private String asd;

        public CustomGetterEntity() {
            // for OPM
        }

        public String getAsd() {
            return asd + "!";
        }
    }

}