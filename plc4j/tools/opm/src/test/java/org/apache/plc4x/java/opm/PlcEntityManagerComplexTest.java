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
package org.apache.plc4x.java.opm;

import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.messages.PlcReader;
import org.apache.plc4x.java.spi.messages.PlcWriter;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcLINT;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentMatchers;
import org.mockito.Mockito;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.CompletableFuture;
import java.util.function.Function;
import java.util.stream.Collectors;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.when;


public class PlcEntityManagerComplexTest implements WithAssertions {

    private PlcDriverManager driverManager;

    @Test
    public void noEntity_throws() {
        PlcEntityManager manager = new PlcEntityManager();

        assertThatThrownBy(() -> manager.read(NoEntity.class, "s7://localhost:5555/0/0"))
            .isInstanceOf(IllegalArgumentException.class);
    }

    @Test
    public void noValidConstructor_throws() {
        PlcEntityManager manager = new PlcEntityManager();

        assertThatThrownBy(() -> manager.read(EntityWithBadConstructor.class, "s7://localhost:5555/0/0"))
            .isInstanceOf(IllegalArgumentException.class);
    }

    @Test
    public void read() throws OPMException, PlcConnectionException {
        Map<String, PlcValue> results = new HashMap<>();
        String prefix = MyEntity.class.getName() + ".";
        results.put(prefix + "counter", new PlcDINT(1));
        results.put(prefix + "counter2", new PlcLINT(1L));
        PlcEntityManager manager = getPlcEntityManager(results);

        MyEntity myEntity = manager.read(MyEntity.class, "s7://localhost:5555/0/0");

        assertEquals(1, (long) myEntity.getCounter());
        assertEquals(1, myEntity.getCounter2());
    }

    @Test
    @Disabled("Sebastian please fix this")
    public void readComplexObject() throws PlcConnectionException, OPMException {
        PlcEntityManager manager = getInitializedEntityManager();

        ConnectedEntity connect = manager.read(ConnectedEntity.class, "s7://localhost:5555/0/0");

        assertNotNull(connect);

        // Call different method
        String s = connect.toString();

        assertEquals("ConnectedEntity{boolVar=true, byteVar=1, shortVar=1, intVar=1, longVar=1, boxedBoolVar=true, boxedByteVar=1, boxedShortVar=1, boxedIntegerVar=1, boxedLongVar=1, bigIntegerVar=1, floatVar=1.0, doubleVar=1.0, bigDecimalVar=1, localTimeVar=01:01, localDateVar=0001-01-01, localDateTimeVar=0001-01-01T01:01, byteArrayVar=[0, 1], bigByteArrayVar=[0, 1], stringVar='Hallo'}", s);
    }

    @Test
    @Disabled("Sebastian please fix this")
    public void connect_callComplexMethod() throws PlcConnectionException, OPMException {
        PlcEntityManager manager = getInitializedEntityManager();

        ConnectedEntity connect = manager.connect(ConnectedEntity.class, "s7://localhost:5555/0/0");

        assertNotNull(connect);

        // Call different method
        String s = connect.toString();

        assertEquals("ConnectedEntity{boolVar=true, byteVar=1, shortVar=1, intVar=1, longVar=1, boxedBoolVar=true, boxedByteVar=1, boxedShortVar=1, boxedIntegerVar=1, boxedLongVar=1, bigIntegerVar=1, floatVar=1.0, doubleVar=1.0, bigDecimalVar=1, localTimeVar=01:01, localDateVar=0001-01-01, localDateTimeVar=0001-01-01T01:01, byteArrayVar=[0, 1], bigByteArrayVar=[0, 1], stringVar='Hallo'}", s);
    }

