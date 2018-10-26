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

package org.apache.plc4x.java.opm;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;
import org.apache.plc4x.java.base.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.base.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.base.messages.InternalPlcReadRequest;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.apache.plc4x.java.base.messages.items.*;
import org.junit.Assert;
import org.junit.Test;
import org.mockito.ArgumentMatchers;
import org.mockito.Mockito;

import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.CompletableFuture;
import java.util.function.Function;
import java.util.stream.Collectors;

import static org.junit.Assert.assertEquals;
import static org.mockito.Mockito.when;


public class PlcEntityManagerTest {

    private PlcDriverManager driverManager;

    @Test(expected = IllegalArgumentException.class)
    public void noEntity_throws() throws OPMException {
        PlcEntityManager manager = new PlcEntityManager();

        manager.read(NoEntity.class);
    }

    @Test(expected = IllegalArgumentException.class)
    public void noValidConstructor_throws() throws OPMException {
        PlcEntityManager manager = new PlcEntityManager();

        manager.read(EntityWithBadConstructor.class);
    }

    @Test
    public void read() throws OPMException, PlcConnectionException {
        Map<String, BaseDefaultFieldItem> results = new HashMap<>();
        String prefix = MyEntity.class.getName() + ".";
        results.put(prefix + "counter", new DefaultIntegerFieldItem(1));
        results.put(prefix + "counter2", new DefaultLongFieldItem(1l));
        PlcEntityManager manager = getPlcEntityManager(results);

        MyEntity myEntity = manager.read(MyEntity.class);

        assertEquals(1, (long) myEntity.getCounter());
        assertEquals(1, myEntity.getCounter2());
    }

    @Test
    public void readComplexObject() throws PlcConnectionException, OPMException {
        Map<String, BaseDefaultFieldItem> map = new HashMap<>();
        String prefix = ConnectedEntity.class.getName() + ".";
        map.put(prefix + "boolVar", new DefaultBooleanFieldItem(true));
        map.put(prefix + "byteVar", new DefaultByteFieldItem((byte) 1));
        map.put(prefix + "shortVar", new DefaultShortFieldItem((short) 1));
        map.put(prefix + "intVar", new DefaultIntegerFieldItem(1));
        map.put(prefix + "longVar", new DefaultLongFieldItem(1l));
        map.put(prefix + "boxedLongVar", new DefaultLongFieldItem(1L));
        map.put(prefix + "stringVar", new DefaultStringFieldItem("Hallo"));
        PlcEntityManager manager = getPlcEntityManager(map);

        ConnectedEntity connect = manager.read(ConnectedEntity.class);

        Assert.assertNotNull(connect);

        // Call different mehtod
        String s = connect.toString();

        assertEquals("ConnectedEntity{boolVar=true, byteVar=1, shortVar=1, intVar=1, longVar=1, boxedLongVar=1, stringVar='Hallo'}", s);
    }

    @Test
    public void connect_callComplexMethod() throws PlcConnectionException, OPMException {
        Map<String, BaseDefaultFieldItem> map = new HashMap<>();
        String prefix = ConnectedEntity.class.getName() + ".";
        map.put(prefix + "boolVar", new DefaultBooleanFieldItem(true));
        map.put(prefix + "byteVar", new DefaultByteFieldItem((byte) 1));
        map.put(prefix + "shortVar", new DefaultShortFieldItem((short) 1));
        map.put(prefix + "intVar", new DefaultIntegerFieldItem(1));
        map.put(prefix + "longVar", new DefaultLongFieldItem(1l));
        map.put(prefix + "boxedLongVar", new DefaultLongFieldItem(1L));
        map.put(prefix + "stringVar", new DefaultStringFieldItem("Hallo"));
        PlcEntityManager manager = getPlcEntityManager(map);

        ConnectedEntity connect = manager.connect(ConnectedEntity.class);

        Assert.assertNotNull(connect);

        // Call different mehtod
        String s = connect.toString();

        assertEquals("ConnectedEntity{boolVar=true, byteVar=1, shortVar=1, intVar=1, longVar=1, boxedLongVar=1, stringVar='Hallo'}", s);
    }

