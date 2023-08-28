/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.iec608705104.readwrite.protocol;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.iec608705104.readwrite.*;
import org.apache.plc4x.java.iec608705104.readwrite.configuration.Iec608705014Configuration;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.messages.PlcBrowser;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;

import java.time.Duration;
import java.util.Collection;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

public class Iec608705104Protocol extends Plc4xProtocolBase<APDU> implements HasConfiguration<Iec608705014Configuration>, PlcSubscriber, PlcBrowser {

    private Iec608705014Configuration configuration;
    private final RequestTransactionManager tm;

    public Iec608705104Protocol() {
        // We're starting with allowing only one message in-flight.
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void setConfiguration(Iec608705014Configuration configuration) {
        this.configuration = configuration;
    }

    @Override
    public void onConnect(ConversationContext<APDU> context) {
        // First we exchange a test-frame
        APDUUFormatTestFrameActivation testFrameActivation = new APDUUFormatTestFrameActivation(0x43);
        RequestTransactionManager.RequestTransaction testFrameTx = tm.startRequest();
        testFrameTx.submit(() -> context.sendRequest(testFrameActivation)
            .expectResponse(APDU.class, Duration.ofMillis(configuration.getTimeoutRequest()))
            .onTimeout(e -> context.getChannel().pipeline().fireExceptionCaught(e))
            .onError((p, e) -> context.getChannel().pipeline().fireExceptionCaught(e))
            .check(apdu -> apdu instanceof APDUUFormatTestFrameConfirmation)
            .unwrap(apdu -> (APDUUFormatTestFrameConfirmation) apdu)
            .handle(testFrameResponse -> {
                testFrameTx.endRequest();

                // Next send the start-data-transfer packet.
                APDUUFormatStartDataTransferActivation startDataTransferActivation = new APDUUFormatStartDataTransferActivation(0x07);
                RequestTransactionManager.RequestTransaction startDataTransferTx = tm.startRequest();
                startDataTransferTx.submit(() -> context.sendRequest(startDataTransferActivation)
                    .expectResponse(APDU.class, Duration.ofMillis(configuration.getTimeoutRequest()))
                    .onTimeout(e -> context.getChannel().pipeline().fireExceptionCaught(e))
                    .onError((p, e) -> context.getChannel().pipeline().fireExceptionCaught(e))
                    .check(apdu -> apdu instanceof APDUUFormatStartDataTransferConfirmation)
                    .unwrap(apdu -> (APDUUFormatStartDataTransferConfirmation) apdu)
                    .handle(startDataTransferResponse -> {
                        startDataTransferTx.endRequest();
                        context.fireConnected();
                    }));
            }));
    }

    @Override
    public void close(ConversationContext<APDU> context) {

    }


    @Override
    protected void decode(ConversationContext<APDU> context, APDU msg) throws Exception {
        if (msg instanceof APDUUFormatTestFrameActivation) {
            APDUUFormatTestFrameConfirmation testFrameConfirmation = new APDUUFormatTestFrameConfirmation(0x83);
            context.sendToWire(testFrameConfirmation);
        }
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        return null;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {

    }

}
