/*
 Licensed to the Apache Software Foundation (ASF) under one or more
 contributor license agreements.  See the NOTICE file distributed with
 this work for additional information regarding copyright ownership.
 The ASF licenses this file to You under the Apache License, Version 2.0
 (the "License"); you may not use this file except in compliance with
 the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 */

package org.apache.plc4x.java.utils.connectionpool;

import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.commons.pool2.impl.GenericObjectPool;
import org.apache.commons.pool2.impl.GenericObjectPoolConfig;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.spi.PlcDriver;
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
import java.util.stream.IntStream;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class PooledPlcDriverManagerTest implements WithAssertions {

    private static Logger LOGGER = LoggerFactory.getLogger(PooledPlcDriverManagerTest.class);

    private PooledPlcDriverManager SUT = new PooledPlcDriverManager(pooledPlcConnectionFactory -> {
        GenericObjectPoolConfig<PlcConnection> plcConnectionGenericObjectPoolConfig = new GenericObjectPoolConfig<>();
        plcConnectionGenericObjectPoolConfig.setMinIdle(1);
        return new GenericObjectPool<>(pooledPlcConnectionFactory, plcConnectionGenericObjectPoolConfig);
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

        assertThat(SUT.getStatistics()).isEmpty();
    }

    @AfterEach
    void tearDown() {
        executorService.shutdown();
    }

    @Test
    void getConnection() throws Exception {
        when(plcDriver.connect(anyString())).then(invocationOnMock -> new DummyPlcConnection(invocationOnMock.getArgument(0)));

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
        verify(plcDriver, times(13)).connect(anyString());

        assertThat(SUT.getStatistics()).contains(
            entry("dummydummy:single/socket1/socket2?fancyOption=true.numActive", 8),
            entry("dummydummy:single/socket1/socket2?fancyOption=true.numIdle", 0)
        );

        futures.forEach(plcConnectionFuture -> {
            try {
                plcConnectionFuture.get().close();
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        });

        assertThat(SUT.getStatistics()).contains(
            entry("dummydummy:single/socket1/socket2?fancyOption=true.numActive", 0),
            entry("dummydummy:single/socket1/socket2?fancyOption=true.numIdle", 8)
        );
    }

    @Test
    void getConnectionWithAuth() throws Exception {
        when(plcDriver.connect(anyString(), any())).then(invocationOnMock -> new DummyPlcConnection(invocationOnMock.getArgument(0), invocationOnMock.getArgument(1)));

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
                return SUT.getConnection("dummydummy:single-" + i + "/socket1/socket2?fancyOption=true", new PlcUsernamePasswordAuthentication("user", "passwordp954368564098ß"));
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
        verify(plcDriver, times(13)).connect(anyString(), any());

        assertThat(SUT.getStatistics()).contains(
            entry("dummydummy:single/socket1/socket2?fancyOption=true/PlcUsernamePasswordAuthentication{username='user', password='*****************'}.numActive", 8),
            entry("dummydummy:single/socket1/socket2?fancyOption=true/PlcUsernamePasswordAuthentication{username='user', password='*****************'}.numIdle", 0)
        );

        futures.forEach(plcConnectionFuture -> {
            try {
                plcConnectionFuture.get().close();
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        });

        assertThat(SUT.getStatistics()).contains(
            entry("dummydummy:single/socket1/socket2?fancyOption=true/PlcUsernamePasswordAuthentication{username='user', password='*****************'}.numActive", 0),
            entry("dummydummy:single/socket1/socket2?fancyOption=true/PlcUsernamePasswordAuthentication{username='user', password='*****************'}.numIdle", 8)
        );
    }

    @Test
    void connectionInvalidation() throws Exception {
        when(plcDriver.connect(anyString())).then(invocationOnMock -> new DummyPlcConnection(invocationOnMock.getArgument(0)));

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
        public void connect() throws PlcConnectionException {
            connected = true;
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
        public void close() throws Exception {
            throw new UnsupportedOperationException("this should never be called due to pool");
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
        public String toString() {
            return "DummyPlcConnection{" +
                "url='" + url + '\'' +
                ", plcAuthentication=" + plcAuthentication +
                ", connected=" + connected +
                '}';
        }
    }
}