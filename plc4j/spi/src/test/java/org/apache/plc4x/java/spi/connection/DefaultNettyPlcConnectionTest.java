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

package org.apache.plc4x.java.spi.connection;

import static java.util.concurrent.ForkJoinPool.commonPool;
import static org.junit.jupiter.api.Assertions.assertEquals;

import java.util.function.Function;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xNettyWrapper;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.netty.NettyHashTimerTimeoutManager;
import org.junit.jupiter.api.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class DefaultNettyPlcConnectionTest {

    private final Logger logger = LoggerFactory.getLogger(DefaultNettyPlcConnectionTest.class);

    final GateKeeper discovery = new GateKeeper("discovery");
    final GateKeeper connect = new GateKeeper("connect");
    final GateKeeper disconnect = new GateKeeper("disconnect");
    final GateKeeper close = new GateKeeper("close");

    @Test
    void checkInitializationSequence() throws Exception {
        ChannelFactory channelFactory = new TestChannelFactory();

        ProtocolStackConfigurer<Message> stackConfigurer = (configuration, pipeline, authentication, passive, listeners) -> {
            TestProtocolBase base = new TestProtocolBase(discovery, connect, disconnect, close);
            Plc4xNettyWrapper<Message> context = new Plc4xNettyWrapper<>(new NettyHashTimerTimeoutManager(), pipeline, passive, base, authentication, Message.class);
            pipeline.addLast(context);
            return base;
        };

        DefaultNettyPlcConnection connection = new PlcConnectionFactory().withDiscovery().create(channelFactory, stackConfigurer);

        logger.info("Warming up");
        verifyEntry(false, false, false, false);
        verifyExits(false, false, false, false);

        commonPool().submit(() -> {
            try {
                logger.info("Activating connection");
                connection.connect();
            } catch (PlcConnectionException e) {
                throw new RuntimeException(e);
            }
        });

        discovery.permitEntry();
        logger.info("Verify discovery event gets fired");
        verifyEntry(true, false, false, false);
        discovery.awaitExit();
        verifyExits(true, false, false, false);

        logger.info("Verify discovery connection disconnect event get fired");
        disconnect.permitEntry();
        verifyEntry(true, false, true, false);
        disconnect.awaitExit();
        verifyExits(true, false, true, false);

        logger.info("Verify discovery connection close event get fired");
        close.permitEntry();
        verifyEntry(true, false, true, true);
        close.awaitExit();
        verifyExits(true, false, true, true);

        // end of discovery phase

        discovery.reset();
        connect.reset();
        disconnect.reset();
        close.reset();

        connect.permitEntry();
        logger.info("Verify connection handshake");
        verifyEntry(false, true, false, false);
        connect.awaitExit();
        verifyExits(false, true, false, false);

        logger.info("Close connection");
        commonPool().submit(() -> {
            try {
                logger.info("Closing connection");
                connection.close();
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        });

        disconnect.permitEntry();
        verifyEntry(false, true, true, false);
        disconnect.awaitExit();
        verifyEntry(false, true, true, false);

        logger.info("Verify connection termination");
        close.permitEntry();
        verifyEntry(false, true, true, true);
        close.awaitExit();
        verifyExits(false, true, true, true);

        logger.info("Connection lifecycle sequence has been confirmed");
    }

    private void verifyEntry(boolean discovered, boolean connected, boolean disconnected, boolean closed) {
        verifyGates(GateKeeper::entered, discovered, connected, disconnected, closed, discovery, connect, disconnect, close);
    }

    private void verifyExits(boolean discovered, boolean connected, boolean disconnected, boolean closed) {
        verifyGates(GateKeeper::exited, discovered, connected, disconnected, closed, discovery, connect, disconnect, close);
    }

    private static void verifyGates(Function<GateKeeper, Boolean> mapper, boolean discovered, boolean connected, boolean disconnected, boolean closed,
        GateKeeper discovery, GateKeeper connect, GateKeeper disconnect, GateKeeper close) {

        boolean actualDiscovered = mapper.apply(discovery);
        boolean actualConnected = mapper.apply(connect);
        boolean actualDisconnected = mapper.apply(disconnect);
        boolean actualClosed = mapper.apply(close);

        assertEquals(
            discovered + "," + connected + "," +  disconnected + "," +  closed,
            actualDiscovered + "," + actualConnected + "," + actualDisconnected + "," + actualClosed,
            "Expectation for state flags (discover, connect, disconnect, close) failed"
        );
    }

    static class TestProtocolBase extends Plc4xProtocolBase<Message> {

        private final Logger logger = LoggerFactory.getLogger(TestProtocolBase.class);;
        private final GateKeeper discover;
        private final GateKeeper connect;
        private final GateKeeper close;
        private final GateKeeper disconnect;

        public TestProtocolBase(GateKeeper discover, GateKeeper connect, GateKeeper disconnect, GateKeeper close) {
            this.discover = discover;
            this.connect = connect;
            this.close = close;
            this.disconnect = disconnect;
        }

        @Override
        public PlcTagHandler getTagHandler() {
            return null;
        }

        @Override
        public void onDiscover(ConversationContext<Message> context) {
            logger.info("On Discover");
            awaitIn(discover);
            context.fireDiscovered(null);
            discover.reportExit();
        }


        @Override
        public void onConnect(ConversationContext<Message> context) {
            logger.info("On Connect");
            awaitIn(connect);
            super.onConnect(context);
            context.fireConnected();
            connect.reportExit();
        }

        @Override
        public void onDisconnect(ConversationContext<Message> context) {
            logger.info("On Disconnect");
            awaitIn(disconnect);
            super.onDisconnect(context);
            context.fireDisconnected();
            disconnect.reportExit();
        }

        @Override
        public void close(ConversationContext<Message> context) {
            logger.info("On Close");
            awaitIn(close);
            close.reportExit();
        }

        private void awaitIn(GateKeeper signal) {
            try {
                if (!signal.awaitIn()) {
                    throw new RuntimeException("Await for " + signal.gate() + " lock failed");
                }
            } catch (InterruptedException e) {
                logger.error("Failed to await for a signal " + signal.gate());
                throw new RuntimeException(e);
            }
        }
    }

}