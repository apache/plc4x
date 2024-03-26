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

    @Test
    void checkInitializationSequence() throws Exception {
        ChannelFactory channelFactory = new TestChannelFactory();

        final GateKeeper discovery = new GateKeeper("discovery");
        final GateKeeper connect = new GateKeeper("connect");
        final GateKeeper disconnect = new GateKeeper("disconnect");
        final GateKeeper close = new GateKeeper("close");

        ProtocolStackConfigurer<Message> stackConfigurer = (configuration, pipeline, authentication, passive, listeners) -> {
            TestProtocolBase base = new TestProtocolBase(discovery, connect, disconnect, close);
            Plc4xNettyWrapper<Message> context = new Plc4xNettyWrapper<>(new NettyHashTimerTimeoutManager(), pipeline, passive, base, authentication, Message.class);
            pipeline.addLast(context);
            return base;
        };

        DefaultNettyPlcConnection connection = new PlcConnectionFactory().withDiscovery().create(channelFactory, stackConfigurer);
        commonPool().submit(() -> {
            try {
                logger.info("Activating connection");
                connection.connect();
            } catch (PlcConnectionException e) {
                throw new RuntimeException(e);
            }
        });

        logger.info("Warming up");
        expect(false, false, false, false, discovery, connect, disconnect, close);
        discovery.permitIn();

        discovery.awaitOut();
        logger.info("Verify discovery phase completion");
        expect(true, false, false, false, discovery, connect, disconnect, close);
        connect.permitIn();

        connect.awaitOut();
        logger.info("Verify connection completion");
        expect(true, true, false, false, discovery, connect, disconnect, close);

        logger.info("Close connection");
        commonPool().submit(new Runnable() {
            @Override
            public void run() {
                try {
                    logger.info("Closing connection");
                    connection.close();
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            }
        });

        disconnect.permitIn();
        expect(true, true, true, false, discovery, connect, disconnect, close);
        disconnect.awaitOut();

        logger.info("Verify connection termination");
        close.permitIn();
        expect(true, true, true, true, discovery, connect, disconnect, close);
        close.awaitOut();

        logger.info("Connection lifecycle sequence has been confirmed");
    }

    private static void expect(boolean discovered, boolean connected, boolean disconnected, boolean closed,
        GateKeeper discovery, GateKeeper connect, GateKeeper disconnect, GateKeeper close) {

        assertEquals(
            discovered + "," + connected + "," +  disconnected + "," +  closed,
            (discovery.entered()) + "," +
            (connect.entered()) + "," +
            (disconnect.entered() + "," +
            (close.entered())),
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
        public void onDiscover(ConversationContext<Message> context) {
            logger.info("On Discover");
            await(discover);
            context.fireDiscovered(null);
            discover.permitOut();
        }


        @Override
        public void onConnect(ConversationContext<Message> context) {
            logger.info("On Connect");
            await(connect);
            super.onConnect(context);
            context.fireConnected();
            connect.permitOut();
        }

        @Override
        public void onDisconnect(ConversationContext<Message> context) {
            logger.info("On Disconnect");
            await(disconnect);
            super.onDisconnect(context);
            context.fireDisconnected();
            disconnect.permitOut();
        }

        @Override
        public void close(ConversationContext<Message> context) {
            logger.info("On Close");
            await(close);
            close.permitOut();
        }

        private void await(GateKeeper signal) {
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