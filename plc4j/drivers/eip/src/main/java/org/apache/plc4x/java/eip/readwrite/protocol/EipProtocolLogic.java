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

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcInteger;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.eip.readwrite.*;
import org.apache.plc4x.java.eip.readwrite.configuration.EIPConfiguration;
import org.apache.plc4x.java.eip.readwrite.field.EipField;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.InternalPlcReadRequest;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.UnsupportedEncodingException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.logging.Level;

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
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        String field = readRequest.getFieldNames().iterator().next();
        //TODO adapt mpec to everything is generated in CipReadRequest
        // and does not have to be created seperately
        //  Handle multiple fields
        EipField plcField = (EipField) readRequest.getField(field);
        String tag = plcField.getTag();

        //We need the size of the request in words (0x91, tagLength, ... tag + possible pad)
        // Taking half to get word size
        byte requestPathSize = (byte)((tag.length()+2+(tag.length()%2)) / 2);
        CipReadRequest req = new CipReadRequest(requestPathSize,toAnsi(tag),1,(byte)1,(byte)4);

        return toPlcReadResponse((InternalPlcReadRequest)readRequest, readInternal(req));
    }

    private byte[] toAnsi(String tag){
        int arrayIndex = 0;
        boolean isArray = false;
        if(tag.contains("[")){
            isArray = true;
            String index = tag.substring(tag.indexOf("[")+1, tag.indexOf("]"));
            arrayIndex = Integer.parseInt(index);
            tag = tag.substring(0, tag.indexOf("["));
        }

        boolean isPadded = tag.length() % 2 != 0;
        int dataSegLength = 2 + tag.length() + (isPadded==true? 1 : 0) + (isArray==true? 2 : 0);

        ByteBuffer buffer = ByteBuffer.allocate(dataSegLength).order(ByteOrder.LITTLE_ENDIAN);

        buffer.put((byte)0x91);
        buffer.put((byte)tag.length());
        byte[] tagBytes = null;
        try {
            tagBytes = tag.getBytes("US-ASCII");
        } catch (UnsupportedEncodingException ex) {
            java.util.logging.Logger.getLogger(EipProtocolLogic.class.getName()).log(Level.SEVERE, null, ex);
        }

        buffer.put(tagBytes);
        buffer.position(2 +tagBytes.length);


        if(isPadded)
        {
            buffer.put((byte)0x00);
        }

        if(isArray){
            buffer.put((byte)0x28);
            buffer.put((byte)arrayIndex);
        }
        return buffer.array();
    }

    private CompletableFuture<PlcReadResponse> toPlcReadResponse(InternalPlcReadRequest readRequest, CompletableFuture<CipReadResponse> response) {
        return response
            .thenApply(p -> {
                try {
                    return ((PlcReadResponse) decodeReadResponse(p, readRequest));
                } catch (PlcProtocolException e) {
                    throw new PlcRuntimeException("Unable to decode Response", e);
                }
            });
    }
    private CompletableFuture<CipReadResponse> readInternal(CipReadRequest request) {
        //TODO check if this is right
        CompletableFuture<CipReadResponse> future = new CompletableFuture<>();
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        CipUnconnectedRequest exchange = new CipUnconnectedRequest(request);
        CipRRData rrData = new CipRRData(sessionHandle,0L,emptySenderContext,0L,exchange);
        transaction.submit(() -> context.sendRequest(rrData)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check( p -> p  instanceof CipRRData)
            .check( p -> p.getSessionHandle()==sessionHandle)
            .unwrap(p ->(CipRRData) p)
            .unwrap(p -> p.getExchange()).check(p -> p instanceof CipReadResponse)
            .unwrap(p -> (CipReadResponse) p)
            .handle(p -> {
                future.complete(p);
                // Finish the request-transaction.
                transaction.endRequest();
            }));
        return future;

    }
    private PlcResponse decodeReadResponse(CipReadResponse p, InternalPlcReadRequest readRequest) throws PlcProtocolException {
        //TODO Check if this is right
        Map<String, Pair<PlcResponseCode, PlcValue>> values = new HashMap<>();
        // only 1 field
        String fieldName = readRequest.getFieldNames().iterator().next();
        EipField field  = (EipField)readRequest.getField(fieldName);
        PlcResponseCode code;
        if(p.getStatus()==0){
            code = PlcResponseCode.OK;
        }
        else{
            code = PlcResponseCode.INTERNAL_ERROR;
        }
        PlcValue plcValue = null;
        ByteBuf data = Unpooled.wrappedBuffer(p.getData());
        if (code == PlcResponseCode.OK) {
            plcValue = parsePlcValue(field, data);
        }
        Pair<PlcResponseCode,PlcValue> result = new ImmutablePair<>(code,plcValue);
        values.put(fieldName,result);
        return new DefaultPlcReadResponse(readRequest,values);
    }

    private PlcValue parsePlcValue(EipField field, ByteBuf data) {
        int dataType = data.getShort(0);
        switch( dataType){
            case 0xC4: return new PlcInteger(data.getInt(0));
            default:
                return null;
        }
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