    private PlcEntityManager getInitializedEntityManager() throws PlcConnectionException {
        Map<String, PlcValue> map = new HashMap<>();
        String prefix = ConnectedEntity.class.getName() + ".";
        map.put(prefix + "boolVar", IEC61131ValueHandler.of(true));
        map.put(prefix + "byteVar", IEC61131ValueHandler.of((byte) 1));
        map.put(prefix + "shortVar", IEC61131ValueHandler.of((short) 1));
        map.put(prefix + "intVar", IEC61131ValueHandler.of(1));
        map.put(prefix + "longVar", IEC61131ValueHandler.of(1L));
        map.put(prefix + "boxedBoolVar", IEC61131ValueHandler.of(1L));
        map.put(prefix + "boxedByteVar", IEC61131ValueHandler.of((byte) 1));
        map.put(prefix + "boxedShortVar", IEC61131ValueHandler.of((short) 1));
        map.put(prefix + "boxedIntegerVar", IEC61131ValueHandler.of(1));
        map.put(prefix + "boxedLongVar", IEC61131ValueHandler.of(1L));
        map.put(prefix + "bigIntegerVar", IEC61131ValueHandler.of(BigInteger.ONE));
        map.put(prefix + "floatVar", IEC61131ValueHandler.of(1f));
        map.put(prefix + "doubleVar", IEC61131ValueHandler.of(1d));
        map.put(prefix + "bigDecimalVar", IEC61131ValueHandler.of(BigDecimal.ONE));
        map.put(prefix + "localTimeVar", IEC61131ValueHandler.of(LocalTime.of(1, 1)));
        map.put(prefix + "localDateVar", IEC61131ValueHandler.of(LocalDate.of(1, 1, 1)));
        map.put(prefix + "localDateTimeVar", IEC61131ValueHandler.of(LocalDateTime.of(1, 1, 1, 1, 1)));
        map.put(prefix + "byteArrayVar", IEC61131ValueHandler.of(new Byte[]{0x0, 0x1}));
        map.put(prefix + "bigByteArrayVar", IEC61131ValueHandler.of(new Byte[]{0x0, 0x1}));
        map.put(prefix + "stringVar", IEC61131ValueHandler.of("Hallo"));
        return getPlcEntityManager(map);
    }

    @Test
    @Disabled("Sebastian please fix this")
    public void connect_callGetter() throws PlcConnectionException, OPMException {
        PlcEntityManager manager = getInitializedEntityManager();

        ConnectedEntity connect = manager.connect(ConnectedEntity.class, "s7://localhost:5555/0/0");

        assertNotNull(connect);

        // Call getter
        assertEquals(1, connect.getIntVar());
        assertEquals("Hallo", connect.getStringVar());
        assertEquals(true, connect.isBoolVar());
    }

    @Test
    @Disabled("Sebastian please fix this")
    public void disconnect() throws PlcConnectionException, OPMException, IllegalAccessException {
        PlcEntityManager manager = getInitializedEntityManager();

        ConnectedEntity connected = manager.connect(ConnectedEntity.class, "s7://localhost:5555/0/0");

        manager.disconnect(connected);

        // Assert disconnected
        Object o = FieldUtils.readDeclaredField(connected, PlcEntityManager.DRIVER_MANAGER_FIELD_NAME, true);
        assertNull(o);

        // Call a method and receive the result
        // We are ok if a result is received and no NPE is thrown, then everything works as expected
        assertNotNull(connected.toString());
        assertNotNull(connected.getByteVar());
    }

    @Test
    @Disabled("Sebastian please fix this")
    public void disconnectTwice_throwsException() throws PlcConnectionException, OPMException {
        PlcEntityManager manager = getInitializedEntityManager();

        ConnectedEntity connected = manager.connect(ConnectedEntity.class, "s7://localhost:5555/0/0");

        manager.disconnect(connected);
        assertThatThrownBy(() -> manager.disconnect(connected))
            .isInstanceOf(OPMException.class);
    }