    @Test
    public void connect_callGetter() throws PlcConnectionException, OPMException {
        Map<String, BaseDefaultFieldItem> map = new HashMap<>();
        map.put("getIntVar", new DefaultIntegerFieldItem(1));
        map.put("getStringVar", new DefaultStringFieldItem("Hello"));
        PlcEntityManager manager = getPlcEntityManager(map);

        ConnectedEntity connect = manager.connect(ConnectedEntity.class);

        Assert.assertNotNull(connect);

        // Call getter
        assertEquals(1, connect.getIntVar());
        assertEquals("Hello", connect.getStringVar());
    }

    private PlcEntityManager getPlcEntityManager(final Map<String, BaseDefaultFieldItem> responses) throws PlcConnectionException {
        driverManager = Mockito.mock(PlcDriverManager.class);
        PlcDriverManager mock = driverManager;
        PlcConnection connection = Mockito.mock(PlcConnection.class);
        when(mock.getConnection(ArgumentMatchers.anyString())).thenReturn(connection);
        when(connection.getMetadata()).thenReturn(new PlcConnectionMetadata() {

            @Override
            public boolean canRead() {
                return true;
            }

            @Override
            public boolean canWrite() {
                return true;
            }

            @Override
            public boolean canSubscribe() {
                return true;
            }
        });

        PlcReader reader = readRequest -> {
            Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> map = readRequest.getFieldNames().stream()
                .collect(Collectors.toMap(
                    Function.identity(),
                    s -> Pair.of(PlcResponseCode.OK, Objects.requireNonNull(responses.get(s), s + " not found"))
                ));
            return CompletableFuture.completedFuture(new DefaultPlcReadResponse((InternalPlcReadRequest) readRequest, map));
        };
        when(connection.readRequestBuilder()).then(invocation -> new DefaultPlcReadRequest.Builder(reader, getFieldHandler()));

        return new PlcEntityManager(mock);
    }

    private PlcFieldHandler getFieldHandler() {
        return new NoOpPlcFieldHandler();
    }

    private static class NoOpPlcFieldHandler implements PlcFieldHandler {
        @Override
        public org.apache.plc4x.java.api.model.PlcField createField(String fieldQuery) throws PlcInvalidFieldException {
            return new org.apache.plc4x.java.api.model.PlcField() {
            };
        }

        @Override
        public BaseDefaultFieldItem encodeBoolean(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeByte(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeShort(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeInteger(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeBigInteger(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeLong(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeFloat(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeBigDecimal(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeDouble(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeString(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeTime(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeDate(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeDateTime(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }

        @Override
        public BaseDefaultFieldItem encodeByteArray(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            return null;
        }
    }

    private static class NoEntity {

    }

    @PlcEntity("source")
    public static class EntityWithBadConstructor {

        @PlcField("asdf")
        private long field;

        public EntityWithBadConstructor(long field) {
            this.field = field;
        }

        public long getField() {
            return field;
        }
    }

    @PlcEntity("s7://localhost:5555/0/0")
    public static class MyEntity {

        @PlcField("%DB3.DBW500")
        private Long counter;

        @PlcField("%DB3.DBW504")
        private long counter2;

        public Long getCounter() {
            return counter;
        }

        public long getCounter2() {
            return counter2;
        }

    }

    @PlcEntity("s7://localhost:5555/0/0")
    public static class ConnectedEntity {

        @PlcField("%DB1.DW111:BOOL")
        private boolean boolVar;
        @PlcField("%DB1.DW111:BYTE")
        private byte byteVar;
        @PlcField("%DB1.DW111:SHORT")
        private short shortVar;
        @PlcField("%DB1.DW111:INT")
        private int intVar;
        @PlcField("%DB1.DW111:LONG")
        private long longVar;
        @PlcField("%DB1.DW111:STRING")
        private Long boxedLongVar;
        @PlcField("%DB1.DW111:STRING")
        private String stringVar;

        public ConnectedEntity() {
            // Default
        }


        public byte getByteVar() {
            return byteVar;
        }


        public short getShortVar() {
            return shortVar;
        }


        public int getIntVar() {
            return intVar;
        }


        public long getLongVar() {
            return longVar;
        }


        public String getStringVar() {
            return stringVar;
        }

        public void someMethod() {
            System.out.println("I do nothing");
        }

        @Override
        public String toString() {
            return "ConnectedEntity{" +
                "boolVar=" + boolVar +
                ", byteVar=" + byteVar +
                ", shortVar=" + shortVar +
                ", intVar=" + intVar +
                ", longVar=" + longVar +
                ", boxedLongVar=" + boxedLongVar +
                ", stringVar='" + stringVar + '\'' +
                '}';
        }
    }
}