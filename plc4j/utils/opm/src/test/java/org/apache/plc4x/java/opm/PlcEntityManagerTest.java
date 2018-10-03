package org.apache.plc4x.java.opm;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;
import org.apache.plc4x.java.base.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.base.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.base.messages.InternalPlcReadRequest;
import org.apache.plc4x.java.base.messages.items.DefaultIntegerFieldItem;
import org.apache.plc4x.java.base.messages.items.DefaultStringFieldItem;
import org.apache.plc4x.java.base.messages.items.FieldItem;
import org.junit.Assert;
import org.junit.Test;
import org.mockito.ArgumentCaptor;
import org.mockito.ArgumentMatchers;
import org.mockito.Mockito;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.function.Function;
import java.util.stream.Collectors;

import static org.junit.Assert.assertEquals;
import static org.mockito.Mockito.verify;
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
    public void find() throws OPMException, PlcConnectionException {
        Map<String, FieldItem> results = new HashMap<>();
        results.put("counter", new DefaultIntegerFieldItem(1L));
        results.put("counter2", new DefaultIntegerFieldItem(1L));
        PlcEntityManager manager = getPlcEntityManager(results);

        MyEntity myEntity = manager.read(MyEntity.class);

        assertEquals(1, (long) myEntity.getCounter());
        assertEquals(1, (long) myEntity.getCounter2());
    }

    @Test
    public void connect() throws PlcConnectionException, OPMException {
        PlcEntityManager manager = getPlcEntityManager(Collections.singletonMap("%DB1.DW111", new DefaultIntegerFieldItem(1L)));

        ConnectedEntity connect = manager.connect(ConnectedEntity.class);

        Assert.assertNotNull(connect);

        long value = connect.getLongVar();

        assertEquals(1, value);
    }

    @Test
    public void connec_callComplexMethodt() throws PlcConnectionException, OPMException {
        Map<String, FieldItem> map = new HashMap<>();
        map.put("byteVar", new DefaultIntegerFieldItem(1L));
        map.put("shortVar", new DefaultIntegerFieldItem(1L));
        map.put("intVar", new DefaultIntegerFieldItem(1L));
        map.put("longVar", new DefaultIntegerFieldItem(1L));
        map.put("boxedLongVar", new DefaultIntegerFieldItem(1L));
        map.put("stringVar", new DefaultStringFieldItem("Hallo"));
        PlcEntityManager manager = getPlcEntityManager(map);

        ConnectedEntity connect = manager.connect(ConnectedEntity.class);

        Assert.assertNotNull(connect);

        String s = connect.toString();

        assertEquals("ConnectedEntity{byteVar=1, shortVar=1, intVar=1, longVar=1, boxedLongVar=1, stringVar='Hallo'}", s);
    }

    @Test
    public void callRandomMeathod() throws PlcConnectionException, OPMException {
        driverManager = Mockito.mock(PlcDriverManager.class);
        PlcDriverManager mock = driverManager;
        PlcConnection connection = Mockito.mock(PlcConnection.class);
        when(mock.getConnection(ArgumentMatchers.anyString())).thenReturn(connection);
        PlcReader reader = Mockito.mock(PlcReader.class);
        when(connection.getReader()).thenReturn(Optional.of(reader));
        PlcEntityManager manager = new PlcEntityManager(driverManager);

        ConnectedEntity connect = manager.connect(ConnectedEntity.class);

        try {
            connect.someMethod();
        } catch (NullPointerException e) {
            // ignore
        }

        ArgumentCaptor<PlcReadRequest> captor = ArgumentCaptor.forClass(PlcReadRequest.class);
        verify(reader).read(captor.capture());

        System.out.println(captor);
    }

    private PlcEntityManager getPlcEntityManager(final Map<String, FieldItem> responses) throws PlcConnectionException {
        driverManager = Mockito.mock(PlcDriverManager.class);
        PlcDriverManager mock = driverManager;
        PlcConnection connection = Mockito.mock(PlcConnection.class);
        when(mock.getConnection(ArgumentMatchers.anyString())).thenReturn(connection);
        PlcReader reader = new PlcReader() {
            @Override
            public CompletableFuture<PlcReadResponse<?>> read(PlcReadRequest readRequest) {
                Map<String, Pair<PlcResponseCode, FieldItem>> map = readRequest.getFieldNames().stream()
                    .collect(Collectors.toMap(
                        Function.identity(),
                        new Function<String, Pair<PlcResponseCode, FieldItem>>() {
                            @Override
                            public Pair<PlcResponseCode, FieldItem> apply(String s) {
                                return Pair.of(PlcResponseCode.OK, responses.get(s));
                            }
                        }
                    ));
                return CompletableFuture.completedFuture(new DefaultPlcReadResponse(((InternalPlcReadRequest) readRequest), map));
            }

            @Override
            public PlcReadRequest.Builder readRequestBuilder() {
                return new DefaultPlcReadRequest.Builder(getFieldHandler());
            }
        };
        when(connection.getReader()).thenReturn(Optional.of(reader));

        return new PlcEntityManager(mock);
    }

    private PlcFieldHandler getFieldHandler() {
        return new PlcFieldHandler() {
            @Override
            public org.apache.plc4x.java.api.model.PlcField createField(String fieldQuery) throws PlcInvalidFieldException {
                return new org.apache.plc4x.java.api.model.PlcField() {
                };
            }

            @Override
            public FieldItem encodeBoolean(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeByte(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeShort(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeInteger(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeBigInteger(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeLong(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeFloat(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeDouble(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeString(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeTime(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeDate(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }

            @Override
            public FieldItem encodeDateTime(org.apache.plc4x.java.api.model.PlcField field, Object[] values) {
                return null;
            }
        };
    }

    public static class NoEntity {

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
                "byteVar=" + byteVar +
                ", shortVar=" + shortVar +
                ", intVar=" + intVar +
                ", longVar=" + longVar +
                ", boxedLongVar=" + boxedLongVar +
                ", stringVar='" + stringVar + '\'' +
                '}';
        }
    }
}