    private PlcEntityManager getPlcEntityManager(final Map<String, PlcValue> responses) throws PlcConnectionException {
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

            @Override
            public boolean canBrowse() {
                return true;
            }

        });

        PlcReader reader = readRequest -> {
            Map<String, ResponseItem<PlcValue>> map = readRequest.getFieldNames().stream()
                .collect(Collectors.toMap(
                    Function.identity(),
                    s -> new ResponseItem<>(PlcResponseCode.OK, Objects.requireNonNull(responses.get(s), s + " not found"))
                ));
            return CompletableFuture.completedFuture(new DefaultPlcReadResponse(readRequest, map));
        };
        when(connection.readRequestBuilder()).then(invocation -> new DefaultPlcReadRequest.Builder(reader, getFieldHandler()));
        PlcWriter writer = writeRequest -> {
            Map<String, PlcResponseCode> map = writeRequest.getFieldNames().stream()
                .collect(Collectors.toMap(
                    Function.identity(),
                    s -> PlcResponseCode.OK
                ));
            return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest, map));
        };
        when(connection.writeRequestBuilder()).then(invocation -> new DefaultPlcWriteRequest.Builder(writer, getFieldHandler(), getValueHandler()));

        return new PlcEntityManager(mock);
    }

    private PlcFieldHandler getFieldHandler() {
        return new NoOpPlcFieldHandler();
    }

    private PlcValueHandler getValueHandler() {
        return new NoOpPlcValueHandler();
    }

    private static class NoOpPlcFieldHandler implements PlcFieldHandler {
        @Override
        public org.apache.plc4x.java.api.model.PlcField createField(String fieldQuery) throws PlcInvalidFieldException {
            return new org.apache.plc4x.java.api.model.PlcField() {
            };
        }
    }

    private static class NoOpPlcValueHandler implements PlcValueHandler {
        @Override
        public PlcValue newPlcValue(Object value) {
            throw new RuntimeException("Data Type " + value.getClass().getSimpleName() + "Is not supported");
        }

        @Override
        public PlcValue newPlcValue(Object[] values) {
            throw new RuntimeException("Data Type " + values.getClass().getSimpleName() + "Is not supported");
        }

        @Override
        public PlcValue newPlcValue(org.apache.plc4x.java.api.model.PlcField field, Object value) {
            throw new RuntimeException("Data Type " + value.getClass().getSimpleName() + "Is not supported");
        }

        @Override
        public PlcValue newPlcValue(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
            throw new RuntimeException("Data Type " + values.getClass().getSimpleName() + "Is not supported");
        }
    }

    private static class NoEntity {

    }

    @PlcEntity()
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

    @PlcEntity()
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

    @PlcEntity()
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
        @PlcField("%DB1.DW111:BOOL")
        private Boolean boxedBoolVar;
        @PlcField("%DB1.DW111:BYTE")
        private Byte boxedByteVar;
        @PlcField("%DB1.DW111:SHORT")
        private Short boxedShortVar;
        @PlcField("%DB1.DW111:SHORT")
        private Integer boxedIntegerVar;
        @PlcField("%DB1.DW111:LONG")
        private Long boxedLongVar;
        @PlcField("%DB1.DW111:BIGINT")
        private BigInteger bigIntegerVar;
        @PlcField("%DB1.DW111:FLOAT")
        private Float floatVar;
        @PlcField("%DB1.DW111:DOUBLE")
        private Double doubleVar;
        @PlcField("%DB1.DW111:BIGDECIMAL")
        private BigDecimal bigDecimalVar;
        @PlcField("%DB1.DW111:LOCALTIME")
        private LocalTime localTimeVar;
        @PlcField("%DB1.DW111:LOCALDATE")
        private LocalDate localDateVar;
        @PlcField("%DB1.DW111:LOCALDATETIME")
        private LocalDateTime localDateTimeVar;
        @PlcField("%DB1.DW111:BYTEARRAY")
        private byte[] byteArrayVar;
        @PlcField("%DB1.DW111:BYTEARRAY")
        private Byte[] bigByteArrayVar;

        @PlcField("%DB1.DW111:STRING")
        private String stringVar;

        public ConnectedEntity() {
            // Default
        }

        public boolean isBoolVar() {
            return boolVar;
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
                ", boxedBoolVar=" + boxedBoolVar +
                ", boxedByteVar=" + boxedByteVar +
                ", boxedShortVar=" + boxedShortVar +
                ", boxedIntegerVar=" + boxedIntegerVar +
                ", boxedLongVar=" + boxedLongVar +
                ", bigIntegerVar=" + bigIntegerVar +
                ", floatVar=" + floatVar +
                ", doubleVar=" + doubleVar +
                ", bigDecimalVar=" + bigDecimalVar +
                ", localTimeVar=" + localTimeVar +
                ", localDateVar=" + localDateVar +
                ", localDateTimeVar=" + localDateTimeVar +
                ", byteArrayVar=" + Arrays.toString(byteArrayVar) +
                ", bigByteArrayVar=" + Arrays.toString(bigByteArrayVar) +
                ", stringVar='" + stringVar + '\'' +
                '}';
        }
    }

}
