package org.apache.plc4x.java.opm;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidAddressException;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.api.types.ResponseCode;
import org.junit.Assert;
import org.junit.Test;
import org.mockito.ArgumentMatchers;
import org.mockito.Mockito;

import java.util.List;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;

import static org.mockito.Mockito.when;

public class PlcEntityManagerTest {

    @Test(expected = IllegalArgumentException.class)
    public void noEntity_throws() throws OPMException {
        PlcEntityManager manager = new PlcEntityManager();

        manager.find(NoEntity.class);
    }

    @Test(expected = IllegalArgumentException.class)
    public void noValidConstructor_throws() throws OPMException {
        PlcEntityManager manager = new PlcEntityManager();

        manager.find(EntityWithBadConstructor.class);
    }

    @Test
    public void find() throws OPMException, PlcConnectionException, PlcInvalidAddressException {
        PlcEntityManager manager = getPlcEntityManager();

        MyEntity myEntity = manager.find(MyEntity.class);

        Assert.assertEquals(1, (long) myEntity.getCounter());
        Assert.assertEquals(1, (long) myEntity.getCounter2());
    }

    private PlcEntityManager getPlcEntityManager() throws PlcConnectionException, PlcInvalidAddressException {
        PlcDriverManager mock = Mockito.mock(PlcDriverManager.class);
        PlcConnection connection = Mockito.mock(PlcConnection.class);
        when(mock.getConnection(ArgumentMatchers.anyString())).thenReturn(connection);
        PlcReadResponse response = Mockito.mock(PlcReadResponse.class);
        PlcReader reader = plcReadRequest -> {
            List<ReadResponseItem<?>> responseItems = plcReadRequest.getRequestItems().stream()
                .map(item -> {
                    if (item.getDatatype() == Long.class || item.getDatatype() == long.class) {
                        return new ReadResponseItem<Long>((ReadRequestItem<Long>) item, ResponseCode.OK, new Long(1));
                    } else {
                        throw new IllegalArgumentException("Unable to mock request with class " + item.getDatatype());
                    }
                })
                .collect(Collectors.toList());
            return CompletableFuture.completedFuture(new PlcReadResponse(plcReadRequest, responseItems));
        };
        when(connection.getReader()).thenReturn(Optional.of(reader));
        // Return two different addresses to keep different request items
        when(connection.parseAddress(ArgumentMatchers.anyString()))
            .thenReturn(new Address() {
            })
            .thenReturn(new Address() {
            });

        return new PlcEntityManager(mock);
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
}