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
package org.apache.plc4x.java.eip.readwrite.protocol;

import org.apache.plc4x.java.eip.readwrite.EipConnectionRequest;
import org.apache.plc4x.java.eip.readwrite.EipDisconnectRequest;
import org.apache.plc4x.java.eip.readwrite.EipPacket;
import org.apache.plc4x.java.eip.readwrite.configuration.EIPConfiguration;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.concurrent.atomic.AtomicInteger;

public class EipProtocolLogic extends Plc4xProtocolBase<EipPacket>implements HasConfiguration<EIPConfiguration> {

    private static final Logger logger = LoggerFactory.getLogger(EipProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private static final short[] emptySenderContext = new short[] {(short) 0x00 ,(short) 0x00 ,(short) 0x00,
        (short) 0x00,(short) 0x00,(short) 0x00, (short) 0x00,(short) 0x00};
    private EIPConfiguration configuration;

    private final AtomicInteger transactionCounterGenerator = new AtomicInteger(10);
    private RequestTransactionManager tm;
    private long sessionHandle;

    @Override
    public void setConfiguration(EIPConfiguration configuration){
        this.configuration = configuration;
        // Set the transaction manager to allow only one message at a time.
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void onConnect(ConversationContext<EipPacket> context) {
        /**Send a ENIP Message with Register Session Code '0x0065',
         * empty Session Handle and Sender Context
         * Then we need to accept the response with the same Code
         * and save the assigned Session Handle
         * PS: Check Status for Success : 0x00000000*/
        logger.info("Sending RegisterSession EIP Package");
        EipConnectionRequest connectionRequest =
            new EipConnectionRequest(0L, 0L, emptySenderContext, 0L);
        context.sendRequest(connectionRequest)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
            .check(p -> p instanceof EipConnectionRequest)
            .unwrap(p -> (EipConnectionRequest) p)
            .handle(EipConnectionRequest -> {
                if(EipConnectionRequest.getStatus()==0L){
                    sessionHandle = EipConnectionRequest.getSessionHandle();
                    logger.trace("Got assigned with Session {}", sessionHandle);
                    // Send an event that connection setup is complete.
                    context.fireConnected();
                }
                else{
                    logger.warn("Got status code [{}]", EipConnectionRequest.getStatus());
                }

            });
    }





    @Override
    public void close(ConversationContext<EipPacket> context) {
        /**Send a ENIP Message with Unregister Session Code '0x0066' */
        logger.info("Sending UnregisterSession EIP Pakcet");
        EipDisconnectRequest disconnectRequest =
            new EipDisconnectRequest(sessionHandle,0L,emptySenderContext,0L);
        context.sendRequest(disconnectRequest); //Unregister gets no response
        logger.trace("Unregistred Session {}", sessionHandle);
    }
}
