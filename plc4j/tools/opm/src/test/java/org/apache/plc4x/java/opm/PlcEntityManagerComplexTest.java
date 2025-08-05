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
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.mock.tag.MockTag;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.messages.PlcReader;
import org.apache.plc4x.java.spi.messages.PlcWriter;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcLINT;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
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
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Function;
import java.util.stream.Collectors;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.when;


public class PlcEntityManagerComplexTest implements WithAssertions {

    private DefaultPlcDriverManager driverManager;

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
        PlcValueHandler valueHandler = new DefaultPlcValueHandler();
        map.put(prefix + "boolVar", valueHandler.newPlcValue(new MockTag(prefix + "boolVar", PlcValueType.BOOL), true));
        map.put(prefix + "byteVar", valueHandler.newPlcValue(new MockTag(prefix + "byteVar", PlcValueType.SINT), (byte) 1));
        map.put(prefix + "shortVar", valueHandler.newPlcValue(new MockTag(prefix + "shortVar", PlcValueType.INT), (short) 1));
        map.put(prefix + "intVar", valueHandler.newPlcValue(new MockTag(prefix + "intVar", PlcValueType.DINT), 1));
        map.put(prefix + "longVar", valueHandler.newPlcValue(new MockTag(prefix + "longVar", PlcValueType.LINT), 1L));
        map.put(prefix + "boxedBoolVar", valueHandler.newPlcValue(new MockTag(prefix + "boxedBoolVar", PlcValueType.BOOL), 1L));
        map.put(prefix + "boxedByteVar", valueHandler.newPlcValue(new MockTag(prefix + "boxedByteVar", PlcValueType.SINT), (byte) 1));
        map.put(prefix + "boxedShortVar", valueHandler.newPlcValue(new MockTag(prefix + "boxedShortVar", PlcValueType.INT), (short) 1));
        map.put(prefix + "boxedIntegerVar", valueHandler.newPlcValue(new MockTag(prefix + "boxedIntegerVar", PlcValueType.DINT), 1));
        map.put(prefix + "boxedLongVar", valueHandler.newPlcValue(new MockTag(prefix + "boxedLongVar", PlcValueType.LINT), 1L));
        map.put(prefix + "bigIntegerVar", valueHandler.newPlcValue(new MockTag(prefix + "bigIntegerVar", PlcValueType.LINT), BigInteger.ONE));
        map.put(prefix + "floatVar", valueHandler.newPlcValue(new MockTag(prefix + "floatVar", PlcValueType.REAL), 1f));
        map.put(prefix + "doubleVar", valueHandler.newPlcValue(new MockTag(prefix + "doubleVar", PlcValueType.LREAL), 1d));
        map.put(prefix + "bigDecimalVar", valueHandler.newPlcValue(new MockTag(prefix + "bigDecimalVar", PlcValueType.LREAL), BigDecimal.ONE));
        map.put(prefix + "localTimeVar", valueHandler.newPlcValue(new MockTag(prefix + "localTimeVar", PlcValueType.TIME), LocalTime.of(1, 1)));
        map.put(prefix + "localDateVar", valueHandler.newPlcValue(new MockTag(prefix + "localDateVar", PlcValueType.DATE), LocalDate.of(1, 1, 1)));
        map.put(prefix + "localDateTimeVar", valueHandler.newPlcValue(new MockTag(prefix + "localDateTimeVar", PlcValueType.DATE_AND_TIME), LocalDateTime.of(1, 1, 1, 1, 1)));
        map.put(prefix + "byteArrayVar", valueHandler.newPlcValue(new MockTag(prefix + "byteArrayVar", PlcValueType.RAW_BYTE_ARRAY), new Byte[]{0x0, 0x1}));
        map.put(prefix + "bigByteArrayVar", valueHandler.newPlcValue(new MockTag(prefix + "bigByteArrayVar", PlcValueType.RAW_BYTE_ARRAY), new Byte[]{0x0, 0x1}));
        map.put(prefix + "stringVar", valueHandler.newPlcValue(new MockTag(prefix + "stringVar", PlcValueType.STRING), "Hallo"));
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
        assertTrue(connect.isBoolVar());
    }

    @Test
    @Disabled("Sebastian please fix this")
    public void disconnect() throws PlcConnectionException, OPMException, IllegalAccessException {
        PlcEntityManager manager = getInitializedEntityManager();

        ConnectedEntity connected = manager.connect(ConnectedEntity.class, "s7://localhost:5555/0/0");

        manager.disconnect(connected);

        // Assert disconnected
        Object o = FieldUtils.readDeclaredField(connected, PlcEntityManager.CONNECTION_MANAGER_FIELD_NAME, true);
        assertNull(o);

        // Call a method and receive the result
        // We are ok if a result is received and no NPE is thrown, then everything works as expected
        assertNotNull(connected.toString());
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
        driverManager = Mockito.mock(DefaultPlcDriverManager.class);
        DefaultPlcDriverManager mock = driverManager;
        PlcConnection connection = Mockito.mock(PlcConnection.class);
        when(mock.getConnection(ArgumentMatchers.anyString())).thenReturn(connection);
        when(connection.getMetadata()).thenReturn(new PlcConnectionMetadata() {

            @Override
            public boolean isReadSupported() {
                return true;
            }

            @Override
            public boolean isWriteSupported() {
                return true;
            }

            @Override
            public boolean isSubscribeSupported() {
                return true;
            }

            @Override
            public boolean isBrowseSupported() {
                return true;
            }

        });

        PlcReader reader = readRequest -> {
            Map<String, PlcResponseItem<PlcValue>> map = readRequest.getTagNames().stream()
                .collect(Collectors.toMap(
                    Function.identity(),
                    s -> new DefaultPlcResponseItem<>(PlcResponseCode.OK, Objects.requireNonNull(responses.get(s), s + " not found"))
                ));
            return CompletableFuture.completedFuture(new DefaultPlcReadResponse(readRequest, map));
        };
        when(connection.readRequestBuilder()).then(invocation -> new DefaultPlcReadRequest.Builder(reader, getTagHandler()));
        PlcWriter writer = writeRequest -> {
            Map<String, PlcResponseCode> map = writeRequest.getTagNames().stream()
                .collect(Collectors.toMap(
                    Function.identity(),
                    s -> PlcResponseCode.OK
                ));
            return CompletableFuture.completedFuture(new DefaultPlcWriteResponse(writeRequest, map));
        };
        when(connection.writeRequestBuilder()).then(invocation -> new DefaultPlcWriteRequest.Builder(writer, getTagHandler(), new DefaultPlcValueHandler()));

        return new PlcEntityManager(mock);
    }

    private PlcTagHandler getTagHandler() {
        return new NoOpPlcTagHandler();
    }

    private static class NoOpPlcTagHandler implements PlcTagHandler {
        @Override
        public org.apache.plc4x.java.api.model.PlcTag parseTag(String tagAddress) throws PlcInvalidTagException {
            return new org.apache.plc4x.java.api.model.PlcTag() {
                @Override
                public String getAddressString() {
                    return "address";
                }

                @Override
                public PlcValueType getPlcValueType() {
                    return org.apache.plc4x.java.api.model.PlcTag.super.getPlcValueType();
                }

                @Override
                public List<ArrayInfo> getArrayInfo() {
                    return org.apache.plc4x.java.api.model.PlcTag.super.getArrayInfo();
                }
            };
        }

        @Override
        public PlcQuery parseQuery(String query) {
            throw new UnsupportedOperationException("This driver doesn't support browsing");
        }
    }

    private static class NoEntity {

    }

    @PlcEntity()
    public static class EntityWithBadConstructor {

        @PlcTag("asdf")
        private final long field;

        public EntityWithBadConstructor(long field) {
            this.field = field;
        }

        public long getField() {
            return field;
        }
    }

    @PlcEntity()
    public static class MyEntity {

        @PlcTag("%DB3.DBW500")
        private Long counter;

        @PlcTag("%DB3.DBW504")
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

        @PlcTag("%DB1.DW111:BOOL")
        private boolean boolVar;
        @PlcTag("%DB1.DW111:BYTE")
        private byte byteVar;
        @PlcTag("%DB1.DW111:SHORT")
        private short shortVar;
        @PlcTag("%DB1.DW111:INT")
        private int intVar;
        @PlcTag("%DB1.DW111:LONG")
        private long longVar;
        @PlcTag("%DB1.DW111:BOOL")
        private Boolean boxedBoolVar;
        @PlcTag("%DB1.DW111:BYTE")
        private Byte boxedByteVar;
        @PlcTag("%DB1.DW111:SHORT")
        private Short boxedShortVar;
        @PlcTag("%DB1.DW111:SHORT")
        private Integer boxedIntegerVar;
        @PlcTag("%DB1.DW111:LONG")
        private Long boxedLongVar;
        @PlcTag("%DB1.DW111:BIGINT")
        private BigInteger bigIntegerVar;
        @PlcTag("%DB1.DW111:FLOAT")
        private Float floatVar;
        @PlcTag("%DB1.DW111:DOUBLE")
        private Double doubleVar;
        @PlcTag("%DB1.DW111:BIGDECIMAL")
        private BigDecimal bigDecimalVar;
        @PlcTag("%DB1.DW111:LOCALTIME")
        private LocalTime localTimeVar;
        @PlcTag("%DB1.DW111:LOCALDATE")
        private LocalDate localDateVar;
        @PlcTag("%DB1.DW111:LOCALDATETIME")
        private LocalDateTime localDateTimeVar;
        @PlcTag("%DB1.DW111:BYTEARRAY")
        private byte[] byteArrayVar;
        @PlcTag("%DB1.DW111:BYTEARRAY")
        private Byte[] bigByteArrayVar;

        @PlcTag("%DB1.DW111:STRING")
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
