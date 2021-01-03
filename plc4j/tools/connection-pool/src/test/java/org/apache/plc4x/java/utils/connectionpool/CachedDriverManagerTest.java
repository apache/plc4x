package org.apache.plc4x.java.utils.connectionpool;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;

import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.utils.connectionpool.CachedDriverManager;
import org.apache.plc4x.java.utils.connectionpool.PlcConnectionFactory;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

/**
 *
 * @author julian
 * Created by julian on 24.02.20
 */
class CachedDriverManagerTest implements WithAssertions {

    @Test
    void noConnectionWithoutRequest() throws PlcException {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        CachedDriverManager driverManager = new CachedDriverManager("", mock);

        verify(mock, never()).create();
    }

    @Test
    void establishConnectionAtFirstRequest() throws Exception {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        CachedDriverManager driverManager = new CachedDriverManager("", mock);

        driverManager.getConnection("").close();

        verify(mock, timeout(1_000).times(1)).create();
    }

    @Test
    void returnConnectionWhenIsActive() throws Exception {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        CachedDriverManager driverManager = new CachedDriverManager("", mock);

        driverManager.getConnection("").close();
//
//        verify(mock, timeout(500).times(1)).get();

        // State should now be Connected
        assertThat(driverManager.getState()).isEqualByComparingTo(CachedDriverManager.ConnectionState.AVAILABLE);
        Assertions.assertDoesNotThrow(() -> driverManager.getConnection(""));
        assertThat(driverManager.getState()).isEqualByComparingTo(CachedDriverManager.ConnectionState.BORROWED);
    }

    @Test
    void freeConnectionAfterReturn() throws Exception {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        CachedDriverManager driverManager = new CachedDriverManager("", mock);

        // Get Connmection
        PlcConnection connection = driverManager.getConnection("");
        // Close the Connection
        assertThat(driverManager.getState()).isEqualByComparingTo(CachedDriverManager.ConnectionState.BORROWED);
        connection.close();
        assertThat(driverManager.getState()).isEqualByComparingTo(CachedDriverManager.ConnectionState.AVAILABLE);
    }

    @Test
    void useClosedConnection() throws Exception {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        CachedDriverManager driverManager = new CachedDriverManager("", mock);

        // Get Connmection
        PlcConnection connection = driverManager.getConnection("");
        // Close the Connection
        connection.close();
        assertThrows(IllegalStateException.class, () -> connection.readRequestBuilder());
    }

    @Test
    void useClosedConnection2() throws Exception {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        MockConnection plcMockConnection = mock(MockConnection.class);
        when(mock.create()).thenReturn(plcMockConnection);
        when(plcMockConnection.readRequestBuilder()).thenReturn(Mockito.mock(PlcReadRequest.Builder.class));

        CachedDriverManager driverManager = new CachedDriverManager("", mock);

        // Get Connmection
        PlcConnection connection = driverManager.getConnection("");
        // Close the Connection
        PlcReadRequest.Builder builder = connection.readRequestBuilder();
        PlcReadRequest request = builder.addItem("", "").build();
        connection.close();
        assertThrows(IllegalStateException.class, () -> request.execute());
    }

    @Test
    void multipleRequests_allbutfirstFail() throws PlcException {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        MockConnection plcMockConnection = mock(MockConnection.class);
        when(mock.create()).thenReturn(plcMockConnection);
        when(plcMockConnection.readRequestBuilder()).thenReturn(Mockito.mock(PlcReadRequest.Builder.class));

        CachedDriverManager driverManager = new CachedDriverManager("", mock);

        // Get Connmection
        PlcConnection connection = driverManager.getConnection("");

        // Try to get another one -> should failt
        assertThrows(PlcConnectionException.class, () -> driverManager.getConnection(""));
    }

    @Test
    void initialRequests_doesNotFail_dueToQueue() throws PlcConnectionException {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        MockConnection plcMockConnection = mock(MockConnection.class);
        when(mock.create()).thenReturn(plcMockConnection);

        CachedDriverManager driverManager = new CachedDriverManager("", mock, 3_000);

        PlcConnection connection = driverManager.getConnection("");

        assertNotNull(connection);
    }

    @Test
    @Disabled
    void twoRequests_firstTakesLong_secondsTimesOut() throws PlcConnectionException, InterruptedException, ExecutionException, TimeoutException {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        MockConnection plcMockConnection = mock(MockConnection.class);
        when(mock.create()).thenReturn(plcMockConnection);
        when(plcMockConnection.readRequestBuilder()).thenReturn(Mockito.mock(PlcReadRequest.Builder.class));

        CachedDriverManager driverManager = new CachedDriverManager("", mock, 5_000);

        CompletableFuture<PlcConnection> future1 = CompletableFuture.supplyAsync(() -> {
            try {
                return driverManager.getConnection("");
            } catch (PlcConnectionException e) {
                throw new RuntimeException();
            }
        });

        CompletableFuture<PlcConnection> future2 = CompletableFuture.supplyAsync(() -> {
            try {
                return driverManager.getConnection("");
            } catch (PlcConnectionException e) {
                throw new RuntimeException();
            }
        });

        PlcConnection conn1 = future1.get(1, TimeUnit.SECONDS);
        assertThrows(Exception.class, () -> future2.get());
    }

    @Test
    @Disabled
    void twoRequests_firstIsFast_secondWorksAlso() throws Exception {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        MockConnection plcMockConnection = mock(MockConnection.class);
        when(mock.create()).thenReturn(plcMockConnection);
        when(plcMockConnection.readRequestBuilder()).thenReturn(Mockito.mock(PlcReadRequest.Builder.class));

        CachedDriverManager driverManager = new CachedDriverManager("", mock, 5_000);

        CompletableFuture<PlcConnection> future1 = CompletableFuture.supplyAsync(() -> {
            try {
                return driverManager.getConnection("");
            } catch (PlcConnectionException e) {
                throw new RuntimeException(e);
            }
        });

        CompletableFuture<PlcConnection> future2 = CompletableFuture.supplyAsync(() -> {
            try {
                return driverManager.getConnection("");
            } catch (PlcConnectionException e) {
                throw new RuntimeException(e);
            }
        });

        // Get and directly close the first one
        assertDoesNotThrow(() -> {
            future1.get(1, TimeUnit.SECONDS).close();
            future2.get(1, TimeUnit.SECONDS).close();
        });
    }

    @Test
    void killBorrowedConnectionWhenRunningLong() throws PlcConnectionException, InterruptedException {
        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        MockConnection plcMockConnection = mock(MockConnection.class);
        when(mock.create()).thenReturn(plcMockConnection);
        when(plcMockConnection.readRequestBuilder()).thenReturn(Mockito.mock(PlcReadRequest.Builder.class));

        CachedDriverManager driverManager = new CachedDriverManager("", mock);

        // Get Connmection
        PlcConnection connection = driverManager.getConnection("");

        // This should work
        connection.readRequestBuilder();

        TimeUnit.SECONDS.sleep(6);

        // If we wait to long, the connection should become invalid
        assertThrows(IllegalStateException.class, () -> connection.readRequestBuilder());

        // And the Pool should once again have a connection
        assertThat(driverManager.getState())
            .isEqualTo(CachedDriverManager.ConnectionState.DISCONNECTED);
    }
}