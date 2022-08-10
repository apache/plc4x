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
package org.apache.plc4x.java.utils.connectionpool;

import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.commons.pool2.impl.GenericKeyedObjectPool;
import org.apache.commons.pool2.impl.GenericKeyedObjectPoolConfig;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.PlcDriver;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Answers;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.stream.IntStream;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class PooledPlcDriverManagerTest implements WithAssertions {

    private static Logger LOGGER = LoggerFactory.getLogger(PooledPlcDriverManagerTest.class);

    private PooledPlcDriverManager SUT = new PooledPlcDriverManager(pooledPlcConnectionFactory -> {
        GenericKeyedObjectPoolConfig<PlcConnection> config = new GenericKeyedObjectPoolConfig<>();
        config.setMinIdlePerKey(1);
        config.setTestOnBorrow(true);
        config.setTestOnReturn(true);
        return new GenericKeyedObjectPool<>(pooledPlcConnectionFactory, config);
    });

    @Mock(answer = Answers.RETURNS_DEEP_STUBS)
    PlcDriver plcDriver;

    private ExecutorService executorService;

    @SuppressWarnings("unchecked")
    @BeforeEach
    void setUp() throws Exception {
        Map<String, PlcDriver> driverMap = (Map) FieldUtils.getField(PooledPlcDriverManager.class, "driverMap", true).get(SUT);
        driverMap.put("dummydummy", plcDriver);
        executorService = Executors.newFixedThreadPool(100);

        assertThat(SUT.getStatistics()).containsOnly(
            entry("pools.count", 0),
            entry("numActive", 0),
            entry("numIdle", 0)
        );
    }

    @AfterEach
    void tearDown() {
        executorService.shutdown();
    }

    @Test
    void getConnection() throws Exception {
        when(plcDriver.getConnection(anyString())).then(invocationOnMock -> new DummyPlcConnection(invocationOnMock.getArgument(0)));

        LinkedList<Callable<PlcConnection>> callables = new LinkedList<>();

        // This: should result in one open connection
        IntStream.range(0, 8).forEach(i -> callables.add(() -> {
            try {
                return SUT.getConnection("dummydummy:single/socket1/socket2?fancyOption=true");
            } catch (PlcConnectionException e) {
                throw new RuntimeException(e);
            }
        }));

        // This should result in five open connections
        IntStream.range(0, 5).forEach(i -> callables.add(() -> {
            try {
                return SUT.getConnection("dummydummy:multi-" + i + "/socket1/socket2?fancyOption=true");
            } catch (PlcConnectionException e) {
                throw new RuntimeException(e);
            }
        }));

        List<Future<PlcConnection>> futures = executorService.invokeAll(callables);

        // Wait for existing connections
        futures.forEach(plcConnectionFuture1 -> {
            try {
                plcConnectionFuture1.get();
            } catch (InterruptedException | ExecutionException e) {
                throw new RuntimeException(e);
            }
        });
        LOGGER.info("Statistics after execution {}", SUT.getStatistics());

        // As we have a pool size of 8 we should have only 8 + 5 calls for the separate pools
        verify(plcDriver, times(13)).getConnection(anyString());

        assertThat(SUT.getStatistics()).contains(
            entry("PoolKey{url='dummydummy:single/socket1/socket2?fancyOption=true'}.numActive", 8)
        );

        futures.forEach(plcConnectionFuture -> {
            try {
                plcConnectionFuture.get().close();
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        });

        assertThat(SUT.getStatistics()).contains(
            entry("PoolKey{url='dummydummy:single/socket1/socket2?fancyOption=true'}.numActive", 0)
        );
    }

    @Test
    void getConnectionWithAuth() throws Exception {
        when(plcDriver.getConnection(anyString(), any())).then(invocationOnMock -> new DummyPlcConnection(invocationOnMock.getArgument(0), invocationOnMock.getArgument(1)));

        LinkedList<Callable<PlcConnection>> callables = new LinkedList<>();

        // This: should result in one open connection
        IntStream.range(0, 8).forEach(i -> callables.add(() -> {
            try {
                return SUT.getConnection("dummydummy:single/socket1/socket2?fancyOption=true", new PlcUsernamePasswordAuthentication("user", "passwordp954368564098ß"));
            } catch (PlcConnectionException e) {
                throw new RuntimeException(e);
            }
        }));

        // This should result in five open connections
        IntStream.range(0, 5).forEach(i -> callables.add(() -> {
            try {
                return SUT.getConnection("dummydummy:multi-" + i + "/socket1/socket2?fancyOption=true", new PlcUsernamePasswordAuthentication("user", "passwordp954368564098ß"));
            } catch (PlcConnectionException e) {
                throw new RuntimeException(e);
            }
        }));

        List<Future<PlcConnection>> futures = executorService.invokeAll(callables);

        futures.forEach(plcConnectionFuture1 -> {
            try {
                plcConnectionFuture1.get();
            } catch (InterruptedException | ExecutionException e) {
                throw new RuntimeException(e);
            }
        });

        LOGGER.info("Statistics after execution {}", SUT.getStatistics());

        // As we have a pool size of 8 we should have only 8 + 5 calls for the separate pools
        verify(plcDriver, times(13)).getConnection(anyString(), any());

        assertThat(SUT.getStatistics()).contains(
            entry("PoolKey{url='dummydummy:single/socket1/socket2?fancyOption=true', plcAuthentication=PlcUsernamePasswordAuthentication{username='user', password='*****************'}}.numActive", 8)
        );

        futures.forEach(plcConnectionFuture -> {
            try {
                plcConnectionFuture.get().close();
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        });

        assertThat(SUT.getStatistics()).contains(
            entry("PoolKey{url='dummydummy:single/socket1/socket2?fancyOption=true', plcAuthentication=PlcUsernamePasswordAuthentication{username='user', password='*****************'}}.numActive", 0)
        );
    }

    @Test
    void connectionInvalidation() throws Exception {
        when(plcDriver.getConnection(anyString())).then(invocationOnMock -> new DummyPlcConnection(invocationOnMock.getArgument(0)));

        PlcConnection connection = SUT.getConnection("dummydummy:single/socket1/socket2?fancyOption=true");
        assertThat(connection.isConnected()).isEqualTo(true);
        assertThat(connection.getMetadata().canRead()).isEqualTo(false);
        assertThat(connection.getMetadata().canWrite()).isEqualTo(false);
        assertThat(connection.getMetadata().canSubscribe()).isEqualTo(false);

        connection.close();
        assertThatThrownBy(connection::connect).isInstanceOf(IllegalStateException.class).hasMessage("Proxy not valid anymore");
        assertThatThrownBy(connection::isConnected).isInstanceOf(IllegalStateException.class).hasMessage("Proxy not valid anymore");
        assertThatThrownBy(connection::close).isInstanceOf(IllegalStateException.class).hasMessage("Proxy not valid anymore");
        assertThatThrownBy(connection::getMetadata).isInstanceOf(IllegalStateException.class).hasMessage("Proxy not valid anymore");
        assertThatThrownBy(connection::readRequestBuilder).isInstanceOf(IllegalStateException.class).hasMessage("Proxy not valid anymore");
        assertThatThrownBy(connection::writeRequestBuilder).isInstanceOf(IllegalStateException.class).hasMessage("Proxy not valid anymore");
        assertThatThrownBy(connection::subscriptionRequestBuilder).isInstanceOf(IllegalStateException.class).hasMessage("Proxy not valid anymore");
        assertThatThrownBy(connection::unsubscriptionRequestBuilder).isInstanceOf(IllegalStateException.class).hasMessage("Proxy not valid anymore");
    }

    @Test
    void cleanupOfBrokenConnections() throws Exception {
        AtomicBoolean failNow = new AtomicBoolean(false);
        when(plcDriver.getConnection(anyString())).then(invocationOnMock -> {
            DummyPlcConnection dummyPlcConnection = spy(new DummyPlcConnection(invocationOnMock.getArgument(0)));
            // we fake an connection which breaks at this call
            doAnswer(invocation -> {
                if (failNow.get()) {
                    throw new PlcConnectionException("blub");
                }
                return invocation.callRealMethod();
            }).when(dummyPlcConnection).connect();
            return dummyPlcConnection;
        });

        assertThat(SUT.getStatistics()).containsOnly(
            entry("pools.count", 0),
            entry("numActive", 0),
            entry("numIdle", 0)
        );
        PlcConnection connection = SUT.getConnection("dummydummy:breakIt");
        assertThat(SUT.getStatistics()).containsOnly(
            entry("pools.count", 1),
            entry("numActive", 1),
            entry("numIdle", 0),
            entry("PoolKey{url='dummydummy:breakIt'}.numActive", 1)
        );
        failNow.set(true);
        try {
            connection.connect();
            fail("This should throw an exception");
        } catch (Exception e) {
            // TODO: currently UndeclaredThrowableException is the top one which should be InvocationTargetException
            //assertThat(e).isInstanceOf(InvocationTargetException.class);
            assertThat(e).hasRootCauseInstanceOf(PlcConnectionException.class);
        }
        // Faulty connection should have been discarded
        assertThat(SUT.getStatistics()).containsOnly(
            entry("pools.count", 0),
            entry("numActive", 0),
            entry("numIdle", 0)
        );
    }

    @Test
    public void testOtherConstructors() {
        assertThat(new PooledPlcDriverManager()).isNotNull();
        assertThat(new PooledPlcDriverManager(new PoolKeyFactory())).isNotNull();
        assertThat(new PooledPlcDriverManager(PooledPlcDriverManager.class.getClassLoader())).isNotNull();
        assertThat(new PooledPlcDriverManager(
            PooledPlcDriverManager.class.getClassLoader(), new PoolKeyFactory())).isNotNull();
    }

    class DummyPlcConnection implements PlcConnection, PlcConnectionMetadata {

        private final String url;

        private final PlcAuthentication plcAuthentication;

        boolean connected = false;

        public DummyPlcConnection(String url) {
            this(url, null);
        }

        public DummyPlcConnection(String url, PlcAuthentication plcAuthentication) {
            this.url = url;
            this.plcAuthentication = plcAuthentication;
        }

        @Override
        public void connect() {
            connected = true;
        }

        @Override
        public CompletableFuture<Void> ping() {
            CompletableFuture<Void> future = new CompletableFuture<>();
            future.completeExceptionally(new UnsupportedOperationException());
            return future;
        }

        @Override
        public boolean isConnected() {
            return connected;
        }

        @Override
        public PlcConnectionMetadata getMetadata() {
            return this;
        }

        @Override
        public boolean canRead() {
            return false;
        }

        @Override
        public boolean canWrite() {
            return false;
        }

        @Override
        public boolean canSubscribe() {
            return false;
        }

        @Override
        public boolean canBrowse() {
            return false;
        }

        @Override
        public void close() {
            connected = false;
        }

        @Override
        public PlcReadRequest.Builder readRequestBuilder() {
            throw new PlcUnsupportedOperationException("The connection does not support reading");
        }

        @Override
        public PlcWriteRequest.Builder writeRequestBuilder() {
            throw new PlcUnsupportedOperationException("The connection does not support writing");
        }

        @Override
        public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
            throw new PlcUnsupportedOperationException("The connection does not support subscription");
        }

        @Override
        public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
            throw new PlcUnsupportedOperationException("The connection does not support subscription");
        }

        @Override
        public PlcBrowseRequest.Builder browseRequestBuilder() {
            throw new PlcUnsupportedOperationException("The connection does not support browsing");
        }

        @Override
        public String toString() {
            return "DummyPlcConnection{" +
                "url='" + url + '\'' +
                ", plcAuthentication=" + plcAuthentication +
                ", connected=" + connected +
                '}';
        }
    }
}