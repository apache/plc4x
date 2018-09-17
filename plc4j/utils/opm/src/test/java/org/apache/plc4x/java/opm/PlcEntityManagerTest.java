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
import org.apache.plc4x.java.base.messages.items.FieldItem;
import org.junit.Assert;
import org.junit.Test;
import org.mockito.ArgumentMatchers;
import org.mockito.Mockito;

import java.util.Map;
import java.util.Optional;
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
    public void find() throws OPMException, PlcConnectionException {
        PlcEntityManager manager = getPlcEntityManager();

        MyEntity myEntity = manager.read(MyEntity.class);

        assertEquals(1, (long) myEntity.getCounter());
        assertEquals(1, (long) myEntity.getCounter2());
    }

    @Test
    public void connect() throws PlcConnectionException, OPMException {
        PlcEntityManager manager = getPlcEntityManager();

        ConnectedEntity connect = manager.connect(ConnectedEntity.class);

        Assert.assertNotNull(connect);

        long value = connect.getValue();

        assertEquals(1, value);
    }

    private PlcEntityManager getPlcEntityManager() throws PlcConnectionException {
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
                                return Pair.of(PlcResponseCode.OK, new DefaultIntegerFieldItem(1L));
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
    public interface ConnectedEntity {

        @PlcField("%DB1.DW111")
        public Long getValue();

    }
}